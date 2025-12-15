
# grpc-demo (Go + gRPC + Postgres)

Demo project showcasing a simple gRPC API in Go backed by Postgres.

- Spanish version: [README_ES.md](README_ES.md)

## What’s inside

- **StudentService**: create a student and fetch by id.
- **TestService**: create/get tests, add questions (client-streaming), enroll students (client-streaming), list students per test (server-streaming), and take a test (bidirectional streaming; requires `x-test-id` metadata).
- **Postgres**: raw `database/sql` queries (no ORM) + init SQL script.

## Architecture (high level)

- Entry point: [main.go](main.go)
	- Loads `.env` via `godotenv` (optional).
	- Builds a Postgres DSN from environment variables.
	- Creates a repository and registers gRPC servers.
	- Enables gRPC reflection (so `grpcurl ... list` works).
- gRPC adapters: [servers/student.go](servers/student.go), [servers/test.go](servers/test.go)
	- Manual conversion between protobuf messages (`pb/*`) and Go models (`models/*`).
- Repository interface: [repositories/repository.go](repositories/repository.go)
	- Primary usage is dependency injection into servers (`NewStudentServer(repo)`), but there are also optional global wrappers around a singleton `repoInstance`.
- Postgres implementation: [repositories/postgres.go](repositories/postgres.go)
	- SQL is written against `public.*` tables.
- Protobuf contracts (generated code lives in `pb/*.pb.go`; do not edit generated files):
	- [pb/student.proto](pb/student.proto)
	- [pb/test.proto](pb/test.proto) (imports `pb/student.proto`)

## Prerequisites

- Go (as declared in [go.mod](go.mod))
- Docker + Docker Compose (for the dev database)
- Optional tooling:
	- `protoc` + Go plugins (only needed if you change `.proto` files)
	- `grpcurl` (recommended for quick manual testing)
	- Bruno (optional, for running the included request collection)

## Quickstart

### 1) Start Postgres (Docker)

```bash
docker compose -f compose.yml up -d
```

By default, Postgres is exposed on host port `54322` (see [compose.yml](compose.yml)).

### 2) Configure env vars (ports via env)

The server reads these environment variables (defaults shown):

- `GRPC_PORT` (default: `50051`)
- `POSTGRES_USER` (default: `user`)
- `POSTGRES_PASSWORD` (default: `password`)
- `POSTGRES_HOST` (default: `localhost`)
- `POSTGRES_PORT` (default: `5432`) ← when using Docker Compose in this repo, you typically set it to `54322`
- `POSTGRES_DB` (default: `grpc_db`)
- `POSTGRES_SSLMODE` (default: `disable`)

Example (recommended for Docker Compose DB):

```bash
export POSTGRES_PORT=54322
export GRPC_PORT=50051
```

You can also put them in a `.env` file (loaded automatically).

### 3) Run the gRPC server

```bash
go run .
```

## Try it (grpcurl)

Reflection is enabled in [main.go](main.go), so these should work:

```bash
grpcurl -plaintext localhost:50051 list
grpcurl -plaintext localhost:50051 describe student.StudentService
grpcurl -plaintext localhost:50051 describe test.TestService
```

### Create a student

```bash
grpcurl -plaintext \
	-d '{"name":"Alice","email":"alice@example.com"}' \
	localhost:50051 student.StudentService/CreateStudent
```

### Fetch a student by id

```bash
grpcurl -plaintext \
	-d '{"id":1}' \
	localhost:50051 student.StudentService/GetStudent
```

### Create a test

```bash
grpcurl -plaintext \
	-d '{"title":"Math basics","description":"Demo test"}' \
	localhost:50051 test.TestService/CreateTest
```

### Stream questions into a test (client-streaming)

`grpcurl` can send a stream from stdin using `-d @`:

```bash
cat <<'JSON' | grpcurl -plaintext -d @ localhost:50051 test.TestService/CreateQuestions
{"question":"2+2?","answer":"4","test_id":1}
{"question":"3+5?","answer":"8","test_id":1}
JSON
```

### Enroll students (client-streaming)

```bash
cat <<'JSON' | grpcurl -plaintext -d @ localhost:50051 test.TestService/EnrollStudents
{"student_id":1,"test_id":1}
{"student_id":2,"test_id":1}
JSON
```

### List students per test (server-streaming)

```bash
grpcurl -plaintext \
	-d '{"test_id":1}' \
	localhost:50051 test.TestService/GetStudentsPerTest
```

### Take a test (bidirectional streaming)

`grpcurl` is not a great fit for bidirectional streaming. This repo includes a demo client in [client/main.go](client/main.go) that sets the required metadata `x-test-id`.

```bash
go run ./client
```

Note: the demo client uses hardcoded IDs (for example `test_id=1`). If your DB already has data, you may need to update the IDs in [client/main.go](client/main.go) or reset the Docker volume.

## DB schema

The DB is initialized from [database/01.start.sql](database/01.start.sql) via Docker’s `docker-entrypoint-initdb.d` (mounted in [compose.yml](compose.yml)).

Notable constraints:

- Students email is unique **case-insensitively** via a generated column `email_lower`.
- Enrollments are unique by `(student_id, test_id)`.

## Protobuf generation

Generated files live under `pb/*.pb.go` and should not be edited manually.

- VS Code: run the task “Compile this proto” while a `.proto` file is open.
- Manual example:

```bash
protoc --proto_path=. \
	--go_out=. --go_opt=paths=source_relative \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
	pb/student.proto
```

If you change [pb/test.proto](pb/test.proto), remember it imports `pb/student.proto`.

## Example requests (Bruno)

There is a Bruno collection under [bruno_collection/](bruno_collection/) you can use to exercise the RPCs.

## Development checks

```bash
go vet
go build ./...
```

## Troubleshooting

- **DB connection refused**: if you use Docker Compose, set `POSTGRES_PORT=54322` (compose maps host `54322` → container `5432`).
- **Unique constraint errors**:
	- Student emails must be unique case-insensitively.
	- Enrollment `(student_id, test_id)` must be unique.
- **Changed `.proto` but Go code didn’t update**: re-run `protoc` (or the VS Code task) and avoid manual edits in `pb/*.pb.go`.

