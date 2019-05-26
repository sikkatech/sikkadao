package democraticvalset

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// Governance Keeper
type Keeper struct {
	// The reference to the Param Keeper to get and set Global Params
	paramsKeeper params.Keeper

	// The reference to the Paramstore to get and set gov specific params
	paramSpace params.Subspace

	// The (unexposed) keys used to access the stores from the Context.
	storeKey sdk.StoreKey

	// The (unexposed) keys used to access the transient stores from the Context.
	transientStoreKey sdk.StoreKey

	// The codec codec for binary encoding/decoding.
	cdc *codec.Codec

	// Reserved codespace
	codespace sdk.CodespaceType
}

func (keeper Keeper) SetValidator(ctx sdk.Context, validator Validator) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryBare(Validator)
	store.Set(validator.OperatorAddress, bz)

	tStore := ctx.TransientStore(keeper.transientStoreKey)
	tStore.Set(validator.ConsPubKey, []byte(validator.Power))
}

func (keeper Keeper) RemoveValidator(ctx sdk.Context, operAddr sdk.ValAddress) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(operAddr)
	var val Validator
	keeper.cdc.MustUnmarshalBinaryBare(bz, &val)
	store.Delete(operAddr)

	tStore := ctx.TransientStore(keeper.transientStoreKey)
	tStore.Set(val.ConsPubKey, []byte(0))
}

func (keeper Keeper) UpdateValidatorPower(ctx sdk.Context, operAddr sdk.ValAddress, newPower int64) sdk.Error {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(operAddr)
	if bz == nil {
		return ErrNonexistantValidator(DefaultCodespace, operAddr)
	}
	var val Validator
	keeper.cdc.MustUnmarshalBinaryBare(bz, &val)
	val.Power = newPower

	keeper.SetValidator(ctx, val)
}

func (keeper Keeper) UpdateValidatorConsPubKey(ctx sdk.Context, operAddr sdk.ValAddress, newConsPubKey sdk.ConsPubKey) sdk.Error {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(operAddr)
	if bz == nil {
		return ErrNonexistantValidator(DefaultCodespace, operAddr)
	}
	var val Validator
	keeper.cdc.MustUnmarshalBinaryBare(bz, &val)
	val.ConsPubKey = newConsPubKey

	keeper.SetValidator(ctx, val)
}

func (keeper Keeper) UpdateValidatorDescription(ctx sdk.Context, operAddr sdk.ValAddress, updateDescription Description) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(operAddr)
	if bz == nil {
		return ErrNonexistantValidator(DefaultCodespace, operAddr)
	}
	var val Validator
	keeper.cdc.MustUnmarshalBinaryBare(bz, &val)
	newDescription, err := val.Description.UpdateDescription(updateDescription)
	if err != nil {
		return err
	}

	keeper.SetValidator(ctx, val)
}
