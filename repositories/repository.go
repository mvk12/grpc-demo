package repositories

import (
	"context"

	"github.com/mvk12/grpc-demo/models"
)

type Repository interface {
	GetStudentByID(ctx context.Context, id int32) (*models.Student, error)
	CreateStudent(ctx context.Context, student *models.Student) (*models.Student, error)
	GetTestByID(ctx context.Context, id int32) (*models.Test, error)
	CreateTest(ctx context.Context, test *models.Test) (*models.Test, error)
	CreateQuestion(ctx context.Context, question *models.Question) (*models.Question, error)
}

var repoInstance Repository

func SetRepository(r Repository) {
	repoInstance = r
}

func CreateStudent(ctx context.Context, student *models.Student) (*models.Student, error) {
	return repoInstance.CreateStudent(ctx, student)
}

func GetStudentByID(ctx context.Context, id int32) (*models.Student, error) {
	return repoInstance.GetStudentByID(ctx, id)
}

func CreateTest(ctx context.Context, test *models.Test) (*models.Test, error) {
	return repoInstance.CreateTest(ctx, test)
}

func GetTestByID(ctx context.Context, id int32) (*models.Test, error) {
	return repoInstance.GetTestByID(ctx, id)
}

func CreateQuestion(ctx context.Context, question *models.Question) (*models.Question, error) {
	return repoInstance.CreateQuestion(ctx, question)
}