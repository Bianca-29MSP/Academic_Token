package ipfs

import (
	"encoding/json"
	"fmt"
	subjectipfs "academictoken/x/subject/ipfs"
)

// TokenDefIPFSClient wraps the base IPFS client with tokendef-specific functionality
type TokenDefIPFSClient struct {
	*subjectipfs.IPFSClient
}

// NewTokenDefIPFSClient creates a new IPFS client for tokendef module
func NewTokenDefIPFSClient(apiEndpoint, localPath string, useHTTP bool) *TokenDefIPFSClient {
	return &TokenDefIPFSClient{
		IPFSClient: subjectipfs.NewIPFSClient(apiEndpoint, localPath, useHTTP),
	}
}

// AddTokenDefData adds tokendef-specific structured data to IPFS
func (c *TokenDefIPFSClient) AddTokenDefData(data interface{}) (string, string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal tokendef data: %w", err)
	}
	return c.Add(jsonData)
}

// GetTokenDefData retrieves and unmarshals tokendef data from IPFS
func (c *TokenDefIPFSClient) GetTokenDefData(ipfsLink string, target interface{}) error {
	data, err := c.Get(ipfsLink)
	if err != nil {
		return fmt.Errorf("failed to get tokendef data from IPFS: %w", err)
	}
	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("failed to unmarshal tokendef data: %w", err)
	}
	return nil
}

// AddTokenMetadata adds token metadata to IPFS
func (c *TokenDefIPFSClient) AddTokenMetadata(metadata map[string]interface{}) (string, string, error) {
	return c.AddTokenDefData(metadata)
}

// GetTokenMetadata retrieves token metadata from IPFS
func (c *TokenDefIPFSClient) GetTokenMetadata(ipfsLink string) (map[string]interface{}, error) {
	var metadata map[string]interface{}
	err := c.GetTokenDefData(ipfsLink, &metadata)
	return metadata, err
}

// AddTokenAttributes adds token attributes to IPFS
func (c *TokenDefIPFSClient) AddTokenAttributes(attributes interface{}) (string, string, error) {
	return c.AddTokenDefData(attributes)
}

// GetTokenAttributes retrieves token attributes from IPFS
func (c *TokenDefIPFSClient) GetTokenAttributes(ipfsLink string) (map[string]interface{}, error) {
	var attributes map[string]interface{}
	err := c.GetTokenDefData(ipfsLink, &attributes)
	return attributes, err
}

// MockTokenDefIPFSClient for testing
type MockTokenDefIPFSClient struct {
	*subjectipfs.MockIPFSClient
}

// NewMockTokenDefIPFSClient creates a new mock IPFS client for testing
func NewMockTokenDefIPFSClient() *MockTokenDefIPFSClient {
	return &MockTokenDefIPFSClient{
		MockIPFSClient: subjectipfs.NewMockIPFSClient(),
	}
}

// AddTokenDefData mocks adding tokendef data
func (m *MockTokenDefIPFSClient) AddTokenDefData(data interface{}) (string, string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", "", err
	}
	return m.Add(jsonData)
}

// GetTokenDefData mocks getting tokendef data
func (m *MockTokenDefIPFSClient) GetTokenDefData(ipfsLink string, target interface{}) error {
	data, err := m.Get(ipfsLink)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}
