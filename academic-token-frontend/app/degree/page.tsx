"use client"
import { useState } from 'react'
import Link from 'next/link'

export default function DiplomaPage() {
  const [selectedStudent, setSelectedStudent] = useState<string>('')
  const [showDiploma, setShowDiploma] = useState(false)
  
  const eligibleStudents = [
    { 
      id: 'maria-santos', 
      name: 'Maria Santos', 
      course: 'Computer Science',
      progress: 100,
      credits: 240,
      totalCredits: 240,
      missingSubjects: [],
      canGraduate: true
    },
    { 
      id: 'joao-silva', 
      name: 'John Silva', 
      course: 'Computer Science',
      progress: 98,
      credits: 236,
      totalCredits: 240,
      missingSubjects: ['CI1100 - Supervised Internship'],
      canGraduate: false
    }
  ]

  const issuedDiplomas = [
    { id: 1, studentName: 'Ana Lima', degreeHash: '0x7f8c2e1a...', issueDate: '2024-05-15', verified: true },
    { id: 2, studentName: 'Carlos Souza', degreeHash: '0x9d4b6f2e...', issueDate: '2024-05-10', verified: true }
  ]

  const handleIssueDiploma = (studentId: string) => {
    setSelectedStudent(studentId)
    setShowDiploma(true)
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header with Navigation */}
      <header className="bg-gradient-to-r from-blue-600 to-purple-600 text-white p-4">
        <div className="container mx-auto">
          <div className="flex items-center justify-between mb-4">
            <div>
              <h1 className="text-3xl font-bold">📜 Degree System</h1>
              <p className="text-blue-100">Curriculum verification and degree NFT issuance</p>
            </div>
            
            <div className="text-sm bg-white/20 px-3 py-1 rounded-full">
              🔗 Connected
            </div>
          </div>

          {/* NAVIGATION */}
          <nav className="flex flex-wrap gap-2">
            <Link href="/" className="bg-white/20 hover:bg-white/30 px-4 py-2 rounded-lg transition-colors">
              🎓 Curriculum Tree
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
            <Link href="/degree" className="bg-white text-blue-600 px-4 py-2 rounded-lg font-semibold">
              📜 Degrees
            </Link>
          </nav>
        </div>
      </header>

      {showDiploma && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-xl p-8 max-w-2xl w-full mx-4 relative">
            <button 
              onClick={() => setShowDiploma(false)}
              className="absolute top-4 right-4 text-gray-500 hover:text-gray-700 text-2xl"
            >
              ×
            </button>
            
            <div className="text-center">
              <div className="text-6xl mb-4">🎓</div>
              <h2 className="text-2xl font-bold text-gray-800 mb-4">Digital Degree Issued!</h2>
              
              <div className="bg-gradient-to-r from-blue-500 to-purple-600 text-white p-6 rounded-lg mb-6">
                <h3 className="text-xl font-semibold mb-2">
                  {eligibleStudents.find(s => s.id === selectedStudent)?.name}
                </h3>
                <p className="mb-2">Bachelor of Computer Science</p>
                <p className="text-sm">Federal University of Juiz de Fora</p>
                <p className="text-sm">Hash: 0x8f4e2d9a5c7b1e6f3a8d2e5c9b4a7e6f</p>
              </div>

              <div className="space-y-3 text-sm text-gray-600">
                <p>✅ Degree registered on Cosmos blockchain</p>
                <p>✅ NFT transferred to student</p>
                <p>✅ Available for public verification</p>
              </div>

              <button
                onClick={() => setShowDiploma(false)}
                className="mt-6 bg-blue-500 text-white px-6 py-3 rounded-lg hover:bg-blue-600 transition-colors"
              >
                🔗 View on Blockchain
              </button>
            </div>
          </div>
        </div>
      )}

      <div className="container mx-auto p-6">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          {/* Eligibility Verification */}
          <div className="bg-white rounded-xl shadow-lg p-6">
            <h2 className="text-xl font-semibold text-gray-800 mb-6">🔍 Eligibility Verification</h2>
            
            <div className="space-y-4">
              {eligibleStudents.map((student, index) => (
                <div key={`student-${student.id || index}`} className="border rounded-lg p-4 hover:shadow-md transition-shadow">
                  <div className="flex justify-between items-start mb-3">
                    <div>
                      <h3 className="font-semibold text-gray-800">{student.name}</h3>
                      <p className="text-sm text-gray-600">{student.course}</p>
                    </div>
                    <span className={`px-3 py-1 rounded-full text-sm font-semibold ${
                      student.canGraduate 
                        ? 'bg-green-100 text-green-700' 
                        : 'bg-orange-100 text-orange-700'
                    }`}>
                      {student.canGraduate ? '✅ Eligible' : '⏳ Pending'}
                    </span>
                  </div>

                  <div className="mb-4">
                    <div className="flex justify-between text-sm text-gray-600 mb-1">
                      <span>Progress:</span>
                      <span>{student.credits}/{student.totalCredits} credits ({student.progress}%)</span>
                    </div>
                    <div className="w-full bg-gray-200 rounded-full h-2">
                      <div 
                        className={`h-2 rounded-full ${
                          student.progress === 100 ? 'bg-green-500' : 'bg-blue-500'
                        }`}
                        style={{width: `${student.progress}%`}}
                      ></div>
                    </div>
                  </div>

                  {student.missingSubjects.length > 0 && (
                    <div className="mb-4">
                      <p className="text-sm font-semibold text-gray-700 mb-2">Missing subjects:</p>
                      <div className="space-y-1">
                        {student.missingSubjects.map((subject, subjectIndex) => (
                          <div key={`missing-${subjectIndex}`} className="text-xs bg-orange-50 text-orange-700 px-2 py-1 rounded">
                            {subject}
                          </div>
                        ))}
                      </div>
                    </div>
                  )}

                  <div className="flex space-x-2">
                    <Link
                      href={`/student`}
                      className="flex-1 px-3 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 transition-colors text-sm text-center"
                    >
                      👀 View Transcript
                    </Link>
                    {student.canGraduate && (
                      <button
                        onClick={() => handleIssueDiploma(student.id)}
                        className="flex-1 px-3 py-2 bg-green-500 text-white rounded hover:bg-green-600 transition-colors text-sm"
                      >
                        🎓 Issue Degree
                      </button>
                    )}
                  </div>
                </div>
              ))}
            </div>
          </div>

          {/* Issued Degrees */}
          <div className="bg-white rounded-xl shadow-lg p-6">
            <h2 className="text-xl font-semibold text-gray-800 mb-6">🏆 Issued Degrees</h2>
            
            <div className="space-y-4">
              {issuedDiplomas.map((degree, index) => (
                <div key={`degree-${degree.id || index}`} className="border border-green-200 bg-green-50 rounded-lg p-4">
                  <div className="flex justify-between items-start mb-3">
                    <div>
                      <h3 className="font-semibold text-gray-800">{degree.studentName}</h3>
                      <p className="text-sm text-gray-600">Issued on {degree.issueDate}</p>
                    </div>
                    <span className="px-3 py-1 bg-green-100 text-green-700 rounded-full text-sm font-semibold">
                      ✅ Verified
                    </span>
                  </div>

                  <div className="bg-white p-3 rounded mb-3">
                    <p className="text-xs font-mono text-gray-600">Degree Hash:</p>
                    <p className="text-sm font-mono break-all">{degree.degreeHash}</p>
                  </div>

                  <button className="w-full px-3 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 transition-colors text-sm">
                    🔗 View on Blockchain
                  </button>
                </div>
              ))}
            </div>

            <div className="mt-6 pt-4 border-t">
              <div className="text-center">
                <div className="text-2xl font-bold text-green-600">{issuedDiplomas.length}</div>
                <div className="text-sm text-gray-600">Degrees issued this semester</div>
              </div>
            </div>
          </div>
        </div>

        {/* Statistics */}
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mt-8">
          <div className="bg-white rounded-lg shadow p-4 text-center">
            <div className="text-2xl font-bold text-green-600">127</div>
            <div className="text-sm text-gray-600">Degrees Issued</div>
          </div>
          <div className="bg-white rounded-lg shadow p-4 text-center">
            <div className="text-2xl font-bold text-blue-600">34</div>
            <div className="text-sm text-gray-600">Ready to Graduate</div>
          </div>
          <div className="bg-white rounded-lg shadow p-4 text-center">
            <div className="text-2xl font-bold text-purple-600">856</div>
            <div className="text-sm text-gray-600">Active Students</div>
          </div>
          <div className="bg-white rounded-lg shadow p-4 text-center">
            <div className="text-2xl font-bold text-orange-600">100%</div>
            <div className="text-sm text-gray-600">Verification Rate</div>
          </div>
        </div>
      </div>
    </div>
  )
}