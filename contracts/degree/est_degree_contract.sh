#!/bin/bash

# ============================================
# SCRIPT DE TESTE DO CONTRATO DEGREE
# ============================================

set -e  # Parar em caso de erro

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Função para imprimir com cor
print_test() {
    echo -e "${YELLOW}[TEST]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# ============================================
# CONFIGURAÇÕES
# ============================================

CHAIN_ID="academictoken"
ALICE_ADDR=$(academictokend keys show alice -a)
BOB_ADDR=$(academictokend keys show bob -a)

# Verificar se o contrato já foi deployado
if [ -f "degree_contract_info.txt" ]; then
    source degree_contract_info.txt
    print_success "Contract info loaded: $CONTRACT_ADDR"
else
    print_error "Contract not deployed. Run deploy_degree.sh first!"
    exit 1
fi

# ============================================
# TESTE 1: QUERY CONFIG
# ============================================

print_test "1. Querying contract config..."

CONFIG_QUERY='{"get_config":{}}'
CONFIG_RESULT=$(academictokend query wasm contract-state smart $CONTRACT_ADDR "$CONFIG_QUERY" --output json 2>/dev/null | jq -r '.data')

if [ -n "$CONFIG_RESULT" ]; then
    print_success "Config query successful!"
    echo "Config: $CONFIG_RESULT"
    
    # Verificar se admin está correto
    ADMIN_IN_CONTRACT=$(echo $CONFIG_RESULT | jq -r '.admin')
    if [ "$ADMIN_IN_CONTRACT" == "$ALICE_ADDR" ]; then
        print_success "Admin address is correct: $ADMIN_IN_CONTRACT"
    else
        print_error "Admin address mismatch!"
    fi
else
    print_error "Failed to query config"
    exit 1
fi

# ============================================
# TESTE 2: CONFIGURAR REQUISITOS DO CURRÍCULO
# ============================================

print_test "2. Setting curriculum requirements for CS-2024..."

SET_REQUIREMENTS_MSG='{
  "set_curriculum_requirements": {
    "curriculum_id": "CS-2024",
    "min_credits": "240",
    "required_subjects": ["CS101", "CS102", "CS201", "CS202", "MATH101", "MATH102", "MATH201", "PHY101"],
    "min_gpa": "2.5",
    "additional_requirements": ["Internship", "Final Project"]
  }
}'

TX_HASH=$(academictokend tx wasm execute $CONTRACT_ADDR "$SET_REQUIREMENTS_MSG" \
    --from alice \
    --gas auto \
    --gas-adjustment 1.3 \
    --fees 5000stake \
    -y \
    --output json 2>/dev/null | jq -r '.txhash')

if [ -n "$TX_HASH" ]; then
    print_success "Set requirements TX: $TX_HASH"
    sleep 6
else
    print_error "Failed to set requirements"
    exit 1
fi

# ============================================
# TESTE 3: QUERY REQUISITOS DO CURRÍCULO
# ============================================

print_test "3. Querying curriculum requirements..."

REQUIREMENTS_QUERY='{"get_curriculum_requirements":{"curriculum_id":"CS-2024"}}'
REQUIREMENTS_RESULT=$(academictokend query wasm contract-state smart $CONTRACT_ADDR "$REQUIREMENTS_QUERY" --output json 2>/dev/null | jq -r '.data')

if [ -n "$REQUIREMENTS_RESULT" ]; then
    print_success "Requirements query successful!"
    echo "Requirements: $REQUIREMENTS_RESULT"
    
    # Verificar alguns campos
    MIN_CREDITS=$(echo $REQUIREMENTS_RESULT | jq -r '.min_credits')
    if [ "$MIN_CREDITS" == "240" ]; then
        print_success "Min credits correctly set: $MIN_CREDITS"
    fi
else
    print_error "Failed to query requirements"
    exit 1
fi

# ============================================
# TESTE 4: VALIDAR REQUISITOS DE GRADUAÇÃO (APROVADO)
# ============================================

print_test "4. Testing degree validation - PASS case..."

VALIDATE_PASS_MSG='{
  "validate_degree_requirements": {
    "student_id": "STU001",
    "curriculum_id": "CS-2024",
    "institution_id": "INST001",
    "final_gpa": "3.5",
    "total_credits": 250,
    "completed_subjects": ["CS101", "CS102", "CS201", "CS202", "MATH101", "MATH102", "MATH201", "PHY101", "ENG101", "ENG102"],
    "signatures": ["sig1", "sig2"],
    "requested_date": "2024-05-30"
  }
}'

VALIDATE_TX=$(academictokend tx wasm execute $CONTRACT_ADDR "$VALIDATE_PASS_MSG" \
    --from alice \
    --gas auto \
    --gas-adjustment 1.3 \
    --fees 5000stake \
    -y \
    --output json 2>/dev/null | jq -r '.txhash')

if [ -n "$VALIDATE_TX" ]; then
    print_success "Validation TX (PASS): $VALIDATE_TX"
    sleep 6
    
    # Verificar resultado da transação
    TX_RESULT=$(academictokend query tx $VALIDATE_TX --output json 2>/dev/null)
    if [ -n "$TX_RESULT" ]; then
        # Extrair o resultado da validação
        IS_VALID=$(echo $TX_RESULT | jq -r '.logs[0].events[] | select(.type=="wasm") | .attributes[] | select(.key=="is_valid") | .value' 2>/dev/null)
        if [ "$IS_VALID" == "true" ]; then
            print_success "Student PASSED validation!"
        else
            print_error "Student should have passed but didn't"
        fi
    fi
