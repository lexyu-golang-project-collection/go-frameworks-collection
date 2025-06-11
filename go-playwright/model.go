package main

import (
	"time"

	"github.com/playwright-community/playwright-go"
)

// =============================================================================
// 基礎結構定義
// =============================================================================

type PlaywrightContext struct {
	pw      *playwright.Playwright
	browser playwright.Browser
	context playwright.BrowserContext
	page    playwright.Page
}

// AudioRequest 音頻請求結構
type AudioRequest struct {
	URL         string            `json:"url"`
	Method      string            `json:"method"`
	ContentType string            `json:"content_type"`
	Headers     map[string]string `json:"headers"`
}

// NetworkCapture 網路捕獲結構
type NetworkCapture struct {
	URL             string            `json:"url"`
	Method          string            `json:"method"`
	Headers         map[string]string `json:"headers"`
	Status          int               `json:"status"`
	StatusText      string            `json:"status_text"`
	ResponseHeaders map[string]string `json:"response_headers"`
	Timestamp       time.Time         `json:"timestamp"`
}
