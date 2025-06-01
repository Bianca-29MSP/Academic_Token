# Instalar dependências do servidor REST

# 1. Inicializar módulo Go (se não existir)
go mod init academictoken-rest || echo "go.mod já existe"

# 2. Instalar dependências
go get github.com/gorilla/mux
go get github.com/gorilla/handlers

# 3. Executar servidor REST
go run cmd/rest-server/main.go
