#!/bin/bash

echo "🔧 REBUILDING PROJECT WITH TOKENDEF FIX..."
echo "============================================"

cd /Users/biancamsp/Desktop/Academic_Token/academictoken/academictoken

# Rebuild the project
echo "Building..."
make install

if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
else
    echo "❌ Build failed!"
    exit 1
fi

echo ""
echo "🧪 TESTING TOKENDEF FIX..."
echo "=========================="

# Stop current chain
echo "Stopping current chain..."
pkill -f academictokend || true
sleep 3

# Quick reset
echo "Quick reset..."
academictokend unsafe-reset-all --home ~/.academictoken 2>/dev/null || true

# Start chain
echo "Starting chain..."
academictokend start --home ~/.academictoken > /tmp/academictoken.log 2>&1 &
CHAIN_PID=$!

echo "Chain started with PID: $CHAIN_PID"
echo "Waiting for chain to be ready..."
sleep 10

# Quick test setup
echo "Creating test institution..."
academictokend tx institution register-institution \
  "Test University" \
  "Test Address" \
  --from alice \
  --chain-id academictoken \
  --gas auto \
  --gas-adjustment 1.5 \
  --fees 1000stake \
  --yes

sleep 3

echo "Authorizing institution..."
academictokend tx institution update-institution \
  "institution-1" \
  "Test University" \
  "Test Address" \
  "true" \
  --from alice \
  --chain-id academictoken \
  --gas auto \
  --gas-adjustment 1.5 \
  --fees 1000stake \
  --yes

sleep 3

echo "Creating test course..."
academictokend tx course create-course \
  "institution-1" \
  "Test Course" \
  "TC" \
  "Test Description" \
  "120" \
  "undergraduate" \
  --from alice \
  --chain-id academictoken \
  --gas auto \
  --gas-adjustment 1.5 \
  --fees 1000stake \
  --yes

sleep 3

echo "Creating test subject..."
academictokend tx subject create-subject \
  "institution-1" \
  "course-1" \
  "Test Subject" \
  "TS1" \
  60 \
  4 \
  "Test subject description" \
  "required" \
  "Test Area" \
  --from alice \
  --gas auto \
  --gas-adjustment 2.0 \
  --fees 2000stake \
  --yes

sleep 5

echo "Checking subject creation..."
academictokend query subject list-subjects

echo ""
echo "🪙 CREATING TOKEN DEFINITION..."
echo "==============================="

echo "Creating token definition with fixed courseId and institutionId..."
academictokend tx tokendef create-token-definition \
  "subject-1" \
  "Test Token" \
  "TT1-TKN" \
  "Test token description" \
  "ACHIEVEMENT" \
  --is-transferable=false \
  --is-burnable=false \
  --max-supply=1000 \
  --image-uri="https://example.com/test.png" \
  --from alice \
  --chain-id academictoken \
  --gas auto \
  --gas-adjustment 1.5 \
  --fees 1000stake \
  --yes

sleep 5

echo ""
echo "📋 CHECKING RESULTS..."
echo "======================"

echo "Token Definitions (should now have courseId and institutionId):"
academictokend query tokendef list-token-definitions

echo ""
echo "🎯 VERIFICATION..."
echo "=================="

TOKEN_RESULT=$(academictokend query tokendef list-token-definitions --output json 2>/dev/null)
COURSE_ID=$(echo "$TOKEN_RESULT" | jq -r '.tokenDefinitions[0].courseId // "EMPTY"' 2>/dev/null)
INSTITUTION_ID=$(echo "$TOKEN_RESULT" | jq -r '.tokenDefinitions[0].institutionId // "EMPTY"' 2>/dev/null)

echo "CourseId in token: '$COURSE_ID'"
echo "InstitutionId in token: '$INSTITUTION_ID'"

if [ "$COURSE_ID" != "EMPTY" ] && [ "$COURSE_ID" != "" ] && [ "$COURSE_ID" != "null" ]; then
    echo "✅ SUCCESS: courseId is now populated!"
else
    echo "❌ FAILED: courseId is still empty"
fi

if [ "$INSTITUTION_ID" != "EMPTY" ] && [ "$INSTITUTION_ID" != "" ] && [ "$INSTITUTION_ID" != "null" ]; then
    echo "✅ SUCCESS: institutionId is now populated!"
else
    echo "❌ FAILED: institutionId is still empty"
fi

echo ""
echo "Chain PID: $CHAIN_PID"
echo "To stop: kill $CHAIN_PID"
echo "Logs: tail -f /tmp/academictoken.log"
