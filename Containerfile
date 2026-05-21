FROM docker.io/golang:1.26.3@sha256:6df14f4a4bc9d979a3721f488981e0d1b318006377e473ed23d026796f5f4c0a AS builder
WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY ./main.go ./main.go
COPY ./internal ./internal

RUN go build -o ./build/main ./main.go

# ---

FROM docker.io/debian:13.4@sha256:e2d08da6f42ef4b09b165d55528a12727aeed8240dc9edf888e3ec07e10ef9da
WORKDIR /app

RUN apt update \
  && apt install -y --no-install-recommends \
  git \
  ca-certificates \
  && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/build/main /usr/local/bin/semver-tagger

CMD ["/usr/local/bin/semver-tagger", "--help"]

LABEL image.name=markormesher/semver-tagger
LABEL image.registry=ghcr.io
LABEL org.opencontainers.image.description=""
LABEL org.opencontainers.image.documentation=""
LABEL org.opencontainers.image.title="semver-tagger"
LABEL org.opencontainers.image.url="https://github.com/markormesher/semver-tagger"
LABEL org.opencontainers.image.vendor=""
LABEL org.opencontainers.image.version=""
