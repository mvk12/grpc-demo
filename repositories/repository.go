package repositories

import (
	"context"

	"github.com/mvk12/grpc-demo/models"
)

type Repository interface {
	GetStudentByID(ctx context.Context, id int32) (*models.Student, error)
	CreateStudent(ctx context.Context, student *models.Student) (*models.Student, error)
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
