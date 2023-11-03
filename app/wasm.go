package app

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	cosmos_sdk_math "cosmossdk.io/math"
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	query "github.com/cosmos/cosmos-sdk/types/query"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ammkeeper "github.com/elys-network/elys/x/amm/keeper"
	ammtype "github.com/elys-network/elys/x/amm/types"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
	marginkeeper "github.com/elys-network/elys/x/margin/keeper"
	margintypes "github.com/elys-network/elys/x/margin/types"
	oraclekeeper "github.com/elys-network/elys/x/oracle/keeper"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
	paramtypes "github.com/elys-network/elys/x/parameter/types"
)

// AllCapabilities returns all capabilities available with the current wasmvm
// See https://github.com/CosmWasm/cosmwasm/blob/main/docs/CAPABILITIES-BUILT-IN.md
// This functionality is going to be moved upstream: https://github.com/CosmWasm/wasmvm/issues/425
func AllCapabilities() []string {
	return []string{
		"iterator",
		"staking",
		"stargate",
		"cosmwasm_1_1",
		"cosmwasm_1_2",
	}
}

type QueryPlugin struct {
	ammKeeper        *ammkeeper.Keeper
	oracleKeeper     *oraclekeeper.Keeper
	bankKeeper       *bankkeeper.BaseKeeper
	stakingKeeper    *stakingkeeper.Keeper
	commitmentKeeper *commitmentkeeper.Keeper
	marginKeeper     *marginkeeper.Keeper
}

// NewQueryPlugin returns a reference to a new QueryPlugin.
func NewQueryPlugin(
	amm *ammkeeper.Keeper,
	oracle *oraclekeeper.Keeper,
	bank *bankkeeper.BaseKeeper,
	staking *stakingkeeper.Keeper,
	commitment *commitmentkeeper.Keeper,
	margin *marginkeeper.Keeper,
) *QueryPlugin {
	return &QueryPlugin{
		ammKeeper:        amm,
		oracleKeeper:     oracle,
		bankKeeper:       bank,
		stakingKeeper:    staking,
		commitmentKeeper: commitment,
		marginKeeper:     margin,
	}
}

func RegisterCustomPlugins(
	amm *ammkeeper.Keeper,
	oracle *oraclekeeper.Keeper,
	margin *marginkeeper.Keeper,
	bank *bankkeeper.BaseKeeper,
	staking *stakingkeeper.Keeper,
	commitment *commitmentkeeper.Keeper,
) []wasmkeeper.Option {
	wasmQueryPlugin := NewQueryPlugin(amm, oracle, bank, staking, commitment, margin)

	queryPluginOpt := wasmkeeper.WithQueryPlugins(&wasmkeeper.QueryPlugins{
		Custom: CustomQuerier(wasmQueryPlugin),
	})

	messengerDecoratorOpt := wasmkeeper.WithMessageHandlerDecorator(
		CustomMessageDecorator(amm, margin, staking, commitment),
	)
	return []wasm.Option{
		queryPluginOpt,
		messengerDecoratorOpt,
	}
}

