#!/bin/bash

echo "=========================================="
echo "Academic Token - FIXED Test Setup Script"
echo "=========================================="

# Add some debugging
set -e  # Exit on first error
export DEBUG=1

echo ""
echo "0. CLEANING UP PREVIOUS DATA..."
echo "=========================================="

# Optional: Reset the chain data if needed
# academictokend unsafe-reset-all --home ~/.academictoken

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
    --yes || echo "Authorization failed for institution-$i"
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
academictokend tx course create-course \
  "institution-2" \
  "Medicina" \
  "MED" \
  "Bacharelado em Medicina com formação médica completa" \
  "360" \
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
echo "3. CREATING SUBJECTS CAREFULLY..."
echo "=========================================="

echo "Testing single subject creation first..."

# Test with minimal subject first
echo "Creating Cálculo I (test)..."
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
  --yes || echo "First subject creation failed"

sleep 5

echo "Checking if subject was created successfully..."
academictokend query subject list-subjects

echo ""
echo "Creating second subject..."
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
  --yes || echo "Second subject creation failed"

sleep 5

echo "Creating third subject..."
academictokend tx subject create-subject \
  "institution-2" \
  "course-2" \
  "Anatomia Humana" \
  "ANAT1" \
  120 \
  8 \
  "Estudo da estrutura do corpo humano" \
  "required" \
  "Medicine" \
  --from alice \
  --gas auto \
  --gas-adjustment 2.0 \
  --fees 2000stake \
  --yes || echo "Third subject creation failed"

sleep 5

echo ""
echo "Querying all subjects:"
academictokend query subject list-subjects

echo ""
echo "4. CREATING CURRICULUM TREES WITH UNIQUE VERSIONS..."
echo "=========================================="

echo "Creating curriculum tree for Computer Science..."
academictokend tx curriculum create-curriculum-tree \
    "course-1" \
    "2025.1" \
    "20" \
    "3600" \
    --from alice \
    --gas auto \
    --gas-adjustment 1.5 \
    --fees 1000stake \
    --chain-id academictoken \
    --yes || echo "Curriculum creation failed for course-1"

sleep 3

echo "Creating curriculum tree for Medicine..."
academictokend tx curriculum create-curriculum-tree \
    "course-2" \
    "2025.1" \
    "30" \
    "5400" \
    --from alice \
    --gas auto \
    --gas-adjustment 1.5 \
    --fees 1000stake \
    --chain-id academictoken \
    --yes || echo "Curriculum creation failed for course-2"

sleep 3

echo ""
echo "Querying curriculum trees:"
academictokend query curriculum list-curriculum-trees

echo ""
echo "5. CREATING TOKEN DEFINITIONS ONLY FOR EXISTING SUBJECTS..."
echo "=========================================="

echo "Checking which subjects exist first..."
SUBJECTS=$(academictokend query subject list-subjects --output json | jq -r '.subjects[].index' 2>/dev/null || echo "")

echo "Found subjects: $SUBJECTS"

# Only create tokens for subjects that actually exist
for subject_id in $SUBJECTS; do
  echo "Creating token definition for $subject_id..."
  
  case $subject_id in
    "subject-1")
      academictokend tx tokendef create-token-definition \
        "$subject_id" \
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
        --yes || echo "Token creation failed for $subject_id"
      ;;
    "subject-2")
      academictokend tx tokendef create-token-definition \
        "$subject_id" \
        "Programming I Achievement" \
        "PROG1-TKN" \
        "Token for Programming I course completion" \
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
        --yes || echo "Token creation failed for $subject_id"
      ;;
    "subject-3")
      academictokend tx tokendef create-token-definition \
        "$subject_id" \
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
        --yes || echo "Token creation failed for $subject_id"
      ;;
  esac
  
  sleep 3
done

echo ""
echo "6. QUERYING FINAL STATE..."
echo "=========================================="

echo "All token definitions:"
academictokend query tokendef list-token-definitions

echo ""
echo "Final subjects list:"
academictokend query subject list-subjects

echo ""
echo "=========================================="
echo "FIXED SETUP COMPLETED!"
echo "=========================================="

echo ""
echo "🎯 RESULTS SUMMARY:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

INSTITUTION_COUNT=$(academictokend query institution list-institutions --output json 2>/dev/null | jq '.institution | length' 2>/dev/null || echo "0")
COURSE_COUNT=$(academictokend query course list-courses --output json 2>/dev/null | jq '.course | length' 2>/dev/null || echo "0")
SUBJECT_COUNT=$(academictokend query subject list-subjects --output json 2>/dev/null | jq '.subjects | length' 2>/dev/null || echo "0")
TOKEN_COUNT=$(academictokend query tokendef list-token-definitions --output json 2>/dev/null | jq '.tokenDefinitions | length' 2>/dev/null || echo "0")

echo "🏛️  INSTITUTIONS: $INSTITUTION_COUNT created and authorized"
echo "🎓 COURSES: $COURSE_COUNT created"
echo "📚 SUBJECTS: $SUBJECT_COUNT created successfully"
echo "🪙  TOKEN DEFINITIONS: $TOKEN_COUNT created"

echo ""
echo "✅ IF SUBJECT COUNT > 1: System is working properly!"
echo "❌ IF SUBJECT COUNT = 1: There are still issues with subject creation"
echo ""
echo "🚀 Next step: If subjects are working, run the full script again!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
