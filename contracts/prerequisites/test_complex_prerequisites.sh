#!/bin/bash

# Script Corrigido: Pré-requisitos com delays entre transações
CONTRACT_ADDR="academic14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9sva26tx"
SENDER="academic1yrq6z28t2gxekm74n2xlegfvzzaz2adkyr5ntt"

echo "🎓 TESTES DE PRÉ-REQUISITOS COMPLEXOS (Versão Corrigida)"
echo "======================================================"

# Função para aguardar e verificar
wait_and_verify() {
    local subject_id=$1
    echo "⏳ Aguardando 3 segundos..."
    sleep 3
    echo "✅ Verificando: $subject_id"
    academictokend query wasm contract-state smart $CONTRACT_ADDR "{\"get_prerequisites\":{\"subject_id\":\"$subject_id\"}}"
    echo ""
}

# 1. Cálculo I (corrigir sequence)
echo "1️⃣ Registrando: Cálculo I"
academictokend tx wasm execute $CONTRACT_ADDR \
  '{"register_prerequisites":{"subject_id":"CALC1","prerequisites":[{"id":"calc1-none","subject_id":"CALC1","group_type":"none","minimum_credits":0,"minimum_completed_subjects":0,"subject_ids":[],"logic":"none","priority":1,"confidence":100}]}}' \
  --from $SENDER --gas auto --gas-adjustment 1.3 --yes

wait_and_verify "CALC1"

# 2. Cálculo II → requer Cálculo I
echo "2️⃣ Registrando: Cálculo II → requer Cálculo I"
academictokend tx wasm execute $CONTRACT_ADDR \
  '{"register_prerequisites":{"subject_id":"CALC2-NEW","prerequisites":[{"id":"calc2-req","subject_id":"CALC2-NEW","group_type":"all","minimum_credits":0,"minimum_completed_subjects":1,"subject_ids":["CALC1"],"logic":"and","priority":1,"confidence":100}]}}' \
  --from $SENDER --gas auto --gas-adjustment 1.3 --yes

wait_and_verify "CALC2-NEW"

# 3. Cálculo III → requer Cálculo II (JSON CORRIGIDO)
echo "3️⃣ Registrando: Cálculo III → requer Cálculo II"
academictokend tx wasm execute $CONTRACT_ADDR \
  '{"register_prerequisites":{"subject_id":"CALC3","prerequisites":[{"id":"calc3-req","subject_id":"CALC3","group_type":"all","minimum_credits":0,"minimum_completed_subjects":1,"subject_ids":["CALC2-NEW"],"logic":"and","priority":1,"confidence":100}]}}' \
  --from $SENDER --gas auto --gas-adjustment 1.3 --yes

wait_and_verify "CALC3"

# 4. Estatística → Cálculo I OU Álgebra Linear
echo "4️⃣ Registrando: Estatística → Cálculo I OU Álgebra Linear"
academictokend tx wasm execute $CONTRACT_ADDR \
  '{"register_prerequisites":{"subject_id":"STAT001","prerequisites":[{"id":"stat-alt","subject_id":"STAT001","group_type":"any","minimum_credits":0,"minimum_completed_subjects":1,"subject_ids":["CALC1","ALG001"],"logic":"or","priority":1,"confidence":90}]}}' \
  --from $SENDER --gas auto --gas-adjustment 1.3 --yes

wait_and_verify "STAT001"

# 5. Equações Diferenciais → Cálculo II E Álgebra Linear
echo "5️⃣ Registrando: Equações Diferenciais → Cálculo II E Álgebra"
academictokend tx wasm execute $CONTRACT_ADDR \
  '{"register_prerequisites":{"subject_id":"DIFF001","prerequisites":[{"id":"diff-multi","subject_id":"DIFF001","group_type":"all","minimum_credits":0,"minimum_completed_subjects":2,"subject_ids":["CALC2-NEW","ALG001"],"logic":"and","priority":1,"confidence":95}]}}' \
  --from $SENDER --gas auto --gas-adjustment 1.3 --yes

wait_and_verify "DIFF001"

# 6. TCC → mínimo 120 créditos
echo "6️⃣ Registrando: TCC → mínimo 120 créditos"
academictokend tx wasm execute $CONTRACT_ADDR \
  '{"register_prerequisites":{"subject_id":"TCC001","prerequisites":[{"id":"tcc-credits","subject_id":"TCC001","group_type":"minimum","minimum_credits":120,"minimum_completed_subjects":0,"subject_ids":[],"logic":"and","priority":1,"confidence":100}]}}' \
  --from $SENDER --gas auto --gas-adjustment 1.3 --yes

wait_and_verify "TCC001"

# FASE 7: VERIFICAR REGISTROS
echo "🔍 FASE 7: Verificando Registros"
echo "-------------------------------"

subjects=("ALG001" "CALC1" "CALC2" "CALC3" "STAT001" "PHYS301" "DIFF001" "NUMER001" "TCC001" "ELEC-ADV" "MEST001")

for subject in "${subjects[@]}"; do
    echo "📋 Verificando: $subject"
    academictokend query wasm contract-state smart $CONTRACT_ADDR "{\"get_prerequisites\":{\"subject_id\":\"$subject\"}}"
    echo ""
done

echo "📊 ESTADO FINAL:"
academictokend query wasm contract-state smart $CONTRACT_ADDR '{"get_state":{}}'

echo ""
echo "🎉 SCRIPT CORRIGIDO CONCLUÍDO!"
