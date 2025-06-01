#!/bin/bash

# Test Academic Token Contract Integration
# This script tests the integration between all deployed contracts

set -e

# Load contract addresses
source ./contracts_config.sh

echo "🧪 Testing Academic Token Contract Integration..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Function to test contract query
test_contract_query() {
    local contract_name="$1"
    local contract_address="$2"
    local query="$3"
    
    echo "🔍 Testing $contract_name query..."
    
    result=$($DAEMON query wasm contract-state smart "$contract_address" "$query" \
        --output json 2>/dev/null || echo "ERROR")
    
    if [[ "$result" == "ERROR" ]]; then
        echo "❌ Query failed for $contract_name"
        return 1
    else
        echo "✅ Query successful for $contract_name"
        echo "   Result: $result"
        return 0
    fi
}

echo "1️⃣ Testing Academic-NFT Contract..."
test_contract_query "Academic-NFT" "$ACADEMIC_NFT_CONTRACT" '{"get_config": {}}'
test_contract_query "Academic-NFT" "$ACADEMIC_NFT_CONTRACT" '{"get_statistics": {}}'

echo ""
echo "2️⃣ Testing Progress Contract..."
test_contract_query "Progress" "$PROGRESS_CONTRACT" '{"get_state": {}}'
test_contract_query "Progress" "$PROGRESS_CONTRACT" '{"get_config": {}}'

echo ""
echo "3️⃣ Testing Degree Contract..."
test_contract_query "Degree" "$DEGREE_CONTRACT" '{"get_config": {}}'

echo ""
echo "4️⃣ Testing Schedule Contract..."
test_contract_query "Schedule" "$SCHEDULE_CONTRACT" '{"get_state": {}}'
test_contract_query "Schedule" "$SCHEDULE_CONTRACT" '{"get_config": {}}'

echo ""
echo "5️⃣ Testing Equivalence Contract..."
test_contract_query "Equivalence" "$EQUIVALENCE_CONTRACT" '{"get_state": {}}'

echo ""
echo "📋 Integration Test Summary:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "All basic contract queries tested."
echo ""
echo "🔗 Ready for Module Integration:"
echo "   • Student Module → Progress Contract"
echo "   • AcademicNFT Module → Academic-NFT Contract"
echo "   • Degree Module → Degree Contract"
echo "   • Schedule Module → Schedule Contract"
echo "   • Equivalence Module → Equivalence Contract"
echo ""
echo "✨ Contract integration tests completed!"
