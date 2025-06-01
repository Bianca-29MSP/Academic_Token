// Simple fallback component for missing CurriculumGraph
"use client"

interface Subject {
  id: string
  name: string
  code: string
  credits: number
  completed?: boolean
  grade?: number
  metadata?: string
}

interface SimpleCurriculumDisplayProps {
  subjects: Subject[]
  onSubjectClick?: (subject: Subject) => void
}

export function SimpleCurriculumDisplay({ subjects, onSubjectClick }: SimpleCurriculumDisplayProps) {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      {subjects.map((subject) => (
        <div
          key={subject.id}
          onClick={() => onSubjectClick?.(subject)}
          className={`p-4 rounded-lg border-2 cursor-pointer transition-all duration-200 hover:shadow-lg ${
            subject.completed
              ? 'bg-green-50 border-green-200 hover:border-green-300'
              : 'bg-white border-gray-200 hover:border-blue-300'
          }`}
        >
          <div className="flex items-start justify-between mb-2">
            <div className="flex-1">
              <h3 className="font-semibold text-gray-800">{subject.code}</h3>
              <p className="text-sm text-gray-600">{subject.name}</p>
            </div>
            <div className="ml-2">
              {subject.completed ? (
                <div className="text-green-600">
                  <div className="text-lg">✅</div>
                  {subject.grade && (
                    <div className="text-xs font-semibold">
                      {subject.grade.toFixed(1)}
                    </div>
                  )}
                </div>
              ) : (
                <div className="text-gray-400">
                  <div className="text-lg">⏳</div>
                </div>
              )}
            </div>
          </div>
          
          <div className="flex items-center justify-between text-sm">
            <span className="text-gray-500">{subject.credits} créditos</span>
            {subject.completed && (
              <span className="text-green-600 text-xs">Completed</span>
            )}
          </div>
          
          {subject.metadata && (
            <p className="text-xs text-gray-500 mt-2 line-clamp-2">
              {subject.metadata}
            </p>
          )}
        </div>
      ))}
    </div>
  )
}
