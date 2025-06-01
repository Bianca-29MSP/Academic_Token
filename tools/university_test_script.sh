#!/bin/bash

# =============================================================================
# ACADEMIC TOKEN - UNIVERSITY TEST SCRIPT
# =============================================================================
# Este script simula o fluxo completo de uma instituição de ensino no sistema
# AcademicToken, cadastrando uma universidade com curso completo de Ciência da 
# Computação com disciplinas, pré-requisitos e grade curricular.
# =============================================================================

set -e  # Para o script se algum comando falhar

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

# Função para logging
log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}"
}

error() {
    echo -e "${RED}[ERROR] $1${NC}"
    exit 1
}

warning() {
    echo -e "${YELLOW}[WARNING] $1${NC}"
}

info() {
    echo -e "${BLUE}[INFO] $1${NC}"
}

# =============================================================================
# CONFIGURAÇÕES INICIAIS
# =============================================================================
CHAIN_ID="academictoken"
KEYRING="test"
ADMIN_KEY="admin"
UNIVERSITY_KEY="university"
HOME_DIR="$HOME/.academictoken"
BINARY="academictokend"

# Verificar se o binário existe
if ! command -v $BINARY &> /dev/null; then
    error "Binário $BINARY não encontrado. Por favor, compile o projeto primeiro."
fi

log "🚀 Iniciando Teste do Fluxo Institucional - AcademicToken"
log "======================================================"

# =============================================================================
# FASE 0: PREPARAÇÃO DO AMBIENTE
# =============================================================================
log "📋 FASE 0: Preparação do Ambiente"

# Criar chaves se não existirem
if ! $BINARY keys show $ADMIN_KEY --keyring-backend $KEYRING --home $HOME_DIR > /dev/null 2>&1; then
    log "Criando chave do administrador..."
    $BINARY keys add $ADMIN_KEY --keyring-backend $KEYRING --home $HOME_DIR
fi

if ! $BINARY keys show $UNIVERSITY_KEY --keyring-backend $KEYRING --home $HOME_DIR > /dev/null 2>&1; then
    log "Criando chave da universidade..."
    $BINARY keys add $UNIVERSITY_KEY --keyring-backend $KEYRING --home $HOME_DIR
fi

# Obter endereços
ADMIN_ADDR=$($BINARY keys show $ADMIN_KEY -a --keyring-backend $KEYRING --home $HOME_DIR)
UNIVERSITY_ADDR=$($BINARY keys show $UNIVERSITY_KEY -a --keyring-backend $KEYRING --home $HOME_DIR)

log "Admin Address: $ADMIN_ADDR"
log "University Address: $UNIVERSITY_ADDR"

# =============================================================================
# FASE 1: CADASTRO DA INSTITUIÇÃO
# =============================================================================
log "🏛️  FASE 1: Cadastro da Instituição"

info "Registrando Universidade Federal de Tecnologia (UFT)..."

$BINARY tx institution register-institution \
    "UFT" \
    "Universidade Federal de Tecnologia" \
    "BR" \
    "Brasil" \
    "RS" \
    "Porto Alegre" \
    "Av. Tecnologia, 1000" \
    "public" \
    "active" \
    "www.uft.edu.br" \
    "contato@uft.edu.br" \
    "+55-51-1234-5678" \
    "Universidade pública federal especializada em tecnologia e inovação. Fundada em 1980, oferece cursos de graduação e pós-graduação nas áreas de computação, engenharia e ciências exatas." \
    --from $UNIVERSITY_ADDR \
    --keyring-backend $KEYRING \
    --home $HOME_DIR \
    --chain-id $CHAIN_ID \
    --gas auto \
    --gas-adjustment 1.5 \
    --fees 100stake \
    --yes

sleep 2
log "✅ Instituição UFT registrada com sucesso!"

