import React, { useMemo } from 'react';
import { SubjectWithPrerequisites, SemesterWithSubjects } from '../../types/curriculum.types';

interface SubjectCardProps {
  subject: SubjectWithPrerequisites;
  isCompleted: boolean;
  onClick: () => void;
  semesters: SemesterWithSubjects[];
  completedSubjects: string[];
}

const SubjectCard: React.FC<SubjectCardProps> = ({ 
  subject, 
  isCompleted, 
  onClick,
  semesters,
  completedSubjects 
}) => {
  // Check if prerequisites are met
  const prerequisitesMet = useMemo(() => {
    if (!subject.prerequisites || subject.prerequisites.length === 0) {
      return true;
    }
    return subject.prerequisites.every(prereqId => completedSubjects.includes(prereqId));
  }, [subject.prerequisites, completedSubjects]);

  // Get prerequisite subjects details
  const prerequisiteSubjects = useMemo(() => {
    if (!subject.prerequisites) return [];
    
    return subject.prerequisites.map(prereqId => {
      for (const semester of semesters) {
        const found = semester.subjects.find(s => s.subjectId === prereqId);
        if (found) return found;
      }
      return null;
    }).filter(Boolean) as SubjectWithPrerequisites[];
  }, [subject.prerequisites, semesters]);

  const cardClassName = `subject-card ${isCompleted ? 'completed' : ''} ${!prerequisitesMet ? 'locked' : ''}`;

  return (
    <div className={cardClassName} onClick={onClick}>
      <div className="subject-header">
        <h4 className="subject-title">{subject.title}</h4>
        <span className="subject-code">{subject.code}</span>
      </div>
      
      <div className="subject-info">
        <div className="info-row">
          <span className="info-label">Créditos:</span>
          <span className="info-value">{subject.credits}</span>
        </div>
        <div className="info-row">
          <span className="info-label">Carga Horária:</span>
          <span className="info-value">{subject.workloadHours}h</span>
        </div>
      </div>

      {prerequisiteSubjects.length > 0 && (
        <div className="prerequisites">
          <span className="prereq-label">Pré-requisitos:</span>
          <ul className="prereq-list">
            {prerequisiteSubjects.map(prereq => (
              <li 
                key={prereq.subjectId} 
                className={completedSubjects.includes(prereq.subjectId) ? 'completed' : 'pending'}
              >
                {prereq.code} - {prereq.title}
              </li>
            ))}
          </ul>
        </div>
      )}

      {isCompleted && (
        <div className="status-badge completed-badge">
          ✓ Concluída
        </div>
      )}

      {!prerequisitesMet && !isCompleted && (
        <div className="status-badge locked-badge">
          🔒 Bloqueada
        </div>
      )}
    </div>
  );
};

export default SubjectCard;
