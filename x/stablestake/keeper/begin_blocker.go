package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) BeginBlocker(ctx sdk.Context) {
	// check if epoch has passed then execute
	epochLength := k.GetEpochLength(ctx)
	epochPosition := k.GetEpochPosition(ctx, epochLength)

	if epochPosition == 0 { // if epoch has passed
		params := k.GetParams(ctx)
		rate := k.InterestRateComputation(ctx)
		params.InterestRate = rate
		k.SetParams(ctx, params)

		debts := k.AllDebts(ctx)
		for _, debt := range debts {
			k.UpdateInterestStacked(ctx, debt)
		}
	}

}
