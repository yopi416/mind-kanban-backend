####################### Build stage #######################
FROM golang:1.25.1-bookworm AS builder

# go.mod と go.sum を app ディレクトリにコピーし、モジュールをダウンロード
WORKDIR /app
# COPY go.mod go.sum ./
COPY go.mod ./
RUN go mod download

# ルートディレクトリの中身を app フォルダにコピーする
COPY . .

# 実行ファイルの作成
# -o はアウトプット名
WORKDIR /app/cmd/api
RUN go build -trimpath -ldflags="-w -s" -o /app/bin/app

####################### Run stage #######################
FROM gcr.io/distroless/base-debian12

# FROM golang:1.25-slim-bullseye

# 依存（必要に応じてca-certificatesなど）
# RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates \
#     && rm -rf /var/lib/apt/lists/*

# Build stage からビルドされた app だけを Run stage にコピーする。（重要）
WORKDIR /app
COPY --from=builder app/bin/app /app/app

# EXPOSE 命令は、実際にポートを公開するわけではない。
# 今回、docker-compose.yml において、api コンテナは 8080 ポートを解放するため「8080」とする。
EXPOSE 8080

# バイナリファイルの実行
CMD [ "/app/app" ]