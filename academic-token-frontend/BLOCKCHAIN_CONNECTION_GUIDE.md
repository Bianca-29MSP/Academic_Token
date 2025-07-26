# Academic Token - Setup and Connection Guide

## 🚀 Quick Start

### Prerequisites
- Node.js (v18 or higher)
- Go (v1.21 or higher)
- Ignite CLI (install with `curl https://get.ignite.com/ | bash`)

### 1. Start the Blockchain

First, you need to start the Academic Token blockchain:

```bash
# Navigate to the blockchain directory
cd /Users/biancamsp/Desktop/Academic_Token/academictoken/academictoken

# Start the blockchain with Ignite
ignite chain serve
```

This will:
- Build the blockchain
- Initialize the chain
- Start the node
- Expose the API on port 1318

### 2. Test the Connection

Once the blockchain is running, test the connection:

```bash
# In the frontend directory
cd /Users/biancamsp/Desktop/Academic_Token/academictoken/academictoken/academic-token-frontend

# Test all endpoints
npm run test:connection
```

You should see something like:
```
✅ Node Info: /cosmos/base/tendermint/v1beta1/node_info - OK
✅ Institution List: /academictoken/institution/institution - OK
...
```

### 3. Start the Frontend

With the blockchain running:

```bash
# In the frontend directory
npm run dev
```

Then open [http://localhost:3000](http://localhost:3000) in your browser.

## 📝 Common Issues and Solutions

### Issue: "Failed to query institutions/courses/subjects"

**Solution**: The blockchain might not be running or the endpoints are not accessible.

1. Check if the blockchain is running:
   ```bash
   npm run test:connection
   ```

2. If no endpoints are responding, start the blockchain:
   ```bash
   cd ../.. # Go to blockchain root
   ignite chain serve
   ```

### Issue: "Error: connect ECONNREFUSED"

**Solution**: The API server is not running on port 1318.

1. Make sure you're running `ignite chain serve` in the blockchain directory
2. Check the terminal output for any errors
3. Verify the API is exposed on port 1318

### Issue: Empty lists (no institutions/courses/subjects)

**Solution**: This is normal for a fresh blockchain. You need to create data.

1. The blockchain starts empty
2. Use the admin interface to create institutions, courses, and subjects
3. Or use the CLI to populate test data

## 🏗️ Architecture Overview

```
Academic Token Project
├── blockchain (Cosmos SDK)
│   ├── x/institution   - Institution module
│   ├── x/course       - Course module  
│   ├── x/subject      - Subject module
│   ├── x/student      - Student module
│   ├── x/academicnft  - NFT module
│   └── ...
└── academic-token-frontend (Next.js)
    ├── app/lib/api.ts - API client
    ├── app/hooks/     - React hooks
    └── app/components/ - UI components
```

## 🔧 API Endpoints

The frontend connects to these blockchain endpoints:

- **Institutions**: `/academictoken/institution/institution`
- **Courses**: `/academictoken/course/course`
- **Subjects**: `/academictoken/subject/subjects`
- **Students**: `/academictoken/student/students`
- **NFTs**: `/academictoken/academicnft/student/{studentId}/tokens`

## 🎯 Next Steps

1. **Create Test Data**: Use the admin interface to create institutions and courses
2. **Register Students**: Create student accounts and enroll them in courses
3. **Mint NFTs**: Complete subjects to receive academic NFTs
4. **Check Eligibility**: Monitor degree progress and eligibility

## 📚 Additional Resources

- [Cosmos SDK Documentation](https://docs.cosmos.network/)
- [Ignite CLI Documentation](https://docs.ignite.com/)
- [Next.js Documentation](https://nextjs.org/docs)
