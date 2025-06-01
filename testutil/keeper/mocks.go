package keeper

import (
	"context"
	"fmt"

	"academictoken/x/course/types"
	degreetype "academictoken/x/degree/types"
	institutiontypes "academictoken/x/institution/types"
	scheduletypes "academictoken/x/schedule/types"
	studenttypes "academictoken/x/student/types"
	subjecttypes "academictoken/x/subject/types"
	tokendefTypes "academictoken/x/tokendef/types"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MockAccountKeeper is a mock implementation of AccountKeeper for testing
type MockAccountKeeper struct{}

func (m MockAccountKeeper) GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI {
	return nil // Return nil for testing
}

// MockBankKeeper is a mock implementation of BankKeeper for testing
type MockBankKeeper struct{}

func (m MockBankKeeper) SpendableCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins {
	return sdk.NewCoins() // Return empty coins for testing
}

// MockWasmMsgServer implements the WasmMsgServer interface for testing
type MockWasmMsgServer struct{}

func (m MockWasmMsgServer) ExecuteContract(ctx context.Context, req *wasmtypes.MsgExecuteContract) (*wasmtypes.MsgExecuteContractResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request is nil")
	}
	// Return mock successful execution
	return &wasmtypes.MsgExecuteContractResponse{
		Data: []byte(`{"success": true, "message": "Mock contract execution successful"}`),
	}, nil
}

// MockWasmQuerier implements the WasmQuerier interface for testing
type MockWasmQuerier struct{}

// SmartContractState mocks a smart contract state query
func (q MockWasmQuerier) SmartContractState(ctx context.Context, req *wasmtypes.QuerySmartContractStateRequest) (*wasmtypes.QuerySmartContractStateResponse, error) {
	// Return mock data for tests
	return &wasmtypes.QuerySmartContractStateResponse{
		Data: []byte(`{"is_eligible":true,"missing_prerequisites":[]}`),
	}, nil
}

// MockSubjectKeeper implements tokendefTypes.SubjectKeeper for testing
type MockSubjectKeeper struct{}

func (m MockSubjectKeeper) HasSubject(ctx sdk.Context, subjectId string) bool {
	return subjectId != "" // Simple mock: exists if not empty
}

func (m MockSubjectKeeper) GetSubject(ctx sdk.Context, subjectId string) (tokendefTypes.SubjectContent, bool) {
	if subjectId == "" {
		return tokendefTypes.SubjectContent{}, false
	}
	return tokendefTypes.SubjectContent{
		Index:         subjectId,
		SubjectId:     subjectId,
		Title:         "Mock Subject",
		Code:          "MOCK101",
		Credits:       4,
		WorkloadHours: 60,
		Institution:   "mock-institution",
		CourseId:      "mock-course",
		Creator:       "academic1test",
	}, true
}

// MockInstitutionKeeper implements ALL InstitutionKeeper interfaces
type MockInstitutionKeeper struct{}

// For tokendef module - FIXED: return correct type
func (m MockInstitutionKeeper) GetInstitution(ctx sdk.Context, institutionID string) (tokendefTypes.Institution, bool) {
	if institutionID == "" {
		return tokendefTypes.Institution{}, false
	}
	return tokendefTypes.Institution{
		Index:        institutionID,
		Address:      "Mock Address",
		Name:         "Mock Institution",
		IsAuthorized: "true",
		Creator:      "academic1test",
	}, true
}

func (m MockInstitutionKeeper) IsInstitutionAuthorized(ctx sdk.Context, institutionID string) bool {
	return institutionID != ""
}

func (m MockInstitutionKeeper) InstitutionExists(ctx sdk.Context, institutionID string) bool {
	return institutionID != ""
}

// Missing method for course module
func (m MockInstitutionKeeper) GetAuthorizedInstitutions(ctx sdk.Context) []institutiontypes.Institution {
	return []institutiontypes.Institution{
		{
			Index:        "institution-1",
			Address:      "Mock Address",
			Name:         "Mock Institution",
			IsAuthorized: "true",
			Creator:      "academic1test",
		},
	}
}

