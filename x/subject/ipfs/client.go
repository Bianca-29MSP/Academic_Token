package ipfs

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// IPFSClient provides functionality to interact with IPFS
type IPFSClient struct {
	ApiEndpoint string
	LocalPath   string
	UseHTTP     bool
	Timeout     time.Duration
}

// NewIPFSClient creates a new IPFS client
func NewIPFSClient(apiEndpoint, localPath string, useHTTP bool) *IPFSClient {
	return &IPFSClient{
		ApiEndpoint: apiEndpoint,
		LocalPath:   localPath,
		UseHTTP:     useHTTP,
		Timeout:     time.Second * 30,
	}
}

// Add adds content to IPFS and returns its hash and link
func (c *IPFSClient) Add(content []byte) (string, string, error) {
	// Calculate content hash before sending to IPFS
	hash := sha256.Sum256(content)
	contentHash := hex.EncodeToString(hash[:])

	var ipfsHash string
	var err error

	if c.UseHTTP {
		ipfsHash, err = c.addViaHTTP(content)
	} else {
		ipfsHash, err = c.addViaLocalFS(content)
	}

	if err != nil {
		return "", "", err
	}

	ipfsLink := fmt.Sprintf("ipfs://%s", ipfsHash)
	return contentHash, ipfsLink, nil
}

// Get retrieves content from IPFS by its hash
func (c *IPFSClient) Get(ipfsLink string) ([]byte, error) {
	// Extract hash from IPFS link (format ipfs://HASH)
	var ipfsHash string
	if len(ipfsLink) > 7 && ipfsLink[:7] == "ipfs://" {
		ipfsHash = ipfsLink[7:]
	} else {
		ipfsHash = ipfsLink
	}

	if c.UseHTTP {
		return c.getViaHTTP(ipfsHash)
	}
	return c.getViaLocalFS(ipfsHash)
}

// addViaHTTP adds content to IPFS via HTTP API
func (c *IPFSClient) addViaHTTP(content []byte) (string, error) {
	url := fmt.Sprintf("%s/api/v0/add", c.ApiEndpoint)

	// Create multipart/form-data request
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "content")
	if err != nil {
		return "", fmt.Errorf("error creating form-data: %v", err)
	}

	_, err = part.Write(content)
	if err != nil {
		return "", fmt.Errorf("error writing content: %v", err)
	}

	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("error finalizing form-data: %v", err)
	}

	// Send request
	client := &http.Client{
		Timeout: c.Timeout,
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("IPFS server error: %s", resp.Status)
	}

	// Read and parse JSON response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	var result struct {
		Hash string `json:"hash"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("error parsing response: %v", err)
	}

	return result.Hash, nil
}

// getViaHTTP retrieves content from IPFS via HTTP API
func (c *IPFSClient) getViaHTTP(ipfsHash string) ([]byte, error) {
	url := fmt.Sprintf("%s/api/v0/cat?arg=%s", c.ApiEndpoint, ipfsHash)

	client := &http.Client{
		Timeout: c.Timeout,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error requesting content: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("IPFS server error: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}

// addViaLocalFS adds content to local filesystem (for development/testing)
func (c *IPFSClient) addViaLocalFS(content []byte) (string, error) {
	// Generate unique hash for the file
	hash := sha256.Sum256(content)
	ipfsHash := hex.EncodeToString(hash[:])

	// Create directory if it doesn't exist
	err := os.MkdirAll(c.LocalPath, 0755)
	if err != nil {
		return "", fmt.Errorf("error creating directory: %v", err)
	}

	// Save content to a local file
	filePath := filepath.Join(c.LocalPath, ipfsHash)
	err = os.WriteFile(filePath, content, 0644)
	if err != nil {
		return "", fmt.Errorf("error saving content: %v", err)
	}

	return ipfsHash, nil
}

// getViaLocalFS retrieves content from local filesystem
func (c *IPFSClient) getViaLocalFS(ipfsHash string) ([]byte, error) {
	filePath := filepath.Join(c.LocalPath, ipfsHash)

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file content: %v", err)
	}

	return content, nil
}

// MockIPFSClient is a mock client for testing
type MockIPFSClient struct {
	storage   map[string][]byte
	addCalled bool
	getCalled bool
}

// NewMockIPFSClient creates a new mock IPFS client for testing
func NewMockIPFSClient() *MockIPFSClient {
	return &MockIPFSClient{
		storage: make(map[string][]byte),
	}
}

// Add mocks adding content to IPFS
func (m *MockIPFSClient) Add(content []byte) (string, string, error) {
	m.addCalled = true
	hash := sha256.Sum256(content)
	contentHash := hex.EncodeToString(hash[:])
	ipfsHash := contentHash[0:16] // Simplified for testing
	ipfsLink := fmt.Sprintf("ipfs://%s", ipfsHash)

	m.storage[ipfsHash] = content
	return contentHash, ipfsLink, nil
}

// Get mocks retrieving content from IPFS
func (m *MockIPFSClient) Get(ipfsLink string) ([]byte, error) {
	m.getCalled = true
	var ipfsHash string
	if len(ipfsLink) > 7 && ipfsLink[:7] == "ipfs://" {
		ipfsHash = ipfsLink[7:]
	} else {
		ipfsHash = ipfsLink
	}

	content, exists := m.storage[ipfsHash]
	if !exists {
		return nil, fmt.Errorf("content not found: %s", ipfsHash)
	}

	return content, nil
}

// WasAddCalled returns whether Add was called
func (m *MockIPFSClient) WasAddCalled() bool {
	return m.addCalled
}

// WasGetCalled returns whether Get was called
func (m *MockIPFSClient) WasGetCalled() bool {
	return m.getCalled
}
