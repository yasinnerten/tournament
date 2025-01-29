package router

import (
	"net/http"
	"strconv"

	"tournament-app/model"
	"tournament-app/service"

	"github.com/gin-gonic/gin"
)

// TournamentRoutes sets up the tournament routes
func TournamentRoutes(router *gin.Engine) {
	router.POST("/tournaments", createTournament)
	router.DELETE("/tournaments/:id", deleteTournament)
	router.PUT("/tournaments/:id", updateTournament)
	router.GET("/tournaments/:id", getTournamentByID)
	router.GET("/tournaments/ongoing", getOngoingTournaments)
	router.POST("/tournaments/join", joinTournament)
	router.POST("/tournaments/:id/end", endTournament)
	router.GET("/tournaments", getAllTournaments)

	router.GET("/leaderboard", getLeaderboard)
	router.GET("/leaderboard/tournament/:id", getLeaderboardByTournamentID)
	router.GET("/leaderboard/user/:id", getLeaderboardByUserID)
	router.GET("/leaderboard/tournament/:id/finished", getFinishedLeaderboardByTournamentID)

	router.GET("/leaderboard/active", getActiveLeaderboard)
	router.GET("/leaderboard/user/:id/active", getActiveLeaderboardByUserID)
	router.GET("/leaderboard/tournament/:id/active", getActiveLeaderboardByTournamentID)
}

// @Summary Create a new tournament
// @Description Create a new tournament with the input payload
// @Tags tournaments
// @Accept  json
// @Produce  json
// @Param   tournament  body    object{name=string, prize=int}  true  "Tournament"
// @Success 201 {object} model.Tournament
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /tournaments [post]
func createTournament(c *gin.Context) {
	var tournament model.Tournament
	if err := c.ShouldBindJSON(&tournament); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.CreateTournament(&tournament); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, tournament)
}

// @Summary Delete a tournament
// @Description Delete a tournament by ID
// @Tags tournaments
// @Produce  json
// @Param   id  path  int  true  "Tournament ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /tournaments/{id} [delete]
func deleteTournament(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tournament ID"})
		return
	}

	if err := service.DeleteTournament(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Tournament deleted successfully"})
}

// @Summary Update a tournament
// @Description Update a tournament by ID with the input payload
// @Tags tournaments
// @Accept  json
// @Produce  json
// @Param   id          path    int              true  "Tournament ID"
// @Param   tournament  body    model.Tournament  true  "Tournament"
// @Success 200 {object} model.Tournament
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /tournaments/{id} [put]
func updateTournament(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tournament ID"})
		return
	}

	var tournament model.Tournament
	if err := c.ShouldBindJSON(&tournament); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tournament.ID = uint(id)
	if err := service.UpdateTournament(&tournament); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tournament)
}

// @Summary Get a tournament by ID
// @Description Get a tournament by its ID
// @Tags tournaments
// @Produce  json
// @Param   id  path  int  true  "Tournament ID"
// @Success 200 {object} model.Tournament
// @Failure 500 {object} map[string]interface{}
// @Router /tournaments/{id} [get]
func getTournamentByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	tournament, err := service.GetTournamentByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tournament)
}

// @Summary Get ongoing tournaments
// @Description Get a list of ongoing tournaments
// @Tags tournaments
// @Produce  json
// @Success 200 {array} model.Tournament
// @Failure 500 {object} map[string]interface{}
// @Router /tournaments/ongoing [get]
func getOngoingTournaments(c *gin.Context) {
	tournaments, err := service.GetOngoingTournaments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tournaments)
}

