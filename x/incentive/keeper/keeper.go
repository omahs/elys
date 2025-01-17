package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	ammtypes "github.com/elys-network/elys/x/amm/types"
	ctypes "github.com/elys-network/elys/x/commitment/types"
	"github.com/elys-network/elys/x/incentive/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeKey     storetypes.StoreKey
		memKey       storetypes.StoreKey
		paramstore   paramtypes.Subspace
		cmk          types.CommitmentKeeper
		stk          types.StakingKeeper
		tci          *types.TotalCommitmentInfo
		authKeeper   types.AccountKeeper
		bankKeeper   types.BankKeeper
		amm          types.AmmKeeper
		oracleKeeper types.OracleKeeper

		feeCollectorName    string // name of the FeeCollector ModuleAccount
		dexRevCollectorName string // name of the Dex Revenue ModuleAccount
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	ck types.CommitmentKeeper,
	sk types.StakingKeeper,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	amm types.AmmKeeper,
	ok types.OracleKeeper,
	feeCollectorName string,
	dexRevCollectorName string,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:                 cdc,
		storeKey:            storeKey,
		memKey:              memKey,
		paramstore:          ps,
		cmk:                 ck,
		stk:                 sk,
		tci:                 &types.TotalCommitmentInfo{},
		feeCollectorName:    feeCollectorName,
		dexRevCollectorName: dexRevCollectorName,
		authKeeper:          ak,
		bankKeeper:          bk,
		amm:                 amm,
		oracleKeeper:        ok,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// Update uncommitted token amount
