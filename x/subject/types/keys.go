package types

import "fmt"

const (
	// ModuleName defines the module name
	ModuleName = "subject"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_subject"
)

// Store prefixes and keys
var (
	ParamsKey                      = []byte("p_subject")
	SubjectCountKey                = []byte("SubjectCount")
	PrerequisiteGroupCountKey      = []byte("PrereqGroupCount")
	SubjectPrefix                  = []byte("Subject/")
	PrerequisiteGroupPrefix        = []byte("PrereqGroup/")
	SubjectPrerequisiteGroupPrefix = []byte("SubjectPrereq/")
)

// Event types
const (
	EventTypeCreateSubject        = "create_subject"
	EventTypeUpdateSubjectContent = "update_subject_content"
	EventTypeAddPrerequisiteGroup = "add_prerequisite_group"
	EventTypeUpdateParams         = "update_params"
)

// Event attribute keys
const (
	AttributeKeyIndex                         = "index"
	AttributeKeyTitle                         = "title"
	AttributeKeyCode                          = "code"
	AttributeKeyInstitution                   = "institution"
	AttributeKeyCourseId                      = "course_id" // Added course_id attribute
	AttributeKeySubjectId                     = "subject_id"
	AttributeKeyContentHash                   = "content_hash"
	AttributeKeyIpfsLink                      = "ipfs_link"
	AttributeKeyGroupType                     = "group_type"
	AttributeKeyGroupIndex                    = "group_index"
	AttributeKeyIPFSEnabled                   = "ipfs_enabled"
	AttributeKeyIPFSGateway                   = "ipfs_gateway"
	AttributeKeyPrerequisiteValidatorContract = "prerequisite_validator_contract"
	AttributeKeyEquivalenceValidatorContract  = "equivalence_validator_contract"
)

const (
	SubjectByCourseKeyPrefix      = "subject_by_course/"
	SubjectByInstitutionKeyPrefix = "subject_by_institution/"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

// SubjectKey returns the store key to retrieve a Subject by ID
func SubjectKey(id string) []byte {
	return append(SubjectPrefix, []byte(id)...)
}

// PrerequisiteGroupKey returns the key for a prerequisite group
func PrerequisiteGroupKey(id string) []byte {
	return append(PrerequisiteGroupPrefix, []byte(id)...)
}

// SubjectPrerequisiteGroupKey returns the key for a subject-prerequisite group relationship
func SubjectPrerequisiteGroupKey(subjectId, groupId string) []byte {
	key := append(SubjectPrerequisiteGroupPrefix, []byte(subjectId)...)
	return append(key, []byte(groupId)...)
}

// GetSubjectPrerequisiteGroupPrefix returns the prefix for all prerequisite groups of a subject
func GetSubjectPrerequisiteGroupPrefix(subjectId string) []byte {
	return append(SubjectPrerequisiteGroupPrefix, []byte(subjectId)...)
}

// SubjectByCourseKey creates a key for subject by course index
func SubjectByCourseKey(courseId string, subjectIndex string) []byte {
	return []byte(fmt.Sprintf("%s%s/%s", SubjectByCourseKeyPrefix, courseId, subjectIndex))
}

// SubjectByInstitutionKey creates a key for subject by institution index
func SubjectByInstitutionKey(institutionId string, subjectIndex string) []byte {
	return []byte(fmt.Sprintf("%s%s/%s", SubjectByInstitutionKeyPrefix, institutionId, subjectIndex))
}

// GetSubjectByCoursePrefix returns the prefix for subjects by course
func GetSubjectByCoursePrefix(courseId string) []byte {
	return []byte(fmt.Sprintf("%s%s/", SubjectByCourseKeyPrefix, courseId))
}

// GetSubjectByInstitutionPrefix returns the prefix for subjects by institution
func GetSubjectByInstitutionPrefix(institutionId string) []byte {
	return []byte(fmt.Sprintf("%s%s/", SubjectByInstitutionKeyPrefix, institutionId))
}

// GetSubjectPrefix returns the prefix for all subjects
func GetSubjectPrefix() []byte {
	return SubjectPrefix // Fixed: was referencing undefined SubjectKeyPrefix
}
