#!/bin/bash
set -e

echo "Building contract..."
cargo wasm

echo "Optimizing contract..."
docker run --rm -v "$(pwd)":/code \
  --mount type=volume,source="$(basename "$(pwd)")_cache",target=/code/target \
  --mount type=volume,source=registry_cache,target=/usr/local/cargo/registry \
  cosmwasm/rust-optimizer:0.16.0

echo "Contract built successfully!"
echo "Wasm file: artifacts/prerequisites_contract.wasm"

# Deploy commands
echo ""
echo "To deploy, run:"
echo "academictokend tx wasm store artifacts/prerequisites_contract.wasm --from alice --chain-id academictoken --gas auto --gas-adjustment 1.5"