package dashboard

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/algo-shield/algo-shield/src/api/internal/shared"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_Service_GetMetrics_WhenCacheHit_ThenReturnsCachedMetrics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedMetrics := &DashboardMetrics{
		StatusDistribution: []StatusCount{{Status: "approved", Count: 10}},
		Temporal24h:        []TemporalCount{{Bucket: time.Now(), Count: 5}},
		Temporal7d:         []TemporalCount{{Bucket: time.Now(), Count: 50}},
		Temporal30d:        []TemporalCount{{Bucket: time.Now(), Count: 200}},
		TotalCount:         100,
		CachedAt:           time.Now(),
	}
	cachedData, _ := json.Marshal(expectedMetrics)
	mockRepo := NewMockRepository(ctrl)
	mockRedis := NewMockRedisClient(ctrl)
	mockRedis.EXPECT().Get(gomock.Any(), "dashboard:metrics").Return(redis.NewStringResult(string(cachedData), nil))
	service := NewService(mockRepo, mockRedis)

	result, err := service.GetMetrics(context.Background())

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedMetrics.TotalCount, result.TotalCount)
	assert.Len(t, result.StatusDistribution, 1)
}

func Test_Service_GetMetrics_WhenCacheMiss_ThenFetchesFromDatabase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	statusDist := []StatusCount{{Status: "approved", Count: 10}, {Status: "rejected", Count: 5}}
	temporal24h := []TemporalCount{{Bucket: time.Now(), Count: 5}}
	temporal7d := []TemporalCount{{Bucket: time.Now(), Count: 50}}
	temporal30d := []TemporalCount{{Bucket: time.Now(), Count: 200}}
	totalCount := int64(100)
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().GetStatusDistribution(gomock.Any()).Return(statusDist, nil)
	mockRepo.EXPECT().GetTemporalCounts(gomock.Any(), "24h").Return(temporal24h, nil)
	mockRepo.EXPECT().GetTemporalCounts(gomock.Any(), "7d").Return(temporal7d, nil)
	mockRepo.EXPECT().GetTemporalCounts(gomock.Any(), "30d").Return(temporal30d, nil)
	mockRepo.EXPECT().GetTotalCount(gomock.Any()).Return(totalCount, nil)
	mockRedis := NewMockRedisClient(ctrl)
	mockRedis.EXPECT().Get(gomock.Any(), "dashboard:metrics").Return(redis.NewStringResult("", redis.Nil))
	mockRedis.EXPECT().Set(gomock.Any(), "dashboard:metrics", gomock.Any(), gomock.Any()).Return(redis.NewStatusResult("OK", nil))
	service := NewService(mockRepo, mockRedis)

	result, err := service.GetMetrics(context.Background())

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, totalCount, result.TotalCount)
	assert.Len(t, result.StatusDistribution, 2)
	assert.Len(t, result.Temporal24h, 1)
	assert.Len(t, result.Temporal7d, 1)
	assert.Len(t, result.Temporal30d, 1)
}

