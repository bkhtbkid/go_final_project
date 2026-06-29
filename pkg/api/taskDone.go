package api

import (
	"net/http"
	"time"

	"github.com/bkhtbkid/go_final_project/pkg/db"
)

func taskDoneHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		writeJson(w, http.StatusMethodNotAllowed, map[string]string{"error": "метод не поддерживается"})
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "Задача не найдена"})
		return
	}

	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeJson(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

	} else {
		now := time.Now()
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

		next, err := NextDate(today, task.Date, task.Repeat)
		if err != nil {
			writeJson(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		if err := db.UpdateDate(next, task.ID); err != nil {
			writeJson(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
	}

	writeJson(w, http.StatusOK, map[string]string{})
}
