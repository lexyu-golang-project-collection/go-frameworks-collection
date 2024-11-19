package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
	// Leetcode API endpoint
	ALGORITHMS_ENDPOINT_URL = "https://leetcode.com/api/problems/algorithms/"
	ALGORITHMS_BASE_URL     = "https://leetcode.com/problems/"
)

type Stat struct {
	QuestionTitleSlug string `json:"question__title_slug"`
	QuestionTitle     string `json:"question__title"`
	FrontendID        int    `json:"frontend_question_id"`
}

type Difficulty struct {
	Level int `json:"level"`
}

// Problem structure
type Problem struct {
	Stat       Stat       `json:"stat"`
	Difficulty Difficulty `json:"difficulty"`
	PaidOnly   bool       `json:"paid_only"`
}

type Data struct {
	StatStatusPairs []Problem `json:"stat_status_pairs"`
}

func (p Problem) Simplify() map[string]interface{} {
	return map[string]interface{}{
		"ID":         p.Stat.FrontendID,
		"Title":      p.Stat.QuestionTitle,
		"Slug":       p.Stat.QuestionTitleSlug,
		"Difficulty": p.Difficulty.Level,
		"PaidOnly":   p.PaidOnly,
	}
}

// Playwright Backlog
func main() {

	// Make HTTP request to fetch Leetcode problems
	resp, err := http.Get(ALGORITHMS_ENDPOINT_URL)
	if err != nil {
		log.Fatalf("could not fetch Leetcode API: %v", err)
	}
	defer resp.Body.Close()

	fmt.Println(resp)

	// Create a JSON-serializable representation of http.Response
	responseData := map[string]interface{}{
		"Status":        resp.Status,
		"StatusCode":    resp.StatusCode,
		"Proto":         resp.Proto,
		"Headers":       resp.Header,
		"ContentLength": resp.ContentLength,
	}

	// Marshal it to pretty JSON
	prettyJSON, err := json.MarshalIndent(responseData, "", " ")
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	fmt.Println(string(prettyJSON))

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("could not read response body: %v", err)
	}
	// Parse JSON response
	data := &Data{}
	err = json.Unmarshal(body, data)
	if err != nil {
		log.Fatalf("could not parse JSON: %v", err)
	}
	fmt.Println(len(data.StatStatusPairs))
	problem := data.StatStatusPairs[len(data.StatStatusPairs)-1]
	jsonB, _ := json.MarshalIndent(problem, "", " ")
	fmt.Println(string(jsonB))

	// Print problem URLs
	// for _, problem := range data.StatStatusPairs {
	// 	if !problem.PaidOnly {
	// 		// Build the problem URL and print it
	// 		url := ALGORITHMS_BASE_URL + problem.Stat.QuestionTitleSlug
	// 		fmt.Println(url)
	// 	}
	// }
}
