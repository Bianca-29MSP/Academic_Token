#!/bin/bash
# remaining_tests.sh

# Cores
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

print_test() {
    echo -e "${YELLOW}[TEST]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

# Carregar informações do contrato
source degree_contract_info.txt

# TESTE 5: Query resultado da validação
print_test "5. Querying validation result for STU001..."
academictokend query wasm contract-state smart $CONTRACT_ADDR '{"get_validation_result":{"student_id":"STU001","curriculum_id":"CS-2024"}}'

echo ""

# TESTE 6: Validação com falha (GPA baixo)
print_test "6. Testing validation FAIL case (low GPA)..."
academictokend tx wasm execute $CONTRACT_ADDR '{
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
}' --from alice --gas auto --gas-adjustment 1.3 --fees 5000stake -y

sleep 6
echo ""

# TESTE 7: Query resultado da validação com falha
print_test "7. Querying validation result for STU002 (should fail)..."
academictokend query wasm contract-state smart $CONTRACT_ADDR '{"get_validation_result":{"student_id":"STU002","curriculum_id":"CS-2024"}}'

echo ""

# TESTE 8: Verificar disciplinas faltantes
print_test "8. Checking missing subjects..."
academictokend query wasm contract-state smart $CONTRACT_ADDR '{
  "get_missing_requirements": {
    "student_id": "STU003",
    "curriculum_id": "CS-2024",
    "completed_subjects": ["CS101", "CS102", "MATH101"]
  }
}'

echo ""

# TESTE 9: Tentar ação não autorizada com Bob
print_test "9. Testing unauthorized action (Bob trying to set requirements)..."
academictokend tx wasm execute $CONTRACT_ADDR '{
  "set_curriculum_requirements": {
    "curriculum_id": "HACK-2024",
    "min_credits": "1",
    "required_subjects": [],
    "min_gpa": "1.0",
    "additional_requirements": []
  }
}' --from bob --gas auto --gas-adjustment 1.3 --fees 5000stake -y 2>&1 | grep -E "(Unauthorized|failed)"

echo ""

# TESTE 10: Limpar cache
print_test "10. Clearing validation cache for STU001..."
academictokend tx wasm execute $CONTRACT_ADDR '{
  "clear_validation_cache": {
    "student_id": "STU001",
    "curriculum_id": "CS-2024"
  }
}' --from alice --gas auto --gas-adjustment 1.3 --fees 5000stake -y

sleep 6
echo ""

# TESTE 11: Verificar se o cache foi limpo
print_test "11. Checking if cache was cleared (should fail)..."
academictokend query wasm contract-state smart $CONTRACT_ADDR '{"get_validation_result":{"student_id":"STU001","curriculum_id":"CS-2024"}}' 2>&1 | grep -E "(not found|error)"

echo ""
print_success "All tests completed!"
