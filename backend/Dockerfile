FROM golang:1.25.2-alpine3.21

# 作業ディレクトリ
WORKDIR /workspace

# 必要なパッケージをインストール（alpine用）
RUN apk update && apk add --no-cache \
    git \
    curl \
    wget \
    ca-certificates \
    postgresql-client \
    redis \
    make \
    vim \
    nano \
    jq \
    bash

# CGOを有効化
ENV CGO_ENABLED=1

# Go開発ツールをインストール
RUN go install github.com/air-verse/air@latest && \
    go install github.com/google/wire/cmd/wire@latest

# PATHにGOPATH/binを追加
ENV PATH="/root/go/bin:${PATH}"

CMD ["sleep", "infinity"]