FROM docker.io/golang:1.25.3@sha256:8c945d3e25320e771326dafc6fb72ecae5f87b0f29328cbbd87c4dff506c9135 AS builder
WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY ./main.go ./main.go
COPY ./internal ./internal

RUN go build -o ./build/main ./main.go

# ---

FROM debian:13.1@sha256:72547dd722cd005a8c2aa2079af9ca0ee93aad8e589689135feaed60b0a8c08d
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
