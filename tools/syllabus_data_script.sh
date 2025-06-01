#!/bin/bash

# =============================================================================
# ACADEMIC TOKEN - SYLLABUS DATA SCRIPT
# =============================================================================
# Este script complementa o university_test_script.sh adicionando ementas
# detalhadas para todas as disciplinas do curso de Ciência da Computação.
# Simula a integração com IPFS para armazenamento de conteúdo extenso.
# =============================================================================

set -e

# Cores para output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}"
}

info() {
    echo -e "${BLUE}[INFO] $1${NC}"
}

# =============================================================================
# CONFIGURAÇÕES
# =============================================================================
CHAIN_ID="academictoken"
KEYRING="test"
UNIVERSITY_KEY="university"
HOME_DIR="$HOME/.academictoken"
BINARY="academictokend"

UNIVERSITY_ADDR=$($BINARY keys show $UNIVERSITY_KEY -a --keyring-backend $KEYRING --home $HOME_DIR)

log "📚 Iniciando Atualização de Ementas Detalhadas"
log "============================================="

# =============================================================================
# FUNÇÃO PARA ATUALIZAR CONTEÚDO DA DISCIPLINA
# =============================================================================
update_subject_content() {
    local subject_code=$1
    local ementa=$2
    local objetivos=$3
    local topicos=$4
    local bibliografia=$5
    local metodologia=$6
    local avaliacao=$7
    
    info "Atualizando ementa de $subject_code..."
    
    # Simular hash IPFS (em produção seria o hash real do conteúdo no IPFS)
    local ipfs_hash="Qm$(echo -n "$subject_code$ementa" | sha256sum | cut -c1-44)"
    local ipfs_link="https://ipfs.io/ipfs/$ipfs_hash"
    
    $BINARY tx subject update-subject-content \
        "$subject_code" \
        "$ementa" \
        "$objetivos" \
        "$topicos" \
        "$bibliografia" \
        "$metodologia" \
        "$avaliacao" \
        "$ipfs_hash" \
        "$ipfs_link" \
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

# =============================================================================
# EMENTAS DETALHADAS - 1º SEMESTRE
# =============================================================================
log "📖 Atualizando ementas do 1º semestre..."

update_subject_content "MAT101" \
    "Estudo do cálculo diferencial e integral de funções de uma variável real. Conceitos fundamentais de limite, continuidade, derivabilidade e integrabilidade. Aplicações em problemas práticos de otimização e análise de funções." \
    "Desenvolver competências em cálculo diferencial e integral para resolução de problemas matemáticos e aplicações em computação. Estabelecer bases matemáticas sólidas para disciplinas subsequentes." \
    "Números reais e funções; Limites e continuidade; Derivadas e aplicações; Máximos e mínimos; Integrais definidas e indefinidas; Teorema fundamental do cálculo; Aplicações da integral; Técnicas de integração" \
    "STEWART, J. Cálculo Vol 1. 8ª ed. Cengage Learning, 2017; GUIDORIZZI, H. L. Um Curso de Cálculo Vol 1. 5ª ed. LTC, 2013; ANTON, H. Cálculo Vol 1. 10ª ed. Bookman, 2014" \
    "Aulas expositivas com demonstrações teóricas; Resolução de exercícios práticos; Laboratórios de aplicação computacional; Estudos dirigidos individuais" \
    "Provas escritas (60%); Listas de exercícios (25%); Trabalhos práticos (15%)"

update_subject_content "ALG101" \
    "Introdução aos conceitos fundamentais de programação de computadores. Desenvolvimento da lógica de programação através da linguagem C. Estruturas básicas, algoritmos e resolução de problemas computacionais." \
    "Desenvolver raciocínio lógico-matemático para resolução de problemas; Dominar os fundamentos da programação estruturada; Implementar algoritmos básicos em linguagem C; Preparar base para disciplinas avançadas de programação." \
    "Introdução à computação e algoritmos; Linguagem C: sintaxe e semântica; Tipos de dados e operadores; Estruturas de controle (if, while, for); Funções e procedimentos; Arrays e strings; Ponteiros e alocação de memória; Arquivos; Debugging e boas práticas" \
    "KERNIGHAN, B.; RITCHIE, D. C: A Linguagem de Programação. 2ª ed. Campus, 1989; DEITEL, P. C: Como Programar. 8ª ed. Pearson, 2015; ASCENCIO, A. F. G. Fundamentos da Programação de Computadores. 3ª ed. Pearson, 2012" \
    "Aulas teórico-práticas em laboratório; Programação individual e em duplas; Projetos práticos incrementais; Code review e pair programming" \
    "Projetos de programação (50%); Provas práticas (30%); Exercícios em laboratório (20%)"

update_subject_content "MAT102" \
    "Estudo de sistemas lineares, matrizes, determinantes, espaços vetoriais e transformações lineares. Aplicações em computação gráfica, machine learning e otimização." \
    "Compreender estruturas algébricas lineares fundamentais; Resolver sistemas lineares e trabalhar com matrizes; Aplicar conceitos em problemas computacionais; Desenvolver intuição geométrica para espaços vetoriais." \
    "Sistemas de equações lineares; Matrizes e operações matriciais; Determinantes e propriedades; Espaços vetoriais e subespaços; Dependência e independência linear; Base e dimensão; Transformações lineares; Autovalores e autovetores; Aplicações computacionais" \
    "ANTON, H. Álgebra Linear Contemporânea. 10ª ed. Bookman, 2012; LAY, D. Álgebra Linear e suas Aplicações. 5ª ed. LTC, 2018; BOLDRINI, J. L. Álgebra Linear. 3ª ed. Harbra, 1986" \
    "Aulas expositivas com exemplos práticos; Resolução de problemas em grupo; Uso de software matemático (MATLAB/Octave); Aplicações em computação gráfica" \
    "Provas teóricas (50%); Trabalhos computacionais (30%); Exercícios aplicados (20%)"

# =============================================================================
# EMENTAS DETALHADAS - 2º SEMESTRE
# =============================================================================
log "📖 Atualizando ementas do 2º semestre..."

update_subject_content "ALG201" \
    "Estudo detalhado das principais estruturas de dados lineares e hierárquicas. Implementação e análise de algoritmos de manipulação, busca e ordenação. Introdução à análise de complexidade." \
    "Dominar estruturas de dados fundamentais; Implementar algoritmos eficientes de manipulação; Analisar complexidade computacional; Escolher estruturas adequadas para cada problema; Desenvolver habilidades de debugging e otimização." \
    "Revisão de ponteiros e alocação dinâmica; Listas ligadas (simples, dupla, circular); Pilhas e filas; Árvores binárias e traversal; Árvores de busca binária; Algoritmos de ordenação (bubble, selection, insertion, merge, quick); Algoritmos de busca; Análise de complexidade O(n); Hash tables básico" \
    "CORMEN, T. H. Algoritmos: Teoria e Prática. 3ª ed. Campus, 2012; TENENBAUM, A. M. Estruturas de Dados Usando C. 1ª ed. Pearson, 1995; SZWARCFITER, J. L. Estruturas de Dados e seus Algoritmos. 3ª ed. LTC, 2010" \
    "Aulas práticas em laboratório; Implementação incremental de estruturas; Projetos de aplicação real; Análise comparativa de algoritmos; Code review em grupo" \
    "Implementações práticas (60%); Provas de algoritmos (25%); Relatórios de análise (15%)"

update_subject_content "MAT202" \
    "Fundamentos da matemática discreta aplicada à ciência da computação. Lógica proposicional, teoria dos conjuntos, relações, funções, indução matemática e combinatória." \
    "Desenvolver raciocínio lógico-matemático rigoroso; Compreender fundamentos teóricos da computação; Aplicar métodos de prova matemática; Resolver problemas combinatórios; Preparar base para algoritmos e complexidade." \
    "Lógica proposicional e predicados; Métodos de demonstração; Teoria dos conjuntos; Relações e suas propriedades; Funções e cardinalidade; Indução matemática; Princípios de contagem; Combinações e permutações; Grafos básicos; Introdução à teoria dos números" \
    "ROSEN, K. H. Matemática Discreta e suas Aplicações. 7ª ed. McGraw Hill, 2013; GERSTING, J. L. Fundamentos Matemáticos para a Ciência da Computação. 7ª ed. LTC, 2017; SCHEINERMAN, E. R. Matemática Discreta. 3ª ed. Cengage Learning, 2016" \
    "Aulas expositivas com demonstrações; Resolução de problemas em classe; Exercícios de demonstração; Aplicações computacionais" \
    "Provas teóricas (70%); Exercícios de demonstração (20%); Trabalhos aplicados (10%)"

# =============================================================================
# EMENTAS DETALHADAS - 3º SEMESTRE
# =============================================================================
log "📖 Atualizando ementas do 3º semestre..."

update_subject_content "ALG301" \
    "Paradigma de programação orientada a objetos utilizando Java. Conceitos de encapsulamento, herança, polimorfismo e abstração. Design patterns e boas práticas de desenvolvimento." \
    "Dominar o paradigma orientado a objetos; Desenvolver aplicações usando Java; Aplicar princípios SOLID; Implementar design patterns; Compreender frameworks e APIs; Desenvolver sistemas modulares e reutilizáveis." \
    "Introdução à OOP e Java; Classes, objetos e métodos; Encapsulamento e modificadores de acesso; Herança e superclasses; Polimorfismo e sobrescrita; Classes abstratas e interfaces; Exceções e tratamento de erros; Collections Framework; Design patterns (Factory, Observer, Strategy); UML e modelagem; Introdução a frameworks; Testes unitários (JUnit)" \
    "DEITEL, P. Java: Como Programar. 10ª ed. Pearson, 2016; SIERRA, K. Use a Cabeça! Java. 2ª ed. Alta Books, 2007; BLOCH, J. Java Efetivo. 3ª ed. Alta Books, 2018; GAMMA, E. Padrões de Projeto. 1ª ed. Bookman, 2000" \
    "Aulas práticas em IDE (Eclipse/IntelliJ); Desenvolvimento de projetos incrementais; Code review e refactoring; Pair programming; Workshops de design patterns" \
    "Projetos orientados a objetos (50%); Implementação de patterns (25%); Provas práticas (25%)"

update_subject_content "ALG302" \
    "Estruturas de dados avançadas e algoritmos eficientes. Árvores balanceadas, grafos, hashing e algoritmos de otimização. Análise rigorosa de complexidade." \
    "Implementar estruturas de dados complexas; Analisar e otimizar algoritmos; Resolver problemas usando grafos; Compreender trade-offs de design; Aplicar técnicas avançadas de programação." \
    "Árvores AVL e rotações; Árvores Red-Black; B-Trees e aplicações; Hash tables avançado; Grafos: representação e algoritmos; Busca em largura (BFS) e profundidade (DFS); Algoritmos de caminho mínimo (Dijkstra, Floyd-Warshall); Árvore geradora mínima (Kruskal, Prim); Análise amortizada; Programação dinâmica introdutória" \
    "CORMEN, T. H. Algoritmos: Teoria e Prática. 3ª ed. Campus, 2012; SEDGEWICK, R. Algorithms. 4ª ed. Addison-Wesley, 2011; KLEINBERG, J. Algorithm Design. 1ª ed. Pearson, 2005" \
    "Implementação de estruturas complexas; Análise experimental de algoritmos; Projetos com grandes datasets; Competições de programação; Seminários de algoritmos" \
    "Implementações avançadas (40%); Análise de complexidade (30%); Projeto final (30%)"

# =============================================================================
# EMENTAS DETALHADAS - 4º SEMESTRE
# =============================================================================
log "📖 Atualizando ementas do 4º semestre..."

update_subject_content "ALG401" \
    "Análise rigorosa da complexidade de algoritmos. Técnicas avançadas de projeto de algoritmos: divisão e conquista, programação dinâmica, algoritmos gulosos e backtracking." \
    "Analisar complexidade temporal e espacial rigorosamente; Dominar técnicas clássicas de projeto; Desenvolver algoritmos eficientes; Compreender limites computacionais; Aplicar matemática avançada em análise." \
    "Análise assintótica avançada; Técnicas de divisão e conquista; Relações de recorrência; Programação dinâmica; Algoritmos gulosos e corretude; Backtracking e branch-and-bound; Algoritmos de ordenação avançados; Análise probabilística; Introdução à complexidade computacional; Lower bounds e optimalidade" \
    "CORMEN, T. H. Algoritmos: Teoria e Prática. 3ª ed. Campus, 2012; KLEINBERG, J. Algorithm Design. 1ª ed. Pearson, 2005; DASGUPTA, S. Algorithms. 1ª ed. McGraw-Hill, 2006" \
    "Análise matemática rigorosa; Implementação e benchmarking; Competições algorítmicas; Seminários de pesquisa; Projetos de otimização" \
    "Provas teóricas (50%); Implementações otimizadas (30%); Projeto de pesquisa (20%)"

update_subject_content "ALG402" \
    "Fundamentos de sistemas de banco de dados relacionais. Modelagem conceitual, modelo relacional, linguagem SQL e técnicas de normalização." \
    "Projetar bancos de dados eficientes; Dominar SQL avançado; Aplicar técnicas de normalização; Compreender arquitetura de SGBD; Desenvolver aplicações database-driven." \
    "Conceitos fundamentais de BD; Modelagem ER (Entidade-Relacionamento); Modelo relacional e álgebra relacional; SQL: DDL, DML, DCL e TCL; Consultas complexas e subqueries; Joins e operações de conjunto; Funções agregadas e grouping; Normalização (1FN, 2FN, 3FN, BCNF); Índices e otimização; Introdução a NoSQL; Conexão com aplicações (JDBC)" \
    "ELMASRI, R. Sistemas de Banco de Dados. 7ª ed. Pearson, 2018; DATE, C. J. Introdução a Sistemas de Bancos de Dados. 8ª ed. Campus, 2003; RAMAKRISHNAN, R. Database Management Systems. 3ª ed. McGraw-Hill, 2002" \
    "Modelagem de casos reais; Laboratórios com MySQL/PostgreSQL; Projetos de aplicação web; Otimização de consultas; Análise de performance" \
    "Projeto de modelagem (40%); Provas de SQL (35%); Laboratórios práticos (25%)"

# =============================================================================
# EMENTAS DETALHADAS - 5º SEMESTRE  
# =============================================================================
log "📖 Atualizando ementas do 5º semestre..."

update_subject_content "ALG501" \
    "Processo de desenvolvimento de software, análise de requisitos, design de sistemas e arquiteturas. Metodologias ágeis e tradicionais. Qualidade e métricas de software." \
    "Compreender o ciclo de vida do software; Aplicar metodologias de desenvolvimento; Realizar análise e especificação de requisitos; Projetar arquiteturas robustas; Gerenciar projetos de software." \
    "Engenharia de software e processo; Ciclo de vida do software; Metodologias: Waterfall, Agile, Scrum; Análise e especificação de requisitos; Casos de uso e user stories; Design de software e arquitetura; Padrões arquiteturais (MVC, Layered, Microservices); UML avançado; Métricas e qualidade; Gerenciamento de configuração; DevOps introdutório" \
    "SOMMERVILLE, I. Engenharia de Software. 10ª ed. Pearson, 2018; PRESSMAN, R. S. Engenharia de Software. 8ª ed. McGraw Hill, 2016; FOWLER, M. UML Essencial. 3ª ed. Bookman, 2005" \
    "Projeto de sistema completo; Metodologias ágeis na prática; Workshops de requisitos; Desenvolvimento em equipe; Apresentações executivas" \
    "Projeto de engenharia (50%); Documentação técnica (25%); Provas conceituais (25%)"

update_subject_content "ALG503" \
    "Fundamentos da inteligência artificial. Algoritmos de busca, representação do conhecimento, sistemas especialistas e introdução ao aprendizado de máquina." \
    "Compreender fundamentos teóricos da IA; Implementar algoritmos de busca; Desenvolver sistemas especialistas; Aplicar técnicas de representação; Preparar base para IA avançada." \
    "História e filosofia da IA; Agentes inteligentes; Busca cega (BFS, DFS, UCS); Busca heurística (A*, IDA*); Busca local e meta-heurísticas; Jogos e minimax; Representação do conhecimento; Lógica proposicional e predicados; Sistemas especialistas; Redes semânticas; Introdução ao machine learning; Algoritmos genéticos; Processamento de linguagem natural básico" \
    "RUSSELL, S. Inteligência Artificial. 3ª ed. Campus, 2013; LUGER, G. F. Inteligência Artificial. 6ª ed. Pearson, 2013; NORVIG, P. Paradigms of Artificial Intelligence Programming. 1ª ed. Morgan Kaufmann, 1992" \
    "Implementação de agentes; Projetos de busca; Desenvolvimento de sistema especialista; Competições de IA; Seminários de pesquisa" \
    "Implementações de algoritmos (40%); Sistema especialista (35%); Provas teóricas (25%)"

# =============================================================================
# EMENTAS DETALHADAS - 7º SEMESTRE
# =============================================================================
log "📖 Atualizando ementas do 7º semestre..."

update_subject_content "ALG701" \
    "Desenvolvimento completo de um projeto de software em equipe. Aplicação prática de metodologias de engenharia de software, gerenciamento de projeto e trabalho colaborativo." \
    "Integrar conhecimentos de engenharia de software; Trabalhar efetivamente em equipe; Gerenciar projetos complexos; Aplicar boas práticas de desenvolvimento; Entregar software funcionando." \
    "Definição de escopo e requisitos; Planejamento e estimativas; Arquitetura e design detalhado; Implementação colaborativa; Versionamento com Git; Integração contínua; Testes automatizados; Documentação técnica; Apresentação para stakeholders; Deploy e manutenção; Retrospectivas e lições aprendidas" \
    "SOMMERVILLE, I. Engenharia de Software. 10ª ed. Pearson, 2018; PMI. Guia PMBOK. 6ª ed. PMI, 2017; MARTIN, R. C. Código Limpo. 1ª ed. Alta Books, 2009" \
    "Desenvolvimento em equipes de 4-6 pessoas; Reuniões de acompanhamento semanais; Code review obrigatório; Apresentações quinzenais; Metodologia Scrum" \
    "Projeto final (70%); Documentação (15%); Apresentação (15%)"

update_subject_content "ALG702" \
    "Técnicas avançadas de inteligência artificial. Machine learning, redes neurais, deep learning e aplicações práticas em visão computacional e processamento de linguagem natural." \
    "Dominar algoritmos de machine learning; Implementar redes neurais; Aplicar deep learning; Desenvolver sistemas inteligentes; Compreender estado da arte em IA." \
    "Revisão de IA I; Aprendizado supervisionado; Regressão linear e logística; Árvores de decisão e random forests; SVM e kernel methods; Clustering e aprendizado não-supervisionado; Redes neurais artificiais; Backpropagation e otimização; Deep learning e CNNs; RNNs e processamento sequencial; Transfer learning; Aplicações práticas: visão computacional, NLP; Ética em IA" \
    "GOODFELLOW, I. Deep Learning. 1ª ed. MIT Press, 2016; MURPHY, K. P. Machine Learning: A Probabilistic Perspective. 1ª ed. MIT Press, 2012; RUSSELL, S. Inteligência Artificial. 3ª ed. Campus, 2013" \
    "Implementação com Python/TensorFlow; Projetos de aplicação real; Competições de ML; Papers de pesquisa; Seminários avançados" \
    "Implementações de ML (45%); Projeto de aplicação (35%); Paper review (20%)"

# =============================================================================
# EMENTAS DETALHADAS - 8º SEMESTRE
# =============================================================================
log "📖 Atualizando ementas do 8º semestre..."

update_subject_content "ALG801" \
    "Desenvolvimento de projeto final de graduação com orientação individual. Pesquisa, implementação e documentação de solução inovadora para problema real da área de computação." \
    "Demonstrar domínio técnico da área; Conduzir pesquisa científica; Implementar solução completa; Redigir documentação acadêmica; Apresentar resultados profissionalmente." \
    "Definição do tema e orientador; Revisão bibliográfica; Metodologia de pesquisa; Desenvolvimento da proposta; Implementação do sistema/algoritmo; Testes e validação; Redação da monografia; Preparação da apresentação; Defesa pública; Reflexão sobre contribuições" \
    "Bibliografia específica do tema escolhido; WAZLAWICK, R. S. Metodologia de Pesquisa para Ciência da Computação. 2ª ed. Campus, 2014; ECO, U. Como se faz uma tese. 26ª ed. Perspectiva, 2016" \
    "Orientação individual semanal; Seminários de acompanhamento; Apresentações parciais; Peer review entre alunos; Workshops de redação" \
    "Monografia (60%); Implementação técnica (25%); Apresentação final (15%)"

update_subject_content "ALG802" \
    "Conceitos fundamentais de empreendedorismo aplicados à área tecnológica. Inovação, modelos de negócio, startups e gestão de empresas de base tecnológica." \
    "Desenvolver visão empreendedora; Elaborar planos de negócio; Compreender ecossistema de startups; Aplicar conceitos de inovação; Preparar para liderança tecnológica." \
    "Empreendedorismo e inovação; Identificação de oportunidades; Lean startup e MVP; Modelos de negócio (Canvas); Validação de mercado; Financiamento e investimento; Propriedade intelectual; Marketing digital; Liderança e gestão de equipes; Aspectos legais; Pitch e apresentação; Cases de sucesso no Brasil" \
    "RIES, E. A Startup Enxuta. 1ª ed. Lua de Papel, 2012; OSTERWALDER, A. Business Model Generation. 1ª ed. Alta Books, 2011; BLANK, S. Manual do Empreendedor. 1ª ed. Alta Books, 2014" \
    "Desenvolvimento de plano de negócio; Pitches para investidores; Visitas a incubadoras; Palestras com empreendedores; Simulação de startup" \
    "Plano de negócio (50%); Pitch presentation (30%); Participação e networking (20%)"

log "✅ Todas as ementas foram atualizadas com sucesso!"

# =============================================================================
# RELATÓRIO DE CONTEÚDO
# =============================================================================
log "📊 RELATÓRIO DE CONTEÚDO ATUALIZADO"
log "==================================="

echo ""
echo -e "${BLUE}📚 Resumo das ementas adicionadas:${NC}"
echo "• Ementas completas com objetivos pedagógicos"
echo "• Conteúdo programático detalhado"
echo "• Bibliografia especializada atualizada"  
echo "• Metodologias de ensino específicas"
echo "• Critérios de avaliação definidos"
echo "• Simulação de hashes IPFS para cada disciplina"
echo ""
echo -e "${GREEN}✅ Integração IPFS simulada - Conteúdo extenso armazenado off-chain${NC}"
echo -e "${YELLOW}📌 Em produção, o conteúdo seria realmente armazenado no IPFS${NC}"
echo ""

log "🎉 Script de ementas detalhadas concluído!"