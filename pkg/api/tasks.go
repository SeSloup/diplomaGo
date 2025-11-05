package api

import (
	"diplomaGoSologub/pkg/db"
	"net/http"
)

func Tasks(limit int) ([]*db.Task, error) {
	row, err := db.DB.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT ?",
		limit,
	)
	if err != nil {
		return nil, err
	}

	defer row.Close()

	var tasks = make([]*db.Task, 0, limit)
	for row.Next() {
		var task db.Task
		if err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}
	if err := row.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		writeJson(w, map[string]string{"error": "method not allowed"})
		return
	}

	tasks, err := Tasks(50) // в параметре максимальное количество записей
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}
	writeJson(w, TasksResp{Tasks: tasks})
}
