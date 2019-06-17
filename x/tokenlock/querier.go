package tokenlock

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the governance Querier
const (
	QueryLocks      = "locks"
	QueryUnlocks    = "unlocks"
	QueryOwnerLocks = "ownerlocks"
)

// NewQuerier creates a new querier
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryLocks:
			return queryLocks(ctx, path[1:], req, keeper)
		case QueryUnlocks:
			return queryUnlocks(ctx, path[1:], req, keeper)
		case QueryOwnerLocks:
			return queryOwnerLocks(ctx, path[1:], req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest(fmt.Sprintf("unknown %s query endpoint: %T", QuerierRoute, path[0]))
		}
	}
}

type QueryUserLocks struct {
	Owner sdk.AccAddress
}

// nolint: unparam
func queryOwnerLocks(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryUserLocks
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}

	locks := keeper.GetOwnerLocks(ctx, params.Owner)
	bz, err := codec.MarshalJSONIndent(keeper.cdc, locks)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func queryLocks(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	locks := keeper.GetAllLocks(ctx)
	bz, err := codec.MarshalJSONIndent(keeper.cdc, locks)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func queryUnlocks(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	unlocks := keeper.GetAllUnlocks(ctx)
	bz, err := codec.MarshalJSONIndent(keeper.cdc, unlocks)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}
