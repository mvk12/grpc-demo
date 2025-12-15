# Instrucciones para agentes (Copilot)

Propósito: que un agente IA sea productivo rápido en este repo Go + gRPC + Postgres.

## Arquitectura (big picture)
- Entry-point: [main.go](../main.go) carga `.env` (godotenv), arma el DSN de Postgres y registra gRPC servers.
- Capas (flujo típico): gRPC → conversión proto↔modelo → repo SQL → Postgres.
  - Adaptadores gRPC: [servers/student.go](../servers/student.go), [servers/test.go](../servers/test.go)
  - Persistencia: [repositories/postgres.go](../repositories/postgres.go) (SQL crudo con `database/sql` sobre `public.*`, sin ORM)
  - Contratos: [pb/student.proto](../pb/student.proto), [pb/test.proto](../pb/test.proto) → generado en `pb/*pb.go` (no editar)
  - Modelos: [models/models.go](../models/models.go)

## Workflows críticos
- DB dev: `docker compose -f compose.yml up -d` (host `54322` → contenedor `5432`; normalmente exporta `POSTGRES_PORT=54322`).
- Servidor: `go run .` (gRPC en `GRPC_PORT=50051` por defecto).
- Checks rápidos: `go vet` y `go build ./...`.
- Protos: si cambias `.proto`, regenera (no toques `pb/*.pb.go` a mano).
  - VS Code Task: “Compile this proto” (desde el `.proto` abierto).
  - Manual:
    - `protoc --proto_path=. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pb/student.proto`
    - `pb/test.proto` importa `pb/student.proto`, así que recompila si cambias imports.

## Convenciones del repo (observadas)
- Conversión manual proto↔modelo dentro de `servers/*` (ej.: `GetStudent`/`CreateStudent` en [servers/student.go](../servers/student.go)).
- Repo se inyecta en servers (`NewStudentServer(repo)` / `NewTestServer(repo)`); existe singleton opcional en [repositories/repository.go](../repositories/repository.go) (evítalo si no hace falta).
- Errores: se devuelven “tal cual” (no hay mapeo consistente a `status.*`); respeta ese patrón.
- Streaming en [servers/test.go](../servers/test.go):
  - client-streaming: `CreateQuestions`, `EnrollStudents` (responden `SimpleStreamResponse{ok}`)
  - server-streaming: `GetStudentsPerTest` (incluye `time.Sleep(1s)` entre envíos)
  - bidireccional: `TakeTest` requiere metadata `x-test-id` (ver también [client/main.go](../client/main.go))

## Integraciones y datos
- DB init/esquema: [database/01.start.sql](../database/01.start.sql) (unicidad por email case-insensitive con `email_lower` y por `(student_id,test_id)` en enrollments).
