// Curriculum types for grid visualization
export interface CurriculumSemester {
  semesterNumber: string;
  subjectIds: string[];
}

export interface CurriculumTree {
  index: string;
  courseId: string;
  version: string;
  totalWorkloadHours: string;
  requiredSubjects: string[];
  electiveMin: number;
  electiveSubjects: string[];
  semesterStructure: CurriculumSemester[];
  graduationRequirements: any;
  electiveGroups: any[];
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
  ipfsLink: string;
  creator: string;
}

export interface SubjectWithPrerequisites extends SubjectContent {
  prerequisites?: string[]; // Subject IDs
}

export interface SemesterWithSubjects {
  semesterNumber: string;
  subjects: SubjectWithPrerequisites[];
}
