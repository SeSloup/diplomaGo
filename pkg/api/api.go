package api

import (
	"diplomaGoSologub/pkg/auth"
	"net/http"
)

func Init() {
	http.HandleFunc("/api/signin", auth.SigninHandler)
	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/task", auth.Auth(taskHandler))
	http.HandleFunc("/api/tasks", auth.Auth(tasksHandler))
	http.HandleFunc("/api/task/done", auth.Auth(doneTaskHandler))

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
