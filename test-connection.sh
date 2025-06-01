#!/bin/bash

echo "🧪 TESTE DE CONEXÃO - ACADEMIC TOKEN"
echo "=================================="
echo ""

echo "1. 🔍 Testando Blockchain Cosmos (porta 1317)..."
curl -s http://localhost:1317/cosmos/base/tendermint/v1beta1/node_info | jq '.default_node_info.network' 2>/dev/null || echo "❌ Blockchain não respondeu"
echo ""

echo "2. 🔍 Testando REST Server (porta 1318)..."
curl -s http://localhost:1318/health || echo "❌ REST Server não está rodando"
echo ""

echo "3. 🔍 Testando IPFS (porta 5001)..."
curl -s -X POST http://localhost:5001/api/v0/version | jq '.Version' 2>/dev/null || echo "❌ IPFS não respondeu"
echo ""

echo "4. 🔍 Testando Frontend (porta 3001)..."
curl -s -o /dev/null -w "%{http_code}" http://localhost:3001 | grep -q "200" && echo "✅ Frontend OK" || echo "❌ Frontend não respondeu"
echo ""

echo "=================================="
echo "✅ Para iniciar tudo, use:"
echo "Terminal 1: go run cmd/rest-server/main.go"
echo "Terminal 2: cd academic-token-frontend && npm run dev" 
echo "Terminal 3: ipfs daemon"
echo "Terminal 4: ignite chain serve"
