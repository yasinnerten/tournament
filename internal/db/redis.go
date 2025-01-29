package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
	"tournament-app/model"

	"github.com/go-redis/redis/v8"
)

var (
	rdb  *redis.Client
	ctx  = context.Background()
	once sync.Once
)

// InitRedis initializes the Redis client
func InitRedis(db int) {
	once.Do(func() {
		addr := os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT")
		rdb = redis.NewClient(&redis.Options{
			Addr: addr,
			DB:   db,
		})

		_, err := rdb.Ping(ctx).Result()
		if err != nil {
			log.Fatalf("Failed to connect to Redis: %v", err)
		}
	})
}

// PingRedis pings the Redis server to check the connection
func PingRedis() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	return err
}

// Redis related functions

func CreateLeaderboardEntry(entry *model.Leaderboard) error {
	ctx := context.Background()
	pipe := rdb.TxPipeline()

	score := float64(entry.Score)
	pipe.ZAdd(ctx, "leaderboard", &redis.Z{Score: score, Member: entry.UserID})

	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func GetLeaderboard(start, stop int64) ([]model.Leaderboard, error) {
	results, err := rdb.ZRevRangeWithScores(context.Background(), "leaderboard", start, stop).Result()
	if err != nil {
		return nil, err
	}

	var leaderboard []model.Leaderboard
	for _, result := range results {
		userIDStr, ok := result.Member.(string)
		if !ok {
			return nil, fmt.Errorf("invalid user ID type")
		}
		userID, err := strconv.ParseUint(userIDStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid user ID format")
		}
		leaderboard = append(leaderboard, model.Leaderboard{
			UserID: uint(userID),
			Score:  result.Score,
		})
	}
	return leaderboard, nil
}

func UpdateLeaderboard(userID string, score float64) error {
	ctx := context.Background()
	pipe := rdb.TxPipeline()

	pipe.ZAdd(ctx, "leaderboard", &redis.Z{Score: score, Member: userID})

	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func RemoveLeaderboardFromRedis(tournamentID uint) error {
	ctx := context.Background()
	pipe := rdb.TxPipeline()

	pipe.Del(ctx, fmt.Sprintf("leaderboard:%d", tournamentID))

	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
