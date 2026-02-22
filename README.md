# API Users (In-Memory)

API simples em Go para gerenciamento de usuários, criada para fins de estudo.  
Os dados são armazenados **em memória** utilizando um `map[string]User` (sem banco de dados).

A API utiliza:
- Go
- Chi Router
- UUID para geração de IDs

---

## 🚀 Executando o projeto

```bash
go run main.go