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
	docker build -t $(IMAGE_NAME) .

run:
	docker run --rm -it -p 8080:8080 $(IMAGE_NAME)

up:
	docker compose up -d

up-build:
	docker compose up -d --build

down:
	docker compose down

# コンテナとボリュームを消す
down-v:
	docker compose down -v

logs:
	docker logs -f minkan-api

# docker compose down -v # ボリュームを削除する

# ################
# DB関連
# ################

DB_NAME = minkan-mysql
DB_USER_NAME = app
DB_PASSWORD = password
DB_DATABASE = minkan

connect-db:
	docker exec -it $(DB_NAME) \
		mysql -u$(DB_USER_NAME) -p$(DB_PASSWORD) $(DB_DATABASE)

connect-db-root:
	docker exec -it $(DB_NAME) mysql


# ################
# OpenAPI
# ################
api-gen:
	go tool oapi-codegen -config ./api/config.yaml ./api/openapi.yaml