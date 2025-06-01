package simulation

import (
	"math/rand"

	"academictoken/x/schedule/keeper"
	"academictoken/x/schedule/types"

	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

const (
	// Operation weights for different message types
	OpWeightMsgCreateStudyPlan            = "op_weight_msg_create_study_plan"
	OpWeightMsgCreateSubjectRecommendation = "op_weight_msg_create_subject_recommendation"
	OpWeightMsgAddPlannedSemester         = "op_weight_msg_add_planned_semester"
	OpWeightMsgUpdateStudyPlanStatus      = "op_weight_msg_update_study_plan_status"

	// Default operation weights
	DefaultWeightMsgCreateStudyPlan            = 100
	DefaultWeightMsgCreateSubjectRecommendation = 100
	DefaultWeightMsgAddPlannedSemester         = 50
	DefaultWeightMsgUpdateStudyPlanStatus      = 50
)

// WeightedOperations returns all the operations from the schedule module with their respective weights
func WeightedOperations(
	appParams simtypes.AppParams, 
	cdc interface{}, 
	ak types.AccountKeeper, 
	bk types.BankKeeper, 
	k keeper.Keeper,
) simulation.WeightedOperations {
	var (
		weightMsgCreateStudyPlan            int
		weightMsgCreateSubjectRecommendation int
		weightMsgAddPlannedSemester         int
		weightMsgUpdateStudyPlanStatus      int
	)

	// Generate weights for each operation
	appParams.GetOrGenerate(OpWeightMsgCreateStudyPlan, &weightMsgCreateStudyPlan, nil,
		func(r *rand.Rand) { weightMsgCreateStudyPlan = DefaultWeightMsgCreateStudyPlan })

	appParams.GetOrGenerate(OpWeightMsgCreateSubjectRecommendation, &weightMsgCreateSubjectRecommendation, nil,
		func(r *rand.Rand) { weightMsgCreateSubjectRecommendation = DefaultWeightMsgCreateSubjectRecommendation })

	appParams.GetOrGenerate(OpWeightMsgAddPlannedSemester, &weightMsgAddPlannedSemester, nil,
		func(r *rand.Rand) { weightMsgAddPlannedSemester = DefaultWeightMsgAddPlannedSemester })

	appParams.GetOrGenerate(OpWeightMsgUpdateStudyPlanStatus, &weightMsgUpdateStudyPlanStatus, nil,
		func(r *rand.Rand) { weightMsgUpdateStudyPlanStatus = DefaultWeightMsgUpdateStudyPlanStatus })

	// Return all weighted operations
	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgCreateStudyPlan,
			SimulateMsgCreateStudyPlan(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgCreateSubjectRecommendation,
			SimulateMsgCreateSubjectRecommendation(ak, bk, k),
		),
		// TODO: Enable these operations after fixing their implementations
		// simulation.NewWeightedOperation(
		// 	weightMsgAddPlannedSemester,
		// 	SimulateMsgAddPlannedSemester(ak, bk, k),
		// ),
		// simulation.NewWeightedOperation(
		// 	weightMsgUpdateStudyPlanStatus,
		// 	SimulateMsgUpdateStudyPlanStatus(ak, bk, k),
		// ),
	}
}
