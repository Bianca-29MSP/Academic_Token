// components/curriculum/CurriculumGraph.tsx
"use client"
import React, { useEffect, useRef, useState } from 'react'

interface Subject {
  id: string
  name: string
  code: string
  credits: number
  prerequisites: string[]
  semester: number
  status: 'completed' | 'available' | 'locked' | 'enrolled'
  grade?: number
}

interface CurriculumGraphProps {
  subjects: Subject[]
  onSubjectClick?: (subject: Subject) => void
  studentNFTs?: any[]
}

export function CurriculumGraph({ subjects, onSubjectClick, studentNFTs = [] }: CurriculumGraphProps) {
  const svgRef = useRef<SVGSVGElement>(null)
  const containerRef = useRef<HTMLDivElement>(null)
  const [selectedSubject, setSelectedSubject] = useState<string | null>(null)
  const [dimensions, setDimensions] = useState({ width: 1200, height: 800 })

  // Update dimensions based on container size
  useEffect(() => {
    const updateDimensions = () => {
      if (containerRef.current) {
        const container = containerRef.current
        const containerWidth = container.clientWidth
        const containerHeight = container.clientHeight || Math.max(600, window.innerHeight - 400)
        
        setDimensions({
          width: Math.max(1200, containerWidth), // Ensure enough width for 8 semesters
          height: Math.max(600, containerHeight - 120) // Account for header and footer
        })
      }
    }

    updateDimensions()
    window.addEventListener('resize', updateDimensions)
    
    // Also update when container becomes visible
    const observer = new ResizeObserver(updateDimensions)
    if (containerRef.current) {
      observer.observe(containerRef.current)
    }
    
    return () => {
      window.removeEventListener('resize', updateDimensions)
      observer.disconnect()
    }
  }, [])

  // Calculate subject positions based on semester and dependencies
  const calculatePositions = () => {
    const positions: Record<string, { x: number; y: number }> = {}
    const maxSemesters = 8
    const semesterWidth = Math.max(180, (dimensions.width - 100) / maxSemesters)
    const subjectHeight = 80
    const padding = 20

    // Group subjects by semester
    const semesters = subjects.reduce((acc, subject) => {
      if (!acc[subject.semester]) acc[subject.semester] = []
      acc[subject.semester].push(subject)
      return acc
    }, {} as Record<number, Subject[]>)

    // Position subjects
    Object.entries(semesters).forEach(([semester, semesterSubjects]) => {
      const semesterNum = parseInt(semester)
      const x = semesterNum * semesterWidth + 50
      
      semesterSubjects.forEach((subject, index) => {
        const y = index * (subjectHeight + padding) + 80
        positions[subject.id] = { x, y }
      })
    })

    return positions
  }

  // Get subject status based on completion and prerequisites
  const getSubjectStatus = (subject: Subject): Subject['status'] => {
    // Check if completed (has NFT)
    const hasNFT = studentNFTs.some(nft => nft.subjectId === subject.id)
    if (hasNFT) return 'completed'

    // Check if all prerequisites are met
    const prerequisitesMet = subject.prerequisites.every(prereqId => 
      studentNFTs.some(nft => nft.subjectId === prereqId)
    )

    if (!prerequisitesMet) return 'locked'
    
    // For demo, show one subject as enrolled
    if (subject.id === 'ci1002') return 'enrolled'
    
    return 'available'
  }

  // Enhanced subjects with calculated status
  const enhancedSubjects = subjects.map(subject => ({
    ...subject,
    status: getSubjectStatus(subject)
  }))

  const positions = calculatePositions()

  // Generate connection lines between prerequisites and subjects
  const generateConnections = () => {
    return enhancedSubjects.flatMap(subject => 
      subject.prerequisites.map(prereqId => {
        const fromPos = positions[prereqId]
        const toPos = positions[subject.id]
        
        if (!fromPos || !toPos) return null

        // Create curved path
        const midX = (fromPos.x + toPos.x) / 2
        const midY = (fromPos.y + toPos.y) / 2
        
        return {
          id: `${prereqId}-${subject.id}`,
          from: prereqId,
          to: subject.id,
          path: `M ${fromPos.x + 80} ${fromPos.y + 30} Q ${midX} ${midY - 20} ${toPos.x} ${toPos.y + 30}`,
          fromPos,
          toPos
        }
      })
    ).filter(Boolean)
  }

  const connections = generateConnections()

  // Subject color based on status
  const getSubjectColor = (status: Subject['status']) => {
    switch (status) {
      case 'completed': return { bg: '#dcfce7', border: '#22c55e', text: '#15803d' } // green
      case 'enrolled': return { bg: '#dbeafe', border: '#3b82f6', text: '#1d4ed8' } // blue
      case 'available': return { bg: '#fef3c7', border: '#f59e0b', text: '#d97706' } // yellow
      case 'locked': return { bg: '#f3f4f6', border: '#9ca3af', text: '#6b7280' } // gray
      default: return { bg: '#f3f4f6', border: '#9ca3af', text: '#6b7280' }
    }
  }

  // Status icon
  const getStatusIcon = (status: Subject['status']) => {
    switch (status) {
      case 'completed': return '✅'
      case 'enrolled': return '📚'
      case 'available': return '🔓'
      case 'locked': return '🔒'
      default: return '❓'
    }
  }

  const handleSubjectClick = (subject: Subject) => {
    setSelectedSubject(subject.id)
    onSubjectClick?.(subject)
  }

  return (
    <div ref={containerRef} className="w-full h-full bg-white rounded-lg shadow-lg flex flex-col">
      <div className="flex items-center justify-between p-6 border-b border-gray-200">
        <h2 className="text-2xl font-semibold text-gray-800">🎓 Curriculum Tree (DAG)</h2>
        <div className="flex items-center space-x-4 text-sm">
          <div className="flex items-center space-x-1">
            <div className="w-3 h-3 bg-green-200 border border-green-400 rounded"></div>
            <span>Completed</span>
          </div>
          <div className="flex items-center space-x-1">
            <div className="w-3 h-3 bg-blue-200 border border-blue-400 rounded"></div>
            <span>Enrolled</span>
          </div>
          <div className="flex items-center space-x-1">
            <div className="w-3 h-3 bg-yellow-200 border border-yellow-400 rounded"></div>
            <span>Available</span>
          </div>
          <div className="flex items-center space-x-1">
            <div className="w-3 h-3 bg-gray-200 border border-gray-400 rounded"></div>
            <span>Locked</span>
          </div>
        </div>
      </div>

      <div className="flex-1 overflow-auto">
        <svg
          ref={svgRef}
          width="100%"
          height="100%"
          viewBox={`0 0 ${dimensions.width} ${dimensions.height}`}
          className="border-0"
          preserveAspectRatio="xMidYMid meet"
        >
          {/* Background grid */}
          <defs>
            <pattern id="grid" width="50" height="50" patternUnits="userSpaceOnUse">
              <path d="M 50 0 L 0 0 0 50" fill="none" stroke="#f1f5f9" strokeWidth="1"/>
            </pattern>
          </defs>
          <rect width="100%" height="100%" fill="url(#grid)" />

          {/* Semester labels */}
          {[1, 2, 3, 4, 5, 6, 7, 8].map(semester => {
            const maxSemesters = 8
            const semesterWidth = Math.max(180, (dimensions.width - 100) / maxSemesters)
            return (
              <text
                key={semester}
                x={semester * semesterWidth + 50 + semesterWidth / 2}
                y={30}
                textAnchor="middle"
                className="text-sm font-semibold fill-gray-600"
              >
                Semester {semester}
              </text>
            )
          })}

          {/* Connection lines */}
          {connections.map(connection => (
            <g key={connection.id}>
              <path
                d={connection.path}
                fill="none"
                stroke="#9ca3af"
                strokeWidth="2"
                markerEnd="url(#arrowhead)"
                className="transition-all duration-200"
              />
            </g>
          ))}

          {/* Arrow marker definition */}
          <defs>
            <marker
              id="arrowhead"
              markerWidth="10"
              markerHeight="7"
              refX="9"
              refY="3.5"
              orient="auto"
            >
              <polygon
                points="0 0, 10 3.5, 0 7"
                fill="#9ca3af"
              />
            </marker>
          </defs>

          {/* Subject nodes */}
          {enhancedSubjects.map(subject => {
            const pos = positions[subject.id]
            if (!pos) return null

            const colors = getSubjectColor(subject.status)
            const isSelected = selectedSubject === subject.id

            return (
              <g
                key={subject.id}
                transform={`translate(${pos.x}, ${pos.y})`}
                className="cursor-pointer"
                onClick={() => handleSubjectClick(subject)}
              >
                {/* Subject box */}
                <rect
                  width={Math.min(160, Math.max(140, (dimensions.width - 100) / 8 - 20))}
                  height="60"
                  rx="8"
                  fill={colors.bg}
                  stroke={colors.border}
                  strokeWidth={isSelected ? "3" : "2"}
                  className="transition-all duration-200 hover:shadow-lg"
                />

                {/* Subject code */}
                <text
                  x="8"
                  y="16"
                  className="text-xs font-mono"
                  fill={colors.text}
                >
                  {subject.code}
                </text>

                {/* Status icon */}
                <text
                  x={Math.min(160, Math.max(140, (dimensions.width - 100) / 8 - 20)) - 20}
                  y="16"
                  className="text-sm"
                >
                  {getStatusIcon(subject.status)}
                </text>

                {/* Subject name */}
                <text
                  x="8"
                  y="32"
                  className="text-sm font-semibold"
                  fill={colors.text}
                >
                  {subject.name.length > 18 ? subject.name.substring(0, 15) + '...' : subject.name}
                </text>

                {/* Credits and grade */}
                <text
                  x="8"
                  y="48"
                  className="text-xs"
                  fill={colors.text}
                >
                  {subject.credits} credits
                  {subject.status === 'completed' && subject.grade && (
                    <tspan> • Grade: {subject.grade.toFixed(1)}</tspan>
                  )}
                </text>

                {/* Hover effect */}
                <rect
                  width={Math.min(160, Math.max(140, (dimensions.width - 100) / 8 - 20))}
                  height="60"
                  rx="8"
                  fill="transparent"
                  className="hover:fill-black hover:fill-opacity-5 transition-all duration-200"
                />
              </g>
            )
          })}
        </svg>
      </div>

      {/* Selected subject details */}
      {selectedSubject && (
      <div className="p-4 bg-gray-50 border-t border-gray-200">
      {(() => {
      const subject = enhancedSubjects.find(s => s.id === selectedSubject)
      if (!subject) return null

      return (
      <div>
      <div className="flex items-center justify-between mb-2">
      <h3 className="text-lg font-semibold">{subject.name}</h3>
      <button 
      onClick={() => setSelectedSubject(null)}
      className="text-gray-500 hover:text-gray-700"
      >
      ✕
      </button>
      </div>
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
      <div>
      <span className="font-semibold">Code:</span> {subject.code}
      </div>
      <div>
      <span className="font-semibold">Credits:</span> {subject.credits}
      </div>
      <div>
      <span className="font-semibold">Semester:</span> {subject.semester}
      </div>
      <div>
      <span className="font-semibold">Status:</span> {getStatusIcon(subject.status)} {subject.status}
      </div>
      </div>
      {subject.prerequisites.length > 0 && (
      <div className="mt-2">
      <span className="font-semibold text-sm">Prerequisites: </span>
      <span className="text-sm">
      {subject.prerequisites.map(prereq => {
      const prereqSubject = enhancedSubjects.find(s => s.id === prereq)
      return prereqSubject?.code || prereq
      }).join(', ')}
      </span>
      </div>
      )}
      {subject.status === 'completed' && subject.grade && (
      <div className="mt-2">
      <span className="font-semibold text-sm">Grade: </span>
      <span className="text-sm">{subject.grade.toFixed(1)}</span>
      </div>
      )}
      </div>
      )
      })()} 
      </div>
      )}
    </div>
  )
}
