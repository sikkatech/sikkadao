package democraticvalset

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmtypes "github.com/tendermint/tendermint/types"
)


// Handle all "democraticvalset" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgUpdateDescription:
			return handleMsgDeposit(ctx, keeper, msg)
		case MsgSubmitProposal:
			return handleMsgSubmitProposal(ctx, keeper, msg)
		case MsgVote:
			return handleMsgVote(ctx, keeper, msg)

		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", RouterKey, msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}


// Called every block, update validator set
func EndBlocker(ctx sdk.Context, k Keeper) ([]abci.ValidatorUpdate, sdk.Tags) {
	tStore := ctx.TransientStore(k.transientStoreKey)

	updates := []abci.ValidatorUpdate
	iterator := sdk.KVStorePrefixIterator(tStore, []byte)

	for ; iterator.Valid(); iterator.Next() {

		updates = append(updates, abci.ValidatorUpdate{
			PubKey: tmtypes.TM2PB.PubKey(sdk.ConsPubKey(iterator.Key())),
			Power: int64(iterator.Value())
		})
	}

	return updates, nil
}
