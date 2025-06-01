// context/BlockchainContext.tsx
'use client';
import { createContext, useContext, useReducer, ReactNode, useEffect, useState } from 'react';
import { 
  Institution, 
  Course, 
  Subject, 
  Student, 
  AcademicNFT,
  EquivalenceRequest 
} from '../services/blockchain';

// Types for global state
interface BlockchainState {
  // Connection
  connected: boolean;
  chainId: string | null;
  
  // Current user
  currentUser: {
    type: 'student' | 'institution' | null;
    id: string | null;
    data: Student | Institution | null;
  };
  
  // Data cache
  institutions: Institution[];
  courses: Course[];
  subjects: Subject[];
  studentNFTs: AcademicNFT[];
  equivalenceRequests: EquivalenceRequest[];
  
  // Loading states
  loading: {
    institutions: boolean;
    courses: boolean;
    subjects: boolean;
    nfts: boolean;
    equivalences: boolean;
  };
  
  // Errors
  errors: {
    connection: string | null;
    general: string | null;
  };
  
  // Configuration
  config: {
    apiUrl: string;
    rpcUrl: string;
    chainId: string;
  };
}

// Actions for the reducer
type BlockchainAction = 
  | { type: 'SET_CONNECTION'; payload: { connected: boolean; chainId: string | null } }
  | { type: 'SET_CURRENT_USER'; payload: { type: 'student' | 'institution' | null; id: string | null; data: any } }
  | { type: 'SET_INSTITUTIONS'; payload: Institution[] }
  | { type: 'ADD_INSTITUTION'; payload: Institution }
  | { type: 'SET_COURSES'; payload: Course[] }
  | { type: 'ADD_COURSE'; payload: Course }
  | { type: 'SET_SUBJECTS'; payload: Subject[] }
  | { type: 'ADD_SUBJECT'; payload: Subject }
  | { type: 'SET_STUDENT_NFTS'; payload: AcademicNFT[] }
  | { type: 'ADD_NFT'; payload: AcademicNFT }
  | { type: 'SET_EQUIVALENCE_REQUESTS'; payload: EquivalenceRequest[] }
  | { type: 'ADD_EQUIVALENCE_REQUEST'; payload: EquivalenceRequest }
  | { type: 'SET_LOADING'; payload: { key: keyof BlockchainState['loading']; value: boolean } }
  | { type: 'SET_ERROR'; payload: { key: keyof BlockchainState['errors']; value: string | null } }
  | { type: 'CLEAR_ERRORS' }
  | { type: 'RESET_STATE' };

// Initial state (configuration will be defined after hydration)
const createInitialState = (): BlockchainState => ({
  connected: false,
  chainId: null,
  currentUser: {
    type: null,
    id: null,
    data: null,
  },
  institutions: [],
  courses: [],
  subjects: [],
  studentNFTs: [],
  equivalenceRequests: [],
  loading: {
    institutions: false,
    courses: false,
    subjects: false,
    nfts: false,
    equivalences: false,
  },
  errors: {
    connection: null,
    general: null,
  },
  config: {
    apiUrl: 'http://localhost:1317', // Default value
    rpcUrl: 'http://localhost:26657', // Default value
    chainId: 'academictoken', // Default value
  },
});

