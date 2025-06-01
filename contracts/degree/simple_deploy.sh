#!/bin/bash

echo "=== DEGREE CONTRACT DEPLOYMENT ==="
echo ""

# 1. Get Alice address
ALICE_ADDR=$(academictokend keys show alice -a)
echo "Alice address: $ALICE_ADDR"

# 2. Store the contract
echo ""
echo "Storing contract..."
academictokend tx wasm store artifacts/degree_opt.wasm \
  --from alice \
  --gas auto \
  --gas-adjustment 1.3 \
  --fees 50000stake \
  --chain-id academictoken \
  -y

echo ""
echo "Waiting for transaction to be included..."
sleep 10

# 3. Get the latest code_id
echo ""
echo "Getting code ID..."
CODE_ID=$(academictokend query wasm list-code --output json | jq -r '.code_infos[-1].code_id')
echo "Code ID: $CODE_ID"

if [ -z "$CODE_ID" ]; then
    echo "ERROR: Could not get code ID"
    exit 1
fi

# 4. Instantiate the contract
echo ""
echo "Instantiating contract..."
INIT_MSG="{\"module_address\":\"$ALICE_ADDR\"}"

academictokend tx wasm instantiate $CODE_ID "$INIT_MSG" \
  --label "degree-contract" \
  --from alice \
  --admin $ALICE_ADDR \
  --gas auto \
  --gas-adjustment 1.3 \
  --fees 50000stake \
  --chain-id academictoken \
  -y

echo ""
echo "Waiting for instantiation..."
sleep 10

# 5. Get contract address
echo ""
echo "Getting contract address..."
CONTRACT_ADDR=$(academictokend query wasm list-contract-by-code $CODE_ID --output json | jq -r '.contracts[0]')
echo "Contract address: $CONTRACT_ADDR"

if [ -z "$CONTRACT_ADDR" ]; then
    echo "ERROR: Could not get contract address"
    exit 1
fi

# 6. Save contract info
echo "CODE_ID=$CODE_ID" > degree_contract_info.txt
echo "CONTRACT_ADDR=$CONTRACT_ADDR" >> degree_contract_info.txt
echo "ADMIN=$ALICE_ADDR" >> degree_contract_info.txt

echo ""
echo "=== DEPLOYMENT COMPLETE ==="
echo "Contract address: $CONTRACT_ADDR"
echo "Info saved to: degree_contract_info.txt"
