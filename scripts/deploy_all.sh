#!/bin/bash

# Academic Token - Deploy All Contracts Script
# This script builds, stores, and instantiates all 5 contracts

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
CHAIN_ID="academictoken-1"
WALLET_NAME="deployer"
GAS_PRICES="0.025utoken"
GAS_ADJUSTMENT="1.3"

# Contract directories
CONTRACTS=(
    "prerequisites"
    "schedule" 
    "progress"
    "equivalence"
    "degree"
)

# Store contract code IDs
typeset -A CODE_IDS

echo -e "${BLUE}🚀 Academic Token - Deploying All Contracts${NC}"
echo -e "${BLUE}=============================================${NC}"
echo ""

# Function to build contract
build_contract() {
    local contract_name=$1
    echo -e "${YELLOW}📦 Building $contract_name contract...${NC}"
    
    cd contracts/$contract_name
    
    # Build the contract
    cargo build --release --target wasm32-unknown-unknown
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ $contract_name built successfully${NC}"
    else
        echo -e "${RED}❌ Failed to build $contract_name${NC}"
        exit 1
    fi
    
    cd ../..
}

# Function to store contract
store_contract() {
    local contract_name=$1
    echo -e "${YELLOW}📤 Storing $contract_name contract...${NC}"
    
    cd contracts/$contract_name
    
    # Store the contract
    local store_result=$(academictokend tx wasm store \
        target/wasm32-unknown-unknown/release/${contract_name}.wasm \
        --from $WALLET_NAME \
        --chain-id $CHAIN_ID \
        --gas-prices $GAS_PRICES \
        --gas auto \
        --gas-adjustment $GAS_ADJUSTMENT \
        --output json \
        --yes)
    
    # Extract code ID from transaction
    local tx_hash=$(echo $store_result | jq -r '.txhash')
    echo "Transaction hash: $tx_hash"
    
    # Wait for transaction to be included
    echo "Waiting for transaction to be included..."
    sleep 6
    
    # Get code ID from transaction result
    local code_id=$(academictokend query tx $tx_hash --output json | jq -r '.logs[0].events[] | select(.type=="store_code") | .attributes[] | select(.key=="code_id") | .value')
    
    if [ "$code_id" != "null" ] && [ "$code_id" != "" ]; then
        CODE_IDS[$contract_name]=$code_id
        echo -e "${GREEN}✅ $contract_name stored with Code ID: $code_id${NC}"
    else
        echo -e "${RED}❌ Failed to get code ID for $contract_name${NC}"
        exit 1
    fi
    
    cd ../..
}

# Function to instantiate contract
instantiate_contract() {
    local contract_name=$1
    local code_id=${CODE_IDS[$contract_name]}
    
    echo -e "${YELLOW}🎯 Instantiating $contract_name contract...${NC}"
    
    # Define instantiate messages for each contract
    local init_msg=""
    case $contract_name in
        "prerequisites")
            init_msg='{"owner": null}'
            ;;
        "schedule")
            init_msg='{"owner": null}'
            ;;
        "progress")
            init_msg='{"owner": null, "analytics_enabled": true, "update_frequency": "Daily", "analytics_depth": "Standard"}'
            ;;
        "equivalence")
            init_msg='{"owner": null, "similarity_threshold": 80, "auto_approval_threshold": 95}'
            ;;
        "degree")
            init_msg='{"owner": null}'
            ;;
    esac
    
    # Instantiate the contract
    local instantiate_result=$(academictokend tx wasm instantiate \
        $code_id \
        "$init_msg" \
        --from $WALLET_NAME \
        --label "${contract_name}-contract-v1" \
        --chain-id $CHAIN_ID \
        --gas-prices $GAS_PRICES \
        --gas auto \
        --gas-adjustment $GAS_ADJUSTMENT \
        --output json \
        --yes)
    
    # Extract contract address
    local tx_hash=$(echo $instantiate_result | jq -r '.txhash')
    echo "Transaction hash: $tx_hash"
    
    # Wait for transaction to be included
    echo "Waiting for transaction to be included..."
    sleep 6
    
    # Get contract address from transaction result
    local contract_address=$(academictokend query tx $tx_hash --output json | jq -r '.logs[0].events[] | select(.type=="instantiate") | .attributes[] | select(.key=="_contract_address") | .value')
    
    if [ "$contract_address" != "null" ] && [ "$contract_address" != "" ]; then
        echo -e "${GREEN}✅ $contract_name instantiated at: $contract_address${NC}"
        
        # Save contract address to file
        echo "$contract_name=$contract_address" >> deployed_contracts.txt
    else
        echo -e "${RED}❌ Failed to get contract address for $contract_name${NC}"
        exit 1
    fi
}

