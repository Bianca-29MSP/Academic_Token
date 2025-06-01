# 🚨 TROUBLESHOOTING: Instituições e Subjects não aparecem

## 🔍 Passo-a-Passo para Diagnóstico

### 1. **PRIMEIRO: Verificar se o servidor REST está rodando**

```bash
# Terminal 1: Iniciar o servidor REST
cd /Users/biancamsp/Desktop/Academic_Token/academictoken/academictoken
go run cmd/rest-server/main.go
```

Você deve ver:
```
🚀 Academic Token REST Server
🌍 API Server: http://localhost:1318
📡 Academic API: http://localhost:1318/academic
💡 Health: http://localhost:1318/health

✅ Pronto para conectar com o frontend!
```

### 2. **SEGUNDO: Testar a API manualmente**

```bash
# Terminal 2: Testar endpoints
cd /Users/biancamsp/Desktop/Academic_Token/academictoken/academictoken/academic-token-frontend

# Executar script de teste
chmod +x test-api-simple.sh
./test-api-simple.sh
```

**Resultados Esperados:**
- ✅ Health check: OK
- ✅ Node info: OK  
- ✅ Institutions loaded successfully (deve mostrar UFJF e USP)
- ✅ Subjects loaded successfully (deve mostrar Cálculo I, Programação 1, etc.)

### 3. **TERCEIRO: Iniciar o frontend**

```bash
# Terminal 3: Frontend
cd /Users/biancamsp/Desktop/Academic_Token/academictoken/academictoken/academic-token-frontend
npm run dev
```

### 4. **QUARTO: Acessar a página de debug**

1. Acesse: `http://localhost:3000/equivalences`
2. Clique na aba **"🔧 Debug"**
3. Verifique:
   - **Connected:** Deve ser ✅ Yes
   - **Institutions:** Deve mostrar "2 loaded"
   - **Subjects:** Deve mostrar "4 loaded"

---

## 🐛 Problemas Comuns e Soluções

### ❌ Problema: "Connection failed - blockchain node may be offline"

**Solução:**
```bash
# Verificar se o servidor REST está rodando
curl http://localhost:1318/health

# Se retornar "OK", o servidor está funcionando
# Se não retornar nada, iniciar o servidor:
cd /Users/biancamsp/Desktop/Academic_Token/academictoken/academictoken
go run cmd/rest-server/main.go
```

### ❌ Problema: "0 institutions loaded" mas servidor está rodando

**Solução:**
1. Abrir o Developer Tools do browser (F12)
2. Ir na aba "Console"
3. Procurar por erros de CORS ou API
4. Verificar se há logs como:
   ```
   🔗 API Request: GET /academic/institution/list
   ✅ API Response: 200 /academic/institution/list
   ```

### ❌ Problema: CORS Error

**Verificar se o servidor REST tem CORS configurado:**
O arquivo `cmd/rest-server/main.go` deve ter:
```go
originsOk := handlers.AllowedOrigins([]string{"http://localhost:3000", "http://localhost:3001"})
```

### ❌ Problema: "Cannot read properties of undefined"

**Solução:** Verificar se os dados estão no formato correto:
```bash
# Testar formato dos dados
curl http://localhost:1318/academic/institution/list | head

# Deve retornar algo como:
# [{"id":"inst_1","name":"Universidade Federal de Juiz de Fora",...}]
```

---

## 🔧 Debug Avançado

### Verificar logs detalhados no browser:

1. Abrir Developer Tools (F12)
2. Ir na aba "Console"
3. Procurar por logs:
   ```
   🚀 Initializing blockchain hook...
   🔄 Refreshing blockchain connection...
   📡 Connection Status: {connected: true, ...}
   ✅ Connection established, loading data...
   🔄 Loading blockchain data...
   📊 Fetching institutions...
   🏛️ Institutions loaded: [{...}]
   ```

### Se não vir esses logs:

1. **Verificar se o hook está sendo usado corretamente**
2. **Verificar se há erros de importação**
3. **Verificar se o componente está renderizando**

### Teste manual no browser console:

```javascript
// Testar API diretamente no browser
fetch('http://localhost:1318/academic/institution/list')
  .then(r => r.json())
  .then(data => console.log('Institutions:', data))

fetch('http://localhost:1318/academic/subject/list')
  .then(r => r.json())
  .then(data => console.log('Subjects:', data))
```

---

## ✅ Lista de Verificação

- [ ] Servidor REST rodando na porta 1318
- [ ] Script de teste da API passou
- [ ] Frontend rodando na porta 3000
- [ ] Não há erros no console do browser
- [ ] Aba Debug mostra "Connected: ✅ Yes"
- [ ] Aba Debug mostra "Institutions: 2 loaded"
- [ ] Aba Debug mostra "Subjects: 4 loaded"
- [ ] Dropdowns de instituição mostram UFJF e USP
- [ ] Dropdowns de disciplinas mostram as matérias

---

## 🆘 Se ainda não funcionar

1. **Pare tudo:**
   ```bash
   # Parar frontend (Ctrl+C)
   # Parar servidor REST (Ctrl+C)
   ```

2. **Limpe o cache:**
   ```bash
   # Limpar cache do Next.js
   rm -rf .next
   npm run dev
   ```

3. **Reinicie na ordem:**
   ```bash
   # Terminal 1: Backend primeiro
   cd /Users/biancamsp/Desktop/Academic_Token/academictoken/academictoken
   go run cmd/rest-server/main.go
   
   # Terminal 2: Teste API
   cd academic-token-frontend
   ./test-api-simple.sh
   
   # Terminal 3: Frontend por último
   npm run dev
   ```

4. **Verifique a aba Debug** em `http://localhost:3000/equivalences`

---

**Se seguir todos esses passos, o sistema deve funcionar perfeitamente! 🎯**
