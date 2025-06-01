package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"academictoken/x/subject/ipfs"

	"github.com/jdkato/prose/v2"
	"github.com/sirupsen/logrus"
	"gopkg.in/neurosnap/sentences.v1/english"
)

// =============================================================================
// Data Structures for the enhanced architecture
// =============================================================================

// ExtractedSyllabusData represents the complete data extracted from a syllabus
type ExtractedSyllabusData struct {
	Title                     string          `json:"title"`
	WorkloadHours             uint64          `json:"workload_hours"`
	Credits                   uint64          `json:"credits"`
	Description               string          `json:"description"`
	Objectives                []string        `json:"objectives"`
	Methodologies             []string        `json:"methodologies"`
	EvaluationMethods         []string        `json:"evaluation_methods"`
	BasicBibliography         []string        `json:"basic_bibliography"`
	ComplementaryBibliography []string        `json:"complementary_bibliography"`
	Topics                    []string        `json:"topics"`
	Keywords                  []string        `json:"keywords"`
	SubjectType               string          `json:"subject_type"`
	KnowledgeArea             string          `json:"knowledge_area"`
	Prerequisites             []string        `json:"prerequisites"`
	QualityScore              float64         `json:"quality_score"`
	ExtractionConfidence      float64         `json:"extraction_confidence"`
	DetailedTopics            []string        `json:"detailed_topics"` // Enhanced topics from program
	ProgramContent            []WeeklyContent `json:"program_content"` // Weekly program structure
}

// WeeklyContent represents structured weekly program content
type WeeklyContent struct {
	Week    int      `json:"week"`
	Content string   `json:"content"`
	Topics  []string `json:"topics"`
}

// SubjectContent matches the protobuf structure
type SubjectContent struct {
	Index         string `json:"index"`
	SubjectId     string `json:"subject_id"`
	Institution   string `json:"institution"`
	CourseId      string `json:"course_id"`
	Title         string `json:"title"`
	Code          string `json:"code"`
	WorkloadHours uint64 `json:"workload_hours"`
	Credits       uint64 `json:"credits"`
	Description   string `json:"description"`
	ContentHash   string `json:"content_hash"`
	SubjectType   string `json:"subject_type"`
	KnowledgeArea string `json:"knowledge_area"`
	IpfsLink      string `json:"ipfs_link"`
	Creator       string `json:"creator"`
}

// PrerequisiteGroup matches the protobuf structure for DAG formation
type PrerequisiteGroup struct {
	Id                       string   `json:"id"`
	SubjectId                string   `json:"subject_id"`
	GroupType                string   `json:"group_type"` // "ALL", "ANY", "MINIMUM", "NONE"
	MinimumCredits           uint64   `json:"minimum_credits"`
	MinimumCompletedSubjects uint64   `json:"minimum_completed_subjects"`
	SubjectIds               []string `json:"subject_ids"`
	Logic                    string   `json:"logic"`    // "AND", "OR", "XOR", "THRESHOLD"
	Priority                 int      `json:"priority"` // Priority for DAG traversal
	Confidence               float64  `json:"confidence"`
}

// ProcessingResult represents the result of processing a syllabus
type ProcessingResult struct {
	ContentHash       string                `json:"content_hash"`
	IPFSCID           string                `json:"ipfs_cid"`
	SubjectContent    SubjectContent        `json:"subject_content"`
	ExtractedData     ExtractedSyllabusData `json:"extracted_data"`
	Prerequisites     []PrerequisiteGroup   `json:"prerequisites"`
	LocalJSONPath     string                `json:"local_json_path"`
	IPFSJSONPath      string                `json:"ipfs_json_path"`
	ProcessingMetrics ProcessingMetrics     `json:"processing_metrics"`
}

// ProcessingMetrics provides quality metrics for the extraction process
type ProcessingMetrics struct {
	ProcessingTimeMs     int64   `json:"processing_time_ms"`
	TextLengthChars      int     `json:"text_length_chars"`
	SectionsIdentified   int     `json:"sections_identified"`
	TopicsExtracted      int     `json:"topics_extracted"`
	KeywordsExtracted    int     `json:"keywords_extracted"`
	BibliographyEntries  int     `json:"bibliography_entries"`
	PrerequisiteGroups   int     `json:"prerequisite_groups"`
	QualityScore         float64 `json:"quality_score"`
	ExtractionConfidence float64 `json:"extraction_confidence"`
	WeeksExtracted       int     `json:"weeks_extracted"`
}

// ProcessingCache provides caching for processed syllabi
type ProcessingCache struct {
	mu    sync.RWMutex
	cache map[string]ProcessingResult
}

func NewProcessingCache() *ProcessingCache {
	return &ProcessingCache{
		cache: make(map[string]ProcessingResult),
	}
}

func (c *ProcessingCache) Get(key string) (ProcessingResult, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	result, found := c.cache[key]
	return result, found
}

func (c *ProcessingCache) Set(key string, result ProcessingResult) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = result
}

// IPFSConfig represents IPFS configuration
type IPFSConfig struct {
	Enabled   bool   `json:"enabled"`
	Endpoint  string `json:"endpoint"`
	LocalPath string `json:"local_path"`
	UseHTTP   bool   `json:"use_http"`
	Pin       bool   `json:"pin"`
}

// CLIConfig represents CLI configuration
type CLIConfig struct {
	ChainID      string `json:"chain_id"`
	From         string `json:"from"`
	Node         string `json:"node"`
	Student      string `json:"student"`
	ProcessOnly  bool   `json:"process_only"`
	CreateOnly   bool   `json:"create_only"`
	TokenizeOnly bool   `json:"tokenize_only"`
	Debug        bool   `json:"debug"`
	CacheEnabled bool   `json:"cache_enabled"`
}

// =============================================================================
// Enhanced SyllabusProcessor - Main content processor
// =============================================================================

type SyllabusProcessor struct {
	ipfsConnector *IPFSConnector
	cache         *ProcessingCache
	logger        *logrus.Logger
}

func NewSyllabusProcessor(ipfsConnector *IPFSConnector, enableCache bool) *SyllabusProcessor {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	var cache *ProcessingCache
	if enableCache {
		cache = NewProcessingCache()
	}

	return &SyllabusProcessor{
		ipfsConnector: ipfsConnector,
		cache:         cache,
		logger:        logger,
	}
}

func (p *SyllabusProcessor) SetLogLevel(level logrus.Level) {
	p.logger.SetLevel(level)
}

// ProcessSyllabus processes a syllabus file and extracts academic content
func (p *SyllabusProcessor) ProcessSyllabus(rawContent, institution, courseId, subjectCode, creator string) (ProcessingResult, error) {
	startTime := time.Now()

	log := p.logger.WithFields(logrus.Fields{
		"institution":  institution,
		"course_id":    courseId,
		"subject_code": subjectCode,
		"creator":      creator,
	})

	log.Info("Starting syllabus processing")

	// Generate content hash for caching
	contentHash := p.generateContentHash(rawContent)

	// Check cache first
	if p.cache != nil {
		if cached, found := p.cache.Get(contentHash); found {
			log.Info("Using cached processing result")
			return cached, nil
		}
	}

	// Extract comprehensive data from syllabus with enhanced extraction
	log.Debug("Extracting syllabus data")
	extractedData, err := p.extractSyllabusData(rawContent)
	if err != nil {
		log.WithError(err).Error("Failed to extract syllabus data")
		return ProcessingResult{}, fmt.Errorf("error extracting syllabus data: %w", err)
	}

	// Validate extracted data
	if err := p.validateExtractedData(extractedData); err != nil {
		log.WithError(err).Warn("Data validation failed, but continuing with available data")
	}

	log.WithFields(logrus.Fields{
		"title":          extractedData.Title,
		"credits":        extractedData.Credits,
		"hours":          extractedData.WorkloadHours,
		"knowledge_area": extractedData.KnowledgeArea,
		"subject_type":   extractedData.SubjectType,
		"quality_score":  extractedData.QualityScore,
	}).Info("Syllabus data extracted successfully")

	// Generate unique subject ID
	subjectId := fmt.Sprintf("%s-%s-%s", institution, courseId, subjectCode)

	// Create SubjectContent for blockchain storage
	subjectContent := SubjectContent{
		Index:         subjectId,
		SubjectId:     subjectId,
		Institution:   institution,
		CourseId:      courseId,
		Title:         extractedData.Title,
		Code:          subjectCode,
		WorkloadHours: extractedData.WorkloadHours,
		Credits:       extractedData.Credits,
		Description:   p.truncateText(extractedData.Description, 200),
		ContentHash:   contentHash,
		SubjectType:   extractedData.SubjectType,
		KnowledgeArea: extractedData.KnowledgeArea,
		Creator:       creator,
	}

	// Extract prerequisites with enhanced DAG logic
	prerequisites := p.extractPrerequisites(extractedData.Prerequisites, subjectId, institution, courseId)

	// Calculate processing metrics
	processingTime := time.Since(startTime).Milliseconds()
	metrics := ProcessingMetrics{
		ProcessingTimeMs:     processingTime,
		TextLengthChars:      len(rawContent),
		SectionsIdentified:   p.countSections(rawContent),
		TopicsExtracted:      len(extractedData.Topics),
		KeywordsExtracted:    len(extractedData.Keywords),
		BibliographyEntries:  len(extractedData.BasicBibliography) + len(extractedData.ComplementaryBibliography),
		PrerequisiteGroups:   len(prerequisites),
		QualityScore:         extractedData.QualityScore,
		ExtractionConfidence: extractedData.ExtractionConfidence,
		WeeksExtracted:       len(extractedData.ProgramContent), // FIXED HERE
	}

	// Save complete data to local JSON file
	completeDataFile := fmt.Sprintf("%s_%s_complete.json", institution, subjectCode)
	completeData := map[string]interface{}{
		"subject_content":    subjectContent,
		"extracted_data":     extractedData,
		"prerequisites":      prerequisites,
		"processing_metrics": metrics,
		"raw_content":        rawContent,
		"processing_date":    time.Now().Format(time.RFC3339),
		"processor_version":  "2.1.0",
	}

	completeDataJSON, _ := json.MarshalIndent(completeData, "", "  ")
	if err := ioutil.WriteFile(completeDataFile, completeDataJSON, 0644); err != nil {
		log.WithError(err).Error("Failed to write complete data file")
		return ProcessingResult{}, fmt.Errorf("error writing complete data file: %w", err)
	}

	log.WithField("file", completeDataFile).Info("Complete data saved locally")

	result := ProcessingResult{
		ContentHash:       contentHash,
		SubjectContent:    subjectContent,
		ExtractedData:     extractedData,
		Prerequisites:     prerequisites,
		LocalJSONPath:     completeDataFile,
		ProcessingMetrics: metrics,
	}

	// Upload to IPFS if enabled
	if p.ipfsConnector != nil && p.ipfsConnector.IsEnabled() {
		log.Debug("Uploading to IPFS")
		ipfsCID, ipfsLink, err := p.ipfsConnector.UploadToIPFS(completeDataJSON)
		if err != nil {
			log.WithError(err).Warn("Could not upload to IPFS")
		} else {
			log.WithFields(logrus.Fields{
				"cid":  ipfsCID,
				"link": ipfsLink,
			}).Info("Data uploaded to IPFS successfully")

			// Update subject content with IPFS link
			subjectContent.IpfsLink = ipfsLink
			result.IPFSCID = ipfsCID
			result.SubjectContent = subjectContent

			// Store IPFS mapping
			if err := p.ipfsConnector.StoreIPFSMapping(institution, subjectCode, ipfsCID); err != nil {
				log.WithError(err).Warn("Could not store IPFS mapping")
			}
		}
	}

	// Cache result if caching is enabled
	if p.cache != nil {
		p.cache.Set(contentHash, result)
		log.Debug("Result cached for future requests")
	}

	log.WithField("processing_time_ms", processingTime).Info("Syllabus processing completed successfully")
	return result, nil
}