func Test_Service_GetMetrics_WhenRepositoryMethodFails_ThenReturnsError(t *testing.T) {
	statusDist := []StatusCount{{Status: "approved", Count: 10}}
	temporal24h := []TemporalCount{{Bucket: time.Now(), Count: 5}}
	temporal7d := []TemporalCount{{Bucket: time.Now(), Count: 50}}
	temporal30d := []TemporalCount{{Bucket: time.Now(), Count: 200}}

	testCases := []struct {
		name      string
		setupMock func(*MockRepository)
	}{
		{
			name: "status distribution fails",
			setupMock: func(m *MockRepository) {
				m.EXPECT().GetStatusDistribution(gomock.Any()).Return(nil, errors.New("database error"))
			},
		},
		{
			name: "temporal 24h fails",
			setupMock: func(m *MockRepository) {
				m.EXPECT().GetStatusDistribution(gomock.Any()).Return(statusDist, nil)
				m.EXPECT().GetTemporalCounts(gomock.Any(), "24h").Return(nil, errors.New("temporal error"))
			},
		},
		{
			name: "temporal 7d fails",
			setupMock: func(m *MockRepository) {
				m.EXPECT().GetStatusDistribution(gomock.Any()).Return(statusDist, nil)
				m.EXPECT().GetTemporalCounts(gomock.Any(), "24h").Return(temporal24h, nil)
				m.EXPECT().GetTemporalCounts(gomock.Any(), "7d").Return(nil, errors.New("7d error"))
			},
		},
		{
			name: "temporal 30d fails",
			setupMock: func(m *MockRepository) {
				m.EXPECT().GetStatusDistribution(gomock.Any()).Return(statusDist, nil)
				m.EXPECT().GetTemporalCounts(gomock.Any(), "24h").Return(temporal24h, nil)
				m.EXPECT().GetTemporalCounts(gomock.Any(), "7d").Return(temporal7d, nil)
				m.EXPECT().GetTemporalCounts(gomock.Any(), "30d").Return(nil, errors.New("30d error"))
			},
		},
		{
			name: "total count fails",
			setupMock: func(m *MockRepository) {
				m.EXPECT().GetStatusDistribution(gomock.Any()).Return(statusDist, nil)
				m.EXPECT().GetTemporalCounts(gomock.Any(), "24h").Return(temporal24h, nil)
				m.EXPECT().GetTemporalCounts(gomock.Any(), "7d").Return(temporal7d, nil)
				m.EXPECT().GetTemporalCounts(gomock.Any(), "30d").Return(temporal30d, nil)
				m.EXPECT().GetTotalCount(gomock.Any()).Return(int64(0), errors.New("count error"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := NewMockRepository(ctrl)
			tc.setupMock(mockRepo)
			mockRedis := NewMockRedisClient(ctrl)
			mockRedis.EXPECT().Get(gomock.Any(), "dashboard:metrics").Return(redis.NewStringResult("", redis.Nil))
			service := NewService(mockRepo, mockRedis)

			result, err := service.GetMetrics(context.Background())

			assert.Nil(t, result)
			assert.Error(t, err)
		})
	}
}

func Test_Service_InvalidateCache_WhenCalled_ThenDeletesBothCacheKeys(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepository(ctrl)
	mockRedis := NewMockRedisClient(ctrl)
	mockRedis.EXPECT().Del(gomock.Any(), "dashboard:metrics").Return(redis.NewIntResult(1, nil))
	mockRedis.EXPECT().Del(gomock.Any(), "dashboard:metrics:synthetic").Return(redis.NewIntResult(1, nil))
	service := NewService(mockRepo, mockRedis)

	err := service.InvalidateCache(context.Background())

	assert.NoError(t, err)
}

func Test_Service_GetMetrics_WhenSyntheticMode_ThenUsesSyntheticCacheKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	statusDist := []StatusCount{{Status: "pending", Count: 10}}
	temporal24h := []TemporalCount{{Bucket: time.Now(), Count: 5}}
	temporal7d := []TemporalCount{{Bucket: time.Now(), Count: 50}}
	temporal30d := []TemporalCount{{Bucket: time.Now(), Count: 200}}
	totalCount := int64(100)
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().GetStatusDistribution(gomock.Any()).Return(statusDist, nil)
	mockRepo.EXPECT().GetTemporalCounts(gomock.Any(), "24h").Return(temporal24h, nil)
	mockRepo.EXPECT().GetTemporalCounts(gomock.Any(), "7d").Return(temporal7d, nil)
	mockRepo.EXPECT().GetTemporalCounts(gomock.Any(), "30d").Return(temporal30d, nil)
	mockRepo.EXPECT().GetTotalCount(gomock.Any()).Return(totalCount, nil)
	mockRedis := NewMockRedisClient(ctrl)
	mockRedis.EXPECT().Get(gomock.Any(), "dashboard:metrics:synthetic").Return(redis.NewStringResult("", redis.Nil))
	mockRedis.EXPECT().Set(gomock.Any(), "dashboard:metrics:synthetic", gomock.Any(), gomock.Any()).Return(redis.NewStatusResult("OK", nil))
	service := NewService(mockRepo, mockRedis)
	ctx := shared.WithSyntheticMode(context.Background(), true)

	result, err := service.GetMetrics(ctx)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, totalCount, result.TotalCount)
}

func Test_Service_GetMetrics_WhenCacheSetFails_ThenStillReturnsMetrics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	statusDist := []StatusCount{{Status: "approved", Count: 10}}
	temporal24h := []TemporalCount{{Bucket: time.Now(), Count: 5}}
	temporal7d := []TemporalCount{{Bucket: time.Now(), Count: 50}}
	temporal30d := []TemporalCount{{Bucket: time.Now(), Count: 200}}
	totalCount := int64(100)
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().GetStatusDistribution(gomock.Any()).Return(statusDist, nil)
	mockRepo.EXPECT().GetTemporalCounts(gomock.Any(), "24h").Return(temporal24h, nil)
	mockRepo.EXPECT().GetTemporalCounts(gomock.Any(), "7d").Return(temporal7d, nil)
	mockRepo.EXPECT().GetTemporalCounts(gomock.Any(), "30d").Return(temporal30d, nil)
	mockRepo.EXPECT().GetTotalCount(gomock.Any()).Return(totalCount, nil)
	mockRedis := NewMockRedisClient(ctrl)
	mockRedis.EXPECT().Get(gomock.Any(), "dashboard:metrics").Return(redis.NewStringResult("", redis.Nil))
	mockRedis.EXPECT().Set(gomock.Any(), "dashboard:metrics", gomock.Any(), gomock.Any()).Return(redis.NewStatusResult("", errors.New("redis error")))
	service := NewService(mockRepo, mockRedis)

	result, err := service.GetMetrics(context.Background())

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, totalCount, result.TotalCount)
}