# =============================================================================
# FASE 2: CRIAÇÃO DO CURSO DE CIÊNCIA DA COMPUTAÇÃO
# =============================================================================
log "💻 FASE 2: Criação do Curso de Ciência da Computação"

info "Criando curso de Bacharelado em Ciência da Computação..."

$BINARY tx course create-course \
    "UFT" \
    "CC-UFT-2024" \
    "Ciência da Computação" \
    "Bacharelado em Ciência da Computação" \
    "bachelor" \
    "8" \
    "240" \
    "O curso de Ciência da Computação da UFT forma profissionais aptos a atuar no desenvolvimento de software, pesquisa em computação, gestão de projetos tecnológicos e inovação digital. O currículo integra teoria e prática, com foco em algoritmos, estruturas de dados, engenharia de software, inteligência artificial e sistemas distribuídos." \
    "pt-BR" \
    --from $UNIVERSITY_ADDR \
    --keyring-backend $KEYRING \
    --home $HOME_DIR \
    --chain-id $CHAIN_ID \
    --gas auto \
    --gas-adjustment 1.5 \
    --fees 100stake \
    --yes

sleep 2
log "✅ Curso de Ciência da Computação criado com sucesso!"

# =============================================================================
# FASE 3: CADASTRO DAS DISCIPLINAS (SUBJECTS)
# =============================================================================
log "📚 FASE 3: Cadastro das Disciplinas"

# Array com todas as disciplinas do curso
declare -A SUBJECTS

# 1º SEMESTRE
SUBJECTS["MAT101"]="Cálculo I|6|90|Introdução ao cálculo diferencial e integral de funções de uma variável real. Limites, derivadas, integrais e aplicações.|required|basic"
SUBJECTS["ALG101"]="Introdução à Programação|6|90|Fundamentos de programação utilizando linguagem C. Algoritmos, estruturas de controle, funções e ponteiros.|required|basic"
SUBJECTS["MAT102"]="Álgebra Linear|4|60|Sistemas lineares, matrizes, determinantes, espaços vetoriais e transformações lineares.|required|basic"
SUBJECTS["FIS101"]="Física Geral I|4|60|Mecânica clássica, leis de Newton, energia, momento e oscilações.|required|basic"
SUBJECTS["HUM101"]="Comunicação e Expressão|2|30|Desenvolvimento de habilidades de comunicação oral e escrita no contexto acadêmico e profissional.|required|basic"

# 2º SEMESTRE  
SUBJECTS["MAT201"]="Cálculo II|6|90|Cálculo de funções de várias variáveis, integrais múltiplas e cálculo vetorial.|required|basic"
SUBJECTS["ALG201"]="Estruturas de Dados I|6|90|Listas, pilhas, filas, árvores e algoritmos de ordenação e busca.|required|intermediate"
SUBJECTS["MAT202"]="Matemática Discreta|4|60|Lógica proposicional, teoria dos conjuntos, relações, funções e combinatória.|required|basic"
SUBJECTS["FIS201"]="Física Geral II|4|60|Eletromagnetismo, ondas e óptica geométrica.|required|basic"
SUBJECTS["HUM201"]="Filosofia da Ciência|2|30|Epistemologia, método científico e ética na ciência da computação.|required|basic"

# 3º SEMESTRE
SUBJECTS["ALG301"]="Programação Orientada a Objetos|6|90|Paradigma orientado a objetos usando Java. Classes, herança, polimorfismo e interfaces.|required|intermediate"
SUBJECTS["ALG302"]="Estruturas de Dados II|6|90|Árvores balanceadas, grafos, hashing e algoritmos avançados.|required|intermediate"
SUBJECTS["MAT301"]="Estatística e Probabilidade|4|60|Probabilidade, distribuições, inferência estatística e análise de dados.|required|intermediate"
SUBJECTS["ALG303"]="Arquitetura de Computadores|4|60|Organização de computadores, processadores, memória e sistemas de E/S.|required|intermediate"
SUBJECTS["HUM301"]="Inglês Técnico|2|30|Leitura e compreensão de textos técnicos em inglês na área de computação.|required|basic"