// validateExtractedData validates the quality of extracted data
func (p *SyllabusProcessor) validateExtractedData(data ExtractedSyllabusData) error {
	var errors []string

	if data.Title == "" || data.Title == "Disciplina Não Identificada" {
		errors = append(errors, "failed to extract subject title")
	}

	// Warn if title seems to contain formatting issues
	if strings.Contains(data.Title, "Código:") || len(strings.TrimSpace(data.Title)) != len(data.Title) {
		p.logger.Warn("Title may contain formatting issues")
	}

	if data.WorkloadHours == 0 {
		errors = append(errors, "failed to extract workload hours")
	}

	if data.Credits == 0 {
		errors = append(errors, "failed to extract credits")
	}

	if data.Description == "" {
		errors = append(errors, "failed to extract description/ementa")
	}

	if data.KnowledgeArea == "GENERAL" {
		errors = append(errors, "could not determine specific knowledge area")
	}

	// Check if topic extraction was too limited
	if len(data.Topics) < 5 {
		p.logger.Warn("Few topics extracted, extraction may need improvement")
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation issues: %s", strings.Join(errors, "; "))
	}

	return nil
}

// extractSyllabusData extracts all relevant data from the syllabus text with enhanced logic
func (p *SyllabusProcessor) extractSyllabusData(rawContent string) (ExtractedSyllabusData, error) {
	data := ExtractedSyllabusData{}

	// Identify sections of the syllabus
	sections := p.identifySections(rawContent)
	p.logger.WithField("sections_found", len(sections)).Debug("Sections identified")

	// Extract title with improved cleaning
	data.Title = p.extractAndCleanTitle(rawContent)

	// Extract numeric data with validation
	data.WorkloadHours = p.extractWorkloadHours(rawContent)
	data.Credits = p.extractCredits(rawContent)

	// Extract description/ementa with multiple patterns
	data.Description = p.extractDescription(sections, rawContent)

	// Extract structured content
	data.Objectives = p.extractObjectives(sections)
	data.Methodologies = p.extractMethodologies(sections)
	data.EvaluationMethods = p.extractEvaluationMethods(sections)

	// Extract bibliography with enhanced parsing
	data.BasicBibliography, data.ComplementaryBibliography = p.extractBibliography(sections)

	// Enhanced program content extraction
	data.ProgramContent = p.extractProgramContent(rawContent)
	p.logger.WithField("program_weeks_extracted", len(data.ProgramContent)).Debug("Program content extracted")

	// Extract topics using enhanced logic that combines basic and detailed extraction
	topics, detailedTopics, err := p.extractTopicsEnhanced(rawContent, data.ProgramContent)
	if err != nil {
		p.logger.WithError(err).Warn("Error extracting topics")
	} else {
		data.Topics = topics
		data.DetailedTopics = detailedTopics
	}

	// Extract keywords with better filtering
	keywords, err := p.extractKeywords(rawContent)
	if err != nil {
		p.logger.WithError(err).Warn("Error extracting keywords")
	} else {
		data.Keywords = keywords
	}

	// Determine subject type and knowledge area using enhanced algorithms
	data.SubjectType = p.determineSubjectType(rawContent, data.Keywords, sections)
	data.KnowledgeArea = p.determineKnowledgeArea(rawContent, data.Keywords, data.DetailedTopics)

	// Extract prerequisites with enhanced detection
	data.Prerequisites = p.extractPrerequisiteTexts(rawContent)

	// Calculate quality metrics
	data.QualityScore = p.calculateQualityScore(data)
	data.ExtractionConfidence = p.calculateExtractionConfidence(data, len(sections))

	return data, nil
}

