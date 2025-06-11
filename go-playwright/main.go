package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"
)

// =============================================================================
// M3U8 捕獲核心功能
// =============================================================================

// ExtractM3U8URLFromClick 點擊後從 iframe src 構造 m3u8 URL
func (pc *PlaywrightContext) ExtractM3U8URLFromClick(clickSelector string) (string, error) {
	// 點擊音頻按鈕
	fmt.Printf("點擊元素: %s\n", clickSelector)
	if err := pc.page.Locator(clickSelector).Click(); err != nil {
		return "", fmt.Errorf("點擊失敗: %v", err)
	}

	// 等待 iframe 出現
	fmt.Println("等待 iframe 出現...")
	iframeSelector := "iframe"
	iframeLocator := pc.page.Locator(iframeSelector)
	if err := iframeLocator.WaitFor(playwright.LocatorWaitForOptions{
		Timeout: playwright.Float(10000),
	}); err != nil {
		return "", fmt.Errorf("iframe 未出現: %v", err)
	}

	// 獲取 iframe 的 src 屬性
	src, err := iframeLocator.GetAttribute("src")
	if err != nil || src == "" {
		return "", fmt.Errorf("無法獲取 iframe src")
	}

	fmt.Printf("iframe src: %s\n", src)

	// 解析 voiceId 參數
	// 預期格式: /news/easy/player/audio-v5.html?voiceId=ne2025060919337_oUAsAp6iLwRQnhIj4ecXzFKbcUThqGqz06pLAzGB.m4a&title=...
	if !strings.Contains(src, "voiceId=") {
		return "", fmt.Errorf("iframe src 中找不到 voiceId 參數")
	}

	// 提取 voiceId
	parts := strings.Split(src, "voiceId=")
	if len(parts) < 2 {
		return "", fmt.Errorf("無法解析 voiceId")
	}

	voiceIdPart := parts[1]
	// 移除 .m4a 後綴和後續參數
	voiceId := strings.Split(voiceIdPart, ".m4a")[0]
	voiceId = strings.Split(voiceId, "&")[0]

	if voiceId == "" {
		return "", fmt.Errorf("voiceId 為空")
	}

	// 構造 m3u8 URL
	m3u8URL := fmt.Sprintf("https://vod-stream.nhk.jp/news/easy_audio/%s/index.m3u8", voiceId)

	fmt.Printf("✅ 提取的 voiceId: %s\n", voiceId)
	fmt.Printf("✅ 構造的 m3u8 URL: %s\n", m3u8URL)

	return m3u8URL, nil
}

// VerifyM3U8URL 驗證構造的 m3u8 URL 是否有效
func (pc *PlaywrightContext) VerifyM3U8URL(m3u8URL string) (bool, error) {
	// 建立新頁面測試 URL
	testPage, err := pc.context.NewPage()
	if err != nil {
		return false, err
	}
	defer testPage.Close()

	// 嘗試訪問 m3u8 URL
	response, err := testPage.Goto(m3u8URL)
	if err != nil {
		return false, fmt.Errorf("無法訪問 m3u8 URL: %v", err)
	}

	status := response.Status()
	fmt.Printf("M3U8 URL 狀態碼: %d\n", status)

	if status == 200 {
		// 檢查內容是否為 m3u8 格式
		content, err := testPage.Content()
		if err != nil {
			return false, err
		}

		if strings.Contains(content, "#EXTM3U") {
			fmt.Printf("✅ M3U8 URL 驗證成功！\n")
			return true, nil
		}
	}

	return false, fmt.Errorf("M3U8 URL 無效或無法訪問")
}

// =============================================================================
// Demo 函數
// =============================================================================

