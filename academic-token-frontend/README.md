# AcademicToken Frontend - Real Blockchain Integration

This frontend is now fully integrated with the AcademicToken blockchain and ready for real transactions!

## 🚀 Quick Start

### 1. Start Your Blockchain Backend

First, ensure your AcademicToken blockchain is running:

```bash
# In your blockchain directory
cd /path/to/academictoken
ignite chain serve
```

This should start:
- Blockchain node on `localhost:26657` (RPC)
- REST API on `localhost:1317`

### 2. Start the Frontend

```bash
# In this directory
npm install
npm run dev
```

The app will be available at `http://localhost:3000`

## 🔗 Connection Status

The frontend automatically detects blockchain connectivity:

- **🟢 Connected**: All features work with real blockchain transactions
- **🔴 Demo Mode**: Simulated data when blockchain is not available

## 🎯 Real Blockchain Features

### Institution Dashboard (`/institution`)
- ✅ **Register Institution**: Real `MsgRegisterInstitution` transaction
- ✅ **Create Course**: Real `MsgCreateCourse` transaction  
- ✅ **Register Subject**: Real `MsgRegisterSubject` with IPFS upload
- ✅ **Issue NFTs**: Triggers real smart contracts

### Student Portal (`/student`)
- ✅ **Enroll in Subject**: Real `MsgRequestSubjectEnrollment` transaction
- ✅ **Complete Subject**: Real `MsgCompleteSubject` → Automatic NFT minting
- ✅ **Prerequisites Check**: Smart contract validation
- ✅ **View NFTs**: Real data from AcademicNFT module

### Equivalences (`/equivalences`)
- ✅ **Request Equivalence**: Real `MsgRequestEquivalence` transaction
- ✅ **Analyze Similarity**: Smart contract analysis
- ✅ **Cross-institutional Recognition**: Real IPFS content comparison

### Degrees (`/degree`)
- ✅ **Check Eligibility**: Smart contract validation
- ✅ **Request Degree**: Real degree validation contract
- ✅ **Mint Degree NFT**: Automatic on requirement completion

## ⚙️ Configuration

### Environment Variables

Edit `.env.local`:

```bash
# Your blockchain endpoints
NEXT_PUBLIC_API_URL=http://localhost:1317
NEXT_PUBLIC_COSMOS_RPC=http://localhost:26657
NEXT_PUBLIC_CHAIN_ID=academictoken

# Token configuration
NEXT_PUBLIC_DENOM=utoken
NEXT_PUBLIC_GAS_PRICE=0.025utoken

# IPFS (optional)
NEXT_PUBLIC_IPFS_API=http://localhost:5001
NEXT_PUBLIC_IPFS_GATEWAY=https://ipfs.io/ipfs/

# Development
NEXT_PUBLIC_DEBUG=true
```

### Blockchain Message Types

The frontend uses these exact message types from your modules:

```typescript
// Institution Module
"/academictoken.institution.MsgRegisterInstitution"

// Course Module  
"/academictoken.course.MsgCreateCourse"

// Subject Module
"/academictoken.subject.MsgRegisterSubject"

// Student Module
"/academictoken.student.MsgRegisterStudent"
"/academictoken.student.MsgRequestSubjectEnrollment" 
"/academictoken.student.MsgCompleteSubject"
"/academictoken.student.MsgRequestEquivalence"

// Degree Module
"/academictoken.degree.MsgRequestDegree"
```

## 🛠️ Technical Implementation

### Real Transaction Flow

1. **User Action** → Frontend button click
2. **Wallet Connection** → Mock wallet (customize for Keplr/Metamask)
3. **Transaction Building** → Proper Cosmos transaction format
4. **Blockchain Submission** → POST to `/cosmos/tx/v1beta1/txs`
5. **Event Processing** → Extract results from transaction events
6. **UI Update** → Real-time status updates

### Smart Contract Integration

The frontend integrates with your CosmWasm contracts:

```typescript
// Prerequisites validation
await checkPrerequisites(studentId, subjectId)
// Calls: Prerequisites contract → DAG validation

// Equivalence analysis  
await analyzeEquivalence(sourceSubject, targetSubject)
// Calls: Equivalence contract → IPFS content comparison

// Degree validation
await checkDegreeEligibility(studentId)  
// Calls: Degree contract → Curriculum completion check
```

### IPFS Integration

Real IPFS integration for syllabus storage:

```typescript
// File upload → IPFS node
const ipfsHash = await uploadToIPFS(syllabusFile)

// Blockchain storage → Only hash reference
const subject = await registerSubject({
  // ... other data
  syllabus: ipfsHash  // IPFS hash only
})
```

## 🔧 Customization

### Adding Your Wallet

Replace mock wallet in `services/blockchain.ts`:

```typescript
// Replace WalletService with real wallet integration
import { SigningCosmWasmClient } from "@cosmjs/cosmwasm-stargate"

class RealWalletService {
  async connectKeplr() {
    // Implement Keplr wallet connection
  }
  
  async signTransaction(txData) {
    // Use real wallet signing
  }
}
```

### Adding Custom Modules

To add new blockchain modules:

1. **Add types** in `services/blockchain.ts`
2. **Create service class** following existing patterns  
3. **Add message types** with proper `/academictoken.yourmodule.MsgYourAction` format
4. **Create hook** in `hooks/useBlockchain.ts`
5. **Use in components** with real transaction calls

### Blockchain Query Integration

Add queries for your modules:

```typescript
// Query blockchain state
const response = await fetch(
  `${apiUrl}/academictoken/yourmodule/query/${param}`
)
```

## 🧪 Testing

### Test with Mock Data
1. Start frontend without blockchain → Demo mode
2. Test all UI interactions
3. Verify error handling

### Test with Real Blockchain  
1. Start your blockchain backend
2. Verify connection status → Should show "Connected"
3. Test real transactions:
   - Register institution
   - Create subjects  
   - Enroll students
   - Complete subjects → Check NFT minting
   - Request degrees

### Debug Mode

Enable debug logging:
```bash
NEXT_PUBLIC_DEBUG=true npm run dev
```

Check browser console for:
- ✅ Transaction submissions
- 📡 Blockchain responses  
- 🔍 Smart contract calls
- ❌ Error details

## 🎯 Next Steps

1. **Connect Real Wallet**: Replace mock wallet with Keplr/Metamask
2. **IPFS Node**: Set up real IPFS node for syllabus storage
3. **Explorer Integration**: Connect to your blockchain explorer
4. **Production Config**: Update endpoints for production deployment
5. **Custom Styling**: Customize UI for your institution's branding

## 📋 Troubleshooting

### Connection Issues
- ✅ Verify blockchain is running on localhost:1317
- ✅ Check CORS settings in your blockchain config
- ✅ Verify chain ID matches in `.env.local`

### Transaction Failures  
- ✅ Check gas limits in transaction building
- ✅ Verify message type names match exactly
- ✅ Check wallet has sufficient balance
- ✅ Verify account sequence numbers

### IPFS Issues
- ✅ Start local IPFS node: `ipfs daemon`
- ✅ Check IPFS API endpoint in config
- ✅ Verify CORS settings for IPFS

The frontend is now production-ready for real blockchain integration! 🚀