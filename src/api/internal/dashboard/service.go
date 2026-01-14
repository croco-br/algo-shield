package dashboard

import (
	"context"
	"encoding/json"
	"time"

	"github.com/algo-shield/algo-shield/src/api/internal"
	"github.com/algo-shield/algo-shield/src/api/internal/shared"
	"github.com/redis/go-redis/v9"
)

// DashboardMetrics contains all dashboard data
type DashboardMetrics struct {
	StatusDistribution []StatusCount   `json:"status_distribution"`
	Temporal24h        []TemporalCount `json:"temporal_24h"`
	Temporal7d         []TemporalCount `json:"temporal_7d"`
	Temporal30d        []TemporalCount `json:"temporal_30d"`
	TotalCount         int64           `json:"total_count"`
	CachedAt           time.Time       `json:"cached_at"`
}

// Service defines the interface for dashboard business logic
type Service interface {
	GetMetrics(ctx context.Context) (*DashboardMetrics, error)
	InvalidateCache(ctx context.Context) error
}

// RedisClient defines the Redis interface we need
type RedisClient interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
}

type service struct {
	repo  Repository
	redis RedisClient
}

func NewService(repo Repository, redis RedisClient) Service {
	return &service{
		repo:  repo,
		redis: redis,
	}
}

func (s *service) GetMetrics(ctx context.Context) (*DashboardMetrics, error) {
	// Try cache first
	cached, err := s.getFromCache(ctx)
	if err == nil && cached != nil {
		return cached, nil
	}

	// Fetch from database
	metrics := &DashboardMetrics{
		CachedAt: time.Now(),
	}

	statusDist, err := s.repo.GetStatusDistribution(ctx)
	if err != nil {
		return nil, err
	}
	metrics.StatusDistribution = statusDist

	temporal24h, err := s.repo.GetTemporalCounts(ctx, "24h")
	if err != nil {
		return nil, err
	}
	metrics.Temporal24h = temporal24h

	temporal7d, err := s.repo.GetTemporalCounts(ctx, "7d")
	if err != nil {
		return nil, err
	}
	metrics.Temporal7d = temporal7d

	temporal30d, err := s.repo.GetTemporalCounts(ctx, "30d")
	if err != nil {
		return nil, err
	}
	metrics.Temporal30d = temporal30d

	totalCount, err := s.repo.GetTotalCount(ctx)
	if err != nil {
		return nil, err
	}
	metrics.TotalCount = totalCount

	// Cache the result
	_ = s.setCache(ctx, metrics)

	return metrics, nil
}

func (s *service) getCacheKey(ctx context.Context) string {
	if shared.IsSyntheticMode(ctx) {
		return "dashboard:metrics:synthetic"
	}
	return "dashboard:metrics"
}

func (s *service) InvalidateCache(ctx context.Context) error {
	// Invalidate both caches
	s.redis.Del(ctx, "dashboard:metrics")
	s.redis.Del(ctx, "dashboard:metrics:synthetic")
	return nil
}

func (s *service) getFromCache(ctx context.Context) (*DashboardMetrics, error) {
	cacheKey := s.getCacheKey(ctx)
	data, err := s.redis.Get(ctx, cacheKey).Bytes()
	if err != nil {
		return nil, err
	}

	var metrics DashboardMetrics
	if err := json.Unmarshal(data, &metrics); err != nil {
		return nil, err
	}

	return &metrics, nil
}

func (s *service) setCache(ctx context.Context, metrics *DashboardMetrics) error {
	cacheKey := s.getCacheKey(ctx)
	data, err := json.Marshal(metrics)
	if err != nil {
		return err
	}

	return s.redis.Set(ctx, cacheKey, data, internal.GetCacheTTL("dashboard")).Err()
}
