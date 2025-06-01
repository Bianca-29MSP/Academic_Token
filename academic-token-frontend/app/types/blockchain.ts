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
  syllabus: string;
  metadata: string;
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
  completedCredits: number;
  requiredCredits: number;
  missingSubjects: string[];
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
