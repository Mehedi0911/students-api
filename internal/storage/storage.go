package storage

import (
	"github.com/mehedi0911/students-api/internal/models"
)

type Storage interface {
	CreateStudent(payload models.Student) (int64, error)
}