// extractAndCleanTitle - Enhanced title extraction with proper cleaning
func (p *SyllabusProcessor) extractAndCleanTitle(text string) string {
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)(?:DISCIPLINA|MATÉRIA|SUBJECT):\s*(.+)`),
		regexp.MustCompile(`(?i)(?:NOME|NAME):\s*(.+)`),
		regexp.MustCompile(`(?i)(?:TÍTULO|TITLE):\s*(.+)`),
		regexp.MustCompile(`(?m)^Disciplina:\s*(.+)$`),
	}

	for _, pattern := range patterns {
		match := pattern.FindStringSubmatch(text)
		if len(match) > 1 {
			title := strings.TrimSpace(match[1])
			// Clean the title properly
			cleanedTitle := p.cleanTitle(title)
			if len(cleanedTitle) > 3 && !strings.Contains(strings.ToLower(cleanedTitle), "universidade") {
				return cleanedTitle
			}
		}
	}

	// Fallback: look for title-like lines
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) > 10 && len(line) < 100 &&
			!strings.Contains(strings.ToLower(line), "universidade") &&
			!strings.Contains(line, ":") &&
			regexp.MustCompile(`^[A-ZÁÉÍÓÚÇ]`).MatchString(line) {
			return p.cleanTitle(line)
		}
	}

	return "Disciplina Não Identificada"
}

// cleanTitle - Properly clean title removing codes and excessive whitespace
func (p *SyllabusProcessor) cleanTitle(title string) string {
	// Remove code pattern from title if present
	codeRegex := regexp.MustCompile(`\s+Código:\s*[A-Z0-9\s]+$`)
	title = codeRegex.ReplaceAllString(title, "")

	// Remove patterns like "MAT 001" at the end
	endCodeRegex := regexp.MustCompile(`\s+[A-Z]{2,4}\s*\d{3,4}\s*$`)
	title = endCodeRegex.ReplaceAllString(title, "")

	// Clean excessive whitespace while preserving single spaces
	title = regexp.MustCompile(`\s+`).ReplaceAllString(title, " ")

	return strings.TrimSpace(title)
}

// extractProgramContent - Extract weekly program structure (FIXED VERSION)
func (p *SyllabusProcessor) extractProgramContent(text string) []WeeklyContent {
	var program []WeeklyContent

	p.logger.Debug("Starting program content extraction")

	// Look for program section - more specific regex for this format
	programRegex := regexp.MustCompile(`(?s)Programa:\s*\n.*?Assunto:\s*\n(.*?)(?:\n\s*\nCritérios de Avaliação|\nCritérios de Avaliação|Bibliografia:|$)`)
	matches := programRegex.FindStringSubmatch(text)

	if len(matches) < 2 {
		// Try alternative pattern
		programRegex2 := regexp.MustCompile(`(?s)Programa:\s*\n(.*?)(?:\nCritérios de Avaliação|\nBibliografia:|$)`)
		matches = programRegex2.FindStringSubmatch(text)
		if len(matches) < 2 {
			p.logger.Debug("Program section not found with any pattern")
			return program
		}
	}

	programText := matches[1]
	p.logger.WithField("program_text_length", len(programText)).Debug("Program section found")

	// Split into lines and process each line
	lines := strings.Split(programText, "\n")

	for i, line := range lines {
		originalLine := line
		line = strings.TrimSpace(line)

		// Skip empty lines and header lines
		if line == "" || strings.Contains(line, "Semana:") || strings.Contains(line, "Assunto:") {
			continue
		}

		p.logger.WithFields(logrus.Fields{
			"line_number":  i,
			"line_content": line,
		}).Debug("Processing line")

		// Match the specific format: "   1      Content..." or "   10     Content..."
		weekRegex := regexp.MustCompile(`^\s*(\d+)\s+(.+)$`)
		if matches := weekRegex.FindStringSubmatch(originalLine); len(matches) > 2 {
			if week, err := strconv.Atoi(matches[1]); err == nil {
				content := strings.TrimSpace(matches[2])

				// Handle multi-line content (week 10 seems to span multiple lines)
				if i+1 < len(lines) {
					nextLine := strings.TrimSpace(lines[i+1])
					// If next line doesn't start with a number, it's continuation
					if nextLine != "" && !regexp.MustCompile(`^\s*\d+\s+`).MatchString(lines[i+1]) {
						content += " " + nextLine
						// Skip the next line since we've consumed it
						i++
					}
				}

				// Extract topics from this week's content
				topics := p.extractTopicsFromWeekContent(content)

				weekContent := WeeklyContent{
					Week:    week,
					Content: content,
					Topics:  topics,
				}

				program = append(program, weekContent)

				p.logger.WithFields(logrus.Fields{
					"week":         week,
					"content":      content,
					"topics_count": len(topics),
				}).Debug("Week parsed successfully")
			}
		}
	}

	p.logger.WithField("total_weeks", len(program)).Info("Program content extraction completed")
	return program
}

// extractTopicsFromWeekContent - Extract topics from individual week content (IMPROVED)
func (p *SyllabusProcessor) extractTopicsFromWeekContent(content string) []string {
	var topics []string

	// Split by common separators and clean
	separators := regexp.MustCompile(`[,;:]\s*`)
	parts := separators.Split(content, -1)

	for _, part := range parts {
		part = strings.TrimSpace(part)
		// Filter out very short parts and common words
		if len(part) > 4 && len(part) < 100 && !p.isCommonStopWord(part) {
			topics = append(topics, part)
		}
	}

	// If no good splits, use the whole content as a topic
	if len(topics) == 0 && len(content) > 3 {
		topics = append(topics, content)
	}

	return topics
}

// extractTopicsEnhanced - Enhanced topic extraction combining basic and detailed approaches
func (p *SyllabusProcessor) extractTopicsEnhanced(rawContent string, programContent []WeeklyContent) ([]string, []string, error) {
	var basicTopics []string
	var detailedTopics []string
	topicSet := make(map[string]bool)

	// Extract from program content first (most detailed)
	for _, week := range programContent {
		for _, topic := range week.Topics {
			if !topicSet[topic] && len(topic) > 3 {
				detailedTopics = append(detailedTopics, topic)
				topicSet[topic] = true
			}
		}

		// Also add the full week content as a topic if it's reasonable length
		if len(week.Content) > 5 && len(week.Content) < 200 && !topicSet[week.Content] {
			detailedTopics = append(detailedTopics, week.Content)
			topicSet[week.Content] = true
		}
	}

	// Extract using NLP for additional topics
	tokenizer, err := english.NewSentenceTokenizer(nil)
	if err != nil {
		return basicTopics, detailedTopics, err
	}

	sentences := tokenizer.Tokenize(rawContent)

	// Enhanced topic patterns
	topicPatterns := []*regexp.Regexp{
		regexp.MustCompile(`(?m)^(?:\d+\.?\s*)(.+)$`),         // Numbered topics
		regexp.MustCompile(`(?m)^(?:[\-•\*]\s*)(.+)$`),        // Bulleted topics
		regexp.MustCompile(`(?m)^(?:Semana\s+\d+:\s*)(.+)$`),  // Weekly topics
		regexp.MustCompile(`(?m)^(?:Unidade\s+\d+:\s*)(.+)$`), // Unit topics
	}

	// Extract from structured content
	for _, s := range sentences {
		for _, pattern := range topicPatterns {
			matches := pattern.FindAllStringSubmatch(s.Text, -1)
			for _, match := range matches {
				if len(match) > 1 {
					topic := strings.TrimSpace(match[1])
					if len(topic) > 5 && len(topic) < 200 && !topicSet[topic] {
						detailedTopics = append(detailedTopics, topic)
						topicSet[topic] = true
					}
				}
			}
		}
	}

	// If we still have few detailed topics, extract noun phrases using NLP
	if len(detailedTopics) < 10 {
		doc, err := prose.NewDocument(rawContent)
		if err == nil {
			for _, token := range doc.Tokens() {
				if strings.HasPrefix(token.Tag, "NN") && len(token.Text) > 3 {
					topic := token.Text
					if !topicSet[strings.ToLower(topic)] && !p.isCommonStopWord(topic) {
						detailedTopics = append(detailedTopics, topic)
						topicSet[strings.ToLower(topic)] = true
					}
				}
			}
		}
	}

	// Create basic topics (shorter list for basic display)
	basicTopicCount := 20
	if len(detailedTopics) < basicTopicCount {
		basicTopicCount = len(detailedTopics)
	}

	basicTopics = make([]string, basicTopicCount)
	copy(basicTopics, detailedTopics[:basicTopicCount])

	// Filter and clean topics
	detailedTopics = p.filterTopics(detailedTopics)
	basicTopics = p.filterTopics(basicTopics)

	return basicTopics, detailedTopics, nil
}

// Enhanced section identification with better patterns
func (p *SyllabusProcessor) identifySections(rawContent string) map[string]string {
	sections := make(map[string]string)
	lines := strings.Split(rawContent, "\n")

	var currentSection string
	var currentContent []string

	// Enhanced regex patterns for section headers
	headerPatterns := []*regexp.Regexp{
		regexp.MustCompile(`^([A-ZÁÉÍÓÚÇ\s]+)[:]\s*$`),
		regexp.MustCompile(`^(Ementa|Bibliografia|Programa|Critérios de Avaliação|Objetivos|Metodologia|Avaliação)[:]\s*$`),
		regexp.MustCompile(`^(EMENTA|BIBLIOGRAFIA|PROGRAMA|OBJETIVOS|METODOLOGIA|AVALIAÇÃO)[:]\s*$`),
		regexp.MustCompile(`^(\d+\.?\s*[A-ZÁÉÍÓÚÇ\s]+)[:]\s*$`),
	}

	for i, line := range lines {
		line = strings.TrimSpace(line)

		// Handle special sections
		if strings.HasPrefix(line, "Ementa:") {
			p.saveCurrentSection(&sections, currentSection, currentContent)
			currentSection = "EMENTA"
			currentContent = []string{}
			continue
		}

		// Check all header patterns
		isHeader := false
		for _, pattern := range headerPatterns {
			matches := pattern.FindStringSubmatch(line)
			if len(matches) > 1 {
				p.saveCurrentSection(&sections, currentSection, currentContent)
				currentSection = strings.ToUpper(strings.TrimSpace(matches[1]))
				currentContent = []string{}
				isHeader = true
				break
			}
		}

		if !isHeader && currentSection != "" && line != "" && !p.isHeaderLine(line) {
			currentContent = append(currentContent, line)
		}

		// Save last section
		if i == len(lines)-1 {
			p.saveCurrentSection(&sections, currentSection, currentContent)
		}
	}

	return sections
}

func (p *SyllabusProcessor) saveCurrentSection(sections *map[string]string, sectionName string, content []string) {
	if sectionName != "" && len(content) > 0 {
		(*sections)[sectionName] = strings.Join(content, "\n")
	}
}

func (p *SyllabusProcessor) isHeaderLine(line string) bool {
	headerIndicators := []string{"Código:", "Disciplina:", "Departamento:", "Unidade:", "Carga Horária:"}
	for _, indicator := range headerIndicators {
		if strings.HasPrefix(line, indicator) {
			return true
		}
	}
	return false
}

// Enhanced workload hours extraction
func (p *SyllabusProcessor) extractWorkloadHours(text string) uint64 {
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)(?:CARGA HORÁRIA TOTAL|Carga Horária Total):\s*(\d+)`),
		regexp.MustCompile(`(?i)(?:CARGA HORÁRIA|CH|WORKLOAD|HOURS?):\s*(\d+)`),
		regexp.MustCompile(`(?i)Total:\s*(\d+)`),
		regexp.MustCompile(`(?i)(\d+)\s*(?:horas?|hrs?|h)\b`),
		regexp.MustCompile(`(?m)^Carga Horária Total:\s*(\d+)$`),
	}

	for _, pattern := range patterns {
		match := pattern.FindStringSubmatch(text)
		if len(match) > 1 {
			var hours uint64
			if n, err := fmt.Sscanf(match[1], "%d", &hours); n == 1 && err == nil && hours > 0 && hours <= 500 {
				return hours
			}
		}
	}

	return 60 // Default hours if not found
}

// Enhanced credits extraction
func (p *SyllabusProcessor) extractCredits(text string) uint64 {
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)(?:CRÉDITOS?|CR|CREDITS?):\s*(\d+)`),
		regexp.MustCompile(`(?i)(\d+)\s*(?:créditos?|credits?)\b`),
		regexp.MustCompile(`(?m)^[Nn]o de créditos:\s*(\d+)$`),
		regexp.MustCompile(`(?m)^Créditos:\s*(\d+)$`),
	}

	for _, pattern := range patterns {
		match := pattern.FindStringSubmatch(text)
		if len(match) > 1 {
			var credits uint64
			if n, err := fmt.Sscanf(match[1], "%d", &credits); n == 1 && err == nil && credits > 0 && credits <= 20 {
				return credits
			}
		}
	}

	return 4 // Default credits if not found
}

// Enhanced description extraction
func (p *SyllabusProcessor) extractDescription(sections map[string]string, rawContent string) string {
	// Priority order for description sources
	sectionKeys := []string{"EMENTA", "SYLLABUS", "DESCRIÇÃO", "DESCRIPTION"}

	for _, key := range sectionKeys {
		if section, ok := sections[key]; ok && strings.TrimSpace(section) != "" {
			return strings.TrimSpace(section)
		}
	}

	// Fallback: extract from raw content
	ementaRegex := regexp.MustCompile(`(?i)Ementa:\s*\n((?:[^\n]+\n?)+?)(?:\n\n|\nPrograma:|\nObjetivos:|\nCritérios|$)`)
	match := ementaRegex.FindStringSubmatch(rawContent)
	if len(match) > 1 {
		return strings.TrimSpace(match[1])
	}

	return ""
}

// Extract objectives from sections
func (p *SyllabusProcessor) extractObjectives(sections map[string]string) []string {
	keys := []string{"OBJETIVOS", "OBJECTIVES", "OBJETIVO GERAL", "OBJETIVOS ESPECÍFICOS"}

	for _, key := range keys {
		if section, ok := sections[key]; ok {
			return p.extractListItems(section)
		}
	}

	return nil
}

// Extract methodologies from sections
func (p *SyllabusProcessor) extractMethodologies(sections map[string]string) []string {
	keys := []string{"METODOLOGIA", "METHODOLOGY", "METODOLOGIAS", "MÉTODO DE ENSINO"}

	for _, key := range keys {
		if section, ok := sections[key]; ok {
			return p.extractListItems(section)
		}
	}

	return nil
}

// Extract evaluation methods from sections
func (p *SyllabusProcessor) extractEvaluationMethods(sections map[string]string) []string {
	keys := []string{"AVALIAÇÃO", "EVALUATION", "CRITÉRIOS DE AVALIAÇÃO", "MÉTODOS DE AVALIAÇÃO"}

	for _, key := range keys {
		if section, ok := sections[key]; ok {
			return p.extractListItems(section)
		}
	}

	return nil
}

// Enhanced bibliography extraction
func (p *SyllabusProcessor) extractBibliography(sections map[string]string) ([]string, []string) {
	var basic, complementary []string

	// Enhanced patterns for bibliography detection
	basicPatterns := []string{
		"BIBLIOGRAFIA BÁSICA", "BASIC BIBLIOGRAPHY", "REFERÊNCIAS BÁSICAS",
		"BIBLIOGRAFIA", "BIBLIOGRAPHY", "REFERÊNCIAS",
	}

	complementaryPatterns := []string{
		"BIBLIOGRAFIA COMPLEMENTAR", "COMPLEMENTARY BIBLIOGRAPHY",
		"REFERÊNCIAS COMPLEMENTARES", "BIBLIOGRAFIA ADICIONAL",
	}

	// Check for basic bibliography
	for _, pattern := range basicPatterns {
		if content, ok := sections[pattern]; ok {
			basic = p.parseBibliographyEntries(content)
			break
		}
	}

	// Check for complementary bibliography
	for _, pattern := range complementaryPatterns {
		if content, ok := sections[pattern]; ok {
			complementary = p.parseBibliographyEntries(content)
			break
		}
	}

	return basic, complementary
}

func (p *SyllabusProcessor) parseBibliographyEntries(content string) []string {
	var entries []string

	// Split by common bibliography patterns
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(?m)^([A-ZÁÉÍÓÚÇ][A-ZÁÉÍÓÚÇ\s,\.]+)\s*[-–]`), // Author names
		regexp.MustCompile(`(?m)^\d+\.\s*(.+)$`),                         // Numbered entries
		regexp.MustCompile(`(?m)^[•\-\*]\s*(.+)$`),                       // Bulleted entries
	}

	for _, pattern := range patterns {
		matches := pattern.FindAllStringSubmatch(content, -1)
		for _, match := range matches {
			if len(match) > 1 {
				entry := strings.TrimSpace(match[1])
				if len(entry) > 10 { // Filter out very short entries
					entries = append(entries, entry)
				}
			}
		}
		if len(entries) > 0 {
			break // Use first successful pattern
		}
	}

	// Fallback: split by lines if no pattern worked
	if len(entries) == 0 {
		lines := strings.Split(content, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if len(line) > 10 && !p.isCommonStopWord(line) {
				entries = append(entries, line)
			}
		}
	}

	return p.removeDuplicates(entries)
}

