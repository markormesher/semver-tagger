FROM docker.io/golang:1.24.0@sha256:3f7444391c51a11a039bf0359ee81cc64e663c17d787ad0e637a4de1a3f62a71 AS builder
WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY ./main.go ./main.go
COPY ./internal ./internal

RUN go build -o ./build/main ./main.go

# ---

FROM debian:bookworm@sha256:35286826a88dc879b4f438b645ba574a55a14187b483d09213a024dc0c0a64ed
WORKDIR /app

LABEL image.registry=ghcr.io
LABEL image.name=markormesher/semver-tagger

RUN apt update \
  && apt install -y --no-install-recommends \
  git \
  ca-certificates \
  && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/build/main /usr/local/bin/semver-tagger

CMD ["/usr/local/bin/semver-tagger", "--help"]
