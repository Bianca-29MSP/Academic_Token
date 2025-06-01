package types

import (
	errorsmod "cosmossdk.io/errors"
)

// x/student module sentinel errors
var (
	ErrInvalidInstitution   = errorsmod.Register(ModuleName, 2060, "invalid institution")
	ErrInvalidCourse        = errorsmod.Register(ModuleName, 2061, "invalid course")
	ErrInvalidSubject       = errorsmod.Register(ModuleName, 2062, "invalid subject")
	ErrInvalidPacketTimeout = errorsmod.Register(ModuleName, 1500, "invalid packet timeout")
	ErrInvalidVersion       = errorsmod.Register(ModuleName, 1501, "invalid version")
	ErrInvalidSigner        = errorsmod.Register(ModuleName, 1502, "invalid signer")

	// Student-specific errors
	ErrStudentNotFound      = errorsmod.Register(ModuleName, 2001, "student not found")
	ErrStudentAlreadyExists = errorsmod.Register(ModuleName, 2002, "student already exists")

	// Enrollment errors
	ErrEnrollmentNotFound      = errorsmod.Register(ModuleName, 2010, "enrollment not found")
	ErrEnrollmentAlreadyExists = errorsmod.Register(ModuleName, 2011, "enrollment already exists")
	ErrInvalidEnrollmentStatus = errorsmod.Register(ModuleName, 2012, "invalid enrollment status")

	// Academic tree errors
	ErrAcademicTreeNotFound      = errorsmod.Register(ModuleName, 2020, "academic tree not found")
	ErrAcademicTreeAlreadyExists = errorsmod.Register(ModuleName, 2021, "academic tree already exists")

	// Curriculum errors
	ErrCurriculumNotFound = errorsmod.Register(ModuleName, 2030, "curriculum not found")

	// Subject enrollment errors
	ErrSubjectEnrollmentNotAllowed = errorsmod.Register(ModuleName, 2040, "subject enrollment not allowed")
	ErrPrerequisitesNotMet         = errorsmod.Register(ModuleName, 2041, "prerequisites not met")
	ErrSubjectNotAvailable         = errorsmod.Register(ModuleName, 2042, "subject not available")

	// Authorization errors
	ErrNotAuthorized           = errorsmod.Register(ModuleName, 2050, "not authorized")
	ErrInsufficientPermissions = errorsmod.Register(ModuleName, 2051, "insufficient permissions")

	ErrInvalidBasicValidation  = errorsmod.Register(ModuleName, 1101, "invalid basic validation")
	ErrStudyPlanNotFound       = errorsmod.Register(ModuleName, 1102, "study plan not found")
	ErrStudyPlanAlreadyExists  = errorsmod.Register(ModuleName, 1103, "study plan already exists")
	ErrMaxCreditsExceeded      = errorsmod.Register(ModuleName, 1106, "maximum credits per semester exceeded")
	ErrPlannedSemesterConflict = errorsmod.Register(ModuleName, 1107, "planned semester conflict")
	ErrSubjectAlreadyCompleted = errorsmod.Register(ModuleName, 1108, "subject already completed")
	ErrSubjectInProgress       = errorsmod.Register(ModuleName, 1109, "subject already in progress")
)
