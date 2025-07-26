// Configuration file for the Academic Token Frontend

export const config = {
  // API Configuration
  api: {
    baseUrl: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:1317',
    timeout: 10000,
  },
  
  // Blockchain Configuration
  blockchain: {
    network: process.env.NEXT_PUBLIC_BLOCKCHAIN_NETWORK || 'academictoken',
    version: process.env.NEXT_PUBLIC_BLOCKCHAIN_VERSION || 'v1.0.0',
  },
  
  // Development Settings
  dev: {
    enableDebugLogs: process.env.NODE_ENV === 'development',
    mockData: false,
  }
};

// Log configuration on load
if (config.dev.enableDebugLogs) {
  console.log('🔧 Academic Token Configuration:', config);
}
