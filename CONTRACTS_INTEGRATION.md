# Academic Token Contract Addresses

This file contains all deployed contract addresses and integration configuration for the Academic Token project.

## 📋 Deployed Contracts

| Contract | Code ID | Address | Description |
|----------|---------|---------|-------------|
| Equivalence | 1 | `academic14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9sva26tx` | Subject equivalence analysis |
| Progress | 2 | `academic1nc5tatafv6eyq7llkr2gv50ff9e22mnf70qgjlv737ktmt4eswrqqg0cfq` | Student progress tracking |
| Degree | 3 | `academic17p9rzwnnfxcjp32un9ug7yhhzgtkhvl9jfksztgw5uh69wac2pgscn3aqy` | Degree validation |
| Schedule | 4 | `academic1ghd753shjuwexxywmgs4xz7x2q732vcnkm6h2pyv9s6ah3hylvrqy9uyq9` | Schedule recommendations |
| Academic-NFT | 5 | `academic1eyfccmjm6732k7wp4p6gdjwhxjwsvje44j0hfx8nkgrm8fs7vqfsljrz85` | NFT minting and management |

## 👤 Admin Address

- **Alice**: `academic1urs658n5ddln24g23gj27wr9y7rnurxfw6rcsc`

## 🔗 Contract Integration Flow

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│ Student Module  │───▶│ Progress Contract│───▶│Academic-NFT     │
│                 │    │                 │    │Contract         │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                        ▲
┌─────────────────┐    ┌─────────────────┐              │
│ Degree Module   │───▶│ Degree Contract │──────────────┘
│                 │    │                 │
└─────────────────┘    └─────────────────┘

┌─────────────────┐    ┌─────────────────┐
│Schedule Module  │───▶│Schedule Contract│
│                 │    │                 │
└─────────────────┘    └─────────────────┘

┌─────────────────┐    ┌─────────────────┐
│Equivalence      │───▶│Equivalence      │
│Module           │    │Contract         │
└─────────────────┘    └─────────────────┘
```

## 🚀 Setup Scripts

1. **Configure Integration**: `./setup_contract_integration.sh`
2. **Test Integration**: `./test_contract_integration.sh`
3. **Load Config**: `source ./contracts_config.sh`

## 📝 Usage Examples

### Query Contract State
```bash
# Load addresses
source ./contracts_config.sh

# Query Academic-NFT config
academictokend query wasm contract-state smart $ACADEMIC_NFT_CONTRACT '{"get_config": {}}'

# Query Progress state  
academictokend query wasm contract-state smart $PROGRESS_CONTRACT '{"get_state": {}}'
```

### Update Contract Configuration
```bash
# Update Academic-NFT minter
academictokend tx wasm execute $ACADEMIC_NFT_CONTRACT \
  '{"update_config": {"minter": "'$ALICE_ADDRESS'"}}' \
  --from alice --chain-id academictoken -y

# Update Progress analytics
academictokend tx wasm execute $PROGRESS_CONTRACT \
  '{"update_config": {"analytics_enabled": true}}' \
  --from alice --chain-id academictoken -y
```

## 🔧 Module Integration Points

### Student Module Integration
- **Progress Contract**: Update student progress when subjects completed
- **Academic-NFT Contract**: Mint NFTs for completed subjects

### Degree Module Integration  
- **Degree Contract**: Validate graduation requirements
- **Academic-NFT Contract**: Mint degree NFTs

### Schedule Module Integration
- **Schedule Contract**: Generate recommendations
- **Progress Contract**: Get student progress data

### Equivalence Module Integration
- **Equivalence Contract**: Analyze subject equivalences

### AcademicNFT Module Integration
- **Academic-NFT Contract**: Direct NFT management

## 📊 Configuration Summary

| Setting | Value |
|---------|-------|
| Chain ID | `academictoken` |
| IPFS Gateway | `https://ipfs.io` |
| Max Subjects/Semester | `6` |
| Auto Approval Threshold | `85%` |
| Analytics Enabled | `true` |
| Update Frequency | `daily` |
| Analytics Depth | `standard` |

## ✅ Next Steps

1. Run integration setup: `chmod +x *.sh && ./setup_contract_integration.sh`
2. Test contracts: `./test_contract_integration.sh`
3. Update module configurations with contract addresses
4. Test end-to-end flows
5. Deploy frontend with contract integration