// ============================================================================
// MOCKS ESPECÍFICOS PARA O MÓDULO DEGREE
// ============================================================================

// MockDegreeStudentKeeper implements degree module's StudentKeeper interface
type MockDegreeStudentKeeper struct{}

func (m MockDegreeStudentKeeper) GetStudent(ctx sdk.Context, id string) (degreetype.Student, bool) {
	if id == "" {
		return degreetype.Student{}, false
	}
	return degreetype.Student{
		Id:            id,
		InstitutionId: "institution-1",
		Status:        "active",
	}, true
}

func (m MockDegreeStudentKeeper) GetStudentAcademicRecord(ctx sdk.Context, studentId string) (degreetype.AcademicRecord, bool) {
	if studentId == "" {
		return degreetype.AcademicRecord{}, false
	}
	return degreetype.AcademicRecord{
		StudentId:         studentId,
		CompletedCredits:  120,
		GPA:               "8.5",
		CompletedSubjects: []string{"subject-1", "subject-2", "subject-3"},
	}, true
}

func (m MockDegreeStudentKeeper) GetStudentsByInstitution(ctx sdk.Context, institutionId string) []degreetype.Student {
	if institutionId == "" {
		return []degreetype.Student{}
	}
	return []degreetype.Student{
		{
			Id:            "student-1",
			InstitutionId: institutionId,
			Status:        "active",
		},
	}
}

func (m MockDegreeStudentKeeper) ValidateStudentExists(ctx sdk.Context, studentId string) error {
	if studentId == "" {
		return fmt.Errorf("student not found")
	}
	return nil
}

func (m MockDegreeStudentKeeper) GetStudentGPA(ctx sdk.Context, studentId string) (string, error) {
	if studentId == "" {
		return "", fmt.Errorf("student not found")
	}
	return "8.5", nil
}

func (m MockDegreeStudentKeeper) GetStudentTotalCredits(ctx sdk.Context, studentId string) (uint64, error) {
	if studentId == "" {
		return 0, fmt.Errorf("student not found")
	}
	return 120, nil
}

func (m MockDegreeStudentKeeper) GetCompletedSubjects(ctx sdk.Context, studentId string) ([]string, error) {
	if studentId == "" {
		return nil, fmt.Errorf("student not found")
	}
	return []string{"subject-1", "subject-2", "subject-3"}, nil
}

// MockDegreeCurriculumKeeper implements degree module's CurriculumKeeper interface
type MockDegreeCurriculumKeeper struct{}

func (m MockDegreeCurriculumKeeper) GetCurriculum(ctx sdk.Context, id string) (degreetype.Curriculum, bool) {
	if id == "" {
		return degreetype.Curriculum{}, false
	}
	return degreetype.Curriculum{
		Id:               id,
		InstitutionId:    "institution-1",
		RequiredCredits:  120,
		RequiredSubjects: []string{"subject-1", "subject-2", "subject-3"},
	}, true
}

func (m MockDegreeCurriculumKeeper) ValidateCurriculumRequirements(ctx sdk.Context, curriculumId string, completedSubjects []string) error {
	if curriculumId == "" {
		return fmt.Errorf("curriculum not found")
	}
	// Simple validation: check if we have at least 3 subjects
	if len(completedSubjects) < 3 {
		return fmt.Errorf("insufficient completed subjects")
	}
	return nil
}

func (m MockDegreeCurriculumKeeper) GetCurriculumRequiredCredits(ctx sdk.Context, curriculumId string) (uint64, error) {
	if curriculumId == "" {
		return 0, fmt.Errorf("curriculum not found")
	}
	return 120, nil
}

