//nolint
package tokenlock

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = "tokenlock"

	CodeInsufficientCoins sdk.CodeType = 1
)

func ErrInsufficientCoins(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInsufficientCoins, "insufficient coins in lock")
}
