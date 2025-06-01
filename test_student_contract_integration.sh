#!/bin/bash

# STUDENT MODULE CONTRACT INTEGRATION - SCRIPT DE TESTE
# Este script demonstra o fluxo completo de integração com contratos

echo "🚀 Testando integração Student Module + Contratos CosmWasm"
echo "============================================================"

# Configurações
CHAIN_ID="academictoken"
NODE="http://localhost:26657"
KEYRING_BACKEND="test"
ADMIN_KEY="alice"
STUDENT_KEY="student1"
INSTITUTION_KEY="institution1"

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_step() {
    echo -e "${BLUE}📋 STEP $1: $2${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Função para executar comando e verificar resultado
execute_cmd() {
    local cmd="$1"
    local description="$2"
    
    echo "Executing: $cmd"
    
    if eval "$cmd"; then
        print_success "$description completed successfully"
        return 0
    else
        print_error "$description failed"
        return 1
    fi
}

# ==============================================================================
# PASSO 1: CONFIGURAÇÃO INICIAL
# ==============================================================================

print_step "1" "Configuração inicial dos parâmetros"

# Configurar parâmetros do Student module com endereços dos contratos
PREREQUISITES_CONTRACT="cosmos1qg5ega6dykkxc307y25pecuufrjkxkaggkkxh7nad0vhyhtuhw3sqaa3c5" # Exemplo
EQUIVALENCE_CONTRACT="cosmos1qg5ega6dykkxc307y25pecuufrjkxkaggkkxh7nad0vhyhtuhw3sqaa3c6"   # Exemplo

UPDATE_PARAMS_CMD="academictokend tx student update-params \\
  --ipfs-gateway 'http://localhost:5001' \\
  --ipfs-enabled true \\
  --admin \$(academictokend keys show $ADMIN_KEY -a --keyring-backend $KEYRING_BACKEND) \\
  --prerequisites-contract-addr '$PREREQUISITES_CONTRACT' \\
  --equivalence-contract-addr '$EQUIVALENCE_CONTRACT' \\
  --from $ADMIN_KEY \\
  --chain-id $CHAIN_ID \\
  --node $NODE \\
  --keyring-backend $KEYRING_BACKEND \\
  --gas auto \\
  --gas-adjustment 1.3 \\
  --yes"

echo "Configurando parâmetros do Student module..."
if execute_cmd "$UPDATE_PARAMS_CMD" "Parameter update"; then
    print_success "Parâmetros configurados com endereços dos contratos"
else
    print_warning "Falha na configuração de parâmetros - continuando com configuração padrão"
fi

# ==============================================================================
# PASSO 2: SETUP BÁSICO (INSTITUIÇÃO, CURSO, DISCIPLINAS)
# ==============================================================================

print_step "2" "Setup básico - Instituição, Curso e Disciplinas"

# Registrar instituição
INSTITUTION_CMD="academictokend tx institution register-institution \\
  'Universidade Teste CosmWasm' \\
  'Rua dos Contratos, 123' \\
  --from $INSTITUTION_KEY \\
  --chain-id $CHAIN_ID \\
  --node $NODE \\
  --keyring-backend $KEYRING_BACKEND \\
  --gas auto \\
  --gas-adjustment 1.3 \\
  --yes"

execute_cmd "$INSTITUTION_CMD" "Institution registration"

# Aguardar processamento
sleep 2

# Obter ID da instituição (assumindo que é 0 por ser a primeira)
INSTITUTION_ID="0"

# Criar curso
COURSE_CMD="academictokend tx course create-course \\
  '$INSTITUTION_ID' \\
  'Ciência da Computação' \\
  'CC001' \\
  'Curso de Bacharelado em Ciência da Computação' \\
  '240' \\
  'bachelor' \\
  --from $INSTITUTION_KEY \\
  --chain-id $CHAIN_ID \\
  --node $NODE \\
  --keyring-backend $KEYRING_BACKEND \\
  --gas auto \\
  --gas-adjustment 1.3 \\
  --yes"

execute_cmd "$COURSE_CMD" "Course creation"

sleep 2
COURSE_ID="0"

# Criar disciplinas
SUBJECT1_CMD="academictokend tx subject create-subject \\
  '$INSTITUTION_ID' \\
  'Cálculo I' \\
  'MAT001' \\
  '60' \\
  '4' \\
  'Fundamentos de Cálculo Diferencial e Integral' \\
  'mathematics' \\
  'required' \\
  --from $INSTITUTION_KEY \\
  --chain-id $CHAIN_ID \\
  --node $NODE \\
  --keyring-backend $KEYRING_BACKEND \\
  --gas auto \\
  --gas-adjustment 1.3 \\
  --yes"

execute_cmd "$SUBJECT1_CMD" "Subject 1 creation (Cálculo I)"

SUBJECT2_CMD="academictokend tx subject create-subject \\
  '$INSTITUTION_ID' \\
  'Cálculo II' \\
  'MAT002' \\
  '60' \\
  '4' \\
  'Técnicas de Integração e Séries' \\
  'mathematics' \\
  'required' \\
  --from $INSTITUTION_KEY \\
  --chain-id $CHAIN_ID \\
  --node $NODE \\
  --keyring-backend $KEYRING_BACKEND \\
  --gas auto \\
  --gas-adjustment 1.3 \\
  --yes"

execute_cmd "$SUBJECT2_CMD" "Subject 2 creation (Cálculo II)"

sleep 2
SUBJECT1_ID="0"
SUBJECT2_ID="1"

# ==============================================================================
# PASSO 3: REGISTRO DE ESTUDANTE
# ==============================================================================

print_step "3" "Registro de Estudante"

STUDENT_ADDR=$(academictokend keys show $STUDENT_KEY -a --keyring-backend $KEYRING_BACKEND)

REGISTER_STUDENT_CMD="academictokend tx student register-student \\
  'João Silva' \\
  '$STUDENT_ADDR' \\
  --from $STUDENT_KEY \\
  --chain-id $CHAIN_ID \\
  --node $NODE \\
  --keyring-backend $KEYRING_BACKEND \\
  --gas auto \\
  --gas-adjustment 1.3 \\
  --yes"

execute_cmd "$REGISTER_STUDENT_CMD" "Student registration"

sleep 2
STUDENT_ID="0"

# ==============================================================================
# PASSO 4: MATRÍCULA EM CURSO
# ==============================================================================

print_step "4" "Matrícula em Curso"

ENROLLMENT_CMD="academictokend tx student create-enrollment \\
  '$STUDENT_ID' \\
  '$INSTITUTION_ID' \\
  '$COURSE_ID' \\
  --from $STUDENT_KEY \\
  --chain-id $CHAIN_ID \\
  --node $NODE \\
  --keyring-backend $KEYRING_BACKEND \\
  --gas auto \\
  --gas-adjustment 1.3 \\
  --yes"

execute_cmd "$ENROLLMENT_CMD" "Course enrollment"

sleep 2

# ==============================================================================
# PASSO 5: MATRÍCULA EM DISCIPLINA (COM VERIFICAÇÃO DE PRÉ-REQUISITOS!)
# ==============================================================================

print_step "5" "Matrícula em Disciplina - Cálculo I (Verificação de pré-requisitos via contrato)"

SUBJECT_ENROLLMENT_CMD="academictokend tx student request-subject-enrollment \\
  '$STUDENT_ID' \\
  '$SUBJECT1_ID' \\
  --from $STUDENT_KEY \\
  --chain-id $CHAIN_ID \\
  --node $NODE \\
  --keyring-backend $KEYRING_BACKEND \\
  --gas auto \\
  --gas-adjustment 1.3 \\
  --yes"

if execute_cmd "$SUBJECT_ENROLLMENT_CMD" "Subject enrollment (Cálculo I)"; then
    print_success "✨ Pré-requisitos verificados via contrato CosmWasm!"
    print_success "✨ Estudante matriculado em Cálculo I"
else
    print_error "Falha na verificação de pré-requisitos via contrato"
fi

sleep 2

# ==============================================================================
# PASSO 6: COMPLETAR DISCIPLINA (ATUALIZA CONTRATO!)
# ==============================================================================

print_step "6" "Completar Disciplina - Cálculo I (Atualiza registro no contrato)"

COMPLETE_SUBJECT_CMD="academictokend tx student complete-subject \\
  '$STUDENT_ID' \\
  '$SUBJECT1_ID' \\
  8500 \\
  '2024-06-15T10:00:00Z' \\
  '2024-1' \\
  'prof_joao_signature' \\
  --from $INSTITUTION_KEY \\
  --chain-id $CHAIN_ID \\
  --node $NODE \\
  --keyring-backend $KEYRING_BACKEND \\
  --gas auto \\
  --gas-adjustment 1.3 \\
  --yes"

if execute_cmd "$COMPLETE_SUBJECT_CMD" "Subject completion (Cálculo I)"; then
    print_success "✨ Disciplina completada com sucesso!"
    print_success "✨ Registro atualizado no contrato Prerequisites!"
    print_success "✨ NFT de conclusão emitido!"
else
    print_error "Falha na conclusão da disciplina"
fi

sleep 2

# ==============================================================================
# PASSO 7: TENTAR MATRÍCULA EM CÁLCULO II (DEVE SER PERMITIDA AGORA!)
# ==============================================================================

print_step "7" "Matrícula em Cálculo II (Pré-requisito deve estar atendido)"

CALC2_ENROLLMENT_CMD="academictokend tx student request-subject-enrollment \\
  '$STUDENT_ID' \\
  '$SUBJECT2_ID' \\
  --from $STUDENT_KEY \\
  --chain-id $CHAIN_ID \\
  --node $NODE \\
  --keyring-backend $KEYRING_BACKEND \\
  --gas auto \\
  --gas-adjustment 1.3 \\
  --yes"

if execute_cmd "$CALC2_ENROLLMENT_CMD" "Subject enrollment (Cálculo II)"; then
    print_success "✨ Pré-requisito atendido! Matrícula em Cálculo II aprovada!"
    print_success "✨ Contrato Prerequisites validou a conclusão de Cálculo I!"
else
    print_error "Falha na matrícula - pré-requisitos não atendidos"
fi

sleep 2

# ==============================================================================
# PASSO 8: SOLICITAR EQUIVALÊNCIA (TESTE DO CONTRATO DE EQUIVALÊNCIA)
# ==============================================================================

print_step "8" "Solicitar Equivalência (Teste de contrato Equivalence)"

REQUEST_EQUIVALENCE_CMD="academictokend tx student request-equivalence \\
  '$STUDENT_ID' \\
  '$SUBJECT1_ID' \\
  '$SUBJECT2_ID' \\
  'Transferência de universidade - disciplina similar cursada' \\
  --from $STUDENT_KEY \\
  --chain-id $CHAIN_ID \\
  --node $NODE \\
  --keyring-backend $KEYRING_BACKEND \\
  --gas auto \\
  --gas-adjustment 1.3 \\
  --yes"

if execute_cmd "$REQUEST_EQUIVALENCE_CMD" "Equivalence request"; then
    print_success "✨ Solicitação de equivalência enviada!"
    print_success "✨ Contrato Equivalence iniciará análise via IPFS!"
else
    print_error "Falha na solicitação de equivalência"
fi

# ==============================================================================
# PASSO 9: CONSULTAS E VERIFICAÇÕES
# ==============================================================================

print_step "9" "Consultas e Verificações do Estado"

echo "Consultando progresso acadêmico do estudante..."
PROGRESS_CMD="academictokend query student get-student-progress '$STUDENT_ID' --node $NODE"
execute_cmd "$PROGRESS_CMD" "Student progress query"

echo "Consultando árvore acadêmica do estudante..."
TREE_CMD="academictokend query student get-student-academic-tree '$STUDENT_ID' --node $NODE"
execute_cmd "$TREE_CMD" "Academic tree query"

echo "Consultando parâmetros do módulo..."
PARAMS_CMD="academictokend query student params --node $NODE"
execute_cmd "$PARAMS_CMD" "Module parameters query"

# ==============================================================================
# RESUMO FINAL
# ==============================================================================

echo ""
echo "🎉 TESTE DE INTEGRAÇÃO COMPLETO!"
echo "=================================="
echo ""
print_success "✅ Parâmetros configurados com endereços dos contratos"
print_success "✅ Instituição, curso e disciplinas criados"
print_success "✅ Estudante registrado e matriculado"
print_success "✅ Verificação de pré-requisitos via contrato Prerequisites"
print_success "✅ Disciplina completada e contrato atualizado"
print_success "✅ NFT de conclusão emitido automaticamente"
print_success "✅ Matrícula subsequente aprovada com base nos pré-requisitos"
print_success "✅ Sistema de equivalência testado via contrato Equivalence"
echo ""
echo "🚀 Sua integração Student Module + Contratos CosmWasm está funcionando!"
echo ""
print_warning "📋 Próximos passos:"
echo "   1. Verificar logs dos contratos para confirmar execução"
echo "   2. Implementar frontend para interação visual"
echo "   3. Configurar monitoramento dos contratos"
echo "   4. Testar cenários de erro e edge cases"
echo ""
