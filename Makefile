.PHONY: fmt vet lint test build

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


# ################
# Docker関連
# ################

IMAGE_NAME=minkan-backend

build:
	docker build -f build/Dockerfile -t $(IMAGE_NAME) .

run:
	docker run --rm -it -p 8080:8080 $(IMAGE_NAME)