# Academic Token Frontend - Blockchain Integration

## 🚀 Quick Start

### 1. Start Backend REST Server
```bash
# Navigate to backend directory
cd /Users/biancamsp/Desktop/Academic_Token/academictoken/academictoken

# Start the REST server
go run cmd/rest-server/main.go
```

The REST server will be available at: `http://localhost:1318`

### 2. Start Frontend
```bash
# Navigate to frontend directory
cd /Users/biancamsp/Desktop/Academic_Token/academictoken/academictoken/academic-token-frontend

# Install dependencies (if needed)
npm install

# Start development server
npm run dev
```

Frontend will be available at: `http://localhost:3000`

## 🔗 Blockchain Integration Features

### ✅ What's Working Now

1. **Real-time Connection Status**
   - ✅ Connection monitoring to blockchain node
   - ✅ Network info display (academictoken network)
   - ✅ Error handling and user feedback

2. **Institution Management**
   - ✅ Load institutions from blockchain API
   - ✅ Display in dropdown selectors
   - ✅ Real-time data sync

3. **Subject Management**
   - ✅ Load subjects from blockchain API
   - ✅ Filter by institution
   - ✅ Display with credits and metadata

4. **Smart Contract Simulation**
   - ✅ Equivalence analysis simulation
   - ✅ Similarity percentage calculation
   - ✅ Approval/rejection recommendations

### 🏗️ Technical Architecture

```
Frontend (Next.js)
├── /app/types/blockchain.ts          # TypeScript interfaces
├── /app/lib/api.ts                   # API service layer
├── /app/hooks/useBlockchain.ts       # State management hook
└── /app/equivalences/page.tsx        # Updated equivalences page

Backend (Go REST Server)
├── cmd/rest-server/main.go           # REST API endpoints
└── Endpoints:
    ├── GET /health                   # Health check
    ├── GET /academic/institution/list # Institutions
    ├── GET /academic/subject/list    # Subjects
    └── GET /cosmos/base/tendermint/v1beta1/node_info # Node info
```

### 📊 Available Data

**Institutions:**
- UFJF (Universidade Federal de Juiz de Fora)
- USP (Universidade de São Paulo)

**Subjects (UFJF):**
- MAT101 - Cálculo I (4 créditos)
- CI1001 - Programação 1 (4 créditos)
- CI1002 - Programação 2 (4 créditos)
- MAT201 - Cálculo II (4 créditos)

### 🔄 How Equivalence Analysis Works

1. **User selects source and target institutions**
2. **Dropdowns populate with real blockchain data**
3. **User selects subjects from each institution**
4. **Smart contract simulation runs automatically**
5. **System provides similarity % and recommendation**
6. **User can submit equivalence request**

### 🐛 Troubleshooting

**Problem: "Blockchain Connection Error"**
```bash
# Solution: Make sure REST server is running
cd /Users/biancamsp/Desktop/Academic_Token/academictoken/academictoken
go run cmd/rest-server/main.go
```

**Problem: "No institutions loading"**
```bash
# Check if server is responding
curl http://localhost:1318/health
curl http://localhost:1318/academic/institution/list
```

**Problem: TypeScript errors**
```bash
# Restart the dev server
npm run dev
```

### 🎯 Next Steps for Full Blockchain Integration

1. **Real Smart Contracts**
   - Replace simulation with actual CosmWasm contract calls
   - Implement IPFS content retrieval for similarity analysis

2. **Transaction Handling**
   - Add wallet connection (Keplr/Cosmostation)
   - Implement transaction signing for equivalence requests

3. **Real-time Updates**
   - WebSocket connection for live blockchain events
   - Automatic UI updates when blockchain state changes

4. **Enhanced Error Handling**
   - Better error messages for failed transactions
   - Retry mechanisms for network issues

### 📝 Development Notes

- All API calls are logged to browser console
- Connection status updates automatically
- Form validation prevents invalid submissions
- Loading states provide user feedback
- Error boundaries handle API failures gracefully

### 🔧 Environment Variables

```env
NEXT_PUBLIC_API_URL=http://localhost:1318
NEXT_PUBLIC_BLOCKCHAIN_NETWORK=academictoken
NEXT_PUBLIC_BLOCKCHAIN_VERSION=v1.0.0
```

---

**Status:** ✅ Basic blockchain integration complete
**Next:** Implement real smart contract calls and wallet integration