// Enhanced list item extraction
func (p *SyllabusProcessor) extractListItems(text string) []string {
	var items []string

	// Multiple regex patterns for different list formats
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(?m)^(?:\d+\.[\d\.]*\s+)(.+)$`), // Numbered lists
		regexp.MustCompile(`(?m)^(?:[\-•\*]\s+)(.+)$`),      // Bullet lists
		regexp.MustCompile(`(?m)^(?:[a-z]\)\s+)(.+)$`),      // Lettered lists
		regexp.MustCompile(`(?m)^(?:[IVX]+\.\s+)(.+)$`),     // Roman numeral lists
	}

	for _, pattern := range patterns {
		matches := pattern.FindAllStringSubmatch(text, -1)
		for _, match := range matches {
			if len(match) > 1 {
				item := strings.TrimSpace(match[1])
				if len(item) > 3 {
					items = append(items, item)
				}
			}
		}
		if len(items) > 0 {
			return items // Return first successful pattern
		}
	}

	// Fallback: split by lines and filter
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) > 5 && !p.isCommonStopWord(line) {
			items = append(items, line)
		}
	}

	return items
}

// Enhanced keywords extraction
func (p *SyllabusProcessor) extractKeywords(rawContent string) ([]string, error) {
	doc, err := prose.NewDocument(rawContent)
	if err != nil {
		return nil, err
	}

	wordFreq := make(map[string]int)

	// Extract entities with frequency counting
	for _, ent := range doc.Entities() {
		word := strings.ToLower(ent.Text)
		if len(word) > 2 && !p.isCommonStopWord(word) {
			wordFreq[word]++
		}
	}

	// Extract important nouns
	for _, tok := range doc.Tokens() {
		if strings.HasPrefix(tok.Tag, "NN") && len(tok.Text) > 2 {
			word := strings.ToLower(tok.Text)
			if !p.isCommonStopWord(word) {
				wordFreq[word]++
			}
		}
	}

	// Sort by frequency and select top keywords
	type wordFreqPair struct {
		word string
		freq int
	}

	var pairs []wordFreqPair
	for word, freq := range wordFreq {
		if freq >= 2 { // Only consider words that appear at least twice
			pairs = append(pairs, wordFreqPair{word, freq})
		}
	}

	// Sort by frequency (descending)
	for i := 0; i < len(pairs)-1; i++ {
		for j := i + 1; j < len(pairs); j++ {
			if pairs[i].freq < pairs[j].freq {
				pairs[i], pairs[j] = pairs[j], pairs[i]
			}
		}
	}

	// Extract top 15 keywords
	var keywords []string
	maxKeywords := 15
	if len(pairs) < maxKeywords {
		maxKeywords = len(pairs)
	}

	for i := 0; i < maxKeywords; i++ {
		keywords = append(keywords, pairs[i].word)
	}

	return keywords, nil
}

// Enhanced subject type determination
func (p *SyllabusProcessor) determineSubjectType(content string, keywords []string, sections map[string]string) string {
	contentLower := strings.ToLower(content)

	// Scoring system for different types
	scores := map[string]int{
		"PRACTICAL":   0,
		"THEORETICAL": 0,
		"PROJECT":     0,
		"LAB":         0,
	}

	// Check keywords first (higher weight)
	for _, keyword := range keywords {
		keywordLower := strings.ToLower(keyword)

		if p.containsAny(keywordLower, []string{"laboratório", "prática", "experimental", "lab", "laboratory", "practical"}) {
			scores["PRACTICAL"] += 3
			scores["LAB"] += 2
		}

		if p.containsAny(keywordLower, []string{"projeto", "desenvolvimento", "project", "development"}) {
			scores["PROJECT"] += 3
		}

		if p.containsAny(keywordLower, []string{"teoria", "conceitos", "fundamentos", "theory", "concepts", "fundamentals"}) {
			scores["THEORETICAL"] += 2
		}
	}

	// Check content patterns
	practicalTerms := []string{"laboratório", "prática", "experimental", "hands-on", "aplicação prática"}
	theoreticalTerms := []string{"teoria", "conceitos", "fundamentos", "princípios", "análise teórica"}
	projectTerms := []string{"projeto", "desenvolvimento", "implementação", "criação"}

	for _, term := range practicalTerms {
		if strings.Contains(contentLower, term) {
			scores["PRACTICAL"] += 2
		}
	}

	for _, term := range theoreticalTerms {
		if strings.Contains(contentLower, term) {
			scores["THEORETICAL"] += 1
		}
	}

	for _, term := range projectTerms {
		if strings.Contains(contentLower, term) {
			scores["PROJECT"] += 2
		}
	}

	// Check section content
	if metodologia, ok := sections["METODOLOGIA"]; ok {
		metodologiaLower := strings.ToLower(metodologia)
		if p.containsAny(metodologiaLower, practicalTerms) {
			scores["PRACTICAL"] += 2
		}
	}

	// Find highest scoring type
	maxScore := 0
	bestType := "THEORETICAL" // Default

	for subjectType, score := range scores {
		if score > maxScore {
			maxScore = score
			bestType = subjectType
		}
	}

	return bestType
}

// Enhanced knowledge area determination
func (p *SyllabusProcessor) determineKnowledgeArea(content string, keywords []string, detailedTopics []string) string {
	contentLower := strings.ToLower(content)

	// Enhanced knowledge areas with more comprehensive keywords
	areas := map[string][]string{
		"MATHEMATICS": {
			"matemática", "cálculo", "álgebra", "geometria", "estatística", "probabilidade",
			"mathematics", "calculus", "algebra", "geometry", "statistics", "probability",
			"derivadas", "integrais", "funções", "trigonométricas", "números", "equações",
			"derivatives", "integrals", "functions", "trigonometric", "numbers", "equations",
		},
		"COMPUTER_SCIENCE": {
			"computação", "programação", "algoritmo", "software", "dados", "sistemas",
			"computer", "programming", "algorithm", "data", "systems", "informática",
			"java", "python", "c++", "javascript", "database", "web", "artificial intelligence",
		},
		"PHYSICS": {
			"física", "mecânica", "termodinâmica", "eletromagnetismo", "óptica", "quântica",
			"physics", "mechanics", "thermodynamics", "electromagnetism", "optics", "quantum",
			"energia", "energia", "força", "movimento", "ondas", "particles",
		},
		"CHEMISTRY": {
			"química", "orgânica", "inorgânica", "físico-química", "analítica",
			"chemistry", "organic", "inorganic", "physical chemistry", "analytical",
			"molecular", "elementos", "elements", "reações", "reactions", "compostos",
		},
		"BIOLOGY": {
			"biologia", "genética", "ecologia", "anatomia", "fisiologia", "microbiologia",
			"biology", "genetics", "ecology", "anatomy", "physiology", "microbiology",
			"célula", "cell", "vida", "life", "evolução", "evolution", "biodiversidade",
		},
		"ENGINEERING": {
			"engenharia", "mecânica", "elétrica", "civil", "industrial", "química",
			"engineering", "mechanical", "electrical", "civil", "industrial", "chemical",
			"estruturas", "materiais", "projeto", "design", "construção", "manufacturing",
		},
		"BUSINESS": {
			"administração", "economia", "gestão", "marketing", "finanças", "contabilidade",
			"business", "management", "economics", "finance", "accounting", "entrepreneurship",
			"estratégia", "strategy", "recursos humanos", "human resources", "vendas", "sales",
		},
		"LAW": {
			"direito", "jurídica", "legal", "constituição", "código", "tribunal", "justiça",
			"law", "legal", "constitution", "code", "court", "justice", "jurisprudence",
			"civil", "penal", "administrativo", "criminal", "constitutional", "commercial",
		},
		"MEDICINE": {
			"medicina", "saúde", "clínica", "diagnóstico", "tratamento", "farmacologia",
			"medical", "health", "clinical", "diagnosis", "treatment", "pharmacology",
			"patologia", "pathology", "terapia", "therapy", "cirurgia", "surgery",
		},
		"HUMANITIES": {
			"filosofia", "história", "literatura", "sociologia", "antropologia", "psicologia",
			"philosophy", "history", "literature", "sociology", "anthropology", "psychology",
			"cultura", "culture", "sociedade", "society", "língua", "language", "arte", "art",
		},
	}

	// Scoring system with weights
	areaScores := make(map[string]float64)

	// Check keywords (highest weight)
	for _, keyword := range keywords {
		keywordLower := strings.ToLower(keyword)
		for area, areaKeywords := range areas {
			for _, areaKeyword := range areaKeywords {
				if strings.Contains(keywordLower, areaKeyword) || keywordLower == areaKeyword {
					areaScores[area] += 3.0
				}
			}
		}
	}

	// Check detailed topics (medium weight)
	for _, topic := range detailedTopics {
		topicLower := strings.ToLower(topic)
		for area, areaKeywords := range areas {
			for _, areaKeyword := range areaKeywords {
				if strings.Contains(topicLower, areaKeyword) {
					areaScores[area] += 2.0
				}
			}
		}
	}

	// Check content (lower weight)
	for area, areaKeywords := range areas {
		for _, keyword := range areaKeywords {
			if strings.Contains(contentLower, keyword) {
				areaScores[area] += 1.0
			}
		}
	}

	// Find area with highest score
	maxScore := 0.0
	bestArea := "GENERAL"

	for area, score := range areaScores {
		if score > maxScore {
			maxScore = score
			bestArea = area
		}
	}

	// Require minimum score for classification
	if maxScore < 2.0 {
		return "GENERAL"
	}

	return bestArea
}

// Enhanced prerequisite text extraction with better patterns
func (p *SyllabusProcessor) extractPrerequisiteTexts(content string) []string {
	var prerequisites []string

	// Enhanced patterns for prerequisite extraction
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)(?:PRÉ-REQUISITOS?|PREREQUISITES?):\s*\n((?:[^\n]+\n)*?)(?:\n\n|\n(?:Ementa|Bibliografia|Disciplina:|$))`),
		regexp.MustCompile(`(?i)(?:REQUER|REQUIRES?|NECESSITA|NEEDS?):\s*(.+?)(?:\n\n|\n(?:[A-Z]|$))`),
		regexp.MustCompile(`(?i)(?:DEPENDÊNCIAS?|DEPENDENCIES?|DEPENDE\s+DE):\s*(.+?)(?:\n\n|\n(?:[A-Z]|$))`),
		regexp.MustCompile(`(?i)(?:REQUISITOS?):\s*(.+?)(?:\n\n|\n(?:[A-Z]|$))`),
	}

	// Check for explicit "no prerequisites" patterns first
	noPrereqPatterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)Pré-requisitos:\s*\n\s*Código:\s*\n\s*Disciplina:\s*\n`),
		regexp.MustCompile(`(?i)(?:PRÉ-REQUISITOS?|PREREQUISITES?):\s*(?:Não|None|Nenhum|N/A|\-|Não há)\s*`),
		regexp.MustCompile(`(?i)(?:PRÉ-REQUISITOS?|PREREQUISITES?):\s*$`),
	}

	for _, pattern := range noPrereqPatterns {
		if pattern.MatchString(content) {
			p.logger.Debug("No prerequisites detected")
			return []string{} // Explicitly no prerequisites
		}
	}

	// Extract prerequisites using patterns
	prereqSet := make(map[string]bool) // To avoid duplicates

	for _, pattern := range patterns {
		matches := pattern.FindAllStringSubmatch(content, -1)
		for _, match := range matches {
			if len(match) > 1 {
				prereqSection := strings.TrimSpace(match[1])

				// Process each line in the prerequisite section
				lines := strings.Split(prereqSection, "\n")
				for _, line := range lines {
					line = strings.TrimSpace(line)

					// Skip empty lines, headers, and common non-prerequisite text
					if p.isValidPrerequisite(line) && !prereqSet[line] {
						prerequisites = append(prerequisites, line)
						prereqSet[line] = true
					}
				}
			}
		}
	}

	p.logger.WithField("prerequisites_found", len(prerequisites)).Debug("Prerequisites extracted")
	return prerequisites
}

// Enhanced prerequisite validation
func (p *SyllabusProcessor) isValidPrerequisite(line string) bool {
	if line == "" || len(line) < 3 {
		return false
	}

	// Skip common non-prerequisite patterns
	skipPatterns := []string{
		"Código:", "Disciplina:", "código:", "disciplina:",
		"Não há", "None", "Nenhum", "N/A", "-",
	}

	for _, pattern := range skipPatterns {
		if strings.Contains(line, pattern) {
			return false
		}
	}

	// Skip lines that are just punctuation or very short
	if len(strings.TrimSpace(strings.Trim(line, ".,;:-"))) < 3 {
		return false
	}

	return !p.isCommonStopWord(line)
}

// Enhanced prerequisites conversion to DAG groups
func (p *SyllabusProcessor) extractPrerequisites(prereqTexts []string, subjectId, institution, courseId string) []PrerequisiteGroup {
	var groups []PrerequisiteGroup

	if len(prereqTexts) == 0 {
		// Explicitly no prerequisites
		groups = append(groups, PrerequisiteGroup{
			Id:         fmt.Sprintf("prereq-%s-none", subjectId),
			SubjectId:  subjectId,
			GroupType:  "NONE",
			Logic:      "NONE",
			Priority:   0,
			SubjectIds: []string{},
			Confidence: 1.0,
		})
		return groups
	}

	// Process each prerequisite text with enhanced logic
	for i, prereqText := range prereqTexts {
		groupType := "ALL" // Default
		logic := "AND"     // Default logic for DAG
		priority := 1      // Default priority
		confidence := 0.8  // Default confidence

		// Detect different types of prerequisites
		prereqLower := strings.ToLower(prereqText)

		if p.containsAny(prereqLower, []string{"ou", "or", "either"}) {
			groupType = "ANY"
			logic = "OR"
			confidence = 0.9
		}

		if p.containsAny(prereqLower, []string{"pelo menos", "at least", "mínimo de", "minimum of"}) {
			groupType = "MINIMUM"
			logic = "THRESHOLD"
			confidence = 0.7
		}

		if p.containsAny(prereqLower, []string{"todos", "all", "ambos", "both"}) {
			groupType = "ALL"
			logic = "AND"
			confidence = 0.9
		}

		// Extract subject codes with enhanced patterns
		subjectIds := p.extractSubjectCodes(prereqText, institution, courseId)
		minRequired := p.extractMinimumCount(prereqText, uint64(len(subjectIds)))

		// Adjust confidence based on extraction success
		if len(subjectIds) == 0 {
			confidence *= 0.5 // Lower confidence if no subjects extracted
		}

		group := PrerequisiteGroup{
			Id:                       fmt.Sprintf("prereq-%s-%d", subjectId, i+1),
			SubjectId:                subjectId,
			GroupType:                groupType,
			MinimumCredits:           p.extractMinimumCredits(prereqText),
			MinimumCompletedSubjects: minRequired,
			SubjectIds:               subjectIds,
			Logic:                    logic,
			Priority:                 priority,
			Confidence:               confidence,
		}

		groups = append(groups, group)
	}

	return groups
}

// Enhanced subject code extraction with more patterns
func (p *SyllabusProcessor) extractSubjectCodes(prereqText, institution, courseId string) []string {
	var subjectIds []string
	codeSet := make(map[string]bool) // Avoid duplicates

	// Enhanced regex patterns for subject codes
	codePatterns := []*regexp.Regexp{
		regexp.MustCompile(`\b([A-Z]{2,4}\d{3,4})\b`),           // MAT1234, COMP123
		regexp.MustCompile(`\b([A-Z]{3}\s?\d{3})\b`),            // MAT 123, MAT123
		regexp.MustCompile(`\b(\d{3}[A-Z]{2,3})\b`),             // 123MAT
		regexp.MustCompile(`\b([A-Z]{2,4}-\d{3,4})\b`),          // MAT-1234
		regexp.MustCompile(`\b([A-Z]{2,4}\s*\d{3,4})\b`),        // MAT 1234
		regexp.MustCompile(`(?i)\b([A-Z]{2,4}\d{2,4}[A-Z]?)\b`), // More flexible pattern
	}

	for _, pattern := range codePatterns {
		matches := pattern.FindAllStringSubmatch(prereqText, -1)
		for _, match := range matches {
			if len(match) > 1 {
				code := strings.ToUpper(strings.TrimSpace(match[1]))
				// Clean code (remove spaces)
				code = strings.ReplaceAll(code, " ", "")

				if !codeSet[code] && len(code) >= 5 && len(code) <= 8 {
					fullSubjectId := fmt.Sprintf("%s-%s-%s", institution, courseId, code)
					subjectIds = append(subjectIds, fullSubjectId)
					codeSet[code] = true
				}
			}
		}
	}

	// If no codes found, try to extract meaningful names
	if len(subjectIds) == 0 {
		// Clean and format the prerequisite text as a subject identifier
		cleanText := regexp.MustCompile(`[^a-zA-Z0-9\s]`).ReplaceAllString(prereqText, "")
		cleanText = regexp.MustCompile(`\s+`).ReplaceAllString(strings.TrimSpace(cleanText), "_")

		if len(cleanText) > 3 && len(cleanText) < 50 {
			subjectIds = append(subjectIds, fmt.Sprintf("%s-%s-%s", institution, courseId, cleanText))
		}
	}

	return subjectIds
}

// Extract minimum credit requirements
func (p *SyllabusProcessor) extractMinimumCredits(prereqText string) uint64 {
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)(\d+)\s*créditos?`),
		regexp.MustCompile(`(?i)(\d+)\s*credits?`),
		regexp.MustCompile(`(?i)mínimo\s+de\s+(\d+)\s*créditos?`),
		regexp.MustCompile(`(?i)minimum\s+of\s+(\d+)\s*credits?`),
	}

	for _, pattern := range patterns {
		match := pattern.FindStringSubmatch(prereqText)
		if len(match) > 1 {
			var credits uint64
			if n, err := fmt.Sscanf(match[1], "%d", &credits); n == 1 && err == nil {
				return credits
			}
		}
	}

	return 0
}

