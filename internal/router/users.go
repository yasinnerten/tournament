package router

import (
	"net/http"
	"strconv"

	"tournament-app/model"
	"tournament-app/service"

	"github.com/gin-gonic/gin"
)

// UserRoutes sets up the user routes
func UserRoutes(router *gin.Engine) {
	router.POST("/users", createUser)
	router.DELETE("/users/:id", deleteUser)
	router.PUT("/users/:id", updateUser)
	router.GET("/users/:id", getUserByID)
	router.GET("/users", getUsers)
	router.POST("/users/:id/levelup", levelUpUser)
	router.GET("/health", getHealth)
	router.POST("/clear-database", clearDatabase)
}

// @Summary Create a new user
// @Description Create a new user with the input payload
// @Tags users
// @Accept  json
// @Produce  json
// @Param   user  body    object{name=string, money=integer, level=integer}  true  "User"
// @Success 201 {object} model.User
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users [post]
func createUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}
	if err := service.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

// @Summary Delete a user
// @Description Delete a user by ID
// @Tags users
// @Produce  json
// @Param   id  path  integer  true  "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users/{id} [delete]
func deleteUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid user ID"})
		return
	}

	if err := service.DeleteUser(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{"message": "User deleted successfully"})
}

// @Summary Update a user
// @Description Update a user by ID with the input payload
// @Tags users
// @Accept  json
// @Produce  json
// @Param   id    path    integer        true  "User ID"
// @Param   user  body    object{name=string, money=integer, level=integer}  true  "User"
// @Success 200 {object} model.User
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users/{id} [put]
func updateUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid user ID"})
		return
	}

	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}
	user.ID = uint(id)
	if err := service.UpdateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// @Summary Get a user by ID
// @Description Get a user by their ID
// @Tags users
// @Produce  json
// @Param   id  path  integer  true  "User ID"
// @Success 200 {object} model.User
// @Failure 500 {object} map[string]interface{}
// @Router /users/{id} [get]
func getUserByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	user, err := service.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// @Summary Get all users
// @Description Get a list of all users
// @Tags users
// @Produce  json
// @Success 200 {array} model.User
// @Failure 500 {object} map[string]interface{}
// @Router /users [get]
func getUsers(c *gin.Context) {
	users, err := service.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// @Summary Level up a user
// @Description Level up a user by ID
// @Tags users
// @Produce  json
// @Param   id  path  integer  true  "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users/{id}/levelup [post]
func levelUpUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid user ID"})
		return
	}

	if err := service.LevelUpUser(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{"message": "User leveled up successfully"})
}

// @Summary Get health status
// @Description Check the health status of the service
// @Tags health
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /health [get]
func getHealth(c *gin.Context) {
	message, err := service.PerformHealthCheck()
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{"message": message})
}

// @Summary Clear the database
// @Description Clear all data from the database
// @Tags health
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /clear-database [post]
func clearDatabase(c *gin.Context) {
	if err := service.ClearDatabase(); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{"message": "Database cleared successfully"})
}
