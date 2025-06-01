package types

const (
	// ModuleName defines the module name
	ModuleName = "institution"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_institution"
)

const (
	InstitutionKeyPrefix = "Institution/value/"
	InstitutionCountKey  = "Institution/count/"
)

const (
	// Event types
	EventTypeRegisterInstitution    = "register_institution"
	EventTypeUpdateInstitution      = "update_institution"
	EventTypeAuthorizeInstitution   = "authorize_institution"
	EventTypeUnauthorizeInstitution = "unauthorize_institution"

	// Event attribute keys
	AttributeKeyInstitutionIndex = "institution_index"
	AttributeKeyInstitutionName  = "institution_name"
	AttributeKeyCreator          = "creator"
	AttributeKeyOriginalCreator  = "original_creator"
	AttributeKeyIsAuthorized     = "is_authorized"
)

var (
	ParamsKey = []byte("params")
)

// InstitutionKey returns the store key to retrieve an Institution from the index fields
func InstitutionKey(index string) []byte {
	var key []byte

	key = append(key, KeyPrefix(InstitutionKeyPrefix)...)

	indexBytes := []byte(index)
	key = append(key, indexBytes...)

	return key
}

// KeyPrefix creates a key prefix from a string
func KeyPrefix(p string) []byte {
	return []byte(p)
}
