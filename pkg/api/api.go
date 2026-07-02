package api

import (
	"net/http"
	"os"
)

func Init() {
	todoPassword = os.Getenv("TODO_PASSWORD")

	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/signin", signinHandler)

	http.HandleFunc("/api/task", taskHandler)
	http.HandleFunc("/api/tasks", tasksHandler)
	http.HandleFunc("/api/task/done", taskDoneHandler)
}
