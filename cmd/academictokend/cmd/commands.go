package cmd

import (
	"context"
	"errors"
	"io"
	"strconv"

	academicnfttypes "academictoken/x/academicnft/types"
	coursetypes "academictoken/x/course/types"
	curriculumtypes "academictoken/x/curriculum/types"
	degreetypes "academictoken/x/degree/types"
	equivalencetypes "academictoken/x/equivalence/types"
	scheduletypes "academictoken/x/schedule/types"
	studenttypes "academictoken/x/student/types"
	subjecttypes "academictoken/x/subject/types"
	tokendefmoduletypes "academictoken/x/tokendef/types"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"

	"cosmossdk.io/log"
	confixcmd "cosmossdk.io/tools/confix/cmd"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/debug"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/pruning"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/client/snapshot"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"academictoken/app"
	"academictoken/x/institution/types"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmcli "github.com/CosmWasm/wasmd/x/wasm/client/cli"
)

func initRootCmd(
	rootCmd *cobra.Command,
	txConfig client.TxConfig,
	basicManager module.BasicManager,
) {
	rootCmd.AddCommand(
		genutilcli.InitCmd(basicManager, app.DefaultNodeHome),
		NewInPlaceTestnetCmd(addModuleInitFlags),
		NewTestnetMultiNodeCmd(basicManager, banktypes.GenesisBalancesIterator{}),
		debug.Cmd(),
		confixcmd.ConfigCommand(),
		pruning.Cmd(newApp, app.DefaultNodeHome),
		snapshot.Cmd(newApp),
	)

	server.AddCommands(rootCmd, app.DefaultNodeHome, newApp, appExport, addModuleInitFlags)

	// add keybase, auxiliary RPC, query, genesis, and tx child commands
	rootCmd.AddCommand(
		server.StatusCommand(),
		genesisCommand(txConfig, basicManager),
		queryCommand(),
		txCommand(),
		keys.Commands(),
	)
	wasmcli.ExtendUnsafeResetAllCmd(rootCmd)
}

func addModuleInitFlags(startCmd *cobra.Command) {
	crisis.AddModuleInitFlags(startCmd)
	wasm.AddModuleInitFlags(startCmd)
}

// genesisCommand builds genesis-related `academictokend genesis` command. Users may provide application specific commands as a parameter
func genesisCommand(txConfig client.TxConfig, basicManager module.BasicManager, cmds ...*cobra.Command) *cobra.Command {
	cmd := genutilcli.Commands(txConfig, basicManager, app.DefaultNodeHome)

	for _, subCmd := range cmds {
		cmd.AddCommand(subCmd)
	}
	return cmd
}

func queryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "query",
		Aliases:                    []string{"q"},
		Short:                      "Querying subcommands",
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		rpc.QueryEventForTxCmd(),
		rpc.ValidatorCommand(),
		server.QueryBlockCmd(),
		authcmd.QueryTxsByEventsCmd(),
		server.QueryBlocksCmd(),
		authcmd.QueryTxCmd(),
		server.QueryBlockResultsCmd(),
		getCourseQueryCmd(),
		getSubjectQueryCmd(),
		getCurriculumQueryCmd(),
		getTokendefQueryCmd(),
		getAcademicNftQueryCmd(),
		getStudentQueryCmd(),
		getDegreeQueryCmd(),
		getEquivalenceQueryCmd(),
		getScheduleQueryCmd(),
	)

	// Add institution module commands
	cmd.AddCommand(
		getInstitutionQueryCmd(),
	)

	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

func txCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "tx",
		Short:                      "Transactions subcommands",
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		authcmd.GetSignCommand(),
		authcmd.GetSignBatchCommand(),
		authcmd.GetMultiSignCommand(),
		authcmd.GetMultiSignBatchCmd(),
		authcmd.GetValidateSignaturesCommand(),
		flags.LineBreak,
		authcmd.GetBroadcastCommand(),
		authcmd.GetEncodeCommand(),
		authcmd.GetDecodeCommand(),
		authcmd.GetSimulateCmd(),
		// Add module-specific transaction commands
		getInstitutionTxCmd(),
		getCourseTxCmd(),
		getSubjectTxCmd(),
		getCurriculumTxCmd(),
		getTokendefTxCmd(),
		getAcademicNftTxCmd(),
		getStudentTxCmd(),
		getDegreeTxCmd(),
		getEquivalenceTxCmd(),
		getScheduleTxCmd(),
	)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

// ============================================================================
// INSTITUTION MODULE COMMANDS
// ============================================================================

// getInstitutionQueryCmd creates the institution query commands
func getInstitutionQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "institution",
		Short:                      "Querying commands for the institution module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		cmdQueryInstitutionParams(),
		cmdQueryListInstitutions(),           // Renamed from cmdQueryInstitutionAll
		cmdQueryGetInstitution(),             // Renamed from cmdQueryInstitution
		cmdQueryListAuthorizedInstitutions(), // Renamed from cmdQueryAuthorizedInstitutions
		cmdQueryGetInstitutionCount(),        // Renamed from cmdQueryInstitutionCount
	)

	return cmd
}

