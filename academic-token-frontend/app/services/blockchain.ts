// services/blockchain.ts
// Base configuration for Cosmos backend connection

export interface CosmosConfig {
  rpcEndpoint: string;
  chainId: string;
  denom: string;
  gasPrice: string;
}

export const COSMOS_CONFIG: CosmosConfig = {
  rpcEndpoint: process.env.NEXT_PUBLIC_COSMOS_RPC || 'http://localhost:26657',
  chainId: process.env.NEXT_PUBLIC_CHAIN_ID || 'academictoken',
  denom: process.env.NEXT_PUBLIC_DENOM || 'utoken',
  gasPrice: process.env.NEXT_PUBLIC_GAS_PRICE || '0.025utoken'
};

// API Configuration - uses REST server for academic data
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:1318';

// Function to get API URL
function getApiUrl(): string {
  return API_BASE_URL;
}

// Types based on your architecture
export interface Institution {
  id: string;
  name: string;
  code: string;
  country: string;
  createdAt: string;
}

export interface Course {
  id: string;
  institutionId: string;
  name: string;
  code: string;
  duration: number;
  totalCredits: number;
}

export interface Subject {
  id: string;
  courseId: string;
  institutionId: string;
  name: string;
  code: string;
  credits: number;
  syllabus: string; // IPFS hash
  metadata: string; // on-chain metadata
}

export interface Student {
  id: string;
  institutionId: string;
  name: string;
  email: string;
  courseId: string;
  curriculumId: string;
  enrollmentDate: string;
}

export interface AcademicNFT {
  id: string;
  studentId: string;
  subjectId: string;
  grade: number;
  completionDate: string;
  nftHash: string;
  metadata: {
    subject: string;
    credits: number;
    institution: string;
  };
}

export interface SubjectCompletion {
  subjectId: string;
  grade: number;
  completionProof: string;
}

export interface EquivalenceRequest {
  id: string;
  studentId: string;
  sourceSubjectId: string;
  targetSubjectId: string;
  status: 'pending' | 'approved' | 'rejected';
  similarity: number;
}

// Wallet Service for transaction signing
export class WalletService {
  private walletAddress: string | null = null;

  async connectWallet(): Promise<string> {
    // TODO: Connect to real wallet (Keplr/Leap)
    // For now, use a demo address until wallet integration is complete
    this.walletAddress = "cosmos1demo123address456";
    console.log('⚠️ Using demo wallet address - implement real wallet connection');
    return this.walletAddress;
  }

  getWalletAddress(): string | null {
    return this.walletAddress;
  }

  async signTransaction(txData: any): Promise<any> {
    // TODO: Use real wallet signing
    console.log('⚠️ Using mock signing - implement real wallet signing');
    console.log('Signing transaction:', txData);
    return {
      ...txData,
      signatures: ["mock_signature"]
    };
  }
}

// Transaction Service for blockchain operations
export class TransactionService {
  private baseUrl: string;
  private walletService: WalletService;

  constructor(baseUrl: string = getApiUrl()) {
    this.baseUrl = baseUrl;
    this.walletService = new WalletService();
  }

  async submitTransaction(msg: any, msgType: string): Promise<any> {
    try {
      console.log('📝 Preparing transaction:', { msgType, msg });
      
      // Validate message data
      if (!msg || Object.keys(msg).length === 0) {
        throw new Error('Empty message data provided');
      }
      
      // Ensure wallet is connected
      const walletAddress = this.walletService.getWalletAddress() || await this.walletService.connectWallet();
      
      if (!walletAddress) {
        throw new Error('No wallet address available');
      }

      const txBody = {
        body: {
          messages: [{
            "@type": msgType,
            creator: walletAddress,
            ...msg
          }],
          memo: "",
          timeout_height: "0",
          extension_options: [],
          non_critical_extension_options: []
        },
        auth_info: {
          signer_infos: [],
          fee: {
            amount: [{ denom: COSMOS_CONFIG.denom, amount: "5000" }],
            gas_limit: "200000",
            payer: "",
            granter: ""
          }
        },
        signatures: []
      };

      console.log('📝 Transaction body prepared:', txBody);

      // Sign transaction
      const signedTx = await this.walletService.signTransaction(txBody);

      console.log('✍️ Transaction signed, submitting to blockchain...');

      // Submit to blockchain
      const response = await fetch(`${this.baseUrl}/cosmos/tx/v1beta1/txs`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(signedTx),
      });

      console.log('📞 Blockchain response status:', response.status);

      if (!response.ok) {
        const error = await response.text();
        console.error('❌ Blockchain rejected transaction:', error);
        throw new Error(`Transaction failed: ${response.status} - ${error}`);
      }

      const result = await response.json();
      
      console.log('📝 Blockchain response:', result);
      
      // Check if transaction was successful
      if (result.tx_response?.code !== 0) {
        console.error('❌ Transaction execution failed:', result.tx_response);
        throw new Error(`Transaction failed: ${result.tx_response?.raw_log || 'Unknown error'}`);
      }

      console.log('✅ Transaction successful!');
      return result;
    } catch (error) {
      console.error('❌ Transaction submission failed:', error);
      throw error; // Don't return mock data - let the error propagate
    }
  }
}

