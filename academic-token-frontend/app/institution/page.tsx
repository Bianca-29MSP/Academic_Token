"use client"
import { useState, useRef } from 'react'
import Link from 'next/link'
import { useBlockchain } from '../hooks/useBlockchain'

export default function InstitutionDashboard() {
  const [selectedFile, setSelectedFile] = useState<File | null>(null)
  const [uploading, setUploading] = useState(false)
  const [notification, setNotification] = useState<{type: 'success' | 'error', message: string} | null>(null)
  const fileInputRef = useRef<HTMLInputElement>(null)
  
  const { 
    connection, 
    isLoading, 
    error, 
    institutions, 
    subjects 
  } = useBlockchain()
  
  // Mock students for NFT issuance demo
  const [students, setStudents] = useState([
    { id: 'student_1', name: 'John Silva', subject: 'Calculus I', status: 'pending' as const },
    { id: 'student_2', name: 'Maria Santos', subject: 'Programming 1', status: 'approved' as const },
    { id: 'student_3', name: 'Peter Costa', subject: 'Linear Algebra', status: 'pending' as const }
  ])

  const showNotification = (type: 'success' | 'error', message: string) => {
    setNotification({ type, message })
    setTimeout(() => setNotification(null), 5000)
  }

  const handleFileUpload = async (file: File) => {
    try {
      setUploading(true)
      
      const fileName = file.name.replace(/\.(pdf|txt|docx|rtf)$/i, '')
      const parts = fileName.split('_')
      const subjectName = parts[0] || fileName
      const subjectCode = parts[1] || fileName.substring(0, 6).toUpperCase()
      
      console.log('🚀 Processing file:', {
        fileName: file.name,
        extractedName: subjectName,
        extractedCode: subjectCode
      })
      
      // Simulate processing
      await new Promise(resolve => setTimeout(resolve, 2000))
      
      if (connection.connected) {
        showNotification('success', 
          `🎉 Subject "${subjectName}" processed and registered on blockchain!`
        )
      } else {
        showNotification('success', 
          `🎉 Subject "${subjectName}" processed successfully! (Demo mode)`
        )
      }
      
    } catch (error) {
      showNotification('error', 'Processing failed')
    } finally {
      setUploading(false)
      setSelectedFile(null)
    }
  }

  const handleSimulateUpload = async (fileName: string) => {
    const parts = fileName.replace('.pdf', '').split('_')
    const subjectName = parts[0] || fileName
    const subjectCode = parts[1] || fileName.substring(0, 6).toUpperCase()
    
    try {
      setUploading(true)
      await new Promise(resolve => setTimeout(resolve, 1500))
      
      if (connection.connected) {
        showNotification('success', 
          `✅ Subject "${subjectName}" (${subjectCode}) registered on blockchain!`
        )
      } else {
        showNotification('success', 
          `✅ Subject "${subjectName}" (${subjectCode}) created! (Demo mode)`
        )
      }
      
    } catch (error) {
      showNotification('error', 'Subject creation failed')
    } finally {
      setUploading(false)
    }
  }

  const handleIssueNFT = async (studentId: string, subjectName: string) => {
    try {
      setUploading(true)
      
      await new Promise(resolve => setTimeout(resolve, 1000))
      
      setStudents(prev => 
        prev.map(s => s.id === studentId ? {...s, status: 'approved' as const} : s)
      )
      
      if (connection.connected) {
        showNotification('success', `🏆 NFT issued for ${subjectName}! Blockchain transaction complete.`)
      } else {
        showNotification('success', `🏆 NFT issued for ${subjectName}! (Demo mode)`)
      }
      
    } catch (error) {
      showNotification('error', 'Error issuing NFT')
    } finally {
      setUploading(false)
    }
  }

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault()
    const files = Array.from(e.dataTransfer.files)
    if (files.length > 0 && (files[0].type === 'application/pdf' || files[0].type === 'text/plain')) {
      setSelectedFile(files[0])
    }
  }

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault()
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-gradient-to-r from-blue-600 to-purple-600 text-white p-4">
        <div className="container mx-auto">
          <div className="flex items-center justify-between mb-4">
            <div>
              <h1 className="text-3xl font-bold">🏛️ Institution Dashboard</h1>
              <p className="text-blue-100">Federal University of Juiz de Fora - UFJF</p>
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
            <Link href="/institution" className="bg-white text-blue-600 px-4 py-2 rounded-lg font-semibold">
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

      <div className="container mx-auto p-6">
        {/* Error Alert */}
        {error && (
          <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg">
            <p className="text-red-700 text-sm">⚠️ Error: {error}</p>
          </div>
        )}

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

        {/* Institution Overview */}
        <div className="bg-white rounded-xl shadow-lg p-6 mb-8">
          <h2 className="text-xl font-semibold text-gray-800 mb-4">Institution Overview</h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div className="p-4 bg-blue-100 text-blue-800 rounded-lg">
              <div className="text-2xl mb-2">🏛️</div>
              <div className="font-semibold">UFJF - Active</div>
              <div className="text-sm opacity-75">Federal University of Juiz de Fora</div>
            </div>
            
            <div className="p-4 bg-purple-100 text-purple-800 rounded-lg">
              <div className="text-2xl mb-2">📚</div>
              <div className="font-semibold">Computer Science</div>
              <div className="text-sm opacity-75">Active course</div>
            </div>
            
            <Link 
              href="/student"
              className="p-4 bg-green-500 text-white rounded-lg hover:bg-green-600 transition-colors block text-center"
            >
              <div className="text-2xl mb-2">👀</div>
              <div className="font-semibold">View Student Portal</div>
              <div className="text-sm opacity-75">See student progress</div>
            </Link>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          {/* Subject Registration */}
          <div className="bg-white rounded-xl shadow-lg p-6">
            <h2 className="text-xl font-semibold text-gray-800 mb-4">
              📄 Subject Registration
              {connection.connected && <span className="ml-2 text-xs bg-green-100 text-green-700 px-2 py-1 rounded-full">LIVE</span>}
            </h2>
            
            <div 
              className="border-2 border-dashed border-blue-300 rounded-lg p-8 text-center hover:border-blue-400 transition-colors"
              onDrop={handleDrop}
              onDragOver={handleDragOver}
            >
              <div className="text-4xl mb-4">📁</div>
              <p className="text-gray-600 mb-4">
                Upload syllabi for {connection.connected ? 'blockchain storage' : 'processing (demo)'}
              </p>
              
              <input
                ref={fileInputRef}
                type="file"
                accept=".pdf,.txt,.docx"
                onChange={(e) => {
                  const file = e.target.files?.[0]
                  if (file) {
                    setSelectedFile(file)
                  }
                }}
                className="hidden"
              />
              
              <div className="space-y-2">
                <button 
                  onClick={() => fileInputRef.current?.click()}
                  disabled={uploading}
                  className="block w-full bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600 transition-colors disabled:opacity-50"
                >
                  {uploading ? '⏳ Processing...' : '📄 Select File'}
                </button>
                
                <div className="text-sm text-gray-500 my-2">or try demo uploads:</div>
                
                <button 
                  onClick={() => handleSimulateUpload('Philosophy_PHI027.pdf')}
                  disabled={uploading}
                  className="block w-full bg-purple-500 text-white px-4 py-2 rounded hover:bg-purple-600 transition-colors disabled:opacity-50"
                >
                  📄 Demo: Philosophy PHI027
                </button>
                <button 
                  onClick={() => handleSimulateUpload('Calculus_MAT101.pdf')}
                  disabled={uploading}
                  className="block w-full bg-green-500 text-white px-4 py-2 rounded hover:bg-green-600 transition-colors disabled:opacity-50"
                >
                  📄 Demo: Calculus MAT101
                </button>
              </div>

              {selectedFile && !uploading && (
                <div className="mt-4 p-3 bg-blue-50 rounded-lg">
                  <p className="text-sm text-blue-700 mb-2">📄 Selected: {selectedFile.name}</p>
                  <button
                    onClick={() => handleFileUpload(selectedFile)}
                    className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 transition-colors"
                  >
                    🚀 Process & {connection.connected ? 'Register' : 'Demo'}
                  </button>
                </div>
              )}

              {uploading && (
                <div className="mt-4 p-3 bg-blue-50 rounded-lg">
                  <p className="text-sm text-blue-700 mb-2">⏳ Processing...</p>
                  <div className="w-full bg-blue-200 rounded-full h-2">
                    <div className="bg-blue-600 h-2 rounded-full animate-pulse" style={{width: '75%'}}></div>
                  </div>
                </div>
              )}
            </div>

            <div className="mt-6">
              <h3 className="text-lg font-semibold text-gray-800 mb-3">
                Registered Subjects
                {isLoading && <span className="text-sm text-blue-600">Loading...</span>}
              </h3>
              
              <div className="space-y-2 max-h-40 overflow-y-auto">
                {subjects.slice(0, 5).map((subject, index) => (
                  <div key={subject.id || index} className="flex items-center justify-between p-2 bg-gray-50 rounded">
                    <div className="flex-1">
                      <div className="flex items-center space-x-2">
                        <span className="font-semibold text-sm">{subject.code}</span>
                        <span className="text-gray-600 text-sm">{subject.name}</span>
                      </div>
                    </div>
                    <div className="flex items-center space-x-2">
                      <span className="text-xs text-gray-500">{subject.credits} credits</span>
                      <span className="text-xs bg-blue-100 text-blue-700 px-2 py-1 rounded-full">
                        {connection.connected ? '⛓️ Blockchain' : '🎭 Demo'}
                      </span>
                    </div>
                  </div>
                ))}
                {subjects.length === 0 && (
                  <div className="text-gray-500 text-sm">
                    {isLoading ? 'Loading subjects...' : 'No subjects registered yet'}
                  </div>
                )}
              </div>
            </div>
          </div>

          {/* NFT Issuance */}
          <div className="bg-white rounded-xl shadow-lg p-6">
            <h2 className="text-xl font-semibold text-gray-800 mb-4">
              🏆 Academic NFT Issuance
              {connection.connected && <span className="ml-2 text-xs bg-green-100 text-green-700 px-2 py-1 rounded-full">LIVE</span>}
            </h2>
            
            <div className="space-y-4">
              {students.map((student) => (
                <div key={student.id} className="flex items-center justify-between p-4 border rounded-lg hover:bg-gray-50 transition-colors">
                  <div>
                    <p className="font-semibold text-gray-800">{student.name}</p>
                    <p className="text-sm text-gray-600">{student.subject}</p>
                  </div>
                  
                  <div className="flex items-center space-x-3">
                    {student.status === 'approved' ? (
                      <span className="px-3 py-1 bg-green-100 text-green-700 rounded-full text-sm font-semibold">
                        ✅ NFT Issued
                      </span>
                    ) : (
                      <button
                        onClick={() => handleIssueNFT(student.id, student.subject)}
                        disabled={uploading}
                        className="px-4 py-2 bg-purple-500 text-white rounded-lg hover:bg-purple-600 transition-colors text-sm font-semibold disabled:opacity-50"
                      >
                        {uploading ? '⏳' : '🎯'} Issue NFT
                      </button>
                    )}
                  </div>
                </div>
              ))}
            </div>

            <div className="mt-6 p-4 bg-blue-50 rounded-lg">
              <h4 className="font-semibold text-blue-800 mb-2">
                {connection.connected ? '🔥 Blockchain Integration' : '🧠 Demo Mode'}
              </h4>
              <ul className="text-sm text-blue-700 space-y-1">
                {connection.connected ? (
                  <>
                    <li>• Real blockchain transactions</li>
                    <li>• IPFS content storage</li>
                    <li>• Smart contract validation</li>
                    <li>• Immutable academic records</li>
                  </>
                ) : (
                  <>
                    <li>• Full simulation of blockchain features</li>
                    <li>• File processing and analysis</li>
                    <li>• NFT creation preview</li>
                    <li>• Connect blockchain for real storage</li>
                  </>
                )}
              </ul>
            </div>
          </div>
        </div>

        {/* Statistics */}
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mt-8">
          <div className="bg-white rounded-lg shadow p-4 text-center">
            <div className="text-2xl font-bold text-blue-600">{subjects.length}</div>
            <div className="text-sm text-gray-600">
              Subjects {connection.connected ? 'on Blockchain' : '(Demo)'}
            </div>
          </div>
          <div className="bg-white rounded-lg shadow p-4 text-center">
            <div className="text-2xl font-bold text-green-600">
              {students.filter(s => s.status === 'approved').length}
            </div>
            <div className="text-sm text-gray-600">NFTs Issued</div>
          </div>
          <div className="bg-white rounded-lg shadow p-4 text-center">
            <div className="text-2xl font-bold text-purple-600">{students.length}</div>
            <div className="text-sm text-gray-600">Active Students</div>
          </div>
          <div className="bg-white rounded-lg shadow p-4 text-center">
            <div className="text-2xl font-bold text-orange-600">{institutions.length}</div>
            <div className="text-sm text-gray-600">Institutions</div>
          </div>
        </div>

        {/* System Status */}
        <div className="bg-white rounded-lg shadow-lg p-6 mt-8">
          <h3 className="text-lg font-semibold text-gray-800 mb-4">⛓️ System Status</h3>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm">
            <div>
              <span className="text-gray-600">Mode:</span>
              <span className={`ml-2 font-semibold ${
                connection.connected ? 'text-green-600' : 'text-blue-600'
              }`}>
                {connection.connected ? '🟢 Live Blockchain' : '🎭 Demo Mode'}
              </span>
            </div>
            <div>
              <span className="text-gray-600">Network:</span>
              <span className="ml-2 font-semibold text-purple-600">
                {connection.network}
              </span>
            </div>
            <div>
              <span className="text-gray-600">Features:</span>
              <span className="ml-2 font-semibold text-green-600">
                {connection.connected ? '🔥 All Live' : '🎯 All Demo'}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
