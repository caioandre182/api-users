# API Users (PostgreSQL + Docker)

API simples em Go para gerenciamento de usuários, criada para fins de estudo.
Os dados agora são persistidos em **PostgreSQL**, utilizando Docker para containerização da aplicação e do banco.

A API utiliza:

* Go
* Chi Router
* UUID para geração de IDs
* database/sql
* pgx (driver PostgreSQL)
* Docker
* Docker Compose
* pgAdmin

---

## 🚀 Executando o projeto

### 1️⃣ Subir os containers

```bash
docker compose up -d --build
```

Isso irá subir:

* API → [http://localhost:8080](http://localhost:8080)
* PostgreSQL
* pgAdmin → [http://localhost:5050](http://localhost:5050)

---

## 🗄 Configuração do Banco

O banco é configurado via `docker-compose.yml` com:

* POSTGRES_USER=app
* POSTGRES_PASSWORD=app
* POSTGRES_DB=appdb

A aplicação utiliza a variável:

```
DATABASE_URL=postgres://app:app@postgres:5432/appdb?sslmode=disable
```

---

## 🗃 Estrutura da Tabela

```sql
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    biography TEXT
);
```

---

## 📌 Endpoints

### Criar usuário

```
POST /user
```

Body:

```json
{
  "first_name": "Caio",
  "last_name": "Macedo",
  "biography": "Software Developer"
}
```

---

### Listar usuários

```
GET /users
```

---

### Buscar usuário por ID

```
GET /user/{id}
```

---

### Atualizar usuário

```
PUT /user/{id}
```

---

### Deletar usuário

```
DELETE /user/{id}
```

---

## 📦 Estrutura do Projeto

```
api-users/
│
├── api/                 # Camada HTTP (handlers e rotas)
├── domain/              # Entidade User
├── store/               # Interface de persistência
│   └── postgres/        # Implementação PostgreSQL
├── main.go              # Bootstrap da aplicação
├── Dockerfile
└── docker-compose.yml
```

---

## 🧠 Conceitos aplicados

* Injeção de dependência
* Interface para desacoplamento
* Context propagation
* ExecContext / QueryContext
* Tratamento de erro customizado (ErrNotFound)
* Containerização com multi-stage build
* Separação por camadas

---

## 🔄 Evoluções futuras

* Migrations (Tern ou golang-migrate)
* Logs estruturados
* Testes unitários com mock do store
* CI/CD
* Autenticação JWT
* Deploy em cloud
