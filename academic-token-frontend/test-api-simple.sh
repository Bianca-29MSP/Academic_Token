#!/bin/bash

# Set execution permission for test script
chmod +x test-api.sh

echo "🔍 Testing Academic Token REST API..."
echo "======================================"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if server is running
echo -e "\n${YELLOW}Checking if REST server is running...${NC}"
if curl -s http://localhost:1318/health > /dev/null; then
    echo -e "${GREEN}✅ REST server is running${NC}"
else
    echo -e "${RED}❌ REST server is not running${NC}"
    echo -e "${YELLOW}Starting REST server...${NC}"
    echo "Please run in another terminal:"
    echo "cd /Users/biancamsp/Desktop/Academic_Token/academictoken/academictoken"
    echo "go run cmd/rest-server/main.go"
    echo ""
    exit 1
fi

# Test health endpoint
echo -e "\n${YELLOW}1. Testing Health Endpoint${NC}"
echo "GET http://localhost:1318/health"
response=$(curl -s -w "HTTP_STATUS:%{http_code}" http://localhost:1318/health 2>/dev/null)
if [[ $response == *"HTTP_STATUS:200"* ]]; then
    echo -e "${GREEN}✅ Health check: OK${NC}"
else
    echo -e "${RED}❌ Health check failed${NC}"
    exit 1
fi

# Test node info
echo -e "\n${YELLOW}2. Testing Node Info${NC}"
echo "GET http://localhost:1318/cosmos/base/tendermint/v1beta1/node_info"
nodeinfo=$(curl -s http://localhost:1318/cosmos/base/tendermint/v1beta1/node_info)
echo "$nodeinfo"
if [[ $nodeinfo == *"academictoken"* ]]; then
    echo -e "${GREEN}✅ Node info: OK${NC}"
else
    echo -e "${RED}❌ Node info failed${NC}"
fi

# Test institutions
echo -e "\n${YELLOW}3. Testing Institutions Endpoint${NC}"
echo "GET http://localhost:1318/academic/institution/list"
institutions=$(curl -s http://localhost:1318/academic/institution/list)
echo "$institutions"

# Count institutions
if [[ $institutions == *"UFJF"* ]]; then
    echo -e "${GREEN}✅ Institutions loaded successfully${NC}"
else
    echo -e "${RED}❌ No institutions found${NC}"
fi

# Test courses
echo -e "\n${YELLOW}4. Testing Courses Endpoint${NC}"
echo "GET http://localhost:1318/academic/course/list"
courses=$(curl -s http://localhost:1318/academic/course/list)
echo "$courses"

if [[ $courses == *"Ciência da Computação"* ]]; then
    echo -e "${GREEN}✅ Courses loaded successfully${NC}"
else
    echo -e "${RED}❌ No courses found${NC}"
fi

# Test subjects
echo -e "\n${YELLOW}5. Testing Subjects Endpoint${NC}"
echo "GET http://localhost:1318/academic/subject/list"
subjects=$(curl -s http://localhost:1318/academic/subject/list)
echo "$subjects"

if [[ $subjects == *"Cálculo I"* ]]; then
    echo -e "${GREEN}✅ Subjects loaded successfully${NC}"
else
    echo -e "${RED}❌ No subjects found${NC}"
fi

echo -e "\n${GREEN}🎉 API Testing Complete!${NC}"
echo -e "${YELLOW}Now test the frontend at: http://localhost:3000/equivalences${NC}"
echo -e "${YELLOW}Check the Debug tab to see real-time connection status${NC}"
