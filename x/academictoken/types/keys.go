package types

const (
	// ModuleName defines the module name
	ModuleName = "academictoken"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_academictoken"
)

var (
	ParamsKey = []byte("p_academictoken")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
