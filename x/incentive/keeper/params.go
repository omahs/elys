package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	ctypes "github.com/elys-network/elys/x/commitment/types"
	etypes "github.com/elys-network/elys/x/epochs/types"
	"github.com/elys-network/elys/x/incentive/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	var params types.Params
	k.paramstore.GetParamSet(ctx, &params)
	return params
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

// GetCommunityTax returns the current distribution community tax.
func (k Keeper) GetCommunityTax(ctx sdk.Context) (percent sdk.Dec) {
	k.paramstore.Get(ctx, types.ParamStoreKeyCommunityTax, &percent)
	return percent
}

// GetWithdrawAddrEnabled returns the current distribution withdraw address
// enabled parameter.
func (k Keeper) GetWithdrawAddrEnabled(ctx sdk.Context) (enabled bool) {
	k.paramstore.Get(ctx, types.ParamStoreKeyWithdrawAddrEnabled, &enabled)
	return enabled
}

// GetDEXRewardPortionForLPs returns the dex revenue percent for Lps
func (k Keeper) GetDEXRewardPortionForLPs(ctx sdk.Context) (percent sdk.Dec) {
	k.paramstore.Get(ctx, types.ParamStoreKeyRewardPortionForLps, &percent)
	return percent
}

// GetPoolInfo
func (k Keeper) GetPoolInfo(ctx sdk.Context, poolId uint64) (types.PoolInfo, bool) {
	// Fetch incentive params
	params := k.GetParams(ctx)

	poolInfos := params.PoolInfos
	for _, ps := range poolInfos {
		if ps.PoolId == poolId {
			return ps, true
		}
	}

	return types.PoolInfo{}, false
}

// InitPoolMultiplier: create a pool information responding to the pool creation.
func (k Keeper) InitPoolMultiplier(ctx sdk.Context, poolId uint64) bool {
	// Fetch incentive params
	params := k.GetParams(ctx)
	poolInfos := params.PoolInfos

	for _, ps := range poolInfos {
		if ps.PoolId == poolId {
			return true
		}
	}

	// Initiate a new pool info
	poolInfo := types.PoolInfo{
		// reward amount
		PoolId: poolId,
		// reward wallet address
		RewardWallet: ammtypes.NewPoolRevenueAddress(poolId).String(),
		// multiplier for lp rewards
		Multiplier: sdk.NewDec(1),
	}

	// Update pool information
	params.PoolInfos = append(params.PoolInfos, poolInfo)
	k.SetParams(ctx, params)

	return true
}

// UpdatePoolMultipliers updates pool multipliers through gov proposal
func (k Keeper) UpdatePoolMultipliers(ctx sdk.Context, poolMultipliers []types.PoolMultipliers) bool {
	if len(poolMultipliers) < 1 {
		return false
	}

	// Fetch incentive params
	params := k.GetParams(ctx)

	// Update pool multiplier
	for _, pm := range poolMultipliers {
		for i, p := range params.PoolInfos {
			// If we found matching poolId
			if p.PoolId == pm.PoolId {
				params.PoolInfos[i].Multiplier = pm.Multiplier
			}
		}
	}

	// Update parameter
	k.SetParams(ctx, params)

	return true
}

