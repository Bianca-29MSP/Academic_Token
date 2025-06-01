#!/bin/bash

echo "🚀 Starting Academic Token System"
echo "=================================="

# Check if we're in the right directory
if [ ! -f "package.json" ]; then
    echo "❌ Error: package.json not found"
    echo "Please run this script from the frontend directory:"
    echo "cd /Users/biancamsp/Desktop/Academic_Token/academictoken/academictoken/academic-token-frontend"
    exit 1
fi

echo "✅ Found package.json - we're in the right directory"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Check if backend is running
echo -e "\n${YELLOW}Checking backend status...${NC}"
if curl -s http://localhost:1318/health > /dev/null; then
    echo -e "${GREEN}✅ Backend is running at http://localhost:1318${NC}"
else
    echo -e "${RED}❌ Backend is not running${NC}"
    echo -e "${YELLOW}Starting backend in background...${NC}"
    echo "Please run this command in another terminal:"
    echo -e "${BLUE}cd /Users/biancamsp/Desktop/Academic_Token/academictoken/academictoken${NC}"
    echo -e "${BLUE}go run cmd/rest-server/main.go${NC}"
    echo ""
    echo "Then press any key to continue..."
    read -n 1 -s
fi

# Test API endpoints
echo -e "\n${YELLOW}Testing API endpoints...${NC}"

# Test institutions
institutions=$(curl -s http://localhost:1318/academic/institution/list)
if echo "$institutions" | grep -q "UFJF"; then
    echo -e "${GREEN}✅ Institutions endpoint working${NC}"
else
    echo -e "${RED}❌ Institutions endpoint failed${NC}"
    echo "Response: $institutions"
fi

# Test subjects
subjects=$(curl -s http://localhost:1318/academic/subject/list)
if echo "$subjects" | grep -q "Cálculo I"; then
    echo -e "${GREEN}✅ Subjects endpoint working${NC}"
else
    echo -e "${RED}❌ Subjects endpoint failed${NC}"
    echo "Response: $subjects"
fi

# Start frontend
echo -e "\n${YELLOW}Starting frontend...${NC}"
echo -e "${BLUE}Frontend will be available at: http://localhost:3000${NC}"
echo -e "${BLUE}Equivalences page: http://localhost:3000/equivalences${NC}"
echo ""
echo -e "${YELLOW}Press Ctrl+C to stop${NC}"
echo ""

npm run dev