// Enhanced minimum count extraction
func (p *SyllabusProcessor) extractMinimumCount(prereqText string, defaultCount uint64) uint64 {
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)pelo\s+menos\s+(\d+)`),
		regexp.MustCompile(`(?i)at\s+least\s+(\d+)`),
		regexp.MustCompile(`(?i)mínimo\s+de\s+(\d+)`),
		regexp.MustCompile(`(?i)minimum\s+of\s+(\d+)`),
		regexp.MustCompile(`(?i)(\d+)\s+ou\s+mais`),
		regexp.MustCompile(`(?i)(\d+)\s+or\s+more`),
	}

	for _, pattern := range patterns {
		match := pattern.FindStringSubmatch(prereqText)
		if len(match) > 1 {
			var count uint64
			if n, err := fmt.Sscanf(match[1], "%d", &count); n == 1 && err == nil {
				return count
			}
		}
	}

	return defaultCount
}

// Quality score calculation
func (p *SyllabusProcessor) calculateQualityScore(data ExtractedSyllabusData) float64 {
	score := 0.0
	maxScore := 10.0

	// Title quality (0-1 points)
	if data.Title != "" && data.Title != "Disciplina Não Identificada" {
		score += 1.0
	}

	// Description quality (0-1 points)
	if len(data.Description) > 20 {
		score += 1.0
	}

	// Numeric data (0-1 points)
	if data.WorkloadHours > 0 && data.Credits > 0 {
		score += 1.0
	}

	// Content richness (0-2 points)
	if len(data.DetailedTopics) > 10 {
		score += 1.0
	}
	if len(data.Keywords) > 3 {
		score += 1.0
	}

	// Bibliography (0-2 points)
	if len(data.BasicBibliography) > 0 {
		score += 1.0
	}
	if len(data.ComplementaryBibliography) > 0 {
		score += 1.0
	}

	// Structured content (0-3 points)
	if len(data.Objectives) > 0 {
		score += 1.0
	}
	if len(data.Methodologies) > 0 {
		score += 1.0
	}
	if len(data.EvaluationMethods) > 0 {
		score += 1.0
	}

	return (score / maxScore) * 10.0 // Scale to 0-10
}

// Extraction confidence calculation
func (p *SyllabusProcessor) calculateExtractionConfidence(data ExtractedSyllabusData, sectionsFound int) float64 {
	confidence := 0.0

	// Base confidence from sections identified
	if sectionsFound > 3 {
		confidence += 0.3
	} else if sectionsFound > 1 {
		confidence += 0.2
	} else {
		confidence += 0.1
	}

	// Confidence from successful extractions
	if data.Title != "Disciplina Não Identificada" {
		confidence += 0.2
	}
	if data.WorkloadHours > 0 && data.Credits > 0 {
		confidence += 0.2
	}
	if len(data.Description) > 10 {
		confidence += 0.1
	}
	if data.KnowledgeArea != "GENERAL" {
		confidence += 0.1
	}
	if len(data.DetailedTopics) > 5 {
		confidence += 0.1
	}

	return confidence
}

// Helper functions
func (p *SyllabusProcessor) generateContentHash(content string) string {
	hasher := sha256.New()
	hasher.Write([]byte(content))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (p *SyllabusProcessor) truncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen-3] + "..."
}

func (p *SyllabusProcessor) countSections(content string) int {
	return len(p.identifySections(content))
}

func (p *SyllabusProcessor) containsAny(text string, terms []string) bool {
	for _, term := range terms {
		if strings.Contains(text, term) {
			return true
		}
	}
	return false
}

func (p *SyllabusProcessor) isCommonStopWord(word string) bool {
	stopWords := []string{
		"não", "nenhum", "nenhuma", "none", "no", "sim", "yes", "ter", "have",
		"concluído", "concluded", "completed", "aprovado", "approved", "cursado", "taken",
		"the", "and", "or", "but", "in", "on", "at", "to", "for", "of", "with", "by",
		"que", "de", "da", "do", "das", "dos", "para", "por", "com", "em", "na", "no",
	}

	wordLower := strings.ToLower(word)
	for _, stopWord := range stopWords {
		if wordLower == stopWord {
			return true
		}
	}
	return false
}

// filterTopics - Enhanced filtering to remove administrative data
func (p *SyllabusProcessor) filterTopics(topics []string) []string {
	var filtered []string
	seen := make(map[string]bool)

	// Administrative terms to skip
	adminTerms := []string{
		"universidade", "federal", "minas", "gerais", "escola", "engenharia",
		"colegiado", "curso", "graduação", "mecânica", "civil", "elétrica",
		"pampulha", "belo", "horizonte", "fone", "departamento", "unidade",
		"instituto", "ciências", "exatas", "carga", "horária", "total",
		"créditos", "período", "teórica", "classificação", "obrigatória",
		"prática", "pré-requisitos", "código", "disciplina",
		"ementa", "bibliografia", "professor", "professora", "coordenação",
	}

	for _, topic := range topics {
		topic = strings.TrimSpace(topic)
		topicLower := strings.ToLower(topic)

		// Skip if already seen, too short, or is a stop word
		if seen[topicLower] || len(topic) < 3 || p.isCommonStopWord(topic) {
			continue
		}

		// Skip administrative terms
		isAdmin := false
		for _, adminTerm := range adminTerms {
			if strings.Contains(topicLower, adminTerm) {
				isAdmin = true
				break
			}
		}
		if isAdmin {
			continue
		}

		// Skip if it's just numbers, codes, or single characters
		if regexp.MustCompile(`^[\d\-\.\s]+$`).MatchString(topic) || len(topic) == 1 {
			continue
		}

		filtered = append(filtered, topic)
		seen[topicLower] = true
	}

	return filtered
}

func (p *SyllabusProcessor) removeDuplicates(items []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, item := range items {
		item = strings.TrimSpace(item)
		if !seen[item] && item != "" {
			result = append(result, item)
			seen[item] = true
		}
	}

	return result
}

// =============================================================================
// Enhanced IPFSConnector - Handle IPFS interactions
// =============================================================================

type IPFSConnector struct {
	config     IPFSConfig
	ipfsClient *ipfs.IPFSClient
	logger     *logrus.Logger
}

func NewIPFSConnector(config IPFSConfig) *IPFSConnector {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	var client *ipfs.IPFSClient
	if config.Enabled {
		client = ipfs.NewIPFSClient(config.Endpoint, config.LocalPath, config.UseHTTP)
		logger.WithField("endpoint", config.Endpoint).Info("IPFS client initialized")
	}

	return &IPFSConnector{
		config:     config,
		ipfsClient: client,
		logger:     logger,
	}
}

func (c *IPFSConnector) IsEnabled() bool {
	return c.config.Enabled && c.ipfsClient != nil
}

func (c *IPFSConnector) UploadToIPFS(content []byte) (string, string, error) {
	if c.ipfsClient == nil {
		return "", "", fmt.Errorf("IPFS client not initialized")
	}

	c.logger.WithField("content_size", len(content)).Debug("Uploading content to IPFS")

	_, ipfsLink, err := c.ipfsClient.Add(content)
	if err != nil {
		c.logger.WithError(err).Error("Failed to upload to IPFS")
		return "", "", fmt.Errorf("error uploading to IPFS: %w", err)
	}

	// Extract actual IPFS CID from the link
	ipfsCID := strings.TrimPrefix(ipfsLink, "ipfs://")

	c.logger.WithFields(logrus.Fields{
		"cid":  ipfsCID,
		"link": ipfsLink,
	}).Info("Content uploaded to IPFS successfully")

	// Pin content if configured
	if c.config.Pin {
		// Note: Add pinning logic here if your IPFS client supports it
		c.logger.WithField("cid", ipfsCID).Debug("Content pinned to IPFS")
	}

	return ipfsCID, ipfsLink, nil
}

func (c *IPFSConnector) StoreIPFSMapping(institution, subjectCode, ipfsCID string) error {
	mapping := struct {
		Institution string    `json:"institution"`
		SubjectCode string    `json:"subject_code"`
		IPFSCID     string    `json:"ipfs_cid"`
		Timestamp   time.Time `json:"timestamp"`
		Version     string    `json:"version"`
	}{
		Institution: institution,
		SubjectCode: subjectCode,
		IPFSCID:     ipfsCID,
		Timestamp:   time.Now(),
		Version:     "2.1.0",
	}

	mappingJSON, err := json.MarshalIndent(mapping, "", "  ")
	if err != nil {
		return fmt.Errorf("error serializing IPFS mapping: %w", err)
	}

	mappingsDir := "ipfs_mappings"
	if _, err := os.Stat(mappingsDir); os.IsNotExist(err) {
		if err := os.Mkdir(mappingsDir, 0755); err != nil {
			return fmt.Errorf("error creating mappings directory: %w", err)
		}
	}

	mappingFile := fmt.Sprintf("%s/%s_%s.json", mappingsDir, institution, subjectCode)
	if err := ioutil.WriteFile(mappingFile, mappingJSON, 0644); err != nil {
		return fmt.Errorf("error writing mapping file: %w", err)
	}

	c.logger.WithField("mapping_file", mappingFile).Info("IPFS mapping stored successfully")
	return nil
}

// =============================================================================
// Enhanced SubjectManager - Interact with Subject module
// =============================================================================

type SubjectManager struct {
	acadTokenPath string
	cliConfig     CLIConfig
	logger        *logrus.Logger
}

func NewSubjectManager(acadTokenPath string, cliConfig CLIConfig) *SubjectManager {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	return &SubjectManager{
		acadTokenPath: acadTokenPath,
		cliConfig:     cliConfig,
		logger:        logger,
	}
}

func (m *SubjectManager) CreateSubject(content SubjectContent, prerequisites []PrerequisiteGroup) error {
	m.logger.WithFields(logrus.Fields{
		"subject_id":    content.SubjectId,
		"institution":   content.Institution,
		"course_id":     content.CourseId,
		"prerequisites": len(prerequisites),
	}).Info("Creating subject content")

	// Create the shell script for subject creation
	scriptContent := fmt.Sprintf(`#!/bin/bash