// CustomQuerier dispatches custom CosmWasm bindings queries.
func CustomQuerier(qp *QueryPlugin) func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
	return func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
		var contractQuery ElysQuery
		if err := json.Unmarshal(request, &contractQuery); err != nil {
			return nil, errorsmod.Wrap(err, "elys query")
		}

		switch {
		case contractQuery.PriceAll != nil:
			pagination := contractQuery.PriceAll.Pagination

			// Calling the PriceAll function and handling its response
			priceResponse, err := qp.oracleKeeper.PriceAll(ctx, &oracletypes.QueryAllPriceRequest{Pagination: pagination})
			if err != nil {
				return nil, errorsmod.Wrap(err, "failed to get all prices")
			}

			// copy array priceResponse.Price
			price := make([]oracletypes.Price, len(priceResponse.Price))
			copy(price, priceResponse.Price)

			res := AllPriceResponse{
				Price: price,
				Pagination: &query.PageResponse{
					NextKey: priceResponse.Pagination.NextKey,
				},
			}

			// Serializing the response to a JSON byte array
			responseBytes, err := json.Marshal(res)
			if err != nil {
				return nil, errorsmod.Wrap(err, "failed to serialize price response")
			}

			return responseBytes, nil
		case contractQuery.QuerySwapEstimation != nil:
			return nil, wasmvmtypes.UnsupportedRequest{Kind: "QuerySwapEstimation, not implemented yet"}

		case contractQuery.AssetInfo != nil:
			denom := contractQuery.AssetInfo.Denom

			AssetInfoResp, err := qp.oracleKeeper.AssetInfo(ctx, &oracletypes.QueryGetAssetInfoRequest{Denom: denom})
			if err != nil {
				return nil, errorsmod.Wrap(err, "failed to query asset info")
			}

			res := AssetInfoResponse{
				AssetInfo: &AssetInfoType{
					Denom:      AssetInfoResp.AssetInfo.Denom,
					Display:    AssetInfoResp.AssetInfo.Display,
					BandTicker: AssetInfoResp.AssetInfo.BandTicker,
					ElysTicker: AssetInfoResp.AssetInfo.ElysTicker,
					Decimal:    AssetInfoResp.AssetInfo.Decimal,
				},
			}

			responseBytes, err := json.Marshal(res)
			if err != nil {
				return nil, errorsmod.Wrap(err, "failed to serialize asset info response")
			}
			return responseBytes, nil
		case contractQuery.BalanceOfDenom != nil:
			denom := contractQuery.BalanceOfDenom.Denom
			addr := contractQuery.BalanceOfDenom.Address
			address, err := sdk.AccAddressFromBech32(contractQuery.BalanceOfDenom.Address)
			if err != nil {
				return nil, errorsmod.Wrap(err, "invalid address")
			}
			balance := qp.bankKeeper.GetBalance(ctx, address, denom)
			if denom != paramtypes.Elys {
				commitment, found := qp.commitmentKeeper.GetCommitments(ctx, addr)
				if !found {
					balance = sdk.NewCoin(denom, sdk.ZeroInt())
				} else {
					uncommittedToken, found := commitment.GetUncommittedTokensForDenom(denom)
					if !found {
						return nil, errorsmod.Wrap(nil, "invalid denom")
					}

					balance = sdk.NewCoin(denom, uncommittedToken.Amount)
				}
			}

			res := QueryBalanceResponse{
				Balance: balance,
			}

			responseBytes, err := json.Marshal(res)
			if err != nil {
				return nil, errorsmod.Wrap(err, "failed to get balance response")
			}
			return responseBytes, nil
		default:
			return nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown elys query variant"}
		}
	}
}

type ElysQuery struct {
	PriceAll            *PriceAll                   `json:"price_all,omitempty"`
	QuerySwapEstimation *QuerySwapEstimationRequest `json:"query_swap_estimation,omitempty"`
	AssetInfo           *AssetInfo                  `json:"asset_info,omitempty"`
	BalanceOfDenom      *QueryBalanceRequest        `json:"balance_of_denom,omitempty"`
}

