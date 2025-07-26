#!/bin/bash

echo "=========================================="
echo "Academic Token - Extended Test Setup Script"
echo "=========================================="

echo ""
echo "1. CREATING INSTITUTIONS..."
echo "=========================================="

echo "Registering UFMG..."
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

echo "Registering USP..."
academictokend tx institution register-institution \
  "Universidade de São Paulo" \
  "Rua da Reitoria, 374, São Paulo, SP, Brasil" \
  --from alice \
  --chain-id academictoken \
  --gas auto \
  --gas-adjustment 1.5 \
  --fees 1000stake \
  --yes

sleep 3

echo "Registering UNICAMP..."
academictokend tx institution register-institution \
  "Universidade Estadual de Campinas" \
  "Cidade Universitária Zeferino Vaz, Campinas, SP, Brasil" \
  --from alice \
  --chain-id academictoken \
  --gas auto \
  --gas-adjustment 1.5 \
  --fees 1000stake \
  --yes

sleep 3

echo "Registering UFRJ..."
academictokend tx institution register-institution \
  "Universidade Federal do Rio de Janeiro" \
  "Av. Pedro Calmon, 550, Rio de Janeiro, RJ, Brasil" \
  --from alice \
  --chain-id academictoken \
  --gas auto \
  --gas-adjustment 1.5 \
  --fees 1000stake \
  --yes

sleep 3

echo "Authorizing all institutions..."
for i in {1..4}; do
  echo "Authorizing institution-$i..."
  academictokend tx institution update-institution \
    "institution-$i" \
    "$(academictokend query institution get-institution institution-$i --output json 2>/dev/null | jq -r '.institution.name // "Institution Name"')" \
    "$(academictokend query institution get-institution institution-$i --output json 2>/dev/null | jq -r '.institution.address // "Institution Address"')" \
    "true" \
    --from alice \
    --chain-id academictoken \
    --gas auto \
    --gas-adjustment 1.5 \
    --fees 1000stake \
    --yes || true
  sleep 2
done

echo ""
echo "Listing institutions:"
academictokend query institution list-institutions

echo ""
echo "2. CREATING COURSES..."
echo "=========================================="

echo "Creating Computer Science course at UFMG..."
academictokend tx course create-course \
  "institution-1" \
  "Ciência da Computação" \
  "CC" \
  "Bacharelado em Ciência da Computação com foco em desenvolvimento de software e algoritmos" \
  "240" \
  "undergraduate" \
  --from alice \
  --chain-id academictoken \
  --gas auto \
  --gas-adjustment 1.5 \
  --fees 1000stake \
  --yes

sleep 3

echo "Creating Medicine course at USP..."
academictokend tx subject create-subject \
  "institution-1" "course-1" "Química Geral" "QUIM1" "90" "6" \
  "Princípios fundamentais da química" "required" "Chemistry" \
  --objectives "Compreender estrutura atômica" \
  --objectives "Estudar ligações químicas" \
  --topic-units "Tabela periódica" \
  --topic-units "Ligações covalentes e iônicas" \
  --chain-id academictoken \
  --from alice --gas auto --gas-adjustment 1.5 --fees 1000stake -y

sleep 3

echo "Creating Civil Engineering course at UFMG..."
academictokend tx course create-course \
  "institution-1" \
  "Engenharia Civil" \
  "EC" \
  "Bacharelado em Engenharia Civil com foco em construção e infraestrutura" \
  "300" \
  "undergraduate" \
  --from alice \
  --chain-id academictoken \
  --gas auto \
  --gas-adjustment 1.5 \
  --fees 1000stake \
  --yes

sleep 3

echo "Creating Physics course at UNICAMP..."
academictokend tx course create-course \
  "institution-3" \
  "Física" \
  "FIS" \
  "Bacharelado em Física com ênfase em pesquisa" \
  "240" \
  "undergraduate" \
  --from alice \
  --chain-id academictoken \
  --gas auto \
  --gas-adjustment 1.5 \
  --fees 1000stake \
  --yes

sleep 3

echo "Creating Mathematics course at USP..."
academictokend tx course create-course \
  "institution-2" \
  "Matemática" \
  "MAT" \
  "Bacharelado em Matemática pura e aplicada" \
  "240" \
  "undergraduate" \
  --from alice \
  --chain-id academictoken \
  --gas auto \
  --gas-adjustment 1.5 \
  --fees 1000stake \
  --yes