set -e

echo "Creating subject content..."

# Check if the command exists
if ! command -v %s &> /dev/null; then
    echo "Error: academictokend command not found at %s"
    exit 1
fi

# Create subject content
echo "Executing subject creation command..."
%s tx subject create-subject-content \
    '%s' \
    '%s' \
    '%s' \
    '%s' \
    '%s' \
    '%s' \
    %d \
    %d \
    '%s' \
    '%s' \
    '%s' \
    '%s' \
    '%s' \
    '%s' \
    --from '%s' \
    --chain-id '%s' \
    --node '%s' \
    --gas auto \
    --gas-adjustment 1.5 \
    --yes

if [ $? -ne 0 ]; then
    echo "Error creating subject content."
    exit 1
fi

echo "Subject content created successfully!"
`,
		m.acadTokenPath,       // Check command
		m.acadTokenPath,       // Error message
		m.acadTokenPath,       // Execute command
		content.Index,         // 1
		content.SubjectId,     // 2
		content.Institution,   // 3
		content.CourseId,      // 4
		content.Title,         // 5
		content.Code,          // 6
		content.WorkloadHours, // 7
		content.Credits,       // 8
		content.Description,   // 9
		content.ContentHash,   // 10
		content.SubjectType,   // 11
		content.KnowledgeArea, // 12
		content.IpfsLink,      // 13
		content.Creator,       // 14
		m.cliConfig.From,      // --from
		m.cliConfig.ChainID,   // --chain-id
		m.cliConfig.Node,      // --node
	)

	return m.executeScript(scriptContent, "create-subject")
}

func (m *SubjectManager) executeScript(scriptContent, scriptName string) error {
	scriptFile, err := ioutil.TempFile("", fmt.Sprintf("%s-*.sh", scriptName))
	if err != nil {
		return fmt.Errorf("error creating temporary file for script: %w", err)
	}
	defer os.Remove(scriptFile.Name())

	if _, err := scriptFile.Write([]byte(scriptContent)); err != nil {
		return fmt.Errorf("error writing script: %w", err)
	}
	if err := scriptFile.Close(); err != nil {
		return fmt.Errorf("error closing script file: %w", err)
	}

	if err := os.Chmod(scriptFile.Name(), 0755); err != nil {
		return fmt.Errorf("error making script executable: %w", err)
	}

	m.logger.WithField("script", scriptName).Info("Executing blockchain transaction script")
	cmd := exec.Command("bash", scriptFile.Name())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// =============================================================================
// TokenDefManager - Interact with TokenDef module
// =============================================================================

type TokenDefManager struct {
	acadTokenPath string
	cliConfig     CLIConfig
	logger        *logrus.Logger
}

func NewTokenDefManager(acadTokenPath string, cliConfig CLIConfig) *TokenDefManager {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	return &TokenDefManager{
		acadTokenPath: acadTokenPath,
		cliConfig:     cliConfig,
		logger:        logger,
	}
}

func (m *TokenDefManager) CreateTokenDefinition(content SubjectContent) error {
	tokenDefId := fmt.Sprintf("tokendef-%s", content.SubjectId)
	tokenName := fmt.Sprintf("%s Completion Token", content.Title)
	tokenSymbol := fmt.Sprintf("%s-TOKEN", strings.ToUpper(content.Code))

	m.logger.WithFields(logrus.Fields{
		"token_def_id": tokenDefId,
		"token_name":   tokenName,
		"token_symbol": tokenSymbol,
	}).Info("Creating token definition")

	scriptContent := fmt.Sprintf(`#!/bin/bash
