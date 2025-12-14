package servers

import (
	"context"
	"errors"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/mvk12/grpc-demo/models"
	"github.com/mvk12/grpc-demo/pb"
	"github.com/mvk12/grpc-demo/repositories"
	"google.golang.org/grpc/metadata"
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

func (s *TestServer) TakeTest(stream pb.TestService_TakeTestServer) error {
	metadata, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return errors.New("no metadata found")
	}

	rawTestId := metadata.Get("x-test-id")

	if len(rawTestId) == 0 {
		return errors.New("\"x-test-id\" is required in metadata")
	}

	idStr := rawTestId[0]
	parsed, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}
	testId := int32(parsed)

	if testId <= 0 {
		return errors.New("invalid test id")
	}

	questions, err := s.repo.GetQuestionsByTestID(stream.Context(), testId)
	if err != nil {
		return err
	}

	questionIndex := 0
	current := &models.Question{}

	for {
		if questionIndex < len(questions) {
			current = questions[questionIndex]

			question := &pb.Question{
				Id:       current.ID,
				Question: current.Question,
			}

			err := stream.Send(question)

			if err != nil {
				return err
			}

			questionIndex++
		} else {
			return nil
		}

		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		log.Printf("Answer: %s", req.GetAnswer())
	}
}
