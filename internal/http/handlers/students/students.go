package students

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/mehedi0911/students-api/internal/models"
	"github.com/mehedi0911/students-api/internal/storage"
	"github.com/mehedi0911/students-api/internal/utils/response"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Creating a new student..")
		var student models.Student
		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError((err)))
		}

		//request validation

		if err := validator.New().Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors) //type casting to get expected type
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		lastInsertedId, err := storage.CreateStudent(
			student,
		)

		if err != nil {
			fmt.Printf("error from creating %s", err)
			fmt.Println(response.WriteJson(w, http.StatusInternalServerError, err))
			response.WriteJson(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprint(err)})
			return
		}
		slog.Info("User created successfully", slog.String("UserId", fmt.Sprint(lastInsertedId)))

		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastInsertedId})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := r.PathValue("id")
		slog.Info("Getting a student by", slog.String("id", id))
		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			slog.Error("error parsing id!", slog.String("id", id))
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
		}
		student, err := storage.GetStudentById(intId)
		if err != nil {
			slog.Error("error getting user!", slog.String("id", id))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
		}

		response.WriteJson(w, http.StatusOK, student)
	}
}
func GetStudentList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		slog.Info("getting all students")

		students, err := storage.GetStudentList()

		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}

		response.WriteJson(w, http.StatusOK, students)

	}
}
