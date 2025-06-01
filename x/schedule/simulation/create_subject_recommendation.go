package simulation

import (
	"math/rand"
	"strconv"
	"time"

	"academictoken/x/schedule/keeper"
	"academictoken/x/schedule/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

func SimulateMsgCreateSubjectRecommendation(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		// Select random accounts for creator and student
		creatorAcc, _ := simtypes.RandomAcc(r, accs)
		studentAcc, _ := simtypes.RandomAcc(r, accs)

		// Generate random subject IDs for recommendations
		subjectIds := []string{
			"subject-math-calculus", "subject-cs-algorithms", "subject-eng-physics",
			"subject-bio-genetics", "subject-chem-organic", "subject-art-history",
			"subject-econ-micro", "subject-stat-intro", "subject-lit-classic",
		}

		// Generate random recommendation semester
		semesters := []string{"2024-1", "2024-2", "2025-1", "2025-2", "2026-1"}
		recommendationSemester := semesters[r.Intn(len(semesters))]

		// Generate random recommendation metadata
		metadataOptions := []string{
			"AI-generated based on academic performance",
			"Recommended by academic advisor",
			"Generated from curriculum optimization",
			"Based on student preferences",
			"Prerequisite chain recommendation",
		}
		recommendationMetadata := metadataOptions[r.Intn(len(metadataOptions))]

		// Create 1-4 recommended subjects
		numRecommendations := r.Intn(4) + 1
		var recommendedSubjects []*types.RecommendedSubject

		for i := 0; i < numRecommendations; i++ {
			subjectId := subjectIds[r.Intn(len(subjectIds))]

			// Random values for each recommended subject
			reasons := []string{
				"Strong prerequisite alignment",
				"Matches student performance profile",
				"Optimal path to graduation",
				"High success probability",
				"Complements current course load",
			}
			reason := reasons[r.Intn(len(reasons))]

			// Random difficulty levels
			difficultyLevels := []string{"easy", "medium", "hard", "advanced"}
			difficultyLevel := difficultyLevels[r.Intn(len(difficultyLevels))]

			// Random required status
			isRequiredOptions := []string{"true", "false"}
			isRequired := isRequiredOptions[r.Intn(len(isRequiredOptions))]

			recommendedSubject := &types.RecommendedSubject{
				SubjectId:          subjectId,
				RecommendationRank: strconv.Itoa(i + 1),
				Reason:             reason,
				IsRequired:         isRequired,
				SemesterAlignment:  recommendationSemester,
				DifficultyLevel:    difficultyLevel,
			}

			recommendedSubjects = append(recommendedSubjects, recommendedSubject)
		}

		// Create message with correct fields from protobuf
		msg := &types.MsgCreateSubjectRecommendation{
			Creator:                creatorAcc.Address.String(),
			Student:                studentAcc.Address.String(),
			RecommendationSemester: recommendationSemester,
			RecommendedSubjects:    recommendedSubjects,
			RecommendationMetadata: recommendationMetadata,
			GeneratedDate:          time.Now().Format(time.RFC3339),
		}

		// Validate the message
		if err := msg.ValidateBasic(); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "invalid message"), nil, err
		}

		// Check account balance
		creatorAccount := ak.GetAccount(ctx, creatorAcc.Address)
		spendable := bk.SpendableCoins(ctx, creatorAccount.GetAddress())

		// Create operation input with correct fields for newer SDK versions
		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           nil, // TxGen is handled internally in newer versions
			Cdc:             nil,
			Msg:             msg,
			Context:         ctx,
			SimAccount:      creatorAcc,
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: spendable,
		}

		// Deliver the transaction with correct API
		opMsg, fops, err := simulation.GenAndDeliverTxWithRandFees(txCtx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "failed to deliver tx"), nil, err
		}

		return opMsg, fops, nil
	}
}
