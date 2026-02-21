# syntax=docker/dockerfile:1.7

FROM golang:1.25-bookworm AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
	go mod download

COPY . .

ARG TARGETOS=linux
ARG TARGETARCH=amd64
RUN --mount=type=cache,target=/go/pkg/mod \
	--mount=type=cache,target=/root/.cache/go-build \
	CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
		go build \
			-ldflags "-s -w" \
			-o build/hapm \
			hapm.go

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /tmp
COPY --from=builder /src/build/hapm /hapm

ENTRYPOINT ["/hapm"]
