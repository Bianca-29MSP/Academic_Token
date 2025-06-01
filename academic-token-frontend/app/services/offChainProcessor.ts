// services/offChainProcessor.ts
// Integration with off-chain subject content processor

export interface ProcessingResult {
  contentHash: string;
  ipfsCid?: string;
  subjectContent: SubjectContent;
  extractedData: ExtractedSyllabusData;
  prerequisites: PrerequisiteGroup[];
  localJsonPath: string;
  ipfsJsonPath?: string;
  processingMetrics: ProcessingMetrics;
}

export interface SubjectContent {
  index: string;
  subjectId: string;
  institution: string;
  courseId: string;
  title: string;
  code: string;
  workloadHours: number;
  credits: number;
  description: string;
  contentHash: string;
  subjectType: string;
  knowledgeArea: string;
  ipfsLink?: string;
  creator: string;
}

export interface ExtractedSyllabusData {
  title: string;
  workloadHours: number;
  credits: number;
  description: string;
  objectives: string[];
  methodologies: string[];
  evaluationMethods: string[];
  basicBibliography: string[];
  complementaryBibliography: string[];
  topics: string[];
  keywords: string[];
  subjectType: string;
  knowledgeArea: string;
  prerequisites: string[];
  qualityScore: number;
  extractionConfidence: number;
  detailedTopics: string[];
  programContent: WeeklyContent[];
}

export interface WeeklyContent {
  week: number;
  content: string;
  topics: string[];
}

export interface PrerequisiteGroup {
  id: string;
  subjectId: string;
  groupType: string; // "ALL", "ANY", "MINIMUM", "NONE"
  minimumCredits: number;
  minimumCompletedSubjects: number;
  subjectIds: string[];
  logic: string; // "AND", "OR", "XOR", "THRESHOLD"
  priority: number;
  confidence: number;
}

export interface ProcessingMetrics {
  processingTimeMs: number;
  textLengthChars: number;
  sectionsIdentified: number;
  topicsExtracted: number;
  keywordsExtracted: number;
  bibliographyEntries: number;
  prerequisiteGroups: number;
  qualityScore: number;
  extractionConfidence: number;
  weeksExtracted: number;
}

export interface ProcessorConfig {
  processorPath: string;
  tempDir: string;
  enableIPFS: boolean;
  chainId: string;
  from: string;
  node: string;
  debug: boolean;
}

export class OffChainProcessor {
  private config: ProcessorConfig;
  private isProcessing: boolean = false;

  constructor(config: Partial<ProcessorConfig> = {}) {
    this.config = {
      processorPath: "/Users/biancamsp/Desktop/Academic_Token/academictoken/academictoken/tools/off_chain_processor/subject_content_processor",
      tempDir: "/tmp/academic_token_processing",
      enableIPFS: true,
      chainId: "academictoken",
      from: "admin",
      node: "tcp://localhost:26657",
      debug: true,
      ...config
    };
  }

  /**
   * Process a syllabus file using the off-chain processor
   */
  async processFile(
    file: File, 
    institution: string, 
    courseId: string, 
    subjectCode: string,
    options: {
      processOnly?: boolean;
      createOnly?: boolean;
      tokenizeOnly?: boolean;
      student?: string;
    } = {}
  ): Promise<ProcessingResult> {
    if (this.isProcessing) {
      throw new Error("Another file is currently being processed");
    }

    this.isProcessing = true;

    try {
      console.log("🔄 Starting off-chain processing:", {
        fileName: file.name,
        institution,
        courseId,
        subjectCode,
        options
      });

      // Create a temporary file for processing
      const tempFilePath = await this.saveFileTemporarily(file);

      try {
        // Build command arguments
        const args = this.buildProcessorArguments(
          tempFilePath,
          institution,
          courseId,
          subjectCode,
          options
        );

        console.log("🚀 Executing processor with args:", args);

        // Execute the processor
        const result = await this.executeProcessor(args);

        console.log("✅ Processing completed successfully:", result);

        return result;

      } finally {
        // Clean up temporary file
        await this.cleanupTempFile(tempFilePath);
      }

    } catch (error) {
      console.error("❌ Processing failed:", error);
      throw new Error(`Off-chain processing failed: ${error instanceof Error ? error.message : 'Unknown error'}`);
    } finally {
      this.isProcessing = false;
    }
  }

