FROM docker.io/golang:1.26.1@sha256:e2ddb153f786ee6210bf8c40f7f35490b3ff7d38be70d1a0d358ba64225f6428 AS builder
WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY ./main.go ./main.go
COPY ./internal ./internal

RUN go build -o ./build/main ./main.go

# ---

FROM docker.io/debian:13.3@sha256:3615a749858a1cba49b408fb49c37093db813321355a9ab7c1f9f4836341e9db
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
