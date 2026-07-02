package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	dateLayout = "20060102"
	maxDays    = 400
)

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("правило повторения не может быть пустым")
	}

	date, err := time.Parse(dateLayout, dstart)
	if err != nil {
		return "", err
	}

	parts := strings.Split(repeat, " ")
	rule := parts[0]

	switch rule {
	case "y":
		return nextYear(now, date), nil
	case "d":
		return nextDay(now, date, parts)
	case "w":
		return nextWeek(now, date, parts)
	case "m":
		return nextMonth(now, date, parts)
	default:
		return "", fmt.Errorf("неверный формат для правила: %q", repeat)
	}
}

func afterNow(date, now time.Time) bool {
	return date.After(now)
}

func nextYear(now, date time.Time) string {
	for {
		date = date.AddDate(1, 0, 0)
		if afterNow(date, now) {
			break
		}
	}
	return date.Format(dateLayout)
}

func nextDay(now, date time.Time, parts []string) (string, error) {
	if len(parts) < 2 {
		return "", errors.New("не указан интервал")
	}
	days, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", err
	}
	if days < 1 || days > maxDays {
		return "", fmt.Errorf("неверный интервал: %d", days)
	}
	for {
		date = date.AddDate(0, 0, days)
		if afterNow(date, now) {
			break
		}
	}
	return date.Format(dateLayout), nil
}

func nextWeek(now, date time.Time, parts []string) (string, error) {
	if len(parts) < 2 {
		return "", errors.New("не указаны дни недели")
	}
	var allowed [8]bool
	for _, s := range strings.Split(parts[1], ",") {
		wd, err := strconv.Atoi(s)
		if err != nil {
			return "", err
		}
		if wd < 1 || wd > 7 {
			return "", fmt.Errorf("неверный день недели: %d", wd)
		}
		allowed[wd] = true
	}
	for !(afterNow(date, now) && allowed[weekdayISO(date)]) {
		date = date.AddDate(0, 0, 1)
	}
	return date.Format(dateLayout), nil
}

func weekdayISO(t time.Time) int {
	wd := int(t.Weekday())
	if wd == 0 {
		return 7
	}
	return wd
}

func nextMonth(now, date time.Time, parts []string) (string, error) {
	if len(parts) < 2 {
		return "", errors.New("не указаны дни месяца")
	}

	var days [32]bool
	var fromEnd [3]bool
	for _, s := range strings.Split(parts[1], ",") {
		d, err := strconv.Atoi(s)
		if err != nil {
			return "", err
		}
		switch {
		case d >= 1 && d <= 31:
			days[d] = true
		case d == -1 || d == -2:
			fromEnd[-d] = true
		default:
			return "", fmt.Errorf("неверный день месяца: %d", d)
		}
	}

	var months [13]bool
	monthFilter := len(parts) >= 3
	if monthFilter {
		for _, s := range strings.Split(parts[2], ",") {
			m, err := strconv.Atoi(s)
			if err != nil {
				return "", err
			}
			if m < 1 || m > 12 {
				return "", fmt.Errorf("неверный месяц: %d", m)
			}
			months[m] = true
		}
	}

	matches := func(t time.Time) bool {
		if monthFilter && !months[int(t.Month())] {
			return false
		}
		last := lastDayOfMonth(t)
		switch {
		case days[t.Day()]:
			return true
		case fromEnd[1] && t.Day() == last:
			return true
		case fromEnd[2] && t.Day() == last-1:
			return true
		}
		return false
	}

	for !(afterNow(date, now) && matches(date)) {
		date = date.AddDate(0, 0, 1)
	}
	return date.Format(dateLayout), nil
}

func lastDayOfMonth(t time.Time) int {
	return time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, t.Location()).Day()
}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	if v := r.FormValue("now"); v != "" {
		parsed, err := time.Parse(dateLayout, v)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		now = parsed
	}

	next, err := NextDate(now, r.FormValue("date"), r.FormValue("repeat"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte(next))
}
