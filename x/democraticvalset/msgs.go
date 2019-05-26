package democraticvalset

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgUpdateDescription - struct for updating a validator's description
type MsgUpdateDescription struct {
	Description
	ValidatorAddress sdk.ValAddress `json:"address"`
}

func NewMsgUpdateDescription(valAddr sdk.ValAddress, description Description) MsgUpdateDescription {
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

// MsgUpdateDescription - struct for updating a validator's description
type MsgUpdateConsPubKey struct {
	ConsPubKey       sdk.ConsPubKey
	ValidatorAddress sdk.ValAddress `json:"address"`
}

func NewMsgUpdateConsPubKey(valAddr sdk.ValAddress, newConsPubKey sdk.ConsPubKey) MsgUpdateConsPubKey {
	return MsgUpdateConsPubKey{
		ConsPubKey:       newConsPubKey,
		ValidatorAddress: valAddr,
	}
}

//nolint
func (msg MsgUpdateConsPubKey) Route() string { return RouterKey }
func (msg MsgUpdateConsPubKey) Type() string  { return "update_conspubkey" }
func (msg MsgUpdateConsPubKey) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.ValidatorAddress)}
}

// get the bytes for the message signer to sign on
func (msg MsgUpdateConsPubKey) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgUpdateConsPubKey) ValidateBasic() sdk.Error {
	if msg.ValidatorAddress.Empty() {
		return sdk.NewError(DefaultCodespace, CodeInvalidInput, "nil validator address")
	}

	if msg.ConsPubKey.Empty() {
		return ErrInvalidValidatorConsPubKey(DefaultCodespace)
	}

	return nil
}
