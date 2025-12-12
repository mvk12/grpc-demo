package repositories

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/mvk12/grpc-demo/models"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)

	if err != nil {
		return nil, err
	}

	return &PostgresRepository{db}, nil
}

func (r *PostgresRepository) Close() error {
	return r.db.Close()
}

func (r *PostgresRepository) CreateStudent(ctx context.Context, student *models.Student) (*models.Student, error) {
	var id int64
	err := r.db.QueryRowContext(ctx, "INSERT INTO public.students (name, email) VALUES ($1, $2) RETURNING id", student.Name, student.Email).Scan(&id)
	if err != nil {
		return nil, err
	}

	student.ID = int32(id)

	return student, nil
}

func (r *PostgresRepository) GetStudentByID(ctx context.Context, id int32) (*models.Student, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, name, email FROM public.students WHERE id = $1 LIMIT 1", id)
	var student = models.Student{}
	err := row.Scan(&student.ID, &student.Name, &student.Email)
	if err != nil {
		return nil, err
	}

	return &student, nil
}
