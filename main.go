package main

import (
	"fmt"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/mvk12/grpc-demo/pb"
	"github.com/mvk12/grpc-demo/repositories"
	"github.com/mvk12/grpc-demo/servers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func main() {
	_ = godotenv.Load()
	grpcPort := getenv("GRPC_PORT", "50051")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))

	if err != nil {
		panic(err)
	}

	defer lis.Close()

	user := getenv("POSTGRES_USER", "user")
	password := getenv("POSTGRES_PASSWORD", "password")
	host := getenv("POSTGRES_HOST", "localhost")
	port := getenv("POSTGRES_PORT", "5432")
	dbName := getenv("POSTGRES_DB", "grpc_db")
	sslmode := getenv("POSTGRES_SSLMODE", "disable")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, password, host, port, dbName, sslmode)

	repo, err := repositories.NewPostgresRepository(dsn)
	if err != nil {
		panic(err)
	}

	studentServer := servers.NewStudentServer(repo)
	testServer := servers.NewTestServer(repo)

	s := grpc.NewServer()
	pb.RegisterStudentServiceServer(s, studentServer)
	pb.RegisterTestServiceServer(s, testServer)

	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}
