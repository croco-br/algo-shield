package rules

import (
	"context"
	"errors"
	"testing"

	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_RuleService_LoadRules_WhenSuccess_ThenStoresRules(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedRules := []models.Rule{
		{Name: "rule1", Action: models.ActionBlock, Conditions: map[string]any{"field": "amount"}},
		{Name: "rule2", Action: models.ActionAllow, Conditions: map[string]any{"field": "currency"}},
	}
	mockRepo := NewMockRuleReader(ctrl)
	mockRepo.EXPECT().LoadRules(gomock.Any()).Return(expectedRules, nil)
	service := NewRuleService(mockRepo)

	err := service.LoadRules(context.Background())

	require.NoError(t, err)
	rules := service.GetRules()
	assert.Equal(t, expectedRules, rules)
}

func Test_RuleService_LoadRules_WhenRepositoryFails_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRuleReader(ctrl)
	mockRepo.EXPECT().LoadRules(gomock.Any()).Return(nil, errors.New("database error"))
	service := NewRuleService(mockRepo)

	err := service.LoadRules(context.Background())

	assert.Error(t, err)
}

func Test_RuleService_GetRules_WhenEmpty_ThenReturnsEmptySlice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRuleReader(ctrl)
	service := NewRuleService(mockRepo)

	rules := service.GetRules()

	assert.Empty(t, rules)
}

func Test_RuleService_GetRules_WhenRulesLoaded_ThenReturnsCopy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedRules := []models.Rule{
		{Name: "rule1", Action: models.ActionAllow, Conditions: map[string]any{}},
	}
	mockRepo := NewMockRuleReader(ctrl)
	mockRepo.EXPECT().LoadRules(gomock.Any()).Return(expectedRules, nil)
	service := NewRuleService(mockRepo)
	err := service.LoadRules(context.Background())
	require.NoError(t, err)

	rules1 := service.GetRules()
	rules2 := service.GetRules()

	assert.Equal(t, rules1, rules2)
	assert.NotSame(t, &rules1, &rules2, "should return a copy, not the same slice")
}

func Test_RuleService_LoadRules_WhenCalledMultipleTimes_ThenUpdatesRules(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	firstRules := []models.Rule{{Name: "rule1"}}
	secondRules := []models.Rule{{Name: "rule2"}, {Name: "rule3"}}
	mockRepo := NewMockRuleReader(ctrl)
	mockRepo.EXPECT().LoadRules(gomock.Any()).Return(firstRules, nil)
	mockRepo.EXPECT().LoadRules(gomock.Any()).Return(secondRules, nil)
	service := NewRuleService(mockRepo)

	err := service.LoadRules(context.Background())
	require.NoError(t, err)
	rules1 := service.GetRules()
	assert.Len(t, rules1, 1)

	err = service.LoadRules(context.Background())
	require.NoError(t, err)
	rules2 := service.GetRules()

	assert.Len(t, rules2, 2)
	assert.Equal(t, secondRules, rules2)
}

func Test_RuleService_GetRules_WhenModifyingReturned_ThenDoesNotAffectStored(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedRules := []models.Rule{{Name: "rule1", Action: models.ActionAllow, Conditions: map[string]any{}}}
	mockRepo := NewMockRuleReader(ctrl)
	mockRepo.EXPECT().LoadRules(gomock.Any()).Return(expectedRules, nil)
	service := NewRuleService(mockRepo)
	err := service.LoadRules(context.Background())
	require.NoError(t, err)

	rules := service.GetRules()
	rules[0].Name = "modified"

	storedRules := service.GetRules()
	assert.Equal(t, "rule1", storedRules[0].Name, "stored rules should not be affected by external modification")
}
