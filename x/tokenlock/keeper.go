package tokenlock

import (
	"time"

	"github.com/cosmos/cosmos-sdk/x/bank"

	"github.com/cosmos/cosmos-sdk/store/prefix"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/tendermint/libs/log"

)

// Keeper is the model object for the package tokenlock module
type Keeper struct {
	bankKeeper bank.BaseKeeper

	// The (unexposed) keys used to access the stores from the Context.
	storeKey sdk.StoreKey

	// The codec codec for binary encoding/decoding.
	cdc *codec.Codec

	// Reserved codespace
	codespace sdk.CodespaceType
}


func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, bk bank.BaseKeeper, codespace sdk.CodespaceType) Keeper {
	return Keeper {
		cdc: cdc,
		storeKey: storeKey,
		bankKeeper: bk,
		codespace: codespace,
	}
}

// Logger returns a module-specific logger.
func (keeper Keeper) Logger(ctx sdk.Context) log.Logger { return ctx.Logger().With("module", "x/tokenlock") }

func (keeper Keeper) GetOwnerLocks(ctx sdk.Context, owner sdk.AccAddress) (locks []TokenLock) {
	store := prefix.NewStore(ctx.KVStore(keeper.storeKey), PrefixLocks)
	iterator := store.Iterator(owner, sdk.PrefixEndBytes(owner))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var lock TokenLock

		keeper.cdc.MustUnmarshalBinaryBare(iterator.Value(), &lock)

		locks = append(locks, lock)
	}
	return locks
}

func (keeper Keeper) GetAllLocks(ctx sdk.Context) (locks []TokenLock) {
	store := prefix.NewStore(ctx.KVStore(keeper.storeKey), PrefixLocks)
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var lock TokenLock

		keeper.cdc.MustUnmarshalBinaryBare(iterator.Value(), &lock)

		locks = append(locks, lock)
	}
	return locks
}


func (keeper Keeper) GetAllUnlocks(ctx sdk.Context) (unlocks []TokenUnlock) {
	store := prefix.NewStore(ctx.KVStore(keeper.storeKey), PrefixUnlockQueue)
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var unlock TokenUnlock

		keeper.cdc.MustUnmarshalBinaryBare(iterator.Value(), &unlock)

		unlocks = append(unlocks, unlock)
	}
	return unlocks
}

func (keeper Keeper) GetLock(ctx sdk.Context, owner sdk.AccAddress, unlockTime time.Duration) (lock TokenLock) {
	store := prefix.NewStore(ctx.KVStore(keeper.storeKey), PrefixLocks)
	bz := store.Get(KeyLock(owner, unlockTime))
	if bz == nil {
		return TokenLock {
			Owner: owner,
			UnlockTime: unlockTime,
		}
	}
	keeper.cdc.MustUnmarshalBinaryBare(bz, &lock)
	return
}

func (keeper Keeper) setLock(ctx sdk.Context, lock TokenLock) {
	store := prefix.NewStore(ctx.KVStore(keeper.storeKey), PrefixLocks)
	if lock.Amount.IsZero() {
		store.Delete(KeyLock(lock.Owner, lock.UnlockTime))
		return
	}
	bz := keeper.cdc.MustMarshalBinaryBare(lock.Amount)
	store.Set(KeyLock(lock.Owner, lock.UnlockTime), bz)
}

func (keeper Keeper) LockCoins(ctx sdk.Context, owner sdk.AccAddress, unlockTime time.Duration, amount sdk.Coins) sdk.Error {
	_, err := keeper.bankKeeper.SubtractCoins(ctx, owner, amount)
	if err != nil {
		return err
	}
	lock := keeper.GetLock(ctx, owner, unlockTime)
	lock.Amount = lock.Amount.Add(amount)
	keeper.setLock(ctx, lock)
	return nil
}

func (keeper Keeper) BeginUnlock(ctx sdk.Context, owner sdk.AccAddress, unlockTime time.Duration, amount sdk.Coins) sdk.Error {
	lock := keeper.GetLock(ctx, owner, unlockTime)
	newAmount, errBool := lock.Amount.SafeSub(amount)
	if errBool {
		return ErrInsufficientCoins(keeper.codespace)
	}
	lock.Amount = newAmount
	keeper.setLock(ctx, lock)

	keeper.InsertUnlockQueue(ctx,
		TokenUnlock{
			Amount:     amount,
			CompletionTime: ctx.BlockHeader().Time.Add(unlockTime),
			Owner:      owner,
		},
	)
	return nil
}

func (keeper Keeper) FinishUnlock(ctx sdk.Context, unlock TokenUnlock) sdk.Error {
	if unlock.CompletionTime.After(ctx.BlockHeader().Time) {
		panic("unlocked too soon")
	}

	keeper.bankKeeper.AddCoins(ctx, unlock.Owner, unlock.Amount)
	return nil
}

// Returns an iterator for all the unlocks in the Unlock Queue that expire by endTime
func (keeper Keeper) UnlockQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := prefix.NewStore(ctx.KVStore(keeper.storeKey), PrefixUnlockQueue)
	return store.Iterator(nil, sdk.PrefixEndBytes(PrefixUnlockQueueTime(endTime)))
}

// Inserts a ProposalID into the active proposal queue at endTime
func (keeper Keeper) InsertUnlockQueue(ctx sdk.Context, unlock TokenUnlock) {
	store := prefix.NewStore(ctx.KVStore(keeper.storeKey), PrefixUnlockQueue)
	bz := keeper.cdc.MustMarshalBinaryBare(unlock)
	store.Set(KeyUnlock(unlock), bz)
}