  /**
   * Check if the processor is available and configured correctly
   */
  async checkProcessorAvailability(): Promise<{
    available: boolean;
    version?: string;
    ipfsEnabled?: boolean;
    error?: string;
  }> {
    try {
      // This would typically check if the Go binary exists and is executable
      // For now, we'll simulate the check
      console.log("🔍 Checking processor availability...");

      // In a real implementation, you would:
      // 1. Check if the binary exists at the specified path
      // 2. Try to execute it with --help or --version
      // 3. Check IPFS availability
      
      // Simulated response
      return {
        available: true,
        version: "2.1.0",
        ipfsEnabled: true
      };

    } catch (error) {
      return {
        available: false,
        error: error instanceof Error ? error.message : 'Unknown error'
      };
    }
  }

  /**
   * Get processing status for monitoring long operations
   */
  getProcessingStatus(): {
    isProcessing: boolean;
    canProcess: boolean;
  } {
    return {
      isProcessing: this.isProcessing,
      canProcess: !this.isProcessing
    };
  }

  private async saveFileTemporarily(file: File): Promise<string> {
    // In a real implementation, this would save the file to the file system
    // For this demo, we'll create a mock path
    const timestamp = Date.now();
    const fileName = file.name.replace(/[^a-zA-Z0-9.-]/g, '_');
    const tempPath = `${this.config.tempDir}/${timestamp}_${fileName}`;
    
    console.log("💾 Saving file temporarily to:", tempPath);
    
    // In production, you would actually write the file:
    // const buffer = await file.arrayBuffer();
    // await fs.writeFile(tempPath, new Uint8Array(buffer));
    
    return tempPath;
  }

  private buildProcessorArguments(
    filePath: string,
    institution: string,
    courseId: string,
    subjectCode: string,
    options: {
      processOnly?: boolean;
      createOnly?: boolean;
      tokenizeOnly?: boolean;
      student?: string;
    }
  ): string[] {
    const args = [
      filePath,
      institution,
      courseId,
      subjectCode
    ];

    // Add blockchain configuration
    args.push("--chain-id", this.config.chainId);
    args.push("--from", this.config.from);
    args.push("--node", this.config.node);

    // Add processing options
    if (options.processOnly) {
      args.push("--process-only");
    }
    if (options.createOnly) {
      args.push("--create-only");
    }
    if (options.tokenizeOnly) {
      args.push("--tokenize-only");
    }
    if (options.student) {
      args.push("--student", options.student);
    }

    // Add IPFS configuration
    if (!this.config.enableIPFS) {
      args.push("--no-ipfs");
    }

    // Add debug flag
    if (this.config.debug) {
      args.push("--debug");
    }

    return args;
  }

