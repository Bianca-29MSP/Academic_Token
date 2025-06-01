#!/bin/bash

echo "=========================================="
echo "Academic Token - CLEAN RESET Setup Script"
echo "=========================================="

echo ""
echo "⚠️  RESETTING BLOCKCHAIN STATE..."
echo "=========================================="

# Stop the chain if running
echo "Stopping any running academictokend processes..."
pkill -f academictokend || true
sleep 3

# Clean reset the blockchain state
echo "Cleaning blockchain state..."
academictokend unsafe-reset-all --home ~/.academictoken 2>/dev/null || true

# Remove any corrupted data
rm -rf ~/.academictoken/data/ 2>/dev/null || true
rm -rf ~/.academictoken/config/addrbook.json 2>/dev/null || true

echo "✅ Blockchain state reset completed"

echo ""
echo "🚀 STARTING FRESH CHAIN..."
echo "=========================================="

# Start the chain in background
echo "Starting academictokend..."
academictokend start --home ~/.academictoken > /tmp/academictoken.log 2>&1 &
CHAIN_PID=$!

echo "Chain started with PID: $CHAIN_PID"
echo "Waiting for chain to be ready..."

# Wait for chain to be ready
sleep 10

# Check if chain is responding
for i in {1..30}; do
    if curl -s http://localhost:1317/cosmos/base/tendermint/v1beta1/node_info > /dev/null 2>&1; then
        echo "✅ Chain is ready!"
        break
    fi
    echo "Waiting for chain... ($i/30)"
    sleep 2
done

echo ""
echo "📝 CREATING TEST DATA..."
echo "=========================================="

echo ""
echo "1. CREATING INSTITUTION..."
echo "=================================="

academictokend tx institution register-institution \
  "Universidade Federal de Minas Gerais" \
  "Av. Antônio Carlos, 6627, Belo Horizonte, MG, Brasil" \
  --from alice \
  --chain-id academictoken \
  --gas auto \
  --gas-adjustment 1.5 \
  --fees 1000stake \
  --yes

sleep 3

echo "Authorizing institution-1..."
academictokend tx institution update-institution \
  "institution-1" \
  "Universidade Federal de Minas Gerais" \
  "Av. Antônio Carlos, 6627, Belo Horizonte, MG, Brasil" \
  "true" \
  --from alice \
  --chain-id academictoken \
  --gas auto \
  --gas-adjustment 1.5 \
  --fees 1000stake \
  --yes

sleep 3

echo ""
echo "2. CREATING COURSE..."
echo "=================================="

academictokend tx course create-course \
  "institution-1" \
  "Ciência da Computação" \
  "CC" \
  "Bacharelado em Ciência da Computação" \
  "240" \
  "undergraduate" \
  --from alice \
  --chain-id academictoken \
  --gas auto \
  --gas-adjustment 1.5 \
  --fees 1000stake \
  --yes

sleep 3

echo ""
echo "3. CREATING SUBJECTS ONE BY ONE..."
echo "=================================="

echo "Creating Subject 1: Cálculo I"
academictokend tx subject create-subject \
  "institution-1" \
  "course-1" \
  "Cálculo I" \
  "CALC1" \
  90 \
  6 \
  "Fundamentos do cálculo diferencial e integral" \
  "required" \
  "Mathematics" \
  --from alice \
  --gas auto \
  --gas-adjustment 2.0 \
  --fees 2000stake \
  --yes

echo "Waiting after first subject..."
sleep 5

echo "Testing if subject creation worked..."
academictokend query subject list-subjects

echo ""
echo "Creating Subject 2: Programação I"
academictokend tx subject create-subject \
  "institution-1" \
  "course-1" \
  "Programação I" \
  "PROG1" \
  60 \
  4 \
  "Introdução à programação estruturada" \
  "required" \
  "Computer Science" \
  --from alice \
  --gas auto \
  --gas-adjustment 2.0 \
  --fees 2000stake \
  --yes

echo "Waiting after second subject..."
sleep 5

echo "Testing subjects again..."
academictokend query subject list-subjects

echo ""
echo "Creating Subject 3: Álgebra Linear"
academictokend tx subject create-subject \
  "institution-1" \
  "course-1" \
  "Álgebra Linear" \
  "ALG1" \
  75 \
  5 \
  "Espaços vetoriais e transformações lineares" \
  "required" \
  "Mathematics" \
  --from alice \
  --gas auto \
  --gas-adjustment 2.0 \
  --fees 2000stake \
  --yes

