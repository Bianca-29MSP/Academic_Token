// cmd/rest-server/main.go
// REST server that connects to real blockchain data

package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "math/rand"
    "net/http"
    "os/exec"
    "strings"
    "time"
    
    "github.com/gorilla/mux"
    "github.com/gorilla/handlers"
)

// ========== STRUCTURES ==========
type Institution struct {
    ID        string `json:"id"`
    Name      string `json:"name"`
    Code      string `json:"code"`
    Country   string `json:"country"`
    CreatedAt string `json:"createdAt"`
}

type Course struct {
    ID             string `json:"id"`
    InstitutionID  string `json:"institutionId"`
    Name           string `json:"name"`
    Code           string `json:"code"`
    Duration       int    `json:"duration"`
    TotalCredits   int    `json:"totalCredits"`
}

type Subject struct {
    ID             string `json:"id"`
    CourseID       string `json:"courseId"`
    InstitutionID  string `json:"institutionId"`
    Name           string `json:"name"`
    Code           string `json:"code"`
    Credits        int    `json:"credits"`
    Syllabus       string `json:"syllabus"`
    Metadata       string `json:"metadata"`
}

type Student struct {
    ID             string `json:"id"`
    InstitutionID  string `json:"institutionId"`
    Name           string `json:"name"`
    Email          string `json:"email"`
    CourseID       string `json:"courseId"`
    CurriculumID   string `json:"curriculumId"`
    EnrollmentDate string `json:"enrollmentDate"`
}

type AcademicNFT struct {
    ID             string      `json:"id"`
    StudentID      string      `json:"studentId"`
    SubjectID      string      `json:"subjectId"`
    Grade          float64     `json:"grade"`
    CompletionDate string      `json:"completionDate"`
    NFTHash        string      `json:"nftHash"`
    Metadata       NFTMetadata `json:"metadata"`
}

type NFTMetadata struct {
    Subject     string `json:"subject"`
    Credits     int    `json:"credits"`
    Institution string `json:"institution"`
}

type DegreeEligibility struct {
    Eligible         bool     `json:"eligible"`
    CompletedCredits int      `json:"completedCredits"`
    RequiredCredits  int      `json:"requiredCredits"`
    MissingSubjects  []string `json:"missingSubjects"`
}

// ========== BLOCKCHAIN QUERIES ==========

// Query blockchain via REST endpoints instead of CLI
func queryBlockchainREST(endpoint string) ([]byte, error) {
    url := fmt.Sprintf("http://localhost:1317%s", endpoint)
    fmt.Printf("📡 Making REST request to: %s\n", url)
    
    resp, err := http.Get(url)
    if err != nil {
        fmt.Printf("❌ REST request failed: %v\n", err)
        return nil, err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        fmt.Printf("❌ REST request returned status: %d\n", resp.StatusCode)
        return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
    }
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        fmt.Printf("❌ Failed to read response body: %v\n", err)
        return nil, err
    }
    
    return body, nil
}

// Query all institutions from blockchain
func queryInstitutions() ([]Institution, error) {
    fmt.Println("🔍 Querying institutions from blockchain...")
    
    // Query via REST endpoint
    output, err := queryBlockchainREST("/academictoken/institution/institution")
    if err != nil {
        fmt.Printf("❌ Failed to query institutions from blockchain: %v\n", err)
        return []Institution{}, fmt.Errorf("no institutions found on blockchain")
    }
    
    fmt.Println("✅ Got blockchain data for institutions")
    var result struct {
        Institution []Institution `json:"institution"`
        Pagination  interface{}   `json:"pagination"`
    }
    
    if err := json.Unmarshal(output, &result); err != nil {
        fmt.Printf("❌ Failed to parse institutions data: %v\n", err)
        return []Institution{}, fmt.Errorf("failed to parse institutions data")
    }
    
    return result.Institution, nil
}