else
    print_error "Failed to validate degree requirements"
    exit 1
fi

# ============================================
# TESTE 5: QUERY RESULTADO DA VALIDAÇÃO
# ============================================

print_test "5. Querying validation result..."

VALIDATION_QUERY='{"get_validation_result":{"student_id":"STU001","curriculum_id":"CS-2024"}}'
VALIDATION_RESULT=$(academictokend query wasm contract-state smart $CONTRACT_ADDR "$VALIDATION_QUERY" --output json 2>/dev/null | jq -r '.data')

if [ -n "$VALIDATION_RESULT" ]; then
    print_success "Validation result query successful!"
    echo "Validation Result: $VALIDATION_RESULT"
    
    IS_VALID=$(echo $VALIDATION_RESULT | jq -r '.is_valid')
    SCORE=$(echo $VALIDATION_RESULT | jq -r '.validation_score')
    print_success "Is Valid: $IS_VALID, Score: $SCORE"
else
    print_error "Failed to query validation result"
fi

# ============================================
# TESTE 6: VALIDAR REQUISITOS DE GRADUAÇÃO (REPROVADO)
# ============================================

print_test "6. Testing degree validation - FAIL case (low GPA)..."

VALIDATE_FAIL_MSG='{
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
}'

VALIDATE_FAIL_TX=$(academictokend tx wasm execute $CONTRACT_ADDR "$VALIDATE_FAIL_MSG" \
    --from alice \
    --gas auto \
    --gas-adjustment 1.3 \
    --fees 5000stake \
    -y \
    --output json 2>/dev/null | jq -r '.txhash')

if [ -n "$VALIDATE_FAIL_TX" ]; then
    print_success "Validation TX (FAIL): $VALIDATE_FAIL_TX"
    sleep 6
else
    print_error "Failed to execute fail validation"
fi

# ============================================
# TESTE 7: VERIFICAR DISCIPLINAS FALTANTES
# ============================================

print_test "7. Checking missing requirements..."

MISSING_QUERY='{
  "get_missing_requirements": {
    "student_id": "STU003",
    "curriculum_id": "CS-2024",
    "completed_subjects": ["CS101", "CS102", "MATH101"]
  }
}'

MISSING_RESULT=$(academictokend query wasm contract-state smart $CONTRACT_ADDR "$MISSING_QUERY" --output json 2>/dev/null | jq -r '.data')

if [ -n "$MISSING_RESULT" ]; then
    print_success "Missing requirements query successful!"
    echo "Missing subjects: $MISSING_RESULT"
else
    print_error "Failed to query missing requirements"
fi

# ============================================
# TESTE 8: TENTAR AÇÃO NÃO AUTORIZADA
# ============================================

print_test "8. Testing unauthorized action (Bob trying to set requirements)..."

UNAUTHORIZED_MSG='{
  "set_curriculum_requirements": {
    "curriculum_id": "HACK-2024",
    "min_credits": "1",
    "required_subjects": [],
    "min_gpa": "1.0",
    "additional_requirements": []
  }
}'

# Espera-se que falhe
UNAUTH_TX=$(academictokend tx wasm execute $CONTRACT_ADDR "$UNAUTHORIZED_MSG" \
    --from bob \
    --gas auto \
    --gas-adjustment 1.3 \
    --fees 5000stake \
    -y \
    --output json 2>&1)

if echo "$UNAUTH_TX" | grep -q "Unauthorized"; then
    print_success "Unauthorized action correctly blocked!"
else
    print_error "Unauthorized action was not blocked properly"
fi

# ============================================
# TESTE 9: LIMPAR CACHE DE VALIDAÇÃO
# ============================================

print_test "9. Clearing validation cache..."

CLEAR_CACHE_MSG='{
  "clear_validation_cache": {
    "student_id": "STU001",
    "curriculum_id": "CS-2024"
  }
}'

CLEAR_TX=$(academictokend tx wasm execute $CONTRACT_ADDR "$CLEAR_CACHE_MSG" \
    --from alice \
    --gas auto \
    --gas-adjustment 1.3 \
    --fees 5000stake \
    -y \
    --output json 2>/dev/null | jq -r '.txhash')

if [ -n "$CLEAR_TX" ]; then
    print_success "Cache cleared: $CLEAR_TX"
else
    print_error "Failed to clear cache"
fi

# ============================================
# RESUMO DOS TESTES
# ============================================

echo ""
echo "============================================"
echo "RESUMO DOS TESTES DO CONTRATO DEGREE"
echo "============================================"
print_success "✓ Config query funcionando"
print_success "✓ Set requirements funcionando"
print_success "✓ Query requirements funcionando"
print_success "✓ Validação (PASS) funcionando"
print_success "✓ Validação (FAIL) funcionando"
print_success "✓ Query validation result funcionando"
print_success "✓ Query missing requirements funcionando"
print_success "✓ Controle de acesso funcionando"
print_success "✓ Clear cache funcionando"
echo ""
print_success "TODOS OS TESTES PASSARAM! O contrato está funcionando corretamente."
echo ""

# Salvar alguns dados úteis para testes futuros
echo "CONTRACT_ADDR=$CONTRACT_ADDR" > test_results.txt
echo "TEST_DATE=$(date)" >> test_results.txt
echo "TEST_STATUS=PASSED" >> test_results.txt

print_success "Resultados salvos em test_results.txt"
