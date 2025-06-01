package ipfs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client represents an IPFS HTTP client for prerequisites contract
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a new IPFS client instance for prerequisites
func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// PrerequisiteIPFSData represents the structure stored in IPFS for prerequisites
type PrerequisiteIPFSData struct {
	SubjectID     string                    `json:"subject_id"`
	Title         string                    `json:"title"`
	Code          string                    `json:"code"`
	Description   string                    `json:"description"`
	Learning      []string                  `json:"learning_objectives"`
	Methodology   string                    `json:"methodology"`
	Topics        []string                  `json:"topics"`
	Bibliography  []string                  `json:"bibliography"`
	Prerequisites []PrerequisiteRequirement `json:"prerequisites"`
	CreatedAt     string                    `json:"created_at"`
	UpdatedAt     string                    `json:"updated_at"`
}

// PrerequisiteRequirement represents a single prerequisite requirement
type PrerequisiteRequirement struct {
	Type     string   `json:"type"`      // "ALL" or "ANY"
	Subjects []string `json:"subjects"`  // List of subject IDs
	MinGrade float64  `json:"min_grade"` // Minimum grade required
	Required bool     `json:"required"`  // Whether this prerequisite is mandatory
}

// AddPrerequisiteData stores prerequisite data in IPFS
func (c *Client) AddPrerequisiteData(data PrerequisiteIPFSData) (string, error) {
	// Serialize data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal prerequisite data: %w", err)
	}

	// Create multipart form data
	var buf bytes.Buffer
	buf.WriteString("--boundary\r\n")
	buf.WriteString("Content-Disposition: form-data; name=\"file\"; filename=\"prerequisite.json\"\r\n")
	buf.WriteString("Content-Type: application/json\r\n\r\n")
	buf.Write(jsonData)
	buf.WriteString("\r\n--boundary--\r\n")

	// Create HTTP request
	req, err := http.NewRequest("POST", c.BaseURL+"/api/v0/add", &buf)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "multipart/form-data; boundary=boundary")

	// Execute request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("IPFS request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var result struct {
		Hash string `json:"Hash"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Hash, nil
}

// GetPrerequisiteData retrieves prerequisite data from IPFS
func (c *Client) GetPrerequisiteData(hash string) (*PrerequisiteIPFSData, error) {
	// Create HTTP request
	req, err := http.NewRequest("POST", c.BaseURL+"/api/v0/cat?arg="+hash, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Execute request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("IPFS request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var data PrerequisiteIPFSData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode prerequisite data: %w", err)
	}

	return &data, nil
}

// ValidatePrerequisiteHash verifies if the hash corresponds to valid prerequisite data
func (c *Client) ValidatePrerequisiteHash(hash string) (bool, error) {
	data, err := c.GetPrerequisiteData(hash)
	if err != nil {
		return false, err
	}

	// Basic validation - check if essential fields are present
	if data.SubjectID == "" || data.Title == "" || data.Code == "" {
		return false, fmt.Errorf("invalid prerequisite data: missing essential fields")
	}

	return true, nil
}

// UpdatePrerequisiteData updates existing prerequisite data in IPFS
// Returns the new hash
func (c *Client) UpdatePrerequisiteData(existingHash string, updates PrerequisiteIPFSData) (string, error) {
	// Get existing data
	existingData, err := c.GetPrerequisiteData(existingHash)
	if err != nil {
		return "", fmt.Errorf("failed to get existing data: %w", err)
	}

	// Merge updates with existing data
	if updates.Title != "" {
		existingData.Title = updates.Title
	}
	if updates.Description != "" {
		existingData.Description = updates.Description
	}
	if len(updates.Learning) > 0 {
		existingData.Learning = updates.Learning
	}
	if updates.Methodology != "" {
		existingData.Methodology = updates.Methodology
	}
	if len(updates.Topics) > 0 {
		existingData.Topics = updates.Topics
	}
	if len(updates.Bibliography) > 0 {
		existingData.Bibliography = updates.Bibliography
	}
	if len(updates.Prerequisites) > 0 {
		existingData.Prerequisites = updates.Prerequisites
	}

	// Update timestamp
	existingData.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	// Store updated data
	return c.AddPrerequisiteData(*existingData)
}

// CheckPrerequisiteCompletion verifies if prerequisites are met based on IPFS data
func (c *Client) CheckPrerequisiteCompletion(hash string, completedSubjects []string, grades map[string]float64) (bool, []string, error) {
	data, err := c.GetPrerequisiteData(hash)
	if err != nil {
		return false, nil, fmt.Errorf("failed to get prerequisite data: %w", err)
	}

	var missingPrereqs []string

	for _, prereq := range data.Prerequisites {
		if !prereq.Required {
			continue // Skip optional prerequisites
		}

		switch prereq.Type {
		case "ALL":
			// All subjects in this group must be completed
			for _, subjectID := range prereq.Subjects {
				if !contains(completedSubjects, subjectID) {
					missingPrereqs = append(missingPrereqs, subjectID)
					continue
				}

				// Check minimum grade if specified
				if prereq.MinGrade > 0 {
					if grade, exists := grades[subjectID]; !exists || grade < prereq.MinGrade {
						missingPrereqs = append(missingPrereqs, fmt.Sprintf("%s (min grade: %.2f)", subjectID, prereq.MinGrade))
					}
				}
			}

		case "ANY":
			// At least one subject in this group must be completed
			found := false
			for _, subjectID := range prereq.Subjects {
				if contains(completedSubjects, subjectID) {
					// Check minimum grade if specified
					if prereq.MinGrade > 0 {
						if grade, exists := grades[subjectID]; exists && grade >= prereq.MinGrade {
							found = true
							break
						}
					} else {
						found = true
						break
					}
				}
			}
			if !found {
				missingPrereqs = append(missingPrereqs, fmt.Sprintf("One of: %v", prereq.Subjects))
			}
		}
	}

	return len(missingPrereqs) == 0, missingPrereqs, nil
}

// Helper function to check if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// GetPrerequisiteSubjects extracts all prerequisite subject IDs from IPFS data
func (c *Client) GetPrerequisiteSubjects(hash string) ([]string, error) {
	data, err := c.GetPrerequisiteData(hash)
	if err != nil {
		return nil, fmt.Errorf("failed to get prerequisite data: %w", err)
	}

	var subjects []string
	for _, prereq := range data.Prerequisites {
		subjects = append(subjects, prereq.Subjects...)
	}

	// Remove duplicates
	unique := make(map[string]bool)
	var result []string
	for _, subject := range subjects {
		if !unique[subject] {
			unique[subject] = true
			result = append(result, subject)
		}
	}

	return result, nil
}
