package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/mehedi0911/students-api/internal/models"
)

func (s *Sqlite) CreateStudent(payload models.Student) (int64, error) {
	// check if student already exists (email must be unique)
	var exists int
	err := s.Db.QueryRow("SELECT COUNT(1) FROM students WHERE email = ?", payload.Email).Scan(&exists)
	if err != nil {
		return 0, err
	}

	if exists > 0 {
		return 0, fmt.Errorf("student with email %s already exists", payload.Email)
	}

	stmt, err := s.Db.Prepare("INSERT INTO students (name, email, age) VALUES(?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(payload.Name, payload.Email, payload.Age)

	if err != nil {
		return 0, err
	}

	lastInsertedId, err := res.LastInsertId()

	if err != nil {
		return 0, err
	}

	return lastInsertedId, nil
}

func (s *Sqlite) GetStudentById(id int64) (models.Student, error) {
	stmt, err := s.Db.Prepare("SELECT * FROM students WHERE id = ? LIMIT 1")
	if err != nil {
		return models.Student{}, err
	}
	defer stmt.Close()

	var student models.Student

	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Student{}, fmt.Errorf("no students found with id %s", fmt.Sprint(id))
		}
		return models.Student{}, fmt.Errorf("query error %w", err)
	}

	return student, nil
}