// Service Classes
export class InstitutionService {
  private baseUrl: string;
  private txService: TransactionService;
  
  constructor(baseUrl: string = getApiUrl()) {
    this.baseUrl = baseUrl;
    this.txService = new TransactionService(baseUrl);
  }

  async registerInstitution(institutionData: Omit<Institution, 'id' | 'createdAt'>): Promise<Institution> {
    try {
      const msg = {
        name: institutionData.name,
        code: institutionData.code,
        country: institutionData.country
      };

      const result = await this.txService.submitTransaction(
        msg, 
        "/academictoken.institution.MsgRegisterInstitution"
      );

      // Extract institution ID from events
      const events = result.tx_response?.events || [];
      let institutionId = "";
      
      for (const event of events) {
        if (event.type === "institution_registered") {
          const idAttr = event.attributes.find((attr: any) => attr.key === "institution_id");
          if (idAttr) {
            institutionId = idAttr.value;
            break;
          }
        }
      }

      return {
        id: institutionId || `inst-${Date.now()}`,
        ...institutionData,
        createdAt: new Date().toISOString()
      };
    } catch (error) {
      console.error('Failed to register institution:', error);
      throw error;
    }
  }

  async getInstitution(id: string): Promise<Institution> {
    const response = await fetch(`${this.baseUrl}/academictoken/institution/institution/${id}`);
    
    if (!response.ok) {
      throw new Error(`Failed to get institution: ${response.statusText}`);
    }
    
    const result = await response.json();
    return result.institution;
  }

  async listInstitutions(): Promise<Institution[]> {
    const response = await fetch(`${this.baseUrl}/academictoken/institution/institution`);
    
    if (!response.ok) {
      throw new Error(`Failed to list institutions: ${response.statusText}`);
    }
    
    const result = await response.json();
    return result.institutions || [];
  }
}

export class CourseService {
  private baseUrl: string;
  private txService: TransactionService;
  
  constructor(baseUrl: string = getApiUrl()) {
    this.baseUrl = baseUrl;
    this.txService = new TransactionService(baseUrl);
  }

  async createCourse(courseData: Omit<Course, 'id'>): Promise<Course> {
    try {
      const msg = {
        institution_id: courseData.institutionId,
        name: courseData.name,
        code: courseData.code,
        duration: courseData.duration,
        total_credits: courseData.totalCredits
      };

      const result = await this.txService.submitTransaction(
        msg, 
        "/academictoken.course.MsgCreateCourse"
      );

      return {
        id: `course-${Date.now()}`,
        ...courseData
      };
    } catch (error) {
      console.error('Failed to create course:', error);
      throw error;
    }
  }

  async getCourse(id: string): Promise<Course> {
    const response = await fetch(`${this.baseUrl}/academictoken/course/course/${id}`);
    
    if (!response.ok) {
      throw new Error(`Failed to get course: ${response.statusText}`);
    }
    
    const result = await response.json();
    return result.course;
  }

  async listCourses(institutionId?: string): Promise<Course[]> {
    const url = institutionId 
      ? `${this.baseUrl}/academictoken/course/course?institution_id=${institutionId}`
      : `${this.baseUrl}/academictoken/course/course`;
      
    const response = await fetch(url);
    
    if (!response.ok) {
      throw new Error(`Failed to list courses: ${response.statusText}`);
    }
    
    const result = await response.json();
    return result.courses || [];
  }
}

