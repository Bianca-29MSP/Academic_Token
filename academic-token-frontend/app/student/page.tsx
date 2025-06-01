"use client"
import { useState, useCallback } from 'react'
import Link from 'next/link'
import { useBlockchain } from '../hooks/useBlockchain'

export default function StudentPortal() {
  const [selectedNFT, setSelectedNFT] = useState<number | null>(null)
  const [requestingEnrollment, setRequestingEnrollment] = useState<string | null>(null)
  const [notification, setNotification] = useState<{type: 'success' | 'error', message: string} | null>(null)
  
  const { 
    connection, 
    isLoading, 
    error, 
    subjects 
  } = useBlockchain()

  // Mock student data
  const currentStudent = {
    id: "student_123",
    name: "John Silva",
    course: "Computer Science",
    institutionId: "inst_1"
  }

  // Mock NFTs (completed subjects)
  const mockNFTs = [
    { 
      id: "nft_1", 
      subjectId: "subj_1", 
      grade: 8.5, 
      completionDate: "2024-03-15T10:00:00Z",
      nftHash: "0x123abc456def789",
      metadata: { subject: "Calculus I", credits: 4, institution: "UFJF" }
    },
    { 
      id: "nft_2", 
      subjectId: "subj_2", 
      grade: 9.2, 
      completionDate: "2024-04-20T14:30:00Z",
      nftHash: "0x456def789abc123",
      metadata: { subject: "Programming 1", credits: 4, institution: "UFJF" }
    }
  ]

  const showNotification = useCallback((type: 'success' | 'error', message: string) => {
    setNotification({ type, message })
    setTimeout(() => setNotification(null), 5000)
  }, [])

  const handleRequestEnrollment = useCallback(async (subjectId: string) => {
    try {
      setRequestingEnrollment(subjectId)
      await new Promise(resolve => setTimeout(resolve, 1000))
      
      const subjectName = subjects.find(s => s.id === subjectId)?.name || subjectId
      
      if (connection.connected) {
        showNotification('success', `✅ Enrollment request submitted! Subject: ${subjectName}`)
      } else {
        showNotification('success', `✅ Enrollment request simulated! Subject: ${subjectName}`)
      }
      
    } catch (error) {
      showNotification('error', 'Enrollment request failed')
    } finally {
      setRequestingEnrollment(null)
    }
  }, [connection.connected, showNotification, subjects])

  const handleSubmitWork = useCallback(async (subjectId: string) => {
    try {
      const subjectName = subjects.find(s => s.id === subjectId)?.name || subjectId
      await new Promise(resolve => setTimeout(resolve, 1000))
      
      if (connection.connected) {
        showNotification('success', `📝 Work submitted for review! Subject: ${subjectName}`)
      } else {
        showNotification('success', `📝 Work submission simulated! Subject: ${subjectName}`)
      }
      
    } catch (error) {
      showNotification('error', 'Work submission failed')
    }
  }, [subjects, showNotification, connection.connected])

  // Map NFTs for UI
  const studentNFTs = mockNFTs.map((nft, index) => ({
    id: index + 1,
    subject: nft.metadata?.subject || 'Unknown Subject',
    code: subjects.find(s => s.id === nft.subjectId)?.code || 'N/A',
    grade: nft.grade || 0,
    credits: nft.metadata?.credits || 4,
    date: nft.completionDate || new Date().toISOString(),
    status: 'completed' as const,
    nftHash: nft.nftHash || 'mock-hash',
    subjectId: nft.subjectId || 'unknown'
  }))

  // Available subjects (not completed)
  const availableSubjects = subjects.filter(subject => 
    !mockNFTs.some(nft => nft.subjectId === subject.id)
  ).map((subject, index) => ({
    id: studentNFTs.length + index + 1,
    subject: subject.name,
    code: subject.code,
    grade: null,
    credits: subject.credits,
    date: null,
    status: 'available' as const,
    subjectId: subject.id
  }))

  // Mock enrolled subjects
  const enrolledSubjects = [
    {
      id: 999,
      subject: 'Programming 2',
      code: 'CI1002',
      grade: null,
      credits: 4,
      date: null,
      status: 'enrolled' as const,
      subjectId: 'ci1002'
    }
  ]

  const allSubjects = [...studentNFTs, ...enrolledSubjects, ...availableSubjects]
  const totalCredits = studentNFTs.reduce((sum, n) => sum + (n.credits || 0), 0)
  const progress = Math.round((totalCredits / 240) * 100)

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gray-50">
        <header className="bg-gradient-to-r from-blue-600 to-purple-600 text-white p-4">
          <div className="container mx-auto">
            <h1 className="text-3xl font-bold">🎓 Student Portal</h1>
            <p className="text-blue-100">Loading...</p>
          </div>
        </header>
        
        <div className="flex items-center justify-center" style={{height: 'calc(100vh - 140px)'}}>
          <div className="text-center">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
            <p className="text-gray-600">Loading student data...</p>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Notification */}
      {notification && (
        <div className={`fixed top-4 right-4 z-50 p-4 rounded-lg shadow-lg ${
          notification.type === 'success' 
            ? 'bg-green-100 border border-green-200 text-green-800' 
            : 'bg-red-100 border border-red-200 text-red-800'
        }`}>
          <div className="text-sm">{notification.message}</div>
        </div>
      )}

      {/* Header */}
      <header className="bg-gradient-to-r from-blue-600 to-purple-600 text-white p-4">
        <div className="container mx-auto">
          <div className="flex items-center justify-between mb-4">
            <div>
              <h1 className="text-3xl font-bold">🎓 Student Portal</h1>
              <p className="text-blue-100">{currentStudent.name} - {currentStudent.course}</p>
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
            <Link href="/student" className="bg-white text-blue-600 px-4 py-2 rounded-lg font-semibold">
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

      <div className="container mx-auto p-6">
        {/* Error Alert */}
        {error && (
          <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg">
            <p className="text-red-700 text-sm">⚠️ Error: {error}</p>
          </div>
        )}

        {/* Academic Progress */}
        <div className="bg-gradient-to-r from-blue-500 to-purple-600 rounded-xl shadow-lg p-6 text-white mb-8">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-xl font-semibold">📊 Academic Progress</h2>
            {connection.connected && (
              <span className="text-xs bg-white/20 px-2 py-1 rounded-full">⛓️ Live</span>
            )}
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <div>
              <div className="text-3xl font-bold">{totalCredits}/240</div>
              <div className="text-blue-100">Credits Earned</div>
            </div>
            <div>
              <div className="text-3xl font-bold">{progress}%</div>
              <div className="text-blue-100">Course Progress</div>
            </div>
            <div>
              <div className="text-3xl font-bold">{studentNFTs.length}</div>
              <div className="text-blue-100">
                {connection.connected ? 'NFTs on Blockchain' : 'Academic NFTs'}
              </div>
            </div>
          </div>
          
          <div className="mt-4">
            <div className="w-full bg-white/20 rounded-full h-3">
              <div 
                className="bg-white h-3 rounded-full transition-all duration-1000"
                style={{width: `${progress}%`}}
              ></div>
            </div>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Subjects Overview */}
          <div className="lg:col-span-2">
            <h2 className="text-2xl font-semibold text-gray-800 mb-6">📚 My Academic Journey</h2>
            
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {allSubjects.map((subject) => (
                <div 
                  key={subject.id}
                  onClick={() => setSelectedNFT(subject.id)}
                  className={`bg-white rounded-lg shadow-lg p-4 cursor-pointer transition-all duration-200 hover:scale-105 ${
                    selectedNFT === subject.id ? 'ring-2 ring-blue-500' : ''
                  }`}
                >
                  <div className="flex items-center justify-between mb-3">
                    <span className="text-xs font-mono bg-gray-100 px-2 py-1 rounded">
                      {subject.code}
                    </span>
                    <div className="flex items-center space-x-1">
                      <span className="text-lg">
                        {subject.status === 'completed' ? '✅' : 
                         subject.status === 'enrolled' ? '📚' : '🔓'}
                      </span>
                      {subject.status === 'completed' && connection.connected && (
                        <span className="text-xs text-green-600">⛓️</span>
                      )}
                    </div>
                  </div>
                  
                  <h3 className="font-semibold text-gray-800 mb-2">{subject.subject}</h3>
                  
                  <div className="space-y-1 text-sm text-gray-600">
                    <div>📚 {subject.credits} credits</div>
                    {subject.grade && <div>📊 Grade: {subject.grade.toFixed(1)}</div>}
                    {subject.date && <div>📅 {new Date(subject.date).toLocaleDateString()}</div>}
                  </div>

                  <div className="mt-3">
                    {subject.status === 'completed' && (
                      <span className="px-2 py-1 bg-green-100 text-green-700 rounded-full text-xs font-semibold block text-center">
                        🏆 NFT Earned
                      </span>
                    )}
                    
                    {subject.status === 'enrolled' && (
                      <div className="space-y-2">
                        <span className="px-2 py-1 bg-blue-100 text-blue-700 rounded-full text-xs font-semibold block text-center">
                          📚 Currently Enrolled
                        </span>
                        <button 
                          onClick={(e) => {
                            e.stopPropagation()
                            handleSubmitWork((subject as any).subjectId)
                          }}
                          className="w-full px-2 py-1 bg-purple-100 text-purple-700 rounded-full text-xs font-semibold hover:bg-purple-200 transition-colors"
                        >
                          📝 Submit Work
                        </button>
                      </div>
                    )}
                    
                    {subject.status === 'available' && (
                      <button 
                        onClick={(e) => {
                          e.stopPropagation()
                          handleRequestEnrollment((subject as any).subjectId)
                        }}
                        disabled={requestingEnrollment === (subject as any).subjectId}
                        className="w-full px-2 py-1 bg-yellow-100 text-yellow-700 rounded-full text-xs font-semibold hover:bg-yellow-200 transition-colors disabled:opacity-50"
                      >
                        {requestingEnrollment === (subject as any).subjectId ? '⏳' : '📝'} 
                        Request Enrollment
                      </button>
                    )}
                  </div>
                </div>
              ))}
            </div>
          </div>

          {/* Sidebar */}
          <div className="space-y-6">
            {/* Quick Actions */}
            <div className="bg-white rounded-lg shadow-lg p-6">
              <h3 className="text-lg font-semibold text-gray-800 mb-4">⚡ Quick Actions</h3>
              <div className="space-y-3">
                <Link 
                  href="/equivalences"
                  className="block p-3 bg-blue-50 rounded-lg hover:bg-blue-100 transition-colors"
                >
                  <div className="font-semibold text-blue-700">🔄 Request Equivalence</div>
                  <div className="text-sm text-blue-600">Smart contract analysis</div>
                </Link>
                
                <Link 
                  href="/"
                  className="block p-3 bg-green-50 rounded-lg hover:bg-green-100 transition-colors"
                >
                  <div className="font-semibold text-green-700">🎓 View Curriculum</div>
                  <div className="text-sm text-green-600">Check prerequisites</div>
                </Link>
                
                <Link 
                  href="/degree"
                  className="block p-3 bg-purple-50 rounded-lg hover:bg-purple-100 transition-colors"
                >
                  <div className="font-semibold text-purple-700">📜 Graduation Status</div>
                  <div className="text-sm text-purple-600">Check eligibility</div>
                </Link>
              </div>
            </div>

            {/* NFT Details */}
            {selectedNFT && (
              <div className="bg-white rounded-lg shadow-lg p-6">
                <h3 className="text-lg font-semibold text-gray-800 mb-4">🔍 Subject Details</h3>
                {(() => {
                  const subject = allSubjects.find(n => n.id === selectedNFT)
                  const nftData = studentNFTs.find(nft => nft.code === subject?.code)
                  
                  return subject ? (
                    <div className="space-y-2 text-sm">
                      {nftData && (
                        <>
                          <div><strong>NFT Hash:</strong> {nftData.nftHash.slice(0, 12)}...</div>
                          <div><strong>Subject ID:</strong> {nftData.subjectId}</div>
                          {connection.connected && (
                            <>
                              <div><strong>Network:</strong> AcademicToken</div>
                              <div><strong>Contract:</strong> AcademicNFT</div>
                            </>
                          )}
                        </>
                      )}
                      <div><strong>Institution:</strong> UFJF</div>
                      <div><strong>Credits:</strong> {subject.credits}</div>
                      <div><strong>Status:</strong> {
                        subject.status === 'completed' ? '✅ Completed' :
                        subject.status === 'enrolled' ? '📚 Enrolled' : '🔓 Available'
                      }</div>
                      {subject.status === 'completed' && subject.date && (
                        <div><strong>Completion:</strong> {new Date(subject.date).toLocaleDateString()}</div>
                      )}
                      
                      {nftData && (
                        <button 
                          onClick={() => {
                            showNotification('success', 
                              connection.connected 
                                ? `🔗 Opening explorer for NFT: ${nftData.nftHash.slice(0, 8)}...`
                                : '🎭 Demo mode - Explorer not available'
                            )
                          }}
                          className="w-full mt-3 bg-blue-500 text-white px-3 py-2 rounded hover:bg-blue-600 transition-colors"
                        >
                          🔗 {connection.connected ? 'View on Explorer' : 'Demo Explorer'}
                        </button>
                      )}
                    </div>
                  ) : null
                })()}
              </div>
            )}
          </div>
        </div>

        {/* Statistics */}
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mt-8">
          <div className="bg-white rounded-lg shadow p-4 text-center">
            <div className="text-2xl font-bold text-green-600">{studentNFTs.length}</div>
            <div className="text-sm text-gray-600">Completed</div>
          </div>
          <div className="bg-white rounded-lg shadow p-4 text-center">
            <div className="text-2xl font-bold text-blue-600">{enrolledSubjects.length}</div>
            <div className="text-sm text-gray-600">In Progress</div>
          </div>
          <div className="bg-white rounded-lg shadow p-4 text-center">
            <div className="text-2xl font-bold text-purple-600">{totalCredits}</div>
            <div className="text-sm text-gray-600">Credits Earned</div>
          </div>
          <div className="bg-white rounded-lg shadow p-4 text-center">
            <div className="text-2xl font-bold text-orange-600">{progress}%</div>
            <div className="text-sm text-gray-600">Progress</div>
          </div>
        </div>
      </div>
    </div>
  )
}
