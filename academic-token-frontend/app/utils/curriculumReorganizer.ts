import { 
  CurriculumTree, 
  SubjectContent, 
  SemesterWithSubjects, 
  SubjectWithPrerequisites,
  CurriculumSemester
} from '../types/curriculum.types';

interface PrerequisiteGroup {
  id: string;
  subjectId: string;
  groupType: 'ALL' | 'ANY' | 'CREDITS';
  minimumCredits: number;
  minimumCompletedSubjects: number;
  subjectIds: string[];
}

/**
 * Reorganizes subjects by semester respecting prerequisites
 * This ensures that no subject is placed in a semester before its prerequisites
 */
export function reorganizeByPrerequisites(
  curriculumTree: CurriculumTree,
  subjects: SubjectContent[],
  prerequisites: Record<string, PrerequisiteGroup[]>
): SemesterWithSubjects[] {
  // Create a map for quick subject lookup
  const subjectMap = new Map<string, SubjectContent>();
  subjects.forEach(subject => {
    subjectMap.set(subject.subjectId, subject);
  });

  // Build dependency graph
  const dependencies = new Map<string, Set<string>>();
  const dependents = new Map<string, Set<string>>();
  
  // Initialize all subjects in the graph
  subjects.forEach(subject => {
    dependencies.set(subject.subjectId, new Set());
    dependents.set(subject.subjectId, new Set());
  });

  // Build the dependency relationships
  Object.entries(prerequisites).forEach(([subjectId, prereqGroups]) => {
    prereqGroups.forEach(group => {
      if (group.groupType === 'ALL') {
        // For ALL type, subject depends on all listed subjects
        group.subjectIds.forEach(prereqId => {
          dependencies.get(subjectId)?.add(prereqId);
          dependents.get(prereqId)?.add(subjectId);
        });
      } else if (group.groupType === 'ANY') {
        // For ANY type, we'll need more complex handling
        // For now, we'll treat it as optional (no hard dependency)
        // In a real implementation, you might want to choose the earliest available
      }
    });
  });

  // Calculate the minimum semester for each subject using topological sort
  const subjectSemester = new Map<string, number>();
  const visited = new Set<string>();
  const visiting = new Set<string>();

  function calculateMinSemester(subjectId: string): number {
    if (visited.has(subjectId)) {
      return subjectSemester.get(subjectId) || 1;
    }

    if (visiting.has(subjectId)) {
      console.warn(`Circular dependency detected involving ${subjectId}`);
      return 1;
    }

    visiting.add(subjectId);

    const deps = dependencies.get(subjectId) || new Set();
    let maxPrereqSemester = 0;

    deps.forEach(depId => {
      const depSemester = calculateMinSemester(depId);
      maxPrereqSemester = Math.max(maxPrereqSemester, depSemester);
    });

    visiting.delete(subjectId);
    visited.add(subjectId);

    // Subject must be at least one semester after its prerequisites
    const minSemester = deps.size > 0 ? maxPrereqSemester + 1 : 1;
    subjectSemester.set(subjectId, minSemester);

    return minSemester;
  }

  // Calculate minimum semester for all subjects
  subjects.forEach(subject => {
    calculateMinSemester(subject.subjectId);
  });

  // Now reorganize the curriculum respecting both original structure and prerequisites
  const reorganizedSemesters: SemesterWithSubjects[] = [];
  const placedSubjects = new Set<string>();

  // First, try to respect the original curriculum structure
  curriculumTree.semesterStructure.forEach((semester, index) => {
    const semesterNumber = parseInt(semester.semesterNumber);
    const validSubjects: SubjectWithPrerequisites[] = [];

    semester.subjectIds.forEach(subjectId => {
      const subject = subjectMap.get(subjectId);
      if (!subject) return;

      const minSemester = subjectSemester.get(subjectId) || 1;
      
      // Only place subject if it can be taken in this semester
      if (minSemester <= semesterNumber) {
        validSubjects.push({
          ...subject,
          prerequisites: Array.from(dependencies.get(subjectId) || [])
        });
        placedSubjects.add(subjectId);
      }
    });

    reorganizedSemesters.push({
      semesterNumber: semester.semesterNumber,
      subjects: validSubjects
    });
  });

  // Handle subjects that couldn't be placed in their original semesters
  const unplacedSubjects = subjects.filter(s => !placedSubjects.has(s.subjectId));
  
  unplacedSubjects.forEach(subject => {
    const minSemester = subjectSemester.get(subject.subjectId) || 1;
    
    // Find or create the appropriate semester
    let targetSemester = reorganizedSemesters.find(s => parseInt(s.semesterNumber) === minSemester);
    
    if (!targetSemester) {
      // Create new semester if needed
      targetSemester = {
        semesterNumber: minSemester.toString(),
        subjects: []
      };
      reorganizedSemesters.push(targetSemester);
    }

    targetSemester.subjects.push({
      ...subject,
      prerequisites: Array.from(dependencies.get(subject.subjectId) || [])
    });
  });

  // Sort semesters by number
  reorganizedSemesters.sort((a, b) => 
    parseInt(a.semesterNumber) - parseInt(b.semesterNumber)
  );

  return reorganizedSemesters;
}

/**
 * Creates a visual warning for curriculum issues
 */
export function validateCurriculumStructure(
  semesters: SemesterWithSubjects[]
): string[] {
  const warnings: string[] = [];
  const subjectSemesterMap = new Map<string, number>();

  // Build a map of when each subject is offered
  semesters.forEach(semester => {
    const semesterNum = parseInt(semester.semesterNumber);
    semester.subjects.forEach(subject => {
      subjectSemesterMap.set(subject.subjectId, semesterNum);
    });
  });

  // Check for prerequisite violations
  semesters.forEach(semester => {
    semester.subjects.forEach(subject => {
      if (subject.prerequisites && subject.prerequisites.length > 0) {
        subject.prerequisites.forEach(prereqId => {
          const prereqSemester = subjectSemesterMap.get(prereqId);
          const subjectSemester = subjectSemesterMap.get(subject.subjectId);

          if (prereqSemester && subjectSemester && prereqSemester >= subjectSemester) {
            warnings.push(
              `⚠️ ${subject.code} (${subject.title}) no período ${subjectSemester} ` +
              `tem como pré-requisito uma disciplina do período ${prereqSemester} ou posterior`
            );
          }
        });
      }
    });
  });

  return warnings;
}
