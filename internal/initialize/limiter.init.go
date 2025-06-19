package initialize

import (
	"fmt"
	"net"
	"time"

	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	redisStore "github.com/ulule/limiter/v3/drivers/store/redis"
	"go-ecommerce-backend-api.com/global"
	"go-ecommerce-backend-api.com/pkg/setting"
	"go.uber.org/zap"
)

// InitLimiter initializes rate limiters based on configuration
func InitLimiter() {
	global.Logger.Info("Starting rate limiter initialization")

	// Initialize the limiters map
	global.Limiters = make(map[string]*limiter.Limiter)

	// Validate configuration first
	if err := ValidateLimiterConfig(); err != nil {
		global.Logger.Fatal("Rate limiter configuration validation failed", zap.Error(err))
	}

	// Check if limiter rules exist
	if global.Config.Limiter.Rules == nil {
		global.Logger.Warn("No limiter rules configured, skipping initialization")
		return
	}

	global.Logger.Info("Initializing rate limiters",
		zap.Int("total_rules", len(global.Config.Limiter.Rules)),
		zap.String("store_type", getStoreType()))

	enabledCount := 0
	for name, rule := range global.Config.Limiter.Rules {
		if !rule.Enabled {
			global.Logger.Info("Skipping disabled limiter",
				zap.String("name", name),
				zap.String("description", rule.Description))
			continue
		}

		global.Logger.Info("Creating rate limiter",
			zap.String("name", name),
			zap.String("rate", rule.Rate),
			zap.String("description", rule.Description),
			zap.Int("burst_multiplier", rule.BurstMultiplier),
			zap.Bool("strict_mode", rule.StrictMode),
			zap.Int("ip_whitelist_count", len(rule.IPWhitelist)))

		limiterInstance, err := createRateLimiter(rule)
		if err != nil {
			global.Logger.Fatal("Failed to create rate limiter",
				zap.String("name", name),
				zap.String("rate", rule.Rate),
				zap.Error(err))
		}

		global.Limiters[name] = limiterInstance
		enabledCount++

		global.Logger.Debug("Rate limiter created successfully",
			zap.String("name", name))
	}

	global.Logger.Info("Rate limiters initialized successfully",
		zap.Int("active_limiters", enabledCount),
		zap.Int("total_rules", len(global.Config.Limiter.Rules)),
		zap.String("store_type", getStoreType()))
}

// createRateLimiter creates a single rate limiter instance based on rule configuration
func createRateLimiter(rule setting.RuleConfig) (*limiter.Limiter, error) {
	// Parse and validate rate format
	rate, err := limiter.NewRateFromFormatted(rule.Rate)
	if err != nil {
		return nil, fmt.Errorf("invalid rate format '%s': %w", rule.Rate, err)
	}

	// Apply burst multiplier if specified
	if rule.BurstMultiplier > 1 {
		newBurst := int64(rule.BurstMultiplier) * rate.Limit
		global.Logger.Debug("Applied burst multiplier",
			zap.Int64("original_limit", rate.Limit),
			zap.Int64("new_burst", newBurst),
			zap.Int("multiplier", rule.BurstMultiplier))
		// Note: The ulule/limiter library does not support a burst parameter directly in Rate,
		// so you may need to handle burst logic at a higher level or use a different library if strict burst control is required.
	}

	// Validate IP whitelist if configured
	if len(rule.IPWhitelist) > 0 {
		if err := validateIPWhitelist(rule.IPWhitelist); err != nil {
			return nil, fmt.Errorf("invalid IP whitelist: %w", err)
		}
		global.Logger.Debug("IP whitelist validated",
			zap.Strings("whitelist", rule.IPWhitelist))
	}

	// Create store options
	storeOptions, err := createStoreOptions()
	if err != nil {
		return nil, fmt.Errorf("failed to create store options: %w", err)
	}

	// Create store
	store, err := createStore(storeOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create store: %w", err)
	}

	return limiter.New(store, rate), nil
}

