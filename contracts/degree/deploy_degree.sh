#!/bin/bash

# Configurações
CHAIN_ID="academictoken"
ALICE_ADDR=$(academictokend keys show alice -a)
WASM_FILE="artifacts/degree_opt.wasm"

echo "Deploying with Alice address: $ALICE_ADDR"

# 1. Armazenar o contrato
echo "Storing contract..."
STORE_TX=$(academictokend tx wasm store $WASM_FILE \
  --from alice \
  --chain-id $CHAIN_ID \
  --gas auto \
  --gas-adjustment 1.3 \
  --fees 5000stake \
  -y \
  --output json | jq -r '.txhash')

echo "Store TX: $STORE_TX"
sleep 6

# 2. Obter o code_id
CODE_ID=$(academictokend query tx $STORE_TX --output json | jq -r '.logs[0].events[] | select(.type=="store_code") | .attributes[] | select(.key=="code_id") | .value')
echo "Code ID: $CODE_ID"

# 3. Instanciar
INIT_MSG="{\"module_address\":\"$ALICE_ADDR\"}"
echo "Instantiating with: $INIT_MSG"

INIT_TX=$(academictokend tx wasm instantiate $CODE_ID "$INIT_MSG" \
  --label "degree-contract-v1" \
  --from alice \
  --admin $ALICE_ADDR \
  --chain-id $CHAIN_ID \
  --gas auto \
  --gas-adjustment 1.3 \
  --fees 5000stake \
  -y \
  --output json | jq -r '.txhash')

echo "Init TX: $INIT_TX"
sleep 6

# 4. Obter endereço do contrato
CONTRACT_ADDR=$(academictokend query tx $INIT_TX --output json | jq -r '.logs[0].events[] | select(.type=="instantiate") | .attributes[] | select(.key=="_contract_address") | .value')
echo "Contract Address: $CONTRACT_ADDR"

# 5. Salvar informações
echo "CODE_ID=$CODE_ID" > degree_contract_info.txt
echo "CONTRACT_ADDR=$CONTRACT_ADDR" >> degree_contract_info.txt
echo "ADMIN=$ALICE_ADDR" >> degree_contract_info.txt

echo "Done! Contract deployed at: $CONTRACT_ADDR"

