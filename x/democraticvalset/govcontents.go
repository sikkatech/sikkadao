package democraticvalset

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sikkatech/sikkadao/x/gov"
)

// Constants pertaining to a Content object
const (
	AddValidatorType = "add_validator"
	UpdateValidatorPowerType = "update_validator"
	RemoveValidatorType = "remove_validator"
)


// Handle all "democraticvalset" type gov contents.
func NewGovHandler(keeper Keeper) gov.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch content := msg.(content) {
		case ContentAddValidator:
			return handleContentAddValidator(ctx, keeper, content)
		case ContentUpdateValidatorPower:
			return handleContentUpdateValidatorPower(ctx, keeper, content)
		
		default:
			errMsg := fmt.Sprintf("unrecognized %s governance content type: %T", ProposalRoute, msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

type ContentAddValidator struct {
	OperatorAddress     sdk.ValAddress
	ConsPubKey crypto.PubKey
	Description Description
	Power          int64
}

func (content ContentAddValidator) GetTitle() string {
	return fmt.Sprintf("Add New Validator %s with voting power %d", content.Description.Moniker, content.Power)
}

func (content ContentAddValidator) GetDescription() string {
	return fmt.Sprintf("Add %s (%s) as a new validator voting power %d.\nTheir website is %s and their keybase identity is %s.\nTheir self described details are:\n%s",
	content.Description.Moniker, content.OperatorAddress, content.Power, content.Description.Website, content.Description.Identity, content.Description.Details)
}

func (content ContentAddValidator) ProposalRoute() string {
	return ProposalRoute
}

func (content ContentAddValidator) ProposalType() string {
	return AddValidatorType
}

func (content ContentAddValidator) ValidateBasic() sdk.Error {
	if content.OperatorAddress.Empty() {
		return ErrNilValidatorAddr(DefaultCodespace)
	}

	if content.Power <= 0 {
		return ErrNonPositivePower(DefaultCodespace)
	}

	_, err := content.Description.EnsureLength()
	return err
}

func (content ContentAddValidator) String() string {
	return string(ModuleCdc.MustMarshalJSON(content))
}

func handleContentAddValidator(ctx sdk.Context, keeper Keeper,content AddValidatorContent) sdk.Error {
	keeper.SetValidator(ctx, Validator{
		OperatorAddress: content.OperatorAddress,
		ConsPubKey: content.ConsPubKey,
		Description: content.Description,
		Power: content.Power
	})
}

type ContentUpdateValidatorPower struct {
	OperatorAddress     sdk.ValAddress
	Power         	    int64
	Justification string
}

func (content ContentUpdateValidatorPower) GetTitle() string {
	return fmt.Sprintf("Update Validator %s to voting power %d", content.OperatorAddress, content.Power)
}

func (content ContentUpdateValidatorPower) GetDescription() string {
	return fmt.Sprintf("Update the validator whose address is %s to voting power %d.\nThe proposed justification is:\n%s",
	content.OperatorAddress, content.Power, content.Justification)
}

func (content ContentUpdateValidatorPower) ProposalRoute() string {
	return ProposalRoute
}

func (content ContentUpdateValidatorPower) ProposalType() string {
	return UpdateValidatorPowerType
}

func (content ContentUpdateValidatorPower) ValidateBasic() sdk.Error {
	if content.OperatorAddress.Empty() {
		return ErrNilValidatorAddr(DefaultCodespace)
	}

	if content.Power <= 0 {
		return ErrNonPositivePower(DefaultCodespace)
	}

	return nil
}

func (content ContentAddValidator) String() string {
	return string(ModuleCdc.MustMarshalJSON(content))
}

func handleContentUpdateValidator(ctx sdk.Context, keeper Keeper,content AddValidatorContent) sdk.Error {
	keeper.UpdateValidatorPower(ctx, content.OperatorAddress, content.Power)
}


type ContentUpdateValidatorPower struct {
	OperatorAddress     sdk.ValAddress
	Justification string
}

func (content ContentUpdateValidatorPower) GetTitle() string {
	return fmt.Sprintf("Remove Validator %s", content.OperatorAddress)
}

func (content ContentUpdateValidatorPower) GetDescription() string {
	return fmt.Sprintf("Remove the validator whose address is %s.\nThe proposed justification is:\n%s",
	content.OperatorAddress, content.Power, content.Justification)
}

func (content ContentUpdateValidatorPower) ProposalRoute() string {
	return ProposalRoute
}

func (content ContentUpdateValidatorPower) ProposalType() string {
	return RemoveValidatorType
}

func (content ContentUpdateValidatorPower) ValidateBasic() sdk.Error {
	if content.OperatorAddress.Empty() {
		return ErrNilValidatorAddr(DefaultCodespace)
	}

	return nil
}

func (content ContentAddValidator) String() string {
	return string(ModuleCdc.MustMarshalJSON(content))
}

func handleContentUpdateValidator(ctx sdk.Context, keeper Keeper,content AddValidatorContent) sdk.Error {
	keeper.RemoveValidator(ctx, content.OperatorAddress, content.Power)
}	