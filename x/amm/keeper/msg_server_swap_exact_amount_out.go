package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

func (k msgServer) SwapExactAmountOut(goCtx context.Context, msg *types.MsgSwapExactAmountOut) (*types.MsgSwapExactAmountOutResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	// Try executing the tx on cached context environment, to filter invalid transactions out
	cacheCtx, _ := ctx.CacheContext()
	tokenInAmount, err := k.RouteExactAmountOut(cacheCtx, sender, msg.Routes, msg.TokenInMaxAmount, msg.TokenOut)
	if err != nil {
		return nil, err
	}

	lastSwapIndex := k.GetLastSwapRequestIndex(ctx)
	k.SetSwapExactAmountOutRequests(ctx, msg, lastSwapIndex+1)
	k.SetLastSwapRequestIndex(ctx, lastSwapIndex+1)

	// Swap event is handled elsewhere
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})

	return &types.MsgSwapExactAmountOutResponse{TokenInAmount: tokenInAmount}, nil
}
