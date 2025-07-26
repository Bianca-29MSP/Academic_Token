import React from 'react';
import { SubjectWithPrerequisites } from '../../types/curriculum.types';
import './SubjectDetailsModal.css';

interface SubjectDetailsModalProps {
  subject: SubjectWithPrerequisites;
  onClose: () => void;
}

const SubjectDetailsModal: React.FC<SubjectDetailsModalProps> = ({ subject, onClose }) => {
  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <button className="close-button" onClick={onClose}>×</button>
        
        <h2>{subject.title}</h2>
        <p className="subject-code-modal">{subject.code}</p>
        
        <div className="details-section">
          <h3>Informações Gerais</h3>
          <div className="detail-row">
            <span>Créditos:</span>
            <span>{subject.credits}</span>
          </div>
          <div className="detail-row">
            <span>Carga Horária:</span>
            <span>{subject.workloadHours}h</span>
          </div>
          <div className="detail-row">
            <span>Tipo:</span>
            <span>{subject.subjectType}</span>
          </div>
          <div className="detail-row">
            <span>Área de Conhecimento:</span>
            <span>{subject.knowledgeArea}</span>
          </div>
        </div>

        <div className="details-section">
          <h3>Descrição</h3>
          <p>{subject.description}</p>
        </div>

        {subject.prerequisites && subject.prerequisites.length > 0 && (
          <div className="details-section">
            <h3>Pré-requisitos</h3>
            <ul>
              {subject.prerequisites.map(prereq => (
                <li key={prereq}>{prereq}</li>
              ))}
            </ul>
          </div>
        )}

        <div className="modal-actions">
          <button className="primary-button" onClick={() => {
            // Navigate to subject details or enroll
            console.log('View full subject details:', subject.subjectId);
          }}>
            Ver Detalhes Completos
          </button>
        </div>
      </div>
    </div>
  );
};

export default SubjectDetailsModal;
