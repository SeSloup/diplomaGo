package api

import (
	"diplomaGoSologub/pkg/db"
	"net/http"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		writeJson(w, map[string]string{"error": "method not allowed"})
		return
	}

	tasks, err := db.Tasks(50) // в параметре максимальное количество записей
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}
	writeJson(w, TasksResp{Tasks: tasks})
}
