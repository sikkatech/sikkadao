package tokenlock

import sdk "github.com/cosmos/cosmos-sdk/types"

// GenesisState defines genesis data for the module
type GenesisState struct {
	TokenLocks   []TokenLock   `json:"tokenlocks"`
	TokenUnlocks []TokenUnlock `json:"tokenunlocks"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState() GenesisState {
	return GenesisState{
		TokenLocks:   nil,
		TokenUnlocks: nil,
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState { return NewGenesisState() }

// InitGenesis initializes story state from genesis file
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	for _, lock := range data.TokenLocks {
		keeper.setLock(ctx, lock)
	}
	for _, unlock := range data.TokenUnlocks {
		keeper.InsertUnlockQueue(ctx, unlock)
	}
}

// ExportGenesis exports the genesis state
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	return GenesisState{
		TokenLocks:   keeper.GetAllLocks(ctx),
		TokenUnlocks: keeper.GetAllUnlocks(ctx),
	}
}

// ValidateGenesis validates the genesis state data
func ValidateGenesis(data GenesisState) error {
	return nil
}