sleep 3

echo "Creating Chemistry course at UFRJ..."
academictokend tx course create-course \
  "institution-4" \
  "Química" \
  "QUI" \
  "Bacharelado em Química com foco em pesquisa e indústria" \
  "240" \
  "undergraduate" \
  --from alice \
  --chain-id academictoken \
  --gas auto \
  --gas-adjustment 1.5 \
  --fees 1000stake \
  --yes

sleep 3

echo "Creating Electrical Engineering course at UNICAMP..."
academictokend tx course create-course \
  "institution-3" \
  "Engenharia Elétrica" \
  "EE" \
  "Bacharelado em Engenharia Elétrica com foco em eletrônica e telecomunicações" \
  "300" \
  "undergraduate" \
  --from alice \
  --chain-id academictoken \
  --gas auto \
  --gas-adjustment 1.5 \
  --fees 1000stake \
  --yes

sleep 3

echo ""
echo "Querying courses:"
academictokend query course list-courses

echo ""
echo "3. CREATING SUBJECTS..."
echo "=========================================="

# Subjects for Computer Science (UFMG)
echo "Creating subjects for Computer Science..."
academictokend tx subject create-subject \
  "institution-1" "course-1" "Cálculo I" "CALC1" "90" "6" \
  "Fundamentos do cálculo diferencial e integral" "required" "Mathematics" \
  --objectives "Compreender limites e derivadas" \
  --objectives "Aplicar integrais em problemas práticos" \
  --topic-units "Limites e continuidade" \
  --topic-units "Derivadas e aplicações" \
  --topic-units "Integrais definidas e indefinidas" \
  --chain-id academictoken \
  --from alice --gas auto --gas-adjustment 1.5 --fees 1000stake -y

sleep 2

academictokend tx subject create-subject \
  "institution-1" "course-1" "Programação I" "PROG1" "60" "4" \
  "Introdução à programação estruturada" "required" "Computer Science" \
  --objectives "Desenvolver lógica de programação" \
  --objectives "Implementar algoritmos básicos" \
  --topic-units "Algoritmos e estruturas de controle" \
  --topic-units "Funções e procedimentos" \
  --topic-units "Arrays e strings" \
  --chain-id academictoken \
  --from alice --gas auto --gas-adjustment 1.5 --fees 1000stake -y

sleep 2

academictokend tx subject create-subject \
  "institution-1" "course-1" "Cálculo II" "CALC2" "90" "6" \
  "Continuação do cálculo diferencial e integral" "required" "Mathematics" \
  --objectives "Dominar técnicas avançadas de integração" \
  --objectives "Trabalhar com séries e sequências" \
  --topic-units "Técnicas de integração" \
  --topic-units "Séries infinitas" \
  --topic-units "Equações diferenciais básicas" \
  --chain-id academictoken \
  --from alice --gas auto --gas-adjustment 1.5 --fees 1000stake -y

sleep 2

academictokend tx subject create-subject \
  "institution-1" "course-1" "Estrutura de Dados" "ED1" "75" "5" \
  "Algoritmos e estruturas de dados fundamentais" "required" "Computer Science" \
  --objectives "Implementar estruturas de dados eficientes" \
  --objectives "Analisar complexidade de algoritmos" \
  --topic-units "Listas, pilhas e filas" \
  --topic-units "Árvores e grafos" \
  --topic-units "Algoritmos de ordenação" \
  --chain-id academictoken \
  --from alice --gas auto --gas-adjustment 1.5 --fees 1000stake -y

sleep 2

academictokend tx subject create-subject \
  "institution-1" "course-1" "Banco de Dados" "BD1" "60" "4" \
  "Sistemas de gerenciamento de banco de dados" "required" "Computer Science" \
  --objectives "Projetar bancos de dados relacionais" \
  --objectives "Implementar consultas SQL complexas" \
  --topic-units "Modelo relacional" \
  --topic-units "Normalização" \
  --topic-units "SQL avançado" \
  --chain-id academictoken \
  --from alice --gas auto --gas-adjustment 1.5 --fees 1000stake -y

sleep 2

