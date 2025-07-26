#!/bin/bash

echo "🚀 Criando dados de teste para pré-requisitos..."

# Aguardar entre comandos para evitar problemas de sequência
wait_cmd() {
    sleep 2
}

echo "📚 Criando disciplinas..."

# Cálculo I
academictokend tx subject create-subject-content \
  --index="CALC1" --subject-id="CALC1" --institution="institution-1" --course-id="course-1" \
  --title="Cálculo I" --code="MAT001" --workload-hours=60 --credits=4 \
  --description="Introdução ao cálculo" --subject-type="Obrigatória" --knowledge-area="Matemática" \
  --from=alice --chain-id=academictoken --yes
wait_cmd

# Cálculo II
academictokend tx subject create-subject-content \
  --index="CALC2" --subject-id="CALC2" --institution="institution-1" --course-id="course-1" \
  --title="Cálculo II" --code="MAT002" --workload-hours=60 --credits=4 \
  --description="Cálculo avançado" --subject-type="Obrigatória" --knowledge-area="Matemática" \
  --from=alice --chain-id=academictoken --yes
wait_cmd

echo "🔗 Criando pré-requisitos..."

# Pré-requisito: CALC2 requer CALC1
academictokend tx subject add-prerequisite-group \
  --subject-id="CALC2" --group-type="ALL" --minimum-credits=0 --minimum-completed-subjects=0 \
  --subject-ids="CALC1" --from=alice --chain-id=academictoken --yes
wait_cmd

echo "✅ Dados criados! Aguardando processamento..."
sleep 5

echo "🧪 Testando endpoint..."
curl http://localhost:1317/academictoken/subject/prerequisites/course/course-1