// Reducer
function blockchainReducer(state: BlockchainState, action: BlockchainAction): BlockchainState {
  switch (action.type) {
    case 'SET_CONNECTION':
      return {
        ...state,
        connected: action.payload.connected,
        chainId: action.payload.chainId,
        errors: {
          ...state.errors,
          connection: action.payload.connected ? null : state.errors.connection,
        },
      };
      
    case 'SET_CURRENT_USER':
      return {
        ...state,
        currentUser: action.payload,
      };
      
    case 'SET_INSTITUTIONS':
      return {
        ...state,
        institutions: action.payload,
      };
      
    case 'ADD_INSTITUTION':
      return {
        ...state,
        institutions: [...state.institutions, action.payload],
      };
      
    case 'SET_COURSES':
      return {
        ...state,
        courses: action.payload,
      };
      
    case 'ADD_COURSE':
      return {
        ...state,
        courses: [...state.courses, action.payload],
      };
      
    case 'SET_SUBJECTS':
      return {
        ...state,
        subjects: action.payload,
      };
      
    case 'ADD_SUBJECT':
      return {
        ...state,
        subjects: [...state.subjects, action.payload],
      };
      
    case 'SET_STUDENT_NFTS':
      return {
        ...state,
        studentNFTs: action.payload,
      };
      
    case 'ADD_NFT':
      return {
        ...state,
        studentNFTs: [...state.studentNFTs, action.payload],
      };
      
    case 'SET_EQUIVALENCE_REQUESTS':
      return {
        ...state,
        equivalenceRequests: action.payload,
      };
      
    case 'ADD_EQUIVALENCE_REQUEST':
      return {
        ...state,
        equivalenceRequests: [...state.equivalenceRequests, action.payload],
      };
      
    case 'SET_LOADING':
      return {
        ...state,
        loading: {
          ...state.loading,
          [action.payload.key]: action.payload.value,
        },
      };
      
    case 'SET_ERROR':
      return {
        ...state,
        errors: {
          ...state.errors,
          [action.payload.key]: action.payload.value,
        },
      };
      
    case 'CLEAR_ERRORS':
      return {
        ...state,
        errors: {
          connection: null,
          general: null,
        },
      };
      
    case 'RESET_STATE':
      return {
        ...createInitialState(),
        config: state.config, // Keep configuration
      };
      
    default:
      return state;
  }
}

// Context
const BlockchainContext = createContext<{
  state: BlockchainState;
  dispatch: React.Dispatch<BlockchainAction>;
} | null>(null);

// Provider component
interface BlockchainProviderProps {
  children: ReactNode;
}

export function BlockchainProvider({ children }: BlockchainProviderProps) {
  const [state, dispatch] = useReducer(blockchainReducer, createInitialState());
  const [isClient, setIsClient] = useState(false);

  // Wait for hydration to avoid SSR mismatch
  useEffect(() => {
    setIsClient(true);
    
    // Configure URLs from environment after hydration
    const config = {
      apiUrl: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:1317',
      rpcUrl: process.env.NEXT_PUBLIC_COSMOS_RPC || 'http://localhost:26657',
      chainId: process.env.NEXT_PUBLIC_CHAIN_ID || 'academictoken',
    };
    
    dispatch({
      type: 'SET_CONNECTION',
      payload: { connected: false, chainId: config.chainId }
    });
  }, []);

  // Check connection only after hydration
  useEffect(() => {
    if (!isClient) return;
    
    const checkConnection = async () => {
      const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:1317';
      console.log('🔍 Testing connection to:', apiUrl);
      
      try {
        const response = await fetch(`${apiUrl}/cosmos/base/tendermint/v1beta1/node_info`, {
          method: 'GET',
          signal: AbortSignal.timeout(5000)
        });
        
        if (response.ok) {
          const data = await response.json();
          const networkId = data.default_node_info?.network || 'academictoken';
          console.log('✅ Connected to blockchain:', networkId);
          
          dispatch({
            type: 'SET_CONNECTION',
            payload: {
              connected: true,
              chainId: networkId,
            },
          });
        } else {
          throw new Error(`API responded with status: ${response.status}`);
        }
      } catch (error) {
        console.warn('❌ Connection error, using demo mode:', error);
        dispatch({
          type: 'SET_CONNECTION',
          payload: { connected: false, chainId: null },
        });
        dispatch({
          type: 'SET_ERROR',
          payload: {
            key: 'connection',
            value: 'Demo mode - backend not connected',
          },
        });
      }
    };

    checkConnection();
    
    // Recheck connection every 60 seconds
    const interval = setInterval(checkConnection, 60000);
    
    return () => clearInterval(interval);
  }, [isClient]);

  return (
    <BlockchainContext.Provider value={{ state, dispatch }}>
      {children}
    </BlockchainContext.Provider>
  );
}