// Find out active incentive params
func (k Keeper) GetProperIncentiveParam(ctx sdk.Context, epochIdentifier string) (bool, types.IncentiveInfo, types.IncentiveInfo) {
	// Fetch incentive params
	params := k.GetParams(ctx)

	// Update params
	defer k.SetParams(ctx, params)

	// If we don't have enough params
	if len(params.StakeIncentives) < 1 || len(params.LpIncentives) < 1 {
		return false, types.IncentiveInfo{}, types.IncentiveInfo{}
	}

	// Current block timestamp
	timestamp := ctx.BlockTime().Unix()
	foundIncentive := false

	// Incentive params initialize
	stakeIncentive := params.StakeIncentives[0]
	lpIncentive := params.LpIncentives[0]

	// Consider epochIdentifier and start time
	// Consider epochNumber as well
	if stakeIncentive.EpochIdentifier != epochIdentifier || timestamp < stakeIncentive.StartTime.Unix() {
		return false, types.IncentiveInfo{}, types.IncentiveInfo{}
	}

	// Increase current epoch of Stake incentive param
	stakeIncentive.CurrentEpoch = stakeIncentive.CurrentEpoch + 1
	if stakeIncentive.CurrentEpoch == stakeIncentive.NumEpochs {
		if len(params.StakeIncentives) > 1 {
			params.StakeIncentives = params.StakeIncentives[1:]
		} else {
			return false, types.IncentiveInfo{}, types.IncentiveInfo{}
		}
	}

	// Increase current epoch of Lp incentive param
	lpIncentive.CurrentEpoch = lpIncentive.CurrentEpoch + 1
	if lpIncentive.CurrentEpoch == lpIncentive.NumEpochs {
		if len(params.LpIncentives) > 1 {
			params.LpIncentives = params.LpIncentives[1:]
		} else {
			return false, types.IncentiveInfo{}, types.IncentiveInfo{}
		}
	}

	// return found, stake, lp incentive params
	return foundIncentive, stakeIncentive, lpIncentive
}

// Calculate epoch counts per year to be used in APR calculation
func (k Keeper) CalculateEpochCountsPerYear(epochIdentifier string) int64 {
	switch epochIdentifier {
	case etypes.WeekEpochID:
		return ptypes.WeeksPerYear
	case etypes.DayEpochID:
		return ptypes.DaysPerYear
	case etypes.HourEpochID:
		return ptypes.HoursPerYear
	}

	return 0
}

// Update total commitment info
func (k Keeper) UpdateTotalCommitmentInfo(ctx sdk.Context) {
	// Fetch total staked Elys amount again
	k.tci.TotalElysBonded = k.stk.TotalBondedTokens(ctx)
	// Initialize with amount zero
	k.tci.TotalEdenEdenBoostCommitted = sdk.ZeroInt()
	// Initialize with amount zero
	k.tci.TotalFeesCollected = sdk.Coins{}
	// Initialize Lp tokens amount
	k.tci.TotalLpTokensCommitted = make(map[string]sdk.Int)
	// ReInitialize Pool revenue tracker
	k.tci.PoolRevenueTrack = make(map[string]sdk.Dec)

	// Collect gas fees collected
	fees := k.CollectGasFeesToIncentiveModule(ctx)
	// Calculate total fees - Gas fees collected
	k.tci.TotalFeesCollected = k.tci.TotalFeesCollected.Add(fees...)

	// Iterate to calculate total Eden, Eden boost and Lp tokens committed
	k.cmk.IterateCommitments(ctx, func(commitments ctypes.Commitments) bool {
		committedEdenToken := commitments.GetCommittedAmountForDenom(ptypes.Eden)
		committedEdenBoostToken := commitments.GetCommittedAmountForDenom(ptypes.EdenB)

		k.tci.TotalEdenEdenBoostCommitted = k.tci.TotalEdenEdenBoostCommitted.Add(committedEdenToken).Add(committedEdenBoostToken)

		// Iterate to calculate total Lp tokens committed
		k.amm.IterateLiquidityPools(ctx, func(p ammtypes.Pool) bool {
			lpToken := ammtypes.GetPoolShareDenom(p.GetPoolId())

			committedLpToken := commitments.GetCommittedAmountForDenom(lpToken)
			amt, ok := k.tci.TotalLpTokensCommitted[lpToken]
			if !ok {
				k.tci.TotalLpTokensCommitted[lpToken] = committedLpToken
			} else {
				k.tci.TotalLpTokensCommitted[lpToken] = amt.Add(committedLpToken)
			}
			return false
		})
		return false
	})
}
