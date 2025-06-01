#!/bin/bash

echo "🧹 Cleaning up conflicting files..."

# Remove old context files that conflict with our new structure
rm -rf /Users/biancamsp/Desktop/Academic_Token/academictoken/academictoken/academic-token-frontend/app/context
rm -rf /Users/biancamsp/Desktop/Academic_Token/academictoken/academictoken/academic-token-frontend/app/services

# Remove old components that might have dependencies issues
rm -rf /Users/biancamsp/Desktop/Academic_Token/academictoken/academictoken/academic-token-frontend/app/components

echo "✅ Cleanup complete!"
echo "Remaining structure should be:"
echo "├── app/"
echo "│   ├── hooks/useBlockchain.ts"
echo "│   ├── lib/api.ts"
echo "│   ├── types/blockchain.ts"
echo "│   ├── equivalences/page.tsx"
echo "│   └── page.tsx (main dashboard)"