# 4º SEMESTRE
SUBJECTS["ALG401"]="Análise de Algoritmos|6|90|Complexidade computacional, análise assintótica e técnicas de projeto de algoritmos.|required|advanced"
SUBJECTS["ALG402"]="Banco de Dados I|6|90|Modelagem conceitual, modelo relacional, SQL e normalização.|required|intermediate"
SUBJECTS["ALG403"]="Sistemas Operacionais|6|90|Processos, threads, sincronização, gerenciamento de memória e sistemas de arquivos.|required|intermediate"
SUBJECTS["MAT401"]="Métodos Numéricos|4|60|Métodos computacionais para resolução de problemas matemáticos.|required|intermediate"
SUBJECTS["ALG404"]="Redes de Computadores I|4|60|Protocolos de comunicação, modelo OSI/TCP-IP e tecnologias de rede.|required|intermediate"

# 5º SEMESTRE
SUBJECTS["ALG501"]="Engenharia de Software I|6|90|Processo de desenvolvimento, análise de requisitos, design e padrões de projeto.|required|advanced"
SUBJECTS["ALG502"]="Compiladores|6|90|Análise léxica, sintática e semântica. Geração de código e otimização.|required|advanced"
SUBJECTS["ALG503"]="Inteligência Artificial I|6|90|Busca, representação do conhecimento e sistemas especialistas.|required|advanced"
SUBJECTS["ALG504"]="Interface Humano-Computador|4|60|Design de interfaces, usabilidade e experiência do usuário.|required|intermediate"
SUBJECTS["ELE501"]="Computação Gráfica|4|60|Primitivas gráficas, transformações geométricas e renderização.|elective|advanced"

# 6º SEMESTRE
SUBJECTS["ALG601"]="Banco de Dados II|6|90|Transações, controle de concorrência, bancos distribuídos e NoSQL.|required|advanced"
SUBJECTS["ALG602"]="Sistemas Distribuídos|6|90|Comunicação entre processos, tolerância a falhas e sistemas peer-to-peer.|required|advanced"
SUBJECTS["ALG603"]="Segurança Computacional|6|90|Criptografia, autenticação, autorização e segurança em sistemas.|required|advanced"
SUBJECTS["ALG604"]="Engenharia de Software II|4|60|Testes de software, manutenção e evolução de sistemas.|required|advanced"
SUBJECTS["ELE601"]="Processamento de Imagens|4|60|Filtragem, segmentação e reconhecimento de padrões em imagens.|elective|advanced"

# 7º SEMESTRE
SUBJECTS["ALG701"]="Projeto de Software|8|120|Desenvolvimento completo de um projeto de software em equipe.|required|advanced"
SUBJECTS["ALG702"]="Inteligência Artificial II|6|90|Machine learning, redes neurais e deep learning.|required|advanced"
SUBJECTS["ALG703"]="Computação Móvel|4|60|Desenvolvimento para dispositivos móveis e computação ubíqua.|required|advanced"
SUBJECTS["ELE701"]="Bioinformática|4|60|Aplicações computacionais em biologia e análise de sequências.|elective|advanced"
SUBJECTS["ELE702"]="Realidade Virtual|4|60|Ambientes virtuais, interação 3D e aplicações imersivas.|elective|advanced"

# 8º SEMESTRE
SUBJECTS["ALG801"]="Trabalho de Conclusão de Curso|12|180|Desenvolvimento de projeto final de graduação com orientação.|required|advanced"
SUBJECTS["ALG802"]="Empreendedorismo Tecnológico|4|60|Inovação, startups e gestão de negócios tecnológicos.|required|intermediate"
SUBJECTS["ELE801"]="Computação Quântica|4|60|Fundamentos da computação quântica e algoritmos quânticos.|elective|advanced"
SUBJECTS["ELE802"]="Blockchain e Criptomoedas|4|60|Tecnologias distribuídas, consenso e aplicações blockchain.|elective|advanced"
SUBJECTS["ELE803"]="IoT e Sistemas Embarcados|4|60|Internet das Coisas, sensores e sistemas de tempo real.|elective|advanced"