type PriceAll struct {
	Pagination *query.PageRequest `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

type AllPriceResponse struct {
	Price      []oracletypes.Price `protobuf:"bytes,1,rep,name=price,proto3" json:"price"`
	Pagination *query.PageResponse `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

type QuerySwapEstimationRequest struct {
	TokenIn sdk.Coin                    `protobuf:"bytes,2,opt,name=tokenIn,proto3" json:"token_in,omitempty"`
	Routes  []ammtype.SwapAmountInRoute `protobuf:"bytes,1,rep,name=routes,proto3" json:"routes,omitempty"`
}

type QuerySwapEstimationResponse struct {
	SpotPrice sdk.Dec  `protobuf:"bytes,1,opt,name=SpotPrice,proto3" json:"spot_price,omitempty"`
	TokenOut  sdk.Coin `protobuf:"bytes,2,opt,name=tokenOut,proto3" json:"token_out,omitempty"`
}

type AssetInfo struct {
	Denom string `protobuf:"bytes,1,opt,name=Denom,proto3" json:"Denom,omitempty"`
}

type AssetInfoResponse struct {
	AssetInfo *AssetInfoType `protobuf:"bytes,1,opt,name=AssetInfo,proto3" json:"asset_info,omitempty"`
}

type AssetInfoType struct {
	Denom      string `protobuf:"bytes,1,opt,name=denom,proto3" json:"denom,omitempty"`
	Display    string `protobuf:"bytes,2,opt,name=display,proto3" json:"display,omitempty"`
	BandTicker string `protobuf:"bytes,3,opt,name=bandTicker,proto3" json:"band_ticker,omitempty"`
	ElysTicker string `protobuf:"bytes,4,opt,name=elysTicker,proto3" json:"elys_ticker,omitempty"`
	Decimal    uint64 `protobuf:"varint,5,opt,name=decimal,proto3" json:"decimal,omitempty"`
}

type QueryBalanceRequest struct {
	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Denom   string `protobuf:"bytes,2,opt,name=denom,proto3" json:"denom,omitempty"`
}

type QueryBalanceResponse struct {
	Balance sdk.Coin `protobuf:"bytes,1,opt,name=balance,proto3" json:"balance,omitempty"`
}

func CustomMessageDecorator(amm *ammkeeper.Keeper, margin *marginkeeper.Keeper, staking *stakingkeeper.Keeper, commitment *commitmentkeeper.Keeper) func(wasmkeeper.Messenger) wasmkeeper.Messenger {
	return func(old wasmkeeper.Messenger) wasmkeeper.Messenger {
		return &CustomMessenger{
			wrapped:    old,
			amm:        amm,
			margin:     margin,
			staking:    staking,
			commitment: commitment,
		}
	}
}

type CustomMessenger struct {
	wrapped    wasmkeeper.Messenger
	amm        *ammkeeper.Keeper
	margin     *marginkeeper.Keeper
	staking    *stakingkeeper.Keeper
	commitment *commitmentkeeper.Keeper
}

var _ wasmkeeper.Messenger = (*CustomMessenger)(nil)

func (m *CustomMessenger) DispatchMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmvmtypes.CosmosMsg) ([]sdk.Event, [][]byte, error) {
	if msg.Custom != nil {
		// only handle the happy path where this is really creating / minting / swapping ...
		// leave everything else for the wrapped version
		var contractMsg ElysMsg
		if err := json.Unmarshal(msg.Custom, &contractMsg); err != nil {
			return nil, nil, errorsmod.Wrap(err, "elys msg")
		}
		if contractMsg.MsgSwapExactAmountIn != nil {
			return m.msgSwapExactAmountIn(ctx, contractAddr, contractMsg.MsgSwapExactAmountIn)
		}
		if contractMsg.MsgClose != nil {
			return m.msgClose(ctx, contractAddr, contractMsg.MsgClose)
		}
		if contractMsg.MsgOpen != nil {
			return m.msgOpen(ctx, contractAddr, contractMsg.MsgOpen)
		}
		if contractMsg.MsgStake != nil {
			return m.msgStake(ctx, contractAddr, contractMsg.MsgStake)
		}
		if contractMsg.MsgUnstake != nil {
			return m.msgUnstake(ctx, contractAddr, contractMsg.MsgUnstake)
		}
	}
	return m.wrapped.DispatchMsg(ctx, contractAddr, contractIBCPortID, msg)
}

func (m *CustomMessenger) msgSwapExactAmountIn(ctx sdk.Context, contractAddr sdk.AccAddress, msgSwapExactAmountIn *MsgSwapExactAmountIn) ([]sdk.Event, [][]byte, error) {
	res, err := PerformMsgSwapExactAmountIn(m.amm, ctx, contractAddr, msgSwapExactAmountIn)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "perform swap")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize swap response")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}

func PerformMsgSwapExactAmountIn(f *ammkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, msgSwapExactAmountIn *MsgSwapExactAmountIn) (*MsgSwapExactAmountInResponse, error) {
	if msgSwapExactAmountIn == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "swap null swap"}
	}

	msgServer := ammkeeper.NewMsgServerImpl(*f)

	var PoolIds []uint64
	var TokenOutDenoms []string

	for _, route := range msgSwapExactAmountIn.Routes {
		PoolIds = append(PoolIds, route.PoolId)
		TokenOutDenoms = append(TokenOutDenoms, route.TokenOutDenom)
	}

	msgMsgSwapExactAmountIn := ammtype.NewMsgSwapExactAmountIn(msgSwapExactAmountIn.Sender, msgSwapExactAmountIn.TokenIn, msgSwapExactAmountIn.TokenOutMinAmount, PoolIds, TokenOutDenoms)

	if err := msgMsgSwapExactAmountIn.ValidateBasic(); err != nil {
		return nil, errorsmod.Wrap(err, "failed validating MsgMsgSwapExactAmountIn")
	}

	// Swap
	swapResp, err := msgServer.SwapExactAmountIn(
		sdk.WrapSDKContext(ctx),
		msgMsgSwapExactAmountIn,
	)
	if err != nil {
		return nil, errorsmod.Wrap(err, "swap msg")
	}

	var resp = &MsgSwapExactAmountInResponse{
		TokenOutAmount: swapResp.TokenOutAmount,
		MetaData:       msgSwapExactAmountIn.MetaData,
	}
	return resp, nil
}

