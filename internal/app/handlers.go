package app

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type App struct {
	DB *gorm.DB
}

func (a *App) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/questions/", a.questionsHandler)
	mux.HandleFunc("/questions", a.questionsHandler)

	mux.HandleFunc("/answers/", a.answersHandler)
	mux.HandleFunc("/answers", a.answersHandler)

	return mux
}

func readJSON(r *http.Request, v any) error {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if len(b) == 0 {
		return errors.New("empty body")
	}
	return json.Unmarshal(b, v)
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

//
// QUESTIONS
//

func (a *App) questionsHandler(w http.ResponseWriter, r *http.Request) {
	p := strings.TrimPrefix(r.URL.Path, "/questions")
	p = strings.TrimSuffix(p, "/")

	if p == "" {
		switch r.Method {
		case http.MethodGet:
			a.listQuestions(w, r)
		case http.MethodPost:
			a.createQuestion(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		return
	}

	id, err := strconv.Atoi(strings.TrimPrefix(p, "/"))
	if err != nil || id <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		a.getQuestion(w, r, uint(id))
	case http.MethodDelete:
		a.deleteQuestion(w, r, uint(id))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (a *App) listQuestions(w http.ResponseWriter, r *http.Request) {
	var q []Question
	if err := a.DB.Preload("Answers").Find(&q).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	writeJSON(w, q)
}

func (a *App) createQuestion(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Text string `json:"text"`
	}
	if err := readJSON(r, &in); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(in.Text) == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	q := Question{Text: in.Text}
	if err := a.DB.Create(&q).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	writeJSON(w, q)
}

func (a *App) getQuestion(w http.ResponseWriter, r *http.Request, id uint) {
	var q Question
	if err := a.DB.Preload("Answers").First(&q, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	writeJSON(w, q)
}

func (a *App) deleteQuestion(w http.ResponseWriter, r *http.Request, id uint) {
	if err := a.DB.Delete(&Question{}, id).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

//
// ANSWERS
//

func (a *App) answersHandler(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/questions/") &&
		strings.Contains(r.URL.Path, "/answers") {

		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) >= 3 && parts[0] == "questions" && parts[2] == "answers" {
			id, _ := strconv.Atoi(parts[1])
			if r.Method == http.MethodPost {
				a.createAnswer(w, r, uint(id))
				return
			}
		}
	}

	p := strings.TrimPrefix(r.URL.Path, "/answers")
	p = strings.TrimSuffix(p, "/")
	id, err := strconv.Atoi(strings.TrimPrefix(p, "/"))
	if err != nil || id <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		a.getAnswer(w, r, uint(id))
	case http.MethodDelete:
		a.deleteAnswer(w, r, uint(id))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (a *App) createAnswer(w http.ResponseWriter, r *http.Request, questionID uint) {
	var q Question
	if err := a.DB.First(&q, questionID).Error; err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("question does not exist"))
		return
	}

	var in struct {
		UserID string `json:"user_id"`
		Text   string `json:"text"`
	}
	if err := readJSON(r, &in); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if in.UserID == "" || in.Text == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ans := Answer{
		QuestionID: questionID,
		UserID:     in.UserID,
		Text:       in.Text,
		CreatedAt:  time.Now(),
	}

	if err := a.DB.Create(&ans).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	writeJSON(w, ans)
}

func (a *App) getAnswer(w http.ResponseWriter, r *http.Request, id uint) {
	var ans Answer
	if err := a.DB.First(&ans, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	writeJSON(w, ans)
}

func (a *App) deleteAnswer(w http.ResponseWriter, r *http.Request, id uint) {
	if err := a.DB.Delete(&Answer{}, id).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