// Called back through epoch hook
func (k Keeper) UpdateUncommittedTokens(ctx sdk.Context, epochIdentifier string, stakeIncentive types.IncentiveInfo, lpIncentive types.IncentiveInfo) {
	// Recalculate total committed info
	k.UpdateTotalCommitmentInfo(ctx)

	// Collect DEX revenue while tracking 65% of it for LPs reward calculation
	// Assume these are collected in USDC
	dexRevenue, dexRevenueForLps := k.CollectDEXRevenue(ctx)

	// Calculate each portion of DEX revenue - stakers, LPs
	dexRevenueDec := sdk.NewDecCoinsFromCoins(dexRevenue...)
	dexRevenueForStakers := dexRevenueDec.Sub(dexRevenueForLps)

	// Calculate each portion of Gas fees collected - stakers, LPs
	gasFeeCollectedDec := sdk.NewDecCoinsFromCoins(k.tci.TotalFeesCollected...)
	rewardPortionForLps := k.GetDEXRewardPortionForLPs(ctx)
	gasFeesForLps := gasFeeCollectedDec.MulDecTruncate(rewardPortionForLps)
	gasFeesForStakers := gasFeeCollectedDec.Sub(gasFeesForLps)

	// Sum Dex revenue for stakers + Gas fees for stakers and name it dex Revenus for stakers
	// But won't sum dex revenue for LPs and gas fees for LPs as the LP revenue will be rewared by pool.
	dexRevenueForStakers = dexRevenueForStakers.Add(gasFeesForStakers...)

	// Fund community pool based on the communtiy tax
	dexRevenueRemainedForStakers := k.UpdateCommunityPool(ctx, dexRevenueForStakers)

	// TODO:
	// Dummy denom for USDC
	// USDC amount in sdk.Dec type
	dexRevenueLPsAmt := dexRevenueForLps.AmountOf(ptypes.BaseCurrency)
	dexRevenueStakersAmt := dexRevenueRemainedForStakers.AmountOf(ptypes.BaseCurrency)
	gasFeesLPsAmt := gasFeesForLps.AmountOf(ptypes.BaseCurrency)

	// Calculate eden amount per epoch
	edenAmountPerEpochStakers := stakeIncentive.Amount.Quo(sdk.NewInt(stakeIncentive.NumEpochs))
	edenAmountPerEpochLPs := lpIncentive.Amount.Quo(sdk.NewInt(lpIncentive.NumEpochs))
	edenBoostAPR := stakeIncentive.EdenBoostApr

	// Proxy TVL
	// Multiplier on each liquidity pool
	// We have 3 pools of 20, 30, 40 TVL
	// We have mulitplier of 0.3, 0.5, 1.0
	// Proxy TVL = 20*0.3+30*0.5+40*1.0
	totalProxyTVL := k.CalculateProxyTVL(ctx)

	totalEdenGiven := sdk.ZeroInt()
	totalEdenGivenLP := sdk.ZeroInt()
	totalRewardsGiven := sdk.ZeroInt()
	totalRewardsGivenLP := sdk.ZeroInt()
	// Process to increase uncomitted token amount of Eden & Eden boost
	k.cmk.IterateCommitments(
		ctx, func(commitments ctypes.Commitments) bool {
			// Commitment owner
			creator := commitments.Creator
			_, err := sdk.AccAddressFromBech32(creator)
			if err != nil {
				// This could be validator address
				return false
			}

			// Calculate delegated amount per delegator
			delegatedAmt := k.CalculateDelegatedAmount(ctx, creator)

			// Calculate new uncommitted Eden tokens from Eden & Eden boost committed, Dex rewards distribution
			// Distribute gas fees to stakers

			// Calculate new uncommitted Eden tokens from Elys staked
			newUncommittedEdenTokens, dexRewards, dexRewardsByStakers := k.CalculateRewardsForStakersByElysStaked(ctx, delegatedAmt, edenAmountPerEpochStakers, dexRevenueStakersAmt)
			totalEdenGiven = totalEdenGiven.Add(newUncommittedEdenTokens)
			totalRewardsGiven = totalRewardsGiven.Add(dexRewards)

			// Calculate new uncommitted Eden tokens from Eden committed
			edenCommitted := commitments.GetCommittedAmountForDenom(ptypes.Eden)
			newUncommittedEdenTokens, dexRewards = k.CalculateRewardsForStakersByCommitted(ctx, edenCommitted, edenAmountPerEpochStakers, dexRevenueStakersAmt)
			totalEdenGiven = totalEdenGiven.Add(newUncommittedEdenTokens)
			totalRewardsGiven = totalRewardsGiven.Add(dexRewards)

			// Calculate new uncommitted Eden tokens from Eden Boost committed
			edenBoostCommitted := commitments.GetCommittedAmountForDenom(ptypes.EdenB)
			newUncommittedEdenTokens, dexRewards = k.CalculateRewardsForStakersByCommitted(ctx, edenBoostCommitted, edenAmountPerEpochStakers, dexRevenueStakersAmt)
			totalEdenGiven = totalEdenGiven.Add(newUncommittedEdenTokens)
			totalRewardsGiven = totalRewardsGiven.Add(dexRewards)

			// Calculate new uncommitted Eden tokens from LpTokens committed, Dex rewards distribution
			// Distribute gas fees to LPs
			newUncommittedEdenTokensLp, dexRewardsLp := k.CalculateRewardsForLPs(ctx, totalProxyTVL, commitments, edenAmountPerEpochLPs, gasFeesLPsAmt)
			totalEdenGivenLP = totalEdenGivenLP.Add(newUncommittedEdenTokensLp)
			totalRewardsGivenLP = totalRewardsGivenLP.Add(dexRewardsLp)

			// Calculate the total Eden uncommitted amount
			newUncommittedEdenTokens = newUncommittedEdenTokens.Add(newUncommittedEdenTokensLp)

			// Give commission to validators ( Eden from stakers and Dex rewards from stakers. )
			edenCommissionGiven, dexRewardsCommissionGiven := k.GiveCommissionToValidators(ctx, creator, delegatedAmt, newUncommittedEdenTokens, dexRewardsByStakers)

			// Minus the given amount and increase with the remains only
			newUncommittedEdenTokens = newUncommittedEdenTokens.Sub(edenCommissionGiven)

			// Plus LpDexRewards and minus commission given
			dexRewards = dexRewards.Add(dexRewardsLp).Sub(dexRewardsCommissionGiven)

			// Calculate new uncommitted Eden-Boost tokens for staker and Eden token holders
			newUncommittedEdenBoostTokens := k.CalculateEdenBoostRewards(ctx, delegatedAmt, commitments, epochIdentifier, edenBoostAPR)

			// Update Commitments with new uncommitted token amounts
			k.UpdateCommitments(ctx, creator, &commitments, newUncommittedEdenTokens, newUncommittedEdenBoostTokens, dexRewards)

			return false
		},
	)

	// Calcualte the remainings
	edenRemained := edenAmountPerEpochStakers.Sub(totalEdenGiven)
	edenRemainedLP := edenAmountPerEpochLPs.Sub(totalEdenGivenLP)
	dexRewardsRemained := dexRevenueStakersAmt.Sub(sdk.NewDecFromInt(totalRewardsGiven))
	dexRewardsRemainedLP := dexRevenueLPsAmt.Add(gasFeesLPsAmt).Sub(sdk.NewDecFromInt(totalRewardsGivenLP))

	// Fund community the remain coins
	// ----------------------------------
	edenRemainedCoin := sdk.NewDecCoin(ptypes.Eden, edenRemained.Add(edenRemainedLP))
	// TODO:
	// Dummy denom for USDC
	dexRewardsRemainedCoin := sdk.NewDecCoinFromDec(ptypes.BaseCurrency, dexRewardsRemained.Add(dexRewardsRemainedLP))

	feePool := k.GetFeePool(ctx)
	feePool.CommunityPool = feePool.CommunityPool.Add(edenRemainedCoin)
	feePool.CommunityPool = feePool.CommunityPool.Add(dexRewardsRemainedCoin)
	k.SetFeePool(ctx, feePool)
	// ----------------------------------
}