// @Summary Join a tournament
// @Description Join a tournament with the input payload
// @Tags tournaments
// @Accept  json
// @Produce  json
// @Param   joinRequest  body    object{tournament_id=uint, user_id=uint}  true  "Join Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /tournaments/join [post]
func joinTournament(c *gin.Context) {
	var request struct {
		TournamentID uint `json:"tournament_id" binding:"required"`
		UserID       uint `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := service.JoinTournament(request.TournamentID, request.UserID); err != nil {
		if err.Error() == "cannot join a finished tournament" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User joined tournament successfully"})
}

// @Summary End a tournament
// @Description End a tournament by ID
// @Tags tournaments
// @Produce  json
// @Param   id  path  int  true  "Tournament ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /tournaments/{id}/end [post]
func endTournament(c *gin.Context) {
	tournamentID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := service.EndTournament(uint(tournamentID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Tournament ended successfully"})
}

// @Summary Get leaderboard
// @Description Get the leaderboard
// @Tags leaderboard
// @Produce  json
// @Param   start  query  int  false  "Start"
// @Param   stop   query  int  false  "Stop"
// @Success 200 {array} model.User
// @Failure 500 {object} map[string]interface{}
// @Router /leaderboard [get]
func getLeaderboard(c *gin.Context) {
	start, _ := strconv.ParseInt(c.DefaultQuery("start", "0"), 10, 64)
	stop, _ := strconv.ParseInt(c.DefaultQuery("stop", "10"), 10, 64)
	leaderboard, err := service.GetActiveLeaderboard(start, stop)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, leaderboard)
}

// @Summary Get leaderboard by tournament ID
// @Description Get the leaderboard for a specific tournament
// @Tags leaderboard
// @Produce  json
// @Param   id  path  int  true  "Tournament ID"
// @Success 200 {array} model.User
// @Failure 500 {object} map[string]interface{}
// @Router /leaderboard/tournament/{id} [get]
func getLeaderboardByTournamentID(c *gin.Context) {
	tournamentID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	leaderboard, err := service.GetActiveLeaderboardByTournamentID(uint(tournamentID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, leaderboard)
}

// @Summary Get leaderboard by user ID
// @Description Get the leaderboard for a specific user
// @Tags leaderboard
// @Produce  json
// @Param   id  path  int  true  "User ID"
// @Success 200 {array} model.User
// @Failure 500 {object} map[string]interface{}
// @Router /leaderboard/user/{id} [get]
func getLeaderboardByUserID(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	leaderboard, err := service.GetActiveLeaderboardByUserID(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, leaderboard)
}

// @Summary Get finished leaderboard by tournament ID
// @Description Get the finished leaderboard for a specific tournament
// @Tags leaderboard
// @Produce  json
// @Param   id  path  int  true  "Tournament ID"
// @Success 200 {array} model.User
// @Failure 500 {object} map[string]interface{}
// @Router /leaderboard/tournament/{id}/finished [get]
func getFinishedLeaderboardByTournamentID(c *gin.Context) {
	tournamentID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	leaderboard, err := service.GetFinishedLeaderboardByTournamentID(uint(tournamentID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, leaderboard)
}

// @Summary Get all tournaments
// @Description Get a list of all tournaments
// @Tags tournaments
// @Produce  json
// @Success 200 {array} model.Tournament
// @Failure 500 {object} map[string]interface{}
// @Router /tournaments [get]
func getAllTournaments(c *gin.Context) {
	tournaments, err := service.GetAllTournaments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tournaments)
}

// @Summary Get active leaderboard
// @Description Get the active leaderboard
// @Tags leaderboard
// @Produce  json
// @Success 200 {array} model.User
// @Failure 500 {object} map[string]interface{}
// @Router /leaderboard/active [get]
func getActiveLeaderboard(c *gin.Context) {
	start, _ := strconv.ParseInt(c.DefaultQuery("start", "0"), 10, 64)
	stop, _ := strconv.ParseInt(c.DefaultQuery("stop", "10"), 10, 64)
	leaderboard, err := service.GetActiveLeaderboard(start, stop)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, leaderboard)
}

// @Summary Get active leaderboard by user ID
// @Description Get the active leaderboard for a specific user
// @Tags leaderboard
// @Produce  json
// @Param   id  path  int  true  "User ID"
// @Success 200 {array} model.User
// @Failure 500 {object} map[string]interface{}
// @Router /leaderboard/user/{id}/active [get]
func getActiveLeaderboardByUserID(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	leaderboard, err := service.GetActiveLeaderboardByUserID(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, leaderboard)
}

// @Summary Get active leaderboard by tournament ID
// @Description Get the active leaderboard for a specific tournament
// @Tags leaderboard
// @Produce  json
// @Param   id  path  int  true  "Tournament ID"
// @Success 200 {array} model.User
// @Failure 500 {object} map[string]interface{}
// @Router /leaderboard/tournament/{id}/active [get]
func getActiveLeaderboardByTournamentID(c *gin.Context) {
	tournamentID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	leaderboard, err := service.GetActiveLeaderboardByTournamentID(uint(tournamentID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, leaderboard)
}
