package main

import (
	"encoding/json"
	"testing"
)

// 測試 Simplify 方法
func TestProblemSimplify(t *testing.T) {
	problem := Problem{
		Stat: Stat{
			QuestionTitleSlug: "two-sum",
			QuestionTitle:     "Two Sum",
			FrontendID:        1,
		},
		Difficulty: Difficulty{
			Level: 1,
		},
		PaidOnly: false,
	}

	expected := map[string]interface{}{
		"ID":         1,
		"Title":      "Two Sum",
		"Slug":       "two-sum",
		"Difficulty": 1,
		"PaidOnly":   false,
	}

	result := problem.Simplify()
	for key, val := range expected {
		if result[key] != val {
			t.Errorf("expected %v, got %v for key %s", val, result[key], key)
		}
	}
}

// 測試 JSON 解析
func TestUnmarshalJSON(t *testing.T) {
	jsonResponse := `{
		"stat_status_pairs": [
			{
				"stat": {
					"question__title_slug": "two-sum",
					"question__title": "Two Sum",
					"frontend_question_id": 1
				},
				"difficulty": {
					"level": 1
				},
				"paid_only": false
			}
		]
	}`

	data := &Data{}
	err := json.Unmarshal([]byte(jsonResponse), data)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if len(data.StatStatusPairs) != 1 {
		t.Errorf("expected 1 problem, got %d", len(data.StatStatusPairs))
	}

	problem := data.StatStatusPairs[0]
	if problem.Stat.QuestionTitle != "Two Sum" {
		t.Errorf("expected 'Two Sum', got '%s'", problem.Stat.QuestionTitle)
	}
	if problem.Stat.FrontendID != 1 {
		t.Errorf("expected ID 1, got %d", problem.Stat.FrontendID)
	}
	if problem.Difficulty.Level != 1 {
		t.Errorf("expected difficulty 1, got %d", problem.Difficulty.Level)
	}
	if problem.PaidOnly {
		t.Errorf("expected PaidOnly to be false, got true")
	}
}
