package servers

import (
	"context"

	"github.com/mvk12/grpc-demo/models"
	"github.com/mvk12/grpc-demo/pb"
	"github.com/mvk12/grpc-demo/repositories"
)

type StudentServer struct {
	repo repositories.Repository
	pb.UnimplementedStudentServiceServer
}

func NewStudentServer(repo repositories.Repository) *StudentServer {
	return &StudentServer{repo: repo}
}

func (s *StudentServer) GetStudent(ctx context.Context, req *pb.GetStudentRequest) (*pb.Student, error) {
	student, err := s.repo.GetStudentByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &pb.Student{
		Id:    student.ID,
		Name:  student.Name,
		Email: student.Email,
	}, nil
}

func (s *StudentServer) CreateStudent(ctx context.Context, req *pb.Student) (*pb.StudentResponse, error) {
	student, err := s.repo.CreateStudent(ctx, &models.Student{
		Name:  req.GetName(),
		Email: req.GetEmail(),
	})
	if err != nil {
		return nil, err
	}

	return &pb.StudentResponse{
		Id: student.ID,
	}, nil
}
