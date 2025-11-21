package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/go-qa-api/internal/app"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestApp(t *testing.T) *app.App {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	db.AutoMigrate(&app.Question{}, &app.Answer{})
	return &app.App{DB: db}
}

func TestCreateAndListQuestions(t *testing.T) {
	a := setupTestApp(t)

	payload := map[string]string{"text": "Test question"}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/questions/", bytes.NewReader(body))
	res := httptest.NewRecorder()
	a.Routes().ServeHTTP(res, req)

	if res.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", res.Code)
	}

	req2 := httptest.NewRequest(http.MethodGet, "/questions/", nil)
	res2 := httptest.NewRecorder()
	a.Routes().ServeHTTP(res2, req2)

	if res2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", res2.Code)
	}
}