// createStoreOptions creates store options from configuration with safe defaults
func createStoreOptions() (limiter.StoreOptions, error) {
	// Set safe defaults
	options := limiter.StoreOptions{
		Prefix:          "rate_limit",
		MaxRetry:        3,
		CleanUpInterval: 300 * time.Second, // 5 minutes
	}

	// Override with config values if available
	if global.Config.Limiter.DefaultConfig != nil {
		config := global.Config.Limiter.DefaultConfig

		// Get prefix
		if prefix, ok := config["prefix"].(string); ok && prefix != "" {
			options.Prefix = prefix
		}

		// Get max retry
		if maxRetry, ok := config["max_retry"].(int); ok && maxRetry > 0 {
			options.MaxRetry = maxRetry
		}

		// Get cleanup interval
		if cleanupInterval, ok := config["clean_up_interval"].(int); ok && cleanupInterval > 0 {
			options.CleanUpInterval = time.Duration(cleanupInterval) * time.Second
		}
	}

	global.Logger.Debug("Store options created",
		zap.String("prefix", options.Prefix),
		zap.Int("max_retry", options.MaxRetry),
		zap.Duration("cleanup_interval", options.CleanUpInterval))

	return options, nil
}

// createStore creates either Redis or Memory store based on configuration
func createStore(options limiter.StoreOptions) (limiter.Store, error) {
	if global.Config.Limiter.Store == 1 {
		// Redis store
		if global.RDB == nil {
			return nil, fmt.Errorf("redis client not initialized - ensure InitRedis() is called before InitLimiter()")
		}

		store, err := redisStore.NewStoreWithOptions(global.RDB, options)
		if err != nil {
			return nil, fmt.Errorf("failed to create Redis store: %w", err)
		}

		global.Logger.Debug("Redis store created successfully")
		return store, nil
	} else {
		// Memory store
		store := memory.NewStoreWithOptions(options)
		global.Logger.Debug("Memory store created successfully")
		return store, nil
	}
}

// validateIPWhitelist validates IP addresses and CIDR blocks in whitelist
func validateIPWhitelist(whitelist []string) error {
	for i, ipOrCIDR := range whitelist {
		// Try parsing as CIDR first
		if _, _, err := net.ParseCIDR(ipOrCIDR); err != nil {
			// If CIDR parsing fails, try parsing as IP
			if net.ParseIP(ipOrCIDR) == nil {
				return fmt.Errorf("invalid IP or CIDR format at index %d: '%s'", i, ipOrCIDR)
			}
		}
	}
	return nil
}

// getStoreType returns human-readable store type
func getStoreType() string {
	if global.Config.Limiter.Store == 1 {
		return "redis"
	}
	return "memory"
}

// ValidateLimiterConfig validates the complete limiter configuration
func ValidateLimiterConfig() error {
	// Check if limiter config exists
	if global.Config.Limiter.Rules == nil {
		return fmt.Errorf("limiter rules not configured")
	}

	// Validate store type
	if global.Config.Limiter.Store < 0 || global.Config.Limiter.Store > 1 {
		return fmt.Errorf("invalid store type: %d (must be 0=memory or 1=redis)", global.Config.Limiter.Store)
	}

	// Validate each rule
	for name, rule := range global.Config.Limiter.Rules {
		if err := validateRule(rule); err != nil {
			return fmt.Errorf("rule '%s' validation failed: %w", name, err)
		}
	}

	global.Logger.Info("Rate limiter configuration validated successfully")
	return nil
}

// validateRule validates a single rate limiter rule
func validateRule(rule setting.RuleConfig) error {
	// Check rate format
	if rule.Rate == "" {
		return fmt.Errorf("rate cannot be empty")
	}

	// Test rate parsing
	if _, err := limiter.NewRateFromFormatted(rule.Rate); err != nil {
		return fmt.Errorf("invalid rate format '%s': %w", rule.Rate, err)
	}

	// Validate burst multiplier
	if rule.BurstMultiplier < 0 {
		return fmt.Errorf("burst_multiplier cannot be negative: %d", rule.BurstMultiplier)
	}

	// Validate IP whitelist
	if len(rule.IPWhitelist) > 0 {
		if err := validateIPWhitelist(rule.IPWhitelist); err != nil {
			return fmt.Errorf("invalid IP whitelist: %w", err)
		}
	}

	return nil
}

// GetLimiter returns a specific rate limiter by name
func GetLimiter(name string) (*limiter.Limiter, bool) {
	if global.Limiters == nil {
		return nil, false
	}

	limiterInstance, exists := global.Limiters[name]
	return limiterInstance, exists
}

// GetActiveLimiters returns all active rate limiter names
func GetActiveLimiters() []string {
	if global.Limiters == nil {
		return []string{}
	}

	names := make([]string, 0, len(global.Limiters))
	for name := range global.Limiters {
		names = append(names, name)
	}
	return names
}
