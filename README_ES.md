# grpc-demo (Go + gRPC + Postgres)

Proyecto demo para mostrar una API gRPC en Go usando Postgres como persistencia.

- Versión principal (English): [README.md](README.md)

## ¿Qué incluye?

- **StudentService**: crear estudiantes y obtener por id.
- **TestService**: crear/obtener tests, crear preguntas (client-streaming), inscribir estudiantes (client-streaming) y listar estudiantes por test (server-streaming).
- **Postgres**: SQL crudo con `database/sql` (sin ORM) + script de inicialización.

## Arquitectura (alto nivel)

- Entry point: [main.go](main.go)
  - Carga `.env` con `godotenv` (opcional).
  - Construye el DSN de Postgres desde variables de entorno.
  - Crea el repositorio y registra los servidores gRPC.
  - Habilita gRPC reflection (por eso `grpcurl ... list` funciona).
- Adaptadores gRPC: [servers/student.go](servers/student.go), [servers/test.go](servers/test.go)
  - Conversión manual entre mensajes protobuf (`pb/*`) y modelos (`models/*`).
- Interfaz del repositorio: [repositories/repository.go](repositories/repository.go)
  - El patrón principal es inyección en los servers (`NewStudentServer(repo)`), pero también existen wrappers globales sobre un singleton `repoInstance`.
- Implementación Postgres: [repositories/postgres.go](repositories/postgres.go)
  - SQL escrito contra tablas `public.*`.
- Contratos protobuf (el código generado vive en `pb/*.pb.go`; no edites esos archivos a mano):
  - [pb/student.proto](pb/student.proto)
  - [pb/test.proto](pb/test.proto) (importa `pb/student.proto`)

## Requisitos

- Go (según [go.mod](go.mod))
- Docker + Docker Compose (para la DB de desarrollo)
- Opcional:
  - `protoc` + plugins Go (solo si cambias `.proto`)
  - `grpcurl` (recomendado para probar rápido)
  - Bruno (opcional, para la colección incluida)

## Inicio rápido

### 1) Levantar Postgres (Docker)

```bash
docker compose -f compose.yml up -d
```

Por defecto, Postgres se expone en el host en el puerto `54322` (ver [compose.yml](compose.yml)).

### 2) Configurar variables de entorno (puertos por env)

El servidor lee estas variables (con sus defaults):

- `GRPC_PORT` (default: `50051`)
- `POSTGRES_USER` (default: `user`)
- `POSTGRES_PASSWORD` (default: `password`)
- `POSTGRES_HOST` (default: `localhost`)
- `POSTGRES_PORT` (default: `5432`) ← usando Docker Compose en este repo normalmente se usa `54322`
- `POSTGRES_DB` (default: `grpc_db`)
- `POSTGRES_SSLMODE` (default: `disable`)

Ejemplo (recomendado con Docker Compose):

```bash
export POSTGRES_PORT=54322
export GRPC_PORT=50051
```

También puedes ponerlas en un archivo `.env` (se carga automáticamente).

### 3) Ejecutar el servidor gRPC

```bash
go run .
```

## Probar con grpcurl

Reflection está habilitado en [main.go](main.go), así que estos comandos deberían funcionar:

```bash
grpcurl -plaintext localhost:50051 list
grpcurl -plaintext localhost:50051 describe student.StudentService
grpcurl -plaintext localhost:50051 describe test.TestService
```

### Crear estudiante

```bash
grpcurl -plaintext \
  -d '{"name":"Alice","email":"alice@example.com"}' \
  localhost:50051 student.StudentService/CreateStudent
```

### Obtener estudiante por id

```bash
grpcurl -plaintext \
  -d '{"id":1}' \
  localhost:50051 student.StudentService/GetStudent
```

### Crear test

```bash
grpcurl -plaintext \
  -d '{"title":"Math basics","description":"Demo test"}' \
  localhost:50051 test.TestService/CreateTest
```

### Enviar preguntas (client-streaming)

```bash
cat <<'JSON' | grpcurl -plaintext -d @ localhost:50051 test.TestService/CreateQuestions
{"question":"2+2?","answer":"4","test_id":1}
{"question":"3+5?","answer":"8","test_id":1}
JSON
```

### Inscribir estudiantes (client-streaming)

```bash
cat <<'JSON' | grpcurl -plaintext -d @ localhost:50051 test.TestService/EnrollStudents
{"student_id":1,"test_id":1}
{"student_id":2,"test_id":1}
JSON
```

### Listar estudiantes por test (server-streaming)

```bash
grpcurl -plaintext \
  -d '{"test_id":1}' \
  localhost:50051 test.TestService/GetStudentsPerTest
```

## Esquema de BD

La DB se inicializa con [database/01.start.sql](database/01.start.sql) vía `docker-entrypoint-initdb.d` (montado en [compose.yml](compose.yml)).

Restricciones destacables:

- Email de estudiantes único **case-insensitive** (columna generada `email_lower`).
- Inscripciones únicas por `(student_id, test_id)`.

## Regenerar Protobufs

Los generados viven en `pb/*.pb.go` y no deben editarse a mano.

- VS Code: Task “Compile this proto” con el `.proto` abierto.
- Manual:

```bash
protoc --proto_path=. \
  --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  pb/student.proto
```

Si cambias [pb/test.proto](pb/test.proto), recuerda que importa `pb/student.proto`.

## Requests de ejemplo (Bruno)

Hay una colección Bruno en [bruno_collection/](bruno_collection/) para ejecutar requests de ejemplo.

## Checks de desarrollo

```bash
go vet
go build ./...
```

## Troubleshooting

- **No conecta a la DB**: si usas Docker Compose, configura `POSTGRES_PORT=54322` (compose mapea host `54322` → contenedor `5432`).
- **Errores de unicidad**:
  - Email de estudiante único (case-insensitive).
  - Enrollment único por `(student_id, test_id)`.
- **Cambiaste `.proto` pero no se actualiza Go**: vuelve a ejecutar `protoc` (o la task de VS Code) y evita editar `pb/*.pb.go`.
