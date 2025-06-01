package degree

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	modulev1 "academictoken/api/academictoken/degree"
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
					RpcMethod:      "RequestDegree",
					Use:            "request-degree-verification [student] [course-id]",
					Short:          "Send a request-degree-verification tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "student"}, {ProtoField: "courseId"}},
				},
				{
					RpcMethod:      "IssueDegree",
					Use:            "issue-degree [student] [institution] [course-id] [final-grade]",
					Short:          "Send a issue-degree tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "student"}, {ProtoField: "institution"}, {ProtoField: "courseId"}, {ProtoField: "finalGrade"}},
				},
				{
					RpcMethod:      "VerifyDegree",
					Use:            "verify-degree [degree-id]",
					Short:          "Send a verify-degree tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "degreeId"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
