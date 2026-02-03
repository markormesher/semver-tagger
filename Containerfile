FROM docker.io/golang:1.25.6@sha256:0c87ea6991c06552ca5f516e3aeb434056bac3b674f32f612691692668e57074 AS builder
WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY ./main.go ./main.go
COPY ./internal ./internal

RUN go build -o ./build/main ./main.go

# ---

FROM debian:13.3@sha256:5cf544fad978371b3df255b61e209b373583cb88b733475c86e49faa15ac2104
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