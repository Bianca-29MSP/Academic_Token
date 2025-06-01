package types

import "encoding/binary"

const (
	// ModuleName defines the module name
	ModuleName = "academicnft"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// QuerierRoute is the querier route for the module
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_" + ModuleName
)

// Store key constants
const (
	// Keys for KV store (using the old string-based format for compatibility)
	ParamsKeyPrefix                    = "Params-value-"
	TokenInstanceCountKeyPrefix        = "TokenInstanceCount-value-"
	SubjectTokenInstanceKeyPrefixStr   = "SubjectTokenInstance-value-"
	StudentTokenIndexKeyPrefixStr      = "StudentTokenIndex-"
	TokenDefInstanceIndexKeyPrefixStr  = "TokenDefInstanceIndex-"
)

// Binary store prefixes (more efficient for new implementations)
var (
	// SubjectTokenInstanceKeyPrefix is the prefix for SubjectTokenInstance store
	SubjectTokenInstanceKeyPrefix = []byte{0x01}

	// StudentTokenIndexKeyPrefix is the prefix for student token indexes
	StudentTokenIndexKeyPrefix = []byte{0x02}

	// TokenDefInstanceIndexKeyPrefix is the prefix for tokendef instance indexes
	TokenDefInstanceIndexKeyPrefix = []byte{0x03}

	// TokenInstanceCountKey is the key for storing token instance count
	TokenInstanceCountKey = []byte{0x04}

	// ParamsKey is the key for storing module parameters
	ParamsKey = []byte{0x05}
)

// SubjectTokenInstanceKey returns the store key for a subject token instance
// Uses binary prefix for efficiency
func SubjectTokenInstanceKey(index string) []byte {
	return append(SubjectTokenInstanceKeyPrefix, []byte(index)...)
}

// SubjectTokenInstanceKeyCompat returns the store key using legacy string format
// Kept for backward compatibility
func SubjectTokenInstanceKeyCompat(index string) []byte {
	var key []byte
	indexBytes := []byte(index)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)
	return key
}

// StudentTokenIndexKey returns the store key for student token index
func StudentTokenIndexKey(studentAddress, tokenInstanceId string) []byte {
	prefix := StudentTokenIndexPrefix(studentAddress)
	return append(prefix, []byte(tokenInstanceId)...)
}

// StudentTokenIndexKeyCompat returns the store key using legacy string format
func StudentTokenIndexKeyCompat(studentAddress, tokenInstanceId string) []byte {
	return append([]byte(StudentTokenIndexKeyPrefixStr), []byte(studentAddress + "/" + tokenInstanceId)...)
}

// StudentTokenIndexPrefix returns the prefix for a student's token indexes
func StudentTokenIndexPrefix(studentAddress string) []byte {
	prefix := append(StudentTokenIndexKeyPrefix, []byte(studentAddress)...)
	return append(prefix, []byte("/")...)
}

// StudentTokenIndexPrefixCompat returns the prefix using legacy string format
func StudentTokenIndexPrefixCompat(studentAddress string) []byte {
	return append([]byte(StudentTokenIndexKeyPrefixStr), []byte(studentAddress + "/")...)
}

// TokenDefInstanceIndexKey returns the store key for tokendef instance index
func TokenDefInstanceIndexKey(tokenDefId, tokenInstanceId string) []byte {
	prefix := TokenDefInstanceIndexPrefix(tokenDefId)
	return append(prefix, []byte(tokenInstanceId)...)
}

// TokenDefInstanceIndexKeyCompat returns the store key using legacy string format
func TokenDefInstanceIndexKeyCompat(tokenDefId, tokenInstanceId string) []byte {
	return append([]byte(TokenDefInstanceIndexKeyPrefixStr), []byte(tokenDefId + "/" + tokenInstanceId)...)
}

// TokenDefInstanceIndexPrefix returns the prefix for a TokenDef's instance indexes
func TokenDefInstanceIndexPrefix(tokenDefId string) []byte {
	prefix := append(TokenDefInstanceIndexKeyPrefix, []byte(tokenDefId)...)
	return append(prefix, []byte("/")...)
}

// TokenDefInstanceIndexPrefixCompat returns the prefix using legacy string format
func TokenDefInstanceIndexPrefixCompat(tokenDefId string) []byte {
	return append([]byte(TokenDefInstanceIndexKeyPrefixStr), []byte(tokenDefId + "/")...)
}

// GetSubjectTokenInstanceIDBytes returns the byte representation of the ID
func GetSubjectTokenInstanceIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

// GetSubjectTokenInstanceIDFromBytes returns ID in uint64 format from a byte array
func GetSubjectTokenInstanceIDFromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}
