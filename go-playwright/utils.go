package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"
)

// CaptureAllM3U8Requests 捕獲所有 m3u8 請求（調試用）
func (pc *PlaywrightContext) CaptureAllM3U8Requests(clickSelector string, timeout time.Duration) ([]string, error) {
	var m3u8URLs []string

	// 主頁面監聽
	pc.page.OnRequest(func(request playwright.Request) {
		url := request.URL()
		if strings.Contains(url, ".m3u8") {
			fmt.Printf("🔍 捕獲 m3u8: %s\n", url)
			m3u8URLs = append(m3u8URLs, url)
		}
	})

	// 點擊觸發
	fmt.Printf("點擊元素: %s\n", clickSelector)
	if err := pc.page.Locator(clickSelector).Click(); err != nil {
		return nil, fmt.Errorf("點擊失敗: %v", err)
	}

	// 等待一段時間收集所有請求
	fmt.Printf("等待 %v 收集請求...\n", timeout)
	time.Sleep(timeout)

	return m3u8URLs, nil
}

// SmartCaptureM3U8 智能捕獲 - 優先 easy_audio，否則返回所有
func (pc *PlaywrightContext) SmartCaptureM3U8(clickSelector string, timeout time.Duration) ([]string, error) {
	var allM3U8URLs []string
	var easyAudioURLs []string

	// 監聽所有請求
	pc.page.OnRequest(func(request playwright.Request) {
		url := request.URL()
		if strings.Contains(url, ".m3u8") {
			fmt.Printf("🔍 發現 m3u8: %s\n", url)
			allM3U8URLs = append(allM3U8URLs, url)

			if strings.Contains(url, "easy_audio") {
				fmt.Printf("✅ 找到 easy_audio: %s\n", url)
				easyAudioURLs = append(easyAudioURLs, url)
			}
		}
	})

	// 點擊觸發
	fmt.Printf("點擊元素: %s\n", clickSelector)
	if err := pc.page.Locator(clickSelector).Click(); err != nil {
		return nil, fmt.Errorf("點擊失敗: %v", err)
	}

	// 等待收集
	fmt.Printf("等待 %v 收集請求...\n", timeout)
	time.Sleep(timeout)

	// 優先返回 easy_audio，否則返回所有
	if len(easyAudioURLs) > 0 {
		fmt.Printf("✅ 返回 %d 個 easy_audio URL\n", len(easyAudioURLs))
		return easyAudioURLs, nil
	}

	if len(allM3U8URLs) > 0 {
		fmt.Printf("⚠️  未找到 easy_audio，返回所有 %d 個 m3u8 URL\n", len(allM3U8URLs))
		return allM3U8URLs, nil
	}

	return nil, fmt.Errorf("未找到任何 m3u8 URL")
}

// =============================================================================
// 輔助功能
// =============================================================================

// isAudioRequest 判斷是否為音頻相關請求
func isAudioRequest(url, contentType string) bool {
	audioExtensions := []string{".m3u8", ".mp3", ".wav", ".aac", ".m4a", ".flac"}
	audioTypes := []string{"audio/", "application/vnd.apple.mpegurl"}

	// 檢查 URL 副檔名
	for _, ext := range audioExtensions {
		if strings.Contains(strings.ToLower(url), ext) {
			return true
		}
	}

	// 檢查 Content-Type
	for _, audioType := range audioTypes {
		if strings.Contains(strings.ToLower(contentType), audioType) {
			return true
		}
	}

	return false
}

// matchPattern 簡單的模式匹配函數
func matchPattern(url, pattern string) bool {
	// 處理 **/*.m3u8 這樣的模式
	if strings.Contains(pattern, "**/*") {
		suffix := strings.TrimPrefix(pattern, "**/*")
		return strings.HasSuffix(url, suffix)
	}
	return strings.Contains(url, pattern)
}

// WaitForSpecificRequest 等待特定模式的請求
func (pc *PlaywrightContext) WaitForSpecificRequest(pattern string, clickSelector string) (string, error) {
	var capturedURL string
	done := make(chan bool, 1)

	// 監聽請求
	pc.page.OnRequest(func(request playwright.Request) {
		url := request.URL()
		if strings.Contains(url, pattern) || matchPattern(url, pattern) {
			capturedURL = url
			done <- true
		}
	})

	// 點擊觸發
	if err := pc.page.Locator(clickSelector).Click(); err != nil {
		return "", fmt.Errorf("點擊失敗: %v", err)
	}

	// 等待結果或超時
	select {
	case <-done:
		return capturedURL, nil
	case <-time.After(10 * time.Second):
		return "", fmt.Errorf("等待請求超時")
	}
}

// AdvancedNetworkCapture 進階網路捕獲 (包含回應內容)
func (pc *PlaywrightContext) AdvancedNetworkCapture(clickSelector string, timeout time.Duration) ([]NetworkCapture, error) {
	var captures []NetworkCapture
	done := make(chan bool, 1)

	// 同時監聽請求和回應
	pc.page.OnRequest(func(request playwright.Request) {
		contentType, _ := request.HeaderValue("content-type")
		if isAudioRequest(request.URL(), contentType) {
			headers, err := request.AllHeaders()
			if err != nil {
				headers = make(map[string]string)
			}

			capture := NetworkCapture{
				URL:       request.URL(),
				Method:    request.Method(),
				Headers:   headers,
				Timestamp: time.Now(),
			}
			captures = append(captures, capture)
		}
	})

	pc.page.OnResponse(func(response playwright.Response) {
		contentType, _ := response.HeaderValue("content-type")
		if isAudioRequest(response.URL(), contentType) {
			// 更新對應的 capture
			for i := range captures {
				if captures[i].URL == response.URL() {
					captures[i].Status = response.Status()
					captures[i].StatusText = response.StatusText()

					// 處理回應 headers
					responseHeaders, err := response.AllHeaders()
					if err != nil {
						responseHeaders = make(map[string]string)
					}
					captures[i].ResponseHeaders = responseHeaders
					break
				}
			}
		}
	})

	go func() {
		time.Sleep(timeout)
		done <- true
	}()

	if err := pc.page.Locator(clickSelector).Click(); err != nil {
		return nil, fmt.Errorf("點擊失敗: %v", err)
	}

	<-done
	return captures, nil
}