# Função para cadastrar uma disciplina
create_subject() {
    local code=$1
    local data=$2
    
    IFS='|' read -r name credits hours description type difficulty <<< "$data"
    
    info "Cadastrando disciplina: $code - $name"
    
    $BINARY tx subject create-subject \
        "$code" \
        "$name" \
        "UFT" \
        "CC-UFT-2024" \
        "$credits" \
        "$hours" \
        "$description" \
        "pt-BR" \
        "$type" \
        "$difficulty" \
        "active" \
        --from $UNIVERSITY_ADDR \
        --keyring-backend $KEYRING \
        --home $HOME_DIR \
        --chain-id $CHAIN_ID \
        --gas auto \
        --gas-adjustment 1.5 \
        --fees 100stake \
        --yes
    
    sleep 1  # Pequeno delay para não sobrecarregar
}

# Cadastrar todas as disciplinas
log "Iniciando cadastro de disciplinas..."
for code in "${!SUBJECTS[@]}"; do
    create_subject "$code" "${SUBJECTS[$code]}"
done

log "✅ Todas as disciplinas foram cadastradas com sucesso!"

# =============================================================================
# FASE 4: DEFINIÇÃO DE PRÉ-REQUISITOS
# =============================================================================
log "🔗 FASE 4: Definição de Pré-requisitos"

# Função para adicionar pré-requisito
add_prerequisite() {
    local subject=$1
    local prereq_type=$2  # "all" ou "any"
    local prereqs=$3
    
    info "Adicionando pré-requisito $prereq_type para $subject: $prereqs"
    
    $BINARY tx subject add-prerequisite-group \
        "$subject" \
        "$prereq_type" \
        "$prereqs" \
        --from $UNIVERSITY_ADDR \
        --keyring-backend $KEYRING \
        --home $HOME_DIR \
        --chain-id $CHAIN_ID \
        --gas auto \
        --gas-adjustment 1.5 \
        --fees 100stake \
        --yes
    
    sleep 1
}

# Definir pré-requisitos (exemplos dos principais)
info "Definindo pré-requisitos das disciplinas..."

# 2º Semestre
add_prerequisite "MAT201" "all" "MAT101"  # Cálculo II precisa de Cálculo I
add_prerequisite "ALG201" "all" "ALG101"  # Estruturas de Dados I precisa de Intro Programação
add_prerequisite "FIS201" "all" "FIS101,MAT101"  # Física II precisa de Física I e Cálculo I

# 3º Semestre  
add_prerequisite "ALG301" "all" "ALG201"  # POO precisa de Estruturas I
add_prerequisite "ALG302" "all" "ALG201,MAT202"  # Estruturas II precisa de Estruturas I e Mat. Discreta

# 4º Semestre
add_prerequisite "ALG401" "all" "ALG302,MAT202"  # Análise de Algoritmos
add_prerequisite "ALG402" "all" "ALG201"  # Banco de Dados I
add_prerequisite "ALG403" "all" "ALG302"  # Sistemas Operacionais

# 5º Semestre
add_prerequisite "ALG501" "all" "ALG301,ALG402"  # Engenharia de Software I
add_prerequisite "ALG502" "all" "ALG401"  # Compiladores
add_prerequisite "ALG503" "all" "ALG401,MAT301"  # IA I

# 6º Semestre
add_prerequisite "ALG601" "all" "ALG402"  # Banco de Dados II
add_prerequisite "ALG602" "all" "ALG403,ALG404"  # Sistemas Distribuídos
add_prerequisite "ALG603" "all" "ALG403,MAT202"  # Segurança
add_prerequisite "ALG604" "all" "ALG501"  # Engenharia de Software II

