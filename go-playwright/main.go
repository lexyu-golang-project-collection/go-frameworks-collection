package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
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

	writeJSONToFile(data.StatStatusPairs, "problems.json")

	urls := []string{}
	for _, problem := range data.StatStatusPairs {
		if !problem.PaidOnly {
			url := ALGORITHMS_BASE_URL + problem.Stat.QuestionTitleSlug
			urls = append(urls, url)
		}
	}
	writeJSONArrayToFile(urls, "urls.json")
}

func writeJSONToFile(data []Problem, filename string) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory error: %v", err)
	}

	fullPath := filepath.Join(currentDir, filename)

	// using buffer
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)

	if err := enc.Encode(data); err != nil {
		return fmt.Errorf("encoding error: %v", err)
	}

	f, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("file open error: %v", err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	if _, err := w.Write(buf.Bytes()); err != nil {
		return fmt.Errorf("write error: %v", err)
	}

	// file location
	fmt.Printf("File written to: %s\n", fullPath)
	return w.Flush()
}

func writeJSONArrayToFile(urls []string, filename string) error {
	jsonData, err := json.MarshalIndent(urls, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal error: %v", err)
	}

	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("file open error: %v", err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	if _, err := w.Write(jsonData); err != nil {
		return fmt.Errorf("write error: %v", err)
	}

	return w.Flush()
}