func (m MockDegreeCurriculumKeeper) GetCurriculumRequiredSubjects(ctx sdk.Context, curriculumId string) ([]string, error) {
	if curriculumId == "" {
		return nil, fmt.Errorf("curriculum not found")
	}
	return []string{"subject-1", "subject-2", "subject-3"}, nil
}

func (m MockDegreeCurriculumKeeper) GetCurriculumsByInstitution(ctx sdk.Context, institutionId string) []degreetype.Curriculum {
	if institutionId == "" {
		return []degreetype.Curriculum{}
	}
	return []degreetype.Curriculum{
		{
			Id:               "curriculum-1",
			InstitutionId:    institutionId,
			RequiredCredits:  120,
			RequiredSubjects: []string{"subject-1", "subject-2", "subject-3"},
		},
	}
}

// MockDegreeAcademicNFTKeeper implements degree module's AcademicNFTKeeper interface
type MockDegreeAcademicNFTKeeper struct{}

func (m MockDegreeAcademicNFTKeeper) MintDegreeNFT(ctx sdk.Context, recipient string, degreeData degreetype.DegreeNFTData) (string, error) {
	if recipient == "" {
		return "", fmt.Errorf("invalid recipient")
	}
	// Generate a mock token ID
	tokenId := fmt.Sprintf("degree-nft-%s-%s", degreeData.StudentId, degreeData.CurriculumId)
	return tokenId, nil
}

func (m MockDegreeAcademicNFTKeeper) GetNFTByTokenID(ctx sdk.Context, tokenId string) (degreetype.AcademicNFT, bool) {
	if tokenId == "" {
		return degreetype.AcademicNFT{}, false
	}
	return degreetype.AcademicNFT{
		TokenId:   tokenId,
		Owner:     "student-1",
		Metadata:  "Mock degree NFT metadata",
		TokenType: "degree",
	}, true
}

func (m MockDegreeAcademicNFTKeeper) TransferNFT(ctx sdk.Context, from string, to string, tokenId string) error {
	if from == "" || to == "" || tokenId == "" {
		return fmt.Errorf("invalid transfer parameters")
	}
	return nil
}

func (m MockDegreeAcademicNFTKeeper) BurnNFT(ctx sdk.Context, tokenId string) error {
	if tokenId == "" {
		return fmt.Errorf("invalid token ID")
	}
	return nil
}

func (m MockDegreeAcademicNFTKeeper) ValidateNFTOwnership(ctx sdk.Context, owner string, tokenId string) error {
	if owner == "" || tokenId == "" {
		return fmt.Errorf("invalid ownership validation parameters")
	}
	return nil
}

// MockDegreeWasmKeeper implements degree module's WasmKeeper interface
type MockDegreeWasmKeeper struct{}

func (m MockDegreeWasmKeeper) Sudo(ctx sdk.Context, contractAddress sdk.AccAddress, msg []byte) ([]byte, error) {
	// Return mock successful response
	return []byte(`{"success":true}`), nil
}

func (m MockDegreeWasmKeeper) QuerySmart(ctx context.Context, contractAddr sdk.AccAddress, req []byte) ([]byte, error) {
	// Return mock query response
	return []byte(`{"eligible":true,"requirements_met":true}`), nil
}

func (m MockDegreeWasmKeeper) Execute(ctx sdk.Context, contractAddress sdk.AccAddress, caller sdk.AccAddress, msg []byte, coins sdk.Coins) ([]byte, error) {
	// Return mock execution response
	return []byte(`{"success":true,"executed":true}`), nil
}

// MockDegreeInstitutionKeeper implements degree module's InstitutionKeeper interface
type MockDegreeInstitutionKeeper struct{}

func (m MockDegreeInstitutionKeeper) GetInstitution(ctx sdk.Context, id string) (interface{}, bool) {
	if id == "" {
		return nil, false
	}
	// Return a generic institution interface{}
	return map[string]interface{}{
		"id":            id,
		"name":          "Mock Institution",
		"address":       "Mock Address",
		"is_authorized": true,
	}, true
}

