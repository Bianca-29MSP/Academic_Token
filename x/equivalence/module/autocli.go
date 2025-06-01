package equivalence

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	modulev1 "academictoken/api/academictoken/equivalence"
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
					RpcMethod:      "RequestEquivalence",
					Use:            "request-equivalence [source-subject-id] [target-institution] [target-subject-id]",
					Short:          "Send a request-equivalence tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "sourceSubjectId"}, {ProtoField: "targetInstitution"}, {ProtoField: "targetSubjectId"}},
				},
				{
					RpcMethod:      "ApproveEquivalence",
					Use:            "approve-equivalence [equivalence-id] [approval-count] [equivalence-percent]",
					Short:          "Send a approve-equivalence tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "equivalenceId"}, {ProtoField: "approvalCount"}, {ProtoField: "equivalencePercent"}},
				},
				{
					RpcMethod:      "RejectEquivalence",
					Use:            "reject-equivalence [equivalence-id]",
					Short:          "Send a reject-equivalence tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "equivalenceId"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
