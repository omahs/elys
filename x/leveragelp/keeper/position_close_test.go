package keeper_test

import (
	"time"

	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/leveragelp/types"
	stablestakekeeper "github.com/elys-network/elys/x/stablestake/keeper"
	stablestaketypes "github.com/elys-network/elys/x/stablestake/types"
)

func (suite KeeperTestSuite) OpenPosition(addr sdk.AccAddress) (*types.Position, types.Pool) {
	k := suite.app.LeveragelpKeeper
	SetupStableCoinPrices(suite.ctx, suite.app.OracleKeeper)
	poolAddr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	treasuryAddr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	pool := types.Pool{
		AmmPoolId:         1,
		Enabled:           true,
		Closed:            false,
		Health:            sdk.ZeroDec(),
		LeveragedLpAmount: sdk.ZeroInt(),
		LeverageMax:       sdk.ZeroDec(),
	}
	poolInit := sdk.Coins{sdk.NewInt64Coin("uusdc", 100000), sdk.NewInt64Coin("uusdt", 100000)}

	err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, poolInit)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, poolAddr, poolInit)
	suite.Require().NoError(err)

	suite.app.AmmKeeper.SetPool(suite.ctx, ammtypes.Pool{
		PoolId:            1,
		Address:           poolAddr.String(),
		RebalanceTreasury: treasuryAddr.String(),
		PoolParams: ammtypes.PoolParams{
			SwapFee:                     sdk.ZeroDec(),
			ExitFee:                     sdk.ZeroDec(),
			UseOracle:                   true,
			WeightBreakingFeeMultiplier: sdk.ZeroDec(),
			ExternalLiquidityRatio:      sdk.NewDec(1),
			LpFeePortion:                sdk.ZeroDec(),
			StakingFeePortion:           sdk.ZeroDec(),
			WeightRecoveryFeePortion:    sdk.ZeroDec(),
			ThresholdWeightDifference:   sdk.ZeroDec(),
			FeeDenom:                    "uusdc",
		},
		TotalShares: sdk.NewCoin("amm/pool/1", sdk.NewInt(2).Mul(ammtypes.OneShare)),
		PoolAssets: []ammtypes.PoolAsset{
			{
				Token:  poolInit[0],
				Weight: sdk.NewInt(10),
			},
			{
				Token:  poolInit[1],
				Weight: sdk.NewInt(10),
			},
		},
		TotalWeight: sdk.NewInt(20),
	})
	k.SetPool(suite.ctx, pool)
	suite.app.AmmKeeper.SetDenomLiquidity(suite.ctx, ammtypes.DenomLiquidity{
		Denom:     "uusdc",
		Liquidity: sdk.NewInt(100000),
	})
	suite.app.AmmKeeper.SetDenomLiquidity(suite.ctx, ammtypes.DenomLiquidity{
		Denom:     "uusdt",
		Liquidity: sdk.NewInt(100000),
	})

	usdcToken := sdk.NewInt64Coin("uusdc", 100000)
	err = suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, sdk.Coins{usdcToken})
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, addr, sdk.Coins{usdcToken})
	suite.Require().NoError(err)

	stableMsgServer := stablestakekeeper.NewMsgServerImpl(suite.app.StablestakeKeeper)
	_, err = stableMsgServer.Bond(sdk.WrapSDKContext(suite.ctx), &stablestaketypes.MsgBond{
		Creator: addr.String(),
		Amount:  sdk.NewInt(10000),
	})
	suite.Require().NoError(err)

	// open a position
	position, err := k.OpenLong(suite.ctx, &types.MsgOpen{
		Creator:          addr.String(),
		CollateralAsset:  "uusdc",
		CollateralAmount: sdk.NewInt(1000),
		AmmPoolId:        1,
		Leverage:         sdk.NewDec(5),
	})
	suite.Require().NoError(err)
	return position, pool
}

func (suite KeeperTestSuite) TestCloseLong() {
	k := suite.app.LeveragelpKeeper
	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())

	var (
		msg = &types.MsgClose{
			Creator: addr.String(),
			Id:      1,
		}
		repayAmount = math.NewInt(0)
	)

	_, repayAmountOut, err := k.CloseLong(suite.ctx, msg)
	suite.Require().Error(err)
	suite.Require().Equal(repayAmount.String(), repayAmountOut.String())
}

func (suite KeeperTestSuite) TestForceCloseLong() {
	k := suite.app.LeveragelpKeeper
	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	position, pool := suite.OpenPosition(addr)
	repayAmount := math.NewInt(4000)

	repayAmountOut, err := k.ForceCloseLong(suite.ctx, *position, pool)
	suite.Require().NoError(err)
	suite.Require().Equal(repayAmount.String(), repayAmountOut.String())
}

func (suite KeeperTestSuite) TestHealthDecreaseForInterest() {
	k := suite.app.LeveragelpKeeper
	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	position, _ := suite.OpenPosition(addr)
	ammPool, found := suite.app.AmmKeeper.GetPool(suite.ctx, position.AmmPoolId)
	suite.Require().True(found)
	health, err := k.GetPositionHealth(suite.ctx, *position, ammPool)
	suite.Require().NoError(err)
	suite.Require().Equal(health.String(), "1.221000000000000000")

	suite.ctx = suite.ctx.WithBlockTime(suite.ctx.BlockTime().Add(time.Hour * 24 * 365))
	suite.app.StablestakeKeeper.BeginBlocker(suite.ctx)
	health, err = k.GetPositionHealth(suite.ctx, *position, ammPool)
	suite.Require().NoError(err)
	suite.Require().Equal(health.String(), "0.610500000000000000")
}