set -e

echo "Creating token definition..."

%s tx tokendef create-token-definition \
    '%s' \
    '%s' \
    '%s' \
    '%s' \
    '%s' \
    '%s' \
    '%s' \
    'NFT' \
    'true' \
    'false' \
    '0' \
    '%s' \
    '%s' \
    '%s' \
    '%s' \
    --from '%s' \
    --chain-id '%s' \
    --node '%s' \
    --gas auto \
    --gas-adjustment 1.5 \
    --yes

echo "Token definition created successfully!"
`,
		m.acadTokenPath,
		tokenDefId,          // 1
		tokenDefId,          // 2 - tokenDefId
		content.SubjectId,   // 3 - subjectId
		content.Institution, // 4 - institutionId
		content.CourseId,    // 5 - courseId
		tokenName,           // 6 - tokenName
		tokenSymbol,         // 7 - tokenSymbol
		content.Description, // 12 - description (metadata)
		content.ContentHash, // 13 - contentHash
		content.IpfsLink,    // 14 - ipfsLink
		content.Creator,     // 15 - creator
		m.cliConfig.From,    // --from
		m.cliConfig.ChainID, // --chain-id
		m.cliConfig.Node,    // --node
	)

	return m.executeScript(scriptContent, "create-tokendef")
}

func (m *TokenDefManager) executeScript(scriptContent, scriptName string) error {
	scriptFile, err := ioutil.TempFile("", fmt.Sprintf("%s-*.sh", scriptName))
	if err != nil {
		return fmt.Errorf("error creating temporary file for script: %w", err)
	}
	defer os.Remove(scriptFile.Name())

	if _, err := scriptFile.Write([]byte(scriptContent)); err != nil {
		return fmt.Errorf("error writing script: %w", err)
	}
	if err := scriptFile.Close(); err != nil {
		return fmt.Errorf("error closing script file: %w", err)
	}

	if err := os.Chmod(scriptFile.Name(), 0755); err != nil {
		return fmt.Errorf("error making script executable: %w", err)
	}

	m.logger.WithField("script", scriptName).Info("Executing token definition creation script")
	cmd := exec.Command("bash", scriptFile.Name())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// =============================================================================
// AcademicNFTManager - Interact with AcademicNFT module
// =============================================================================

type AcademicNFTManager struct {
	acadTokenPath string
	cliConfig     CLIConfig
	logger        *logrus.Logger
}

func NewAcademicNFTManager(acadTokenPath string, cliConfig CLIConfig) *AcademicNFTManager {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	return &AcademicNFTManager{
		acadTokenPath: acadTokenPath,
		cliConfig:     cliConfig,
		logger:        logger,
	}
}

func (m *AcademicNFTManager) MintSubjectToken(content SubjectContent, student string) error {
	tokenDefId := fmt.Sprintf("tokendef-%s", content.SubjectId)
	completionDate := time.Now().Format("2006-01-02")

	m.logger.WithFields(logrus.Fields{
		"token_def_id":    tokenDefId,
		"student":         student,
		"completion_date": completionDate,
	}).Info("Minting subject token instance")

	scriptContent := fmt.Sprintf(`#!/bin/bash
set -e

echo "Minting subject token instance..."

%s tx academicnft mint-subject-token-instance \
    '%s' \
    '%s' \
    '%s' \
    '%s' \
    '10.0' \
    '%s' \
    '2024-1' \
    'professor-signature' \
    --from '%s' \
    --chain-id '%s' \
    --node '%s' \
    --gas auto \
    --gas-adjustment 1.5 \
    --yes