// cmdQueryInstitutionParams implements the params query command
func cmdQueryInstitutionParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Shows the parameters of the module",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Params(context.Background(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// cmdQueryListInstitutions implements the list-institutions query command
func cmdQueryListInstitutions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-institutions",
		Short: "List all institutions",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			params := &types.QueryAllInstitutionRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.InstitutionAll(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// cmdQueryGetInstitution implements the get-institution query command
func cmdQueryGetInstitution() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-institution [index]",
		Short: "Show a specific institution by index",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			params := &types.QueryGetInstitutionRequest{
				Index: args[0],
			}

			res, err := queryClient.Institution(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// cmdQueryListAuthorizedInstitutions implements the list-authorized-institutions query command
func cmdQueryListAuthorizedInstitutions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-authorized-institutions",
		Short: "List all authorized institutions",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			params := &types.QueryAuthorizedInstitutionsRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.AuthorizedInstitutions(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// cmdQueryGetInstitutionCount implements the get-institution-count query command
func cmdQueryGetInstitutionCount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-institution-count",
		Short: "Get the total count of institutions",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.InstitutionCount(context.Background(), &types.QueryInstitutionCountRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// Add institution transaction commands
func getInstitutionTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "institution",
		Short:                      "Institution transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		cmdTxRegisterInstitution(),
		cmdTxUpdateInstitution(),
	)

	return cmd
}

// Institution transaction command implementations
func cmdTxRegisterInstitution() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-institution [name] [address]",
		Short: "Register a new institution",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			name := args[0]
			address := args[1]

			msg := &types.MsgRegisterInstitution{
				Creator: clientCtx.GetFromAddress().String(),
				Name:    name,
				Address: address,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func cmdTxUpdateInstitution() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-institution [index] [name] [address] [is-authorized]",
		Short: "Update an existing institution",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			index := args[0]
			name := args[1]
			address := args[2]
			isAuthorized := args[3]

			msg := &types.MsgUpdateInstitution{
				Creator:      clientCtx.GetFromAddress().String(),
				Index:        index,
				Name:         name,
				Address:      address,
				IsAuthorized: isAuthorized,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// ============================================================================
// COURSE MODULE COMMANDS
// ============================================================================

// getCourseQueryCmd creates the course query commands
func getCourseQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "course",
		Short:                      "Querying commands for the course module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		cmdQueryCourseParams(),
		cmdQueryListCourses(),              // Renamed from cmdQueryCourseAll
		cmdQueryGetCourse(),                // Renamed from cmdQueryCourse
		cmdQueryListCoursesByInstitution(), // Renamed from cmdQueryCoursesByInstitution
		cmdQueryGetCourseCount(),           // Renamed from cmdQueryCourseCount
	)

	return cmd
}

// cmdQueryCourseParams implements the params query command
func cmdQueryCourseParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Shows the parameters of the course module",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := coursetypes.NewQueryClient(clientCtx)
			res, err := queryClient.Params(context.Background(), &coursetypes.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// cmdQueryListCourses implements the list-courses query command
func cmdQueryListCourses() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-courses",
		Short: "List all courses",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := coursetypes.NewQueryClient(clientCtx)
			params := &coursetypes.QueryAllCourseRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.CourseAll(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// cmdQueryGetCourse implements the get-course query command
func cmdQueryGetCourse() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-course [index]",
		Short: "Show a specific course by index",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := coursetypes.NewQueryClient(clientCtx)
			params := &coursetypes.QueryGetCourseRequest{
				Index: args[0],
			}

			res, err := queryClient.Course(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// cmdQueryListCoursesByInstitution implements the list-courses-by-institution query command
func cmdQueryListCoursesByInstitution() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-courses-by-institution [institution-index]",
		Short: "List courses by institution",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := coursetypes.NewQueryClient(clientCtx)
			params := &coursetypes.QueryCoursesByInstitutionRequest{
				InstitutionIndex: args[0],
				Pagination:       pageReq,
			}

			res, err := queryClient.CoursesByInstitution(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// cmdQueryGetCourseCount implements the get-course-count query command
func cmdQueryGetCourseCount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-course-count",
		Short: "Get the total count of courses",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := coursetypes.NewQueryClient(clientCtx)
			res, err := queryClient.CourseCount(context.Background(), &coursetypes.QueryCourseCountRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// Add course transaction commands
func getCourseTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "course",
		Short:                      "Course transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		cmdTxCreateCourse(),
		cmdTxUpdateCourse(),
	)

	return cmd
}

// Course transaction command implementations
func cmdTxCreateCourse() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-course [institution] [name] [code] [description] [total-credits] [degree-level]",
		Short: "Create a new course",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			institution := args[0]
			name := args[1]
			code := args[2]
			description := args[3]
			totalCredits := args[4]
			degreeLevel := args[5]

			msg := &coursetypes.MsgCreateCourse{
				Creator:      clientCtx.GetFromAddress().String(),
				Institution:  institution,
				Name:         name,
				Code:         code,
				Description:  description,
				TotalCredits: totalCredits,
				DegreeLevel:  degreeLevel,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func cmdTxUpdateCourse() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-course [index] [name] [description] [total-credits]",
		Short: "Update an existing course",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			index := args[0]
			name := args[1]
			description := args[2]
			totalCredits := args[3]

			msg := &coursetypes.MsgUpdateCourse{
				Creator:      clientCtx.GetFromAddress().String(),
				Index:        index,
				Name:         name,
				Description:  description,
				TotalCredits: totalCredits,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// ============================================================================
// SUBJECT MODULE COMMANDS
// ============================================================================

// Subject Query Commands
func getSubjectQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "subject",
		Short:                      "Querying commands for the subject module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		cmdQuerySubjectParams(),
		cmdQueryGetSubject(),
		cmdQueryGetSubjectFull(),
		cmdQueryGetSubjectWithPrerequisites(),
		cmdQueryListSubjects(),
		cmdQueryListSubjectsByCourse(),      // Renamed from cmdQuerySubjectsByCourse
		cmdQueryListSubjectsByInstitution(), // Renamed from cmdQuerySubjectsByInstitution
		cmdQueryCheckPrerequisites(),
		cmdQueryCheckEquivalence(),
	)

	return cmd
}

// Subject query command implementations
func cmdQuerySubjectParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Shows the parameters of the subject module",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := subjecttypes.NewQueryClient(clientCtx)
			res, err := queryClient.Params(context.Background(), &subjecttypes.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryGetSubject() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-subject [index]",
		Short: "Show a specific subject by index",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := subjecttypes.NewQueryClient(clientCtx)
			params := &subjecttypes.QueryGetSubjectRequest{
				Index: args[0],
			}

			res, err := queryClient.GetSubject(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryGetSubjectFull() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-subject-full [index]",
		Short: "Show a subject with full content (including IPFS data) by index",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := subjecttypes.NewQueryClient(clientCtx)
			params := &subjecttypes.QueryGetSubjectFullRequest{
				Index: args[0],
			}

			res, err := queryClient.GetSubjectFull(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryGetSubjectWithPrerequisites() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-subject-with-prerequisites [subject-id]",
		Short: "Show a subject with its prerequisite groups",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := subjecttypes.NewQueryClient(clientCtx)
			params := &subjecttypes.QueryGetSubjectWithPrerequisitesRequest{
				SubjectId: args[0],
			}

			res, err := queryClient.GetSubjectWithPrerequisites(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryListSubjects() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-subjects",
		Short: "List all subjects with pagination",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := subjecttypes.NewQueryClient(clientCtx)
			params := &subjecttypes.QueryListSubjectsRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.ListSubjects(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// cmdQueryListSubjectsByCourse implements the list-subjects-by-course query command
func cmdQueryListSubjectsByCourse() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-subjects-by-course [course-id]",
		Short: "List subjects by course",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := subjecttypes.NewQueryClient(clientCtx)
			params := &subjecttypes.QuerySubjectsByCourseRequest{
				CourseId:   args[0],
				Pagination: pageReq,
			}

			res, err := queryClient.SubjectsByCourse(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// cmdQueryListSubjectsByInstitution implements the list-subjects-by-institution query command
func cmdQueryListSubjectsByInstitution() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-subjects-by-institution [institution-id]",
		Short: "List subjects by institution",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := subjecttypes.NewQueryClient(clientCtx)
			params := &subjecttypes.QuerySubjectsByInstitutionRequest{
				InstitutionId: args[0],
				Pagination:    pageReq,
			}

			res, err := queryClient.SubjectsByInstitution(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryCheckPrerequisites() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check-prerequisites [student-id] [subject-id]",
		Short: "Check if a student meets prerequisites for a subject",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := subjecttypes.NewQueryClient(clientCtx)
			params := &subjecttypes.QueryCheckPrerequisitesRequest{
				StudentId: args[0],
				SubjectId: args[1],
			}

			res, err := queryClient.CheckPrerequisites(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryCheckEquivalence() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check-equivalence [source-subject-id] [target-subject-id]",
		Short: "Check equivalence between two subjects",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			forceRecalculate, _ := cmd.Flags().GetBool("force-recalculate")

			queryClient := subjecttypes.NewQueryClient(clientCtx)
			params := &subjecttypes.QueryCheckEquivalenceRequest{
				SourceSubjectId:  args[0],
				TargetSubjectId:  args[1],
				ForceRecalculate: forceRecalculate,
			}

			res, err := queryClient.CheckEquivalence(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	cmd.Flags().Bool("force-recalculate", false, "Force recalculation of equivalence")
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// Subject Transaction Commands
func getSubjectTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "subject",
		Short:                      "Subject transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		//cmdTxCreateSubject(),
		cmdTxAddPrerequisiteGroup(),
		cmdTxUpdateSubjectContent(),
		CmdCreateSubjectContent(),
	)

	return cmd
}

/*
// Subject transaction command implementations
func cmdTxCreateSubject() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-subject [institution] [course-id] [title] [code] [workload-hours] [credits] [description] [subject-type] [knowledge-area]",
		Short: "Create a new subject",
		Args:  cobra.ExactArgs(9),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// Parse basic arguments
			institution := args[0]
			courseId := args[1]
			title := args[2]
			code := args[3]

			workloadHours, err := strconv.ParseUint(args[4], 10, 64)
			if err != nil {
				return err
			}

			credits, err := strconv.ParseUint(args[5], 10, 64)
			if err != nil {
				return err
			}

			description := args[6]
			subjectType := args[7]
			knowledgeArea := args[8]

			// Parse optional arrays from flags
			objectives, _ := cmd.Flags().GetStringSlice("objectives")
			topicUnits, _ := cmd.Flags().GetStringSlice("topic-units")

			msg := &subjecttypes.MsgCreateSubject{
				Creator:       clientCtx.GetFromAddress().String(),
				Institution:   institution,
				CourseId:      courseId,
				Title:         title,
				Code:          code,
				WorkloadHours: workloadHours,
				Credits:       credits,
				Description:   description,
				SubjectType:   subjectType,
				KnowledgeArea: knowledgeArea,
				Objectives:    objectives,
				TopicUnits:    topicUnits,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().StringSlice("objectives", []string{}, "Subject objectives")
	cmd.Flags().StringSlice("topic-units", []string{}, "Subject topic units")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}*/

func CmdCreateSubjectContent() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-subject [institution] [course-id] [title] [code] [workload-hours] [credits] [description] [subject-type] [knowledge-area]",
		Short: "Create a new subject content",
		Args:  cobra.ExactArgs(9),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// Parse numeric arguments
			workloadHours, err := strconv.ParseUint(args[4], 10, 64)
			if err != nil {
				return err
			}

			credits, err := strconv.ParseUint(args[5], 10, 64)
			if err != nil {
				return err
			}

			// Get flags
			objectives, _ := cmd.Flags().GetStringSlice("objectives")
			topicUnits, _ := cmd.Flags().GetStringSlice("topic-units")

			msg := subjecttypes.NewMsgCreateSubjectContent(
				clientCtx.GetFromAddress().String(),
				args[0], // institution
				args[1], // course_id
				args[2], // title
				args[3], // code
				workloadHours,
				credits,
				args[6], // description
				args[7], // subject_type
				args[8], // knowledge_area
				objectives,
				topicUnits,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().StringSlice("objectives", []string{}, "Subject objectives")
	cmd.Flags().StringSlice("topic-units", []string{}, "Subject topic units")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func cmdTxAddPrerequisiteGroup() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-prerequisite-group [subject-id] [group-type] [minimum-credits] [minimum-completed-subjects] [subject-ids...]",
		Short: "Add a prerequisite group to a subject",
		Args:  cobra.MinimumNArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			subjectId := args[0]
			groupType := args[1]

			minimumCredits, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			minimumCompletedSubjects, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return err
			}

			subjectIds := args[4:]

			msg := &subjecttypes.MsgAddPrerequisiteGroup{
				Creator:                  clientCtx.GetFromAddress().String(),
				SubjectId:                subjectId,
				GroupType:                groupType,
				MinimumCredits:           minimumCredits,
				MinimumCompletedSubjects: minimumCompletedSubjects,
				SubjectIds:               subjectIds,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func cmdTxUpdateSubjectContent() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-subject-content [subject-id]",
		Short: "Update the content of an existing subject",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			subjectId := args[0]

			// Get all the optional fields from flags
			objectives, _ := cmd.Flags().GetStringSlice("objectives")
			topicUnits, _ := cmd.Flags().GetStringSlice("topic-units")
			methodologies, _ := cmd.Flags().GetStringSlice("methodologies")
			evaluationMethods, _ := cmd.Flags().GetStringSlice("evaluation-methods")
			bibliographyBasic, _ := cmd.Flags().GetStringSlice("bibliography-basic")
			bibliographyComplementary, _ := cmd.Flags().GetStringSlice("bibliography-complementary")
			keywords, _ := cmd.Flags().GetStringSlice("keywords")
			contentHash, _ := cmd.Flags().GetString("content-hash")
			ipfsLink, _ := cmd.Flags().GetString("ipfs-link")

			msg := &subjecttypes.MsgUpdateSubjectContent{
				Creator:                   clientCtx.GetFromAddress().String(),
				SubjectId:                 subjectId,
				Objectives:                objectives,
				TopicUnits:                topicUnits,
				Methodologies:             methodologies,
				EvaluationMethods:         evaluationMethods,
				BibliographyBasic:         bibliographyBasic,
				BibliographyComplementary: bibliographyComplementary,
				Keywords:                  keywords,
				ContentHash:               contentHash,
				IpfsLink:                  ipfsLink,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().StringSlice("objectives", []string{}, "Subject objectives")
	cmd.Flags().StringSlice("topic-units", []string{}, "Subject topic units")
	cmd.Flags().StringSlice("methodologies", []string{}, "Teaching methodologies")
	cmd.Flags().StringSlice("evaluation-methods", []string{}, "Evaluation methods")
	cmd.Flags().StringSlice("bibliography-basic", []string{}, "Basic bibliography")
	cmd.Flags().StringSlice("bibliography-complementary", []string{}, "Complementary bibliography")
	cmd.Flags().StringSlice("keywords", []string{}, "Keywords")
	cmd.Flags().String("content-hash", "", "Content hash (optional)")
	cmd.Flags().String("ipfs-link", "", "IPFS link (optional)")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// ============================================================================
// CURRICULUM MODULE COMMANDS
// ============================================================================

// Curriculum Query Commands
func getCurriculumQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "curriculum",
		Short:                      "Querying commands for the curriculum module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		cmdQueryCurriculumParams(),
		cmdQueryListCurriculumTrees(),         // Renamed from cmdQueryCurriculumTreeAll
		cmdQueryGetCurriculumTree(),           // Renamed from cmdQueryCurriculumTree
		cmdQueryListCurriculumTreesByCourse(), // Renamed from cmdQueryCurriculumTreesByCourse
	)

	return cmd
}

// cmdQueryCurriculumParams implements the curriculum params query command
func cmdQueryCurriculumParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Shows the parameters of the curriculum module",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := curriculumtypes.NewQueryClient(clientCtx)
			res, err := queryClient.Params(context.Background(), &curriculumtypes.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// cmdQueryListCurriculumTrees implements the list-curriculum-trees query command
func cmdQueryListCurriculumTrees() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-curriculum-trees",
		Short: "List all curriculum trees",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := curriculumtypes.NewQueryClient(clientCtx)
			params := &curriculumtypes.QueryAllCurriculumTreeRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.CurriculumTreeAll(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// cmdQueryGetCurriculumTree implements the get-curriculum-tree query command
func cmdQueryGetCurriculumTree() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-curriculum-tree [index]",
		Short: "Show a specific curriculum tree by index",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := curriculumtypes.NewQueryClient(clientCtx)
			params := &curriculumtypes.QueryGetCurriculumTreeRequest{
				Index: args[0],
			}

			res, err := queryClient.CurriculumTree(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// cmdQueryListCurriculumTreesByCourse implements the list-curriculum-trees-by-course query command
func cmdQueryListCurriculumTreesByCourse() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-curriculum-trees-by-course [course-id]",
		Short: "List curriculum trees by course",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := curriculumtypes.NewQueryClient(clientCtx)
			params := &curriculumtypes.QueryCurriculumTreesByCourseRequest{
				CourseId: args[0],
			}

			res, err := queryClient.CurriculumTreesByCourse(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// Curriculum Transaction Commands
func getCurriculumTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "curriculum",
		Short:                      "Curriculum transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		cmdTxCreateCurriculumTree(),
		cmdTxAddSemesterToCurriculum(),
		cmdTxAddElectiveGroup(),
		cmdTxSetGraduationRequirements(),
		cmdTxUpdateCurriculumParams(),
	)

	return cmd
}

// cmdTxCreateCurriculumTree implements the create-curriculum-tree transaction command
func cmdTxCreateCurriculumTree() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-curriculum-tree [course-id] [version] [elective-min] [total-workload-hours]",
		Short: "Create a new curriculum tree",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			courseId := args[0]
			version := args[1]

			electiveMin, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			totalWorkloadHours, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return err
			}

			// Parse optional arrays from flags
			requiredSubjects, _ := cmd.Flags().GetStringSlice("required-subjects")
			electiveSubjects, _ := cmd.Flags().GetStringSlice("elective-subjects")

			msg := &curriculumtypes.MsgCreateCurriculumTree{
				Creator:            clientCtx.GetFromAddress().String(),
				CourseId:           courseId,
				Version:            version,
				ElectiveMin:        electiveMin,
				TotalWorkloadHours: totalWorkloadHours,
				RequiredSubjects:   requiredSubjects,
				ElectiveSubjects:   electiveSubjects,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().StringSlice("required-subjects", []string{}, "Required subject IDs")
	cmd.Flags().StringSlice("elective-subjects", []string{}, "Elective subject IDs")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// cmdTxAddSemesterToCurriculum implements the add-semester-to-curriculum transaction command
func cmdTxAddSemesterToCurriculum() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-semester-to-curriculum [curriculum-index] [semester-number] [subject-ids...]",
		Short: "Add a semester to a curriculum tree",
		Args:  cobra.MinimumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			curriculumIndex := args[0]

			semesterNumber, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			subjectIds := args[2:]

			msg := &curriculumtypes.MsgAddSemesterToCurriculum{
				Creator:         clientCtx.GetFromAddress().String(),
				CurriculumIndex: curriculumIndex,
				SemesterNumber:  semesterNumber,
				SubjectIds:      subjectIds,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// cmdTxAddElectiveGroup implements the add-elective-group transaction command
func cmdTxAddElectiveGroup() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-elective-group [curriculum-index] [name] [description] [min-subjects-required] [credits-required] [knowledge-area] [subject-ids...]",
		Short: "Add an elective group to a curriculum tree",
		Args:  cobra.MinimumNArgs(7),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			curriculumIndex := args[0]
			name := args[1]
			description := args[2]

			minSubjectsRequired, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return err
			}

			creditsRequired, err := strconv.ParseUint(args[4], 10, 64)
			if err != nil {
				return err
			}

			knowledgeArea := args[5]
			subjectIds := args[6:]

			msg := &curriculumtypes.MsgAddElectiveGroup{
				Creator:             clientCtx.GetFromAddress().String(),
				CurriculumIndex:     curriculumIndex,
				Name:                name,
				Description:         description,
				MinSubjectsRequired: minSubjectsRequired,
				CreditsRequired:     creditsRequired,
				KnowledgeArea:       knowledgeArea,
				SubjectIds:          subjectIds,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// cmdTxSetGraduationRequirements implements the set-graduation-requirements transaction command
func cmdTxSetGraduationRequirements() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-graduation-requirements [curriculum-index] [total-credits-required] [min-gpa] [required-elective-credits] [minimum-time-years] [maximum-time-years]",
		Short: "Set graduation requirements for a curriculum tree",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			curriculumIndex := args[0]

			totalCreditsRequired, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			minGpa, err := strconv.ParseFloat(args[2], 32)
			if err != nil {
				return err
			}

			requiredElectiveCredits, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return err
			}

			minimumTimeYears, err := strconv.ParseFloat(args[4], 32)
			if err != nil {
				return err
			}

			maximumTimeYears, err := strconv.ParseFloat(args[5], 32)
			if err != nil {
				return err
			}

			// Parse optional array from flags
			requiredActivities, _ := cmd.Flags().GetStringSlice("required-activities")

			msg := &curriculumtypes.MsgSetGraduationRequirements{
				Creator:                 clientCtx.GetFromAddress().String(),
				CurriculumIndex:         curriculumIndex,
				TotalCreditsRequired:    totalCreditsRequired,
				MinGpa:                  float32(minGpa),
				RequiredElectiveCredits: requiredElectiveCredits,
				MinimumTimeYears:        float32(minimumTimeYears),
				MaximumTimeYears:        float32(maximumTimeYears),
				RequiredActivities:      requiredActivities,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().StringSlice("required-activities", []string{}, "Required activities for graduation")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// cmdTxUpdateCurriculumParams implements the update-params transaction command for curriculum module
func cmdTxUpdateCurriculumParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-params [ipfs-gateway] [ipfs-enabled] [admin]",
		Short: "Update curriculum module parameters",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			ipfsGateway := args[0]
			ipfsEnabled := args[1] == "true"
			admin := args[2]

			params := curriculumtypes.NewParams(ipfsGateway, ipfsEnabled, admin)

			msg := &curriculumtypes.MsgUpdateParams{
				Authority: clientCtx.GetFromAddress().String(),
				Params:    params,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// ============================================================================
// TOKENDEF MODULE COMMANDS
// ============================================================================

// getTokendefQueryCmd creates the tokendef query commands
func getTokendefQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "tokendef",
		Short:                      "Querying commands for the tokendef module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		cmdQueryTokendefParams(),
		cmdQueryGetTokenDefinition(),
		cmdQueryGetTokenDefinitionFull(),
		cmdQueryListTokenDefinitions(),
		cmdQueryListTokenDefinitionsBySubject(),
	)

	return cmd
}

// cmdQueryTokendefParams implements the params query command
func cmdQueryTokendefParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Shows the parameters of the tokendef module",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := tokendefmoduletypes.NewQueryClient(clientCtx)
			res, err := queryClient.Params(context.Background(), &tokendefmoduletypes.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// cmdQueryGetTokenDefinition implements the get-token-definition query command
func cmdQueryGetTokenDefinition() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-token-definition [index]",
		Short: "Show a specific token definition by index",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := tokendefmoduletypes.NewQueryClient(clientCtx)
			params := &tokendefmoduletypes.QueryGetTokenDefinitionRequest{
				Index: args[0],
			}

			res, err := queryClient.GetTokenDefinition(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// cmdQueryGetTokenDefinitionFull implements the get-token-definition-full query command
func cmdQueryGetTokenDefinitionFull() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-token-definition-full [index]",
		Short: "Show a token definition with full content (including IPFS data) by index",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := tokendefmoduletypes.NewQueryClient(clientCtx)
			params := &tokendefmoduletypes.QueryGetTokenDefinitionFullRequest{
				Index: args[0],
			}

			res, err := queryClient.GetTokenDefinitionFull(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// cmdQueryListTokenDefinitions implements the list-token-definitions query command
func cmdQueryListTokenDefinitions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-token-definitions",
		Short: "List all token definitions with pagination",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := tokendefmoduletypes.NewQueryClient(clientCtx)
			params := &tokendefmoduletypes.QueryListTokenDefinitionsRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.ListTokenDefinitions(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// cmdQueryListTokenDefinitionsBySubject implements the list-token-definitions-by-subject query command
func cmdQueryListTokenDefinitionsBySubject() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-token-definitions-by-subject [subject-id]",
		Short: "List token definitions by subject",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := tokendefmoduletypes.NewQueryClient(clientCtx)
			params := &tokendefmoduletypes.QueryListTokenDefinitionsBySubjectRequest{
				SubjectId:  args[0],
				Pagination: pageReq,
			}

			res, err := queryClient.ListTokenDefinitionsBySubject(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// getTokendefTxCmd creates the tokendef transaction commands
func getTokendefTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "tokendef",
		Short:                      "Tokendef transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		cmdTxCreateTokenDefinition(),
		cmdTxUpdateTokenDefinition(),
		cmdTxUpdateTokendefParams(),
	)

	return cmd
}

// cmdTxCreateTokenDefinition implements the create-token-definition transaction command
func cmdTxCreateTokenDefinition() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-token-definition [subject-id] [token-name] [token-symbol] [description] [token-type]",
		Short: "Create a new token definition",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			subjectId := args[0]
			tokenName := args[1]
			tokenSymbol := args[2]
			description := args[3]
			tokenType := args[4]

			// Parse optional flags
			isTransferable, _ := cmd.Flags().GetBool("is-transferable")
			isBurnable, _ := cmd.Flags().GetBool("is-burnable")
			maxSupply, _ := cmd.Flags().GetUint64("max-supply")
			imageUri, _ := cmd.Flags().GetString("image-uri")

			msg := &tokendefmoduletypes.MsgCreateTokenDefinition{
				Creator:        clientCtx.GetFromAddress().String(),
				SubjectId:      subjectId,
				TokenName:      tokenName,
				TokenSymbol:    tokenSymbol,
				Description:    description,
				TokenType:      tokenType,
				IsTransferable: isTransferable,
				IsBurnable:     isBurnable,
				MaxSupply:      maxSupply,
				ImageUri:       imageUri,
				// Attributes can be added later if needed
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().Bool("is-transferable", false, "Whether the token is transferable")
	cmd.Flags().Bool("is-burnable", false, "Whether the token is burnable")
	cmd.Flags().Uint64("max-supply", 0, "Maximum supply of tokens")
	cmd.Flags().String("image-uri", "", "Image URI for the token")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// cmdTxUpdateTokendefParams implements the update-params transaction command for tokendef module
func cmdTxUpdateTokendefParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-params [ipfs-gateway] [ipfs-enabled] [admin]",
		Short: "Update tokendef module parameters",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			ipfsGateway := args[0]
			ipfsEnabled := args[1] == "true"
			admin := args[2]

			params := tokendefmoduletypes.NewParams(ipfsGateway, ipfsEnabled, admin)

			msg := &tokendefmoduletypes.MsgUpdateParams{
				Authority: clientCtx.GetFromAddress().String(),
				Params:    params,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// cmdTxUpdateTokenDefinition implements the update-token-definition transaction command
func cmdTxUpdateTokenDefinition() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-token-definition [token-def-id] [token-name] [token-symbol] [description]",
		Short: "Update an existing token definition",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			tokenDefId := args[0]
			tokenName := args[1]
			tokenSymbol := args[2]
			description := args[3]

			// Parse optional flags
			isTransferable, _ := cmd.Flags().GetBool("is-transferable")
			isBurnable, _ := cmd.Flags().GetBool("is-burnable")
			maxSupply, _ := cmd.Flags().GetUint64("max-supply")
			imageUri, _ := cmd.Flags().GetString("image-uri")

			msg := &tokendefmoduletypes.MsgUpdateTokenDefinition{
				Creator:        clientCtx.GetFromAddress().String(),
				TokenDefId:     tokenDefId,
				TokenName:      tokenName,
				TokenSymbol:    tokenSymbol,
				Description:    description,
				IsTransferable: isTransferable,
				IsBurnable:     isBurnable,
				MaxSupply:      maxSupply,
				ImageUri:       imageUri,
				// Attributes can be added later if needed
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().Bool("is-transferable", false, "Whether the token is transferable")
	cmd.Flags().Bool("is-burnable", false, "Whether the token is burnable")
	cmd.Flags().Uint64("max-supply", 0, "Maximum supply of tokens")
	cmd.Flags().String("image-uri", "", "Image URI for the token")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// ============================================================================
// ACADEMICNFT MODULE COMMANDS
// ============================================================================

// getAcademicNftQueryCmd creates the academicnft query commands
func getAcademicNftQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "academicnft",
		Short:                      "Querying commands for the academicnft module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		cmdQueryAcademicNftParams(),
		cmdQueryGetTokenInstance(),
		cmdQueryGetStudentTokens(),
		cmdQueryGetTokenDefInstances(),
		cmdQueryVerifyTokenInstance(),
	)

	return cmd
}

// cmdQueryAcademicNftParams implements the params query command
func cmdQueryAcademicNftParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Shows the parameters of the academicnft module",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := academicnfttypes.NewQueryClient(clientCtx)
			res, err := queryClient.Params(context.Background(), &academicnfttypes.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// cmdQueryGetTokenInstance implements the get-token-instance query command
func cmdQueryGetTokenInstance() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-token-instance [token-instance-id]",
		Short: "Show a specific token instance by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := academicnfttypes.NewQueryClient(clientCtx)
			params := &academicnfttypes.QueryGetTokenInstanceRequest{
				TokenInstanceId: args[0],
			}

			res, err := queryClient.GetTokenInstance(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// cmdQueryGetStudentTokens implements the get-student-tokens query command
func cmdQueryGetStudentTokens() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-student-tokens [student-address]",
		Short: "Get all token instances for a student",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := academicnfttypes.NewQueryClient(clientCtx)
			params := &academicnfttypes.QueryGetStudentTokensRequest{
				StudentAddress: args[0],
				Pagination:     pageReq,
			}

			res, err := queryClient.GetStudentTokens(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// cmdQueryGetTokenDefInstances implements the get-tokendef-instances query command
func cmdQueryGetTokenDefInstances() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-tokendef-instances [token-def-id]",
		Short: "Get all instances of a specific TokenDef",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := academicnfttypes.NewQueryClient(clientCtx)
			params := &academicnfttypes.QueryGetTokenDefInstancesRequest{
				TokenDefId: args[0],
				Pagination: pageReq,
			}

			res, err := queryClient.GetTokenDefInstances(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// cmdQueryVerifyTokenInstance implements the verify-token-instance query command
func cmdQueryVerifyTokenInstance() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify-token-instance [token-instance-id]",
		Short: "Verify if a token instance exists and is valid",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := academicnfttypes.NewQueryClient(clientCtx)
			params := &academicnfttypes.QueryVerifyTokenInstanceRequest{
				TokenInstanceId: args[0],
			}

			res, err := queryClient.VerifyTokenInstance(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// getAcademicNftTxCmd creates the academicnft transaction commands
func getAcademicNftTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "academicnft",
		Short:                      "AcademicNft transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		cmdTxMintSubjectToken(),
		cmdTxVerifyTokenInstance(),
		cmdTxUpdateAcademicNftParams(),
	)

	return cmd
}

// cmdTxMintSubjectToken implements the mint-subject-token transaction command
func cmdTxMintSubjectToken() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint-subject-token [token-def-id] [student] [completion-date] [grade] [issuer-institution] [semester]",
		Short: "Mint a new subject token for a student",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			tokenDefId := args[0]
			student := args[1]
			completionDate := args[2]
			grade := args[3]
			issuerInstitution := args[4]
			semester := args[5]

			// Parse optional professor signature from flag
			professorSignature, _ := cmd.Flags().GetString("professor-signature")

			msg := &academicnfttypes.MsgMintSubjectToken{
				Creator:            clientCtx.GetFromAddress().String(),
				TokenDefId:         tokenDefId,
				Student:            student,
				CompletionDate:     completionDate,
				Grade:              grade,
				IssuerInstitution:  issuerInstitution,
				Semester:           semester,
				ProfessorSignature: professorSignature,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("professor-signature", "", "Professor signature for the token")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// cmdTxVerifyTokenInstance implements the verify-token-instance transaction command
func cmdTxVerifyTokenInstance() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify-token-instance [token-instance-id]",
		Short: "Verify a token instance on-chain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			tokenInstanceId := args[0]

			msg := &academicnfttypes.MsgVerifyTokenInstance{
				Creator:         clientCtx.GetFromAddress().String(),
				TokenInstanceId: tokenInstanceId,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// cmdTxUpdateAcademicNftParams implements the update-params transaction command for academicnft module
func cmdTxUpdateAcademicNftParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-params [ipfs-gateway] [ipfs-enabled] [admin]",
		Short: "Update academicnft module parameters",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			ipfsGateway := args[0]
			ipfsEnabled := args[1] == "true"
			admin := args[2]

			params := academicnfttypes.NewParams(ipfsGateway, ipfsEnabled, admin)

			msg := &academicnfttypes.MsgUpdateParams{
				Authority: clientCtx.GetFromAddress().String(),
				Params:    params,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// ============================================================================
// STUDENT MODULE COMMANDS
// ============================================================================

// getStudentQueryCmd creates the student query commands
func getStudentQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "student",
		Short:                      "Querying commands for the student module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		cmdQueryStudentParams(),
		cmdQueryListStudents(),
		cmdQueryGetStudent(),
		cmdQueryListEnrollments(),
		cmdQueryGetEnrollment(),
		cmdQueryGetEnrollmentsByStudent(),
		cmdQueryGetStudentProgress(),
		cmdQueryGetStudentsByInstitution(),
		cmdQueryGetStudentsByCourse(),
		cmdQueryGetStudentAcademicTree(),
		cmdQueryCheckGraduationEligibility(),
	)

	return cmd
}

// Student query command implementations
func cmdQueryStudentParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Shows the parameters of the student module",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := studenttypes.NewQueryClient(clientCtx)
			res, err := queryClient.Params(context.Background(), &studenttypes.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryListStudents() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-students",
		Short: "List all students with pagination",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := studenttypes.NewQueryClient(clientCtx)
			params := &studenttypes.QueryListStudentsRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.ListStudents(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryGetStudent() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-student [student-id]",
		Short: "Show a specific student by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := studenttypes.NewQueryClient(clientCtx)
			params := &studenttypes.QueryGetStudentRequest{
				StudentId: args[0],
			}

			res, err := queryClient.GetStudent(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryListEnrollments() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-enrollments",
		Short: "List all enrollments with pagination",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := studenttypes.NewQueryClient(clientCtx)
			params := &studenttypes.QueryListEnrollmentsRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.ListEnrollments(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryGetEnrollment() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-enrollment [enrollment-id]",
		Short: "Show a specific enrollment by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := studenttypes.NewQueryClient(clientCtx)
			params := &studenttypes.QueryGetEnrollmentRequest{
				EnrollmentId: args[0],
			}

			res, err := queryClient.GetEnrollment(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryGetEnrollmentsByStudent() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-enrollments-by-student [student-id]",
		Short: "Get all enrollments for a specific student",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := studenttypes.NewQueryClient(clientCtx)
			params := &studenttypes.QueryGetEnrollmentsByStudentRequest{
				StudentId:  args[0],
				Pagination: pageReq,
			}

			res, err := queryClient.GetEnrollmentsByStudent(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryGetStudentProgress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-student-progress [student-id]",
		Short: "Get the academic progress of a student",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := studenttypes.NewQueryClient(clientCtx)
			params := &studenttypes.QueryGetStudentProgressRequest{
				StudentId: args[0],
			}

			res, err := queryClient.GetStudentProgress(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryGetStudentsByInstitution() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-students-by-institution [institution-id]",
		Short: "Get all students enrolled in a specific institution",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := studenttypes.NewQueryClient(clientCtx)
			params := &studenttypes.QueryGetStudentsByInstitutionRequest{
				InstitutionId: args[0],
				Pagination:    pageReq,
			}

			res, err := queryClient.GetStudentsByInstitution(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryGetStudentsByCourse() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-students-by-course [course-id]",
		Short: "Get all students enrolled in a specific course",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := studenttypes.NewQueryClient(clientCtx)
			params := &studenttypes.QueryGetStudentsByCourseRequest{
				CourseId:   args[0],
				Pagination: pageReq,
			}

			res, err := queryClient.GetStudentsByCourse(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryGetStudentAcademicTree() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-student-academic-tree [student-id]",
		Short: "Get the academic tree (curriculum progress) of a student",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := studenttypes.NewQueryClient(clientCtx)
			params := &studenttypes.QueryGetStudentAcademicTreeRequest{
				StudentId: args[0],
			}

			res, err := queryClient.GetStudentAcademicTree(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryCheckGraduationEligibility() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check-graduation-eligibility [student-id]",
		Short: "Check if a student is eligible for graduation",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := studenttypes.NewQueryClient(clientCtx)
			params := &studenttypes.QueryCheckGraduationEligibilityRequest{
				StudentId: args[0],
			}

			res, err := queryClient.CheckGraduationEligibility(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// getStudentTxCmd creates the student transaction commands
func getStudentTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "student",
		Short:                      "Student transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		cmdTxRegisterStudent(),
		cmdTxCreateEnrollment(),
		cmdTxUpdateEnrollmentStatus(),
		cmdTxRequestSubjectEnrollment(),
		cmdTxUpdateAcademicTree(),
		cmdTxUpdateStudentParams(),
	)

	return cmd
}

// Student transaction command implementations
func cmdTxRegisterStudent() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-student [name] [address]",
		Short: "Register a new student",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			name := args[0]
			address := args[1]

			msg := &studenttypes.MsgRegisterStudent{
				Creator: clientCtx.GetFromAddress().String(),
				Name:    name,
				Address: address,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func cmdTxCreateEnrollment() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-enrollment [student] [institution] [course-id]",
		Short: "Create a new enrollment for a student",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			student := args[0]
			institution := args[1]
			courseId := args[2]

			msg := &studenttypes.MsgCreateEnrollment{
				Creator:     clientCtx.GetFromAddress().String(),
				Student:     student,
				Institution: institution,
				CourseId:    courseId,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func cmdTxUpdateEnrollmentStatus() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-enrollment-status [enrollment-id] [status]",
		Short: "Update the status of an enrollment",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			enrollmentId := args[0]
			status := args[1]

			msg := &studenttypes.MsgUpdateEnrollmentStatus{
				Creator:      clientCtx.GetFromAddress().String(),
				EnrollmentId: enrollmentId,
				Status:       status,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func cmdTxRequestSubjectEnrollment() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "request-subject-enrollment [student] [subject-id]",
		Short: "Request enrollment in a specific subject",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			student := args[0]
			subjectId := args[1]

			msg := &studenttypes.MsgRequestSubjectEnrollment{
				Creator:   clientCtx.GetFromAddress().String(),
				Student:   student,
				SubjectId: subjectId,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func cmdTxUpdateAcademicTree() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-academic-tree [student-id]",
		Short: "Update the academic tree (progress) of a student",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			studentId := args[0]

			// Parse token arrays from flags
			completedTokens, _ := cmd.Flags().GetStringSlice("completed-tokens")
			inProgressTokens, _ := cmd.Flags().GetStringSlice("in-progress-tokens")
			availableTokens, _ := cmd.Flags().GetStringSlice("available-tokens")

			msg := &studenttypes.MsgUpdateAcademicTree{
				Creator:          clientCtx.GetFromAddress().String(),
				StudentId:        studentId,
				CompletedTokens:  completedTokens,
				InProgressTokens: inProgressTokens,
				AvailableTokens:  availableTokens,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().StringSlice("completed-tokens", []string{}, "List of completed token IDs")
	cmd.Flags().StringSlice("in-progress-tokens", []string{}, "List of in-progress token IDs")
	cmd.Flags().StringSlice("available-tokens", []string{}, "List of available token IDs")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// cmdTxUpdateStudentParams implements the update-params transaction command for student module
func cmdTxUpdateStudentParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-params [ipfs-gateway] [ipfs-enabled] [admin]",
		Short: "Update student module parameters",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			ipfsGateway := args[0]
			ipfsEnabled := args[1] == "true"
			admin := args[2]

			params := studenttypes.NewParams(ipfsGateway, ipfsEnabled, admin, "", "", "", "", "")

			msg := &studenttypes.MsgUpdateParams{
				Authority: clientCtx.GetFromAddress().String(),
				Params:    params,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// ============================================================================
// DEGREE MODULE COMMANDS
// ============================================================================

// getDegreeQueryCmd creates the degree query commands
func getDegreeQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "degree",
		Short:                      "Querying commands for the degree module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		cmdQueryDegreeParams(),
		cmdQueryListDegrees(),
		cmdQueryGetDegree(),
		cmdQueryListDegreesByStudent(),
		cmdQueryListDegreesByInstitution(),
		cmdQueryListDegreeRequests(),
		cmdQueryGetDegreeValidationStatus(),
	)

	return cmd
}

// Degree query command implementations
func cmdQueryDegreeParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Shows the parameters of the degree module",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := degreetypes.NewQueryClient(clientCtx)
			res, err := queryClient.Params(context.Background(), &degreetypes.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryListDegrees() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-degrees",
		Short: "List all degrees with pagination",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := degreetypes.NewQueryClient(clientCtx)
			params := &degreetypes.QueryAllDegreeRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.DegreeAll(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryGetDegree() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-degree [index]",
		Short: "Show a specific degree by index",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := degreetypes.NewQueryClient(clientCtx)
			params := &degreetypes.QueryGetDegreeRequest{
				Index: args[0],
			}

			res, err := queryClient.Degree(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryListDegreesByStudent() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-degrees-by-student [student-id]",
		Short: "List degrees by student ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := degreetypes.NewQueryClient(clientCtx)
			params := &degreetypes.QueryDegreesByStudentRequest{
				StudentId:  args[0],
				Pagination: pageReq,
			}

			res, err := queryClient.DegreesByStudent(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryListDegreesByInstitution() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-degrees-by-institution [institution-id]",
		Short: "List degrees by institution ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := degreetypes.NewQueryClient(clientCtx)
			params := &degreetypes.QueryDegreesByInstitutionRequest{
				InstitutionId: args[0],
				Pagination:    pageReq,
			}

			res, err := queryClient.DegreesByInstitution(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryListDegreeRequests() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-degree-requests",
		Short: "List degree requests (pending degrees)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			status, _ := cmd.Flags().GetString("status")

			queryClient := degreetypes.NewQueryClient(clientCtx)
			params := &degreetypes.QueryDegreeRequestsRequest{
				Status:     status,
				Pagination: pageReq,
			}

			res, err := queryClient.DegreeRequests(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	cmd.Flags().String("status", "", "Filter by status (optional)")
	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryGetDegreeValidationStatus() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-degree-validation-status [degree-id]",
		Short: "Get degree validation status by degree ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := degreetypes.NewQueryClient(clientCtx)
			params := &degreetypes.QueryDegreeValidationStatusRequest{
				DegreeId: args[0],
			}

			res, err := queryClient.DegreeValidationStatus(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// getDegreeTxCmd creates the degree transaction commands
func getDegreeTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "degree",
		Short:                      "Degree transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		cmdTxRequestDegree(),
		cmdTxValidateDegreeRequirements(),
		cmdTxIssueDegree(),
		cmdTxUpdateDegreeContract(),
		cmdTxCancelDegreeRequest(),
	)

	return cmd
}

// Degree transaction command implementations
func cmdTxRequestDegree() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "request-degree [student-id] [institution-id] [curriculum-id] [expected-graduation-date]",
		Short: "Request a degree for a student",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			studentId := args[0]
			institutionId := args[1]
			curriculumId := args[2]
			expectedGraduationDate := args[3]

			msg := &degreetypes.MsgRequestDegree{
				Creator:                clientCtx.GetFromAddress().String(),
				StudentId:              studentId,
				InstitutionId:          institutionId,
				CurriculumId:           curriculumId,
				ExpectedGraduationDate: expectedGraduationDate,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func cmdTxValidateDegreeRequirements() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate-degree-requirements [degree-request-id] [contract-address] [validation-parameters]",
		Short: "Validate degree requirements via CosmWasm contract",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			degreeRequestId := args[0]
			contractAddress := args[1]
			validationParameters := args[2]

			msg := &degreetypes.MsgValidateDegreeRequirements{
				Creator:              clientCtx.GetFromAddress().String(),
				DegreeRequestId:      degreeRequestId,
				ContractAddress:      contractAddress,
				ValidationParameters: validationParameters,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func cmdTxIssueDegree() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue-degree [degree-request-id] [final-gpa] [total-credits]",
		Short: "Issue a degree after successful validation",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			degreeRequestId := args[0]
			finalGpa := args[1]

			totalCredits, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			// Parse optional flags
			signatures, _ := cmd.Flags().GetStringSlice("signatures")
			additionalNotes, _ := cmd.Flags().GetString("additional-notes")

			msg := &degreetypes.MsgIssueDegree{
				Creator:         clientCtx.GetFromAddress().String(),
				DegreeRequestId: degreeRequestId,
				FinalGpa:        finalGpa,
				TotalCredits:    totalCredits,
				Signatures:      signatures,
				AdditionalNotes: additionalNotes,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().StringSlice("signatures", []string{}, "Digital signatures from authorized personnel")
	cmd.Flags().String("additional-notes", "", "Additional notes for the degree")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func cmdTxUpdateDegreeContract() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-degree-contract [new-contract-address] [contract-version] [migration-reason]",
		Short: "Update the CosmWasm contract address for degree validation",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			newContractAddress := args[0]
			contractVersion := args[1]
			migrationReason := args[2]

			msg := &degreetypes.MsgUpdateDegreeContract{
				Authority:          clientCtx.GetFromAddress().String(),
				NewContractAddress: newContractAddress,
				ContractVersion:    contractVersion,
				MigrationReason:    migrationReason,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func cmdTxCancelDegreeRequest() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-degree-request [degree-request-id] [cancellation-reason]",
		Short: "Cancel a pending degree request",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			degreeRequestId := args[0]
			cancellationReason := args[1]

			msg := &degreetypes.MsgCancelDegreeRequest{
				Creator:            clientCtx.GetFromAddress().String(),
				DegreeRequestId:    degreeRequestId,
				CancellationReason: cancellationReason,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// ============================================================================
// EQUIVALENCE MODULE COMMANDS
// ============================================================================

// getEquivalenceQueryCmd creates the equivalence query commands
func getEquivalenceQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "equivalence",
		Short:                      "Querying commands for the equivalence module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		cmdQueryEquivalenceParams(),
		cmdQueryListEquivalences(),
		cmdQueryGetEquivalence(),
		cmdQueryGetEquivalencesBySourceSubject(),
		cmdQueryGetEquivalencesByTargetSubject(),
		cmdQueryGetEquivalencesByInstitution(),
		cmdQueryCheckEquivalenceStatus(),
		cmdQueryGetPendingAnalysis(),
		cmdQueryGetApprovedEquivalences(),
		cmdQueryGetRejectedEquivalences(),
		cmdQueryGetEquivalencesByContract(),
		cmdQueryGetEquivalencesByContractVersion(),
		cmdQueryGetEquivalenceHistory(),
		cmdQueryGetEquivalenceStats(),
		cmdQueryGetAnalysisMetadata(),
		cmdQueryVerifyAnalysisIntegrity(),
	)

	return cmd
}

// Equivalence query command implementations
func cmdQueryEquivalenceParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Shows the parameters of the equivalence module",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := equivalencetypes.NewQueryClient(clientCtx)
			res, err := queryClient.Params(context.Background(), &equivalencetypes.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryListEquivalences() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-equivalences",
		Short: "List all subject equivalences with pagination",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			statusFilter, _ := cmd.Flags().GetString("status-filter")

			queryClient := equivalencetypes.NewQueryClient(clientCtx)
			params := &equivalencetypes.QueryListEquivalencesRequest{
				Pagination:   pageReq,
				StatusFilter: statusFilter,
			}

			res, err := queryClient.ListEquivalences(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	cmd.Flags().String("status-filter", "", "Filter by status (optional)")
	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryGetEquivalence() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-equivalence [index]",
		Short: "Show a specific subject equivalence by index",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := equivalencetypes.NewQueryClient(clientCtx)
			params := &equivalencetypes.QueryGetEquivalenceRequest{
				Index: args[0],
			}

			res, err := queryClient.GetEquivalence(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryGetEquivalencesBySourceSubject() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-equivalences-by-source-subject [source-subject-id]",
		Short: "Get equivalences by source subject ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			statusFilter, _ := cmd.Flags().GetString("status-filter")

			queryClient := equivalencetypes.NewQueryClient(clientCtx)
			params := &equivalencetypes.QueryGetEquivalencesBySourceSubjectRequest{
				SourceSubjectId: args[0],
				Pagination:      pageReq,
				StatusFilter:    statusFilter,
			}

			res, err := queryClient.GetEquivalencesBySourceSubject(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	cmd.Flags().String("status-filter", "", "Filter by status (optional)")
	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryGetEquivalencesByTargetSubject() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-equivalences-by-target-subject [target-subject-id]",
		Short: "Get equivalences by target subject ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			statusFilter, _ := cmd.Flags().GetString("status-filter")

			queryClient := equivalencetypes.NewQueryClient(clientCtx)
			params := &equivalencetypes.QueryGetEquivalencesByTargetSubjectRequest{
				TargetSubjectId: args[0],
				Pagination:      pageReq,
				StatusFilter:    statusFilter,
			}

			res, err := queryClient.GetEquivalencesByTargetSubject(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	cmd.Flags().String("status-filter", "", "Filter by status (optional)")
	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryGetEquivalencesByInstitution() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-equivalences-by-institution [institution-id]",
		Short: "Get equivalences by target institution",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			statusFilter, _ := cmd.Flags().GetString("status-filter")

			queryClient := equivalencetypes.NewQueryClient(clientCtx)
			params := &equivalencetypes.QueryGetEquivalencesByInstitutionRequest{
				InstitutionId: args[0],
				Pagination:    pageReq,
				StatusFilter:  statusFilter,
			}

			res, err := queryClient.GetEquivalencesByInstitution(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	cmd.Flags().String("status-filter", "", "Filter by status (optional)")
	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryCheckEquivalenceStatus() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check-equivalence-status [source-subject-id] [target-subject-id]",
		Short: "Check if two subjects have an established equivalence",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := equivalencetypes.NewQueryClient(clientCtx)
			params := &equivalencetypes.QueryCheckEquivalenceStatusRequest{
				SourceSubjectId: args[0],
				TargetSubjectId: args[1],
			}

			res, err := queryClient.CheckEquivalenceStatus(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryGetPendingAnalysis() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-pending-analysis",
		Short: "Get equivalences awaiting contract analysis",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := equivalencetypes.NewQueryClient(clientCtx)
			params := &equivalencetypes.QueryGetPendingAnalysisRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.GetPendingAnalysis(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryGetApprovedEquivalences() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-approved-equivalences",
		Short: "Get equivalences with approved status",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			minEquivalencePercent, _ := cmd.Flags().GetString("min-equivalence-percent")

			queryClient := equivalencetypes.NewQueryClient(clientCtx)
			params := &equivalencetypes.QueryGetApprovedEquivalencesRequest{
				Pagination:            pageReq,
				MinEquivalencePercent: minEquivalencePercent,
			}

			res, err := queryClient.GetApprovedEquivalences(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	cmd.Flags().String("min-equivalence-percent", "", "Minimum equivalence percentage filter")
	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryGetRejectedEquivalences() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-rejected-equivalences",
		Short: "Get equivalences rejected by contract analysis",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := equivalencetypes.NewQueryClient(clientCtx)
			params := &equivalencetypes.QueryGetRejectedEquivalencesRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.GetRejectedEquivalences(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryGetEquivalencesByContract() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-equivalences-by-contract [contract-address]",
		Short: "Get equivalences analyzed by a specific contract",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := equivalencetypes.NewQueryClient(clientCtx)
			params := &equivalencetypes.QueryGetEquivalencesByContractRequest{
				ContractAddress: args[0],
				Pagination:      pageReq,
			}

			res, err := queryClient.GetEquivalencesByContract(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryGetEquivalencesByContractVersion() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-equivalences-by-contract-version [contract-version]",
		Short: "Get equivalences by contract version",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := equivalencetypes.NewQueryClient(clientCtx)
			params := &equivalencetypes.QueryGetEquivalencesByContractVersionRequest{
				ContractVersion: args[0],
				Pagination:      pageReq,
			}

			res, err := queryClient.GetEquivalencesByContractVersion(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryGetEquivalenceHistory() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-equivalence-history [subject-id]",
		Short: "Get the analysis history of equivalence requests for a subject",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := equivalencetypes.NewQueryClient(clientCtx)
			params := &equivalencetypes.QueryGetEquivalenceHistoryRequest{
				SubjectId:  args[0],
				Pagination: pageReq,
			}

			res, err := queryClient.GetEquivalenceHistory(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryGetEquivalenceStats() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-equivalence-stats",
		Short: "Get statistics about automated equivalence analysis",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := equivalencetypes.NewQueryClient(clientCtx)
			params := &equivalencetypes.QueryGetEquivalenceStatsRequest{}

			res, err := queryClient.GetEquivalenceStats(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryGetAnalysisMetadata() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-analysis-metadata [equivalence-id]",
		Short: "Get detailed analysis metadata for an equivalence",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := equivalencetypes.NewQueryClient(clientCtx)
			params := &equivalencetypes.QueryGetAnalysisMetadataRequest{
				EquivalenceId: args[0],
			}

			res, err := queryClient.GetAnalysisMetadata(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryVerifyAnalysisIntegrity() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify-analysis-integrity [equivalence-id]",
		Short: "Verify the integrity of an equivalence analysis",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := equivalencetypes.NewQueryClient(clientCtx)
			params := &equivalencetypes.QueryVerifyAnalysisIntegrityRequest{
				EquivalenceId: args[0],
			}

			res, err := queryClient.VerifyAnalysisIntegrity(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// getEquivalenceTxCmd creates the equivalence transaction commands
func getEquivalenceTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "equivalence",
		Short:                      "Equivalence transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		cmdTxRequestEquivalence(),
		cmdTxExecuteEquivalenceAnalysis(),
		cmdTxBatchRequestEquivalence(),
		cmdTxUpdateContractAddress(),
		cmdTxReanalyzeEquivalence(),
	)

	return cmd
}

// Equivalence transaction command implementations
func cmdTxRequestEquivalence() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "request-equivalence [source-subject-id] [target-institution] [target-subject-id]",
		Short: "Request equivalence between subjects",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			sourceSubjectId := args[0]
			targetInstitution := args[1]
			targetSubjectId := args[2]

			forceRecalculation, _ := cmd.Flags().GetBool("force-recalculation")

			msg := &equivalencetypes.MsgRequestEquivalence{
				Creator:            clientCtx.GetFromAddress().String(),
				SourceSubjectId:    sourceSubjectId,
				TargetInstitution:  targetInstitution,
				TargetSubjectId:    targetSubjectId,
				ForceRecalculation: forceRecalculation,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().Bool("force-recalculation", false, "Force recalculation of equivalence")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func cmdTxExecuteEquivalenceAnalysis() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "execute-equivalence-analysis [equivalence-id] [contract-address]",
		Short: "Execute contract analysis for equivalence",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			equivalenceId := args[0]
			contractAddress := args[1]

			msg := &equivalencetypes.MsgExecuteEquivalenceAnalysis{
				Creator:         clientCtx.GetFromAddress().String(),
				EquivalenceId:   equivalenceId,
				ContractAddress: contractAddress,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func cmdTxBatchRequestEquivalence() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "batch-request-equivalence [requests-json]",
		Short: "Submit multiple equivalence requests in batch",
		Long: `Submit multiple equivalence requests in batch. 
The requests-json should be a JSON string containing an array of equivalence requests.
Example: '[{"sourceSubjectId":"subj1","targetInstitution":"inst1","targetSubjectId":"subj2"}]'`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			forceRecalculation, _ := cmd.Flags().GetBool("force-recalculation")

			// Parse the JSON requests
			// For simplicity, we'll accept a comma-separated list of \"source:target_inst:target_subj\"
			// In a real implementation, you'd parse JSON here
			requestStr := args[0]
			_ = requestStr // TODO: Parse actual JSON requests

			// Placeholder - in real implementation, parse JSON to create requests array
			var requests []*equivalencetypes.EquivalenceRequest

			msg := &equivalencetypes.MsgBatchRequestEquivalence{
				Creator:            clientCtx.GetFromAddress().String(),
				Requests:           requests,
				ForceRecalculation: forceRecalculation,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().Bool("force-recalculation", false, "Force recalculation for all equivalences")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func cmdTxUpdateContractAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-contract-address [new-contract-address] [contract-version]",
		Short: "Update the CosmWasm contract address for equivalence analysis",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			newContractAddress := args[0]
			contractVersion := args[1]

			msg := &equivalencetypes.MsgUpdateContractAddress{
				Authority:          clientCtx.GetFromAddress().String(),
				NewContractAddress: newContractAddress,
				ContractVersion:    contractVersion,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func cmdTxReanalyzeEquivalence() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reanalyze-equivalence [equivalence-id] [reanalysis-reason]",
		Short: "Re-analyze an existing equivalence",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			equivalenceId := args[0]
			reanalysisReason := args[1]

			contractAddress, _ := cmd.Flags().GetString("contract-address")

			msg := &equivalencetypes.MsgReanalyzeEquivalence{
				Creator:          clientCtx.GetFromAddress().String(),
				EquivalenceId:    equivalenceId,
				ReanalysisReason: reanalysisReason,
				ContractAddress:  contractAddress,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("contract-address", "", "Specific contract address to use (optional)")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// ============================================================================
// SCHEDULE MODULE COMMANDS
// ============================================================================

// getScheduleQueryCmd creates the schedule query commands
func getScheduleQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "schedule",
		Short:                      "Querying commands for the schedule module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		cmdQueryScheduleParams(),
		cmdQueryGetStudyPlan(),
		cmdQueryListStudyPlans(),
		cmdQueryGetStudyPlansByStudent(),
		cmdQueryGetPlannedSemester(),
		cmdQueryGetSubjectRecommendation(),
		cmdQueryGetSubjectRecommendationsByStudent(),
		cmdQueryGenerateRecommendations(),
		cmdQueryCheckStudentProgress(),
		cmdQueryOptimizeSchedule(),
	)

	return cmd
}

// Schedule query command implementations
func cmdQueryScheduleParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Shows the parameters of the schedule module",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := scheduletypes.NewQueryClient(clientCtx)
			res, err := queryClient.Params(context.Background(), &scheduletypes.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryGetStudyPlan() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-study-plan [study-plan-id]",
		Short: "Get a specific study plan by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := scheduletypes.NewQueryClient(clientCtx)
			params := &scheduletypes.QueryGetStudyPlanRequest{
				StudyPlanId: args[0],
			}

			res, err := queryClient.StudyPlan(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryListStudyPlans() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-study-plans",
		Short: "List all study plans with pagination",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := scheduletypes.NewQueryClient(clientCtx)
			params := &scheduletypes.QueryAllStudyPlanRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.StudyPlanAll(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryGetStudyPlansByStudent() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-study-plans-by-student [student-id]",
		Short: "Get all study plans for a specific student",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := scheduletypes.NewQueryClient(clientCtx)
			params := &scheduletypes.QueryStudyPlansByStudentRequest{
				StudentId: args[0],
			}

			res, err := queryClient.StudyPlansByStudent(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryGetPlannedSemester() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-planned-semester [planned-semester-id]",
		Short: "Get a specific planned semester by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := scheduletypes.NewQueryClient(clientCtx)
			params := &scheduletypes.QueryGetPlannedSemesterRequest{
				PlannedSemesterId: args[0],
			}

			res, err := queryClient.PlannedSemester(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryGetSubjectRecommendation() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-subject-recommendation [recommendation-id]",
		Short: "Get a specific subject recommendation by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := scheduletypes.NewQueryClient(clientCtx)
			params := &scheduletypes.QueryGetSubjectRecommendationRequest{
				RecommendationId: args[0],
			}

			res, err := queryClient.SubjectRecommendation(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryGetSubjectRecommendationsByStudent() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-subject-recommendations-by-student [student-id]",
		Short: "Get all recommendations for a student",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := scheduletypes.NewQueryClient(clientCtx)
			params := &scheduletypes.QuerySubjectRecommendationsByStudentRequest{
				StudentId: args[0],
			}

			res, err := queryClient.SubjectRecommendationsByStudent(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryGenerateRecommendations() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate-recommendations [student-id]",
		Short: "Generate automatic recommendations for a student",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			semesterCode, _ := cmd.Flags().GetString("semester-code")

			queryClient := scheduletypes.NewQueryClient(clientCtx)
			params := &scheduletypes.QueryGenerateRecommendationsRequest{
				StudentId:    args[0],
				SemesterCode: semesterCode,
			}

			res, err := queryClient.GenerateRecommendations(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	cmd.Flags().String("semester-code", "", "Semester code for recommendations (optional)")
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryCheckStudentProgress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check-student-progress [student-id]",
		Short: "Check the academic progress of a student",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			courseId, _ := cmd.Flags().GetString("course-id")

			queryClient := scheduletypes.NewQueryClient(clientCtx)
			params := &scheduletypes.QueryCheckStudentProgressRequest{
				StudentId: args[0],
				CourseId:  courseId,
			}

			res, err := queryClient.CheckStudentProgress(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	cmd.Flags().String("course-id", "", "Course ID to check progress for (optional)")
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryOptimizeSchedule() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "optimize-schedule [study-plan-id]",
		Short: "Optimize a student's schedule",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := scheduletypes.NewQueryClient(clientCtx)
			params := &scheduletypes.QueryOptimizeScheduleRequest{
				StudyPlanId: args[0],
			}

			res, err := queryClient.OptimizeSchedule(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// getScheduleTxCmd creates the schedule transaction commands
func getScheduleTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "schedule",
		Short:                      "Schedule transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		cmdTxCreateStudyPlan(),
		cmdTxAddPlannedSemester(),
		cmdTxCreateSubjectRecommendation(),
		cmdTxUpdateStudyPlanStatus(),
	)

	return cmd
}

// Schedule transaction command implementations with CORRECT field names
func cmdTxCreateStudyPlan() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-study-plan [student] [completion-target]",
		Short: "Create a new study plan",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			student := args[0]
			completionTarget := args[1]

			semesterCodes, _ := cmd.Flags().GetStringSlice("semester-codes")
			additionalNotes, _ := cmd.Flags().GetString("additional-notes")
			status, _ := cmd.Flags().GetString("status")

			msg := &scheduletypes.MsgCreateStudyPlan{
				Creator:          clientCtx.GetFromAddress().String(),
				Student:          student,
				CompletionTarget: completionTarget,
				SemesterCodes:    semesterCodes,
				AdditionalNotes:  additionalNotes,
				Status:           status,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().StringSlice("semester-codes", []string{}, "Semester codes for the study plan")
	cmd.Flags().String("additional-notes", "", "Additional notes for the study plan")
	cmd.Flags().String("status", "active", "Status of the study plan")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func cmdTxAddPlannedSemester() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-planned-semester [study-plan-id] [semester-code] [planned-subjects...]",
		Short: "Add a planned semester to a study plan",
		Args:  cobra.MinimumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			studyPlanId := args[0]
			semesterCode := args[1]
			plannedSubjects := args[2:]

			totalCredits, _ := cmd.Flags().GetUint64("total-credits")
			totalHours, _ := cmd.Flags().GetUint64("total-hours")
			status, _ := cmd.Flags().GetString("status")

			msg := &scheduletypes.MsgAddPlannedSemester{
				Creator:         clientCtx.GetFromAddress().String(),
				StudyPlanId:     studyPlanId,
				SemesterCode:    semesterCode,
				PlannedSubjects: plannedSubjects,
				TotalCredits:    totalCredits,
				TotalHours:      totalHours,
				Status:          status,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().Uint64("total-credits", 0, "Total credits for the semester")
	cmd.Flags().Uint64("total-hours", 0, "Total hours for the semester")
	cmd.Flags().String("status", "planned", "Status of the planned semester")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func cmdTxCreateSubjectRecommendation() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-subject-recommendation [student] [recommendation-semester]",
		Short: "Create a subject recommendation",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			student := args[0]
			recommendationSemester := args[1]

			recommendationMetadata, _ := cmd.Flags().GetString("recommendation-metadata")
			generatedDate, _ := cmd.Flags().GetString("generated-date")

			// For RecommendedSubjects, we'll need to parse from flags or make it empty for now
			// In a real implementation, you'd parse JSON or have a more sophisticated input method
			var recommendedSubjects []*scheduletypes.RecommendedSubject

			msg := &scheduletypes.MsgCreateSubjectRecommendation{
				Creator:                clientCtx.GetFromAddress().String(),
				Student:                student,
				RecommendationSemester: recommendationSemester,
				RecommendedSubjects:    recommendedSubjects,
				RecommendationMetadata: recommendationMetadata,
				GeneratedDate:          generatedDate,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("recommendation-metadata", "", "Metadata for the recommendation")
	cmd.Flags().String("generated-date", "", "Date when the recommendation was generated")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func cmdTxUpdateStudyPlanStatus() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-study-plan-status [study-plan-id] [status]",
		Short: "Update the status of a study plan",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			studyPlanId := args[0]
			status := args[1]

			msg := &scheduletypes.MsgUpdateStudyPlanStatus{
				Creator:     clientCtx.GetFromAddress().String(),
				StudyPlanId: studyPlanId,
				Status:      status,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// ============================================================================
// APP FUNCTIONS
// ============================================================================

// newApp creates the application
func newApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	appOpts servertypes.AppOptions,
) servertypes.Application {
	baseappOptions := server.DefaultBaseappOptions(appOpts)

	app, err := app.New(
		logger, db, traceStore, true,
		appOpts,
		baseappOptions...,
	)
	if err != nil {
		panic(err)
	}
	return app
}

// appExport creates a new app (optionally at a given height) and exports state.
func appExport(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	height int64,
	forZeroHeight bool,
	jailAllowedAddrs []string,
	appOpts servertypes.AppOptions,
	modulesToExport []string,
) (servertypes.ExportedApp, error) {
	var (
		bApp *app.App
		err  error
	)

	// this check is necessary as we use the flag in x/upgrade.
	// we can exit more gracefully by checking the flag here.
	homePath, ok := appOpts.Get(flags.FlagHome).(string)
	if !ok || homePath == "" {
		return servertypes.ExportedApp{}, errors.New("application home not set")
	}

	viperAppOpts, ok := appOpts.(*viper.Viper)
	if !ok {
		return servertypes.ExportedApp{}, errors.New("appOpts is not viper.Viper")
	}

	// overwrite the FlagInvCheckPeriod
	viperAppOpts.Set(server.FlagInvCheckPeriod, 1)
	appOpts = viperAppOpts

	if height != -1 {
		bApp, err = app.New(logger, db, traceStore, false, appOpts)
		if err != nil {
			return servertypes.ExportedApp{}, err
		}

		if err := bApp.LoadHeight(height); err != nil {
			return servertypes.ExportedApp{}, err
		}
	} else {
		bApp, err = app.New(logger, db, traceStore, true, appOpts)
		if err != nil {
			return servertypes.ExportedApp{}, err
		}
	}

	return bApp.ExportAppStateAndValidators(forZeroHeight, jailAllowedAddrs, modulesToExport)
}