# 7º Semestre
add_prerequisite "ALG701" "all" "ALG501,ALG604"  # Projeto de Software
add_prerequisite "ALG702" "all" "ALG503,MAT301"  # IA II

# 8º Semestre
add_prerequisite "ALG801" "all" "ALG701"  # TCC precisa do Projeto

log "✅ Pré-requisitos definidos com sucesso!"

# =============================================================================
# FASE 5: CRIAÇÃO DA GRADE CURRICULAR (CURRICULUM)
# =============================================================================
log "📋 FASE 5: Criação da Grade Curricular"

info "Criando árvore curricular do curso..."

$BINARY tx curriculum create-curriculum-tree \
    "CURR-CC-UFT-2024" \
    "CC-UFT-2024" \
    "UFT" \
    "2024.1" \
    "active" \
    "Grade curricular do curso de Ciência da Computação 2024 - 8 semestres com 240 créditos totais" \
    --from $UNIVERSITY_ADDR \
    --keyring-backend $KEYRING \
    --home $HOME_DIR \
    --chain-id $CHAIN_ID \
    --gas auto \
    --gas-adjustment 1.5 \
    --fees 100stake \
    --yes

sleep 2

# Função para adicionar semestre ao currículo
add_semester() {
    local semester=$1
    local subjects=$2
    
    info "Adicionando $semester ao currículo com disciplinas: $subjects"
    
    $BINARY tx curriculum add-semester-to-curriculum \
        "CURR-CC-UFT-2024" \
        "$semester" \
        "$subjects" \
        --from $UNIVERSITY_ADDR \
        --keyring-backend $KEYRING \
        --home $HOME_DIR \
        --chain-id $CHAIN_ID \
        --gas auto \
        --gas-adjustment 1.5 \
        --fees 100stake \
        --yes
    
    sleep 1
}

# Adicionar semestres com suas disciplinas
add_semester "1" "MAT101,ALG101,MAT102,FIS101,HUM101"
add_semester "2" "MAT201,ALG201,MAT202,FIS201,HUM201"  
add_semester "3" "ALG301,ALG302,MAT301,ALG303,HUM301"
add_semester "4" "ALG401,ALG402,ALG403,MAT401,ALG404"
add_semester "5" "ALG501,ALG502,ALG503,ALG504,ELE501"
add_semester "6" "ALG601,ALG602,ALG603,ALG604,ELE601"
add_semester "7" "ALG701,ALG702,ALG703,ELE701,ELE702"
add_semester "8" "ALG801,ALG802,ELE801,ELE802,ELE803"

# Adicionar grupo de disciplinas eletivas
info "Adicionando grupo de disciplinas eletivas..."

$BINARY tx curriculum add-elective-group \
    "CURR-CC-UFT-2024" \
    "Eletivas Técnicas" \
    "ELE501,ELE601,ELE701,ELE702,ELE801,ELE802,ELE803" \
    "3" \
    "Disciplinas eletivas da área técnica - aluno deve cursar pelo menos 3" \
    --from $UNIVERSITY_ADDR \
    --keyring-backend $KEYRING \
    --home $HOME_DIR \
    --chain-id $CHAIN_ID \
    --gas auto \
    --gas-adjustment 1.5 \
    --fees 100stake \
    --yes

sleep 2

# Definir requisitos de graduação
info "Definindo requisitos de graduação..."

$BINARY tx curriculum set-graduation-requirements \
    "CURR-CC-UFT-2024" \
    "240" \
    "190" \
    "30" \
    "20" \
    "8.0" \
    "ALG801" \
    "Para se formar, o aluno deve: 1) Completar 240 créditos totais; 2) Completar 190 créditos obrigatórios; 3) Completar 30 créditos eletivos; 4) Completar 20 créditos complementares; 5) Ter média geral >= 8.0; 6) Concluir o TCC (ALG801)" \
    --from $UNIVERSITY_ADDR \
    --keyring-backend $KEYRING \
    --home $HOME_DIR \
    --chain-id $CHAIN_ID \
    --gas auto \
    --gas-adjustment 1.5 \
    --fees 100stake \
    --yes

