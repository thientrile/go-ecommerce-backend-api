package middlewares

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"go-ecommerce-backend-api.com/global"
	"go-ecommerce-backend-api.com/pkg/context"
	"go-ecommerce-backend-api.com/pkg/response"
	"go.uber.org/zap"
)

const (
	// Key format constants
	keyFormatThreeParts = "%s:%s:%s"
	keyFormatTwoParts   = "%s:%s"
	keyFormatUserID     = "%s:%d:%s"
	keyFormatUserSimple = "user:%d:%s"
	keyFormatIPSimple   = "ip:%s:%s"
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
			"code":        response.ErrCodeRateLimitExceeded,
		})
		return
	}

	c.Next()
}

func (rl *RateLimiters) filterPathUrl(url string) string {
	// Duyệt qua tất cả các nhóm URL path được cấu hình
	for groupName, paths := range global.Config.Limiter.URLPath {
		for _, path := range paths {
			if strings.HasPrefix(url, path) {
				global.Logger.Debug("URL matched group",
					zap.String("url", url),
					zap.String("matched_path", path),
					zap.String("group", groupName))
				return groupName
			}
		}
	}

	// Fallback to global nếu không match với group nào
	// global.Logger.Debug("URL fallback to global", zap.String("url", url))
	return "global"
}

// DynamicRateLimiter - rate limiting động theo URL path group
func (rl *RateLimiters) DynamicRateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		urlPath := c.Request.URL.Path
		clientIP := c.ClientIP()

		// Xác định group của URL
		rateLimitGroup := rl.filterPathUrl(urlPath)
		// Tạo key dựa trên group và context
		var key string
		switch rateLimitGroup {
		case "private":
			// Sử dụng user ID cho private endpoints
			if userId, err := context.GetUserIdFormUUID(c); err == nil {
				key = fmt.Sprintf(keyFormatUserID, rateLimitGroup, userId, urlPath)
			} else {
				// Fallback to IP nếu không lấy được user ID
				key = fmt.Sprintf(keyFormatThreeParts, rateLimitGroup, clientIP, urlPath)
			}
		case "admin":
			// Admin có thể cần logic riêng
			key = fmt.Sprintf(keyFormatThreeParts, rateLimitGroup, clientIP, urlPath)
		case "upload", "payment":
			// Các endpoint đặc biệt có thể dùng IP + path
			key = fmt.Sprintf(keyFormatThreeParts, rateLimitGroup, clientIP, urlPath)
		default:
			// Public và global endpoints
			key = fmt.Sprintf(keyFormatTwoParts, rateLimitGroup, clientIP)
		}

		global.Logger.Debug("Dynamic rate limiting applied",
			zap.String("url", urlPath),
			zap.String("group", rateLimitGroup),
			zap.String("key", key),
			zap.String("ip", clientIP))

		rl.setRateLimit(c, key, rateLimitGroup)
	}
}

// SmartRateLimiter - rate limiting thông minh với priority
func (rl *RateLimiters) SmartRateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		urlPath := c.Request.URL.Path
		// Xác định group của URL với priority order
		rateLimitGroup := rl.filterPathUrlWithPriority(urlPath)

		// Check if group has custom limiter rule
		if _, exists := global.Config.Limiter.Rules[rateLimitGroup]; !exists {
			global.Logger.Warn("No rate limit rule found for group, using global",
				zap.String("group", rateLimitGroup),
				zap.String("url", urlPath))
			rateLimitGroup = "global"
		}

		// Generate smart key based on group characteristics
		key := rl.generateSmartKey(c, rateLimitGroup)

		rl.setRateLimit(c, key, rateLimitGroup)
	}
}

// filterPathUrlWithPriority filters URL with priority order (simplified)
func (rl *RateLimiters) filterPathUrlWithPriority(url string) string {
	return rl.checkPriorityGroups(url)
}

// checkPriorityGroups checks priority groups first, then others
func (rl *RateLimiters) checkPriorityGroups(url string) string {
	// Priority order: most specific to least specific
	priorityOrder := []string{"admin", "payment", "upload", "private", "public"}

	// Check priority groups first
	for _, priority := range priorityOrder {
		if rl.matchesGroup(url, priority) {
			return priority
		}
	}

	// Check remaining groups
	return rl.checkRemainingGroups(url, priorityOrder)
}

