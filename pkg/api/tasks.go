package api

import (
	"net/http"

	"github.com/bkhtbkid/go_final_project/pkg/db"
)

type TaskResponse struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	tasks, err := db.Tasks(50, search)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, TaskResponse{Tasks: tasks})
}
