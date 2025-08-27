FROM --platform=linux/amd64 golang:1.24.6

ARG TARGETOS=linux
ARG TARGETARCH=amd64
ENV GOOS=${TARGETOS} GOARCH=${TARGETARCH} CGO_ENABLED=1

RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
    btrfs-progs \
    libbtrfs-dev \
    linux-libc-dev \
    pkg-config \
    build-essential \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build

RUN go build -o btrfs-cli cmd/main.go

RUN mkdir -p /mnt/btrfs-test

WORKDIR /app

CMD ["./btrfs-cli"]
