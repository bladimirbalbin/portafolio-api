# Portafolio API

API REST desarrollada en **Go** para exponer proyectos de portafolio personal.  
Diseñada con una arquitectura limpia, conexión real a **PostgreSQL**, migraciones versionadas y preparada para despliegue en entornos productivos.

Este proyecto está pensado como **pieza de presentación técnica** para backend roles.

---

## 🚀 Tech Stack

- **Go** (`net/http`)
- **Chi** (router HTTP)
- **PostgreSQL**
- **pgx / pgxpool** (driver y pool de conexiones)
- **Docker & Docker Compose**
- **golang-migrate** (migraciones SQL)
- **godotenv** (configuración por entorno en desarrollo)

---

## 📦 Características

- Health check con verificación real de base de datos
- CRUD de proyectos (lectura)
- Filtros por query params
- Paginación y ordenamiento
- Migraciones SQL versionadas
- Arquitectura por capas (domain / repository / handlers)
- OpenAPI mínimo para documentación
- Configuración por variables de entorno

---

## 📂 Estructura del proyecto

```text
portafolio-api/
├── cmd/api/                 # Entry point de la aplicación
├── internal/
│   ├── config/              # Configuración por entorno
│   ├── domain/              # Modelos de dominio
│   ├── repository/
│   │   └── postgres/        # Acceso a datos (pgx)
│   └── http/
│       ├── handlers/        # Handlers HTTP
│       └── middleware/
├── migrations/              # Migraciones SQL
├── docs/                    # OpenAPI
├── docker-compose.yml       # PostgreSQL local
├── go.mod
├── go.sum
└── README.md

