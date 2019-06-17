package tags

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Governance tags
const (
	ActionTokenUnlockStarted   = "unlock-started"
	ActionTokenUnlockCompleted = "unlock-completed"
	TxCategory                 = "tokenlock"

	UnlockTime = "unlock-time"
)

// SDK tag aliases
var (
	Action   = sdk.TagAction
	Category = sdk.TagCategory
	Sender   = sdk.TagSender
)
