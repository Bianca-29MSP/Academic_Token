package types

import (
	"crypto/sha256"
	fmt "fmt"

	sdkerrors "cosmossdk.io/errors"
)

const (
	// ModuleName defines the module name
	ModuleName = "tokendef"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_tokendef"

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName
)

var (
	ParamsKey = []byte("p_tokendef")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

const (
	TokenDefinitionKeyPrefix              = "TokenDefinition/value/"
	TokenDefinitionCountKey               = "TokenDefinition/count/"
	TokenDefinitionBySubjectKeyPrefix     = "TokenDefinition/subject/"
	TokenDefinitionByCourseKeyPrefix      = "TokenDefinition/course/"
	TokenDefinitionByInstitutionKeyPrefix = "TokenDefinition/institution/"
)

// x/tokendef module sentinel errors
var (
	ErrInvalidSigner           = sdkerrors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrTokenDefinitionNotFound = sdkerrors.Register(ModuleName, 1101, "token definition not found")
	ErrInvalidSubjectId        = sdkerrors.Register(ModuleName, 1102, "invalid subject ID")
	ErrInvalidTokenName        = sdkerrors.Register(ModuleName, 1103, "invalid token name")
	ErrInvalidTokenSymbol      = sdkerrors.Register(ModuleName, 1104, "invalid token symbol")
	ErrTokenDefinitionExists   = sdkerrors.Register(ModuleName, 1105, "token definition already exists")
	ErrInvalidContentHash      = sdkerrors.Register(ModuleName, 1106, "invalid content hash")
	ErrInvalidIpfsLink         = sdkerrors.Register(ModuleName, 1107, "invalid IPFS link")
	ErrUnauthorized            = sdkerrors.Register(ModuleName, 1108, "unauthorized")
	ErrInvalidMaxSupply        = sdkerrors.Register(ModuleName, 1109, "invalid max supply")
	ErrInvalidMetadata         = sdkerrors.Register(ModuleName, 1110, "invalid metadata")
)

// TokenDefinitionKey returns the store key to retrieve a TokenDefinition from the index fields
func TokenDefinitionKey(index string) []byte {
	var key []byte
	indexBytes := []byte(index)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)
	return key
}

// TokenDefinitionBySubjectKey returns the store key for subject indexing
func TokenDefinitionBySubjectKey(subjectId, tokenDefIndex string) []byte {
	var key []byte
	subjectIdBytes := []byte(subjectId)
	key = append(key, subjectIdBytes...)
	key = append(key, []byte("/")...)
	indexBytes := []byte(tokenDefIndex)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)
	return key
}

// TokenDefinitionByCourseKey returns the store key for course indexing
func TokenDefinitionByCourseKey(courseId, tokenDefIndex string) []byte {
	var key []byte
	courseIdBytes := []byte(courseId)
	key = append(key, courseIdBytes...)
	key = append(key, []byte("/")...)
	indexBytes := []byte(tokenDefIndex)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)
	return key
}

// TokenDefinitionByInstitutionKey returns the store key for institution indexing
func TokenDefinitionByInstitutionKey(institutionId, tokenDefIndex string) []byte {
	var key []byte
	institutionIdBytes := []byte(institutionId)
	key = append(key, institutionIdBytes...)
	key = append(key, []byte("/")...)
	indexBytes := []byte(tokenDefIndex)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)
	return key
}

// GenerateTokenDefId generates a unique token definition ID
func GenerateTokenDefId(subjectId, tokenName string) string {
	data := fmt.Sprintf("%s:%s", subjectId, tokenName)
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("tokendef_%x", hash[:8])
}

// GenerateContentHash generates content hash for IPFS integration
func GenerateContentHash(metadata *TokenMetadata) string {
	if metadata == nil {
		return ""
	}

	content := fmt.Sprintf("%s:%s", metadata.Description, metadata.ImageUri)
	for _, attr := range metadata.Attributes {
		content += fmt.Sprintf(":%s:%s:%t", attr.TraitType, attr.DisplayType, attr.IsDynamic)
	}

	hash := sha256.Sum256([]byte(content))
	return fmt.Sprintf("%x", hash)
}