  private async executeProcessor(args: string[]): Promise<ProcessingResult> {
    console.log("⚙️ Executing processor command:", this.config.processorPath, args.join(" "));

    // In a real implementation, this would execute the Go binary
    // For this demo, we'll simulate the processing and return mock data
    
    // Simulate processing time
    await new Promise(resolve => setTimeout(resolve, 2000));

    // Mock successful processing result
    const mockResult: ProcessingResult = {
      contentHash: "abc123def456789",
      ipfsCid: this.config.enableIPFS ? "QmYwAPJzv5CZsnA625s3Xf2nemtYgPpHdWEz79ojWnPbdG" : undefined,
      subjectContent: {
        index: `${args[1]}-${args[2]}-${args[3]}`,
        subjectId: `${args[1]}-${args[2]}-${args[3]}`,
        institution: args[1],
        courseId: args[2],
        title: this.extractSubjectName(args[0]),
        code: args[3],
        workloadHours: 60,
        credits: 4,
        description: "Enhanced description extracted from syllabus content",
        contentHash: "abc123def456789",
        subjectType: "THEORETICAL",
        knowledgeArea: "COMPUTER_SCIENCE",
        ipfsLink: this.config.enableIPFS ? "ipfs://QmYwAPJzv5CZsnA625s3Xf2nemtYgPpHdWEz79ojWnPbdG" : undefined,
        creator: this.config.from
      },
      extractedData: {
        title: this.extractSubjectName(args[0]),
        workloadHours: 60,
        credits: 4,
        description: "Comprehensive syllabus description extracted using NLP techniques",
        objectives: [
          "Understand fundamental concepts",
          "Apply theoretical knowledge in practice",
          "Develop problem-solving skills"
        ],
        methodologies: [
          "Lectures with visual presentations",
          "Practical exercises and labs",
          "Group discussions and projects"
        ],
        evaluationMethods: [
          "Written examinations (60%)",
          "Practical assignments (30%)",
          "Participation and attendance (10%)"
        ],
        basicBibliography: [
          "Core Textbook - Author Name, 2023",
          "Essential Reference - Another Author, 2022"
        ],
        complementaryBibliography: [
          "Additional Reading - Third Author, 2021",
          "Supplementary Material - Fourth Author, 2020"
        ],
        topics: [
          "Introduction to concepts",
          "Fundamental principles",
          "Advanced applications",
          "Case studies and examples"
        ],
        keywords: [
          "algorithm", "data", "structure", "programming", "analysis"
        ],
        subjectType: "THEORETICAL",
        knowledgeArea: "COMPUTER_SCIENCE",
        prerequisites: [],
        qualityScore: 8.5,
        extractionConfidence: 0.92,
        detailedTopics: [
          "Week 1: Course introduction and overview",
          "Week 2: Basic concepts and terminology",
          "Week 3: Fundamental algorithms",
          "Week 4: Data structures implementation"
        ],
        programContent: [
          {
            week: 1,
            content: "Course introduction and overview",
            topics: ["Introduction", "Course objectives", "Assessment methods"]
          },
          {
            week: 2,
            content: "Basic concepts and terminology",
            topics: ["Definitions", "Basic concepts", "Key terminology"]
          }
        ]
      },
      prerequisites: [
        {
          id: `prereq-${args[1]}-${args[2]}-${args[3]}-1`,
          subjectId: `${args[1]}-${args[2]}-${args[3]}`,
          groupType: "NONE",
          minimumCredits: 0,
          minimumCompletedSubjects: 0,
          subjectIds: [],
          logic: "NONE",
          priority: 0,
          confidence: 1.0
        }
      ],
      localJsonPath: `/tmp/academic_token_processing/${args[1]}_${args[3]}_complete.json`,
      ipfsJsonPath: this.config.enableIPFS ? "ipfs://QmYwAPJzv5CZsnA625s3Xf2nemtYgPpHdWEz79ojWnPbdG" : undefined,
      processingMetrics: {
        processingTimeMs: 2000,
        textLengthChars: 5000,
        sectionsIdentified: 8,
        topicsExtracted: 15,
        keywordsExtracted: 25,
        bibliographyEntries: 12,
        prerequisiteGroups: 1,
        qualityScore: 8.5,
        extractionConfidence: 0.92,
        weeksExtracted: 16
      }
    };

    console.log("📊 Processing metrics:", mockResult.processingMetrics);

    return mockResult;
  }

  private async cleanupTempFile(filePath: string): Promise<void> {
    console.log("🧹 Cleaning up temporary file:", filePath);
    // In production: await fs.unlink(filePath);
  }

  private extractSubjectName(filePath: string): string {
    // Extract a reasonable subject name from filename
    const fileName = filePath.split('/').pop() || '';
    const nameWithoutExt = fileName.replace(/\.(pdf|txt|docx|rtf)$/i, '');
    
    // Clean up the name
    const cleaned = nameWithoutExt
      .replace(/[_-]/g, ' ')
      .replace(/\d+/g, '')
      .trim();
    
    return cleaned || 'Extracted Subject';
  }
}

// Factory for creating processor instances
export class ProcessorFactory {
  private static instance: OffChainProcessor | null = null;

  static getInstance(config?: Partial<ProcessorConfig>): OffChainProcessor {
    if (!ProcessorFactory.instance) {
      ProcessorFactory.instance = new OffChainProcessor(config);
    }
    return ProcessorFactory.instance;
  }

  static createProcessor(config?: Partial<ProcessorConfig>): OffChainProcessor {
    return new OffChainProcessor(config);
  }
}

// Default processor instance
export const offChainProcessor = ProcessorFactory.getInstance();
