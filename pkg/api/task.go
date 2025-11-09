package api

import (
	"database/sql"
	"fmt"
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
