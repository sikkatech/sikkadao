//nolint
package poa

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = "poa"

	CodeNonexistantValidator sdk.CodeType = 1
	CodeInvalidDelegation    sdk.CodeType = 2
	CodeInvalidInput         sdk.CodeType = 3
	CodeNonPositivePower     sdk.CodeType = 4
)

func ErrNonexistantValidator(codespace sdk.CodespaceType, valAddress sdk.ValAddress) sdk.Error {
	return sdk.NewError(codespace, CodeNonexistantValidator, fmt.Sprintf("no validator with with address %s", valAddress))
}

//validator
func ErrNilValidatorAddr(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, "validator address is nil")
}

func ErrDescriptionLength(codespace sdk.CodespaceType, descriptor string, got, max int) sdk.Error {
	msg := fmt.Sprintf("bad description length for %v, got length %v, max is %v", descriptor, got, max)
	return sdk.NewError(codespace, CodeInvalidValidator, msg)
}

func ErrNonPositivePower(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeNonPositivePower, fmt.Sprintf("validator power must be positive"))
}
