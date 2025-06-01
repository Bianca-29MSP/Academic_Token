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
    course: "Loading...",
    institutionId: ""
  }

  // Real NFTs will be loaded from blockchain - for now empty until student system is complete
  const studentNFTs: any[] = []

  // Enhance curriculum with grades from real blockchain data
  const enhancedCurriculum = subjects.map(subject => {
    const nft = studentNFTs.find(nft => nft.subjectId === subject.id)
    return {
      ...subject,
      grade: nft?.grade,
      completed: !!nft
    }
  })

  // Calculate stats from real data
  const completedCount = studentNFTs.length
  const totalCount = subjects.length
  const availableCount = enhancedCurriculum.filter(subject => !subject.completed).length

  const handleSubjectClick = (subject: any) => {
    console.log('Subject clicked:', subject)
  }

  return (
    <div className="min-h-screen flex flex-col">
      {/* Header */}
      <header className="bg-gradient-to-r from-blue-600 to-purple-600 text-white p-4">
        <div className="container mx-auto">
          <div className="flex items-center justify-between mb-4">
            <div>
              <h1 className="text-3xl font-bold">🎓 AcademicToken</h1>
              <p className="text-blue-100">
                {currentStudent.name} - {currentStudent.course}
              </p>
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
            <Link href="/" className="bg-white text-blue-600 px-4 py-2 rounded-lg font-semibold">
              🎓 Dashboard
            </Link>
            <Link href="/institution" className="bg-white/20 hover:bg-white/30 px-4 py-2 rounded-lg transition-colors">
              🏛️ Institution
            </Link>
            <Link href="/student" className="bg-white/20 hover:bg-white/30 px-4 py-2 rounded-lg transition-colors">
              👨‍🎓 Student
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

      {/* Error Alert */}
      {error && (
        <div className="bg-red-50 border border-red-200 p-4">
          <div className="container mx-auto">
            <p className="text-red-700 text-sm">⚠️ Error: {error}</p>
          </div>
        </div>
      )}

      {/* Stats */}
      <div className="bg-gray-50 p-4">
        <div className="container mx-auto">
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div className="bg-white p-3 rounded-lg shadow text-center">
              <div className="text-2xl font-bold text-green-600">{completedCount}</div>
              <div className="text-sm text-gray-600">Completed</div>
            </div>
            <div className="bg-white p-3 rounded-lg shadow text-center">
              <div className="text-2xl font-bold text-blue-600">{availableCount}</div>
              <div className="text-sm text-gray-600">Available</div>
            </div>
            <div className="bg-white p-3 rounded-lg shadow text-center">
              <div className="text-2xl font-bold text-purple-600">
                {totalCount > 0 ? Math.round((completedCount / totalCount) * 100) : 0}%
              </div>
              <div className="text-sm text-gray-600">Progress</div>
            </div>
            <div className="bg-white p-3 rounded-lg shadow text-center">
              <div className="text-2xl font-bold text-orange-600">{institutions.length}</div>
              <div className="text-sm text-gray-600">Institutions</div>
            </div>
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="flex-1 bg-gradient-to-br from-blue-50 to-purple-50">
        <div className="p-6">
          {isLoading ? (
            <div className="flex items-center justify-center h-64">
              <div className="text-center">
                <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
                <p className="text-gray-600">Loading from blockchain...</p>
              </div>
            </div>
          ) : (
            <div className="max-w-6xl mx-auto">
              <h2 className="text-2xl font-bold text-gray-800 mb-6">📚 Curriculum</h2>
              
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                {enhancedCurriculum.map((subject) => (
                  <div
                    key={subject.id}
                    onClick={() => handleSubjectClick(subject)}
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
                      <span className="text-gray-500">{subject.credits} credits</span>
                      {subject.completed && connection.connected && (
                        <span className="text-green-600 text-xs">NFT</span>
                      )}
                    </div>
                    
                    {subject.metadata && (
                      <p className="text-xs text-gray-500 mt-2">
                        {subject.metadata}
                      </p>
                    )}
                  </div>
                ))}
              </div>

              {subjects.length === 0 && !isLoading && (
                <div className="text-center py-12">
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
              )}
            </div>
          )}
        </div>
      </div>

      {/* Footer */}
      <footer className="bg-gray-800 text-white p-3">
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