echo "Subject token instance minted successfully!"
`,
		m.acadTokenPath,
		fmt.Sprintf("instance-%s-%d", content.SubjectId, time.Now().Unix()), // 1 - index
		tokenDefId,          // 2 - tokenDefId
		student,             // 3 - student
		completionDate,      // 4 - completionDate
		content.Institution, // 6 - issuerInstitution
		m.cliConfig.From,    // --from
		m.cliConfig.ChainID, // --chain-id
		m.cliConfig.Node,    // --node
	)

	return m.executeScript(scriptContent, "mint-token")
}

func (m *AcademicNFTManager) executeScript(scriptContent, scriptName string) error {
	scriptFile, err := ioutil.TempFile("", fmt.Sprintf("%s-*.sh", scriptName))
	if err != nil {
		return fmt.Errorf("error creating temporary file for script: %w", err)
	}
	defer os.Remove(scriptFile.Name())

	if _, err := scriptFile.Write([]byte(scriptContent)); err != nil {
		return fmt.Errorf("error writing script: %w", err)
	}
	if err := scriptFile.Close(); err != nil {
		return fmt.Errorf("error closing script file: %w", err)
	}

	if err := os.Chmod(scriptFile.Name(), 0755); err != nil {
		return fmt.Errorf("error making script executable: %w", err)
	}

	m.logger.WithField("script", scriptName).Info("Executing token minting script")
	cmd := exec.Command("bash", scriptFile.Name())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// =============================================================================
// Enhanced AcademicTokenService - Main orchestration service
// =============================================================================

type AcademicTokenService struct {
	syllabusProcessor  *SyllabusProcessor
	subjectManager     *SubjectManager
	tokenDefManager    *TokenDefManager
	academicNFTManager *AcademicNFTManager
	cliConfig          CLIConfig
	logger             *logrus.Logger
}

func NewAcademicTokenService(
	syllabusProcessor *SyllabusProcessor,
	subjectManager *SubjectManager,
	tokenDefManager *TokenDefManager,
	academicNFTManager *AcademicNFTManager,
	cliConfig CLIConfig,
) *AcademicTokenService {
	logger := logrus.New()
	if cliConfig.Debug {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}

	return &AcademicTokenService{
		syllabusProcessor:  syllabusProcessor,
		subjectManager:     subjectManager,
		tokenDefManager:    tokenDefManager,
		academicNFTManager: academicNFTManager,
		cliConfig:          cliConfig,
		logger:             logger,
	}
}

func (s *AcademicTokenService) ProcessAndTokenize(filePath, institution, courseId, subjectCode, creator string) error {
	s.logger.WithFields(logrus.Fields{
		"file_path":    filePath,
		"institution":  institution,
		"course_id":    courseId,
		"subject_code": subjectCode,
		"creator":      creator,
	}).Info("Starting enhanced syllabus processing and tokenization")

	// 1. Extract text from file with multiple format support
	var rawContent string
	var err error

	s.logger.WithField("file_path", filePath).Info("Extracting text from file")
	rawContent, err = extractTextFromFile(filePath)
	if err != nil {
		s.logger.WithError(err).Error("Failed to extract text from file")
		return fmt.Errorf("error extracting text from file: %v", err)
	}

	s.logger.WithField("content_length", len(rawContent)).Debug("Text extracted successfully")

	// 2. Process syllabus content
	s.logger.Info("Processing syllabus content with enhanced extraction")
	result, err := s.syllabusProcessor.ProcessSyllabus(rawContent, institution, courseId, subjectCode, creator)
	if err != nil {
		s.logger.WithError(err).Error("Failed to process syllabus")
		return fmt.Errorf("error processing syllabus: %v", err)
	}

	// Log detailed results
	s.logger.WithFields(logrus.Fields{
		"content_hash":          result.ContentHash,
		"prerequisites_groups":  len(result.Prerequisites),
		"quality_score":         result.ExtractedData.QualityScore,
		"extraction_confidence": result.ExtractedData.ExtractionConfidence,
		"processing_time_ms":    result.ProcessingMetrics.ProcessingTimeMs,
		"topics_extracted":      result.ProcessingMetrics.TopicsExtracted,
		"keywords_extracted":    result.ProcessingMetrics.KeywordsExtracted,
		"weeks_extracted":       result.ProcessingMetrics.WeeksExtracted,
	}).Info("Syllabus processed successfully")

	if result.IPFSCID != "" {
		s.logger.WithField("ipfs_cid", result.IPFSCID).Info("Content uploaded to IPFS")
	}

	// Print detailed prerequisite information for DAG verification
	for i, prereq := range result.Prerequisites {
		s.logger.WithFields(logrus.Fields{
			"group_id":   prereq.Id,
			"group_type": prereq.GroupType,
			"logic":      prereq.Logic,
			"subjects":   prereq.SubjectIds,
			"confidence": prereq.Confidence,
		}).Infof("Prerequisite Group %d details", i+1)
	}

	// 3. Execute based on configuration
	if s.cliConfig.ProcessOnly {
		s.logger.Info("Processing only mode - stopping here")
		return nil
	}

	return s.executeBlockchainOperations(result)
}

func (s *AcademicTokenService) executeBlockchainOperations(result ProcessingResult) error {
	if s.cliConfig.CreateOnly || (!s.cliConfig.TokenizeOnly) {
		// Create subject content
		s.logger.Info("Step 1: Creating subject content on blockchain")
		if err := s.subjectManager.CreateSubject(result.SubjectContent, result.Prerequisites); err != nil {
			s.logger.WithError(err).Error("Failed to create subject")
			return fmt.Errorf("error creating subject: %v", err)
		}

		s.logger.Info("Waiting for transaction to be processed...")
		time.Sleep(5 * time.Second)

		// Create token definition
		s.logger.Info("Step 2: Creating token definition on blockchain")
		if err := s.tokenDefManager.CreateTokenDefinition(result.SubjectContent); err != nil {
			s.logger.WithError(err).Error("Failed to create token definition")
			return fmt.Errorf("error creating token definition: %v", err)
		}

		s.logger.Info("Waiting for transaction to be processed...")
		time.Sleep(5 * time.Second)
	}

	if s.cliConfig.TokenizeOnly || (!s.cliConfig.CreateOnly) {
		if s.cliConfig.Student == "" {
			return fmt.Errorf("student address required for tokenization")
		}

		// Mint subject token
		s.logger.Info("Step 3: Minting subject token on blockchain")
		if err := s.academicNFTManager.MintSubjectToken(result.SubjectContent, s.cliConfig.Student); err != nil {
			s.logger.WithError(err).Error("Failed to mint subject token")
			return fmt.Errorf("error minting subject token: %v", err)
		}
	}

	s.logger.Info("Complete blockchain flow executed successfully!")
	return nil
}

// =============================================================================
// Enhanced file processing functions
// =============================================================================

func extractTextFromFile(filePath string) (string, error) {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".pdf":
		return extractTextFromPDF(filePath)
	case ".docx":
		return extractTextFromDOCX(filePath)
	case ".txt", ".text":
		return extractTextFromTXT(filePath)
	case ".rtf":
		return extractTextFromRTF(filePath)
	default:
		// Try to read as text file anyway
		return extractTextFromTXT(filePath)
	}
}

func extractTextFromPDF(filePath string) (string, error) {
	_, err := exec.LookPath("pdftotext")
	if err != nil {
		return "", fmt.Errorf("pdftotext not found. Install poppler-utils package: %w", err)
	}

	cmd := exec.Command("pdftotext", "-layout", filePath, "-")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error executing pdftotext: %w", err)
	}

	return string(output), nil
}

func extractTextFromDOCX(filePath string) (string, error) {
	// Try to use pandoc as fallback
	if _, err := exec.LookPath("pandoc"); err == nil {
		cmd := exec.Command("pandoc", "-t", "plain", filePath)
		output, err := cmd.Output()
		if err != nil {
			return "", fmt.Errorf("error executing pandoc for DOCX: %w", err)
		}
		return string(output), nil
	}

	// If pandoc not available, suggest conversion
	return "", fmt.Errorf("DOCX format requires pandoc. Install pandoc or convert to PDF/TXT format")
}

func extractTextFromTXT(filePath string) (string, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("error reading text file: %w", err)
	}
	return string(content), nil
}

func extractTextFromRTF(filePath string) (string, error) {
	// Try to use pandoc as fallback
	if _, err := exec.LookPath("pandoc"); err == nil {
		cmd := exec.Command("pandoc", "-t", "plain", filePath)
		output, err := cmd.Output()
		if err != nil {
			return "", fmt.Errorf("error executing pandoc for RTF: %w", err)
		}
		return string(output), nil
	}

	// If pandoc not available, suggest conversion
	return "", fmt.Errorf("RTF format requires pandoc. Install pandoc or convert to PDF/TXT format")
}

// =============================================================================
// Enhanced main function with comprehensive CLI
// =============================================================================

func main() {
	// Check arguments
	if len(os.Args) < 5 {
		printUsage()
		os.Exit(1)
	}

	// Parse arguments
	filePath := os.Args[1]
	institution := os.Args[2]
	courseId := os.Args[3]
	subjectCode := os.Args[4]

	// Enhanced flag parsing
	var (
		fromFlag         = ""
		chainIDFlag      = "academictoken"
		nodeFlag         = "tcp://localhost:26657"
		studentFlag      = ""
		processOnlyFlag  = false
		createOnlyFlag   = false
		tokenizeOnlyFlag = false
		noIPFSFlag       = false
		debugFlag        = false
		noCacheFlag      = false
	)

	// Parse flags with validation
	for i, arg := range os.Args {
		switch arg {
		case "--from":
			if i+1 < len(os.Args) {
				fromFlag = os.Args[i+1]
			}
		case "--chain-id":
			if i+1 < len(os.Args) {
				chainIDFlag = os.Args[i+1]
			}
		case "--node":
			if i+1 < len(os.Args) {
				nodeFlag = os.Args[i+1]
			}
		case "--student":
			if i+1 < len(os.Args) {
				studentFlag = os.Args[i+1]
			}
		case "--process-only":
			processOnlyFlag = true
		case "--create-only":
			createOnlyFlag = true
		case "--tokenize-only":
			tokenizeOnlyFlag = true
		case "--no-ipfs":
			noIPFSFlag = true
		case "--debug":
			debugFlag = true
		case "--no-cache":
			noCacheFlag = true
		case "--help", "-h":
			printUsage()
			os.Exit(0)
		}
	}

	// Validate file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("Error: File %s does not exist\n", filePath)
		os.Exit(1)
	}

	// Validate flags
	if !processOnlyFlag && fromFlag == "" {
		fmt.Println("Error: --from flag is required for blockchain operations")
		os.Exit(1)
	}

	if tokenizeOnlyFlag && studentFlag == "" {
		fmt.Println("Error: --student flag is required for tokenization")
		os.Exit(1)
	}

	// IPFS configuration with enhanced detection
	ipfsConfig := IPFSConfig{
		Enabled:   false,
		Endpoint:  "http://localhost:5001",
		LocalPath: "/Users/biancamsp/.academictoken/ipfs_local",
		UseHTTP:   true,
		Pin:       false,
	}

	// Check IPFS availability if not disabled
	if !noIPFSFlag {
		if err := checkIPFSAvailability(); err == nil {
			ipfsConfig.Enabled = true
			fmt.Println("IPFS enabled for distributed storage")
		} else {
			fmt.Printf("IPFS not available: %v - using local storage only\n", err)
		}
	}

	// CLI configuration
	cliConfig := CLIConfig{
		ChainID:      chainIDFlag,
		From:         fromFlag,
		Node:         nodeFlag,
		Student:      studentFlag,
		ProcessOnly:  processOnlyFlag,
		CreateOnly:   createOnlyFlag,
		TokenizeOnly: tokenizeOnlyFlag,
		Debug:        debugFlag,
		CacheEnabled: !noCacheFlag,
	}

	// Validate binary path
	acadTokenPath := "/Users/biancamsp/go/bin/academictokend"
	if err := validateBinaryPath(acadTokenPath); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Initialize components with enhanced configuration
	ipfsConnector := NewIPFSConnector(ipfsConfig)
	syllabusProcessor := NewSyllabusProcessor(ipfsConnector, cliConfig.CacheEnabled)

	if debugFlag {
		syllabusProcessor.SetLogLevel(logrus.DebugLevel)
	}

	subjectManager := NewSubjectManager(acadTokenPath, cliConfig)
	tokenDefManager := NewTokenDefManager(acadTokenPath, cliConfig)
	academicNFTManager := NewAcademicNFTManager(acadTokenPath, cliConfig)

	// Create main service
	service := NewAcademicTokenService(
		syllabusProcessor,
		subjectManager,
		tokenDefManager,
		academicNFTManager,
		cliConfig,
	)

	// Execute processing
	fmt.Println("Starting enhanced syllabus processing and tokenization...")
	fmt.Printf("Using binary: %s\n", acadTokenPath)
	fmt.Printf("Data directory: /Users/biancamsp/.academictoken\n")
	fmt.Printf("Cache enabled: %v\n", cliConfig.CacheEnabled)
	fmt.Printf("Debug mode: %v\n", debugFlag)

	if err := service.ProcessAndTokenize(filePath, institution, courseId, subjectCode, fromFlag); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Process completed successfully!")
}

// Helper functions for enhanced main
func printUsage() {
	fmt.Println("Enhanced Academic Token Syllabus Processor v2.1.0")
	fmt.Println("Usage: syllabus_processor <syllabus_file> <institution> <course_id> <subject_code> [options]")
	fmt.Println("Example: ./syllabus_processor syllabus.pdf UFBA CS101 CALC1 --chain-id academictoken --from alice")
	fmt.Println("\nSupported file formats: PDF, TXT, DOCX (experimental), RTF (experimental)")
	fmt.Println("\nOptions:")
	fmt.Println("  --process-only        Only process and extract data, don't create blockchain records")
	fmt.Println("  --create-only         Only create subject and token definition, don't mint tokens")
	fmt.Println("  --tokenize-only       Only mint tokens, assumes subject and token def exist")
	fmt.Println("  --student <address>   Student address for token minting")
	fmt.Println("  --chain-id <chain-id> Chain ID (default: academictoken)")
	fmt.Println("  --from <key-name>     Key name for signing transactions")
	fmt.Println("  --node <node>         RPC node address (default: tcp://localhost:26657)")
	fmt.Println("  --no-ipfs             Disable IPFS storage")
	fmt.Println("  --debug               Enable debug logging")
	fmt.Println("  --no-cache            Disable processing cache")
	fmt.Println("  --help, -h            Show this help message")
	fmt.Println("\nNew Features in v2.1.0:")
	fmt.Println("  • Enhanced title extraction with proper cleaning")
	fmt.Println("  • Improved topic extraction from weekly program content")
	fmt.Println("  • Better structured content analysis")
	fmt.Println("  • Enhanced quality validation and reporting")
	fmt.Println("  • Comprehensive program content modeling")
	fmt.Println("  • Detailed logging with processing metrics")
}

func checkIPFSAvailability() error {
	_, err := exec.Command("ipfs", "version").Output()
	if err != nil {
		return fmt.Errorf("IPFS daemon not running or not installed")
	}

	// Test connection to IPFS API
	_, err = exec.Command("curl", "-s", "http://localhost:5001/api/v0/version").Output()
	if err != nil {
		return fmt.Errorf("IPFS API not accessible at localhost:5001")
	}

	return nil
}

func validateBinaryPath(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("binary not found at %s. Please ensure academictokend is installed", path)
	}

	// Test if binary is executable
	if err := exec.Command(path, "version").Run(); err != nil {
		return fmt.Errorf("binary at %s is not executable or not working properly", path)
	}

	return nil
}
