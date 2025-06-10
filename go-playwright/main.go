package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"
)

// PlaywrightContext 包裝 Playwright 相關物件
type PlaywrightContext struct {
	pw      *playwright.Playwright
	browser playwright.Browser
	context playwright.BrowserContext
	page    playwright.Page
}

// NewPlaywrightContext 建立新的 Playwright 上下文
func NewPlaywrightContext() *PlaywrightContext {
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("Playwright 啟動失敗: %v", err)
	}

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	if err != nil {
		log.Fatalf("瀏覽器啟動失敗: %v", err)
	}

	ctx, err := browser.NewContext()
	if err != nil {
		log.Fatalf("建立 context 失敗: %v", err)
	}

	page, err := ctx.NewPage()
	if err != nil {
		log.Fatalf("建立 page 失敗: %v", err)
	}

	return &PlaywrightContext{
		pw:      pw,
		browser: browser,
		context: ctx,
		page:    page,
	}
}

// Close 關閉所有 Playwright 資源
func (pc *PlaywrightContext) Close() {
	_ = pc.browser.Close()
	_ = pc.pw.Stop()
}

// NavigateToPage 導航到指定頁面
func (pc *PlaywrightContext) NavigateToPage(url string) error {
	_, err := pc.page.Goto(url)
	return err
}

// === 基礎文字擷取 API Demo ===

// GetSingleText 取得單一元素的文字內容
func (pc *PlaywrightContext) GetSingleText(selector string) (string, error) {
	return pc.page.Locator(selector).TextContent()
}

// GetAllTexts 取得多個元素的文字內容
func (pc *PlaywrightContext) GetAllTexts(selector string) ([]string, error) {
	locator := pc.page.Locator(selector)
	count, err := locator.Count()
	if err != nil {
		return nil, err
	}

	var texts []string
	for i := 0; i < count; i++ {
		text, err := locator.Nth(i).TextContent()
		if err != nil {
			continue
		}
		texts = append(texts, text)
	}
	return texts, nil
}

// GetTextByIndex 根據索引取得特定元素的文字
func (pc *PlaywrightContext) GetTextByIndex(selector string, index int) (string, error) {
	return pc.page.Locator(selector).Nth(index).TextContent()
}

// === 屬性擷取 API Demo ===

// GetSingleAttribute 取得單一元素的屬性值
func (pc *PlaywrightContext) GetSingleAttribute(selector, attribute string) (string, error) {
	val, err := pc.page.Locator(selector).GetAttribute(attribute)
	if err != nil {
		return "", err
	}
	// Playwright Go 的 GetAttribute 返回 string，空字符串表示屬性不存在
	if val == "" {
		return "", fmt.Errorf("屬性 %s 不存在或為空", attribute)
	}
	return val, nil
}

// GetMultipleAttributes 取得多個元素的相同屬性值
func (pc *PlaywrightContext) GetMultipleAttributes(selector, attribute string) ([]string, error) {
	script := fmt.Sprintf(`elements => elements.map(el => el.getAttribute("%s"))`, attribute)
	result, err := pc.page.Locator(selector).EvaluateAll(script)
	if err != nil {
		return nil, err
	}

	var attributes []string
	for _, v := range result.([]interface{}) {
		if v != nil {
			attributes = append(attributes, v.(string))
		}
	}
	return attributes, nil
}

// === 鏈接擷取 API Demo ===

// ExtractLinksFromClass 從指定 class 底下提取所有 a 元素的 href 鏈接
func (pc *PlaywrightContext) ExtractLinksFromClass(className string) ([]string, error) {
	selector := fmt.Sprintf(".%s a", className)
	return pc.GetMultipleAttributes(selector, "href")
}