func (m MockDegreeInstitutionKeeper) ValidateInstitutionAuthorization(ctx sdk.Context, institutionId string) error {
	if institutionId == "" {
		return fmt.Errorf("institution not found")
	}
	return nil
}

func (m MockDegreeInstitutionKeeper) IsAuthorizedToIssueDegrees(ctx sdk.Context, institutionId string) bool {
	return institutionId != ""
}

// For other modules that need institutiontypes.Institution
func (m MockInstitutionKeeper) GetInstitutionOriginal(ctx sdk.Context, institutionID string) (institutiontypes.Institution, bool) {
	return institutiontypes.Institution{
		Index:        institutionID,
		Address:      "Mock Address",
		Name:         "Mock Institution",
		IsAuthorized: "true",
		Creator:      "academic1test",
	}, true
}

// MockCourseKeeper implements ALL CourseKeeper interfaces
type MockCourseKeeper struct{}

// For tokendef module - FIXED: return correct type
func (m MockCourseKeeper) GetCourse(ctx sdk.Context, index string) (tokendefTypes.Course, bool) {
	if index == "" {
		return tokendefTypes.Course{}, false
	}
	return tokendefTypes.Course{
		Index:        index,
		Institution:  "institution-1",
		Name:         "Mock Course",
		Code:         "MOCK101",
		Description:  "Mock course for testing",
		TotalCredits: "120",
		DegreeLevel:  "bachelor",
	}, true
}

func (m MockCourseKeeper) HasCourse(ctx sdk.Context, courseId string) bool {
	return courseId != ""
}

func (m MockCourseKeeper) CourseExists(ctx sdk.Context, index string) bool {
	return index != ""
}

// For other modules that need types.Course
func (m MockCourseKeeper) GetCourseOriginal(ctx sdk.Context, index string) (types.Course, bool) {
	if index == "" {
		return types.Course{}, false
	}
	return types.Course{
		Index:        index,
		Institution:  "institution-1",
		Name:         "Mock Course",
		Code:         "MOCK101",
		Description:  "Mock course for testing",
		TotalCredits: "120",
		DegreeLevel:  "bachelor",
	}, true
}

// For original subject module
func (m MockSubjectKeeper) GetSubjectOriginal(ctx sdk.Context, subjectId string) (subjecttypes.SubjectContent, bool) {
	if subjectId == "" {
		return subjecttypes.SubjectContent{}, false
	}
	return subjecttypes.SubjectContent{
		Index:         subjectId,
		SubjectId:     subjectId,
		Title:         "Mock Subject",
		Code:          "MOCK101",
		Credits:       4,
		WorkloadHours: 60,
		Institution:   "mock-institution",
		CourseId:      "mock-course",
		Creator:       "academic1test",
	}, true
}

// MockCourseInstitutionKeeper implements course module's InstitutionKeeper interface
type MockCourseInstitutionKeeper struct{}

func (m MockCourseInstitutionKeeper) GetInstitution(ctx sdk.Context, institutionID string) (institutiontypes.Institution, bool) {
	return institutiontypes.Institution{
		Index:        institutionID,
		Address:      "Mock Address",
		Name:         "Mock Institution",
		IsAuthorized: "true",
		Creator:      "academic1test",
	}, true
}

func (m MockCourseInstitutionKeeper) IsInstitutionAuthorized(ctx sdk.Context, institutionID string) bool {
	return institutionID != ""
}

func (m MockCourseInstitutionKeeper) InstitutionExists(ctx sdk.Context, institutionID string) bool {
	return institutionID != ""
}

func (m MockCourseInstitutionKeeper) GetAuthorizedInstitutions(ctx sdk.Context) []institutiontypes.Institution {
	return []institutiontypes.Institution{
		{
			Index:        "institution-1",
			Address:      "Mock Address",
			Name:         "Mock Institution",
			IsAuthorized: "true",
			Creator:      "academic1test",
		},
	}
}

