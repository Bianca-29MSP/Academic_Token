#!/bin/bash

# Academic Token - Simple Deploy Script (macOS compatible)

set -e

# Configuration
WALLET_NAME=${1:-"deployer"}
CHAIN_ID="academictoken-1" 
GAS_PRICES="0.025utoken"

echo "🚀 Academic Token - Simple Deploy"
echo "================================="
echo ""

# Contracts to deploy
contracts=("prerequisites" "schedule" "progress" "equivalence" "degree")

# Create results file
echo "# Academic Token Deployment Results" > deployment_results.txt
echo "# Generated: $(date)" >> deployment_results.txt
echo "" >> deployment_results.txt

for contract in "${contracts[@]}"; do
    echo "📦 Deploying: $contract"
    echo "----------------------------"
    
    # Navigate to contract
    cd contracts/$contract
    
    # Build
    echo "Building..."
    cargo build --release --target wasm32-unknown-unknown
    
    # Store
    echo "Storing..."
    store_tx=$(academictokend tx wasm store \
        target/wasm32-unknown-unknown/release/${contract}.wasm \
        --from $WALLET_NAME \
        --chain-id $CHAIN_ID \
        --gas-prices $GAS_PRICES \
        --gas auto \
        --gas-adjustment 1.3 \
        --output json \
        --yes)
    
    # Get transaction hash
    tx_hash=$(echo $store_tx | jq -r '.txhash')
    echo "Store TX: $tx_hash"
    
    # Wait for confirmation
    echo "Waiting for confirmation..."
    sleep 6
    
    # Get code ID
    code_id=$(academictokend query tx $tx_hash --output json | jq -r '.logs[0].events[] | select(.type=="store_code") | .attributes[] | select(.key=="code_id") | .value')
    echo "Code ID: $code_id"
    
    # Instantiate message
    case $contract in
        "prerequisites"|"schedule"|"degree")
            init_msg='{"owner": null}'
            ;;
        "progress")
            init_msg='{"owner": null, "analytics_enabled": true, "update_frequency": "Daily", "analytics_depth": "Standard"}'
            ;;
        "equivalence")
            init_msg='{"owner": null, "similarity_threshold": 80, "auto_approval_threshold": 95}'
            ;;
    esac
    
    # Instantiate
    echo "Instantiating..."
    inst_tx=$(academictokend tx wasm instantiate \
        $code_id \
        "$init_msg" \
        --from $WALLET_NAME \
        --label "${contract}-v1" \
        --chain-id $CHAIN_ID \
        --gas-prices $GAS_PRICES \
        --gas auto \
        --gas-adjustment 1.3 \
        --output json \
        --yes)
    
    # Get transaction hash
    tx_hash=$(echo $inst_tx | jq -r '.txhash')
    echo "Instantiate TX: $tx_hash"
    
    # Wait for confirmation
    sleep 6
    
    # Get contract address
    contract_addr=$(academictokend query tx $tx_hash --output json | jq -r '.logs[0].events[] | select(.type=="instantiate") | .attributes[] | select(.key=="_contract_address") | .value')
    
    echo "✅ SUCCESS!"
    echo "Contract: $contract"
    echo "Code ID: $code_id"
    echo "Address: $contract_addr"
    echo ""
    
    # Save to results
    echo "## $contract" >> ../../deployment_results.txt
    echo "- Code ID: $code_id" >> ../../deployment_results.txt
    echo "- Address: $contract_addr" >> ../../deployment_results.txt
    echo "- Store TX: $tx_hash" >> ../../deployment_results.txt
    echo "" >> ../../deployment_results.txt
    
    # Go back to root
    cd ../..
done

echo "🎉 ALL CONTRACTS DEPLOYED!"
echo "========================="
echo ""
echo "📋 Results saved to: deployment_results.txt"
echo ""
cat deployment_results.txt
