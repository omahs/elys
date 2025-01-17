package cli

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	"github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/elys-network/elys/x/incentive/types"
)

var (
	FlagCommission                        = "commission"
	FlagValidatorAddress                  = "validator-address"
	FlagMaxMessagesPerTx                  = "max-msgs"
	DefaultRelativePacketTimeoutTimestamp = uint64((time.Duration(10) * time.Minute).Nanoseconds())
)

const (
	flagPacketTimeoutTimestamp = "packet-timeout-timestamp"
	MaxMessagesPerTxDefault    = 0
	listSeparator              = ","
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(
		CmdWithdrawRewardsCmd(),
		CmdUpdatePoolInfoProposal(),
	)

	// this line is used by starport scaffolding # 1

	return cmd
}

type newGenerateOrBroadcastFunc func(client.Context, *pflag.FlagSet, ...sdk.Msg) error

func newSplitAndApply(
	genOrBroadcastFn newGenerateOrBroadcastFunc, clientCtx client.Context,
	fs *pflag.FlagSet, msgs []sdk.Msg, chunkSize int,
) error {
	if chunkSize == 0 {
		return genOrBroadcastFn(clientCtx, fs, msgs...)
	}

	// split messages into slices of length chunkSize
	totalMessages := len(msgs)
	for i := 0; i < len(msgs); i += chunkSize {

		sliceEnd := i + chunkSize
		if sliceEnd > totalMessages {
			sliceEnd = totalMessages
		}

		msgChunk := msgs[i:sliceEnd]
		if err := genOrBroadcastFn(clientCtx, fs, msgChunk...); err != nil {
			return err
		}
	}

	return nil
}

// CmdWithdrawRewardsCmd returns a CLI command handler for creating a MsgWithdrawDelegatorReward transaction.
func CmdWithdrawRewardsCmd() *cobra.Command {
	bech32PrefixValAddr := sdk.GetConfig().GetBech32ValidatorAddrPrefix()

	cmd := &cobra.Command{
		Use:   "withdraw-rewards",
		Short: "Withdraw rewards from a given delegation address, and optionally withdraw validator commission if the delegation address given is a validator operator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Withdraw rewards from a given delegation address,
and optionally withdraw validator commission if the delegation address given is a validator operator.

Example:
$ %s tx incentive withdraw-rewards --from mykey
$ %s tx incentive withdraw-rewards --from mykey --commission --validator-address %s1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj
`,
				version.AppName, bech32PrefixValAddr, bech32PrefixValAddr,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			argDenom := args[0]
			delAddr := clientCtx.GetFromAddress()
			msgs := []sdk.Msg{types.NewMsgWithdrawRewards(delAddr, argDenom)}

			if commission, _ := cmd.Flags().GetBool(FlagCommission); commission {
				if validatorAddr, _ := cmd.Flags().GetString(FlagValidatorAddress); len(validatorAddr) > 0 {
					valAddr, err := sdk.ValAddressFromBech32(validatorAddr)
					if err != nil {
						return err
					}
					msgs = append(msgs, types.NewMsgWithdrawValidatorCommission(delAddr, valAddr, argDenom))
				}
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msgs...)
		},
	}

	cmd.Flags().Bool(FlagCommission, false, "Withdraw the validator's commission in addition to the rewards")
	cmd.Flags().String(FlagValidatorAddress, "", "Validator's operator address to withdraw commission from")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// CmdUpdatePoolInfoProposal returns a CLI command handler for submitting a UpdatePoolInfo proposal.
func CmdUpdatePoolInfoProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-pool-info [pool-ids] [multipliers]",
		Short: "Submit an update-pool-info proposal",
		Long:  "e.g. update-pool-info 1,2,3,4, 1,1,2,2",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argPoolIds := args[0]
			argMultipliers := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			title, err := cmd.Flags().GetString(cli.FlagTitle)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(cli.FlagDescription)
			if err != nil {
				return err
			}

			poolIds := strings.Split(argPoolIds, ",")
			multipliers := strings.Split(argMultipliers, ",")
			if len(poolIds) < 1 || len(poolIds) != len(multipliers) {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "invalid parameter")
			}

			poolMultipliers := make([]types.PoolMultipliers, 0)
			for i := range poolIds {
				poolId, err := strconv.ParseUint(poolIds[i], 10, 64)
				if err != nil {
					return sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "invalid parameter")
				}

				multiplier, err := sdk.NewDecFromStr(multipliers[i])
				if err != nil {
					return sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "invalid parameter")
				}
				poolMultiplier := types.PoolMultipliers{
					PoolId:     poolId,
					Multiplier: multiplier,
				}

				poolMultipliers = append(poolMultipliers, poolMultiplier)
			}

			content := types.NewProposalUpdatePoolMultipliers(
				title,
				description,
				poolMultipliers,
			)

			from := clientCtx.GetFromAddress()

			depositStr, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return err
			}
			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			msg, err := v1beta1.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(cli.FlagTitle, "", "title of proposal")
	cmd.Flags().String(cli.FlagDescription, "", "description of proposal")
	cmd.Flags().String(cli.FlagDeposit, "", "deposit of proposal")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
