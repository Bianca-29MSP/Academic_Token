export interface Institution {
  index: string;
  address: string;
  name: string;
  code?: string;
  country?: string;
  isAuthorized: string;
  creator: string;
  createdAt?: string;
}

export interface Course {
  index: string;
  institution: string;
  name: string;
  code: string;
  description: string;
  totalCredits: string;
  degreeLevel: string;
  duration?: number;
}

export interface Subject {
  index: string;
  subjectId: string;
  institution: string;
  course_id: string;
  title: string;
  code: string;
  workloadHours: string;
  credits: string;
  description: string;
  contentHash: string;
  subjectType: string;
  knowledgeArea: string;
  ipfsLink: string;
  creator: string;
  syllabus?: string;
  metadata?: string;
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

export interface NFTMetadata {
  subject: string;
  credits: number;
  institution: string;
}

export interface AcademicNFT {
  id: string;
  studentId: string;
  subjectId: string;
  grade: number;
  completionDate: string;
  nftHash: string;
  metadata: NFTMetadata;
}

export interface DegreeEligibility {
  eligible: boolean;
  completedCredits?: number;
  requiredCredits?: number;
  missingSubjects?: string[];
  message?: string;
  graduationStatus?: any; // Will be defined by the graduation_status proto message
}

export interface EquivalenceRequest {
  sourceInstitutionId: string;
  sourceSubjectId: string;
  targetInstitutionId: string;
  targetSubjectId: string;
  similarity?: number;
  status: 'pending' | 'approved' | 'rejected';
  requestDate: string;
}

export interface BlockchainConnection {
  connected: boolean;
  nodeUrl: string;
  network: string;
  version?: string;
}