func (m *CustomMessenger) msgOpen(ctx sdk.Context, contractAddr sdk.AccAddress, msgOpen *MsgOpen) ([]sdk.Event, [][]byte, error) {
	res, err := PerformMsgOpen(m.margin, ctx, contractAddr, msgOpen)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "perform open")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize open response")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}

func PerformMsgOpen(f *marginkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, msgOpen *MsgOpen) (*MsgOpenResponse, error) {
	if msgOpen == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "margin open null margin open"}
	}
	msgServer := marginkeeper.NewMsgServerImpl(*f)

	msgMsgOpen := margintypes.NewMsgOpen(msgOpen.Creator, msgOpen.CollateralAsset, cosmos_sdk_math.Int(msgOpen.CollateralAmount), msgOpen.BorrowAsset, msgOpen.Position, msgOpen.Leverage, msgOpen.TakeProfitPrice)

	if err := msgMsgOpen.ValidateBasic(); err != nil {
		return nil, errorsmod.Wrap(err, "failed validating msgMsgOpen")
	}

	_, err := msgServer.Open(ctx, msgMsgOpen) // Discard the response because it's empty

	if err != nil {
		return nil, errorsmod.Wrap(err, "margin open msg")
	}

	var resp = &MsgOpenResponse{
		MetaData: msgOpen.MetaData,
	}
	return resp, nil
}

func (m *CustomMessenger) msgClose(ctx sdk.Context, contractAddr sdk.AccAddress, msgClose *MsgClose) ([]sdk.Event, [][]byte, error) {
	res, err := PerformMsgClose(m.margin, ctx, contractAddr, msgClose)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "perform close")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize close response")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}

func PerformMsgClose(f *marginkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, msgClose *MsgClose) (*MsgCloseResponse, error) {
	if msgClose == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "margin close null margin close"}
	}
	msgServer := marginkeeper.NewMsgServerImpl(*f)

	msgMsgClose := margintypes.NewMsgClose(msgClose.Creator, uint64(msgClose.Id))

	if err := msgMsgClose.ValidateBasic(); err != nil {
		return nil, errorsmod.Wrap(err, "failed validating msgMsgClose")
	}

	_, err := msgServer.Close(ctx, msgMsgClose) // Discard the response because it's empty

	if err != nil {
		return nil, errorsmod.Wrap(err, "margin close msg")
	}

	var resp = &MsgCloseResponse{
		MetaData: msgClose.MetaData,
	}
	return resp, nil
}

func (m *CustomMessenger) msgStake(ctx sdk.Context, contractAddr sdk.AccAddress, msgStake *MsgStake) ([]sdk.Event, [][]byte, error) {
	var res *RequestResponse
	var err error
	if msgStake.Asset == paramtypes.Elys {
		res, err = PerformMsgStakeElys(m.staking, ctx, contractAddr, msgStake)
		if err != nil {
			return nil, nil, errorsmod.Wrap(err, "perform elys stake")
		}
	} else {
		res, err = PerformMsgCommit(m.commitment, ctx, contractAddr, msgStake)
		if err != nil {
			return nil, nil, errorsmod.Wrap(err, "perform elys stake")
		}
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize stake")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}

func PerformMsgStakeElys(f *stakingkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, msgStake *MsgStake) (*RequestResponse, error) {
	if msgStake == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "Invalid staking parameter"}
	}

	msgServer := stakingkeeper.NewMsgServerImpl(f)
	address, err := sdk.AccAddressFromBech32(msgStake.Address)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid address")
	}

	validator_address, err := sdk.ValAddressFromBech32(msgStake.ValidatorAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid address")
	}

	amount := sdk.NewCoin(msgStake.Asset, msgStake.Amount)
	msgMsgDelegate := stakingtypes.NewMsgDelegate(address, validator_address, amount)

	if err := msgMsgDelegate.ValidateBasic(); err != nil {
		return nil, errorsmod.Wrap(err, "failed validating msgMsgDelegate")
	}

	_, err = msgServer.Delegate(ctx, msgMsgDelegate) // Discard the response because it's empty
	if err != nil {
		return nil, errorsmod.Wrap(err, "elys stake msg")
	}

	var resp = &RequestResponse{
		Code:   paramtypes.RES_OK,
		Result: "Staking succeed",
	}

	return resp, nil
}

