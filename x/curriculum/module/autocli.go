package curriculum

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	modulev1 "academictoken/api/academictoken/curriculum"
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
					RpcMethod:      "CreateCurriculumTree",
					Use:            "create-curriculum-tree [course-id] [version] [elective-min] [total-workload-hours]",
					Short:          "Send a create-curriculum-tree tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "courseId"}, {ProtoField: "version"}, {ProtoField: "electiveMin"}, {ProtoField: "totalWorkloadHours"}},
				},
				{
					RpcMethod:      "AddSemesterToCurriculum",
					Use:            "add-semester-to-curriculum [curriculum-index] [semester-number]",
					Short:          "Send a add-semester-to-curriculum tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "curriculumIndex"}, {ProtoField: "semesterNumber"}},
				},
				{
					RpcMethod:      "AddElectiveGroup",
					Use:            "add-elective-group [curriculum-index] [name] [description] [min-subjects-required] [credits-required] [knowledge-area]",
					Short:          "Send a add-elective-group tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "curriculumIndex"}, {ProtoField: "name"}, {ProtoField: "description"}, {ProtoField: "minSubjectsRequired"}, {ProtoField: "creditsRequired"}, {ProtoField: "knowledgeArea"}},
				},
				{
					RpcMethod:      "SetGraduationRequirements",
					Use:            "set-graduation-requirements [curriculum-index] [total-credits-required] [min-gpa] [required-elective-credits] [minimum-time-years] [maximum-time-years]",
					Short:          "Send a set-graduation-requirements tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "curriculumIndex"}, {ProtoField: "totalCreditsRequired"}, {ProtoField: "minGpa"}, {ProtoField: "requiredElectiveCredits"}, {ProtoField: "minimumTimeYears"}, {ProtoField: "maximumTimeYears"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
