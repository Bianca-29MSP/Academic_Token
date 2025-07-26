// Component to help debug React key warnings

export function DebugKeys() {
  // This component can be imported to help debug key warnings
  const debugArrays = [
    { id: 1, name: 'Item 1' },
    { id: 2, name: 'Item 2' },
    { id: 3, name: 'Item 3' }
  ];

  return (
    <div className="p-4 bg-gray-100 rounded-lg">
      <h3 className="font-bold mb-2">Debug Key Component</h3>
      <p className="text-sm text-gray-600 mb-4">
        This component helps debug React key warnings. 
        If you see key warnings in the console, check that all .map() calls have unique key props.
      </p>
      
      <div className="space-y-2">
        {/* Correct usage with key */}
        {debugArrays.map((item) => (
          <div key={item.id} className="p-2 bg-white rounded">
            ✅ Correct: {item.name} (key={item.id})
          </div>
        ))}
      </div>
      
      <div className="mt-4 p-3 bg-yellow-50 rounded text-sm">
        <strong>Common causes of key warnings:</strong>
        <ul className="mt-2 space-y-1 text-xs">
          <li>• Using array index as key when items can be reordered</li>
          <li>• Missing key prop entirely in .map() calls</li>
          <li>• Non-unique keys (duplicate values)</li>
          <li>• Keys on wrong element (should be on direct child of .map())</li>
        </ul>
      </div>
    </div>
  );
}

// Helper function to generate unique keys
export function generateKey(prefix: string, index: number): string {
  return `${prefix}-${index}-${Date.now()}`;
}

// Helper to check if an array has unique IDs
export function hasUniqueIds<T extends { id?: any }>(array: T[]): boolean {
  const ids = array.map(item => item.id).filter(id => id !== undefined);
  return ids.length === new Set(ids).size;
}