export class SubjectService {
  private baseUrl: string;
  private txService: TransactionService;
  
  constructor(baseUrl: string = getApiUrl()) {
    this.baseUrl = baseUrl;
    this.txService = new TransactionService(baseUrl);
  }

  async registerSubject(
    subjectData: Omit<Subject, 'id'>, 
    syllabusFile?: File
  ): Promise<Subject> {
    try {
      // Validate input data
      if (!subjectData.name || !subjectData.code || !subjectData.courseId) {
        throw new Error('Missing required subject data: name, code, or courseId');
      }

      // First, upload syllabus to IPFS if provided
      let ipfsHash = "";
      if (syllabusFile) {
        ipfsHash = await this.uploadToIPFS(syllabusFile);
      }

      const msg = {
        course_id: subjectData.courseId,
        institution_id: subjectData.institutionId,
        name: subjectData.name,
        code: subjectData.code,
        credits: subjectData.credits || 4,
        syllabus: ipfsHash || subjectData.syllabus || '',
        metadata: subjectData.metadata || ''
      };

      console.log('📄 Submitting subject registration:', msg);

      const result = await this.txService.submitTransaction(
        msg, 
        "/academictoken.subject.MsgRegisterSubject"
      );

      return {
        id: `subject-${Date.now()}`,
        ...subjectData,
        syllabus: ipfsHash || subjectData.syllabus || ''
      };
    } catch (error) {
      console.error('Failed to register subject:', error);
      throw error; // Don't return fallback data - let the error propagate
    }
  }

  private async uploadToIPFS(file: File): Promise<string> {
    try {
      // TODO: Implement real IPFS upload
      console.log('⚠️ IPFS upload not yet implemented - need to integrate IPFS client');
      throw new Error('IPFS upload not implemented yet');
    } catch (error) {
      console.error('IPFS upload failed:', error);
      throw new Error('Failed to upload to IPFS');
    }
  }

  async getSubject(id: string): Promise<Subject> {
    const response = await fetch(`${this.baseUrl}/academictoken/subject/subject/${id}`);
    
    if (!response.ok) {
      throw new Error(`Failed to get subject: ${response.statusText}`);
    }
    
    const result = await response.json();
    return result.subject;
  }

  async listSubjects(courseId?: string): Promise<Subject[]> {
    const url = courseId 
      ? `${this.baseUrl}/academictoken/subject/subject?course_id=${courseId}`
      : `${this.baseUrl}/academictoken/subject/subject`;
      
    const response = await fetch(url);
    
    if (!response.ok) {
      throw new Error(`Failed to list subjects: ${response.statusText}`);
    }
    
    const result = await response.json();
    return result.subjects || [];
  }

  async getSyllabus(ipfsHash: string): Promise<string> {
    try {
      // TODO: Implement real IPFS retrieval
      console.log('⚠️ IPFS retrieval not yet implemented - need to integrate IPFS client');
      throw new Error('IPFS retrieval not implemented yet');
    } catch (error) {
      console.error('Failed to get syllabus:', error);
      throw new Error('Failed to retrieve syllabus from IPFS');
    }
  }
}

export class StudentService {
  private baseUrl: string;
  private txService: TransactionService;
  
  constructor(baseUrl: string = getApiUrl()) {
    this.baseUrl = baseUrl;
    this.txService = new TransactionService(baseUrl);
  }

  async registerStudent(studentData: Omit<Student, 'id'>): Promise<Student> {
    try {
      const msg = {
        institution_id: studentData.institutionId,
        name: studentData.name,
        email: studentData.email,
        course_id: studentData.courseId,
        curriculum_id: studentData.curriculumId
      };

      const result = await this.txService.submitTransaction(
        msg, 
        "/academictoken.student.MsgRegisterStudent"
      );

      return {
        id: `student-${Date.now()}`,
        ...studentData,
        enrollmentDate: new Date().toISOString()
      };
    } catch (error) {
      console.error('Failed to register student:', error);
      throw error;
    }
  }

  async enrollInCourse(studentId: string, courseId: string, curriculumId: string): Promise<void> {
    try {
      const msg = {
        student_id: studentId,
        course_id: courseId,
        curriculum_id: curriculumId
      };

      await this.txService.submitTransaction(
        msg, 
        "/academictoken.student.MsgEnrollInCourse"
      );
    } catch (error) {
      console.error('Failed to enroll in course:', error);
      throw error;
    }
  }