// MockSubjectInstitutionKeeper implements subject module's InstitutionKeeper interface
type MockSubjectInstitutionKeeper struct{}

func (m MockSubjectInstitutionKeeper) GetInstitution(ctx sdk.Context, institutionID string) (institutiontypes.Institution, bool) {
	return institutiontypes.Institution{
		Index:        institutionID,
		Address:      "Mock Address",
		Name:         "Mock Institution",
		IsAuthorized: "true",
		Creator:      "academic1test",
	}, true
}

func (m MockSubjectInstitutionKeeper) IsInstitutionAuthorized(ctx sdk.Context, institutionID string) bool {
	return institutionID != ""
}

func (m MockSubjectInstitutionKeeper) InstitutionExists(ctx sdk.Context, institutionID string) bool {
	return institutionID != ""
}

// MockSubjectCourseKeeper implements subject module's CourseKeeper interface
type MockSubjectCourseKeeper struct{}

func (m MockSubjectCourseKeeper) GetCourse(ctx sdk.Context, index string) (types.Course, bool) {
	if index == "" {
		return types.Course{}, false
	}
	return types.Course{
		Index:        index,
		Institution:  "institution-1",
		Name:         "Mock Course",
		Code:         "MOCK101",
		Description:  "Mock course for testing",
		TotalCredits: "120",
		DegreeLevel:  "bachelor",
	}, true
}

func (m MockSubjectCourseKeeper) CourseExists(ctx sdk.Context, index string) bool {
	return index != ""
}

func (m MockSubjectCourseKeeper) HasCourse(ctx sdk.Context, courseId string) bool {
	return courseId != ""
}

// ============================================================================
// MOCKS ESPECÍFICOS PARA O MÓDULO STUDENT (APENAS ADICIONADOS)
// ============================================================================

// MockStudentCurriculumKeeper implements student module's CurriculumKeeper interface
type MockStudentCurriculumKeeper struct{}

func (m MockStudentCurriculumKeeper) GetCurriculumTree(ctx context.Context, curriculumId string) (studenttypes.CurriculumTree, bool) {
	if curriculumId == "" {
		return studenttypes.CurriculumTree{}, false
	}
	return studenttypes.CurriculumTree{
		Index:              curriculumId,
		CourseId:           "course-1",
		Version:            "v1.0",
		RequiredSubjects:   []string{"subject-1", "subject-2"},
		ElectiveMin:        2,
		ElectiveSubjects:   []string{"elective-1", "elective-2"},
		TotalWorkloadHours: 2400,
		GraduationRequirements: studenttypes.GraduationRequirements{
			TotalCreditsRequired:    120,
			MinGPA:                  6.0,
			RequiredElectiveCredits: 20,
			RequiredActivities:      []string{"tcc", "estagio"},
			MinimumTimeYears:        4.0,
			MaximumTimeYears:        8.0,
		},
	}, true
}

func (m MockStudentCurriculumKeeper) GetCurriculumTreesByCourse(ctx context.Context, courseId string) []studenttypes.CurriculumTree {
	if courseId == "" {
		return []studenttypes.CurriculumTree{}
	}
	return []studenttypes.CurriculumTree{
		{
			Index:              "curriculum-1",
			CourseId:           courseId,
			Version:            "v1.0",
			RequiredSubjects:   []string{"subject-1", "subject-2"},
			ElectiveMin:        2,
			ElectiveSubjects:   []string{"elective-1", "elective-2"},
			TotalWorkloadHours: 2400,
			GraduationRequirements: studenttypes.GraduationRequirements{
				TotalCreditsRequired:    120,
				MinGPA:                  6.0,
				RequiredElectiveCredits: 20,
				RequiredActivities:      []string{"tcc", "estagio"},
				MinimumTimeYears:        4.0,
				MaximumTimeYears:        8.0,
			},
		},
	}
}