// ExtractLinksFromClassWithText 從指定 class 底下提取帶有特定文字的 a 元素鏈接
func (pc *PlaywrightContext) ExtractLinksFromClassWithText(className, textFilter string) ([]string, error) {
	selector := fmt.Sprintf(".%s a", className)
	locator := pc.page.Locator(selector)
	filtered := locator.Filter(playwright.LocatorFilterOptions{
		HasText: playwright.String(textFilter),
	})

	count, err := filtered.Count()
	if err != nil {
		return nil, err
	}

	var links []string
	for i := 0; i < count; i++ {
		href, err := filtered.Nth(i).GetAttribute("href")
		if err != nil || href == "" {
			continue
		}
		links = append(links, href)
	}
	return links, nil
}

// ExtractAllLinksWithDetails 提取頁面所有鏈接的詳細資訊
func (pc *PlaywrightContext) ExtractAllLinksWithDetails(selector string) ([]LinkDetails, error) {
	locator := pc.page.Locator(selector)
	count, err := locator.Count()
	if err != nil {
		return nil, err
	}

	var linkDetails []LinkDetails
	for i := 0; i < count; i++ {
		element := locator.Nth(i)

		href, _ := element.GetAttribute("href")
		text, _ := element.TextContent()
		title, _ := element.GetAttribute("title")

		detail := LinkDetails{
			Index: i,
			Text:  strings.TrimSpace(text),
		}

		if href != "" {
			detail.URL = href
		}
		if title != "" {
			detail.Title = title
		}

		linkDetails = append(linkDetails, detail)
	}
	return linkDetails, nil
}

// LinkDetails 鏈接詳細資訊結構
type LinkDetails struct {
	Index int    `json:"index"`
	URL   string `json:"url"`
	Text  string `json:"text"`
	Title string `json:"title"`
}

// === HTML 內容擷取 API Demo ===

// GetInnerHTML 取得元素的 innerHTML
func (pc *PlaywrightContext) GetInnerHTML(selector string) (string, error) {
	return pc.page.Locator(selector).First().InnerHTML()
}

// GetOuterHTML 取得元素的 outerHTML
func (pc *PlaywrightContext) GetOuterHTML(selector string) (string, error) {
	return pc.page.Locator(selector).InnerHTML()
}

// === 表單元素 API Demo ===

// GetInputValue 取得 input 元素的值
func (pc *PlaywrightContext) GetInputValue(selector string) (string, error) {
	return pc.page.Locator(selector).InputValue()
}

// GetSelectValue 取得 select 元素的值
func (pc *PlaywrightContext) GetSelectValue(selector string) (string, error) {
	return pc.page.Locator(selector).InputValue()
}

// IsElementChecked 檢查 checkbox 或 radio 是否被選中
func (pc *PlaywrightContext) IsElementChecked(selector string) (bool, error) {
	return pc.page.Locator(selector).IsChecked()
}

// === 元素過濾 API Demo ===

// FilterElementsByText 根據文字內容過濾元素
func (pc *PlaywrightContext) FilterElementsByText(selector, text string) ([]string, error) {
	filtered := pc.page.Locator(selector).Filter(playwright.LocatorFilterOptions{
		HasText: playwright.String(text),
	})

	count, err := filtered.Count()
	if err != nil {
		return nil, err
	}

	var results []string
	for i := 0; i < count; i++ {
		content, err := filtered.Nth(i).TextContent()
		if err != nil {
			continue
		}
		results = append(results, content)
	}
	return results, nil
}

// FilterElementsByAttribute 根據屬性過濾元素
func (pc *PlaywrightContext) FilterElementsByAttribute(selector, attribute, value string) ([]string, error) {
	script := fmt.Sprintf(`elements => elements.filter(el => el.getAttribute("%s") === "%s").map(el => el.textContent)`, attribute, value)
	result, err := pc.page.Locator(selector).EvaluateAll(script)
	if err != nil {
		return nil, err
	}

	var filteredTexts []string
	for _, v := range result.([]interface{}) {
		if v != nil {
			filteredTexts = append(filteredTexts, v.(string))
		}
	}
	return filteredTexts, nil
}

// === Demo 執行函數 ===

