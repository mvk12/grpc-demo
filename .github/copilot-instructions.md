# Instrucciones para agentes (Copilot)

Propósito: que un agente IA sea productivo rápido en este repo Go + gRPC + Postgres.

## Arquitectura (big picture)
- Entry-point: [main.go](../main.go) carga `.env` (godotenv), construye DSN Postgres, crea repo y registra gRPC servers.
- Capas:
  - gRPC adapters: [servers/student.go](../servers/student.go), [servers/test.go](../servers/test.go)
  - Persistencia SQL cruda: [repositories/postgres.go](../repositories/postgres.go) (sin ORM, `database/sql`, esquema `public.*`)
  - Modelos Go: [models/models.go](../models/models.go)
  - Contratos: [pb/student.proto](../pb/student.proto), [pb/test.proto](../pb/test.proto) → código generado en `pb/*pb.go` (no editar a mano).

## Workflows reproducibles
- DB dev: `docker compose -f compose.yml up -d` (expone Postgres en host `54322` → en local suele requerir `POSTGRES_PORT=54322`).
- Servidor: `go run .` (gRPC por defecto en `GRPC_PORT=50051`).
- Checks rápidos: `go vet` y `go build ./...`.
- Protos:
  - VS Code Task: “Compile this proto” (desde el archivo `.proto` abierto).
  - Manual (ejemplo):
    - `protoc --proto_path=. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pb/student.proto`
    - si cambias [pb/test.proto](../pb/test.proto) recuerda que importa `pb/student.proto`.

## Conveciones del código (observadas)
- Conversión manual proto↔modelo en `servers/*` (ej.: `CreateStudent`, `GetStudent` en [servers/student.go](../servers/student.go)).
- Repositorio se inyecta en servers (`NewStudentServer(repo)`), aunque existe un singleton opcional + wrappers en [repositories/repository.go](../repositories/repository.go).
- Errores: se devuelven “tal cual” (sin mapeo consistente a `status.*`). Respeta ese patrón al editar.
- RPCs streaming: en [servers/test.go](../servers/test.go) hay client-streaming (`CreateQuestions`, `EnrollStudents`) y server-streaming (`GetStudentsPerTest`, con `time.Sleep(1s)` entre envíos).

## Integraciones y ejemplos
- Esquema DB e init: [database/01.start.sql](../database/01.start.sql) (tablas `students/tests/questions/enrollments`, unicidad por email case-insensitive y por `(student_id,test_id)`).
- Reflection gRPC está habilitado (ver [main.go](../main.go)), así que `grpcurl -plaintext localhost:50051 list` funciona.
- Requests de ejemplo: `bruno_collection/*.bru`.
