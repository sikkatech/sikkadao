package democraticvalset

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// query endpoints supported by the governance Querier
const (
	QueryValidator = "validator"
	QueryProposals = "validators"
)

func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryValidator:
			return queryValidator(ctx, path[1:], req, keeper)
		case QueryValidators:
			return queryValidators(ctx, path[1:], req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown %s query endpoint: %T", QuerierRoute, path[0])
		}
	}
}

// Params for querying a specific validator
type QueryValidatorParams struct {
	OperatorAddress sdk.ValAddress `json:"operator_address"` // address of the validator's operator; bech encoded in JSON
}

// nolint: unparam
func queryValidator(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryValidatorParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}

	val, err := keeper.GetValidator(ctx, params.OperatorAddress)
	if err != nil {
		return nil, err
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, val)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

// nolint: unparam
func queryValidators(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	vals := []Validator{}

	iterator := keeper.ValidatorIterator(ctx)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		val, _ := keeper.GetValidator(ctx, iterator.Key())
		vals = append(vals, val)
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, vals)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}
