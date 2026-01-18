# Decisiones técnicas — Portafolio API

Este documento resume las decisiones técnicas principales del proyecto, con sus motivaciones y trade-offs.
El objetivo es evidenciar criterio de ingeniería backend: diseño, mantenibilidad, despliegue y operación.

---

## Objetivo del proyecto

Construir una API REST en Go, desplegada públicamente, conectada a PostgreSQL real, con estructura limpia, documentación mínima (OpenAPI) y enfoque en calidad (migraciones + testing), pensada como pieza verificable para entrevistas.

Demo:
- https://portafolio-api-xiw6.onrender.com

Repositorio:
- https://github.com/bladimirbalbin/portafolio-api

---

## Stack y por qué

- **Go + net/http**
  - Razón: runtime simple, performance y facilidad para construir servicios pequeños y claros.
- **Chi**
  - Razón: router liviano, middleware estándar, ergonomía para APIs.
- **PostgreSQL**
  - Razón: base robusta, común en empresas, queries explícitas, tipos (arrays para tags).
- **pgx / pgxpool**
  - Razón: driver moderno, buen performance, pooling explícito.
- **golang-migrate**
  - Razón: migraciones SQL versionadas y reproducibles en dev/prod/test.
- **Docker Compose (local)**
  - Razón: reproducibilidad, entorno local rápido para DB.
- **Render (deploy)**
  - Razón: setup rápido, DB gestionada y web service sin fricción para demo pública.

---

## Arquitectura (por capas)

Se eligió una separación simple y explícita:

- `cmd/api/`  
  Entry point. Ensambla dependencias (config, DB pool, router) y arranca el servidor.

- `internal/config/`  
  Lectura/validación de variables de entorno (ej. `DATABASE_URL`).

- `internal/domain/`  
  Modelos de dominio (ej. `Project`). Sin dependencias de infraestructura.

- `internal/repository/postgres/`  
  Acceso a datos (SQL + pgx). Encapsula queries y mapea rows → domain.

- `internal/http/handlers/`  
  Handlers HTTP. Parsean input (params/query), llaman repo y devuelven respuestas JSON.

**Trade-off:** no es Clean Architecture “completa” con múltiples interfaces por todo lado; se priorizó claridad y tamaño acotado (portafolio), manteniendo separación suficiente para testear y crecer.

---

## Diseño de endpoints y contratos

- `GET /health`
  - Verifica que el servicio está vivo y que PostgreSQL responde.
- `GET /projects`
  - Lista proyectos con paginación y filtros (limit/offset/tag/featured/sort).
- `GET /projects/{slug}`
  - Recurso individual, con **404 real** cuando no existe.

**Trade-off:** el MVP se enfocó en lectura (read-only). CRUD completo y auth quedaron fuera para evitar complejidad innecesaria.

---

## Manejo de errores y códigos HTTP

- `404` cuando un `slug` no existe (`ErrNotFound` en repository → mapeo a 404 en handler)
- `500` para errores inesperados
- Respuestas JSON consistentes (`{ "data": ... }` o `{ "error": ... }`)

**Motivación:** API predecible y fácil de consumir; errores claramente diferenciados.

---

## SQL explícito (sin ORM)

Se decidió **no usar ORM**.

**Ventajas**
- Control total del SQL, índices, ORDER BY, filtros.
- Debug más directo (lo que se ejecuta es evidente).
- Menos magia en un portafolio de backend.

**Trade-off**
- Más responsabilidad en validación, queries y mapeo.
- Requiere disciplina para evitar SQL injection (se usó whitelist para ORDER BY y parámetros bind).

---

## Migraciones

- Migraciones SQL versionadas en `migrations/`.
- Las migraciones se ejecutan en:
  - Local (Docker)
  - Producción (Render) manualmente (para controlar el momento)

**Trade-off:** Render no ejecuta migraciones automáticamente; se asumió como parte del proceso operacional (similar a CI/CD real).

---

## Testing (base del MVP)

Se implementó una base de tests con dos estrategias:

1) **Repository tests con PostgreSQL real**
   - Con una DB de test en Docker (`docker-compose.test.yml`)
   - Migraciones aplicadas antes de correr tests
   - Objetivo: asegurar SQL, scan y errores (por ejemplo `ErrNotFound`)

2) **Handler tests con `httptest`**
   - Sin DB: repos falsos (fakes) para controlar respuestas
   - Objetivo: validar status codes (200/404/500) y formato JSON

**Motivación:** separar infraestructura (DB) de lógica HTTP, y asegurar confiabilidad con bajo costo.

---

## Operación y configuración

- Configuración por variables de entorno:
  - `DATABASE_URL`
  - `PORT` (inyectado por Render)
- En Render:
  - DB y servicio en la misma región (para internal networking)
  - Se usa la **Internal Database URL** en producción

---

## Posibles mejoras (no incluidas en el MVP)

Estas mejoras se dejaron fuera intencionalmente para mantener el MVP pequeño y completo:

- CI con GitHub Actions (test + lint)
- Lint/format (golangci-lint)
- Observabilidad (structured logs, métricas)
- CRUD completo y auth
- Seeds versionados (datos demo automatizados)

---

## Resumen

Este MVP prioriza:
- Deploy real
- PostgreSQL real
- SQL explícito
- Arquitectura simple por capas
- Migraciones reproducibles
- Tests básicos y útiles
- Documentación mínima pero verificable

El resultado es un servicio pequeño, claro y defendible técnicamente en entrevistas.