  async enrollInSubject(studentId: string, subjectId: string): Promise<void> {
    try {
      const msg = {
        student: studentId,
        subject_id: subjectId
      };

      await this.txService.submitTransaction(
        msg, 
        "/academictoken.student.MsgRequestSubjectEnrollment"
      );
    } catch (error) {
      console.error('Failed to enroll in subject:', error);
      throw error;
    }
  }

  async markSubjectComplete(
    studentId: string, 
    completion: SubjectCompletion & { semester: string; completionDate: string }
  ): Promise<AcademicNFT> {
    try {
      const msg = {
        student_id: studentId,
        subject_id: completion.subjectId,
        grade: Math.round(completion.grade),
        completion_date: completion.completionDate,
        semester: completion.semester
      };

      const result = await this.txService.submitTransaction(
        msg, 
        "/academictoken.student.MsgCompleteSubject"
      );

      // Extract NFT data from events
      const events = result.tx_response?.events || [];
      let nftTokenId = "";
      
      for (const event of events) {
        if (event.type === "subject_completed") {
          const nftAttr = event.attributes.find((attr: any) => attr.key === "nft_token_id");
          if (nftAttr) {
            nftTokenId = nftAttr.value;
            break;
          }
        }
      }

      return {
        id: nftTokenId || `nft-${Date.now()}`,
        studentId: studentId,
        subjectId: completion.subjectId,
        grade: completion.grade,
        completionDate: completion.completionDate,
        nftHash: result.tx_response?.txhash || `hash-${Date.now()}`,
        metadata: {
          subject: completion.subjectId,
          credits: 4, // Default - should be fetched from subject data
          institution: 'UFJF'
        }
      };
    } catch (error) {
      console.error('Failed to complete subject:', error);
      throw error;
    }
  }

  async getStudent(id: string): Promise<Student> {
    const response = await fetch(`${this.baseUrl}/academictoken/student/student/${id}`);
    
    if (!response.ok) {
      throw new Error(`Failed to get student: ${response.statusText}`);
    }
    
    const result = await response.json();
    return result.student;
  }

  async getStudentNFTs(studentId: string): Promise<AcademicNFT[]> {
    const response = await fetch(`${this.baseUrl}/academictoken/academicnft/student/${studentId}/tokens`);
    
    if (!response.ok) {
      throw new Error(`Failed to get student NFTs: ${response.statusText}`);
    }
    
    const result = await response.json();
    
    // Convert SubjectTokenInstance to AcademicNFT
    return (result.token_instances || []).map((token: any) => ({
      id: token.index,
      studentId: token.student,
      subjectId: token.token_def_id,
      grade: parseFloat(token.grade),
      completionDate: token.completion_date,
      nftHash: token.index,
      metadata: {
        subject: token.token_def_id,
        credits: 4, // Default
        institution: token.issuer_institution
      }
    }));
  }

  async checkPrerequisites(studentId: string, subjectId: string): Promise<boolean> {
    try {
      // Call prerequisites contract
      const response = await fetch(`${this.baseUrl}/cosmwasm/wasm/v1/contract/cosmos1prereq.../smart/ewogICJjaGVja19wcmVyZXF1aXNpdGVzIjogewogICAgInN0dWRlbnRfaWQiOiAiJHtzdHVkZW50SWR9IiwKICAgICJzdWJqZWN0X2lkIjogIiR7c3ViamVjdElkfSIKICB9Cn0=`);
      
      if (!response.ok) {
        throw new Error(`Failed to check prerequisites: ${response.statusText}`);
      }
      
      const result = await response.json();
      return result.data.eligible || false;
    } catch (error) {
      console.error('Prerequisites check failed:', error);
      throw new Error('Failed to check prerequisites via smart contract');
    }
  }
}

export class EquivalenceService {
  private baseUrl: string;
  private txService: TransactionService;
  
  constructor(baseUrl: string = getApiUrl()) {
    this.baseUrl = baseUrl;
    this.txService = new TransactionService(baseUrl);
  }

