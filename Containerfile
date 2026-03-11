FROM docker.io/golang:1.26.1@sha256:cdebbd553e5ed852386e9772e429031467fa44ca3a06735e6beb005d615e623d AS builder
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
