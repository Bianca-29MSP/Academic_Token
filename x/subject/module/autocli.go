package subject

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	modulev1 "academictoken/api/academictoken/subject"
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
					RpcMethod:      "CreateSubject",
					Use:            "create-subject [institution] [title] [code] [workload-hours] [credits] [description] [subject-type] [knowledge-area]",
					Short:          "Send a create-subject tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "institution"}, {ProtoField: "title"}, {ProtoField: "code"}, {ProtoField: "workloadHours"}, {ProtoField: "credits"}, {ProtoField: "description"}, {ProtoField: "subjectType"}, {ProtoField: "knowledgeArea"}},
				},
				{
					RpcMethod:      "AddPrerequisiteGroup",
					Use:            "add-prerequisite-group [subject-id] [group-type] [minimum-credits] [minimum-completed-subjects]",
					Short:          "Send a add-prerequisite-group tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "subjectId"}, {ProtoField: "groupType"}, {ProtoField: "minimumCredits"}, {ProtoField: "minimumCompletedSubjects"}},
				},
				{
					RpcMethod:      "UpdateSubjectContent",
					Use:            "update-subject-content [subject-id]",
					Short:          "Send a update-subject-content tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "subjectId"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