// Hook to use the context
export function useBlockchainContext() {
  const context = useContext(BlockchainContext);
  if (!context) {
    throw new Error('useBlockchainContext must be used within a BlockchainProvider');
  }
  return context;
}

// Higher-order component for pages that require connection
interface WithBlockchainConnectionProps {
  children: ReactNode;
  fallback?: ReactNode;
}

export function WithBlockchainConnection({ children, fallback }: WithBlockchainConnectionProps) {
  const { state } = useBlockchainContext();

  // If not connected, but still is demo, let it continue
  if (!state.connected && state.errors.connection !== 'Demo mode - backend not connected') {
    return (
      fallback || (
        <div className="min-h-screen flex items-center justify-center bg-gray-50">
          <div className="text-center">
            <div className="text-6xl mb-4">🔗</div>
            <h2 className="text-2xl font-bold text-gray-800 mb-2">Connecting to Blockchain</h2>
            <p className="text-gray-600 mb-4">Waiting for connection to Cosmos network...</p>
            {state.errors.connection && (
              <div className="bg-red-50 border border-red-200 rounded-lg p-4 text-red-700 max-w-md mx-auto">
                {state.errors.connection}
                <div className="mt-2 text-sm text-red-600">
                  The frontend will continue working in demo mode
                </div>
              </div>
            )}
            <div className="mt-4">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
            </div>
          </div>
        </div>
      )
    );
  }

  return <>{children}</>;
}

// Utility hooks based on context
export function useCurrentUser() {
  const { state, dispatch } = useBlockchainContext();
  
  const setCurrentUser = (type: 'student' | 'institution' | null, id: string | null, data: any = null) => {
    dispatch({
      type: 'SET_CURRENT_USER',
      payload: { type, id, data },
    });
  };
  
  const logout = () => {
    dispatch({
      type: 'SET_CURRENT_USER',
      payload: { type: null, id: null, data: null },
    });
  };
  
  return {
    currentUser: state.currentUser,
    setCurrentUser,
    logout,
    isLoggedIn: state.currentUser.type !== null,
    isStudent: state.currentUser.type === 'student',
    isInstitution: state.currentUser.type === 'institution',
  };
}

export function useConnectionStatus() {
  const { state, dispatch } = useBlockchainContext();
  
  const reconnect = async () => {
    dispatch({ type: 'CLEAR_ERRORS' });
    
    try {
      const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:1317';
      const response = await fetch(`${apiUrl}/cosmos/base/tendermint/v1beta1/node_info`);
      
      if (response.ok) {
        const data = await response.json();
        dispatch({
          type: 'SET_CONNECTION',
          payload: {
            connected: true,
            chainId: data.default_node_info?.network || 'academictoken',
          },
        });
      } else {
        throw new Error('Failed to connect');
      }
    } catch (error) {
      dispatch({
        type: 'SET_ERROR',
        payload: {
          key: 'connection',
          value: error instanceof Error ? error.message : 'Connection error',
        },
      });
    }
  };
  
  return {
    connected: state.connected,
    chainId: state.chainId,
    error: state.errors.connection,
    reconnect,
  };
}

export function useGlobalLoading() {
  const { state, dispatch } = useBlockchainContext();
  
  const setLoading = (key: keyof BlockchainState['loading'], value: boolean) => {
    dispatch({
      type: 'SET_LOADING',
      payload: { key, value },
    });
  };
  
  const isAnyLoading = Object.values(state.loading).some(loading => loading);
  
  return {
    loading: state.loading,
    isAnyLoading,
    setLoading,
  };
}

export function useGlobalError() {
  const { state, dispatch } = useBlockchainContext();
  
  const setError = (key: keyof BlockchainState['errors'], value: string | null) => {
    dispatch({
      type: 'SET_ERROR',
      payload: { key, value },
    });
  };
  
  const clearErrors = () => {
    dispatch({ type: 'CLEAR_ERRORS' });
  };
  
  const hasErrors = Object.values(state.errors).some(error => error !== null);
  
  return {
    errors: state.errors,
    hasErrors,
    setError,
    clearErrors,
  };
}