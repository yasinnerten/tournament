package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"tournament-app/internal/router"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	router.UserRoutes(r)
	router.TournamentRoutes(r)

	return r
}

func TestRoutes(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		method   string
		endpoint string
		body     string
	}{
		// User routes
		{"POST", "/users", `{"name": "User1", "money": 1000, "level": 1}`},
		{"POST", "/users", `{"name": "User2", "money": 2000, "level": 2}`},
		{"POST", "/users", `{"name": "User3", "money": 3000, "level": 3}`},
		{"POST", "/users", `{"name": "User4", "money": 4000, "level": 4}`},
		{"POST", "/users", `{"name": "User5", "money": 5000, "level": 5}`},

		// Tournament routes
		{"POST", "/tournaments", `{"name": "Tournament1", "prize": 1000}`},
		{"POST", "/tournaments", `{"name": "Tournament2", "prize": 2000}`},
		{"POST", "/tournaments", `{"name": "Tournament3", "prize": 3000}`},
		{"DELETE", "/tournaments/1", ""},
		{"PUT", "/tournaments/2", `{"name": "Updated Tournament2", "prize": 2500}`},
		{"GET", "/tournaments/2", ""},
		{"GET", "/tournaments/ongoing", ""},
		{"POST", "/tournaments/join", `{"tournament_id": 2, "user_id": 1}`},
		{"POST", "/tournaments/2/end", ""},
		{"GET", "/tournaments", ""},

		// Leaderboard routes
		{"GET", "/leaderboard", ""},
		{"GET", "/leaderboard/tournament/2", ""},
		{"GET", "/leaderboard/user/1", ""},
		{"GET", "/leaderboard/tournament/2/finished", ""},
		{"GET", "/leaderboard/active", ""},
		{"GET", "/leaderboard/user/1/active", ""},
		{"GET", "/leaderboard/tournament/2/active", ""},
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.endpoint, func(t *testing.T) {
			var req *http.Request
			if tt.method == "POST" || tt.method == "PUT" {
				req = httptest.NewRequest(tt.method, tt.endpoint, bytes.NewBufferString(tt.body))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = httptest.NewRequest(tt.method, tt.endpoint, nil)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.NotEqual(t, http.StatusNotFound, w.Code)
		})
	}
}
