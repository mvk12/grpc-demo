package servers

import (
	"context"

	"github.com/mvk12/grpc-demo/models"
	"github.com/mvk12/grpc-demo/pb"
	"github.com/mvk12/grpc-demo/repositories"
)

type TestServer struct {
	repo repositories.Repository
	pb.UnimplementedTestServiceServer
}

func NewTestServer(repo repositories.Repository) *TestServer {
	return &TestServer{repo: repo}
}

func (s *TestServer) GetTest(ctx context.Context, req *pb.GetTestRequest) (*pb.Test, error) {
	test, err := s.repo.GetTestByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &pb.Test{
		Id:          test.ID,
		Title:       test.Title,
		Description: test.Description,
	}, nil
}

func (s *TestServer) CreateTest(ctx context.Context, req *pb.Test) (*pb.TestResponse, error) {
	test, err := s.repo.CreateTest(ctx, &models.Test{
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
	})
	if err != nil {
		return nil, err
	}

	return &pb.TestResponse{
		Id: test.ID,
	}, nil
}
