#!/usr/bin/env node

// Script to check blockchain data structure
const axios = require('axios');

const API_URL = 'http://localhost:1317';

async function checkDataStructure() {
  console.log('🔍 Checking Academic Token Data Structure...\n');
  
  try {
    // Check institutions
    console.log('📊 Fetching Institutions...');
    const instResponse = await axios.get(`${API_URL}/academictoken/institution/institution`);
    const institutions = instResponse.data.institution || [];
    console.log(`Found ${institutions.length} institutions`);
    if (institutions.length > 0) {
      console.log('First institution structure:');
      console.log(JSON.stringify(institutions[0], null, 2));
    }
    
    // Check courses
    console.log('\n📊 Fetching Courses...');
    const courseResponse = await axios.get(`${API_URL}/academictoken/course/course`);
    const courses = courseResponse.data.course || [];
    console.log(`Found ${courses.length} courses`);
    if (courses.length > 0) {
      console.log('First course structure:');
      console.log(JSON.stringify(courses[0], null, 2));
    }
    
    // Check subjects
    console.log('\n📊 Fetching Subjects...');
    const subjectResponse = await axios.get(`${API_URL}/academictoken/subject/subjects`);
    const subjects = subjectResponse.data.subjects || [];
    console.log(`Found ${subjects.length} subjects`);
    if (subjects.length > 0) {
      console.log('First subject structure:');
      console.log(JSON.stringify(subjects[0], null, 2));
      
      // Check if subjects have institutionId
      console.log('\n🔍 Checking subject fields:');
      const firstSubject = subjects[0];
      console.log('Has institutionId?', 'institutionId' in firstSubject);
      console.log('Has courseId?', 'courseId' in firstSubject);
      console.log('Has course_id?', 'course_id' in firstSubject);
      console.log('Has institution_id?', 'institution_id' in firstSubject);
      
      // List all fields
      console.log('\nAll fields in first subject:');
      console.log(Object.keys(firstSubject));
    }
    
    // Check subjects by institution endpoint
    if (institutions.length > 0) {
      const firstInstId = institutions[0].id;
      console.log(`\n📊 Testing subjects by institution endpoint for ${firstInstId}...`);
      try {
        const instSubjectsResponse = await axios.get(`${API_URL}/academictoken/subject/institutions/${firstInstId}/subjects`);
        const instSubjects = instSubjectsResponse.data.subjects || [];
        console.log(`Found ${instSubjects.length} subjects for institution ${firstInstId}`);
      } catch (error) {
        console.log('❌ Endpoint failed:', error.message);
        console.log('   This endpoint might not be implemented yet');
      }
    }
    
  } catch (error) {
    console.error('❌ Error:', error.message);
  }
}

checkDataStructure();
