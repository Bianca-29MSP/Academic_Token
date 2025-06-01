package types

import "encoding/binary"

const (
	// ModuleName defines the module name
	ModuleName = "course"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_course"
)

var (
	ParamsKey = []byte("p_course")
)

const (
	CourseKeyPrefix = "Course/value/"
	CourseCountKey  = "Course/count/"
)

// Module event types
const (
	EventTypeCourseCreated = "course_created"
	EventTypeCourseUpdated = "course_updated"
)

// Module event attribute keys
const (
	AttributeKeyCreator     = "creator"
	AttributeKeyCourseIndex = "course_index"
	AttributeKeyInstitution = "institution"
	AttributeKeyCourseName  = "course_name"
)

// CourseKey returns the store key to retrieve a Course from the index fields
func CourseKey(index string) []byte {
	var key []byte
	key = append(key, KeyPrefix(CourseKeyPrefix)...) // ✅ ADICIONADO: usa o prefixo
	indexBytes := []byte(index)
	key = append(key, indexBytes...)
	return key
}

// KeyPrefix returns key prefix for module store
func KeyPrefix(p string) []byte {
	return []byte(p)
}

// GetCourseIDBytes returns the byte representation of the ID
func GetCourseIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

// GetCourseIDFromBytes returns ID in uint64 format from a byte array
func GetCourseIDFromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}