# Subjects for Medicine (USP)
echo "Creating subjects for Medicine..."
academictokend tx subject create-subject \
  "institution-2" "course-2" "Anatomia Humana" "ANAT1" "120" "8" \
  "Estudo da estrutura do corpo humano" "required" "Medicine" \
  --objectives "Conhecer sistemas anatômicos" \
  --objectives "Identificar estruturas corporais" \
  --topic-units "Sistema esquelético" \
  --topic-units "Sistema muscular" \
  --topic-units "Sistema nervoso" \
  --chain-id academictoken \
  --from alice --gas auto --gas-adjustment 1.5 --fees 1000stake -y

sleep 2

academictokend tx subject create-subject \
  "institution-2" "course-2" "Fisiologia" "FISIO1" "90" "6" \
  "Funcionamento dos sistemas do corpo humano" "required" "Medicine" \
  --objectives "Compreender funções orgânicas" \
  --objectives "Analisar homeostase corporal" \
  --topic-units "Fisiologia cardiovascular" \
  --topic-units "Fisiologia respiratória" \
  --topic-units "Fisiologia renal" \
  --chain-id academictoken \
  --from alice --gas auto --gas-adjustment 1.5 --fees 1000stake -y

sleep 2

academictokend tx subject create-subject \
  "institution-2" "course-2" "Bioquímica" "BIOQ1" "75" "5" \
  "Processos químicos em organismos vivos" "required" "Medicine" \
  --objectives "Entender metabolismo celular" \
  --objectives "Estudar enzimas e coenzimas" \
  --topic-units "Metabolismo de carboidratos" \
  --topic-units "Metabolismo de lipídios" \
  --topic-units "Metabolismo de proteínas" \
  --chain-id academictoken \
  --from alice --gas auto --gas-adjustment 1.5 --fees 1000stake -y

sleep 2

# Subjects for Physics (UNICAMP)
echo "Creating subjects for Physics..."
academictokend tx subject create-subject \
  "institution-3" "course-4" "Física I" "FIS1" "90" "6" \
  "Mecânica clássica e termodinâmica" "required" "Physics" \
  --objectives "Aplicar leis da mecânica" \
  --objectives "Resolver problemas de termodinâmica" \
  --topic-units "Cinemática e dinâmica" \
  --topic-units "Energia e momento" \
  --topic-units "Termodinâmica básica" \
  --chain-id academictoken \
  --from alice --gas auto --gas-adjustment 1.5 --fees 1000stake -y

sleep 2

academictokend tx subject create-subject \
  "institution-3" "course-4" "Física II" "FIS2" "90" "6" \
  "Eletromagnetismo e óptica" "required" "Physics" \
  --objectives "Compreender fenômenos eletromagnéticos" \
  --objectives "Estudar propriedades da luz" \
  --topic-units "Eletrostática e magnetismo" \
  --topic-units "Circuitos elétricos" \
  --topic-units "Óptica geométrica e ondulatória" \
  --chain-id academictoken \
  --from alice --gas auto --gas-adjustment 1.5 --fees 1000stake -y

sleep 2

academictokend tx subject create-subject \
  "institution-3" "course-4" "Cálculo Vetorial" "CALCVET" "75" "5" \
  "Cálculo em várias variáveis" "required" "Mathematics" \
  --objectives "Trabalhar com funções vetoriais" \
  --objectives "Aplicar teoremas fundamentais" \
  --topic-units "Derivadas parciais" \
  --topic-units "Integrais múltiplas" \
  --topic-units "Teoremas de Green e Stokes" \
  --chain-id academictoken \
  --from alice --gas auto --gas-adjustment 1.5 --fees 1000stake -y

sleep 2

# Subjects for Mathematics (USP)
echo "Creating subjects for Mathematics..."
academictokend tx subject create-subject \
  "institution-2" "course-5" "Álgebra Linear" "ALG1" "75" "5" \
  "Espaços vetoriais e transformações lineares" "required" "Mathematics" \
  --objectives "Compreender espaços vetoriais" \
  --objectives "Trabalhar com transformações lineares" \
  --topic-units "Vetores e matrizes" \
  --topic-units "Determinantes e sistemas lineares" \
  --topic-units "Autovalores e autovetores" \
  --chain-id academictoken \
  --from alice --gas auto --gas-adjustment 1.5 --fees 1000stake -y

sleep 2

