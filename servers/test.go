package servers

import (
	"context"
	"io"
	"time"

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

func (s *TestServer) CreateQuestions(stream pb.TestService_CreateQuestionsServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.SimpleStreamResponse{
				Ok: true,
			})
		}

		if err != nil {
			return err
		}

		questionModel := &models.Question{
			Question: req.GetQuestion(),
			Answer:   req.GetAnswer(),
			TestID:   req.GetTestId(),
		}

		_, err = s.repo.CreateQuestion(stream.Context(), questionModel)
		if err != nil {
			return stream.SendAndClose(&pb.SimpleStreamResponse{
				Ok: false,
			})
		}
	}
}

func (s *TestServer) EnrollStudents(stream pb.TestService_EnrollStudentsServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.SimpleStreamResponse{
				Ok: true,
			})
		}

		if err != nil {
			return err
		}

		enrollment := &models.Enrollment{
			StudentID: req.GetStudentId(),
			TestID:    req.GetTestId(),
		}

		_, err = s.repo.CreateEnrollment(stream.Context(), enrollment)
		if err != nil {
			return stream.SendAndClose(&pb.SimpleStreamResponse{
				Ok: false,
			})
		}
	}
}

func (s *TestServer) GetStudentsPerTest(req *pb.GetStudentsPerTestRequest, stream pb.TestService_GetStudentsPerTestServer) error {
	students, err := s.repo.GetStudentsPerTest(stream.Context(), req.GetTestId())
	if err != nil {
		return err
	}

	for _, student := range students {
		student := &pb.Student{
			Id:    student.ID,
			Name:  student.Name,
			Email: student.Email,
		}

		err := stream.Send(student)

		if err != nil {
			return err
		}

		time.Sleep(1 * time.Second)
	}

	return nil
}
