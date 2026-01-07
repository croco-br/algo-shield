package queue

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_QueueService_PopTransaction_WhenEventAvailable_ThenReturnsEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	event := models.Event{
		"external_id": "ext-123",
		"amount":      100.50,
		"currency":    "USD",
	}
	eventJSON, _ := json.Marshal(event)
	mockRedis := NewMockRedisPopper(ctrl)
	cmd := redis.NewStringSliceCmd(context.Background())
	cmd.SetVal([]string{"transaction:queue", string(eventJSON)})
	mockRedis.EXPECT().BRPop(gomock.Any(), gomock.Any(), "transaction:queue").Return(cmd)
	service := NewQueueService(mockRedis, 5*time.Second)

	result, err := service.PopTransaction(context.Background())

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "ext-123", (*result)["external_id"])
}

func Test_QueueService_PopTransaction_WhenTimeout_ThenReturnsErrTimeout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := NewMockRedisPopper(ctrl)
	cmd := redis.NewStringSliceCmd(context.Background())
	cmd.SetErr(redis.Nil)
	mockRedis.EXPECT().BRPop(gomock.Any(), gomock.Any(), "transaction:queue").Return(cmd)
	service := NewQueueService(mockRedis, 5*time.Second)

	result, err := service.PopTransaction(context.Background())

	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrTimeout)
}

func Test_QueueService_PopTransaction_WhenRedisError_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := NewMockRedisPopper(ctrl)
	cmd := redis.NewStringSliceCmd(context.Background())
	cmd.SetErr(errors.New("redis connection error"))
	mockRedis.EXPECT().BRPop(gomock.Any(), gomock.Any(), "transaction:queue").Return(cmd)
	service := NewQueueService(mockRedis, 5*time.Second)

	result, err := service.PopTransaction(context.Background())

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.NotErrorIs(t, err, ErrTimeout)
}

func Test_QueueService_PopTransaction_WhenInvalidJSON_ThenReturnsErrInvalidData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := NewMockRedisPopper(ctrl)
	cmd := redis.NewStringSliceCmd(context.Background())
	cmd.SetVal([]string{"transaction:queue", "invalid json"})
	mockRedis.EXPECT().BRPop(gomock.Any(), gomock.Any(), "transaction:queue").Return(cmd)
	service := NewQueueService(mockRedis, 5*time.Second)

	result, err := service.PopTransaction(context.Background())

	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrInvalidData)
}

func Test_QueueService_PopTransaction_WhenEmptyResult_ThenReturnsErrInvalidData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := NewMockRedisPopper(ctrl)
	cmd := redis.NewStringSliceCmd(context.Background())
	cmd.SetVal([]string{})
	mockRedis.EXPECT().BRPop(gomock.Any(), gomock.Any(), "transaction:queue").Return(cmd)
	service := NewQueueService(mockRedis, 5*time.Second)

	result, err := service.PopTransaction(context.Background())

	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrInvalidData)
}

func Test_QueueService_PopTransaction_WhenSingleElement_ThenReturnsErrInvalidData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := NewMockRedisPopper(ctrl)
	cmd := redis.NewStringSliceCmd(context.Background())
	cmd.SetVal([]string{"transaction:queue"})
	mockRedis.EXPECT().BRPop(gomock.Any(), gomock.Any(), "transaction:queue").Return(cmd)
	service := NewQueueService(mockRedis, 5*time.Second)

	result, err := service.PopTransaction(context.Background())

	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrInvalidData)
}

func Test_QueueService_PopTransaction_WhenValidComplexEvent_ThenReturnsEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	event := models.Event{
		"external_id": "ext-123",
		"amount":      100.50,
		"currency":    "USD",
		"metadata": map[string]any{
			"source":  "api",
			"version": 2,
		},
	}
	eventJSON, _ := json.Marshal(event)
	mockRedis := NewMockRedisPopper(ctrl)
	cmd := redis.NewStringSliceCmd(context.Background())
	cmd.SetVal([]string{"transaction:queue", string(eventJSON)})
	mockRedis.EXPECT().BRPop(gomock.Any(), gomock.Any(), "transaction:queue").Return(cmd)
	service := NewQueueService(mockRedis, 5*time.Second)

	result, err := service.PopTransaction(context.Background())

	require.NoError(t, err)
	assert.NotNil(t, result)
	metadata := (*result)["metadata"].(map[string]any)
	assert.Equal(t, "api", metadata["source"])
}
