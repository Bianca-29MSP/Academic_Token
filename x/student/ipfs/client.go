package ipfs

import (
	"encoding/json"
	"fmt"
	subjectipfs "academictoken/x/subject/ipfs"
)

// StudentIPFSClient wraps the base IPFS client with student-specific functionality
type StudentIPFSClient struct {
	*subjectipfs.IPFSClient
}

// NewStudentIPFSClient creates a new IPFS client for student module
func NewStudentIPFSClient(apiEndpoint, localPath string, useHTTP bool) *StudentIPFSClient {
	return &StudentIPFSClient{
		IPFSClient: subjectipfs.NewIPFSClient(apiEndpoint, localPath, useHTTP),
	}
}

// AddStudentData adds student-specific structured data to IPFS
func (c *StudentIPFSClient) AddStudentData(data interface{}) (string, string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal student data: %w", err)
	}
	return c.Add(jsonData)
}

// GetStudentData retrieves and unmarshals student data from IPFS
func (c *StudentIPFSClient) GetStudentData(ipfsLink string, target interface{}) error {
	data, err := c.Get(ipfsLink)
	if err != nil {
		return fmt.Errorf("failed to get student data from IPFS: %w", err)
	}
	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("failed to unmarshal student data: %w", err)
	}
	return nil
}

// AddAcademicHistory adds academic history to IPFS
func (c *StudentIPFSClient) AddAcademicHistory(history map[string]interface{}) (string, string, error) {
	return c.AddStudentData(history)
}

// GetAcademicHistory retrieves academic history from IPFS
func (c *StudentIPFSClient) GetAcademicHistory(ipfsLink string) (map[string]interface{}, error) {
	var history map[string]interface{}
	err := c.GetStudentData(ipfsLink, &history)
	return history, err
}

// AddTranscriptData adds transcript data to IPFS
func (c *StudentIPFSClient) AddTranscriptData(transcript interface{}) (string, string, error) {
	return c.AddStudentData(transcript)
}

// GetTranscriptData retrieves transcript data from IPFS
func (c *StudentIPFSClient) GetTranscriptData(ipfsLink string) (map[string]interface{}, error) {
	var transcript map[string]interface{}
	err := c.GetStudentData(ipfsLink, &transcript)
	return transcript, err
}

// AddEnrollmentHistory adds enrollment history to IPFS
func (c *StudentIPFSClient) AddEnrollmentHistory(history interface{}) (string, string, error) {
	return c.AddStudentData(history)
}

// GetEnrollmentHistory retrieves enrollment history from IPFS
func (c *StudentIPFSClient) GetEnrollmentHistory(ipfsLink string) (map[string]interface{}, error) {
	var history map[string]interface{}
	err := c.GetStudentData(ipfsLink, &history)
	return history, err
}

// AddPersonalDocuments adds personal documents to IPFS
func (c *StudentIPFSClient) AddPersonalDocuments(documents interface{}) (string, string, error) {
	return c.AddStudentData(documents)
}

// GetPersonalDocuments retrieves personal documents from IPFS
func (c *StudentIPFSClient) GetPersonalDocuments(ipfsLink string) (map[string]interface{}, error) {
	var documents map[string]interface{}
	err := c.GetStudentData(ipfsLink, &documents)
	return documents, err
}

// MockStudentIPFSClient for testing
type MockStudentIPFSClient struct {
	*subjectipfs.MockIPFSClient
}

// NewMockStudentIPFSClient creates a new mock IPFS client for testing
func NewMockStudentIPFSClient() *MockStudentIPFSClient {
	return &MockStudentIPFSClient{
		MockIPFSClient: subjectipfs.NewMockIPFSClient(),
	}
}

// AddStudentData mocks adding student data
func (m *MockStudentIPFSClient) AddStudentData(data interface{}) (string, string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", "", err
	}
	return m.Add(jsonData)
}

// GetStudentData mocks getting student data
func (m *MockStudentIPFSClient) GetStudentData(ipfsLink string, target interface{}) error {
	data, err := m.Get(ipfsLink)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}