func (k Keeper) UpdateCommitments(ctx sdk.Context, creator string, commitments *ctypes.Commitments, newUncommittedEdenTokens sdk.Int, newUncommittedEdenBoostTokens sdk.Int, dexRewards sdk.Int) {
	// Update uncommitted Eden balances in the Commitments structure
	k.UpdateTokensCommitment(commitments, newUncommittedEdenTokens, ptypes.Eden)
	// Update uncommitted Eden-Boost token balances in the Commitments structure
	k.UpdateTokensCommitment(commitments, newUncommittedEdenBoostTokens, ptypes.EdenB)

	// All dex revenue are collected to incentive module in USDC
	// Gas fees (Elys) are also converted into USDC and collected into total dex revenue wallet of incentive module.
	// Update USDC balances in the Commitments structure.
	// TODO:
	// USDC token denom is dummy one for now until we have real usdc brought through bridge.
	// These are the rewards from each pool, margin, gas fee.
	k.UpdateTokensCommitment(commitments, dexRewards, ptypes.BaseCurrency)

	// Save the updated Commitments
	k.cmk.SetCommitments(ctx, *commitments)
}

// Update the uncommitted Eden token balance
func (k Keeper) UpdateTokensCommitment(commitments *ctypes.Commitments, new_uncommitted_eden_tokens sdk.Int, denom string) {
	uncommittedEden, found := commitments.GetUncommittedTokensForDenom(denom)
	if !found {
		uncommittedTokens := commitments.GetUncommittedTokens()
		uncommittedTokens = append(uncommittedTokens, &ctypes.UncommittedTokens{
			Denom:  denom,
			Amount: new_uncommitted_eden_tokens,
		})
		commitments.UncommittedTokens = uncommittedTokens
	} else {
		uncommittedEden.Amount = uncommittedEden.Amount.Add(new_uncommitted_eden_tokens)
	}
}

// Calculate Proxy TVL
func (k Keeper) CalculateProxyTVL(ctx sdk.Context) sdk.Dec {
	multipliedShareSum := sdk.ZeroDec()
	k.amm.IterateLiquidityPools(ctx, func(p ammtypes.Pool) bool {
		tvl, err := p.TVL(ctx, k.oracleKeeper)
		if err != nil {
			return false
		}

		// Get pool info from incentive param
		poolInfo, found := k.GetPoolInfo(ctx, p.GetPoolId())
		if !found {
			return false
		}

		proxyTVL := tvl.Mul(poolInfo.Multiplier)

		// Calculate total pool share by TVL and multiplier
		multipliedShareSum = multipliedShareSum.Add(proxyTVL)

		return false
	})

	// return total sum of TVL share using multiplier of all pools
	return multipliedShareSum
}

// Caculate total TVL
func (k Keeper) CalculateTVL(ctx sdk.Context) sdk.Dec {
	TVL := sdk.ZeroDec()

	k.amm.IterateLiquidityPools(ctx, func(p ammtypes.Pool) bool {
		tvl, err := p.TVL(ctx, k.oracleKeeper)
		if err != nil {
			return false
		}
		TVL = TVL.Add(tvl)
		return false
	})

	return TVL
}
