package poa

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sikkatech/sikkadao/x/gov"
)

// Constants pertaining to a Content object
const (
	ProposalRoute = "poa"

	AddValidatorType = "add_validator"
	UpdateValidatorPowerType = "update_validator"
	RemoveValidatorType = "remove_validator"
)


// Handle all "poa" type gov contents.
func NewGovHandler(keeper Keeper) gov.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch content := msg.(content) {
		case ContentAddValidator:
			return handleContentAddValidator(ctx, keeper, content)
		
		default:
			errMsg := fmt.Sprintf("unrecognized poa governance content type: %T", msg)
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
	return content.String()
}

func (content ContentAddValidator) ProposalRoute() string {
	return ProposalRoute
}

func (content ContentAddValidator) ProposalType() string {
	return AddValidatorType
}

func (content ContentAddValidator) ValidateBasic() sdk.Error {
	if content.valAddress.Empty() {
		return ErrNilValidatorAddr(DefaultCodespace)
	}

	if content.power <= 0 {
		return ErrNonPositivePower(DefaultCodespace)
	}

	_, err := content.valDescription.EnsureLength()
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