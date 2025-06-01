package types

//"fmt"

// this line is used by starport scaffolding # genesis/types/import

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		// Note: SubjectEquivalenceList and SubjectEquivalenceCount fields
		// need to be defined in the protobuf genesis.proto file
		// For now, only include the params
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Note: Validation for SubjectEquivalenceList would go here
	// when the protobuf fields are properly defined

	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