  async requestEquivalence(
    studentId: string,
    sourceSubjectId: string,
    targetSubjectId: string
  ): Promise<EquivalenceRequest> {
    try {
      const msg = {
        student_id: studentId,
        source_subject_id: sourceSubjectId,
        target_subject_id: targetSubjectId,
        reason: "Student equivalence request"
      };

      const result = await this.txService.submitTransaction(
        msg, 
        "/academictoken.student.MsgRequestEquivalence"
      );

      return {
        id: result.tx_response?.txhash || `eq-${Date.now()}`,
        studentId,
        sourceSubjectId,
        targetSubjectId,
        status: 'pending',
        similarity: 0.85 // Mock - would be calculated by equivalence contract
      };
    } catch (error) {
      console.error('Failed to request equivalence:', error);
      throw error;
    }
  }

  async getEquivalenceRequests(studentId: string): Promise<EquivalenceRequest[]> {
    const response = await fetch(`${this.baseUrl}/academictoken/equivalence/requests/${studentId}`);
    
    if (!response.ok) {
      throw new Error(`Failed to get equivalence requests: ${response.statusText}`);
    }
    
    const result = await response.json();
    return result.requests || [];
  }

  async analyzeEquivalence(
    sourceSubjectId: string,
    targetSubjectId: string
  ): Promise<{ similarity: number; compatible: boolean }> {
    try {
      // Call equivalence analysis contract
      const response = await fetch(`${this.baseUrl}/cosmwasm/wasm/v1/contract/cosmos1equiv.../smart/...`);
      
      if (!response.ok) {
        throw new Error(`Failed to analyze equivalence: ${response.statusText}`);
      }
      
      const result = await response.json();
      return {
        similarity: result.data.similarity || 0.85,
        compatible: result.data.compatible || true
      };
    } catch (error) {
      console.error('Equivalence analysis failed:', error);
      throw new Error('Failed to analyze equivalence via smart contract');
    }
  }
}

export class DegreeService {
  private baseUrl: string;
  private txService: TransactionService;
  
  constructor(baseUrl: string = getApiUrl()) {
    this.baseUrl = baseUrl;
    this.txService = new TransactionService(baseUrl);
  }

  async checkDegreeEligibility(studentId: string): Promise<{
    eligible: boolean;
    completedCredits: number;
    requiredCredits: number;
    missingSubjects: string[];
  }> {
    try {
      // Call degree eligibility contract
      const response = await fetch(`${this.baseUrl}/cosmwasm/wasm/v1/contract/cosmos1degree.../smart/...`);
      
      if (!response.ok) {
        throw new Error(`Failed to check degree eligibility: ${response.statusText}`);
      }
      
      const result = await response.json();
      return result.data;
    } catch (error) {
      console.error('Degree eligibility check failed:', error);
      throw new Error('Failed to check degree eligibility via smart contract');
    }
  }

  async requestDegree(studentId: string): Promise<AcademicNFT> {
    try {
      const msg = {
        student_id: studentId
      };

      const result = await this.txService.submitTransaction(
        msg, 
        "/academictoken.degree.MsgRequestDegree"
      );

      return {
        id: `degree-${Date.now()}`,
        studentId: studentId,
        subjectId: 'DEGREE',
        grade: 0, // Degrees don't have grades
        completionDate: new Date().toISOString(),
        nftHash: result.tx_response?.txhash || `degree-hash-${Date.now()}`,
        metadata: {
          subject: 'Bachelor Degree',
          credits: 240,
          institution: 'UFJF'
        }
      };
    } catch (error) {
      console.error('Failed to request degree:', error);
      throw error;
    }
  }
}

// Factory to create service instances
export class BlockchainServiceFactory {
  private static instance: BlockchainServiceFactory;
  
  private constructor() {}
  
  static getInstance(): BlockchainServiceFactory {
    if (!BlockchainServiceFactory.instance) {
      BlockchainServiceFactory.instance = new BlockchainServiceFactory();
    }
    return BlockchainServiceFactory.instance;
  }
  
  createInstitutionService(): InstitutionService {
    return new InstitutionService();
  }
  
  createCourseService(): CourseService {
    return new CourseService();
  }
  
  createSubjectService(): SubjectService {
    return new SubjectService();
  }
  
  createStudentService(): StudentService {
    return new StudentService();
  }
  
  createEquivalenceService(): EquivalenceService {
    return new EquivalenceService();
  }
  
  createDegreeService(): DegreeService {
    return new DegreeService();
  }
}

// Export singleton
export const blockchainServices = BlockchainServiceFactory.getInstance();