// MockStudentSubjectKeeper implements student module's SubjectKeeper interface
type MockStudentSubjectKeeper struct{}

func (m MockStudentSubjectKeeper) GetSubject(ctx sdk.Context, index string) (studenttypes.SubjectContent, bool) {
	if index == "" {
		return studenttypes.SubjectContent{}, false
	}
	return studenttypes.SubjectContent{
		Index:         index,
		SubjectId:     index,
		Institution:   "institution-1",
		Title:         "Mock Subject",
		Code:          "MOCK101",
		WorkloadHours: 60,
		Credits:       4,
		Description:   "Mock subject for testing",
		ContentHash:   "mock-hash",
		SubjectType:   "required",
		KnowledgeArea: "computer-science",
		IPFSLink:      "ipfs://mock-link",
	}, true
}

func (m MockStudentSubjectKeeper) SubjectExists(ctx sdk.Context, index string) bool {
	return index != ""
}

func (m MockStudentSubjectKeeper) CheckPrerequisitesViaContract(ctx sdk.Context, studentID string, subjectID string) (bool, []string, error) {
	return true, []string{}, nil
}

func (m MockStudentSubjectKeeper) CheckEquivalenceViaContract(ctx sdk.Context, sourceSubjectID string, targetSubjectID string, forceRecalculate bool) (uint64, string, error) {
	return 85, "High similarity", nil
}

// MockStudentTokenDefKeeper implements student module's TokenDefKeeper interface
type MockStudentTokenDefKeeper struct{}

func (m MockStudentTokenDefKeeper) GetTokenDefinitionByIndex(ctx sdk.Context, index string) (studenttypes.TokenDefinition, bool) {
	if index == "" {
		return studenttypes.TokenDefinition{}, false
	}
	return studenttypes.TokenDefinition{
		Index:          index,
		SubjectId:      "subject-1",
		TokenName:      "Mock Token",
		TokenSymbol:    "MTK",
		Description:    "Mock token description",
		TokenType:      "subject",
		IsTransferable: true,
		IsBurnable:     false,
		MaxSupply:      1000,
		ImageUri:       "https://example.com/image.png",
		ContentHash:    "mock-hash",
		IPFSLink:       "ipfs://mock-link",
	}, true
}

func (m MockStudentTokenDefKeeper) GetTokenDefinitionsBySubject(ctx sdk.Context, subjectId string) []studenttypes.TokenDefinition {
	if subjectId == "" {
		return []studenttypes.TokenDefinition{}
	}
	return []studenttypes.TokenDefinition{
		{
			Index:          "token-1",
			SubjectId:      subjectId,
			TokenName:      "Mock Token",
			TokenSymbol:    "MTK",
			Description:    "Mock token description",
			TokenType:      "subject",
			IsTransferable: true,
			IsBurnable:     false,
			MaxSupply:      1000,
			ImageUri:       "https://example.com/image.png",
			ContentHash:    "mock-hash",
			IPFSLink:       "ipfs://mock-link",
		},
	}
}

// MockStudentAcademicNFTKeeper implements student module's AcademicNFTKeeper interface
type MockStudentAcademicNFTKeeper struct{}

func (m MockStudentAcademicNFTKeeper) GetSubjectTokenInstance(ctx sdk.Context, tokenInstanceId string) (studenttypes.SubjectTokenInstance, bool) {
	if tokenInstanceId == "" {
		return studenttypes.SubjectTokenInstance{}, false
	}
	return studenttypes.SubjectTokenInstance{
		TokenInstanceId:    tokenInstanceId,
		TokenDefId:         "token-def-1",
		Student:            "student-1",
		CompletionDate:     "2024-01-01",
		Grade:              "A",
		IssuerInstitution:  "institution-1",
		Semester:           "2024-1",
		ProfessorSignature: "prof-signature",
		MintedAt:           "2024-01-01T00:00:00Z",
		IsValid:            true,
	}, true
}

