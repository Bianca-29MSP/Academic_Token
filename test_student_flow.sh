#!/bin/bash
CONTRACT_ADDR="academic1ufs3tlq4umljk0qfe8k5ya0x6hpavn897u2cnf9k0en9jr7qarqqj7temq"
SENDER="academic1kxc8unwj9zrtsvvtu0pys34he5ht5x08t35wln"

echo "👨‍🎓 TESTANDO FLUXO COMPLETO DE ESTUDANTE"
echo "========================================"


# ESTUDANTE 1: João - Iniciante
echo "1️⃣ JOÃO - Completou Álgebra Linear"
academictokend tx wasm execute $CONTRACT_ADDR \
  '{"update_student_record":{"student_id":"joao","completed_subject":{"subject_id":"ALG001","credits":4,"completion_date":"2024-01-15","grade":8500,"nft_token_id":"nft-alg-joao"}}}' \
  --from $SENDER --gas auto --gas-adjustment 1.3 --yes

sleep 3

# Verificar elegibilidade para Estatística (requer ALG001 OU CALC1)
echo "�� João pode cursar Estatística?"
academictokend query wasm contract-state smart $CONTRACT_ADDR \
  '{"check_eligibility":{"student_id":"joao","subject_id":"STAT001"}}'

sleep 3

# ESTUDANTE 2: Maria - Intermediária  
echo ""
echo "2️⃣ MARIA - Completou Cálculo I e II"
academictokend tx wasm execute $CONTRACT_ADDR \
  '{"update_student_record":{"student_id":"maria","completed_subject":{"subject_id":"CALC1","credits":6,"completion_date":"2024-02-01","grade":9000,"nft_token_id":"nft-calc1-maria"}}}' \
  --from $SENDER --gas auto --gas-adjustment 1.3 --yes

sleep 3

academictokend tx wasm execute $CONTRACT_ADDR \
  '{"update_student_record":{"student_id":"maria","completed_subject":{"subject_id":"CALC2-NEW","credits":6,"completion_date":"2024-06-01","grade":8800,"nft_token_id":"nft-calc2-maria"}}}' \
  --from $SENDER --gas auto --gas-adjustment 1.3 --yes

sleep 3

# Verificar elegibilidade para Cálculo III (requer CALC2-NEW)
echo "📐 Maria pode cursar Cálculo III?"
academictokend query wasm contract-state smart $CONTRACT_ADDR \
  '{"check_eligibility":{"student_id":"maria","subject_id":"CALC3"}}'

sleep 3

# ESTUDANTE 3: Pedro - Avançado
echo ""
echo "3️⃣ PEDRO - Estudante avançado (para testar TCC)"
# Adicionar múltiplas disciplinas para Pedro
academictokend tx wasm execute $CONTRACT_ADDR \
  '{"update_student_record":{"student_id":"pedro","completed_subject":{"subject_id":"CALC1","credits":6,"completion_date":"2023-03-01","grade":9500,"nft_token_id":"nft-calc1-pedro"}}}' \
  --from $SENDER --gas auto --gas-adjustment 1.3 --yes

sleep 3

academictokend tx wasm execute $CONTRACT_ADDR \
  '{"update_student_record":{"student_id":"pedro","completed_subject":{"subject_id":"ALG001","credits":4,"completion_date":"2023-03-01","grade":9200,"nft_token_id":"nft-alg-pedro"}}}' \
  --from $SENDER --gas auto --gas-adjustment 1.3 --yes

sleep 3

# Verificar elegibilidade para Equações Diferenciais (requer CALC2-NEW E ALG001)
echo "🔬 Pedro pode cursar Equações Diferenciais?"
academictokend query wasm contract-state smart $CONTRACT_ADDR \
  '{"check_eligibility":{"student_id":"pedro","subject_id":"DIFF001"}}'

sleep 3

# VERIFICAÇÕES FINAIS
echo ""
echo "📋 HISTÓRICOS DOS ESTUDANTES:"
echo "----------------------------"
echo "👨‍🎓 João:"
academictokend query wasm contract-state smart $CONTRACT_ADDR '{"get_student_record":{"student_id":"joao"}}'

echo ""
echo "👩‍🎓 Maria:"
academictokend query wasm contract-state smart $CONTRACT_ADDR '{"get_student_record":{"student_id":"maria"}}'

echo ""
echo "👨‍💼 Pedro:"
academictokend query wasm contract-state smart $CONTRACT_ADDR '{"get_student_record":{"student_id":"pedro"}}'

echo ""
echo "📊 ESTADO FINAL:"
academictokend query wasm contract-state smart $CONTRACT_ADDR '{"get_state":{}}'

echo ""
echo "🎉 TESTE COMPLETO FINALIZADO!"
