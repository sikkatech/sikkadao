package tokenlock

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

// getMockApp returns an initialized mock application for this module.
func getMockApp(t *testing.T) (*mock.App, Keeper) {
	mApp := mock.NewApp()

	RegisterCodec(mApp.Cdc)

	keyTokenlock := sdk.NewKVStoreKey("tokenlock")

	bankKeeper := bank.NewBaseKeeper(mApp.AccountKeeper, mApp.ParamsKeeper.Subspace(bank.DefaultParamspace), bank.DefaultCodespace)
	keeper := NewKeeper(mApp.Cdc, keyTokenlock, bankKeeper, DefaultCodespace)

	mApp.Router().AddRoute(RouterKey, NewHandler(keeper))
	mApp.SetEndBlocker(getEndBlocker(keeper))
	mApp.SetInitChainer(getInitChainer(mApp, keeper, mApp.AccountKeeper))

	require.NoError(t, mApp.CompleteSetup(keyTokenlock))
	return mApp, keeper
}

// getEndBlocker returns a staking endblocker.
func getEndBlocker(keeper Keeper) sdk.EndBlocker {
	return func(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
		tags := EndBlocker(ctx, keeper)

		return abci.ResponseEndBlock{
			ValidatorUpdates: nil,
			Tags:             tags,
		}
	}
}

// getInitChainer initializes the chainer of the mock app and sets the genesis
// state. It returns an empty ResponseInitChain.
func getInitChainer(mapp *mock.App, keeper Keeper, accountKeeper types.AccountKeeper) sdk.InitChainer {
	return func(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
		mapp.InitChainer(ctx, req)

		tokenlockGenesis := DefaultGenesisState()
		InitGenesis(ctx, keeper, tokenlockGenesis)
		return abci.ResponseInitChain{}
	}
}

//__________________________________________________________________________________________

// func checkValidator(t *testing.T, mapp *mock.App, keeper Keeper,
// 	addr sdk.ValAddress, expFound bool) Validator {

// 	ctxCheck := mapp.BaseApp.NewContext(true, abci.Header{})
// 	validator, found := keeper.GetValidator(ctxCheck, addr)

// 	require.Equal(t, expFound, found)
// 	return validator
// }

// func checkDelegation(
// 	t *testing.T, mapp *mock.App, keeper Keeper, delegatorAddr sdk.AccAddress,
// 	validatorAddr sdk.ValAddress, expFound bool, expShares sdk.Dec,
// ) {

// 	ctxCheck := mapp.BaseApp.NewContext(true, abci.Header{})
// 	delegation, found := keeper.GetDelegation(ctxCheck, delegatorAddr, validatorAddr)
// 	if expFound {
// 		require.True(t, found)
// 		require.True(sdk.DecEq(t, expShares, delegation.Shares))

// 		return
// 	}

// 	require.False(t, found)
// }

func Test(t *testing.T) {
	mApp, keeper := getMockApp(t)
	genCoin := sdk.NewInt64Coin("foocoin", 100)

	acc1 := &auth.BaseAccount{
		Address: addr1,
		Coins:   sdk.Coins{genCoin},
	}
	accs := []auth.Account{acc1}

	mock.SetGenesis(mApp, accs)
	mock.CheckBalance(t, mApp, addr1, sdk.Coins{genCoin})

	// create validator
	lockCoinsMsg := NewMsgLockCoins(sdk.Coins{genCoin}, time.Hour * 5, acc1.Address)

	header := abci.Header{Height: mApp.LastBlockHeight() + 1}
	mock.SignCheckDeliver(t, mApp.Cdc, mApp.BaseApp, header, []sdk.Msg{lockCoinsMsg}, []uint64{0}, []uint64{0}, true, true, priv1)
	mock.CheckBalance(t, mApp, addr1, sdk.Coins{genCoin.Sub(genCoin)})

	header = abci.Header{Height: mApp.LastBlockHeight() + 1}
	mApp.BeginBlock(abci.RequestBeginBlock{Header: header})

	validator := checkValidator(t, mApp, keeper, sdk.ValAddress(addr1), true)
	require.Equal(t, sdk.ValAddress(addr1), validator.OperatorAddress)
	require.Equal(t, sdk.Bonded, validator.Status)
	require.True(sdk.IntEq(t, bondTokens, validator.BondedTokens()))

	header = abci.Header{Height: mApp.LastBlockHeight() + 1}
	mApp.BeginBlock(abci.RequestBeginBlock{Header: header})

	// edit the validator
	description = NewDescription("bar_moniker", "", "", "")
	editValidatorMsg := NewMsgEditValidator(sdk.ValAddress(addr1), description, nil, nil)

	header = abci.Header{Height: mApp.LastBlockHeight() + 1}
	mock.SignCheckDeliver(t, mApp.Cdc, mApp.BaseApp, header, []sdk.Msg{editValidatorMsg}, []uint64{0}, []uint64{1}, true, true, priv1)

	validator = checkValidator(t, mApp, keeper, sdk.ValAddress(addr1), true)
	require.Equal(t, description, validator.Description)

	// delegate
	mock.CheckBalance(t, mApp, addr2, sdk.Coins{genCoin})
	delegateMsg := NewMsgDelegate(addr2, sdk.ValAddress(addr1), bondCoin)

	header = abci.Header{Height: mApp.LastBlockHeight() + 1}
	mock.SignCheckDeliver(t, mApp.Cdc, mApp.BaseApp, header, []sdk.Msg{delegateMsg}, []uint64{1}, []uint64{0}, true, true, priv2)
	mock.CheckBalance(t, mApp, addr2, sdk.Coins{genCoin.Sub(bondCoin)})
	checkDelegation(t, mApp, keeper, addr2, sdk.ValAddress(addr1), true, bondTokens.ToDec())

	// begin unbonding
	beginUnbondingMsg := NewMsgUndelegate(addr2, sdk.ValAddress(addr1), bondCoin)
	header = abci.Header{Height: mApp.LastBlockHeight() + 1}
	mock.SignCheckDeliver(t, mApp.Cdc, mApp.BaseApp, header, []sdk.Msg{beginUnbondingMsg}, []uint64{1}, []uint64{1}, true, true, priv2)

	// delegation should exist anymore
	checkDelegation(t, mApp, keeper, addr2, sdk.ValAddress(addr1), false, sdk.Dec{})

	// balance should be the same because bonding not yet complete
	mock.CheckBalance(t, mApp, addr2, sdk.Coins{genCoin.Sub(bondCoin)})
}
