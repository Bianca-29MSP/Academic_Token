#!/usr/bin/env node

// Script to find React components with missing keys in map functions

const fs = require('fs');
const path = require('path');

const EXTENSIONS = ['.tsx', '.jsx', '.ts', '.js'];
const EXCLUDE_DIRS = ['node_modules', '.next', 'dist', 'build', '.git'];

function searchInFile(filePath) {
  try {
    const content = fs.readFileSync(filePath, 'utf-8');
    const lines = content.split('\n');
    const issues = [];
    
    // Look for .map( patterns without key props
    const mapPattern = /\.map\s*\(/g;
    const keyPattern = /key\s*=/;
    
    lines.forEach((line, index) => {
      if (mapPattern.test(line)) {
        // Check next 5 lines for key prop
        let hasKey = false;
        for (let i = 0; i < 5 && index + i < lines.length; i++) {
          if (keyPattern.test(lines[index + i])) {
            hasKey = true;
            break;
          }
        }
        
        if (!hasKey) {
          issues.push({
            file: filePath,
            line: index + 1,
            code: line.trim()
          });
        }
      }
    });
    
    return issues;
  } catch (error) {
    return [];
  }
}

function searchDirectory(dir) {
  const results = [];
  
  try {
    const files = fs.readdirSync(dir);
    
    files.forEach(file => {
      const fullPath = path.join(dir, file);
      const stat = fs.statSync(fullPath);
      
      if (stat.isDirectory() && !EXCLUDE_DIRS.includes(file)) {
        results.push(...searchDirectory(fullPath));
      } else if (stat.isFile() && EXTENSIONS.includes(path.extname(file))) {
        const issues = searchInFile(fullPath);
        results.push(...issues);
      }
    });
  } catch (error) {
    console.error(`Error reading directory ${dir}:`, error.message);
  }
  
  return results;
}

console.log('🔍 Searching for React components with missing keys...\n');

const appDir = path.join(__dirname, '../app');
const issues = searchDirectory(appDir);

if (issues.length === 0) {
  console.log('✅ No missing key props found!');
} else {
  console.log(`⚠️  Found ${issues.length} potential missing key props:\n`);
  
  issues.forEach(issue => {
    console.log(`📁 ${issue.file}`);
    console.log(`   Line ${issue.line}: ${issue.code}`);
    console.log('');
  });
}
