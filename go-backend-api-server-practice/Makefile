export LANG=zh_TW.UTF-8

# 設置 Go 編譯器
GO := go

# 設置主程序的路徑
MAIN := ./cmd/main.go

# 設置輸出二進制文件的名稱
BINARY_NAME := app

# 設置構建標誌
BUILD_FLAGS := -v

# 默認目標
all: build
    
# 構建目標
build:
	$(GO) build $(BUILD_FLAGS) -o $(BINARY_NAME) $(MAIN)

# 運行目標
run:
	./$(BINARY_NAME)

# 測試目標
test:
	$(GO) test -v ./...

# 清理目標
clean:
	$(GO) clean
	rm -f $(BINARY_NAME)

# 幫助目標
help:
	@echo "執行方法:"
	@echo "make build - 編譯程序"
	@echo "make run   - 編譯並運行程序"
	@echo "make test  - 運行測試"
	@echo "make clean - 清理編譯產物"



.PHONY: all build run test clean help