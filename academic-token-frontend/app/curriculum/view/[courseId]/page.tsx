'use client';

import React, { useState } from 'react';
import { useParams } from 'next/navigation';
import { useCurriculumData } from '../../../hooks/useCurriculumData';
import CurriculumGrid from '../../../components/CurriculumGrid/CurriculumGrid';
import SubjectDetailsModal from '../../../components/SubjectDetailsModal/SubjectDetailsModal';
import { SubjectWithPrerequisites } from '../../../types/curriculum.types';
import './page.css';

export default function CurriculumViewPage() {
  const params = useParams();
  const courseId = params.courseId as string;
  const { curriculumWithSubjects, isLoading, curriculumTree, error, warnings } = useCurriculumData(courseId);
  const [selectedSubject, setSelectedSubject] = useState<SubjectWithPrerequisites | null>(null);
  const [showWarnings, setShowWarnings] = useState(true);
  
  // TODO: Replace with actual student data from context/API
  const completedSubjects: string[] = []; // Mock empty for now

  if (isLoading) {
    return (
      <div className="loading-container">
        <div className="loading">Carregando grade curricular...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="error-container">
        <div className="error">{error}</div>
      </div>
    );
  }

  if (!curriculumTree) {
    return (
      <div className="error-container">
        <div className="error">Grade curricular não encontrada</div>
      </div>
    );
  }

  return (
    <div className="curriculum-view">
      <header className="curriculum-header">
        <h1>Grade Curricular</h1>
        <div className="curriculum-info">
          <span>Versão: {curriculumTree.version}</span>
          <span>Carga Horária Total: {curriculumTree.totalWorkloadHours}h</span>
          <span>Disciplinas Obrigatórias: {curriculumTree.requiredSubjects.length}</span>
          <span>Eletivas Mínimas: {curriculumTree.electiveMin}</span>
        </div>
      </header>

      {/* Warnings Section */}
      {warnings.length > 0 && showWarnings && (
        <div className="warnings-container">
          <div className="warnings-header">
            <h3>⚠️ Avisos sobre a estrutura curricular</h3>
            <button 
              className="close-warnings"
              onClick={() => setShowWarnings(false)}
            >
              ✕
            </button>
          </div>
          <div className="warnings-content">
            <p className="warnings-explanation">
              Os seguintes problemas foram detectados na organização das disciplinas por período:
            </p>
            <ul className="warnings-list">
              {warnings.map((warning, index) => (
                <li key={index}>{warning}</li>
              ))}
            </ul>
            <p className="warnings-note">
              <strong>Nota:</strong> A grade foi automaticamente reorganizada para respeitar os pré-requisitos.
            </p>
          </div>
        </div>
      )}

      <CurriculumGrid
        semesters={curriculumWithSubjects}
        onSubjectClick={setSelectedSubject}
        completedSubjects={completedSubjects}
      />

      {selectedSubject && (
        <SubjectDetailsModal
          subject={selectedSubject}
          onClose={() => setSelectedSubject(null)}
        />
      )}
    </div>
  );
}
