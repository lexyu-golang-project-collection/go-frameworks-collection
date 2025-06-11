package main

import (
	"log"

	"github.com/playwright-community/playwright-go"
)

// =============================================================================
// 基礎 Playwright 操作
// =============================================================================

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
