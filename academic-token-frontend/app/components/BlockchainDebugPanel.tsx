// Debug component to inspect blockchain data structure
import { useBlockchain } from '../hooks/useBlockchain';

export function BlockchainDebugPanel() {
  const { 
    institutions, 
    courses, 
    subjects,
    connection 
  } = useBlockchain();

  if (!connection.connected) {
    return null;
  }

  return (
    <div className="fixed bottom-4 right-4 w-96 max-h-96 overflow-auto bg-black bg-opacity-90 text-green-400 p-4 rounded-lg font-mono text-xs z-50">
      <h3 className="text-yellow-400 mb-2">🔍 Blockchain Debug</h3>
      
      <div className="mb-2">
        <strong className="text-cyan-400">Institutions ({institutions.length}):</strong>
        {institutions.length > 0 && (
          <pre className="text-gray-300">{JSON.stringify(institutions[0], null, 2)}</pre>
        )}
      </div>
      
      <div className="mb-2">
        <strong className="text-cyan-400">Courses ({courses.length}):</strong>
        {courses.length > 0 && (
          <pre className="text-gray-300">{JSON.stringify(courses[0], null, 2)}</pre>
        )}
      </div>
      
      <div className="mb-2">
        <strong className="text-cyan-400">Subjects ({subjects.length}):</strong>
        {subjects.length > 0 && (
          <pre className="text-gray-300">{JSON.stringify(subjects[0], null, 2)}</pre>
        )}
      </div>
      
      <button 
        onClick={() => {
          console.log('📊 Full Data Dump:');
          console.log('Institutions:', institutions);
          console.log('Courses:', courses);
          console.log('Subjects:', subjects);
        }}
        className="mt-2 px-3 py-1 bg-blue-600 text-white rounded hover:bg-blue-700"
      >
        Log Full Data
      </button>
    </div>
  );
}
