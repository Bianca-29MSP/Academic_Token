"use client"
import { useState, useEffect, useMemo } from 'react'
import Link from 'next/link'
import { useBlockchain } from '../hooks/useBlockchain'
import { CurriculumGraph } from '../components/curriculum/CurriculumGraph'

interface EnhancedSubject {
  id: string
  name: string
  code: string
  credits: number
  prerequisites: string[]
  semester: number
  status: 'completed' | 'available' | 'locked' | 'enrolled'
  grade?: number
}

export default function CurriculumPage() {
  const [viewMode, setViewMode] = useState<'graph' | 'list'>('graph')
  const [selectedCourse, setSelectedCourse] = useState<string>('')
  const [selectedSubject, setSelectedSubject] = useState<any>(null)
  
  const { 
    connection, 
    isLoading, 
    error,
    institutions,
    courses,
    subjects,
    getCoursesByInstitution,
    getSubjectsByCourse
  } = useBlockchain()

  // Mock student NFTs for demonstration
  const studentNFTs = [
    { subjectId: 'ci1001', grade: 8.5 },
    { subjectId: 'ma1001', grade: 9.0 },
    { subjectId: 'ci1003', grade: 7.8 }
  ]

  // Get courses for selected institution (if any)
  const availableCourses = selectedCourse ? getCoursesByInstitution(selectedCourse) : courses

  // Transform subjects into enhanced format with prerequisites and semester info
  const enhancedSubjects: EnhancedSubject[] = useMemo(() => {
    return subjects.map((subject, index) => {
      // Extract prerequisites from metadata if available
      let prerequisites: string[] = []
      let semester = 1
      
      // Try to parse metadata for prerequisites
      if (subject.metadata) {
        try {
          const metadata = JSON.parse(subject.metadata)
          prerequisites = metadata.prerequisites || []
          semester = metadata.semester || Math.floor(index / 6) + 1 // Default: 6 subjects per semester
        } catch (e) {
          // If metadata is not JSON, try to extract semester from subject code
          const semesterMatch = subject.code.match(/\d/)
          if (semesterMatch) {
            semester = parseInt(semesterMatch[0])
          }
        }
      }

      // Determine subject status
      const hasNFT = studentNFTs.some(nft => nft.subjectId === subject.subjectId)
      let status: EnhancedSubject['status'] = 'available'
      
      if (hasNFT) {
        status = 'completed'
      } else if (prerequisites.length > 0) {
        // Check if all prerequisites are met
        const allPrereqsMet = prerequisites.every(prereqId => 
          studentNFTs.some(nft => nft.subjectId === prereqId)
        )
        status = allPrereqsMet ? 'available' : 'locked'
      }

      // Mock one subject as enrolled
      if (subject.subjectId === 'ci1002') {
        status = 'enrolled'
      }

      const nft = studentNFTs.find(n => n.subjectId === subject.subjectId)

      return {
        id: subject.subjectId,
        name: subject.title || subject.name,
        code: subject.code,
        credits: parseInt(subject.credits) || 4,
        prerequisites,
        semester,
        status,
        grade: nft?.grade
      }
    })
  }, [subjects, studentNFTs])

  // Create sample prerequisites relationships if none exist
  const subjectsWithPrereqs = useMemo(() => {
    if (enhancedSubjects.some(s => s.prerequisites.length > 0)) {
      return enhancedSubjects
    }

    // Create sample prerequisite relationships for demonstration
    return enhancedSubjects.map(subject => {
      const newSubject = { ...subject }
      
      // Add some logical prerequisites based on subject codes
      if (subject.code.includes('2')) {
        // Level 2 subjects require corresponding level 1
        const baseCode = subject.code.replace('2', '1')
        const prereq = enhancedSubjects.find(s => s.code === baseCode)
        if (prereq) {
          newSubject.prerequisites = [prereq.id]
        }
      } else if (subject.code.includes('3')) {
        // Level 3 subjects require level 2
        const baseCode = subject.code.replace('3', '2')
        const prereq = enhancedSubjects.find(s => s.code === baseCode)
        if (prereq) {
          newSubject.prerequisites = [prereq.id]
        }
      }
      
      // Add cross-discipline prerequisites for advanced subjects
      if (subject.semester >= 5) {
        const foundationSubjects = enhancedSubjects
          .filter(s => s.semester <= 2 && s.id !== subject.id)
          .slice(0, 2)
          .map(s => s.id)
        
        newSubject.prerequisites = [...new Set([...newSubject.prerequisites, ...foundationSubjects])]
      }
      
      return newSubject
    })
  }, [enhancedSubjects])

  const handleSubjectClick = (subject: any) => {
    setSelectedSubject(subject)
    console.log('Subject clicked:', subject)
  }

  // Group subjects by semester for list view
  const subjectsBySemester = useMemo(() => {
    const grouped: Record<number, EnhancedSubject[]> = {}
    subjectsWithPrereqs.forEach(subject => {
      if (!grouped[subject.semester]) {
        grouped[subject.semester] = []
      }
      grouped[subject.semester].push(subject)
    })
    return grouped
  }, [subjectsWithPrereqs])

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-gradient-to-r from-blue-600 to-purple-600 text-white p-4">
        <div className="container mx-auto">
          <div className="flex items-center justify-between mb-4">
            <div>
              <h1 className="text-3xl font-bold">📚 Curriculum Structure</h1>
              <p className="text-blue-100">Visualize your academic journey</p>
            </div>
            
            <div className={`text-sm px-3 py-1 rounded-full ${
              connection.connected 
                ? 'bg-green-500/20 text-green-200' 
                : 'bg-red-500/20 text-red-200'
            }`}>
              {connection.connected ? '🔗 Connected' : '❌ Offline'}
            </div>
          </div>

          {/* Navigation */}
          <nav className="flex flex-wrap gap-2">
            <Link href="/" className="bg-white/20 hover:bg-white/30 px-4 py-2 rounded-lg transition-colors">
              🎓 Dashboard
            </Link>
            <Link href="/institution" className="bg-white/20 hover:bg-white/30 px-4 py-2 rounded-lg transition-colors">
              🏛️ Institution
            </Link>
            <Link href="/student" className="bg-white/20 hover:bg-white/30 px-4 py-2 rounded-lg transition-colors">
              👨‍🎓 Student
            </Link>
            <Link href="/curriculum" className="bg-white text-blue-600 px-4 py-2 rounded-lg font-semibold">
              📚 Curriculum
            </Link>
            <Link href="/equivalences" className="bg-white/20 hover:bg-white/30 px-4 py-2 rounded-lg transition-colors">
              🔄 Equivalences
            </Link>
            <Link href="/degree" className="bg-white/20 hover:bg-white/30 px-4 py-2 rounded-lg transition-colors">
              📜 Degrees
            </Link>
          </nav>
        </div>
      </header>

      <div className="container mx-auto p-6">
        {/* Error Alert */}
        {error && (
          <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg">
            <p className="text-red-700 text-sm">⚠️ Error: {error}</p>
          </div>
        )}

        {/* View Mode Toggle */}
        <div className="bg-white rounded-lg shadow-lg p-4 mb-6">
          <div className="flex items-center justify-between">
            <h2 className="text-xl font-semibold text-gray-800">View Mode</h2>
            <div className="flex space-x-2">
              {/* Grid View Button - For testing */}
              <Link
                href="/curriculum/view/course1"
                className="px-4 py-2 rounded-lg bg-purple-500 text-white hover:bg-purple-600 transition-colors"
              >
                🎯 Grid View (New)
              </Link>
              <button
                onClick={() => setViewMode('graph')}
                className={`px-4 py-2 rounded-lg transition-colors ${
                  viewMode === 'graph'
                    ? 'bg-blue-500 text-white'
                    : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                }`}
              >
                🕸️ Prerequisite Graph
              </button>
              <button
                onClick={() => setViewMode('list')}
                className={`px-4 py-2 rounded-lg transition-colors ${
                  viewMode === 'list'
                    ? 'bg-blue-500 text-white'
                    : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                }`}
              >
                📋 Semester List
              </button>
            </div>
          </div>
        </div>

        {/* Main Content */}
        {isLoading ? (
          <div className="flex items-center justify-center h-96 bg-white rounded-lg shadow-lg">
            <div className="text-center">
              <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
              <p className="text-gray-600">Loading curriculum data...</p>
            </div>
          </div>
        ) : subjects.length === 0 ? (
          <div className="bg-white rounded-lg shadow-lg p-12 text-center">
            <div className="text-gray-400 text-6xl mb-4">📚</div>
            <h3 className="text-xl font-semibold text-gray-600 mb-2">No curriculum data available</h3>
            <p className="text-gray-500">
              {connection.connected 
                ? "No subjects have been registered yet. Go to Institution page to add subjects."
                : "Connect to blockchain to see curriculum data."
              }
            </p>
          </div>
        ) : (
          <>
            {viewMode === 'graph' ? (
              <div className="bg-white rounded-lg shadow-lg" style={{ height: '70vh' }}>
                <CurriculumGraph
                  subjects={subjectsWithPrereqs}
                  onSubjectClick={handleSubjectClick}
                  studentNFTs={studentNFTs}
                />
              </div>
            ) : (
              <div className="bg-white rounded-lg shadow-lg p-6">
                <h2 className="text-2xl font-semibold text-gray-800 mb-6">📅 Curriculum by Semester</h2>
                
                <div className="space-y-8">
                  {Object.entries(subjectsBySemester)
                    .sort(([a], [b]) => parseInt(a) - parseInt(b))
                    .map(([semester, semesterSubjects]) => (
                      <div key={semester} className="border-l-4 border-blue-500 pl-6">
                        <h3 className="text-lg font-semibold text-gray-800 mb-4">
                          Semester {semester}
                        </h3>
                        
                        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                          {semesterSubjects.map(subject => (
                            <div
                              key={subject.id}
                              onClick={() => handleSubjectClick(subject)}
                              className={`p-4 rounded-lg border-2 cursor-pointer transition-all duration-200 hover:shadow-lg ${
                                subject.status === 'completed'
                                  ? 'bg-green-50 border-green-200'
                                  : subject.status === 'enrolled'
                                  ? 'bg-blue-50 border-blue-200'
                                  : subject.status === 'available'
                                  ? 'bg-yellow-50 border-yellow-200'
                                  : 'bg-gray-50 border-gray-200'
                              }`}
                            >
                              <div className="flex items-center justify-between mb-2">
                                <span className="text-xs font-mono bg-white px-2 py-1 rounded">
                                  {subject.code}
                                </span>
                                <span className="text-lg">
                                  {subject.status === 'completed' ? '✅' :
                                   subject.status === 'enrolled' ? '📚' :
                                   subject.status === 'available' ? '🔓' : '🔒'}
                                </span>
                              </div>
                              
                              <h4 className="font-semibold text-gray-800 mb-1">
                                {subject.name}
                              </h4>
                              
                              <div className="text-sm text-gray-600">
                                <div>📊 {subject.credits} credits</div>
                                {subject.grade && (
                                  <div>✨ Grade: {subject.grade.toFixed(1)}</div>
                                )}
                                {subject.prerequisites.length > 0 && (
                                  <div className="mt-1 text-xs">
                                    ⚡ Prerequisites: {subject.prerequisites.length}
                                  </div>
                                )}
                              </div>
                            </div>
                          ))}
                        </div>
                      </div>
                    ))}
                </div>
              </div>
            )}
          </>
        )}

        {/* Selected Subject Details */}
        {selectedSubject && (
          <div className="mt-6 bg-white rounded-lg shadow-lg p-6">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-xl font-semibold text-gray-800">
                {selectedSubject.name}
              </h3>
              <button
                onClick={() => setSelectedSubject(null)}
                className="text-gray-500 hover:text-gray-700"
              >
                ✕
              </button>
            </div>
            
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <p className="text-sm text-gray-600">
                  <strong>Code:</strong> {selectedSubject.code}
                </p>
                <p className="text-sm text-gray-600">
                  <strong>Credits:</strong> {selectedSubject.credits}
                </p>
                <p className="text-sm text-gray-600">
                  <strong>Semester:</strong> {selectedSubject.semester}
                </p>
                <p className="text-sm text-gray-600">
                  <strong>Status:</strong> {selectedSubject.status}
                </p>
                {selectedSubject.grade && (
                  <p className="text-sm text-gray-600">
                    <strong>Grade:</strong> {selectedSubject.grade.toFixed(1)}
                  </p>
                )}
              </div>
              
              <div>
                {selectedSubject.prerequisites.length > 0 && (
                  <div>
                    <p className="text-sm font-semibold text-gray-700 mb-2">
                      Prerequisites:
                    </p>
                    <ul className="list-disc list-inside text-sm text-gray-600">
                      {selectedSubject.prerequisites.map((prereqId: string) => {
                        const prereq = subjectsWithPrereqs.find(s => s.id === prereqId)
                        return (
                          <li key={prereqId}>
                            {prereq ? `${prereq.code} - ${prereq.name}` : prereqId}
                            {studentNFTs.some(nft => nft.subjectId === prereqId) && ' ✅'}
                          </li>
                        )
                      })}
                    </ul>
                  </div>
                )}
              </div>
            </div>
            
            {selectedSubject.status === 'available' && (
              <button className="mt-4 bg-blue-500 text-white px-4 py-2 rounded-lg hover:bg-blue-600 transition-colors">
                📝 Request Enrollment
              </button>
            )}
          </div>
        )}
      </div>
    </div>
  )
}