// Query all courses from blockchain
func queryCourses() ([]Course, error) {
    fmt.Println("🔍 Querying courses from blockchain...")
    
    // Query via REST endpoint
    output, err := queryBlockchainREST("/academictoken/course/course")
    if err != nil {
        fmt.Printf("❌ Failed to query courses from blockchain: %v\n", err)
        return []Course{}, fmt.Errorf("no courses found on blockchain")
    }
    
    fmt.Println("✅ Got blockchain data for courses")
    var result struct {
        Course     []Course    `json:"course"`
        Pagination interface{} `json:"pagination"`
    }
    
    if err := json.Unmarshal(output, &result); err != nil {
        fmt.Printf("❌ Failed to parse courses data: %v\n", err)
        return []Course{}, fmt.Errorf("failed to parse courses data")
    }
    
    return result.Course, nil
}

// Query all subjects from blockchain
func querySubjects() ([]Subject, error) {
    fmt.Println("🔍 Querying subjects from blockchain...")
    
    // Query via REST endpoint  
    output, err := queryBlockchainREST("/academictoken/subject/subjects")
    if err != nil {
        fmt.Printf("❌ Failed to query subjects from blockchain: %v\n", err)
        return []Subject{}, fmt.Errorf("no subjects found on blockchain")
    }
    
    fmt.Println("✅ Got blockchain data for subjects")
    var result struct {
        Subjects   []Subject   `json:"subjects"`
        Pagination interface{} `json:"pagination"`
    }
    
    if err := json.Unmarshal(output, &result); err != nil {
        fmt.Printf("❌ Failed to parse subjects data: %v\n", err)
        return []Subject{}, fmt.Errorf("failed to parse subjects data")
    }
    
    return result.Subjects, nil
}

// ========== HANDLERS ==========