func (m MockStudentAcademicNFTKeeper) GetStudentTokenInstances(ctx sdk.Context, studentAddress string) ([]studenttypes.SubjectTokenInstance, error) {
	if studentAddress == "" {
		return []studenttypes.SubjectTokenInstance{}, nil
	}
	return []studenttypes.SubjectTokenInstance{
		{
			TokenInstanceId:    "token-instance-1",
			TokenDefId:         "token-def-1",
			Student:            studentAddress,
			CompletionDate:     "2024-01-01",
			Grade:              "A",
			IssuerInstitution:  "institution-1",
			Semester:           "2024-1",
			ProfessorSignature: "prof-signature",
			MintedAt:           "2024-01-01T00:00:00Z",
			IsValid:            true,
		},
	}, nil
}

// MockStudentInstitutionKeeper implements student module's InstitutionKeeper interface
type MockStudentInstitutionKeeper struct{}

func (m MockStudentInstitutionKeeper) GetInstitution(ctx sdk.Context, index string) (studenttypes.Institution, bool) {
	if index == "" {
		return studenttypes.Institution{}, false
	}
	return studenttypes.Institution{
		Index:        index,
		Name:         "Mock Institution",
		Address:      "Mock Address",
		IsAuthorized: true,
	}, true
}

func (m MockStudentInstitutionKeeper) IsInstitutionAuthorized(ctx sdk.Context, institutionIndex string) bool {
	return institutionIndex != ""
}

// MockStudentCourseKeeper implements student module's CourseKeeper interface
type MockStudentCourseKeeper struct{}

func (m MockStudentCourseKeeper) GetCourse(ctx sdk.Context, index string) (studenttypes.Course, bool) {
	if index == "" {
		return studenttypes.Course{}, false
	}
	return studenttypes.Course{
		Index:        index,
		Institution:  "institution-1",
		Name:         "Mock Course",
		Code:         "MOCK101",
		Description:  "Mock course for testing",
		TotalCredits: "120",
		DegreeLevel:  "bachelor",
	}, true
}

func (m MockStudentCourseKeeper) GetCoursesByInstitution(ctx sdk.Context, institutionId string) []studenttypes.Course {
	if institutionId == "" {
		return []studenttypes.Course{}
	}
	return []studenttypes.Course{
		{
			Index:        "course-1",
			Institution:  institutionId,
			Name:         "Mock Course",
			Code:         "MOCK101",
			Description:  "Mock course for testing",
			TotalCredits: "120",
			DegreeLevel:  "bachelor",
		},
	}
}

// ============================================================================
// MOCKS ESPECÍFICOS PARA O MÓDULO SCHEDULE
// ============================================================================

// MockScheduleSubjectKeeper implements schedule module's SubjectKeeper interface
type MockScheduleSubjectKeeper struct{}

func (m MockScheduleSubjectKeeper) GetSubject(ctx sdk.Context, subjectID string) (scheduletypes.SubjectContent, bool) {
	if subjectID == "" {
		return scheduletypes.SubjectContent{}, false
	}
	return scheduletypes.SubjectContent{
		Index:         subjectID,
		SubjectId:     subjectID,
		Institution:   "institution-1",
		Title:         "Mock Subject",
		Code:          "MOCK101", 
		WorkloadHours: 60,
		Credits:       4,
		SubjectType:   "Required",
		KnowledgeArea: "Computer Science",
		PrerequisiteGroups: []scheduletypes.PrerequisiteGroup{
			{
				GroupType:                "ALL",
				SubjectIds:               []string{"subject-prerequisite-1"},
				MinimumCredits:           4,
				MinimumCompletedSubjects: 1,
			},
		},
		DifficultyLevel: "Medium",
	}, true
}

func (m MockScheduleSubjectKeeper) CheckPrerequisites(ctx sdk.Context, studentID string, subjectID string) (bool, []string, error) {
	// Mock implementation: assume prerequisites are met for testing
	if studentID == "" || subjectID == "" {
		return false, []string{"Invalid student or subject ID"}, nil
	}
	return true, []string{}, nil
}

