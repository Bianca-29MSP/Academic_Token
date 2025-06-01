package student

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	modulev1 "academictoken/api/academictoken/student"
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
			EnhanceCustomCommand: false, // Disable custom command enhancement
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Use:       "update-params",
					Short:     "Update student module parameters with all contract addresses",
					Skip:      false, // Force not to skip
					FlagOptions: map[string]*autocliv1.FlagOptions{
						"params.ipfs_gateway": {
							Name:      "ipfs-gateway",
							Shorthand: "",
							Usage:     "IPFS gateway URL",
						},
						"params.ipfs_enabled": {
							Name:      "ipfs-enabled",
							Shorthand: "",
							Usage:     "Enable IPFS integration",
						},
						"params.admin": {
							Name:      "admin",
							Shorthand: "",
							Usage:     "Admin address",
						},
						"params.prerequisites_contract_addr": {
							Name:      "prerequisites-contract",
							Shorthand: "",
							Usage:     "Prerequisites contract address",
						},
						"params.equivalence_contract_addr": {
							Name:      "equivalence-contract",
							Shorthand: "",
							Usage:     "Equivalence contract address",
						},
						"params.academic_progress_contract_addr": {
							Name:      "academic-progress-contract",
							Shorthand: "",
							Usage:     "Academic progress contract address",
						},
						"params.degree_contract_addr": {
							Name:      "degree-contract",
							Shorthand: "",
							Usage:     "Degree contract address",
						},
						"params.nft_minting_contract_addr": {
							Name:      "nft-minting-contract",
							Shorthand: "",
							Usage:     "NFT minting contract address",
						},
					},
				},
				{
					RpcMethod:      "RegisterStudent",
					Use:            "register-student [name] [address]",
					Short:          "Send a register-student tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "name"}, {ProtoField: "address"}},
				},
				{
					RpcMethod:      "CreateEnrollment",
					Use:            "create-enrollment [student] [institution] [course-id]",
					Short:          "Send a create-enrollment tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "student"}, {ProtoField: "institution"}, {ProtoField: "courseId"}},
				},
				{
					RpcMethod:      "UpdateEnrollmentStatus",
					Use:            "update-enrollment-status [enrollment-id] [status]",
					Short:          "Send a update-enrollment-status tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "enrollmentId"}, {ProtoField: "status"}},
				},
				{
					RpcMethod:      "RequestSubjectEnrollment",
					Use:            "request-subject-enrollment [student] [subject-id]",
					Short:          "Send a request-subject-enrollment tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "student"}, {ProtoField: "subjectId"}},
				},
				{
					RpcMethod:      "UpdateAcademicTree",
					Use:            "update-academic-tree [student-id]",
					Short:          "Send a update-academic-tree tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "studentId"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
