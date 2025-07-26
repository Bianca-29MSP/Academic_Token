// Equivalence API endpoints
export interface EquivalenceAnalysisRequest {
  studentId: string;
  sourceSubjectId: string;
  targetInstitution: string;
  targetSubjectId: string;
}

export interface EquivalenceAnalysisResponse {
  equivalenceId: string;
  sourceSubjectId: string;
  targetSubjectId: string;
  equivalencePercent: string;
  analysisMetadata: string;
  contractAddress: string;
  status: 'pending' | 'approved' | 'rejected' | 'error';
  recommendation: 'approved' | 'rejected';
}
