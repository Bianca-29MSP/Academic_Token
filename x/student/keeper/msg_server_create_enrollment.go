package keeper

import (
	"context"
	"fmt"
	"time"

	"academictoken/x/student/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateEnrollment(goCtx context.Context, msg *types.MsgCreateEnrollment) (*types.MsgCreateEnrollmentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate that student exists
	student, found := k.Keeper.getStudentByIndex(ctx, msg.Student)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrStudentNotFound, "student with ID %s not found", msg.Student)
	}

	// Validate that institution exists
	_, found = k.Keeper.institutionKeeper.GetInstitution(ctx, msg.Institution)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrInvalidInstitution, "institution with ID %s not found", msg.Institution)
	}

	// Validate that course exists
	_, found = k.Keeper.courseKeeper.GetCourse(ctx, msg.CourseId)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrInvalidCourse, "course with ID %s not found", msg.CourseId)
	}

	// Check permissions - only the student or institution admin can create enrollment
	if !k.canCreateEnrollment(ctx, msg.Creator, student, msg.Institution) {
		return nil, errorsmod.Wrapf(sdkerrors.ErrUnauthorized,
			"address %s is not authorized to create enrollment for student %s",
			msg.Creator, msg.Student)
	}

	// Check if student is already enrolled in this course at this institution
	existingEnrollments := k.Keeper.getEnrollmentsByStudentId(ctx, msg.Student)
	for _, enrollment := range existingEnrollments {
		if enrollment.Institution == msg.Institution && enrollment.CourseId == msg.CourseId {
			// Check if enrollment is active or pending
			if enrollment.Status == types.EnrollmentStatusActive || enrollment.Status == types.EnrollmentStatusPending {
				return nil, errorsmod.Wrapf(types.ErrEnrollmentAlreadyExists,
					"student %s is already enrolled in course %s at institution %s",
					msg.Student, msg.CourseId, msg.Institution)
			}
		}
	}

	// Create new enrollment
	enrollment := types.StudentEnrollment{
		Student:        msg.Student,
		Institution:    msg.Institution,
		CourseId:       msg.CourseId,
		EnrollmentDate: ctx.BlockTime().Format(time.RFC3339),
		Status:         types.EnrollmentStatusPending, // Start as pending, needs approval
		AcademicTreeId: "",                            // Will be set when enrollment is activated
	}

	// Save enrollment
	enrollmentId, err := k.Keeper.AppendStudentEnrollment(ctx, enrollment)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to save enrollment")
	}

	// Get the saved enrollment to get the generated index
	savedEnrollment, found := k.Keeper.getStudentEnrollment(ctx, fmt.Sprintf("%d", enrollmentId))
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrNotFound, "failed to retrieve saved enrollment")
	}

	// Update student's enrollment list
	student.EnrollmentIds = append(student.EnrollmentIds, savedEnrollment.Index)
	k.Keeper.setStudent(ctx, student)

	// Emit enrollment creation event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateEnrollment,
			sdk.NewAttribute(types.AttributeKeyEnrollmentId, savedEnrollment.Index),
			sdk.NewAttribute(types.AttributeKeyStudent, msg.Student),
			sdk.NewAttribute(types.AttributeKeyInstitution, msg.Institution),
			sdk.NewAttribute(types.AttributeKeyCourseId, msg.CourseId),
			sdk.NewAttribute(types.AttributeKeyEnrollmentDate, enrollment.EnrollmentDate),
			sdk.NewAttribute(types.AttributeKeyStatus, enrollment.Status),
			sdk.NewAttribute("created_by", msg.Creator),
		),
	)

	// Log successful enrollment creation
	k.Logger().Info("Enrollment created successfully",
		"enrollment_id", savedEnrollment.Index,
		"student", msg.Student,
		"institution", msg.Institution,
		"course", msg.CourseId,
	)

	return &types.MsgCreateEnrollmentResponse{}, nil
}

// canCreateEnrollment checks if the creator has permission to create an enrollment
func (k msgServer) canCreateEnrollment(ctx sdk.Context, creator string, student types.Student, institutionId string) bool {
	_ = ctx           // Suppress unused parameter warning
	_ = institutionId // Suppress unused parameter warning

	// Student can enroll themselves
	if creator == student.Address {
		return true
	}

	// TODO: Institution admins can enroll students (implement when institution module is ready)
	// For now, allow any address to create enrollments
	return true
}
