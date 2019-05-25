package poa

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgUpdateDescription - struct for updating a validator's description
type MsgUpdateDescription struct {
	Description
	ValidatorAddress sdk.ValAddress `json:"address"`
}

func NewMsgUpdateDescription(valAddr sdk.ValAddress, description Description, newRate *sdk.Dec, newMinSelfDelegation *sdk.Int) MsgEditValidator {
	return MsgUpdateDescription{
		Description:      description,
		ValidatorAddress: valAddr,
	}
}

//nolint
func (msg MsgUpdateDescription) Route() string { return RouterKey }
func (msg MsgUpdateDescription) Type() string  { return "update_description" }
func (msg MsgUpdateDescription) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.ValidatorAddress)}
}

// get the bytes for the message signer to sign on
func (msg MsgUpdateDescription) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgUpdateDescription) ValidateBasic() sdk.Error {
	if msg.ValidatorAddress.Empty() {
		return sdk.NewError(DefaultCodespace, CodeInvalidInput, "nil validator address")
	}

	if msg.Description == (Description{}) {
		return sdk.NewError(DefaultCodespace, CodeInvalidInput, "transaction must include some information to modify")
	}

	return nil
}
