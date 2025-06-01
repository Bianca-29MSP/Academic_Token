import { useState, useEffect, useCallback } from 'react';
import { AcademicTokenAPI } from '../lib/api';
import type { 
  Institution, 
  Course, 
  Subject, 
  BlockchainConnection 
} from '../types/blockchain';

interface UseBlockchainReturn {
  // Connection state
  connection: BlockchainConnection;
  isLoading: boolean;
  error: string | null;
  
  // Data state
  institutions: Institution[];
  courses: Course[];
  subjects: Subject[];
  
  // Debug info
  debug: {
    lastUpdate: string;
    apiCalls: number;
    errors: string[];
  };
  
  // Methods
  refreshConnection: () => Promise<void>;
  refreshData: () => Promise<void>;
  getSubjectsByInstitution: (institutionId: string) => Subject[];
  getCoursesByInstitution: (institutionId: string) => Course[];
  clearErrors: () => void;
}

export function useBlockchain(): UseBlockchainReturn {
  // Connection state
  const [connection, setConnection] = useState<BlockchainConnection>({
    connected: false,
    nodeUrl: 'http://localhost:1318',
    network: 'disconnected'
  });
  
  // Loading and error state
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  
  // Data state
  const [institutions, setInstitutions] = useState<Institution[]>([]);
  const [courses, setCourses] = useState<Course[]>([]);
  const [subjects, setSubjects] = useState<Subject[]>([]);
  
  // Debug state
  const [debug, setDebug] = useState({
    lastUpdate: 'Never',
    apiCalls: 0,
    errors: [] as string[]
  });

  // Clear errors
  const clearErrors = useCallback(() => {
    setError(null);
    setDebug(prev => ({ ...prev, errors: [] }));
  }, []);

  // Update debug info
  const updateDebug = useCallback((action: string, success: boolean, errorMsg?: string) => {
    setDebug(prev => ({
      lastUpdate: new Date().toLocaleTimeString(),
      apiCalls: prev.apiCalls + 1,
      errors: errorMsg ? [...prev.errors.slice(-4), `${action}: ${errorMsg}`] : prev.errors
    }));
    
    console.log(`🔧 DEBUG [${action}]:`, success ? '✅ Success' : '❌ Failed', errorMsg || '');
  }, []);

  // Check blockchain connection
  const refreshConnection = useCallback(async () => {
    console.log('🔄 Refreshing blockchain connection...');
    
    try {
      setError(null);
      updateDebug('Connection Check', true);
      
      const connectionStatus = await AcademicTokenAPI.checkConnection();
      setConnection(connectionStatus);
      
      console.log('📡 Connection Status:', connectionStatus);
      
      if (!connectionStatus.connected) {
        const errorMsg = 'Unable to connect to blockchain node at ' + connectionStatus.nodeUrl;
        setError(errorMsg);
        updateDebug('Connection Check', false, errorMsg);
      } else {
        updateDebug('Connection Check', true);
        console.log('✅ Successfully connected to blockchain');
      }
    } catch (err) {
      const errorMsg = `Connection failed: ${err instanceof Error ? err.message : 'Unknown error'}`;
      console.error('❌ Connection refresh failed:', err);
      setError('Connection failed - blockchain node may be offline');
      updateDebug('Connection Check', false, errorMsg);
      setConnection(prev => ({ ...prev, connected: false }));
    }
  }, [updateDebug]);

  // Load all blockchain data
  const refreshData = useCallback(async () => {
    if (!connection.connected) {
      console.log('⏭️ Skipping data refresh - not connected to blockchain');
      return;
    }

    console.log('🔄 Loading blockchain data...');
    setIsLoading(true);
    
    try {
      setError(null);
      
      console.log('📊 Fetching institutions...');
      const institutionsData = await AcademicTokenAPI.getInstitutions();
      console.log('🏛️ Institutions loaded:', institutionsData.length);
      setInstitutions(institutionsData);
      updateDebug('Load Institutions', true);

      console.log('📊 Fetching courses...');
      const coursesData = await AcademicTokenAPI.getCourses();
      console.log('📚 Courses loaded:', coursesData.length);
      setCourses(coursesData);
      updateDebug('Load Courses', true);

      console.log('📊 Fetching subjects...');
      const subjectsData = await AcademicTokenAPI.getSubjects();
      console.log('📖 Subjects loaded:', subjectsData.length);
      setSubjects(subjectsData);
      updateDebug('Load Subjects', true);
      
      console.log(`✅ Data loading complete: ${institutionsData.length} institutions, ${coursesData.length} courses, ${subjectsData.length} subjects`);
      
      // Show helpful message if blockchain is empty
      if (institutionsData.length === 0 && coursesData.length === 0 && subjectsData.length === 0) {
        console.log('💡 Blockchain appears to be empty - this is normal for a new installation');
      }
      
    } catch (err) {
      const errorMsg = `Failed to load data: ${err instanceof Error ? err.message : 'Unknown error'}`;
      console.error('❌ Data refresh failed:', err);
      setError('Failed to load data from blockchain');
      updateDebug('Load Data', false, errorMsg);
    } finally {
      setIsLoading(false);
    }
  }, [connection.connected, updateDebug]);

  // Helper methods
  const getSubjectsByInstitution = useCallback((institutionId: string): Subject[] => {
    const filtered = subjects.filter(subject => subject.institutionId === institutionId);
    console.log(`🔍 Filtering subjects for institution ${institutionId}:`, filtered);
    return filtered;
  }, [subjects]);

  const getCoursesByInstitution = useCallback((institutionId: string): Course[] => {
    const filtered = courses.filter(course => course.institutionId === institutionId);
    console.log(`🔍 Filtering courses for institution ${institutionId}:`, filtered);
    return filtered;
  }, [courses]);

  // Initialize connection and data on mount
  useEffect(() => {
    console.log('🚀 Initializing blockchain hook...');
    refreshConnection();
  }, [refreshConnection]);

  // Load data when connection is established
  useEffect(() => {
    if (connection.connected) {
      console.log('✅ Connection established, loading data...');
      refreshData();
    } else {
      console.log('❌ Not connected, skipping data load');
      setIsLoading(false);
    }
  }, [connection.connected, refreshData]);

  // Debug log current state
  useEffect(() => {
    console.log('📊 Current State:', {
      connected: connection.connected,
      loading: isLoading,
      error: error,
      institutionCount: institutions.length,
      courseCount: courses.length,
      subjectCount: subjects.length
    });
  }, [connection.connected, isLoading, error, institutions.length, courses.length, subjects.length]);

  return {
    // Connection state
    connection,
    isLoading,
    error,
    
    // Data
    institutions,
    courses,
    subjects,
    
    // Debug
    debug,
    
    // Methods
    refreshConnection,
    refreshData,
    getSubjectsByInstitution,
    getCoursesByInstitution,
    clearErrors,
  };
}