sleep 5

echo ""
echo "4. FINAL VERIFICATION..."
echo "=================================="

echo "Institutions:"
academictokend query institution list-institutions

echo ""
echo "Courses:"
academictokend query course list-courses

echo ""
echo "Subjects:"
academictokend query subject list-subjects

echo ""
echo "5. CREATING TOKEN DEFINITIONS..."
echo "=================================="

# Get actual subject IDs
SUBJECT_1=$(academictokend query subject list-subjects --output json 2>/dev/null | jq -r '.subjects[0].index // "subject-1"')
SUBJECT_2=$(academictokend query subject list-subjects --output json 2>/dev/null | jq -r '.subjects[1].index // "subject-2"')
SUBJECT_3=$(academictokend query subject list-subjects --output json 2>/dev/null | jq -r '.subjects[2].index // "subject-3"')

echo "Found subjects: $SUBJECT_1, $SUBJECT_2, $SUBJECT_3"

if [ "$SUBJECT_1" != "null" ] && [ "$SUBJECT_1" != "" ]; then
    echo "Creating token for $SUBJECT_1..."
    academictokend tx tokendef create-token-definition \
      "$SUBJECT_1" \
      "Cálculo I Token" \
      "CALC1-TKN" \
      "Token for Calculus I completion" \
      "ACHIEVEMENT" \
      --is-transferable=false \
      --is-burnable=false \
      --max-supply=1000 \
      --image-uri="https://academic-tokens.example.com/calc1.png" \
      --from alice \
      --chain-id academictoken \
      --gas auto \
      --gas-adjustment 1.5 \
      --fees 1000stake \
      --yes
    
    sleep 3
fi

if [ "$SUBJECT_2" != "null" ] && [ "$SUBJECT_2" != "" ]; then
    echo "Creating token for $SUBJECT_2..."
    academictokend tx tokendef create-token-definition \
      "$SUBJECT_2" \
      "Programming I Token" \
      "PROG1-TKN" \
      "Token for Programming I completion" \
      "ACHIEVEMENT" \
      --is-transferable=true \
      --is-burnable=false \
      --max-supply=800 \
      --image-uri="https://academic-tokens.example.com/prog1.png" \
      --from alice \
      --chain-id academictoken \
      --gas auto \
      --gas-adjustment 1.5 \
      --fees 1000stake \
      --yes
    
    sleep 3
fi

echo ""
echo "6. FINAL STATUS..."
echo "=================================="

echo "Token Definitions:"
academictokend query tokendef list-token-definitions

echo ""
echo "=========================================="
echo "✅ CLEAN SETUP COMPLETED SUCCESSFULLY!"
echo "=========================================="

INSTITUTION_COUNT=$(academictokend query institution list-institutions --output json 2>/dev/null | jq '.institution | length' 2>/dev/null || echo "0")
COURSE_COUNT=$(academictokend query course list-courses --output json 2>/dev/null | jq '.course | length' 2>/dev/null || echo "0")
SUBJECT_COUNT=$(academictokend query subject list-subjects --output json 2>/dev/null | jq '.subjects | length' 2>/dev/null || echo "0")
TOKEN_COUNT=$(academictokend query tokendef list-token-definitions --output json 2>/dev/null | jq '.tokenDefinitions | length' 2>/dev/null || echo "0")

echo ""
echo "🎯 FINAL SUMMARY:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🏛️  INSTITUTIONS: $INSTITUTION_COUNT"
echo "🎓 COURSES: $COURSE_COUNT"
echo "📚 SUBJECTS: $SUBJECT_COUNT"
echo "🪙  TOKENS: $TOKEN_COUNT"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if [ "$SUBJECT_COUNT" -gt "1" ]; then
    echo "✅ SUCCESS: Multiple subjects created successfully!"
    echo "🚀 System is working properly!"
else
    echo "❌ WARNING: Only $SUBJECT_COUNT subject(s) created"
    echo "💡 Check logs at /tmp/academictoken.log for errors"
fi

echo ""
echo "Chain is running with PID: $CHAIN_PID"
echo "To stop: kill $CHAIN_PID"
echo "To view logs: tail -f /tmp/academictoken.log"