// RunM3U8CaptureDemo M3U8 捕獲示例
func RunM3U8CaptureDemo() {
	fmt.Println("=== M3U8 捕獲測試 ===")

	pc := NewPlaywrightContext()
	defer pc.Close()

	// 目標頁面和選擇器
	url := "https://www3.nhk.or.jp/news/easy/ne2025060919337/ne2025060919337.html"
	clickSelector := ".js-open-audio" // 修正後的選擇器

	// 導航到目標頁面
	fmt.Printf("導航到: %s\n", url)
	if err := pc.NavigateToPage(url); err != nil {
		log.Printf("頁面導航失敗: %v", err)
		return
	}

	// 等待頁面載入完成
	fmt.Println("等待頁面載入...")
	time.Sleep(5 * time.Second)

	// 先檢查所有網路請求
	fmt.Println("\n--- 設置全域網路監聽 ---")
	pc.page.OnRequest(func(request playwright.Request) {
		url := request.URL()
		fmt.Printf("🌐 網路請求: %s\n", url)
		if strings.Contains(url, ".m3u8") {
			fmt.Printf("🎵 發現 M3U8: %s\n", url)
		}
	})

	// 檢查按鈕狀態
	fmt.Printf("\n--- 檢查按鈕狀態 ---\n")
	locator := pc.page.Locator(clickSelector)
	count, err := locator.Count()
	if err != nil {
		fmt.Printf("❌ 檢查按鈕失敗: %v\n", err)
		return
	}
	fmt.Printf("找到 %d 個按鈕\n", count)

	if count == 0 {
		fmt.Printf("❌ 找不到按鈕，嘗試其他選擇器...\n")
		// 嘗試其他可能的選擇器
		alternativeSelectors := []string{
			".article-buttons__audio",
			"[class*='audio']",
			"a[href='#']",
		}

		for _, selector := range alternativeSelectors {
			testLocator := pc.page.Locator(selector)
			testCount, _ := testLocator.Count()
			fmt.Printf("選擇器 '%s': %d 個元素\n", selector, testCount)
		}
		return
	}

	// 檢查按鈕是否可見和可點擊
	isVisible, _ := locator.IsVisible()
	isEnabled, _ := locator.IsEnabled()
	fmt.Printf("按鈕可見: %t, 可點擊: %t\n", isVisible, isEnabled)

	// 等待按鈕可點擊
	fmt.Println("等待按鈕可點擊...")
	if err := locator.WaitFor(playwright.LocatorWaitForOptions{
		State:   playwright.WaitForSelectorStateVisible,
		Timeout: playwright.Float(10000),
	}); err != nil {
		fmt.Printf("❌ 按鈕等待失敗: %v\n", err)
		return
	}

	fmt.Println("\n--- 執行點擊 ---")
	if err := locator.Click(); err != nil {
		fmt.Printf("❌ 點擊失敗: %v\n", err)
		return
	}
	fmt.Println("✅ 點擊成功")

	// 檢查點擊後的變化
	fmt.Println("\n--- 檢查點擊後變化 ---")
	time.Sleep(2 * time.Second)

	// 檢查 iframe 是否出現
	iframeLocator := pc.page.Locator("iframe")
	iframeCount, _ := iframeLocator.Count()
	fmt.Printf("找到 %d 個 iframe\n", iframeCount)

	// 檢查 audio 元素
	audioLocator := pc.page.Locator("audio")
	audioCount, _ := audioLocator.Count()
	fmt.Printf("找到 %d 個 audio 元素\n", audioCount)

	// 檢查是否有 is-active class
	activeLocator := pc.page.Locator(".is-active")
	activeCount, _ := activeLocator.Count()
	fmt.Printf("找到 %d 個 .is-active 元素\n", activeCount)

	// 長時間等待網路請求
	fmt.Println("\n--- 等待網路請求 (20秒) ---")
	time.Sleep(20 * time.Second)

	fmt.Println("✅ 測試完成")
}

// DeepDebugM3U8 深度除錯版本
func DeepDebugM3U8() {
	fmt.Println("=== 深度除錯 M3U8 ===")

	pc := NewPlaywrightContext()
	defer pc.Close()

	url := "https://www3.nhk.or.jp/news/easy/ne2025060919337/ne2025060919337.html"

	// 監聽所有可能的網路活動
	pc.page.OnRequest(func(request playwright.Request) {
		reqURL := request.URL()
		method := request.Method()
		fmt.Printf("📤 請求: %s %s\n", method, reqURL)

		// 檢查任何音頻相關的請求
		if strings.Contains(reqURL, "audio") ||
			strings.Contains(reqURL, "m3u8") ||
			strings.Contains(reqURL, "mp3") ||
			strings.Contains(reqURL, "m4a") ||
			strings.Contains(reqURL, "vod-stream") {
			fmt.Printf("🎵 音頻相關請求: %s\n", reqURL)
		}
	})

	pc.page.OnResponse(func(response playwright.Response) {
		respURL := response.URL()
		status := response.Status()

		if strings.Contains(respURL, "audio") ||
			strings.Contains(respURL, "m3u8") ||
			strings.Contains(respURL, "vod-stream") {
			fmt.Printf("📥 音頻回應: %d %s\n", status, respURL)
		}
	})

	// 導航並執行
	fmt.Printf("導航到: %s\n", url)
	pc.NavigateToPage(url)

	fmt.Println("等待初始載入...")
	time.Sleep(5 * time.Second)

	fmt.Println("執行點擊...")
	if err := pc.page.Locator(".js-open-audio").Click(); err != nil {
		fmt.Printf("點擊失敗: %v\n", err)
		return
	}

	fmt.Println("等待 30 秒觀察...")
	time.Sleep(30 * time.Second)
}

