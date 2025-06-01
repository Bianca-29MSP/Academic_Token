#!/bin/bash

echo "🔍 Testing Academic Token REST API..."
echo "======================================"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test health endpoint
echo -e "\n${YELLOW}1. Testing Health Endpoint${NC}"
echo "GET http://localhost:1318/health"
response=$(curl -s -w "HTTP_STATUS:%{http_code}" http://localhost:1318/health 2>/dev/null)
if [[ $response == *"HTTP_STATUS:200"* ]]; then
    echo -e "${GREEN}✅ Health check: OK${NC}"
else
    echo -e "${RED}❌ Health check failed${NC}"
    echo "Make sure the REST server is running:"
    echo "go run cmd/rest-server/main.go"
    exit 1
fi

# Test node info
echo -e "\n${YELLOW}2. Testing Node Info${NC}"
echo "GET http://localhost:1318/cosmos/base/tendermint/v1beta1/node_info"
curl -s http://localhost:1318/cosmos/base/tendermint/v1beta1/node_info | jq '.' 2>/dev/null
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Node info: OK${NC}"
else
    echo -e "${RED}❌ Node info failed or jq not installed${NC}"
fi

# Test institutions
echo -e "\n${YELLOW}3. Testing Institutions Endpoint${NC}"
echo "GET http://localhost:1318/academic/institution/list"
institutions=$(curl -s http://localhost:1318/academic/institution/list)
echo "$institutions" | jq '.' 2>/dev/null
count=$(echo "$institutions" | jq 'length' 2>/dev/null)
if [ "$count" -gt 0 ]; then
    echo -e "${GREEN}✅ Institutions loaded: $count items${NC}"
else
    echo -e "${RED}❌ No institutions found${NC}"
fi

# Test courses
echo -e "\n${YELLOW}4. Testing Courses Endpoint${NC}"
echo "GET http://localhost:1318/academic/course/list"
courses=$(curl -s http://localhost:1318/academic/course/list)
echo "$courses" | jq '.' 2>/dev/null
count=$(echo "$courses" | jq 'length' 2>/dev/null)
if [ "$count" -gt 0 ]; then
    echo -e "${GREEN}✅ Courses loaded: $count items${NC}"
else
    echo -e "${RED}❌ No courses found${NC}"
fi

# Test subjects
echo -e "\n${YELLOW}5. Testing Subjects Endpoint${NC}"
echo "GET http://localhost:1318/academic/subject/list"
subjects=$(curl -s http://localhost:1318/academic/subject/list)
echo "$subjects" | jq '.' 2>/dev/null
count=$(echo "$subjects" | jq 'length' 2>/dev/null)
if [ "$count" -gt 0 ]; then
    echo -e "${GREEN}✅ Subjects loaded: $count items${NC}"
else
    echo -e "${RED}❌ No subjects found${NC}"
fi

echo -e "\n${GREEN}🎉 API Testing Complete!${NC}"
echo -e "${YELLOW}Now test the frontend at: http://localhost:3000/equivalences${NC}"
