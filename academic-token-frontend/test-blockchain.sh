#!/bin/bash

# Script to test the blockchain connection

echo "🚀 Academic Token - Blockchain Connection Test"
echo "=============================================="
echo ""

# Check if node is installed
if ! command -v node &> /dev/null; then
    echo "❌ Node.js is not installed"
    exit 1
fi

# Run the test script
node scripts/test-connection.js
