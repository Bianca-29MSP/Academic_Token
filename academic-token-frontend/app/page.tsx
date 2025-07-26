"use client"
import Link from 'next/link'
import { useBlockchain } from './hooks/useBlockchain'

export default function AcademicTokenDashboard() {
  const { 
    connection, 
    isLoading, 
    error, 
    institutions, 
    subjects 
  } = useBlockchain()

  // Real student data - should be loaded from blockchain/authentication
  const currentStudent = {
    id: "demo_student",
    name: "Demo Student",
    course: "Computer Science",
    institutionId: ""
  }

  // Mock student NFTs for demonstration
  const studentNFTs = [
    { subjectId: 'calc1', grade: 8.5 },
    { subjectId: 'prog1', grade: 9.0 },
  ]

  // Transform subjects with status and prerequisites
  const enhancedSubjects = subjects.map((subject, index) => {
    const hasNFT = studentNFTs.some(nft => nft.subjectId === subject.subjectId)
    
    // Extract prerequisites from metadata if available
    let prerequisites: string[] = []
    let semester = Math.floor(index / 5) + 1 // Default: 5 subjects per semester
    
    if (subject.metadata) {
      try {
        const metadata = JSON.parse(subject.metadata)
        prerequisites = metadata.prerequisites || []
        semester = metadata.semester || semester
      } catch (e) {
        // If metadata is not JSON, continue with defaults
      }
    }

    // Determine status based on NFT and prerequisites
    let status = 'available'
    if (hasNFT) {
      status = 'completed'
    } else if (prerequisites.length > 0) {
      const allPrereqsMet = prerequisites.every(prereqId => 
        studentNFTs.some(nft => nft.subjectId === prereqId)
      )
      status = allPrereqsMet ? 'available' : 'blocked'
    }

    // Mock one subject as enrolled
    if (subject.subjectId === 'calc2') {
      status = 'enrolled'
    }

    const nft = studentNFTs.find(n => n.subjectId === subject.subjectId)

    return {
      ...subject,
      prerequisites,
      semester,
      status,
      grade: nft?.grade
    }
  })

  // Group subjects by semester
  const subjectsBySemester = enhancedSubjects.reduce((acc, subject) => {
    if (!acc[subject.semester]) {
      acc[subject.semester] = []
    }
    acc[subject.semester].push(subject)
    return acc
  }, {} as Record<number, typeof enhancedSubjects>)

  // Calculate statistics
  const stats = {
    completed: enhancedSubjects.filter(s => s.status === 'completed').length,
    enrolled: enhancedSubjects.filter(s => s.status === 'enrolled').length,
    available: enhancedSubjects.filter(s => s.status === 'available').length,
    blocked: enhancedSubjects.filter(s => s.status === 'blocked').length,
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'completed':
        return 'bg-green-100 border-green-300'
      case 'enrolled':
        return 'bg-blue-100 border-blue-300'
      case 'available':
        return 'bg-yellow-100 border-yellow-300'
      case 'blocked':
        return 'bg-gray-100 border-gray-300'
      default:
        return 'bg-gray-100 border-gray-300'
    }
  }

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'completed':
        return '✅'
      case 'enrolled':
        return '📚'
      case 'available':
        return '🏆'
      case 'blocked':
        return '🔒'
      default:
        return '❓'
    }
  }

  const getStatusText = (status: string) => {
    switch (status) {
      case 'completed':
        return 'Completed'
      case 'enrolled':
        return 'Enrolled'
      case 'available':
        return 'Available'
      case 'blocked':
        return 'Blocked'
      default:
        return 'Unknown'
    }
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <div className="bg-gradient-to-r from-blue-600 via-blue-500 to-purple-600 text-white">
        <div className="container mx-auto px-4 py-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-3">
              <div className="w-8 h-8 bg-white/20 rounded-lg flex items-center justify-center">
                🎓
              </div>
              <div>
                <h1 className="text-xl font-bold">AcademicToken</h1>
                <p className="text-sm text-blue-100">{currentStudent.name} - {currentStudent.course}</p>
              </div>
            </div>
            <div className="flex items-center space-x-2">
              <div className={`w-2 h-2 rounded-full ${connection.connected ? 'bg-green-400' : 'bg-red-400'}`}></div>
              <span className="text-sm">{connection.connected ? 'Connected' : 'Offline'}</span>
            </div>
          </div>
        </div>

        {/* Navigation */}
        <div className="container mx-auto px-4">
          <div className="flex space-x-1 pb-4">
            <button className="bg-white text-blue-600 hover:bg-gray-100 px-4 py-2 rounded text-sm font-medium flex items-center">
              🏠 Dashboard
            </button>
            <Link href="/institution" className="text-white/80 hover:text-white hover:bg-white/10 px-4 py-2 rounded text-sm font-medium flex items-center">
              🏛️ Institution
            </Link>
            <Link href="/student" className="text-white/80 hover:text-white hover:bg-white/10 px-4 py-2 rounded text-sm font-medium flex items-center">
              👨‍🎓 Student
            </Link>
            <Link href="/equivalences" className="text-white/80 hover:text-white hover:bg-white/10 px-4 py-2 rounded text-sm font-medium flex items-center">
              🔄 Equivalences
            </Link>
            <Link href="/degree" className="text-white/80 hover:text-white hover:bg-white/10 px-4 py-2 rounded text-sm font-medium flex items-center">
              📜 Degrees
            </Link>
          </div>
        </div>
      </div>

      {/* Error Alert */}
      {error && (
        <div className="bg-red-50 border border-red-200 p-4">
          <div className="container mx-auto px-4">
            <p className="text-red-700 text-sm">⚠️ Error: {error}</p>
          </div>
        </div>
      )}

      {/* Main Content */}
      <div className="container mx-auto px-4 py-6">
        {/* Statistics Cards */}
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
          <div className="bg-white rounded-lg shadow p-4 text-center">
            <div className="text-2xl font-bold text-green-600">{stats.completed}</div>
            <div className="text-sm text-gray-600">Completed</div>
          </div>
          <div className="bg-white rounded-lg shadow p-4 text-center">
            <div className="text-2xl font-bold text-blue-600">{stats.enrolled}</div>
            <div className="text-sm text-gray-600">Enrolled</div>
          </div>
          <div className="bg-white rounded-lg shadow p-4 text-center">
            <div className="text-2xl font-bold text-yellow-600">{stats.available}</div>
            <div className="text-sm text-gray-600">Available</div>
          </div>
          <div className="bg-white rounded-lg shadow p-4 text-center">
            <div className="text-2xl font-bold text-gray-600">{stats.blocked}</div>
            <div className="text-sm text-gray-600">Blocked</div>
          </div>
        </div>

        {/* Legend */}
        <div className="bg-white rounded-lg shadow p-4 mb-6">
          <h3 className="font-semibold mb-3">Status Legend:</h3>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div className="flex items-center space-x-2">
              <span className="text-green-600">✅</span>
              <span className="text-sm">Completed</span>
            </div>
            <div className="flex items-center space-x-2">
              <span className="text-blue-600">📚</span>
              <span className="text-sm">Enrolled</span>
            </div>
            <div className="flex items-center space-x-2">
              <span className="text-yellow-600">🏆</span>
              <span className="text-sm">Available</span>
            </div>
            <div className="flex items-center space-x-2">
              <span className="text-gray-400">🔒</span>
              <span className="text-sm">Blocked</span>
            </div>
          </div>
        </div>

        {/* Curriculum Content */}
        {isLoading ? (
          <div className="flex items-center justify-center h-64 bg-white rounded-lg shadow">
            <div className="text-center">
              <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
              <p className="text-gray-600">Loading from blockchain...</p>
            </div>
          </div>
        ) : subjects.length === 0 ? (
          <div className="bg-white rounded-lg shadow p-12 text-center">
            <div className="text-gray-400 text-6xl mb-4">📚</div>
            <h3 className="text-xl font-semibold text-gray-600 mb-2">Blockchain is empty</h3>
            <p className="text-gray-500 mb-4">
              {connection.connected 
                ? "No institutions or subjects have been registered yet."
                : "Connect to blockchain to see registered data."
              }
            </p>
            {connection.connected && (
              <div className="text-sm text-gray-400">
                <p>💡 To get started:</p>
                <p>1. Go to Institution page to register institutions</p>
                <p>2. Add courses and subjects</p>
                <p>3. Data will appear here automatically</p>
              </div>
            )}
          </div>
        ) : (
          <div className="space-y-8">
            {Object.entries(subjectsBySemester)
              .sort(([a], [b]) => parseInt(a) - parseInt(b))
              .map(([semester, semesterSubjects]) => (
                <div key={semester}>
                  <h2 className="text-lg font-semibold flex items-center mb-4">
                    📚 Semester {semester}
                  </h2>
                  <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    {semesterSubjects.map((subject, index) => (
                      <div
                        key={`${subject.subjectId}-${index}`}
                        className={`bg-white rounded-lg shadow hover:shadow-md transition-shadow border-2 p-4 ${getStatusColor(subject.status)}`}
                      >
                        <div className="flex justify-between items-start mb-3">
                          <div className="flex-1">
                            <div className="flex items-center space-x-2 mb-1">
                              <span className="text-xs bg-white px-2 py-1 rounded border border-gray-300 font-mono">
                                {subject.code}
                              </span>
                              <span className="text-xs bg-white px-2 py-1 rounded border border-gray-300">
                                {getStatusText(subject.status)}
                              </span>
                            </div>
                            <h3 className="font-medium text-sm mb-1">{subject.title || subject.name}</h3>
                            <p className="text-xs text-gray-600">{subject.credits} credits</p>
                          </div>
                          <div className="ml-2 text-lg">{getStatusIcon(subject.status)}</div>
                        </div>

                        {/* Prerequisites */}
                        {subject.prerequisites && subject.prerequisites.length > 0 && (
                          <div className="mt-3 pt-3 border-t border-gray-200">
                            <p className="text-xs font-medium text-gray-700 mb-1">Prerequisites:</p>
                            <div className="flex flex-wrap gap-1">
                              {subject.prerequisites.map((prereqId) => {
                                const prereqSubject = enhancedSubjects.find(s => s.subjectId === prereqId)
                                const isCompleted = studentNFTs.some(nft => nft.subjectId === prereqId)
                                return (
                                  <span
                                    key={prereqId}
                                    className={`text-xs px-2 py-1 rounded border ${
                                      isCompleted
                                        ? 'bg-green-50 text-green-700 border-green-300'
                                        : 'bg-red-50 text-red-700 border-red-300'
                                    }`}
                                  >
                                    {prereqSubject?.code || prereqId}
                                    {isCompleted && ' ✓'}
                                  </span>
                                )
                              })}
                            </div>
                          </div>
                        )}

                        {/* Grade if completed */}
                        {subject.status === 'completed' && subject.grade && (
                          <div className="mt-3 pt-3 border-t border-gray-200">
                            <p className="text-xs text-gray-600">
                              Grade: <span className="font-semibold text-green-600">{subject.grade.toFixed(1)}</span>
                            </p>
                          </div>
                        )}

                        {/* Action Button */}
                        <div className="mt-3">
                          {subject.status === 'available' && (
                            <button className="w-full bg-yellow-500 hover:bg-yellow-600 text-white px-3 py-2 rounded text-sm font-medium transition-colors">
                              Enroll
                            </button>
                          )}
                          {subject.status === 'enrolled' && (
                            <button className="w-full border border-blue-300 text-blue-700 px-3 py-2 rounded text-sm font-medium">
                              In Progress
                            </button>
                          )}
                          {subject.status === 'completed' && (
                            <button className="w-full border border-green-300 text-green-700 px-3 py-2 rounded text-sm font-medium flex items-center justify-center">
                              ✓ Completed
                              {connection.connected && (
                                <span className="ml-2 text-xs">NFT</span>
                              )}
                            </button>
                          )}
                          {subject.status === 'blocked' && (
                            <button disabled className="w-full border border-gray-300 text-gray-400 px-3 py-2 rounded text-sm font-medium cursor-not-allowed">
                              Blocked
                            </button>
                          )}
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              ))}
          </div>
        )}
      </div>

      {/* Footer */}
      <footer className="bg-gray-800 text-white p-3 mt-8">
        <div className="container mx-auto text-center text-sm">
          <p>🎓 AcademicToken - Decentralized Academic System</p>
          {!connection.connected && (
            <p className="text-yellow-300 text-xs mt-1">
              ⚠️ Demo mode - Blockchain offline
            </p>
          )}
        </div>
      </footer>
    </div>
  )
}
