package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/schedule module sentinel errors
var (
	ErrInvalidSigner              = sdkerrors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrInvalidStudyPlan           = sdkerrors.Register(ModuleName, 1101, "invalid study plan")
	ErrStudentNotFound            = sdkerrors.Register(ModuleName, 1102, "student not found")
	ErrStudyPlanNotFound          = sdkerrors.Register(ModuleName, 1103, "study plan not found")
	ErrPlannedSemesterNotFound    = sdkerrors.Register(ModuleName, 1104, "planned semester not found")
	ErrSubjectRecommendationNotFound = sdkerrors.Register(ModuleName, 1105, "subject recommendation not found")
	ErrInvalidSemesterNumber      = sdkerrors.Register(ModuleName, 1106, "invalid semester number")
	ErrPrerequisitesNotMet        = sdkerrors.Register(ModuleName, 1107, "prerequisites not met for recommended subject")
	ErrInvalidRecommendationType  = sdkerrors.Register(ModuleName, 1108, "invalid recommendation type")
	ErrInvalidPriority            = sdkerrors.Register(ModuleName, 1109, "invalid priority level")
	ErrInvalidCredits             = sdkerrors.Register(ModuleName, 1110, "invalid credits")
	ErrInvalidWorkloadHours       = sdkerrors.Register(ModuleName, 1111, "invalid workload hours")
	ErrStudyPlanAlreadyExists     = sdkerrors.Register(ModuleName, 1112, "study plan already exists for this student")
	ErrPlannedSemesterConflict    = sdkerrors.Register(ModuleName, 1113, "planned semester conflict")
	ErrInvalidDifficultyLevel     = sdkerrors.Register(ModuleName, 1114, "invalid difficulty level")
	ErrMaxCreditsExceeded         = sdkerrors.Register(ModuleName, 1115, "maximum credits per semester exceeded")
	ErrInvalidScheduleStatus      = sdkerrors.Register(ModuleName, 1116, "invalid schedule status")
	ErrCurriculumNotFound         = sdkerrors.Register(ModuleName, 1117, "curriculum not found")
	ErrInvalidRecommendationScore = sdkerrors.Register(ModuleName, 1118, "invalid recommendation score")
	ErrSubjectAlreadyCompleted    = sdkerrors.Register(ModuleName, 1119, "subject already completed")
	ErrSubjectInProgress          = sdkerrors.Register(ModuleName, 1120, "subject already in progress")
)
