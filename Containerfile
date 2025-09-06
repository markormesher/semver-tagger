FROM docker.io/golang:1.25.1@sha256:a5e935dbd8bc3a5ea24388e376388c9a69b40628b6788a81658a801abbec8f2e AS builder
WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY ./main.go ./main.go
COPY ./internal ./internal

RUN go build -o ./build/main ./main.go

# ---

FROM debian:13.0@sha256:6d87375016340817ac2391e670971725a9981cfc24e221c47734681ed0f6c0f5
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
