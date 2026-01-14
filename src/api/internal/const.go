package internal

import (
	"sync"
	"time"

	"github.com/algo-shield/algo-shield/src/pkg/config"
)

var (
	globalConfig *config.Config
	configMutex  sync.RWMutex
)

// SetGlobalConfig sets the global configuration for timeout access
// Should be called once during application startup
func SetGlobalConfig(cfg *config.Config) {
	configMutex.Lock()
	defer configMutex.Unlock()
	globalConfig = cfg
}

// GetHandlerTimeout returns the configured handler timeout from config
// Falls back to default if config is not set (for backward compatibility)
func GetHandlerTimeout() time.Duration {
	configMutex.RLock()
	defer configMutex.RUnlock()
	if globalConfig != nil && globalConfig.API.Timeouts.HandlerTimeout > 0 {
		return globalConfig.API.Timeouts.HandlerTimeout
	}
	return 500 * time.Millisecond // Default fallback
}

// GetHealthCheckTimeout returns the configured health check timeout from config
// Falls back to default if config is not set
func GetHealthCheckTimeout() time.Duration {
	configMutex.RLock()
	defer configMutex.RUnlock()
	if globalConfig != nil && globalConfig.API.Timeouts.HealthCheck > 0 {
		return globalConfig.API.Timeouts.HealthCheck
	}
	return 2 * time.Second // Default fallback
}

// GetCacheTTL returns the configured cache TTL for the given cache type
func GetCacheTTL(cacheType string) time.Duration {
	configMutex.RLock()
	defer configMutex.RUnlock()
	if globalConfig == nil {
		return getDefaultCacheTTL(cacheType)
	}
	switch cacheType {
	case "dashboard":
		if globalConfig.API.Cache.DashboardTTL > 0 {
			return globalConfig.API.Cache.DashboardTTL
		}
	case "branding":
		if globalConfig.API.Cache.BrandingTTL > 0 {
			return globalConfig.API.Cache.BrandingTTL
		}
	case "system":
		if globalConfig.API.Cache.SystemTTL > 0 {
			return globalConfig.API.Cache.SystemTTL
		}
	case "rules":
		if globalConfig.API.Cache.RulesTTL > 0 {
			return globalConfig.API.Cache.RulesTTL
		}
	}
	return getDefaultCacheTTL(cacheType)
}

func getDefaultCacheTTL(cacheType string) time.Duration {
	switch cacheType {
	case "dashboard":
		return 30 * time.Second
	case "branding":
		return 10 * time.Minute
	case "system":
		return 5 * time.Minute
	case "rules":
		return 5 * time.Minute
	default:
		return 5 * time.Minute
	}
}
