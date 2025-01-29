package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
	"golang.org/x/net/context"
)

var (
	db  *sql.DB
	rdb *redis.Client
	ctx = context.Background()
)

func initDB() {
	var err error
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        name TEXT NOT NULL,
        money INT NOT NULL
    )`)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}
}

func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
	})
}

func main() {
	initDB()
	initRedis()

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hey!",
		})
	})

	r.POST("/user", func(c *gin.Context) {
		var user struct {
			Name  string `json:"name"`
			Money int    `json:"money"`
		}
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := db.Exec("INSERT INTO users (name, money) VALUES ($1, $2)", user.Name, user.Money)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err = rdb.ZAdd(ctx, "leaderboard", &redis.Z{
			Score:  float64(user.Money),
			Member: user.Name,
		}).Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "user added"})
	})

	r.PUT("/user", func(c *gin.Context) {
		var user struct {
			Name  string `json:"name"`
			Money int    `json:"money"`
		}
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		tx, err := db.Begin()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to begin transaction"})
			return
		}

		_, err = tx.Exec("UPDATE users SET money = $1 WHERE name = $2", user.Money, user.Name)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user in database"})
			return
		}

		err = rdb.ZAdd(ctx, "leaderboard", &redis.Z{
			Score:  float64(user.Money),
			Member: user.Name,
		}).Err()
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user in leaderboard"})
			return
		}

		err = tx.Commit()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "user updated"})
	})

	r.DELETE("/user/:name", func(c *gin.Context) {
		name := c.Param("name")

		_, err := db.Exec("DELETE FROM users WHERE name = $1", name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err = rdb.ZRem(ctx, "leaderboard", name).Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "user deleted"})
	})

	r.GET("/leaderboard", func(c *gin.Context) {
		leaders, err := rdb.ZRevRangeWithScores(ctx, "leaderboard", 0, 9).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var leaderboard []gin.H
		for _, leader := range leaders {
			leaderboard = append(leaderboard, gin.H{
				"name":  leader.Member,
				"money": leader.Score,
			})
		}

		c.JSON(http.StatusOK, leaderboard)
	})

	r.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")

		score, err := rdb.ZScore(ctx, "leaderboard", name).Result()
		if err == redis.Nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"name":  name,
			"money": score,
		})
	})

	r.Run(":8080")
}