// RunTextExtractionDemo 執行文字擷取示例
func RunTextExtractionDemo(url string) {
	fmt.Println("=== 文字擷取 Demo ===")
	pc := NewPlaywrightContext()
	defer pc.Close()

	if err := pc.NavigateToPage(url); err != nil {
		log.Printf("頁面導航失敗: %v", err)
		return
	}

	// 取得標題文字
	if title, err := pc.GetSingleText("h1"); err == nil {
		fmt.Printf("頁面標題: %s\n", title)
	}

	// 取得所有列表項目
	if texts, err := pc.GetAllTexts("ul > li"); err == nil {
		fmt.Printf("找到 %d 個列表項目:\n", len(texts))
		for i, text := range texts {
			fmt.Printf("  %d: %s\n", i+1, text)
		}
	}
}

// RunLinkExtractionDemo 執行鏈接擷取示例
func RunLinkExtractionDemo(url, className string) {
	fmt.Println("=== 鏈接擷取 Demo ===")
	pc := NewPlaywrightContext()
	defer pc.Close()

	if err := pc.NavigateToPage(url); err != nil {
		log.Printf("頁面導航失敗: %v", err)
		return
	}

	// 從指定 class 提取所有鏈接
	if links, err := pc.ExtractLinksFromClass(className); err == nil {
		fmt.Printf("從 class '%s' 找到 %d 個鏈接:\n", className, len(links))
		for i, link := range links {
			fmt.Printf("  %d: %s\n", i+1, link)
		}
	}

	// 提取所有鏈接的詳細資訊
	if details, err := pc.ExtractAllLinksWithDetails("a"); err == nil {
		fmt.Printf("\n所有鏈接詳細資訊 (前5個):\n")
		for i, detail := range details {
			if i >= 5 {
				break
			}
			fmt.Printf("  %d: %s -> %s\n", detail.Index+1, detail.Text, detail.URL)
		}
	}
}

// RunAttributeExtractionDemo 執行屬性擷取示例
func RunAttributeExtractionDemo(url string) {
	fmt.Println("=== 屬性擷取 Demo ===")
	pc := NewPlaywrightContext()
	defer pc.Close()

	if err := pc.NavigateToPage(url); err != nil {
		log.Printf("頁面導航失敗: %v", err)
		return
	}

	// 取得第一個鏈接的 href
	if href, err := pc.GetSingleAttribute("a", "href"); err == nil && href != "" {
		fmt.Printf("第一個鏈接的 href: %s\n", href)
	}

	// 取得所有圖片的 alt 屬性
	if alts, err := pc.GetMultipleAttributes("img", "alt"); err == nil {
		fmt.Printf("找到 %d 個圖片的 alt 屬性:\n", len(alts))
		for i, alt := range alts {
			if alt != "" {
				fmt.Printf("  %d: %s\n", i+1, alt)
			}
		}
	}
}

// RunFormElementDemo 執行表單元素示例
func RunFormElementDemo(url string) {
	fmt.Println("=== 表單元素 Demo ===")
	pc := NewPlaywrightContext()
	defer pc.Close()

	if err := pc.NavigateToPage(url); err != nil {
		log.Printf("頁面導航失敗: %v", err)
		return
	}

	// 嘗試取得搜尋框的值
	if value, err := pc.GetInputValue("input[name=q]"); err == nil {
		fmt.Printf("搜尋框當前值: '%s'\n", value)
	}
}

// =========================================================================

