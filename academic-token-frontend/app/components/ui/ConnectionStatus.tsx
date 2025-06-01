// components/ui/ConnectionStatus.tsx
'use client';
import { useBlockchainContext } from '../../context/BlockchainContext';
import { useBlockchainConnection } from '../../hooks/useBlockchain';

export function ConnectionStatus() {
  const { state } = useBlockchainContext();
  const { connected, chainId, error } = useBlockchainConnection();

  if (connected) {
    return (
      <div className="text-sm text-green-600 mt-1 flex items-center space-x-2">
        <span className="w-2 h-2 bg-green-500 rounded-full animate-pulse"></span>
        <span>🔗 Connected to {chainId}</span>
      </div>
    );
  }

  return (
    <div className="text-sm text-orange-600 mt-1 flex items-center space-x-2">
      <span className="w-2 h-2 bg-orange-500 rounded-full"></span>
      <span>⚠️ Demo Mode - {error?.message || 'Backend not connected'}</span>
    </div>
  );
}

export function ConnectionBanner() {
  const { state } = useBlockchainContext();
  const { connected, error, reconnect } = useBlockchainConnection();

  if (connected) {
    return (
      <div className="bg-green-50 border border-green-200 rounded-lg p-3 mb-4">
        <div className="flex items-center justify-between">
          <div className="flex items-start">
            <div className="flex-shrink-0">
              <span className="text-green-500">✅</span>
            </div>
            <div className="ml-3">
              <p className="text-sm text-green-700">
                <strong>Blockchain Connected:</strong> All features are fully functional. 
                Transactions will be submitted to the real blockchain.
              </p>
              <p className="text-xs text-green-600 mt-1">
                Chain ID: {state.chainId} | API: {state.config.apiUrl}
              </p>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="bg-orange-50 border border-orange-200 rounded-lg p-3 mb-4">
      <div className="flex items-center justify-between">
        <div className="flex items-start">
          <div className="flex-shrink-0">
            <span className="text-orange-500">⚠️</span>
          </div>
          <div className="ml-3">
            <p className="text-sm text-orange-700">
              <strong>Demo Mode Active:</strong> The blockchain backend is not connected. 
              All functionality is simulated for demonstration.
            </p>
            <p className="text-xs text-orange-600 mt-1">
              To connect with real blockchain, start your Cosmos backend on localhost:1317
            </p>
            {error && (
              <p className="text-xs text-orange-500 mt-1">
                Error: {error.message}
              </p>
            )}
          </div>
        </div>
        <button 
          onClick={reconnect}
          className="ml-4 px-3 py-1 bg-orange-200 text-orange-800 rounded hover:bg-orange-300 transition-colors text-xs"
        >
          Retry Connection
        </button>
      </div>
    </div>
  );
}

// Advanced connection dashboard
export function ConnectionDashboard() {
  const { state } = useBlockchainContext();
  const { connected, chainId, loading, error, reconnect } = useBlockchainConnection();

  return (
    <div className="bg-white rounded-lg shadow-lg p-6">
      <h3 className="text-lg font-semibold text-gray-800 mb-4">🔗 Blockchain Connection</h3>
      
      <div className="space-y-4">
        {/* Connection Status */}
        <div className="flex items-center justify-between p-3 bg-gray-50 rounded">
          <div className="flex items-center space-x-3">
            <div className={`w-3 h-3 rounded-full ${
              loading ? 'bg-yellow-500 animate-pulse' :
              connected ? 'bg-green-500' : 'bg-red-500'
            }`}></div>
            <div>
              <p className="font-semibold text-gray-800">
                {loading ? 'Connecting...' : connected ? 'Connected' : 'Disconnected'}
              </p>
              <p className="text-sm text-gray-600">
                {connected ? `Chain: ${chainId}` : 'Using demo mode'}
              </p>
            </div>
          </div>
          {!connected && !loading && (
            <button 
              onClick={reconnect}
              className="px-3 py-1 bg-blue-500 text-white rounded hover:bg-blue-600 transition-colors text-sm"
            >
              Reconnect
            </button>
          )}
        </div>

        {/* Configuration */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
          <div>
            <p className="font-semibold text-gray-700">API URL:</p>
            <p className="text-gray-600 font-mono break-all">{state.config.apiUrl}</p>
          </div>
          <div>
            <p className="font-semibold text-gray-700">RPC URL:</p>
            <p className="text-gray-600 font-mono break-all">{state.config.rpcUrl}</p>
          </div>
          <div>
            <p className="font-semibold text-gray-700">Chain ID:</p>
            <p className="text-gray-600 font-mono">{state.config.chainId}</p>
          </div>
          <div>
            <p className="font-semibold text-gray-700">Denomination:</p>
            <p className="text-gray-600 font-mono">utoken</p>
          </div>
        </div>

        {/* Error Details */}
        {error && (
          <div className="p-3 bg-red-50 border border-red-200 rounded">
            <p className="font-semibold text-red-800 text-sm">Connection Error:</p>
            <p className="text-red-700 text-xs mt-1">{error.message}</p>
            {process.env.NEXT_PUBLIC_DEBUG === 'true' && error.details && (
              <details className="mt-2">
                <summary className="text-xs text-red-600 cursor-pointer">Debug Details</summary>
                <pre className="text-xs text-red-500 mt-1 whitespace-pre-wrap">
                  {JSON.stringify(error.details, null, 2)}
                </pre>
              </details>
            )}
          </div>
        )}

        {/* Features Status */}
        <div className="pt-4 border-t border-gray-200">
          <p className="font-semibold text-gray-700 mb-2">Available Features:</p>
          <div className="grid grid-cols-2 gap-2 text-xs">
            {[
              { name: 'Student Enrollment', available: true },
              { name: 'Subject Completion', available: true },
              { name: 'NFT Minting', available: connected },
              { name: 'Equivalence Analysis', available: connected },
              { name: 'Degree Validation', available: connected },
              { name: 'Real Transactions', available: connected }
            ].map((feature, index) => (
              <div key={index} className="flex items-center space-x-2">
                <span className={`w-2 h-2 rounded-full ${
                  feature.available ? 'bg-green-500' : 'bg-gray-400'
                }`}></span>
                <span className={feature.available ? 'text-gray-700' : 'text-gray-500'}>
                  {feature.name}
                </span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}