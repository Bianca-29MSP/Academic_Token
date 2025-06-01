#!/bin/bash

# Optimized WASM Build Script

set -e

CONTRACT_NAME=$1

if [ -z "$CONTRACT_NAME" ]; then
    echo "Usage: ./build_optimized.sh <contract_name>"
    exit 1
fi

echo "🔨 Building optimized WASM for: $CONTRACT_NAME"

cd contracts/$CONTRACT_NAME

# Build with optimization flags
echo "1. Building with Rust optimization..."
RUSTFLAGS='-C link-arg=-s -C target-feature=+bulk-memory,+sign-ext' \
cargo build --release --target wasm32-unknown-unknown

# Check if wasm-opt is available
if command -v wasm-opt &> /dev/null; then
    echo "2. Optimizing with wasm-opt..."
    wasm-opt -Oz --enable-reference-types \
        target/wasm32-unknown-unknown/release/${CONTRACT_NAME}.wasm \
        -o target/wasm32-unknown-unknown/release/${CONTRACT_NAME}_optimized.wasm
    
    echo "3. Replacing original with optimized..."
    mv target/wasm32-unknown-unknown/release/${CONTRACT_NAME}_optimized.wasm \
       target/wasm32-unknown-unknown/release/${CONTRACT_NAME}.wasm
else
    echo "2. wasm-opt not found, skipping optimization"
    echo "   Install with: brew install binaryen"
fi

echo "✅ Build complete!"
ls -la target/wasm32-unknown-unknown/release/${CONTRACT_NAME}.wasm

cd ../..
