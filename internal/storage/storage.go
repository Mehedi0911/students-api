package storage

import (
	"github.com/mehedi0911/students-api/internal/models"
)

type Storage interface {
	CreateStudent(payload models.Student) (int64, error)
	GetStudentById(id int64) (models.Student, error)
	GetStudentList() ([]models.Student, error)
	UpdateStudent(payload models.Student, id int64) (int64, error)
	DeleteStudent(id int64) (int64, error)
}