// matchesGroup checks if URL matches any path in the group
func (rl *RateLimiters) matchesGroup(url, groupName string) bool {
	if paths, exists := global.Config.Limiter.URLPath[groupName]; exists {
		for _, path := range paths {
			if strings.HasPrefix(url, path) {
				return true
			}
		}
	}
	return false
}

// checkRemainingGroups checks non-priority groups
func (rl *RateLimiters) checkRemainingGroups(url string, priorityOrder []string) string {
	for groupName := range global.Config.Limiter.URLPath {
		if rl.isPriorityGroup(groupName, priorityOrder) {
			continue
		}
		if rl.matchesGroup(url, groupName) {
			return groupName
		}
	}
	return "global"
}

// isPriorityGroup checks if group is in priority list
func (rl *RateLimiters) isPriorityGroup(groupName string, priorityOrder []string) bool {
	for _, p := range priorityOrder {
		if p == groupName {
			return true
		}
	}
	return false
}

// generateSmartKey generates appropriate key based on group type
func (rl *RateLimiters) generateSmartKey(c *gin.Context, group string) string {
	clientIP := c.ClientIP()
	urlPath := c.Request.URL.Path

	switch group {
	case "private":
		if userId, err := context.GetUserIdFormUUID(c); err == nil {
			return fmt.Sprintf("user:%d:%s", userId, group)
		}
		return fmt.Sprintf("ip:%s:%s", clientIP, group)

	case "admin":
		// Admin endpoints might need stricter per-endpoint limiting
		return fmt.Sprintf("admin:%s:%s", clientIP, urlPath)

	case "upload":
		// Upload endpoints limited per user per endpoint
		if userId, err := context.GetUserIdFormUUID(c); err == nil {
			return fmt.Sprintf("upload:%d:%s", userId, urlPath)
		}
		return fmt.Sprintf("upload:%s:%s", clientIP, urlPath)

	case "payment":
		// Payment endpoints need very strict limiting
		if userId, err := context.GetUserIdFormUUID(c); err == nil {
			return fmt.Sprintf("payment:%d", userId)
		}
		return fmt.Sprintf("payment:%s", clientIP)

	case "public":
		// Public endpoints limited per IP
		return fmt.Sprintf("public:%s", clientIP)

	default:
		// Global or unknown groups
		return fmt.Sprintf("global:%s", clientIP)
	}
}

// GetURLGroups returns all configured URL groups
func (rl *RateLimiters) GetURLGroups() []string {
	groups := make([]string, 0, len(global.Config.Limiter.URLPath))
	for group := range global.Config.Limiter.URLPath {
		groups = append(groups, group)
	}
	return groups
}

// IsURLInGroup checks if URL belongs to a specific group
func (rl *RateLimiters) IsURLInGroup(url, group string) bool {
	paths, exists := global.Config.Limiter.URLPath[group]
	if !exists {
		return false
	}

	for _, path := range paths {
		if strings.HasPrefix(url, path) {
			return true
		}
	}
	return false
}

// GetURLGroupPaths returns all paths for a specific group
func (rl *RateLimiters) GetURLGroupPaths(group string) []string {
	if paths, exists := global.Config.Limiter.URLPath[group]; exists {
		return paths
	}
	return []string{}
}

// ValidateURLPathConfig validates URL path configuration
func (rl *RateLimiters) ValidateURLPathConfig() error {
	if global.Config.Limiter.URLPath == nil {
		return fmt.Errorf("url_path configuration is nil")
	}

	for group, paths := range global.Config.Limiter.URLPath {
		if len(paths) == 0 {
			global.Logger.Warn("Empty path list for group", zap.String("group", group))
			continue
		}

		for _, path := range paths {
			if path == "" {
				return fmt.Errorf("empty path found in group '%s'", group)
			}
			if !strings.HasPrefix(path, "/") {
				return fmt.Errorf("path '%s' in group '%s' must start with '/'", path, group)
			}
		}
	}

	global.Logger.Info("URL path configuration validated",
		zap.Int("groups_count", len(global.Config.Limiter.URLPath)),
		zap.Strings("groups", rl.GetURLGroups()))

	return nil
}
