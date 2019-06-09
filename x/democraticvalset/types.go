package democraticvalset

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

// Simple Validator for a Proof of Authority System
type Validator struct {
	OperatorAddress sdk.ValAddress `json:"operator_address"` // address of the validator's operator; bech encoded in JSON
	ConsPubKey      crypto.PubKey  `json:"consensus_pubkey"` // the consensus public key of the validator; bech encoded in JSON
	Description     Description    `json:"description"`      // description terms for the validator
	Power           int64          `json:"power"`
}

func (val Validator) Equals(val2 Validator) bool {
	if val.OperatorAddress.Equals(val2.OperatorAddress) &&
		val.ConsPubKey.Equals(val2.ConsPubKey) &&
		val.Description.Equals(val2.Description) &&
		val.Power == val2.Power {
		return true
	}
	return false
}

func EmptyValidator() Validator {
	return Validator{}
}

func (val Validator) IsEmpty() bool {
	return val.Equals(EmptyValidator())
}

// nolint
const (
	// TODO: Why can't we just have one string description which can be JSON by convention
	MaxMonikerLength  = 70
	MaxIdentityLength = 3000
	MaxWebsiteLength  = 140
	MaxDetailsLength  = 280
)

// constant used in flags to indicate that description field should not be updated
const DoNotModifyDesc = "[do-not-modify]"

// Description - description fields for a validator
type Description struct {
	Moniker  string `json:"moniker"`  // name
	Identity string `json:"identity"` // optional identity signature (ex. UPort or Keybase)
	Website  string `json:"website"`  // optional website link
	Details  string `json:"details"`  // optional details
}

// NewDescription returns a new Description with the provided values.
func NewDescription(moniker, identity, website, details string) Description {
	return Description{
		Moniker:  moniker,
		Identity: identity,
		Website:  website,
		Details:  details,
	}
}

// UpdateDescription updates the fields of a given description. An error is
// returned if the resulting description contains an invalid length.
func (d Description) UpdateDescription(d2 Description) (Description, sdk.Error) {
	if d2.Moniker == DoNotModifyDesc {
		d2.Moniker = d.Moniker
	}
	if d2.Identity == DoNotModifyDesc {
		d2.Identity = d.Identity
	}
	if d2.Website == DoNotModifyDesc {
		d2.Website = d.Website
	}
	if d2.Details == DoNotModifyDesc {
		d2.Details = d.Details
	}

	return Description{
		Moniker:  d2.Moniker,
		Identity: d2.Identity,
		Website:  d2.Website,
		Details:  d2.Details,
	}.EnsureLength()
}

// EnsureLength ensures the length of a validator's description.
func (d Description) EnsureLength() (Description, sdk.Error) {
	if len(d.Moniker) > MaxMonikerLength {
		return d, ErrDescriptionLength(DefaultCodespace, "moniker", len(d.Moniker), MaxMonikerLength)
	}
	if len(d.Identity) > MaxIdentityLength {
		return d, ErrDescriptionLength(DefaultCodespace, "identity", len(d.Identity), MaxIdentityLength)
	}
	if len(d.Website) > MaxWebsiteLength {
		return d, ErrDescriptionLength(DefaultCodespace, "website", len(d.Website), MaxWebsiteLength)
	}
	if len(d.Details) > MaxDetailsLength {
		return d, ErrDescriptionLength(DefaultCodespace, "details", len(d.Details), MaxDetailsLength)
	}

	return d, nil
}

func (d Description) Equals(d2 Description) bool {
	if d.Moniker.Equals(d2.Moniker) &&
		d.Identity.Equals(d2.Identity) &&
		d.Website.Equals(d2.Website) &&
		d.Details.Equals(d2.Details) {
		return true
	}
	return false
}
