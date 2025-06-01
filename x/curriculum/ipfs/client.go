package ipfs

import (
	"encoding/json"
	"fmt"

	subjectipfs "academictoken/x/subject/ipfs"
)

// CurriculumIPFSClient wraps the base IPFS client with curriculum-specific functionality
type CurriculumIPFSClient struct {
	*subjectipfs.IPFSClient
}

// NewCurriculumIPFSClient creates a new IPFS client for curriculum module
func NewCurriculumIPFSClient(apiEndpoint, localPath string, useHTTP bool) *CurriculumIPFSClient {
	return &CurriculumIPFSClient{
		IPFSClient: subjectipfs.NewIPFSClient(apiEndpoint, localPath, useHTTP),
	}
}

// AddCurriculumData adds curriculum-specific structured data to IPFS
// This method can handle curriculum metadata and convert it to JSON before storing
func (c *CurriculumIPFSClient) AddCurriculumData(data interface{}) (string, string, error) {
	// Convert curriculum data to JSON bytes
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal curriculum data: %w", err)
	}

	// Use the base Add method
	return c.Add(jsonData)
}

// GetCurriculumData retrieves and unmarshals curriculum data from IPFS
func (c *CurriculumIPFSClient) GetCurriculumData(ipfsLink string, target interface{}) error {
	// Get raw data from IPFS
	data, err := c.Get(ipfsLink)
	if err != nil {
		return fmt.Errorf("failed to get curriculum data from IPFS: %w", err)
	}

	// Unmarshal into target structure
	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("failed to unmarshal curriculum data: %w", err)
	}

	return nil
}

// MockCurriculumIPFSClient for testing
type MockCurriculumIPFSClient struct {
	*subjectipfs.MockIPFSClient
}

// NewMockCurriculumIPFSClient creates a new mock IPFS client for testing
func NewMockCurriculumIPFSClient() *MockCurriculumIPFSClient {
	return &MockCurriculumIPFSClient{
		MockIPFSClient: subjectipfs.NewMockIPFSClient(),
	}
}