# Main deployment process
main() {
    echo -e "${BLUE}Starting deployment process...${NC}"
    echo ""
    
    # Clear previous deployment file
    > deployed_contracts.txt
    > code_ids.txt
    
    # Build all contracts
    echo -e "${BLUE}🔨 PHASE 1: Building all contracts${NC}"
    for contract in "${CONTRACTS[@]}"; do
        build_contract $contract
        echo ""
    done
    
    echo -e "${BLUE}📤 PHASE 2: Storing all contracts${NC}"
    for contract in "${CONTRACTS[@]}"; do
        store_contract $contract
        echo ""
    done
    
    echo -e "${BLUE}🎯 PHASE 3: Instantiating all contracts${NC}"
    for contract in "${CONTRACTS[@]}"; do
        instantiate_contract $contract
        echo ""
    done
    
    # Save code IDs to file
    echo -e "${BLUE}💾 Saving deployment information...${NC}"
    for contract in "${CONTRACTS[@]}"; do
        echo "$contract=${CODE_IDS[$contract]}" >> code_ids.txt
    done
    
    echo ""
    echo -e "${GREEN}🎉 ALL CONTRACTS DEPLOYED SUCCESSFULLY!${NC}"
    echo -e "${GREEN}======================================${NC}"
    echo ""
    echo -e "${BLUE}📋 Deployment Summary:${NC}"
    echo ""
    
    echo -e "${YELLOW}Code IDs:${NC}"
    cat code_ids.txt
    echo ""
    
    echo -e "${YELLOW}Contract Addresses:${NC}"
    cat deployed_contracts.txt
    echo ""
    
    echo -e "${BLUE}📁 Files created:${NC}"
    echo "  - code_ids.txt (contract code IDs)"
    echo "  - deployed_contracts.txt (contract addresses)"
    echo ""
    
    echo -e "${GREEN}✅ Deployment completed successfully!${NC}"
}

# Function to show help
show_help() {
    echo "Academic Token Deploy Script"
    echo ""
    echo "Usage: ./deploy_all.sh [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -h, --help     Show this help message"
    echo "  -w, --wallet   Specify wallet name (default: deployer)"
    echo "  -c, --chain    Specify chain ID (default: academictoken-1)"
    echo ""
    echo "Examples:"
    echo "  ./deploy_all.sh"
    echo "  ./deploy_all.sh --wallet alice --chain testnet-1"
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -w|--wallet)
            WALLET_NAME="$2"
            shift 2
            ;;
        -c|--chain)
            CHAIN_ID="$2"
            shift 2
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            show_help
            exit 1
            ;;
    esac
done

# Check if academictokend is available
if ! command -v academictokend &> /dev/null; then
    echo -e "${RED}❌ academictokend command not found${NC}"
    echo "Please ensure the Academic Token daemon is installed and in your PATH"
    exit 1
fi

# Check if jq is available (for JSON parsing)
if ! command -v jq &> /dev/null; then
    echo -e "${RED}❌ jq command not found${NC}"
    echo "Please install jq for JSON parsing: brew install jq (macOS) or apt install jq (Ubuntu)"
    exit 1
fi

# Verify wallet exists
echo -e "${BLUE}🔐 Verifying wallet '$WALLET_NAME'...${NC}"
if ! academictokend keys show $WALLET_NAME &> /dev/null; then
    echo -e "${RED}❌ Wallet '$WALLET_NAME' not found${NC}"
    echo "Available wallets:"
    academictokend keys list
    exit 1
fi

echo -e "${GREEN}✅ Wallet '$WALLET_NAME' found${NC}"
echo ""

# Start deployment
main
