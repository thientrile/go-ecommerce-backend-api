package middlewares

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"go-ecommerce-backend-api.com/global"
	"go-ecommerce-backend-api.com/pkg/context"
	"go.uber.org/zap"
)

type RateLimiters struct{}

// isIPWhitelisted checks if the given IP is in the whitelist.
func (rl *RateLimiters) isIPWhitelisted(ip string, whitelist []string) bool {
	for _, whitelistedIP := range whitelist {
		if ip == whitelistedIP {
			return true
		}
	}
	return false
}

func NewRateLimiter() *RateLimiters {
	return &RateLimiters{}
}

func (rl *RateLimiters) setRateLimit(c *gin.Context, key, limiterName string) {
	rateLimiter, exists := global.Limiters[limiterName]
	if !exists {
		global.Logger.Error("Limiter not found", zap.String("limiter", limiterName))
		c.AbortWithStatusJSON(500, gin.H{"error": "Limiter not configured"})
		return
	}

	// Check IP whitelist if configured
	rule := global.Config.Limiter.Rules[limiterName]
	if len(rule.IPWhitelist) > 0 && !rl.isIPWhitelisted(c.ClientIP(), rule.IPWhitelist) {
		// Apply rate limiting
	} else if len(rule.IPWhitelist) > 0 {
		// Skip rate limiting for whitelisted IPs
		c.Next()
		return
	}

	limitContext, err := rateLimiter.Get(c.Request.Context(), key)
	if err != nil {
		global.Logger.Error("Rate limit check failed", zap.Error(err))
		if rule.StrictMode {
			c.AbortWithStatusJSON(500, gin.H{"error": "Rate limiter error"})
			return
		}
		// In non-strict mode, allow the request
		c.Next()
		return
	}

	// Set enhanced headers
	c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limitContext.Limit))
	c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", limitContext.Remaining))
	c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", limitContext.Reset))
	c.Header("X-RateLimit-Policy", rule.Description)

	if limitContext.Reached {
		global.Logger.Warn("Rate limit exceeded",
			zap.String("limiter", limiterName),
			zap.String("key", key),
			zap.String("ip", c.ClientIP()))

		c.AbortWithStatusJSON(429, gin.H{
			"error":       "Rate limit exceeded",
			"message":     rule.Description,
			"retry_after": limitContext.Reset,
			"policy":      rule.Rate,
		})
		return
	}

	c.Next()
}

func (rl *RateLimiters) filterPathUrl(url string) string {
	for _, path := range global.Config.Limiter.URLPath.Public {
		if strings.HasPrefix(url, path) {
			return "public"
		}
	}

	for _, path := range global.Config.Limiter.URLPath.Private {
		if strings.HasPrefix(url, path) {
			return "private"
		}
	}

	return "global"
}

// GlobalRateLimiter - rate limiting toàn cục theo IP
func (rl *RateLimiters) GlobalRateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		key := fmt.Sprintf("global:%s", clientIP)
		rl.setRateLimit(c, key, "global")
	}
}

// PublicRateLimiter - rate limiting cho public API theo IP + URL path
func (rl *RateLimiters) PublicRateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		urlPath := c.Request.URL.Path
		rateLimitPath := rl.filterPathUrl(urlPath)
		key := fmt.Sprintf("%s:%s", clientIP, urlPath)
		rl.setRateLimit(c, key, rateLimitPath)
	}
}

// UserPrivateRateLimiter - rate limiting cho private API theo user ID + URL path
func (rl *RateLimiters) UserPrivateRateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		urlPath := c.Request.URL.Path

		userId, err := context.GetUserIdFormUUID(c)
		if err != nil {
			global.Logger.Error("Error getting user ID from context", zap.Error(err))
			// Fallback to IP-based rate limiting
			clientIP := c.ClientIP()
			key := fmt.Sprintf("%s:%s", clientIP, urlPath)
			rateLimitPath := rl.filterPathUrl(urlPath)
			rl.setRateLimit(c, key, rateLimitPath)
			return
		}

		rateLimitPath := rl.filterPathUrl(urlPath)
		key := fmt.Sprintf("%d:%s", userId, urlPath)
		rl.setRateLimit(c, key, rateLimitPath)
	}
}
