# AGENTS — Guía rápida para agentes

## Propósito
dar a un agente IA la mínima información necesaria para trabajar productivamente en este repo Go + gRPC siguiendo el estilo SimpleMA.

## Contexto rápido

Servicio gRPC principal: `main.go` carga `.env` (godotenv), arma el DSN de Postgres, crea el repositorio y registra servidores (`Student` y `Test`).

## Capas
`servers/*` (adaptadores gRPC) y `repositories/*` (acceso SQL directo con `database/sql`). Modelos en `models/models.go` y mensajes en `pb/`.

## Patrones importantes:
- Conversión manual entre `models` y mensajes protobuf dentro de `servers/*` (ver `GetStudent` / `CreateStudent` en `servers/student.go`).
- Repositorio: existe la interfaz `Repository` en `repositories/repository.go` y la implementación `PostgresRepository` en `repositories/postgres.go`.
- Hay dos estilos de uso del repo: inyección explícita en servidores (`NewStudentServer(repo)`) y helpers globales que usan un singleton `repoInstance` (funciones wrapper en `repositories/repository.go`).
- No hay ORM: las consultas usan SQL crudo con esquema `public.*`.
- Streaming: `TakeTest` (bidireccional) requiere metadata `x-test-id` (ver `servers/test.go` y el ejemplo en `client/main.go`).

## Workflows reproducibles

### Compilar protos
No edites a mano `pb/*pb.go` (son generados). Si cambias `.proto`, regenera.

Opción VS Code: Task “Compile this proto” (desde el `.proto` abierto).

```shell
protoc --proto_path=. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pb/student.proto
```

Nota: `pb/test.proto` importa `pb/student.proto`, así que si tocas `pb/test.proto` revisa imports y recompila.

### Levantar DB de desarrollo
```shell
docker compose -f compose.yml up -d
```

`compose.yml` expone Postgres en el host en `54322` → en local suele requerir `POSTGRES_PORT=54322`.

### Ejecutar servidor
```shell
go run .
```

Por defecto usa `GRPC_PORT=50051`.

### Compilar binario
```shell
go build -o bin/grpc-demo .
```

### Ejemplos de llamadas
- Ver [README.md](README.md) para `grpcurl` y la colección de Bruno.

### Reglas y precauciones al editar
- Si modificas `.proto`, siempre regenera `pb/` con `protoc` y verifica `option go_package`.
- Si tocas la DB, actualiza `database/01.start.sql` y prueba con `docker compose` para no romper la inicialización.
- Evita cambiar el paquete `pb` sin actualizar `go_package` y regenerar; esto rompe imports.
- El manejo de errores actual propaga errores crudos a gRPC; respeta ese patrón salvo que haya una razón clara para mapear a códigos `status`.

### Debug / checks rápidos:
- Comprobar dependencias y build:
```shell
go vet
go build ./...
```

### Archivos para revisar primero en un PR o cambio:
- `main.go` — inicialización y variables de entorno
- `pb/*.proto` — contratos RPC
- `servers/*` — adaptadores y conversiones proto→model
- `repositories/*` — queries SQL y esquemas `public.*`
- `database/01.start.sql` — scripts de init DB

### Checklist mínimo para PRs que afectan runtime:
- Compilar protos si se cambió `.proto`.
- Verificar `go build` y `go vet` locales.
- Probar flujo básico: `docker compose -f compose.yml up -d` → `go run .` → llamadas RPC (manual o con cliente).
