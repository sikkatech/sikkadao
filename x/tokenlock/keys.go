package tokenlock

import (
	"bytes"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Key for getting a the next available proposalID from the store
var (
	KeyDelimiter = []byte(":")

	PrefixLocks       = []byte("locks")
	PrefixUnlockQueue = []byte("unlocks")
)

func KeyLock(owner sdk.AccAddress, unlockTime time.Duration) []byte {
	return bytes.Join([][]byte{
		owner,
		[]byte(unlockTime.String()),
	}, KeyDelimiter)
}

// Returns the key for a proposalID in the activeProposalQueue
func PrefixUnlockQueueTime(endTime time.Time) []byte {
	return bytes.Join([][]byte{
		sdk.FormatTimeBytes(endTime),
	}, KeyDelimiter)
}

// Returns the key for a proposalID in the activeProposalQueue
func KeyUnlock(unlock TokenUnlock) []byte {
	return bytes.Join([][]byte{
		sdk.FormatTimeBytes(unlock.CompletionTime),
		unlock.Owner,
	}, KeyDelimiter)
}
