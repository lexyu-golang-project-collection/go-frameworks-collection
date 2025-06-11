package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/playwright-community/playwright-go"
)

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
