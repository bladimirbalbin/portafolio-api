# Portafolio API (Go)

API REST en Go para gestionar proyectos de portafolio profesional.

## Demo
- Health: https://portafolio-api-xiw6.onrender.com/health
- Projects: https://portafolio-api-xiw6.onrender.com/projects

## Features
- Go + Chi router
- PostgreSQL (pgx)
- Migraciones SQL con golang-migrate
- JWT Auth (login + rutas protegidas)

## Endpoints
### Públicos
- GET /health
- GET /projects
- GET /projects/{slug}
- GET /docs/openapi.json

### Protegidos (JWT)
- POST /auth/login
- POST /projects
- PUT /projects/{slug}
- DELETE /projects/{slug}

## Variables de entorno
- DATABASE_URL
- PORT
- JWT_SECRET
- ADMIN_USER
- ADMIN_PASS


## 🔧 Test with Postman

Import this file into Postman:

👉 `docs/postman_collection.json`

Steps:
1. Run **Login** request to obtain JWT
2. Token is stored automatically
3. Use protected endpoints

Authentication:

This API uses environment-based admin credentials.

ADMIN_USER and ADMIN_PASS must be configured in the environment.
Login returns a JWT token which must be used as:

Authorization: Bearer <token>
## Local (Docker)
```bash
docker compose up -d
export DATABASE_URL="postgresql://portafolio:portafolio@127.0.0.1:5433/portafolio?sslmode=disable"
export PORT=8081
export JWT_SECRET="change_me"
export ADMIN_USER="admin"
export ADMIN_PASS="1234"

migrate -path ./migrations -database "$DATABASE_URL" up
go run ./cmd/api