// TestClickButton 測試按鈕點擊功能
func TestClickButton() {
	fmt.Println("=== 測試按鈕點擊 ===")

	pc := NewPlaywrightContext()
	defer pc.Close()

	url := "https://www3.nhk.or.jp/news/easy/ne2025060919337/ne2025060919337.html"

	fmt.Printf("導航到: %s\n", url)
	if err := pc.NavigateToPage(url); err != nil {
		log.Printf("頁面導航失敗: %v", err)
		return
	}

	time.Sleep(2 * time.Second)

	// 測試按鈕是否存在
	clickSelector := ".js-open-audio"
	locator := pc.page.Locator(clickSelector)

	count, err := locator.Count()
	if err != nil {
		fmt.Printf("❌ 檢查元素失敗: %v\n", err)
		return
	}

	fmt.Printf("找到 %d 個匹配的元素\n", count)

	if count > 0 {
		fmt.Printf("✅ 嘗試點擊按鈕...\n")
		if err := locator.Click(); err != nil {
			fmt.Printf("❌ 點擊失敗: %v\n", err)
		} else {
			fmt.Printf("✅ 點擊成功！\n")
			time.Sleep(3 * time.Second) // 觀察結果
		}
	} else {
		fmt.Printf("❌ 找不到目標按鈕\n")
	}
}

// DownloadM3U8Content 直接下載 m3u8 內容到本地
func DownloadM3U8Content(m3u8URL, saveFilename string) error {
	fmt.Printf("正在下載 M3U8: %s\n", m3u8URL)

	// 建立 HTTP 客戶端
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// 建立請求
	req, err := http.NewRequest("GET", m3u8URL, nil)
	if err != nil {
		return fmt.Errorf("建立請求失敗: %v", err)
	}

	// 設置 headers 模擬瀏覽器
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept", "application/vnd.apple.mpegurl,*/*")
	req.Header.Set("Referer", "https://www3.nhk.or.jp/")

	// 發送請求
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("請求失敗: %v", err)
	}
	defer resp.Body.Close()

	// 檢查狀態碼
	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP 錯誤: %d", resp.StatusCode)
	}

	// 讀取內容
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("讀取內容失敗: %v", err)
	}

	// 檢查是否為 m3u8 格式
	contentStr := string(content)
	if !strings.Contains(contentStr, "#EXTM3U") {
		return fmt.Errorf("回應不是有效的 m3u8 格式")
	}

	// 儲存到文件
	err = os.WriteFile(saveFilename, content, 0644)
	if err != nil {
		return fmt.Errorf("儲存文件失敗: %v", err)
	}

	fmt.Printf("✅ M3U8 內容已儲存到: %s\n", saveFilename)
	fmt.Printf("📄 內容大小: %d bytes\n", len(content))

	// 顯示前幾行內容
	lines := strings.Split(contentStr, "\n")
	fmt.Println("📋 M3U8 內容預覽:")
	for i, line := range lines {
		if i >= 5 {
			break
		}
		fmt.Printf("  %s\n", line)
	}

	return nil
}

// CompleteM3U8Workflow 完整的 M3U8 提取和下載流程
func CompleteM3U8Workflow() {
	fmt.Println("=== 完整 M3U8 提取流程 ===")

	pc := NewPlaywrightContext()
	defer pc.Close()

	url := "https://www3.nhk.or.jp/news/easy/ne2025060919337/ne2025060919337.html"

	// 導航和提取 URL
	fmt.Printf("導航到: %s\n", url)
	pc.NavigateToPage(url)
	time.Sleep(3 * time.Second)

	m3u8URL, err := pc.ExtractM3U8URLFromClick(".js-open-audio")
	if err != nil {
		log.Printf("❌ 提取失敗: %v\n", err)
		return
	}

	fmt.Printf("🎯 M3U8 URL: %s\n", m3u8URL)

	// 直接下載 m3u8 內容
	// filename := "audio_playlist.m3u8"
	// if err := DownloadM3U8Content(m3u8URL, filename); err != nil {
	// 	log.Printf("❌ 下載失敗: %v\n", err)
	// 	return
	// }

	// fmt.Printf("🎉 流程完成！M3U8 列表已儲存到 %s\n", filename)
}

// =============================================================================
// 主程式
// =============================================================================

func main() {
	// 安裝 Playwright
	fmt.Println("安裝 Playwright...")
	if err := playwright.Install(); err != nil {
		log.Fatalf("Playwright Install() 失敗: %v", err)
	}

	// 選擇要執行的測試
	// fmt.Println("開始深度除錯測試...")
	// 先執行深度除錯
	// DeepDebugM3U8()

	// fmt.Println("\n" + strings.Repeat("=", 60) + "\n")

	// 再執行詳細測試
	// RunM3U8CaptureDemo()

	// 執行最終解決方案
	CompleteM3U8Workflow()
}
