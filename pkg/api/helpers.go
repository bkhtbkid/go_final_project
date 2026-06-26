package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/bkhtbkid/go_final_project/pkg/db"
)

func checkDate(task *db.Task) error {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	if task.Date == "" {
		task.Date = today.Format(dateLayout)
	}

	t, err := time.Parse(dateLayout, task.Date)
	if err != nil {
		return errors.New("дата представлена в неверном формате")
	}

	var next string
	if task.Repeat != "" {
		next, err = NextDate(today, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	if afterNow(today, t) {
		if task.Repeat == "" {
			task.Date = today.Format(dateLayout)
		} else {
			task.Date = next
		}
	}

	return nil
}

func writeJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}