func PerformMsgCommit(f *commitmentkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, msgStake *MsgStake) (*RequestResponse, error) {
	if msgStake == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "Invalid staking parameter"}
	}
	msgServer := commitmentkeeper.NewMsgServerImpl(*f)
	msgMsgCommit := commitmenttypes.NewMsgCommitTokens(msgStake.Address, msgStake.Amount, msgStake.Asset)

	if err := msgMsgCommit.ValidateBasic(); err != nil {
		return nil, errorsmod.Wrap(err, "failed validating msgMsgCommit")
	}

	_, err := msgServer.CommitTokens(ctx, msgMsgCommit) // Discard the response because it's empty
	if err != nil {
		return nil, errorsmod.Wrap(err, "commit msg")
	}

	var resp = &RequestResponse{
		Code:   paramtypes.RES_OK,
		Result: "Staking succeed",
	}
	return resp, nil
}

func (m *CustomMessenger) msgUnstake(ctx sdk.Context, contractAddr sdk.AccAddress, msgUnstake *MsgUnstake) ([]sdk.Event, [][]byte, error) {
	var res *RequestResponse
	var err error
	if msgUnstake.Asset == paramtypes.Elys {
		res, err = PerformMsgUnstakeElys(m.staking, ctx, contractAddr, msgUnstake)
		if err != nil {
			return nil, nil, errorsmod.Wrap(err, "perform elys unstake")
		}
	} else {
		res, err = PerformMsgUncommit(m.commitment, ctx, contractAddr, msgUnstake)
		if err != nil {
			return nil, nil, errorsmod.Wrap(err, "perform elys uncommit")
		}
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize stake")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}

func PerformMsgUnstakeElys(f *stakingkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, msgUnstake *MsgUnstake) (*RequestResponse, error) {
	if msgUnstake == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "Invalid unstaking parameter"}
	}

	msgServer := stakingkeeper.NewMsgServerImpl(f)
	address, err := sdk.AccAddressFromBech32(msgUnstake.Address)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid address")
	}

	validator_address, err := sdk.ValAddressFromBech32(msgUnstake.ValidatorAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid address")
	}

	amount := sdk.NewCoin(msgUnstake.Asset, msgUnstake.Amount)
	msgMsgUndelegate := stakingtypes.NewMsgUndelegate(address, validator_address, amount)

	if err := msgMsgUndelegate.ValidateBasic(); err != nil {
		return nil, errorsmod.Wrap(err, "failed validating msgMsgDelegate")
	}

	_, err = msgServer.Undelegate(ctx, msgMsgUndelegate) // Discard the response because it's empty
	if err != nil {
		return nil, errorsmod.Wrap(err, "elys unstake msg")
	}

	var resp = &RequestResponse{
		Code:   paramtypes.RES_OK,
		Result: "Unstaking succeed",
	}

	return resp, nil
}

func PerformMsgUncommit(f *commitmentkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, msgUnstake *MsgUnstake) (*RequestResponse, error) {
	if msgUnstake == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "Invalid staking parameter"}
	}
	msgServer := commitmentkeeper.NewMsgServerImpl(*f)
	msgMsgUncommit := commitmenttypes.NewMsgUncommitTokens(msgUnstake.Address, msgUnstake.Amount, msgUnstake.Asset)

	if err := msgMsgUncommit.ValidateBasic(); err != nil {
		return nil, errorsmod.Wrap(err, "failed validating msgMsgCommit")
	}

	_, err := msgServer.UncommitTokens(ctx, msgMsgUncommit) // Discard the response because it's empty
	if err != nil {
		return nil, errorsmod.Wrap(err, "commit msg")
	}

	var resp = &RequestResponse{
		Code:   paramtypes.RES_OK,
		Result: "Unstaking succeed",
	}
	return resp, nil
}