sleep 2
log "✅ Grade curricular criada com sucesso!"

# =============================================================================
# FASE 6: CRIAÇÃO DE TOKEN DEFINITIONS
# =============================================================================
log "🪙 FASE 6: Criação de Token Definitions"

# Função para criar token definition
create_token_def() {
    local subject_id=$1
    local token_name=$2
    
    info "Criando token definition para $subject_id: $token_name"
    
    $BINARY tx tokendef create-token-definition \
        "TOKEN-$subject_id" \
        "$token_name" \
        "$subject_id" \
        "UFT" \
        "CC-UFT-2024" \
        "academic_subject" \
        "Token NFT representando a conclusão da disciplina $subject_id do curso de Ciência da Computação da UFT" \
        --from $UNIVERSITY_ADDR \
        --keyring-backend $KEYRING \
        --home $HOME_DIR \
        --chain-id $CHAIN_ID \
        --gas auto \
        --gas-adjustment 1.5 \
        --fees 100stake \
        --yes
    
    sleep 1
}

# Criar tokens para algumas disciplinas principais (exemplo)
info "Criando token definitions para disciplinas principais..."

create_token_def "ALG101" "Token Introdução à Programação"
create_token_def "ALG201" "Token Estruturas de Dados I"
create_token_def "ALG301" "Token Programação Orientada a Objetos"
create_token_def "ALG401" "Token Análise de Algoritmos"
create_token_def "ALG501" "Token Engenharia de Software I"
create_token_def "ALG601" "Token Banco de Dados II"
create_token_def "ALG701" "Token Projeto de Software"
create_token_def "ALG801" "Token Trabalho de Conclusão de Curso"

log "✅ Token definitions criados com sucesso!"

# =============================================================================
# RELATÓRIO FINAL
# =============================================================================
log "📊 RELATÓRIO FINAL"
log "=================="

echo ""
echo -e "${PURPLE}🎓 UNIVERSIDADE FEDERAL DE TECNOLOGIA (UFT) - SETUP COMPLETO${NC}"
echo ""
echo -e "${BLUE}📋 Resumo do que foi criado:${NC}"
echo "• 1 Instituição: Universidade Federal de Tecnologia (UFT)"
echo "• 1 Curso: Bacharelado em Ciência da Computação"  
echo "• 40 Disciplinas distribuídas em 8 semestres"
echo "• Sistema completo de pré-requisitos"
echo "• 1 Grade curricular estruturada (CURR-CC-UFT-2024)"
echo "• Grupo de disciplinas eletivas"
echo "• Requisitos de graduação definidos"
echo "• 8 Token definitions para disciplinas principais"
echo ""
echo -e "${GREEN}✅ Sistema institucional configurado e pronto para receber alunos!${NC}"
echo ""
echo -e "${YELLOW}📌 Próximos passos:${NC}"
echo "1. Criar script para jornada do aluno"
echo "2. Testar matrícula de estudantes"
echo "3. Simular conclusão de disciplinas"
echo "4. Testar emissão de NFTs acadêmicos"
echo "5. Testar processo de graduação"
echo ""
echo -e "${BLUE}🔧 Para consultar dados criados:${NC}"
echo "• Instituições: $BINARY q institution list-institution"
echo "• Cursos: $BINARY q course list-course"  
echo "• Disciplinas: $BINARY q subject list-subject"
echo "• Grade curricular: $BINARY q curriculum show-curriculum-tree CURR-CC-UFT-2024"
echo ""

log "🎉 Teste do Fluxo Institucional concluído com sucesso!"