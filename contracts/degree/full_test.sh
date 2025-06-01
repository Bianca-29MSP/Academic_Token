#!/bin/bash
# full_test.sh

CONTRACT_ADDR="academic1ufs3tlq4umljk0qfe8k5ya0x6hpavn897u2cnf9k0en9jr7qarqqj7temq"

echo "=== DEGREE CONTRACT TESTS ==="
echo ""

# TESTE 1: Verificar resultado da validação anterior
echo "1. Checking validation result for STU001:"
academictokend query wasm contract-state smart $CONTRACT_ADDR '{"get_validation_result":{"student_id":"STU001","curriculum_id":"CS-2024"}}' 2>&1

echo ""
echo "2. Testing validation with FAILING student (low GPA and missing subjects):"
TX_FAIL=$(academictokend tx wasm execute $CONTRACT_ADDR '{
  "validate_degree_requirements": {
    "student_id": "STU002",
    "curriculum_id": "CS-2024",
    "institution_id": "INST001",
    "final_gpa": "2.0",
    "total_credits": 200,
    "completed_subjects": ["CS101", "CS102", "MATH101"],
    "signatures": ["sig1"],
    "requested_date": "2024-05-30"
  }
}' --from alice --gas auto --gas-adjustment 1.3 --fees 5000stake -y --output json 2>/dev/null | jq -r '.txhash')

echo "TX Hash: $TX_FAIL"
sleep 6

echo ""
echo "3. Checking validation result for STU002 (should show failure):"
academictokend query wasm contract-state smart $CONTRACT_ADDR '{"get_validation_result":{"student_id":"STU002","curriculum_id":"CS-2024"}}' 2>&1

echo ""
echo "4. Checking missing requirements for a student:"
academictokend query wasm contract-state smart $CONTRACT_ADDR '{
  "get_missing_requirements": {
    "student_id": "STU003",
    "curriculum_id": "CS-2024",
    "completed_subjects": ["CS101", "CS102", "MATH101", "MATH102"]
  }
}'

echo ""
echo "5. Testing unauthorized access (Bob trying to set requirements):"
academictokend tx wasm execute $CONTRACT_ADDR '{
  "set_curriculum_requirements": {
    "curriculum_id": "UNAUTHORIZED-2024",
    "min_credits": "1",
    "required_subjects": [],
    "min_gpa": "1.0",
    "additional_requirements": []
  }
}' --from bob --gas auto --gas-adjustment 1.3 --fees 5000stake -y 2>&1 | tail -5

echo ""
echo "6. Clear validation cache:"
academictokend tx wasm execute $CONTRACT_ADDR '{
  "clear_validation_cache": {
    "student_id": "STU001",
    "curriculum_id": "CS-2024"
  }
}' --from alice --gas auto --gas-adjustment 1.3 --fees 5000stake -y --output json | jq -r '.txhash'

echo ""
echo "=== TESTS COMPLETE ==="