type ElysMsg struct {
	MsgSwapExactAmountIn *MsgSwapExactAmountIn `json:"msg_swap_exact_amount_in,omitempty"`
	MsgOpen              *MsgOpen              `json:"msg_open,omitempty"`
	MsgClose             *MsgClose             `json:"msg_close,omitempty"`
	MsgStake             *MsgStake             `json:"msg_stake,omitempty"`
	MsgUnstake           *MsgUnstake           `json:"msg_unstake,omitempty"`
}

type MsgSwapExactAmountIn struct {
	Sender            string                      `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender,omitempty"`
	Routes            []ammtype.SwapAmountInRoute `protobuf:"bytes,2,rep,name=routes,proto3" json:"routes,omitempty"`
	TokenIn           sdk.Coin                    `protobuf:"bytes,3,opt,name=tokenIn,proto3" json:"token_in,omitempty"`
	TokenOutMinAmount cosmos_sdk_math.Int         `protobuf:"bytes,4,opt,name=tokenOutMinAmount,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"token_out_min_amount,omitempty"`
	MetaData          *[]byte                     `protobuf:"bytes,5,opt,name=tokenData,proto3" json:"meta_data,omitempty"`
}

type MsgSwapExactAmountInResponse struct {
	TokenOutAmount cosmos_sdk_math.Int `protobuf:"bytes,1,opt,name=tokenOutAmount,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"token_out_amount,omitempty"`
	MetaData       *[]byte             `protobuf:"bytes,2,opt,name=tokenData,proto3" json:"meta_data,omitempty"`
}

type MsgOpen struct {
	Creator          string               `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty"`
	CollateralAsset  string               `protobuf:"bytes,2,opt,name=collateralAsset,proto3" json:"collateral_asset,omitempty"`
	CollateralAmount sdk.Uint             `protobuf:"bytes,3,opt,name=collateralAmount,proto3" json:"collateral_amount,omitempty"`
	BorrowAsset      string               `protobuf:"bytes,4,opt,name=borrowAsset,proto3" json:"borrow_asset,omitempty"`
	Position         margintypes.Position `protobuf:"bytes,5,opt,name=position,proto3" json:"position,omitempty"`
	Leverage         sdk.Dec              `protobuf:"bytes,6,opt,name=leverage,proto3" json:"leverage,omitempty"`
	TakeProfitPrice  sdk.Dec              `protobuf:"bytes,7,opt,name=takeProfitPrice,proto3" json:"take_profit_price,omitempty"`
	MetaData         *[]byte              `protobuf:"bytes,8,opt,name=tokenData,proto3" json:"meta_data,omitempty"`
}

type MsgClose struct {
	Creator  string  `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty"`
	Id       int64   `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	MetaData *[]byte `protobuf:"bytes,3,opt,name=tokenData,proto3" json:"meta_data,omitempty"`
}

type MsgOpenResponse struct {
	MetaData *[]byte `protobuf:"bytes,1,opt,name=tokenData,proto3" json:"meta_data,omitempty"`
}
type MsgCloseResponse struct {
	MetaData *[]byte `protobuf:"bytes,1,opt,name=tokenData,proto3" json:"meta_data,omitempty"`
}

type MsgStake struct {
	Address          string              `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Amount           cosmos_sdk_math.Int `protobuf:"bytes,2,opt,name=amount,proto3" json:"amount,omitempty"`
	Asset            string              `protobuf:"bytes,3,opt,name=asset,proto3" json:"asset,omitempty"`
	ValidatorAddress string              `protobuf:"bytes,4,opt,name=validator_address,proto3" json:"validator_address,omitempty"`
}

type MsgUnstake struct {
	Address          string              `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Amount           cosmos_sdk_math.Int `protobuf:"bytes,2,opt,name=amount,proto3" json:"amount,omitempty"`
	Asset            string              `protobuf:"bytes,3,opt,name=asset,proto3" json:"asset,omitempty"`
	ValidatorAddress string              `protobuf:"bytes,4,opt,name=validator_address,proto3" json:"validator_address,omitempty"`
}

type RequestResponse struct {
	Code   uint64 `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	Result string `protobuf:"bytes,2,opt,name=result,proto3" json:"result,omitempty"`
}
