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

func (r *PostgresRepository) CreateTest(ctx context.Context, test *models.Test) (*models.Test, error) {
	var id int64
	err := r.db.QueryRowContext(ctx, "INSERT INTO public.tests (title, description) VALUES ($1, $2) RETURNING id", test.Title, test.Description).Scan(&id)
	if err != nil {
		return nil, err
	}

	test.ID = int32(id)

	return test, nil
}

func (r *PostgresRepository) GetTestByID(ctx context.Context, id int32) (*models.Test, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, title, description FROM public.tests WHERE id = $1 LIMIT 1", id)
	var test = models.Test{}
	err := row.Scan(&test.ID, &test.Title, &test.Description)
	if err != nil {
		return nil, err
	}

	return &test, nil
}

func (r *PostgresRepository) CreateQuestion(ctx context.Context, question *models.Question) (*models.Question, error) {
	var id int64
	err := r.db.QueryRowContext(ctx, "INSERT INTO public.questions (question, answer, test_id) VALUES ($1, $2, $3) RETURNING id", question.Question, question.Answer, question.TestID).Scan(&id)
	if err != nil {
		return nil, err
	}

	question.ID = int32(id)

	return question, nil
}

func (r *PostgresRepository) CreateEnrollment(ctx context.Context, enrollment *models.Enrollment) (*models.Enrollment, error) {
	var id int64
	err := r.db.QueryRowContext(ctx, "INSERT INTO public.enrollments (student_id, test_id) VALUES ($1, $2) RETURNING id", enrollment.StudentID, enrollment.TestID).Scan(&id)
	if err != nil {
		return nil, err
	}

	enrollment.ID = int32(id)

	return enrollment, nil
}

func (r *PostgresRepository) GetStudentsPerTest(ctx context.Context, testId int32) ([]*models.Student, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT s.id, s.name, s.email
		FROM public.students s
		JOIN public.enrollments e ON s.id = e.student_id
		WHERE e.test_id = $1
	`, testId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var students []*models.Student
	for rows.Next() {
		var student models.Student
		if err := rows.Scan(&student.ID, &student.Name, &student.Email); err == nil {
			students = append(students, &student)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return students, nil
}