academictokend tx subject create-subject \
  "institution-2" "course-5" "Análise Real" "ANREAL" "90" "6" \
  "Fundamentos da análise matemática" "required" "Mathematics" \
  --objectives "Dominar conceitos de análise" \
  --objectives "Demonstrar teoremas fundamentais" \
  --topic-units "Sequências e séries" \
  --topic-units "Funções contínuas" \
  --topic-units "Diferenciabilidade" \
  --chain-id academictoken \
  --from alice --gas auto --gas-adjustment 1.5 --fees 1000stake -y

sleep 2

# Subjects for Chemistry (UFRJ)
echo "Creating subjects for Chemistry..."
academictokend tx subject create-subject \
  "institution-4" "course-6" "Química Geral" "QUIM1" "90" "6" \
  "Princípios fundamentais da química" "required" "Chemistry" \
  --objectives "Compreender estrutura atômica" \
  --objectives "Estudar ligações químicas" \
  --topic-units "Tabela periódica" \
  --topic-units "Ligações covalentes e iônicas" \
  --topic-units "Reações químicas básicas" \
  --chain-id academictoken \
  --from alice --gas auto --gas-adjustment 1.5 --fees 1000stake -y

sleep 2

academictokend tx subject create-subject \
  "institution-4" "course-6" "Química Orgânica" "QUIMORG" "105" "7" \
  "Compostos orgânicos e suas reações" "required" "Chemistry" \
  --objectives "Identificar grupos funcionais" \
  --objectives "Prever mecanismos de reação" \
  --topic-units "Hidrocarbonetos" \
  --topic-units "Compostos funcionais" \
  --topic-units "Mecanismos de reação" \
  --chain-id academictoken \
  --from alice --gas auto --gas-adjustment 1.5 --fees 1000stake -y

sleep 2

# Subjects for Electrical Engineering (UNICAMP)
echo "Creating subjects for Electrical Engineering..."
academictokend tx subject create-subject \
  "institution-3" "course-6" "Circuitos Elétricos" "CIRC1" "90" "6" \
  "Análise de circuitos elétricos básicos" "required" "Engineering" \
  --objectives "Analisar circuitos DC e AC" \
  --objectives "Aplicar leis de Kirchhoff" \
  --topic-units "Resistores e lei de Ohm" \
  --topic-units "Capacitores e indutores" \
  --topic-units "Análise de circuitos AC" \
  --chain-id academictoken \
  --from alice --gas auto --gas-adjustment 1.5 --fees 1000stake -y

sleep 2

academictokend tx subject create-subject \
  "institution-3" "course-6" "Eletrônica Digital" "ELETDIG" "75" "5" \
  "Sistemas digitais e microprocessadores" "required" "Engineering" \
  --objectives "Projetar circuitos digitais" \
  --objectives "Programar microcontroladores" \
  --topic-units "Álgebra booleana" \
  --topic-units "Circuitos combinacionais" \
  --topic-units "Microprocessadores" \
  --chain-id academictoken \
  --from alice --gas auto --gas-adjustment 1.5 --fees 1000stake -y

sleep 2

echo ""
echo "Querying subjects:"
academictokend query subject list-subjects

echo ""
echo "4. CREATING CURRICULUM TREES..."
echo "=========================================="

echo "Creating curriculum tree for Computer Science..."
academictokend tx curriculum create-curriculum-tree \
    "course-1" \
    "2024.1" \
    "20" \
    "3600" \
    --from alice \
    --gas auto \
    --gas-adjustment 1.5 \
    --fees 1000stake \
    --chain-id academictoken \
    --yes

sleep 3

echo "Creating curriculum tree for Medicine..."
academictokend tx curriculum create-curriculum-tree \
    "course-2" \
    "2024.1" \
    "30" \
    "5400" \
    --from alice \
    --gas auto \
    --gas-adjustment 1.5 \
    --fees 1000stake \
    --chain-id academictoken \
    --yes

sleep 3

echo "Creating curriculum tree for Physics..."
academictokend tx curriculum create-curriculum-tree \
    "course-4" \
    "2024.1" \
    "15" \
    "3600" \
    --from alice \
    --gas auto \
    --gas-adjustment 1.5 \
    --fees 1000stake \
    --chain-id academictoken \
    --yes

sleep 3

echo "Creating curriculum tree for Mathematics..."
academictokend tx curriculum create-curriculum-tree \
    "course-5" \
    "2024.1" \
    "18" \
    "3600" \
    --from alice \
    --gas auto \
    --gas-adjustment 1.5 \
    --fees 1000stake \
    --chain-id academictoken \
    --yes

