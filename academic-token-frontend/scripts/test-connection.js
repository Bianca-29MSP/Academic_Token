#!/usr/bin/env node

// Test script to verify blockchain connection

const axios = require('axios');

const API_URL = 'http://localhost:1317';

async function testConnection() {
  console.log('🔍 Testing Academic Token blockchain connection...\n');
  console.log('🌐 API URL:', API_URL);
  console.log('─'.repeat(50));

  const endpoints = [
    {
      name: 'Base URL',
      url: '/',
      description: 'Check if API is running'
    },
    {
      name: 'OpenAPI Spec',
      url: '/static/openapi.yml',
      description: 'Check OpenAPI documentation'
    },
    {
      name: 'Institutions',
      url: '/academictoken/institution/institution',
      description: 'List all institutions'
    },
    {
      name: 'Courses', 
      url: '/academictoken/course/course',
      description: 'List all courses'
    },
    {
      name: 'Subjects',
      url: '/academictoken/subject/subjects',
      description: 'List all subjects'
    },
    {
      name: 'Node Info',
      url: '/cosmos/base/tendermint/v1beta1/node_info',
      description: 'Cosmos node information'
    }
  ];

  for (const endpoint of endpoints) {
    console.log(`\n📌 Testing: ${endpoint.name}`);
    console.log(`   URL: ${API_URL}${endpoint.url}`);
    console.log(`   ${endpoint.description}`);
    
    try {
      const response = await axios.get(`${API_URL}${endpoint.url}`, {
        timeout: 5000,
        validateStatus: () => true // Accept any status
      });
      
      if (response.status === 200) {
        console.log(`   ✅ Status: ${response.status} OK`);
        
        // Show data preview
        if (response.data) {
          const preview = JSON.stringify(response.data, null, 2);
          const lines = preview.split('\n');
          const maxLines = 10;
          
          if (lines.length > maxLines) {
            console.log(`   📊 Response preview (first ${maxLines} lines):`);
            console.log('   ' + lines.slice(0, maxLines).join('\n   '));
            console.log(`   ... (${lines.length - maxLines} more lines)`);
          } else {
            console.log('   📊 Response:');
            console.log('   ' + lines.join('\n   '));
          }
        }
      } else {
        console.log(`   ⚠️  Status: ${response.status} ${response.statusText || ''}`);
        if (response.data) {
          console.log(`   Message: ${JSON.stringify(response.data)}`);
        }
      }
    } catch (error) {
      console.log(`   ❌ Error: ${error.message}`);
      if (error.code === 'ECONNREFUSED') {
        console.log('   💡 Make sure the blockchain node is running on port 1317');
      }
    }
  }

  console.log('\n' + '─'.repeat(50));
  console.log('✨ Test complete!\n');
}

// Run the test
testConnection().catch(error => {
  console.error('Fatal error:', error);
  process.exit(1);
});
