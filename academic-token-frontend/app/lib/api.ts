import axios from 'axios';
import type { 
  Institution, 
  Course, 
  Subject, 
  Student, 
  AcademicNFT, 
  DegreeEligibility,
  BlockchainConnection 
} from '../types/blockchain';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:1318';

// Create axios instance with default config
const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor for logging
api.interceptors.request.use(
  (config) => {
    console.log(`🔗 API Request: ${config.method?.toUpperCase()} ${config.url}`);
    return config;
  },
  (error) => {
    console.error('❌ Request Error:', error);
    return Promise.reject(error);
  }
);

// Response interceptor for error handling
api.interceptors.response.use(
  (response) => {
    console.log(`✅ API Response: ${response.status} ${response.config.url}`);
    return response;
  },
  (error) => {
    console.error('❌ Response Error:', error.response?.data || error.message);
    return Promise.reject(error);
  }
);

export class AcademicTokenAPI {
  // Connection methods
  static async checkConnection(): Promise<BlockchainConnection> {
    try {
      const response = await api.get('/health');
      const nodeInfo = await api.get('/cosmos/base/tendermint/v1beta1/node_info');
      
      return {
        connected: response.status === 200,
        nodeUrl: API_BASE_URL,
        network: nodeInfo.data.default_node_info?.network || 'academictoken',
        version: nodeInfo.data.application_version?.version || 'v1.0.0'
      };
    } catch (error) {
      console.error('Connection check failed:', error);
      return {
        connected: false,
        nodeUrl: API_BASE_URL,
        network: 'disconnected'
      };
    }
  }

  // Institution methods
  static async getInstitutions(): Promise<Institution[]> {
    try {
      const response = await api.get('/academic/institution/list');
      return response.data || [];
    } catch (error) {
      console.log('🔍 No institutions found on blockchain (this is normal for empty blockchain)');
      return []; // Return empty array instead of throwing error
    }
  }

  static async getInstitutionById(id: string): Promise<Institution | null> {
    try {
      const institutions = await this.getInstitutions();
      return institutions.find(inst => inst.id === id) || null;
    } catch (error) {
      console.error('Failed to fetch institution:', error);
      return null;
    }
  }

  // Course methods
  static async getCourses(): Promise<Course[]> {
    try {
      const response = await api.get('/academic/course/list');
      return response.data || [];
    } catch (error) {
      console.log('🔍 No courses found on blockchain (this is normal for empty blockchain)');
      return []; // Return empty array instead of throwing error
    }
  }

  static async getCoursesByInstitution(institutionId: string): Promise<Course[]> {
    try {
      const courses = await this.getCourses();
      return courses.filter(course => course.institutionId === institutionId);
    } catch (error) {
      console.error('Failed to fetch courses by institution:', error);
      return [];
    }
  }

  // Subject methods
  static async getSubjects(): Promise<Subject[]> {
    try {
      const response = await api.get('/academic/subject/list');
      return response.data || [];
    } catch (error) {
      console.log('🔍 No subjects found on blockchain (this is normal for empty blockchain)');
      return []; // Return empty array instead of throwing error
    }
  }

  static async getSubjectsByInstitution(institutionId: string): Promise<Subject[]> {
    try {
      const subjects = await this.getSubjects();
      return subjects.filter(subject => subject.institutionId === institutionId);
    } catch (error) {
      console.error('Failed to fetch subjects by institution:', error);
      return [];
    }
  }

  static async getSubjectsByCourse(courseId: string): Promise<Subject[]> {
    try {
      const subjects = await this.getSubjects();
      return subjects.filter(subject => subject.courseId === courseId);
    } catch (error) {
      console.error('Failed to fetch subjects by course:', error);
      return [];
    }
  }

  // Student methods
  static async getStudent(studentId: string): Promise<Student | null> {
    try {
      const response = await api.get(`/academic/student/${studentId}`);
      return response.data;
    } catch (error) {
      console.error('Failed to fetch student:', error);
      return null;
    }
  }

  static async getStudentNFTs(studentId: string): Promise<AcademicNFT[]> {
    try {
      const response = await api.get(`/academic/student/${studentId}/nfts`);
      return response.data || [];
    } catch (error) {
      console.error('Failed to fetch student NFTs:', error);
      return [];
    }
  }

  // Degree methods
  static async getDegreeEligibility(studentId: string): Promise<DegreeEligibility> {
    try {
      const response = await api.get(`/academic/degree/${studentId}/eligibility`);
      return response.data;
    } catch (error) {
      console.error('Failed to fetch degree eligibility:', error);
      throw new Error('Failed to check degree eligibility');
    }
  }

  // Prerequisites methods
  static async checkPrerequisites(studentId: string, subjectId: string): Promise<boolean> {
    try {
      const response = await api.get(`/academic/student/${studentId}/prerequisites/${subjectId}`);
      return response.data.eligible || false;
    } catch (error) {
      console.error('Failed to check prerequisites:', error);
      throw new Error('Prerequisites check failed');
    }
  }

  // Utility methods
  // Note: Real equivalence analysis should be done via CosmWasm contracts
}

export default AcademicTokenAPI;
