package ipfs

import (
	"encoding/json"
	"fmt"
	subjectipfs "academictoken/x/subject/ipfs"
)

// AcademicNFTIPFSClient wraps the base IPFS client with academicnft-specific functionality
type AcademicNFTIPFSClient struct {
	*subjectipfs.IPFSClient
}

// NewAcademicNFTIPFSClient creates a new IPFS client for academicnft module
func NewAcademicNFTIPFSClient(apiEndpoint, localPath string, useHTTP bool) *AcademicNFTIPFSClient {
	return &AcademicNFTIPFSClient{
		IPFSClient: subjectipfs.NewIPFSClient(apiEndpoint, localPath, useHTTP),
	}
}

// AddAcademicNFTData adds academicnft-specific structured data to IPFS
func (c *AcademicNFTIPFSClient) AddAcademicNFTData(data interface{}) (string, string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal academicnft data: %w", err)
	}
	return c.Add(jsonData)
}

// GetAcademicNFTData retrieves and unmarshals academicnft data from IPFS
func (c *AcademicNFTIPFSClient) GetAcademicNFTData(ipfsLink string, target interface{}) error {
	data, err := c.Get(ipfsLink)
	if err != nil {
		return fmt.Errorf("failed to get academicnft data from IPFS: %w", err)
	}
	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("failed to unmarshal academicnft data: %w", err)
	}
	return nil
}

// AddNFTMetadata adds NFT metadata to IPFS
func (c *AcademicNFTIPFSClient) AddNFTMetadata(metadata map[string]interface{}) (string, string, error) {
	return c.AddAcademicNFTData(metadata)
}

// GetNFTMetadata retrieves NFT metadata from IPFS
func (c *AcademicNFTIPFSClient) GetNFTMetadata(ipfsLink string) (map[string]interface{}, error) {
	var metadata map[string]interface{}
	err := c.GetAcademicNFTData(ipfsLink, &metadata)
	return metadata, err
}

// AddCertificateData adds certificate data to IPFS
func (c *AcademicNFTIPFSClient) AddCertificateData(certificate interface{}) (string, string, error) {
	return c.AddAcademicNFTData(certificate)
}

// GetCertificateData retrieves certificate data from IPFS
func (c *AcademicNFTIPFSClient) GetCertificateData(ipfsLink string) (map[string]interface{}, error) {
	var certificate map[string]interface{}
	err := c.GetAcademicNFTData(ipfsLink, &certificate)
	return certificate, err
}

// AddCompletionDetails adds completion details to IPFS
func (c *AcademicNFTIPFSClient) AddCompletionDetails(details interface{}) (string, string, error) {
	return c.AddAcademicNFTData(details)
}

// GetCompletionDetails retrieves completion details from IPFS
func (c *AcademicNFTIPFSClient) GetCompletionDetails(ipfsLink string) (map[string]interface{}, error) {
	var details map[string]interface{}
	err := c.GetAcademicNFTData(ipfsLink, &details)
	return details, err
}

// MockAcademicNFTIPFSClient for testing
type MockAcademicNFTIPFSClient struct {
	*subjectipfs.MockIPFSClient
}

// NewMockAcademicNFTIPFSClient creates a new mock IPFS client for testing
func NewMockAcademicNFTIPFSClient() *MockAcademicNFTIPFSClient {
	return &MockAcademicNFTIPFSClient{
		MockIPFSClient: subjectipfs.NewMockIPFSClient(),
	}
}

// AddAcademicNFTData mocks adding academicnft data
func (m *MockAcademicNFTIPFSClient) AddAcademicNFTData(data interface{}) (string, string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", "", err
	}
	return m.Add(jsonData)
}

// GetAcademicNFTData mocks getting academicnft data
func (m *MockAcademicNFTIPFSClient) GetAcademicNFTData(ipfsLink string, target interface{}) error {
	data, err := m.Get(ipfsLink)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}