// CaptureM3U8AfterClick 點擊元素後捕獲 m3u8 URL
func (pc *PlaywrightContext) CaptureM3U8AfterClick(clickSelector string, timeout time.Duration) ([]string, error) {
	var m3u8URLs []string
	done := make(chan bool, 1)

	// 監聽網路請求
	pc.page.OnRequest(func(request playwright.Request) {
		url := request.URL()
		if strings.Contains(url, ".m3u8") {
			fmt.Printf("發現 m3u8: %s\n", url)
			m3u8URLs = append(m3u8URLs, url)

			// 優先選擇 easy_audio 版本
			if strings.Contains(url, "easy_audio") {
				fmt.Printf("找到目標音頻: %s\n", url)
				done <- true
			}
		}
	})

	// 點擊觸發元素
	if err := pc.page.Locator(clickSelector).Click(); err != nil {
		return nil, fmt.Errorf("點擊失敗: %v", err)
	}

	// 等待結果或超時
	select {
	case <-done:
		return m3u8URLs, nil
	case <-time.After(timeout):
		if len(m3u8URLs) > 0 {
			return m3u8URLs, nil
		}
		return nil, fmt.Errorf("超時未找到目標 m3u8")
	}
}

// CaptureAllAudioRequests 捕獲所有音頻相關請求
func (pc *PlaywrightContext) CaptureAllAudioRequests(clickSelector string, timeout time.Duration) ([]AudioRequest, error) {
	var audioRequests []AudioRequest
	done := make(chan bool, 1)

	pc.page.OnRequest(func(request playwright.Request) {
		url := request.URL()
		contentType, _ := request.HeaderValue("content-type")

		// 檢查是否為音頻相關請求
		if isAudioRequest(url, contentType) {
			headers, err := request.AllHeaders()
			if err != nil {
				headers = make(map[string]string) // 如果取得 headers 失敗，使用空 map
			}

			audioReq := AudioRequest{
				URL:         url,
				Method:      request.Method(),
				ContentType: contentType,
				Headers:     headers,
			}
			audioRequests = append(audioRequests, audioReq)
			fmt.Printf("音頻請求: %s\n", url)
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
	return audioRequests, nil
}

// AudioRequest 音頻請求結構
type AudioRequest struct {
	URL         string            `json:"url"`
	Method      string            `json:"method"`
	ContentType string            `json:"content_type"`
	Headers     map[string]string `json:"headers"`
}

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

// matchPattern 簡單的模式匹配函數
func matchPattern(url, pattern string) bool {
	// 處理 **/*.m3u8 這樣的模式
	if strings.Contains(pattern, "**/*") {
		suffix := strings.TrimPrefix(pattern, "**/*")
		return strings.HasSuffix(url, suffix)
	}
	return strings.Contains(url, pattern)
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

// Demo 使用範例
func RunM3U8CaptureDemo() {
	pc := NewPlaywrightContext()
	defer pc.Close()

	// 導航到目標頁面
	pc.NavigateToPage("https://www3.nhk.or.jp/news/easy/ne2025060919337/ne2025060919337.html")

	// 方法 1: 簡單捕獲 m3u8
	m3u8URLs, err := pc.CaptureM3U8AfterClick(".article-buttons__audio js-open-audio is-active", 5*time.Second)
	if err != nil {
		log.Printf("捕獲失敗: %v", err)
		return
	}

	fmt.Printf("找到 %d 個 m3u8 URL:\n", len(m3u8URLs))
	for _, url := range m3u8URLs {
		fmt.Printf("- %s\n", url)
	}

	// 方法 2: 等待特定請求模式
	// m3u8URL, err := pc.WaitForSpecificRequest("**/*.m3u8", ".play-button")
	// if err != nil {
	// 	log.Printf("等待請求失敗: %v", err)
	// 	return
	// }
	// fmt.Printf("M3U8 URL: %s\n", m3u8URL)
}

// =========================================================================

func main() {
	// 安裝 Playwright
	if err := playwright.Install(); err != nil {
		log.Fatalf("Playwright Install() 失敗: %v", err)
	}

	// url := "https://www3.nhk.or.jp/news/easy/"

	// 執行各種 Demo
	// RunTextExtractionDemo(url)
	// fmt.Println()

	// RunLinkExtractionDemo(url, "news-list__item")
	// fmt.Println()

	// RunAttributeExtractionDemo(url)
	// fmt.Println()

	// RunFormElementDemo(url)

	RunM3U8CaptureDemo()
}
