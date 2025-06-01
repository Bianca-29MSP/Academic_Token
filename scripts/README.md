# Academic Token - Deployment Scripts

Scripts para fazer deploy dos 5 contratos do sistema Academic Token.

## 📋 Contratos

1. **prerequisites** - Verificação de pré-requisitos
2. **schedule** - Recomendações de horários
3. **progress** - Acompanhamento de progresso
4. **equivalence** - Análise de equivalências
5. **degree** - Emissão de diplomas

## 🚀 Como usar

### Deploy de todos os contratos:
```bash
cd /Users/biancamsp/Desktop/Academic_Token/academictoken/academictoken

# Tornar executável
chmod +x scripts/deploy_all.sh

# Executar
./scripts/deploy_all.sh
```

### Deploy de um contrato específico:
```bash
# Tornar executável
chmod +x scripts/deploy_single.sh

# Executar
./scripts/deploy_single.sh progress
./scripts/deploy_single.sh equivalence
```

### Com wallet específica:
```bash
./scripts/deploy_all.sh --wallet alice --chain testnet-1
./scripts/deploy_single.sh progress alice
```

## 📁 Arquivos gerados

Após o deploy, serão criados:
- `code_ids.txt` - IDs dos códigos armazenados
- `deployed_contracts.txt` - Endereços dos contratos
- `deployed_contracts.log` - Log completo

## ⚡ O que cada script faz

### deploy_all.sh
1. 📦 **Build** todos os 5 contratos
2. 📤 **Store** todos no blockchain  
3. 🎯 **Instantiate** todos com configurações padrão
4. 💾 **Salva** todos os IDs e endereços

### deploy_single.sh  
1. 📦 **Build** um contrato específico
2. 📤 **Store** no blockchain
3. 🎯 **Instantiate** com configuração padrão
4. 💾 **Salva** ID e endereço

## 🔧 Pré-requisitos

- `academictokend` instalado e configurado
- `jq` para parsing JSON: `brew install jq`
- Wallet configurada com fundos
- Node sincronizado

## 📝 Configurações de Instantiate

### Progress Contract:
```json
{
  "owner": null,
  "analytics_enabled": true, 
  "update_frequency": "Daily",
  "analytics_depth": "Standard"
}
```

### Equivalence Contract:
```json
{
  "owner": null,
  "similarity_threshold": 80,
  "auto_approval_threshold": 95  
}
```

### Outros Contratos:
```json
{
  "owner": null
}
```

## 🎯 Exemplo de uso completo

```bash
# 1. Navegar para o projeto
cd /Users/biancamsp/Desktop/Academic_Token/academictoken/academictoken

# 2. Fazer deploy de todos
chmod +x scripts/deploy_all.sh
./scripts/deploy_all.sh

# 3. Verificar resultados
cat code_ids.txt
cat deployed_contracts.txt
```

## 🔍 Verificar deployment

```bash
# Listar códigos armazenados
academictokend query wasm list-code

# Verificar contrato específico
academictokend query wasm contract <CONTRACT_ADDRESS>

# Testar query
academictokend query wasm contract-state smart <CONTRACT_ADDRESS> '{"get_state":{}}'
```
