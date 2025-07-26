import React from 'react';
import { SemesterWithSubjects, SubjectWithPrerequisites } from '../../types/curriculum.types';
import './CurriculumComparison.css';

interface CurriculumComparisonProps {
  originalSemesters: SemesterWithSubjects[];
  reorganizedSemesters: SemesterWithSubjects[];
  onSubjectClick?: (subject: SubjectWithPrerequisites) => void;
}

const CurriculumComparison: React.FC<CurriculumComparisonProps> = ({
  originalSemesters,
  reorganizedSemesters,
  onSubjectClick
}) => {
  // Create a map to track which subjects moved
  const subjectMovements = new Map<string, { from: number; to: number }>();
  
  originalSemesters.forEach(semester => {
    const fromSemester = parseInt(semester.semesterNumber);
    semester.subjects.forEach(subject => {
      const toSemester = reorganizedSemesters.find(s => 
        s.subjects.some(subj => subj.subjectId === subject.subjectId)
      );
      if (toSemester && parseInt(toSemester.semesterNumber) !== fromSemester) {
        subjectMovements.set(subject.subjectId, {
          from: fromSemester,
          to: parseInt(toSemester.semesterNumber)
        });
      }
    });
  });

  const renderSubject = (subject: SubjectWithPrerequisites, isMoved: boolean) => (
    <div
      key={subject.subjectId}
      className={`comparison-subject ${isMoved ? 'moved' : ''}`}
      onClick={() => onSubjectClick?.(subject)}
    >
      <div className="subject-header">
        <span className="subject-code">{subject.code}</span>
        {isMoved && <span className="moved-badge">Movida</span>}
      </div>
      <div className="subject-title">{subject.title}</div>
      <div className="subject-details">
        <span>{subject.credits} créditos</span>
        <span>{subject.workloadHours}h</span>
      </div>
      {subject.prerequisites && subject.prerequisites.length > 0 && (
        <div className="subject-prereqs">
          Pré-req: {subject.prerequisites.length}
        </div>
      )}
    </div>
  );

  return (
    <div className="curriculum-comparison">
      <div className="comparison-side original">
        <h3>Grade Original</h3>
        <div className="semesters-list">
          {originalSemesters.map(semester => (
            <div key={semester.semesterNumber} className="semester-section">
              <h4>{semester.semesterNumber}º Período</h4>
              <div className="subjects-list">
                {semester.subjects.map(subject => 
                  renderSubject(subject, subjectMovements.has(subject.subjectId))
                )}
              </div>
            </div>
          ))}
        </div>
      </div>

      <div className="comparison-divider">
        <span>vs</span>
      </div>

      <div className="comparison-side reorganized">
        <h3>Grade Reorganizada (Respeitando Pré-requisitos)</h3>
        <div className="semesters-list">
          {reorganizedSemesters.map(semester => (
            <div key={semester.semesterNumber} className="semester-section">
              <h4>{semester.semesterNumber}º Período</h4>
              <div className="subjects-list">
                {semester.subjects.map(subject => {
                  const movement = subjectMovements.get(subject.subjectId);
                  const movedFrom = movement ? movement.from : null;
                  
                  return (
                    <div key={subject.subjectId} className="subject-wrapper">
                      {renderSubject(subject, !!movement)}
                      {movement && (
                        <div className="movement-info">
                          ← Veio do {movement.from}º período
                        </div>
                      )}
                    </div>
                  );
                })}
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default CurriculumComparison;
