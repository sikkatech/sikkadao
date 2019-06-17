package dao

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

//nolint
var (
	DaoTokenKey  = []byte{0x00} // key for the token of the dao
	FundsPoolKey = []byte{0x01} // key for the pool of coins
)

// DAO Keeper
type Keeper struct {
	// The reference to the Param Keeper to get and set Global Params
	paramsKeeper params.Keeper

	// The reference to the Paramstore to get and set dao specific params
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

func (keeper Keeper) getFunds(ctx sdk.Context) sdk.Coins {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(FundsPoolKey)
	if bz == nil {
		panic("genesis not initialized properly")
	}
	var funds sdk.Coins
	keeper.cdc.MustUnmarshalBinaryBare(bz, &funds)
	return funds
}

func (keeper Keeper) setFunds(ctx sdk.Context, funds sdk.Coins) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryBare(funds)
	store.Set(FundsPoolKey, bz)
}

func (keeper Keeper) DepositCoins(ctx sdk.Context, depositor sdk.AccAddress, amount sdk.Coins) (Validator, sdk.Error) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(validator.OperatorAddress)
	if bz == nil {
		return Validator, ErrNonexistantValidator(DefaultCodespace, operAddr)
	}
	var val Validator
	keeper.cdc.MustUnmarshalBinaryBare(bz, &val)
	return val
}

func (keeper Keeper) SetValidator(ctx sdk.Context, validator Validator) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryBare(Validator)
	store.Set(validator.OperatorAddress, bz)

	tStore := ctx.TransientStore(keeper.transientStoreKey)
	tStore.Set(validator.ConsPubKey, []byte(validator.Power))
}

func (keeper Keeper) RemoveValidator(ctx sdk.Context, operAddr sdk.ValAddress) sdk.Error {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(operAddr)
	if bz == nil {
		return ErrNonexistantValidator(DefaultCodespace, operAddr)
	}
	var val Validator
	keeper.cdc.MustUnmarshalBinaryBare(bz, &val)
	store.Delete(operAddr)

	tStore := ctx.TransientStore(keeper.transientStoreKey)
	tStore.Set(val.ConsPubKey, []byte(0))
	return nil
}

func (keeper Keeper) UpdateValidatorPower(ctx sdk.Context, operAddr sdk.ValAddress, newPower int64) sdk.Error {
	val, err := keeper.GetValidator(ctx, operAddr)
	if err != nil {
		return err
	}
	val.Power = newPower
	keeper.SetValidator(ctx, val)
}

func (keeper Keeper) UpdateValidatorConsPubKey(ctx sdk.Context, operAddr sdk.ValAddress, newConsPubKey sdk.ConsPubKey) sdk.Error {
	val, err := keeper.GetValidator(ctx, operAddr)
	if err != nil {
		return err
	}
	val.ConsPubKey = newConsPubKey
	keeper.SetValidator(ctx, val)
}

func (keeper Keeper) UpdateValidatorDescription(ctx sdk.Context, operAddr sdk.ValAddress, updateDescription Description) {
	val, err := keeper.GetValidator(ctx, operAddr)
	if err != nil {
		return err
	}
	newDescription, err := val.Description.UpdateDescription(updateDescription)
	if err != nil {
		return err
	}
	val.Description = newDescription
	keeper.SetValidator(ctx, val)
}

func (keeper Keeper) ValidatorIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(keeper.storeKey)
	return store.Iterator(nil, nil)
}
