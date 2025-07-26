import React from 'react';
import { SemesterWithSubjects, SubjectWithPrerequisites } from '../../types/curriculum.types';
import SubjectCard from './SubjectCard';
import './CurriculumGrid.css';

interface CurriculumGridProps {
  semesters: SemesterWithSubjects[];
  onSubjectClick?: (subject: SubjectWithPrerequisites) => void;
  completedSubjects?: string[]; // Array of completed subject IDs
}

const CurriculumGrid: React.FC<CurriculumGridProps> = ({ 
  semesters, 
  onSubjectClick,
  completedSubjects = [] 
}) => {
  return (
    <div className="curriculum-grid">
      {semesters.map((semester) => (
        <div key={semester.semesterNumber} className="semester-column">
          <h3 className="semester-title">
            {semester.semesterNumber}º Período
          </h3>
          <div className="subjects-container">
            {semester.subjects.map((subject) => (
              <SubjectCard
                key={subject.subjectId}
                subject={subject}
                isCompleted={completedSubjects.includes(subject.subjectId)}
                onClick={() => onSubjectClick?.(subject)}
                semesters={semesters}
                completedSubjects={completedSubjects}
              />
            ))}
          </div>
        </div>
      ))}
    </div>
  );
};

export default CurriculumGrid;