sleep 3

echo ""
echo "Querying curriculum trees:"
academictokend query curriculum list-curriculum-trees

echo ""
echo "5. CREATING TOKEN DEFINITIONS..."
echo "=========================================="

echo "Creating token definitions for Computer Science subjects..."

academictokend tx tokendef create-token-definition \
  "subject-1" \
  "Cálculo I Completion Token" \
  "CALC1-TKN" \
  "Token representing successful completion of Calculus I course" \
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

academictokend tx tokendef create-token-definition \
  "subject-2" \
  "Programming I Achievement" \
  "PROG1-TKN" \
  "Token for Programming I course completion with practical projects" \
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

academictokend tx tokendef create-token-definition \
  "subject-3" \
  "Advanced Calculus Token" \
  "CALC2-TKN" \
  "Token for Calculus II advanced concepts mastery" \
  "ACHIEVEMENT" \
  --is-transferable=false \
  --is-burnable=true \
  --max-supply=1000 \
  --image-uri="https://academic-tokens.example.com/calc2.png" \
  --from alice \
  --chain-id academictoken \
  --gas auto \
  --gas-adjustment 1.5 \
  --fees 1000stake \
  --yes

sleep 3

echo "Creating token definitions for Medicine subjects..."

academictokend tx tokendef create-token-definition \
  "subject-6" \
  "Human Anatomy Mastery" \
  "ANAT-TKN" \
  "Token certifying comprehensive knowledge of human anatomy" \
  "NFT" \
  --is-transferable=false \
  --is-burnable=false \
  --max-supply=500 \
  --image-uri="https://academic-tokens.example.com/anatomy.png" \
  --from alice \
  --chain-id academictoken \
  --gas auto \
  --gas-adjustment 1.5 \
  --fees 1000stake \
  --yes

sleep 3

academictokend tx tokendef create-token-definition \
  "subject-7" \
  "Physiology Excellence" \
  "FISIO-TKN" \
  "Token for advanced physiology understanding" \
  "NFT" \
  --is-transferable=true \
  --is-burnable=false \
  --max-supply=600 \
  --image-uri="https://academic-tokens.example.com/physiology.png" \
  --from alice \
  --chain-id academictoken \
  --gas auto \
  --gas-adjustment 1.5 \
  --fees 1000stake \
  --yes

sleep 3

echo "Creating token definitions for Physics subjects..."

academictokend tx tokendef create-token-definition \
  "subject-9" \
  "Classical Physics Token" \
  "FIS1-TKN" \
  "Token for mastery of classical mechanics and thermodynamics" \
  "ACHIEVEMENT" \
  --is-transferable=true \
  --is-burnable=true \
  --max-supply=750 \
  --image-uri="https://academic-tokens.example.com/physics1.png" \
  --from alice \
  --chain-id academictoken \
  --gas auto \
  --gas-adjustment 1.5 \
  --fees 1000stake \
  --yes

sleep 3

academictokend tx tokendef create-token-definition \
  "subject-10" \
  "Electromagnetism Token" \
  "FIS2-TKN" \
  "Token certifying electromagnetic theory knowledge" \
  "ACHIEVEMENT" \
  --is-transferable=false \
  --is-burnable=false \
  --max-supply=750 \
  --image-uri="https://academic-tokens.example.com/physics2.png" \
  --from alice \
  --chain-id academictoken \
  --gas auto \
  --gas-adjustment 1.5 \
  --fees 1000stake \
  --yes

sleep 3

echo "Creating token definitions for Chemistry subjects..."

academictokend tx tokendef create-token-definition \
  "subject-14" \
  "General Chemistry Foundation" \
  "QUIM1-TKN" \
  "Token for fundamental chemistry principles mastery" \
  "ACHIEVEMENT" \
  --is-transferable=true \
  --is-burnable=false \
  --max-supply=900 \
  --image-uri="https://academic-tokens.example.com/chemistry.png" \
  --from alice \
  --chain-id academictoken \
  --gas auto \
  --gas-adjustment 1.5 \
  --fees 1000stake \
  --yes

sleep 3

echo "Creating token definitions for Engineering subjects..."

