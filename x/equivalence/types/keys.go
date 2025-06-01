package types

const (
	// ModuleName defines the module name
	ModuleName = "equivalence"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_equivalence"
)

var (
	ParamsKey = []byte("p_equivalence")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

const (
	SubjectEquivalenceKeyPrefix = "SubjectEquivalence/value/"
	SubjectEquivalenceCountKey  = "SubjectEquivalence/count/"
)

// SubjectEquivalenceKey returns the store key to retrieve a SubjectEquivalence from the index fields
func SubjectEquivalenceKey(index string) []byte {
	var key []byte

	indexBytes := []byte(index)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)

	return key
}

// SubjectEquivalenceBySourceKey returns the store key prefix for equivalences by source subject
func SubjectEquivalenceBySourceKey(sourceSubjectId string) []byte {
	return []byte("source/" + sourceSubjectId + "/")
}

// SubjectEquivalenceByTargetKey returns the store key prefix for equivalences by target subject
func SubjectEquivalenceByTargetKey(targetSubjectId string) []byte {
	return []byte("target/" + targetSubjectId + "/")
}

// SubjectEquivalenceByInstitutionKey returns the store key prefix for equivalences by institution
func SubjectEquivalenceByInstitutionKey(institutionId string) []byte {
	return []byte("institution/" + institutionId + "/")
}

// SubjectEquivalenceByStatusKey returns the store key prefix for equivalences by status
func SubjectEquivalenceByStatusKey(status string) []byte {
	return []byte("status/" + status + "/")
}

// SubjectEquivalenceByContractKey returns the store key prefix for equivalences by contract
func SubjectEquivalenceByContractKey(contractAddress string) []byte {
	return []byte("contract/" + contractAddress + "/")
}

// SubjectEquivalenceByContractVersionKey returns the store key prefix for equivalences by contract version
func SubjectEquivalenceByContractVersionKey(contractVersion string) []byte {
	return []byte("contract-version/" + contractVersion + "/")
}
