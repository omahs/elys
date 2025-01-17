package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/margin/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (k Keeper) CalcMTPConsolidateCollateral(ctx sdk.Context, mtp *types.MTP) error {
	consolidateCollateral := sdk.ZeroInt()
	for _, asset := range mtp.Collaterals {
		if asset.Denom == ptypes.BaseCurrency {
			consolidateCollateral = consolidateCollateral.Add(asset.Amount)
		} else {
			// swap into base currency
			_, ammPool, _, err := k.OpenChecker.PreparePools(ctx, asset.Denom)
			if err != nil {
				return err
			}

			collateralAmtIn := sdk.NewCoin(asset.Denom, asset.Amount)
			C, err := k.EstimateSwapGivenOut(ctx, collateralAmtIn, ptypes.BaseCurrency, ammPool)
			if err != nil {
				return err
			}

			consolidateCollateral = consolidateCollateral.Add(C)
		}
	}

	mtp.SumCollateral = consolidateCollateral

	return nil
}