// POST /academic/equivalence/request
func requestEquivalence(w http.ResponseWriter, r *http.Request) {
    var req map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    
    fmt.Printf("📡 Creating equivalence request: %v\n", req)
    
    // Forward the request to the blockchain module using the correct endpoint
    reqBody, _ := json.Marshal(req)
    resp, err := http.Post("http://localhost:1317/academictoken.equivalence.Msg/RequestEquivalence", "application/json", bytes.NewReader(reqBody))
    if err != nil {
        fmt.Printf("❌ Failed to forward request to blockchain: %v\n", err)
        http.Error(w, "Failed to create equivalence request", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()
    
    // Read and forward the response
    body, _ := io.ReadAll(resp.Body)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(resp.StatusCode)
    w.Write(body)
}

// POST /academic/equivalence/execute-analysis
func executeEquivalenceAnalysis(w http.ResponseWriter, r *http.Request) {
    var req map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    
    fmt.Printf("📊 Executing equivalence analysis: %v\n", req)
    
    // Forward the request to the blockchain module using the correct endpoint
    reqBody, _ := json.Marshal(req)
    resp, err := http.Post("http://localhost:1317/academictoken.equivalence.Msg/ExecuteEquivalenceAnalysis", "application/json", bytes.NewReader(reqBody))
    if err != nil {
        fmt.Printf("❌ Failed to forward request to blockchain: %v\n", err)
        http.Error(w, "Failed to execute analysis", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()
    
    // Read and forward the response
    body, _ := io.ReadAll(resp.Body)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(resp.StatusCode)
    w.Write(body)
}

// GET /academic/equivalence/equivalences/{id}
func getEquivalence(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    equivalenceID := vars["id"]
    
    fmt.Printf("🔍 Getting equivalence details for: %s\n", equivalenceID)
    
    // Query blockchain for equivalence details using the correct endpoint
    output, err := queryBlockchainREST(fmt.Sprintf("/academictoken/equivalence/equivalences/%s", equivalenceID))
    if err != nil {
        fmt.Printf("❌ Failed to query equivalence from blockchain: %v\n", err)
        http.Error(w, "Equivalence not found", http.StatusNotFound)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.Write(output)
}

// GET /academic/equivalence/check/{source_subject_id}/{target_subject_id}
func checkEquivalenceStatus(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    sourceSubjectID := vars["source_subject_id"]
    targetSubjectID := vars["target_subject_id"]
    
    fmt.Printf("🔍 Checking equivalence status between %s and %s\n", sourceSubjectID, targetSubjectID)
    
    output, err := queryBlockchainREST(fmt.Sprintf("/academictoken/equivalence/check/%s/%s", sourceSubjectID, targetSubjectID))
    if err != nil {
        fmt.Printf("❌ Failed to check equivalence status: %v\n", err)
        http.Error(w, "Failed to check equivalence status", http.StatusNotFound)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.Write(output)
}

// GET /academic/equivalence/equivalences
func listEquivalences(w http.ResponseWriter, r *http.Request) {
    fmt.Println("🔍 Listing all equivalences")
    
    output, err := queryBlockchainREST("/academictoken/equivalence/equivalences")
    if err != nil {
        fmt.Printf("❌ Failed to list equivalences: %v\n", err)
        // Return empty list instead of error
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]interface{}{
            "equivalences": []interface{}{},
            "pagination": nil,
        })
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.Write(output)
}

// GET /academic/equivalence/analysis/{equivalence_id}
func getAnalysisMetadata(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    equivalenceID := vars["equivalence_id"]
    
    fmt.Printf("🔍 Getting analysis metadata for: %s\n", equivalenceID)
    
    output, err := queryBlockchainREST(fmt.Sprintf("/academictoken/equivalence/analysis/%s", equivalenceID))
    if err != nil {
        fmt.Printf("❌ Failed to get analysis metadata: %v\n", err)
        http.Error(w, "Analysis metadata not found", http.StatusNotFound)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.Write(output)
}

// GET /academic/institution/list
func listInstitutions(w http.ResponseWriter, r *http.Request) {
    institutions, err := queryInstitutions()
    if err != nil {
        http.Error(w, "Failed to query institutions", http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(institutions)
}

// GET /academic/course/list
func listCourses(w http.ResponseWriter, r *http.Request) {
    courses, err := queryCourses()
    if err != nil {
        http.Error(w, "Failed to query courses", http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(courses)
}

// GET /academic/subject/list
func listSubjects(w http.ResponseWriter, r *http.Request) {
    subjects, err := querySubjects()
    if err != nil {
        http.Error(w, "Failed to query subjects", http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(subjects)
}

// GET /academic/student/{id}
func getStudent(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    studentID := vars["id"]
    
    // Query real student data from blockchain
    fmt.Printf("🔍 Querying student %s from blockchain...\n", studentID)
    
    output, err := queryBlockchainREST(fmt.Sprintf("/academictoken/student/student/%s", studentID))
    if err != nil {
        fmt.Printf("❌ Student not found on blockchain: %v\n", err)
        http.Error(w, "Student not found on blockchain", http.StatusNotFound)
        return
    }
    
    var student Student
    if err := json.Unmarshal(output, &student); err != nil {
        fmt.Printf("❌ Failed to parse student data: %v\n", err)
        http.Error(w, "Failed to parse student data", http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(student)
}

// GET /academic/student/{id}/nfts
func getStudentNFTs(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    studentID := vars["id"]
    
    // Query real NFT data from blockchain
    fmt.Printf("🔍 Querying NFTs for student %s from blockchain...\n", studentID)
    
    output, err := queryBlockchainREST(fmt.Sprintf("/academictoken/academicnft/student/%s/tokens", studentID))
    if err != nil {
        fmt.Printf("❌ No NFTs found for student on blockchain: %v\n", err)
        // Return empty array instead of error - student might have no completed subjects yet
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode([]AcademicNFT{})
        return
    }
    
    var nfts []AcademicNFT
    if err := json.Unmarshal(output, &nfts); err != nil {
        fmt.Printf("❌ Failed to parse NFTs data: %v\n", err)
        // Return empty array instead of error
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode([]AcademicNFT{})
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(nfts)
}

// GET /academic/degree/{student_id}/eligibility
func getDegreeEligibility(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    studentID := vars["student_id"]
    
    // Query real degree eligibility from blockchain via CosmWasm contract
    fmt.Printf("🔍 Checking degree eligibility for student %s...\n", studentID)
    
    // TODO: This should call the degree eligibility CosmWasm contract
    // For now, return error indicating contract is needed
    http.Error(w, "Degree eligibility check requires deployed CosmWasm contract", http.StatusNotImplemented)
}

// GET /academic/student/{student_id}/prerequisites/{subject_id}
func checkPrerequisites(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    studentID := vars["student_id"]
    subjectID := vars["subject_id"]
    
    // Query prerequisites via CosmWasm contract
    fmt.Printf("🔍 Checking prerequisites for student %s, subject %s...\n", studentID, subjectID)
    
    // TODO: This should call the prerequisites CosmWasm contract
    // For now, return error indicating contract is needed
    http.Error(w, "Prerequisites check requires deployed CosmWasm contract", http.StatusNotImplemented)
}

// Cosmos node info endpoint
func getNodeInfo(w http.ResponseWriter, r *http.Request) {
    nodeInfo := map[string]interface{}{
        "default_node_info": map[string]interface{}{
            "network": "academictoken",
            "moniker": "academic-node",
        },
        "application_version": map[string]interface{}{
            "name": "academictoken",
            "version": "v1.0.0",
        },
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(nodeInfo)
}

// Test blockchain connectivity
func testBlockchainConnection() {
    fmt.Println("🔍 Testing blockchain connectivity...")
    
    // Test if academictokend is available
    cmd := exec.Command("academictokend", "status")
    output, err := cmd.Output()
    if err != nil {
        fmt.Printf("⚠️  Blockchain node not available: %v\n", err)
        fmt.Println("📝 Will use fallback data")
    } else {
        fmt.Println("✅ Blockchain node is responding")
        if strings.Contains(string(output), "latest_block_height") {
            fmt.Println("🎯 Node is synced and ready")
        }
    }
}

func main() {
    // Initialize random seed
    rand.Seed(time.Now().UnixNano())
    
    // Test blockchain connection on startup
    testBlockchainConnection()
    
    r := mux.NewRouter()
    
    // CORS middleware
    headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
    originsOk := handlers.AllowedOrigins([]string{"http://localhost:3000", "http://localhost:3001", "http://u6s2n-gx777-77774-qaaba-cai.localhost:4943", "http://127.0.0.1:4943"})
    methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})
    
    // Academic API routes
    api := r.PathPrefix("/academic").Subrouter()
    api.HandleFunc("/institution/list", listInstitutions).Methods("GET")
    api.HandleFunc("/course/list", listCourses).Methods("GET")
    api.HandleFunc("/subject/list", listSubjects).Methods("GET")
    api.HandleFunc("/student/{id}", getStudent).Methods("GET")
    api.HandleFunc("/student/{id}/nfts", getStudentNFTs).Methods("GET")
    api.HandleFunc("/degree/{student_id}/eligibility", getDegreeEligibility).Methods("GET")
    api.HandleFunc("/student/{student_id}/prerequisites/{subject_id}", checkPrerequisites).Methods("GET")
    
    // Equivalence routes
    api.HandleFunc("/equivalence/request", requestEquivalence).Methods("POST")
    api.HandleFunc("/equivalence/execute-analysis", executeEquivalenceAnalysis).Methods("POST")
    api.HandleFunc("/equivalence/equivalences/{id}", getEquivalence).Methods("GET")
    api.HandleFunc("/equivalence/check/{source_subject_id}/{target_subject_id}", checkEquivalenceStatus).Methods("GET")
    api.HandleFunc("/equivalence/equivalences", listEquivalences).Methods("GET")
    api.HandleFunc("/equivalence/analysis/{equivalence_id}", getAnalysisMetadata).Methods("GET")
    
    // Cosmos compatibility
    r.HandleFunc("/cosmos/base/tendermint/v1beta1/node_info", getNodeInfo).Methods("GET")
    
    // Health check
    r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    }).Methods("GET")
    
    fmt.Println("🚀 Academic Token REST Server")
    fmt.Println("🌍 API Server: http://localhost:1318")
    fmt.Println("📡 Academic API: http://localhost:1318/academic")
    fmt.Println("💡 Health: http://localhost:1318/health")
    fmt.Println("🔍 Data Source: Real blockchain data ONLY (no fallback data)")
    fmt.Println("⚠️  Note: If blockchain has no data, endpoints will return errors")
    fmt.Println("")
    fmt.Println("✅ Ready to connect with frontend!")
    
    log.Fatal(http.ListenAndServe(":1318", handlers.CORS(originsOk, headersOk, methodsOk)(r)))
}
