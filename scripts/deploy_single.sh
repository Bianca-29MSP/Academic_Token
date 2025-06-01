#!/bin/bash

# Academic Token - Deploy Single Contract Script

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Check if contract name provided
if [ $# -eq 0 ]; then
    echo -e "${RED}❌ Please provide contract name${NC}"
    echo "Usage: ./deploy_single.sh <contract_name> [wallet_name]"
    echo "Available contracts: prerequisites, schedule, progress, equivalence, degree"
    exit 1
fi

CONTRACT_NAME=$1
WALLET_NAME=${2:-"deployer"}
CHAIN_ID="academictoken-1"
GAS_PRICES="0.025utoken"

echo -e "${BLUE}🚀 Deploying $CONTRACT_NAME contract${NC}"
echo ""

# Navigate to contract directory
cd contracts/$CONTRACT_NAME

# 1. Build
echo -e "${YELLOW}📦 Building...${NC}"
cargo build --release --target wasm32-unknown-unknown

# 2. Store
echo -e "${YELLOW}📤 Storing...${NC}"
STORE_RESULT=$(academictokend tx wasm store \
    target/wasm32-unknown-unknown/release/${CONTRACT_NAME}.wasm \
    --from $WALLET_NAME \
    --chain-id $CHAIN_ID \
    --gas-prices $GAS_PRICES \
    --gas auto \
    --gas-adjustment 1.3 \
    --output json \
    --yes)

TX_HASH=$(echo $STORE_RESULT | jq -r '.txhash')
echo "Store TX: $TX_HASH"

sleep 6

CODE_ID=$(academictokend query tx $TX_HASH --output json | jq -r '.logs[0].events[] | select(.type=="store_code") | .attributes[] | select(.key=="code_id") | .value')

echo -e "${GREEN}✅ Code ID: $CODE_ID${NC}"

# 3. Instantiate
echo -e "${YELLOW}🎯 Instantiating...${NC}"

case $CONTRACT_NAME in
    "prerequisites")
        INIT_MSG='{"owner": null}'
        ;;
    "schedule")
        INIT_MSG='{"owner": null}'
        ;;
    "progress")
        INIT_MSG='{"owner": null, "analytics_enabled": true, "update_frequency": "Daily", "analytics_depth": "Standard"}'
        ;;
    "equivalence")
        INIT_MSG='{"owner": null, "similarity_threshold": 80, "auto_approval_threshold": 95}'
        ;;
    "degree")
        INIT_MSG='{"owner": null}'
        ;;
esac

INSTANTIATE_RESULT=$(academictokend tx wasm instantiate \
    $CODE_ID \
    "$INIT_MSG" \
    --from $WALLET_NAME \
    --label "${CONTRACT_NAME}-v1" \
    --chain-id $CHAIN_ID \
    --gas-prices $GAS_PRICES \
    --gas auto \
    --gas-adjustment 1.3 \
    --output json \
    --yes)

TX_HASH=$(echo $INSTANTIATE_RESULT | jq -r '.txhash')
echo "Instantiate TX: $TX_HASH"

sleep 6

CONTRACT_ADDRESS=$(academictokend query tx $TX_HASH --output json | jq -r '.logs[0].events[] | select(.type=="instantiate") | .attributes[] | select(.key=="_contract_address") | .value')

echo ""
echo -e "${GREEN}🎉 SUCCESS!${NC}"
echo -e "${GREEN}Contract: $CONTRACT_NAME${NC}"
echo -e "${GREEN}Code ID: $CODE_ID${NC}"
echo -e "${GREEN}Address: $CONTRACT_ADDRESS${NC}"

# Save to file
echo "$CONTRACT_NAME:$CODE_ID:$CONTRACT_ADDRESS" >> ../../deployed_contracts.log

cd ../..
