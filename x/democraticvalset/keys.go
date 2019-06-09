package democraticvalset

const (
	// ModuleName is the name of this module
	ModuleName = "democraticvalset"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// TStoreKey is the string transient store representation
	TStoreKey = "transient_" + ModuleName

	// QuerierRoute is the querier route for this module
	QuerierRoute = ModuleName

	// RouterKey is the msg router key for this module
	RouterKey = ModuleName

	// ProposalKey is the gov proposal router key for this module
	ProposalKey = ModuleName
)
