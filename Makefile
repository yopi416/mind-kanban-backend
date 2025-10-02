.PHONY: fmt vet lint test

# コード整形
fmt:
	gofmt -s -w .

# 静的解析
vet:
	go vet ./...

# Lint（golangci-lintを利用する場合）
lint:
	golangci-lint run

# テスト実行
test:
	go test ./...
