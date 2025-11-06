package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"diplomaGoSologub/pkg/db"
)

func AddTask(task *db.Task) (int64, error) {
	var id int64
	// определите запрос
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}

func GetTask(id string) (*db.Task, error) {
	var task db.Task

	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`

	err := db.DB.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

	fmt.Println(task)

	if err != nil {
		if err == sql.ErrNoRows {

			return nil, err
		}

	}
	return &task, nil

}

func UpdateTask(task *db.Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := db.DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("incorrect id for updating task")
	}
	return nil
}

func DeleteTask(id string) error {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	res, err := db.DB.Exec("DELETE FROM scheduler WHERE id = ?", idInt)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("id %v is not founded", id)
	}
	return nil
}

func checkDate(task *db.Task) error {
	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format(dateFormat)
	}

	taskDate, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		return fmt.Errorf("invalid date format, expected YYYYMMDD")
	}

	if taskDate.Before(now) {
		if task.Repeat != "" {
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return err
			}
			task.Date = next
		} else {
			task.Date = now.Format(dateFormat)
		}
	}
	return nil
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method is not allowed", http.StatusMethodNotAllowed)
		return
	}
	var task db.Task

	if err := readJson(r.Body, &task); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	if task.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": "title is empty"})

		return
	}

	if err := checkDate(&task); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": err.Error()})

		return
	}

	id, err := AddTask(&task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, map[string]string{"id": fmt.Sprint(id)})
}

func getByIdTaskHandler(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusInternalServerError)
		writeJson(w, map[string]string{"error": "Не указан идентификатор задачи"})
		return

	}
	task, err := GetTask(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Задача не найдена",
		})
		return
	}

	writeJson(w, task)

}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "method is not allowed", http.StatusMethodNotAllowed)
		return
	}
	var task db.Task

	if err := readJson(r.Body, &task); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	if task.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": "title is empty"})

		return
	}

	if err := checkDate(&task); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": err.Error()})

		return
	}

	err := UpdateTask(&task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, map[string]string{"message": "update scheduler success"})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": "id not defined"})
		return
	}

	err := DeleteTask(id)
	if err != nil {
		if errors.Is(err, fmt.Errorf("id %v is not founded", id)) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}
	writeJson(w, map[string]string{}) // успешное выполнение должно вернуть {}
}

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		writeJson(w, map[string]string{"error": "method not allowed"})
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": "id not defined"})
		return
	}

	nowStr := r.URL.Query().Get("now")
	var now time.Time
	if nowStr == "" {
		now = time.Now()
	} else {
		tm, err := time.Parse(dateFormat, nowStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			writeJson(w, map[string]string{"error": "invalid now format"})
			return
		}
		now = tm
	}

	task, err := GetTask(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	if task.Repeat == "" {
		if err := DeleteTask(id); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			writeJson(w, map[string]string{"error": err.Error()})
			return
		}
		writeJson(w, map[string]string{})
		return
	}

	next, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}
	task.Date = next

	if err := UpdateTask(task); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, map[string]string{"message": "done and update task complete"})
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// обработка других методов будет добавлена на следующих шагах
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getByIdTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	}
}
