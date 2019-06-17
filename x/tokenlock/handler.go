package tokenlock

import (
	"fmt"

	"github.com/sunnya97/cosmos-sdk-modules/x/tokenlock/tags"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler creates a new handler for tokenlock module
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgLockCoins:
			return handleMsgLockCoins(ctx, keeper, msg)
		case MsgUnlockCoins:
			return handleMsgUnlockCoins(ctx, keeper, msg)

		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", RouterKey, msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgLockCoins(ctx sdk.Context, keeper Keeper, msg MsgLockCoins) sdk.Result {
	err := keeper.LockCoins(ctx, msg.Owner, msg.UnlockTime, msg.Amount)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleMsgUnlockCoins(ctx sdk.Context, keeper Keeper, msg MsgUnlockCoins) sdk.Result {
	err := keeper.BeginUnlock(ctx, msg.Owner, msg.UnlockTime, msg.Amount)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: sdk.NewTags(tags.Sender, msg.Owner, tags.Category, tags.TxCategory, tags.Action, tags.ActionTokenUnlockStarted),
	}
}

// Called every block, process inflation, update validator set
func EndBlocker(ctx sdk.Context, keeper Keeper) sdk.Tags {
	logger := keeper.Logger(ctx)
	resTags := sdk.NewTags()

	iterator := keeper.UnlockQueueIterator(ctx, ctx.BlockHeader().Time)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var unlock TokenUnlock

		keeper.cdc.MustUnmarshalBinaryBare(iterator.Value(), &unlock)

		keeper.FinishUnlock(ctx, unlock)

		store := prefix.NewStore(ctx.KVStore(keeper.storeKey), PrefixUnlockQueue)
		store.Delete(iterator.Key())

		resTags = resTags.AppendTag(tags.Action, tags.ActionTokenUnlockCompleted)
		resTags = resTags.AppendTag(tags.Sender, unlock.Owner.String())

		logger.Info(
			fmt.Sprintf("unlocked %d to %s",
				unlock.Amount, unlock.Owner,
			),
		)
	}

	return resTags
}