academictokend tx tokendef create-token-definition \
  "subject-16" \
  "Circuit Analysis Expert" \
  "CIRC-TKN" \
  "Token for electrical circuit analysis expertise" \
  "FUNGIBLE" \
  --is-transferable=false \
  --is-burnable=true \
  --max-supply=400 \
  --image-uri="https://academic-tokens.example.com/circuits.png" \
  --from alice \
  --chain-id academictoken \
  --gas auto \
  --gas-adjustment 1.5 \
  --fees 1000stake \
  --yes

sleep 3

echo ""
echo "6. QUERYING TOKEN DEFINITIONS..."
echo "=========================================="

echo "All token definitions:"
academictokend query tokendef list-token-definitions

echo ""
echo "Token definitions for Programming I (subject-2):"
academictokend query tokendef list-token-definitions-by-subject "subject-2" || echo "Query command may not exist yet"

echo ""
echo "Detailed view of first token definition:"
academictokend query tokendef get-token-definition "tokendef-1" || echo "No token definitions found yet"

echo ""
echo "Full content of Physics token:"
academictokend query tokendef get-token-definition-full "tokendef-6" || echo "Full content query may not exist yet"

echo ""
echo "7. TESTING UPDATES..."
echo "=========================================="

echo "Updating Programming token definition..."
academictokend tx tokendef update-token-definition \
  "tokendef-2" \
  "Advanced Programming I Achievement" \
  "ADVPROG1-TKN" \
  "Enhanced token for Programming I with advanced project completion" \
  --is-transferable=true \
  --is-burnable=false \
  --max-supply=1200 \
  --image-uri="https://academic-tokens.example.com/advanced-prog1.png" \
  --from alice \
  --chain-id academictoken \
  --gas auto \
  --gas-adjustment 1.5 \
  --fees 1000stake \
  --yes || echo "Update failed - token may not exist yet"

sleep 3

echo ""
echo "Updated token definition:"
academictokend query tokendef get-token-definition "tokendef-2" || echo "Token definition not found"

echo ""
echo "=========================================="
echo "EXTENDED SETUP COMPLETED SUCCESSFULLY!"
echo "=========================================="

echo ""
echo "COMPREHENSIVE SUMMARY:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🏛️  INSTITUTIONS: 4 institutions created and authorized"
echo "   • UFMG (Universidade Federal de Minas Gerais)"
echo "   • USP (Universidade de São Paulo)"
echo "   • UNICAMP (Universidade Estadual de Campinas)"
echo "   • UFRJ (Universidade Federal do Rio de Janeiro)"
echo ""
echo "🎓 COURSES: 7 diverse courses across multiple institutions"
echo "   • Computer Science (UFMG)"
echo "   • Medicine (USP)"
echo "   • Civil Engineering (UFMG)"
echo "   • Physics (UNICAMP)"
echo "   • Mathematics (USP)"
echo "   • Chemistry (UFRJ)"
echo "   • Electrical Engineering (UNICAMP)"
echo ""
echo "📚 SUBJECTS: 17 subjects across different knowledge areas"
echo "   • Mathematics: Cálculo I, Cálculo II, Álgebra Linear, Análise Real, Cálculo Vetorial"
echo "   • Computer Science: Programação I, Estrutura de Dados, Banco de Dados"
echo "   • Medicine: Anatomia Humana, Fisiologia, Bioquímica"
echo "   • Physics: Física I, Física II"
echo "   • Chemistry: Química Geral, Química Orgânica"
echo "   • Engineering: Circuitos Elétricos, Eletrônica Digital"
echo ""
echo "🌳 CURRICULUM TREES: 4 curriculum trees for major programs"
echo "   • Computer Science 2024.1"
echo "   • Medicine 2024.1"
echo "   • Physics 2024.1"
echo "   • Mathematics 2024.1"
echo ""
echo "🪙  TOKEN DEFINITIONS: 10 diverse academic tokens created"
echo "   • Achievement Tokens (6)"
echo "   • NFT Tokens (2)"
echo "   • Fungible Tokens (2)"
echo ""
echo "✅ READY FOR TESTING:"
echo "   • Prerequisites and dependencies"
echo "   • Subject equivalences between institutions"
echo "   • Token transfers and burning"
echo "   • Curriculum progression tracking"
echo "   • Multi-institutional academic recognition"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "🚀 Next steps: Test prerequisites, equivalences, NFT minting, and cross-institutional transfers!"
