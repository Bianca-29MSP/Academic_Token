"use client"
import { useState, useEffect } from 'react'
import Link from 'next/link'
import { useBlockchain } from '../hooks/useBlockchain'
import { AcademicTokenAPI } from '../lib/api'

export default function EquivalencesPage() {
  const [activeTab, setActiveTab] = useState<'request' | 'pending' | 'approved' | 'debug'>('request')
  const [showSuccess, setShowSuccess] = useState(false)
  
  // Blockchain data
  const { 
    connection, 
    isLoading: blockchainLoading, 
    error: blockchainError,
    institutions, 
    subjects,
    debug,
    getSubjectsByInstitution,
    refreshConnection,
    refreshData,
    clearErrors
  } = useBlockchain()

  // Form state for equivalence request
  const [sourceInstitution, setSourceInstitution] = useState<string>('')
  const [sourceSubject, setSourceSubject] = useState<string>('')
  const [targetInstitution, setTargetInstitution] = useState<string>('')
  const [targetSubject, setTargetSubject] = useState<string>('')
  const [analysis, setAnalysis] = useState<{
    similarity: number;
    recommendation: 'approved' | 'rejected';
  } | null>(null)
  const [analyzingEquivalence, setAnalyzingEquivalence] = useState(false)

  // Real equivalence data loaded from blockchain - empty until we have real requests
  const pendingEquivalences: any[] = []
  const approvedEquivalences: any[] = []

  // TODO: Load real equivalence requests from blockchain
  // useEffect(() => {
  //   if (connection.connected) {
  //     loadEquivalenceRequests()
  //   }
  // }, [connection.connected])

  // Get subjects for selected institutions
  const sourceSubjects = sourceInstitution ? getSubjectsByInstitution(sourceInstitution) : []
  const targetSubjects = targetInstitution ? getSubjectsByInstitution(targetInstitution) : []

  // Auto-analyze equivalence when both subjects are selected
  useEffect(() => {
    if (sourceSubject && targetSubject && !analyzingEquivalence) {
      analyzeEquivalence()
    }
  }, [sourceSubject, targetSubject])

  const analyzeEquivalence = async () => {
    if (!sourceSubject || !targetSubject) return

    setAnalyzingEquivalence(true)
    try {
      // Call real equivalence contract on blockchain
      console.log('📡 Calling equivalence contract...', { sourceSubject, targetSubject })
      
      // TODO: Replace with actual contract call
      // const contractResult = await callEquivalenceContract(sourceSubject, targetSubject)
      
      // For now, show that contract call is needed
      throw new Error('Equivalence analysis requires deployed CosmWasm contract')
      
    } catch (error) {
      console.error('❌ Equivalence analysis failed:', error)
      setAnalysis(null)
      // Don't show fake analysis - show error instead
    } finally {
      setAnalyzingEquivalence(false)
    }
  }

  const handleApprove = (id: number) => {
    setShowSuccess(true)
    setTimeout(() => setShowSuccess(false), 3000)
  }

  const handleRequestEquivalence = async () => {
    if (!sourceInstitution || !sourceSubject || !targetInstitution || !targetSubject) {
      alert('Please fill all fields')
      return
    }

    try {
      console.log('📡 Submitting equivalence request to blockchain...', {
        sourceInstitution,
        sourceSubject, 
        targetInstitution,
        targetSubject
      })
      
      // TODO: Call real blockchain transaction
      // await blockchainServices.createEquivalenceService().requestEquivalence(
      //   'student_id', sourceSubject, targetSubject
      // )
      
      // For now, show that this needs blockchain integration
      throw new Error('Equivalence request requires blockchain transaction submission')
      
    } catch (error) {
      console.error('❌ Equivalence request failed:', error)
      alert(`Error: ${error instanceof Error ? error.message : 'Unknown error'}`)
      return
    }

    setShowSuccess(true)
    setTimeout(() => setShowSuccess(false), 3000)
    
    // Reset form
    setSourceInstitution('')
    setSourceSubject('')
    setTargetInstitution('')
    setTargetSubject('')
    setAnalysis(null)
  }

  const getConnectionStatusColor = () => {
    if (blockchainLoading) return 'bg-yellow-500/20 text-yellow-700'
    if (connection.connected) return 'bg-green-500/20 text-green-700'
    return 'bg-red-500/20 text-red-700'
  }

  const getConnectionStatusText = () => {
    if (blockchainLoading) return '🔄 Connecting...'
    if (connection.connected) return `🔗 Connected to ${connection.network}`
    return '❌ Disconnected'
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-gradient-to-r from-blue-600 to-purple-600 text-white p-4">
        <div className="container mx-auto">
          <div className="flex items-center justify-between mb-4">
            <div>
              <h1 className="text-3xl font-bold">🔄 Equivalences System</h1>
              <p className="text-blue-100">Automatic recognition between institutions</p>
            </div>
            
            <div className="flex items-center space-x-3">
              <div className={`text-sm px-3 py-1 rounded-full ${getConnectionStatusColor()}`}>
                {getConnectionStatusText()}
              </div>
              <button
                onClick={refreshConnection}
                className="bg-white/20 hover:bg-white/30 px-3 py-1 rounded text-sm transition-colors"
              >
                🔄 Reconnect
              </button>
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
            <Link href="/equivalences" className="bg-white text-blue-600 px-4 py-2 rounded-lg font-semibold">
              🔄 Equivalences
            </Link>
            <Link href="/degree" className="bg-white/20 hover:bg-white/30 px-4 py-2 rounded-lg transition-colors">
              📜 Degrees
            </Link>
          </nav>
        </div>
      </header>

      <div className="container mx-auto p-6">
        {/* Connection Error Alert */}
        {blockchainError && (
          <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg">
            <div className="flex justify-between items-start">
              <div>
                <p className="text-red-700 font-semibold">⚠️ Blockchain Connection Error</p>
                <p className="text-sm text-red-600">{blockchainError}</p>
                <p className="text-xs text-red-500 mt-1">
                  Make sure the REST server is running: <code>go run cmd/rest-server/main.go</code>
                </p>
              </div>
              <button
                onClick={clearErrors}
                className="text-red-500 hover:text-red-700 px-2 py-1 text-sm"
              >
                ✕
              </button>
            </div>
          </div>
        )}

        {showSuccess && (
          <div className="mb-6 p-4 bg-green-50 border border-green-200 rounded-lg">
            <p className="text-green-700 font-semibold">✅ Action completed successfully!</p>
            <p className="text-sm text-green-600">Equivalence processed on blockchain</p>
          </div>
        )}

        {/* Loading State */}
        {blockchainLoading && (
          <div className="mb-6 p-4 bg-blue-50 border border-blue-200 rounded-lg">
            <p className="text-blue-700 font-semibold">🔄 Loading blockchain data...</p>
            <p className="text-sm text-blue-600">Connecting to {connection.nodeUrl}</p>
          </div>
        )}

        {/* Stats */}
        {connection.connected && (
          <div className="mb-6 grid grid-cols-1 md:grid-cols-4 gap-4">
            <div className="bg-white p-4 rounded-lg shadow">
              <h3 className="text-lg font-semibold text-gray-800">🏛️ Institutions</h3>
              <p className="text-2xl font-bold text-blue-600">{institutions.length}</p>
              <p className="text-sm text-gray-600">Registered on blockchain</p>
            </div>
            <div className="bg-white p-4 rounded-lg shadow">
              <h3 className="text-lg font-semibold text-gray-800">📚 Subjects</h3>
              <p className="text-2xl font-bold text-purple-600">{subjects.length}</p>
              <p className="text-sm text-gray-600">Available for equivalence</p>
            </div>
            <div className="bg-white p-4 rounded-lg shadow">
              <h3 className="text-lg font-semibold text-gray-800">🌐 Network</h3>
              <p className="text-sm font-bold text-green-600">{connection.network}</p>
              <p className="text-sm text-gray-600">Version {connection.version}</p>
            </div>
            <div className="bg-white p-4 rounded-lg shadow">
              <h3 className="text-lg font-semibold text-gray-800">🔧 Debug</h3>
              <p className="text-sm font-bold text-orange-600">{debug.apiCalls} calls</p>
              <p className="text-sm text-gray-600">Last update: {debug.lastUpdate}</p>
            </div>
          </div>
        )}

        {/* Tabs */}
        <div className="bg-white rounded-lg shadow-lg mb-8">
          <div className="border-b border-gray-200">
            <nav className="flex space-x-8 px-6">
              {[
                { key: 'request', label: '📝 Request', count: 0 },
                { key: 'pending', label: '⏳ Pending', count: pendingEquivalences.length },
                { key: 'approved', label: '✅ Approved', count: approvedEquivalences.length },
                { key: 'debug', label: '🔧 Debug', count: 0 }
              ].map((tab) => (
                <button
                  key={tab.key}
                  onClick={() => setActiveTab(tab.key as any)}
                  className={`py-4 px-2 border-b-2 font-medium text-sm transition-colors ${
                    activeTab === tab.key
                      ? 'border-blue-500 text-blue-600'
                      : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                  }`}
                >
                  {tab.label}
                  {tab.count > 0 && (
                    <span className="ml-2 bg-blue-100 text-blue-600 px-2 py-1 rounded-full text-xs">
                      {tab.count}
                    </span>
                  )}
                </button>
              ))}
            </nav>
          </div>

          <div className="p-6">
            {/* Debug Tab */}
            {activeTab === 'debug' && (
              <div className="space-y-6">
                <h2 className="text-xl font-semibold text-gray-800">🔧 Debug Information</h2>
                
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div className="bg-gray-50 p-4 rounded-lg">
                    <h3 className="font-semibold mb-3">📡 Connection Status</h3>
                    <div className="space-y-2 text-sm">
                      <div><strong>Connected:</strong> {connection.connected ? '✅ Yes' : '❌ No'}</div>
                      <div><strong>Node URL:</strong> {connection.nodeUrl}</div>
                      <div><strong>Network:</strong> {connection.network}</div>
                      <div><strong>Version:</strong> {connection.version || 'Unknown'}</div>
                      <div><strong>Loading:</strong> {blockchainLoading ? '🔄 Yes' : '✅ No'}</div>
                    </div>
                  </div>

                  <div className="bg-gray-50 p-4 rounded-lg">
                    <h3 className="font-semibold mb-3">📊 Data Status</h3>
                    <div className="space-y-2 text-sm">
                      <div><strong>Institutions:</strong> {institutions.length} loaded</div>
                      <div><strong>Subjects:</strong> {subjects.length} loaded</div>
                      <div><strong>API Calls:</strong> {debug.apiCalls}</div>
                      <div><strong>Last Update:</strong> {debug.lastUpdate}</div>
                      <div><strong>Errors:</strong> {debug.errors.length}</div>
                    </div>
                  </div>
                </div>

                {/* Raw Data Display */}
                <div className="bg-gray-50 p-4 rounded-lg">
                  <h3 className="font-semibold mb-3">🗄️ Raw Data</h3>
                  <div className="space-y-4">
                    <div>
                      <h4 className="font-medium">Institutions ({institutions.length}):</h4>
                      <pre className="text-xs bg-white p-2 rounded mt-1 overflow-auto max-h-32">
                        {JSON.stringify(institutions, null, 2)}
                      </pre>
                    </div>
                    <div>
                      <h4 className="font-medium">Subjects ({subjects.length}):</h4>
                      <pre className="text-xs bg-white p-2 rounded mt-1 overflow-auto max-h-32">
                        {JSON.stringify(subjects, null, 2)}
                      </pre>
                    </div>
                  </div>
                </div>

                {/* Action Buttons */}
                <div className="flex space-x-4">
                  <button
                    onClick={refreshConnection}
                    className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600 transition-colors"
                  >
                    🔄 Refresh Connection
                  </button>
                  <button
                    onClick={refreshData}
                    className="bg-green-500 text-white px-4 py-2 rounded hover:bg-green-600 transition-colors"
                    disabled={!connection.connected}
                  >
                    📊 Refresh Data
                  </button>
                  <button
                    onClick={clearErrors}
                    className="bg-gray-500 text-white px-4 py-2 rounded hover:bg-gray-600 transition-colors"
                  >
                    🗑️ Clear Errors
                  </button>
                </div>
              </div>
            )}

            {/* Request Equivalence */}
            {activeTab === 'request' && (
              <div className="space-y-6">
                <h2 className="text-xl font-semibold text-gray-800">📝 Request New Equivalence</h2>
                <p className="text-sm text-gray-600">AI-powered content analysis via smart contracts</p>
                
                {/* Data availability check */}
                {!connection.connected && (
                  <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
                    <p className="text-yellow-700 font-semibold">⚠️ Blockchain Disconnected</p>
                    <p className="text-sm text-yellow-600">Connect to blockchain to load institutions and subjects</p>
                  </div>
                )}

                {connection.connected && institutions.length === 0 && (
                  <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
                    <p className="text-blue-700 font-semibold">📚 Blockchain is Empty</p>
                    <p className="text-sm text-blue-600">No institutions registered yet. Go to Institution page to add some first.</p>
                  </div>
                )}
                
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  {/* Source Institution & Subject */}
                  <div className="space-y-4">
                    <h3 className="font-semibold text-gray-700">🏫 Source Institution & Subject</h3>
                    <div className="space-y-3">
                      <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                          Institution ({institutions.length} available)
                        </label>
                        <select 
                          value={sourceInstitution}
                          onChange={(e) => {
                            setSourceInstitution(e.target.value)
                            setSourceSubject('') // Reset subject when institution changes
                          }}
                          className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                          disabled={blockchainLoading || !connection.connected}
                        >
                          <option value="">Select source institution...</option>
                          {institutions.map((inst) => (
                            <option key={inst.id} value={inst.id}>
                              {inst.name} ({inst.code})
                            </option>
                          ))}
                        </select>
                        {institutions.length === 0 && connection.connected && (
                          <p className="text-sm text-red-500 mt-1">No institutions found</p>
                        )}
                      </div>
                      
                      <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                          Subject ({sourceSubjects.length} available)
                        </label>
                        <select 
                          value={sourceSubject}
                          onChange={(e) => setSourceSubject(e.target.value)}
                          className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                          disabled={!sourceInstitution || blockchainLoading}
                        >
                          <option value="">Select subject...</option>
                          {sourceSubjects.map((subject) => (
                            <option key={subject.id} value={subject.id}>
                              {subject.code} - {subject.name} ({subject.credits} credits)
                            </option>
                          ))}
                        </select>
                        {sourceInstitution && sourceSubjects.length === 0 && (
                          <p className="text-sm text-gray-500 mt-1">No subjects found for this institution</p>
                        )}
                      </div>
                    </div>
                  </div>

                  {/* Target Institution & Subject */}
                  <div className="space-y-4">
                    <h3 className="font-semibold text-gray-700">🎯 Target Institution & Subject</h3>
                    <div className="space-y-3">
                      <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                          Institution ({institutions.length} available)
                        </label>
                        <select 
                          value={targetInstitution}
                          onChange={(e) => {
                            setTargetInstitution(e.target.value)
                            setTargetSubject('') // Reset subject when institution changes
                          }}
                          className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                          disabled={blockchainLoading || !connection.connected}
                        >
                          <option value="">Select target institution...</option>
                          {institutions.map((inst) => (
                            <option key={inst.id} value={inst.id}>
                              {inst.name} ({inst.code})
                            </option>
                          ))}
                        </select>
                      </div>
                      
                      <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                          Subject ({targetSubjects.length} available)
                        </label>
                        <select 
                          value={targetSubject}
                          onChange={(e) => setTargetSubject(e.target.value)}
                          className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                          disabled={!targetInstitution || blockchainLoading}
                        >
                          <option value="">Select subject...</option>
                          {targetSubjects.map((subject) => (
                            <option key={subject.id} value={subject.id}>
                              {subject.code} - {subject.name} ({subject.credits} credits)
                            </option>
                          ))}
                        </select>
                        {targetInstitution && targetSubjects.length === 0 && (
                          <p className="text-sm text-gray-500 mt-1">No subjects found for this institution</p>
                        )}
                      </div>
                    </div>
                  </div>
                </div>

                {/* AI Analysis */}
                {(analyzingEquivalence || analysis) && (
                  <div className="bg-blue-50 rounded-lg p-4">
                    <h4 className="font-semibold text-blue-800 mb-2">🤖 Automatic Analysis via Smart Contract</h4>
                    {analyzingEquivalence ? (
                      <div className="flex items-center space-x-2">
                        <div className="animate-spin h-4 w-4 border-2 border-blue-600 border-t-transparent rounded-full"></div>
                        <span className="text-sm text-blue-700">Analyzing content similarity...</span>
                      </div>
                    ) : analysis && (
                      <div className="space-y-2 text-sm">
                        <div className="flex justify-between">
                          <span>Content similarity:</span>
                          <span className={`font-semibold ${analysis.similarity >= 80 ? 'text-green-600' : 'text-yellow-600'}`}>
                            {analysis.similarity}%
                          </span>
                        </div>
                        <div className="flex justify-between">
                          <span>Contract recommendation:</span>
                          <span className={`font-semibold ${analysis.recommendation === 'approved' ? 'text-green-600' : 'text-red-600'}`}>
                            {analysis.recommendation === 'approved' ? '✅ Approval Recommended' : '❌ Rejection Recommended'}
                          </span>
                        </div>
                      </div>
                    )}
                  </div>
                )}

                <button
                  onClick={handleRequestEquivalence}
                  disabled={!sourceInstitution || !sourceSubject || !targetInstitution || !targetSubject || analyzingEquivalence || !connection.connected}
                  className="w-full bg-blue-500 text-white px-6 py-3 rounded-lg hover:bg-blue-600 transition-colors font-semibold disabled:bg-gray-300 disabled:cursor-not-allowed"
                >
                  {!connection.connected ? '❌ Blockchain Disconnected' : 
                   analyzingEquivalence ? '🔄 Analyzing...' : 
                   '🚀 Request Equivalence'}
                </button>
              </div>
            )}

            {/* Pending Equivalences */}
            {activeTab === 'pending' && (
              <div className="space-y-6">
                <h2 className="text-xl font-semibold text-gray-800">⏳ Pending Equivalences</h2>
                
                <div className="space-y-4">
                  {pendingEquivalences.map((equiv) => (
                    <div key={equiv.id} className="border border-gray-200 rounded-lg p-4 hover:shadow-md transition-shadow">
                      <div className="flex justify-between items-start mb-3">
                        <div>
                          <h3 className="font-semibold text-gray-800">{equiv.studentName}</h3>
                          <p className="text-sm text-gray-600">Equivalence request</p>
                        </div>
                        <span className="px-3 py-1 bg-yellow-100 text-yellow-700 rounded-full text-sm font-semibold">
                          ⏳ Pending
                        </span>
                      </div>

                      <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
                        <div className="bg-gray-50 p-3 rounded">
                          <p className="text-sm font-semibold text-gray-700">From:</p>
                          <p className="text-sm">{equiv.sourceSubject}</p>
                          <p className="text-xs text-gray-500">{equiv.sourceInstitution}</p>
                        </div>
                        <div className="bg-blue-50 p-3 rounded">
                          <p className="text-sm font-semibold text-blue-700">To:</p>
                          <p className="text-sm">{equiv.targetSubject}</p>
                          <p className="text-xs text-blue-500">{equiv.targetInstitution}</p>
                        </div>
                      </div>

                      <div className="flex items-center justify-between">
                        <div className="flex items-center space-x-2">
                          <span className="text-sm text-gray-600">Similarity:</span>
                          <span className="text-sm font-semibold text-green-600">{equiv.similarity}%</span>
                        </div>

                        <div className="flex space-x-2">
                          <button className="px-4 py-2 bg-red-500 text-white rounded hover:bg-red-600 transition-colors text-sm">
                            ❌ Reject
                          </button>
                          <button 
                            onClick={() => handleApprove(equiv.id)}
                            className="px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600 transition-colors text-sm"
                          >
                            ✅ Approve
                          </button>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            )}

            {/* Approved Equivalences */}
            {activeTab === 'approved' && (
              <div className="space-y-6">
                <h2 className="text-xl font-semibold text-gray-800">✅ Approved Equivalences</h2>
                
                <div className="space-y-4">
                  {approvedEquivalences.map((equiv) => (
                    <div key={equiv.id} className="border border-green-200 bg-green-50 rounded-lg p-4">
                      <div className="flex justify-between items-start mb-3">
                        <div>
                          <h3 className="font-semibold text-gray-800">{equiv.studentName}</h3>
                          <p className="text-sm text-gray-600">Approved equivalence</p>
                        </div>
                        <span className="px-3 py-1 bg-green-100 text-green-700 rounded-full text-sm font-semibold">
                          ✅ Approved
                        </span>
                      </div>

                      <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
                        <div className="bg-white p-3 rounded">
                          <p className="text-sm font-semibold text-gray-700">From:</p>
                          <p className="text-sm">{equiv.sourceSubject}</p>
                          <p className="text-xs text-gray-500">{equiv.sourceInstitution}</p>
                        </div>
                        <div className="bg-blue-50 p-3 rounded">
                          <p className="text-sm font-semibold text-blue-700">To:</p>
                          <p className="text-sm">{equiv.targetSubject}</p>
                          <p className="text-xs text-blue-500">{equiv.targetInstitution}</p>
                        </div>
                      </div>

                      <button className="w-full bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600 transition-colors text-sm">
                        🔗 View on Blockchain
                      </button>
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
