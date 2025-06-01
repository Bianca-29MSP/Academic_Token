package schedule

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	modulev1 "academictoken/api/academictoken/schedule"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: modulev1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              modulev1.Msg_ServiceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod:      "CreateSubjectRecommendation",
					Use:            "create-subject-recommendation [student] [recommendation-semester]",
					Short:          "Send a create-subject-recommendation tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "student"}, {ProtoField: "recommendationSemester"}},
				},
				{
					RpcMethod:      "CreateStudyPlan",
					Use:            "create-study-plan [student] [completion-target]",
					Short:          "Send a create-study-plan tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "student"}, {ProtoField: "completionTarget"}},
				},
				{
					RpcMethod:      "AddPlannedSemester",
					Use:            "add-planned-semester [study-plan-id] [semester-code]",
					Short:          "Send a add-planned-semester tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "studyPlanId"}, {ProtoField: "semesterCode"}},
				},
				{
					RpcMethod:      "UpdateStudyPlanStatus",
					Use:            "update-study-plan-status [study-plan-id] [status]",
					Short:          "Send a update-study-plan-status tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "studyPlanId"}, {ProtoField: "status"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
