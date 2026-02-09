package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/datmedevil17/go-vuln/internal/config"
	"github.com/datmedevil17/go-vuln/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// RateLimitMiddleware implements token bucket rate limiting using Redis
func RateLimitMiddleware(cfg *config.Config, redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !cfg.RateLimit.Enabled {
			c.Next()
			return
		}

		// Get user ID from context or use IP as fallback
		var identifier string
		if userID, exists := c.Get("user_id"); exists {
			identifier = fmt.Sprintf("user:%s", userID.(uuid.UUID).String())
		} else {
			identifier = fmt.Sprintf("ip:%s", c.ClientIP())
		}

		// Check rate limit
		allowed, err := checkRateLimit(c.Request.Context(), redisClient, identifier, cfg.RateLimit.Requests, cfg.RateLimit.Window)
		if err != nil {
			// On error, allow the request but log the error
			c.Next()
			return
		}

		if !allowed {
			utils.ErrorResponse(c, 429, "Rate limit exceeded. Please try again later.")
			c.Abort()
			return
		}

		c.Next()
	}
}

// checkRateLimit implements token bucket algorithm using Redis
func checkRateLimit(ctx context.Context, redisClient *redis.Client, identifier string, maxRequests int, window time.Duration) (bool, error) {
	key := fmt.Sprintf("ratelimit:%s", identifier)
	now := time.Now().Unix()
	windowStart := now - int64(window.Seconds())

	// Use Redis pipeline for atomic operations
	pipe := redisClient.Pipeline()

	// Remove old entries outside the window
	pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart))

	// Count requests in current window
	countCmd := pipe.ZCard(ctx, key)

	// Execute pipeline
	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return false, err
	}

	count := countCmd.Val()

	// Check if limit exceeded
	if int(count) >= maxRequests {
		return false, nil
	}

	// Add current request
	score := float64(now)
	member := fmt.Sprintf("%d:%s", now, uuid.New().String())

	pipe2 := redisClient.Pipeline()
	pipe2.ZAdd(ctx, key, redis.Z{Score: score, Member: member})
	pipe2.Expire(ctx, key, window)

	_, err = pipe2.Exec(ctx)
	if err != nil {
		return false, err
	}

	return true, nil
}

// GetRateLimitInfo returns current rate limit status for a user
func GetRateLimitInfo(ctx context.Context, redisClient *redis.Client, identifier string, maxRequests int, window time.Duration) (remaining int, resetAt time.Time, err error) {
	key := fmt.Sprintf("ratelimit:%s", identifier)
	now := time.Now()
	windowStart := now.Add(-window).Unix()

	// Count requests in current window
	count, err := redisClient.ZCount(ctx, key, fmt.Sprintf("%d", windowStart), "+inf").Result()
	if err != nil && err != redis.Nil {
		return 0, time.Time{}, err
	}

	remaining = maxRequests - int(count)
	if remaining < 0 {
		remaining = 0
	}

	// Calculate reset time (end of current window)
	resetAt = now.Add(window)

	return remaining, resetAt, nil
}
