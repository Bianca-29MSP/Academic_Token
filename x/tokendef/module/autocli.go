package tokendef

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	modulev1 "academictoken/api/academictoken/tokendef"
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
					RpcMethod:      "CreateSubjectToken",
					Use:            "create-subject-token [subject-id] [name] [code]",
					Short:          "Send a create-subject-token tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "subjectId"}, {ProtoField: "name"}, {ProtoField: "code"}},
				},
				{
					RpcMethod:      "UpdateSubjectToken",
					Use:            "update-subject-token [token-id] [name] [code] [prerequisites]",
					Short:          "Send a update-subject-token tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "tokenId"}, {ProtoField: "name"}, {ProtoField: "code"}, {ProtoField: "prerequisites"}},
				},
				{
					RpcMethod:      "CreateTokenDefinition",
					Use:            "create-token-definition [subject-id] [name] [code]",
					Short:          "Send a create-token-definition tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "subjectId"}, {ProtoField: "name"}, {ProtoField: "code"}},
				},
				{
					RpcMethod:      "UpdateTokenDefinition",
					Use:            "update-token-definition [token-id] [name] [code]",
					Short:          "Send a update-token-definition tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "tokenId"}, {ProtoField: "name"}, {ProtoField: "code"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