func (m MockScheduleSubjectKeeper) GetSubjectsByArea(ctx sdk.Context, area string) []scheduletypes.SubjectContent {
	if area == "" {
		return []scheduletypes.SubjectContent{}
	}
	return []scheduletypes.SubjectContent{
		{
			Index:         "subject-area-1",
			SubjectId:     "subject-area-1",
			Institution:   "institution-1",
			Title:         "Mock Subject in " + area,
			Code:          "AREA101",
			WorkloadHours: 60,
			Credits:       4,
			SubjectType:   "Elective",
			KnowledgeArea: area,
			DifficultyLevel: "Easy",
		},
	}
}

// MockScheduleStudentKeeper implements schedule module's StudentKeeper interface
type MockScheduleStudentKeeper struct{}

func (m MockScheduleStudentKeeper) GetAcademicTree(ctx sdk.Context, studentID string, courseID string) (scheduletypes.StudentAcademicTree, bool) {
	if studentID == "" {
		return scheduletypes.StudentAcademicTree{}, false
	}
	return scheduletypes.StudentAcademicTree{
		Index:               fmt.Sprintf("tree_%s_%s", studentID, courseID),
		Student:             studentID,
		Institution:         "institution-1",
		CourseId:            courseID,
		CurriculumVersion:   "v1.0",
		CompletedTokens:     []string{"token-1", "token-2"},
		InProgressTokens:    []string{"token-3"},
		AvailableTokens:     []string{"token-4", "token-5", "token-6"},
		TotalCredits:        60,
		TotalCompletedHours: 900,
		CoefficientGPA:      8.5,
	}, true
}

func (m MockScheduleStudentKeeper) GetCompletedSubjects(ctx sdk.Context, studentID string) []string {
	if studentID == "" {
		return []string{}
	}
	return []string{"subject-1", "subject-2", "subject-completed-3"}
}

func (m MockScheduleStudentKeeper) GetInProgressSubjects(ctx sdk.Context, studentID string) []string {
	if studentID == "" {
		return []string{}
	}
	return []string{"subject-in-progress-1", "subject-in-progress-2"}
}

// MockScheduleCurriculumKeeper implements schedule module's CurriculumKeeper interface
type MockScheduleCurriculumKeeper struct{}

func (m MockScheduleCurriculumKeeper) GetCurriculumTree(ctx sdk.Context, courseID string, version string) (scheduletypes.CurriculumTree, bool) {
	if courseID == "" {
		return scheduletypes.CurriculumTree{}, false
	}
	return scheduletypes.CurriculumTree{
		Index:             fmt.Sprintf("curriculum_%s_%s", courseID, version),
		CourseId:          courseID,
		Version:           version,
		RequiredSubjects:  []string{"subject-1", "subject-2", "subject-3"},
		ElectiveSubjects:  []string{"elective-1", "elective-2", "elective-3"},
		SemesterStructure: []scheduletypes.CurriculumSemester{
			{
				SemesterNumber: 1,
				SubjectIds:     []string{"subject-1", "subject-2"},
			},
			{
				SemesterNumber: 2,
				SubjectIds:     []string{"subject-3", "elective-1"},
			},
		},
		ElectiveGroups: []scheduletypes.ElectiveGroup{
			{
				GroupId:             "elective-group-1",
				Name:                "Computer Science Electives",
				SubjectIds:          []string{"elective-1", "elective-2"},
				MinSubjectsRequired: 1,
				CreditsRequired:     4,
				KnowledgeArea:       "Computer Science",
			},
		},
	}, true
}

func (m MockScheduleCurriculumKeeper) GetCurrentCurriculumVersion(ctx sdk.Context, courseID string) string {
	if courseID == "" {
		return ""
	}
	return "v1.0"
}
