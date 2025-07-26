import { useState, useEffect } from 'react';
import { 
  CurriculumTree, 
  SubjectContent, 
  SemesterWithSubjects, 
  SubjectWithPrerequisites 
} from '../types/curriculum.types';
import { reorganizeByPrerequisites, validateCurriculumStructure } from '../utils/curriculumReorganizer';
import axios from 'axios';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

interface PrerequisiteGroup {
  id: string;
  subjectId: string;
  groupType: 'ALL' | 'ANY' | 'CREDITS';
  minimumCredits: number;
  minimumCompletedSubjects: number;
  subjectIds: string[];
}

export const useCurriculumData = (courseId: string) => {
  const [curriculumWithSubjects, setCurriculumWithSubjects] = useState<SemesterWithSubjects[]>([]);
  const [curriculumTree, setCurriculumTree] = useState<CurriculumTree | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [warnings, setWarnings] = useState<string[]>([]);

  useEffect(() => {
    if (!courseId) return;

    const fetchData = async () => {
      setIsLoading(true);
      setError(null);

      try {
        // Fetch curriculum tree
        const curriculumResponse = await axios.get(`${API_URL}/curriculum/tree/${courseId}`);
        const curriculumData = curriculumResponse.data as CurriculumTree;
        setCurriculumTree(curriculumData);

        // Fetch all subjects for the course
        const subjectsResponse = await axios.get(`${API_URL}/subjects/course/${courseId}`);
        const allSubjects = subjectsResponse.data as SubjectContent[];

        // Try to fetch prerequisites
        let prerequisites: Record<string, PrerequisiteGroup[]> = {};
        try {
          const prerequisitesResponse = await axios.get(`${API_URL}/subject/prerequisites/course/${courseId}`);
          prerequisites = prerequisitesResponse.data;
        } catch (prereqError) {
          console.warn('Prerequisites endpoint not available, using empty prerequisites');
        }

        // Check if we should reorganize
        const hasPrerequisites = Object.keys(prerequisites).length > 0;
        
        if (hasPrerequisites) {
          // Reorganize curriculum respecting prerequisites
          const reorganizedSemesters = reorganizeByPrerequisites(
            curriculumData,
            allSubjects,
            prerequisites
          );
          
          // Validate and get warnings
          const structureWarnings = validateCurriculumStructure(reorganizedSemesters);
          setWarnings(structureWarnings);
          
          setCurriculumWithSubjects(reorganizedSemesters);
        } else {
          // Use original structure without prerequisites
          const semestersWithSubjects: SemesterWithSubjects[] = curriculumData.semesterStructure.map(semester => {
            const subjectsInSemester = semester.subjectIds
              .map(subjectId => {
                const subject = allSubjects.find(s => s.subjectId === subjectId);
                if (subject) {
                  return {
                    ...subject,
                    prerequisites: [],
                  } as SubjectWithPrerequisites;
                }
                return null;
              })
              .filter(Boolean) as SubjectWithPrerequisites[];

            return {
              semesterNumber: semester.semesterNumber,
              subjects: subjectsInSemester,
            };
          });

          setCurriculumWithSubjects(semestersWithSubjects);
        }
      } catch (err) {
        console.error('Error fetching curriculum data:', err);
        setError('Failed to load curriculum data');
      } finally {
        setIsLoading(false);
      }
    };

    fetchData();
  }, [courseId]);

  return {
    curriculumWithSubjects,
    isLoading,
    curriculumTree,
    error,
    warnings,
  };
};
