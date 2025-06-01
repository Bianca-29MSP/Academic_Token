#!/bin/bash

# Academic Token - Auto Deploy (finds wallet with balance)

set -e

echo "🚀 Academic Token - Auto Deploy"
echo "==============================="
echo ""

# Function to check balance
check_balance() {
    local address=$1
    local balance=$(academictokend query bank balances $address --output json 2>/dev/null | jq -r '.balances[0].amount // "0"')
    echo $balance
}

# Find wallet with balance
echo "🔍 Looking for wallet with balance..."

# First check if we have any keys
keys=$(academictokend keys list --output json 2>/dev/null)
if [ "$keys" = "[]" ] || [ "$keys" = "" ]; then
    echo "❌ No wallets found!"
    echo "Create a wallet first:"
    echo "academictokend keys add deployer"
    exit 1
fi

# Check each wallet
found_wallet=""
found_address=""
found_balance=0

echo "$keys" | jq -r '.[].name' | while read wallet_name; do
    if [ "$wallet_name" != "" ]; then
        address=$(academictokend keys show $wallet_name --address 2>/dev/null)
        balance=$(check_balance $address)
        
        echo "🔑 Wallet: $wallet_name"
        echo "   Address: $address" 
        echo "   Balance: $balance tokens"
        
        if [ "$balance" != "0" ] && [ "$balance" != "" ]; then
            echo "✅ Found wallet with balance!"
            echo "USING_WALLET=$wallet_name" > .wallet_config
            echo "USING_ADDRESS=$address" >> .wallet_config
            echo "USING_BALANCE=$balance" >> .wallet_config
            break
        fi
    fi
done

# Source the config if found
if [ -f ".wallet_config" ]; then
    source .wallet_config
    echo ""
    echo "💰 Using wallet: $USING_WALLET ($USING_ADDRESS)"
    echo "💰 Balance: $USING_BALANCE tokens"
    echo ""
else
    echo "❌ No wallet with balance found!"
    echo ""
    echo "💡 Solutions:"
    echo "1. Fund an existing wallet"
    echo "2. Transfer tokens from another source"
    echo "3. Use a faucet if on testnet"
    exit 1
fi

# Proceed with deployment
contracts=("prerequisites" "schedule" "progress" "equivalence" "degree")

echo "# Academic Token Deployment Results" > deployment_results.txt
echo "# Generated: $(date)" >> deployment_results.txt
echo "# Wallet: $USING_WALLET" >> deployment_results.txt
echo "# Address: $USING_ADDRESS" >> deployment_results.txt
echo "" >> deployment_results.txt

for contract in "${contracts[@]}"; do
    echo "📦 Deploying: $contract"
    echo "----------------------------"
    
    cd contracts/$contract
    
    # Build with specific WASM flags
    echo "Building..."
    RUSTFLAGS='-C link-arg=-s' cargo build --release --target wasm32-unknown-unknown
    
    echo "Storing..."
    store_tx=$(academictokend tx wasm store \
        target/wasm32-unknown-unknown/release/${contract}.wasm \
        --from $USING_WALLET \
        --chain-id academictoken-1 \
        --gas-prices 0.025utoken \
        --gas auto \
        --gas-adjustment 1.3 \
        --output json \
        --yes)
    
    tx_hash=$(echo $store_tx | jq -r '.txhash')
    echo "Store TX: $tx_hash"
    
    sleep 6
    
    code_id=$(academictokend query tx $tx_hash --output json | jq -r '.logs[0].events[] | select(.type=="store_code") | .attributes[] | select(.key=="code_id") | .value')
    echo "Code ID: $code_id"
    
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
    
    echo "Instantiating..."
    inst_tx=$(academictokend tx wasm instantiate \
        $code_id \
        "$init_msg" \
        --from $USING_WALLET \
        --label "${contract}-v1" \
        --chain-id academictoken-1 \
        --gas-prices 0.025utoken \
        --gas auto \
        --gas-adjustment 1.3 \
        --output json \
        --yes)
    
    tx_hash=$(echo $inst_tx | jq -r '.txhash')
    echo "Instantiate TX: $tx_hash"
    
    sleep 6
    
    contract_addr=$(academictokend query tx $tx_hash --output json | jq -r '.logs[0].events[] | select(.type=="instantiate") | .attributes[] | select(.key=="_contract_address") | .value')
    
    echo "✅ SUCCESS!"
    echo "Contract: $contract"
    echo "Code ID: $code_id"
    echo "Address: $contract_addr"
    echo ""
    
    echo "## $contract" >> ../../deployment_results.txt
    echo "- Code ID: $code_id" >> ../../deployment_results.txt
    echo "- Address: $contract_addr" >> ../../deployment_results.txt
    echo "" >> ../../deployment_results.txt
    
    cd ../..
done

echo "🎉 ALL CONTRACTS DEPLOYED!"
cat deployment_results.txt

# Cleanup
rm -f .wallet_config
