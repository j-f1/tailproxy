FROM --platform=$BUILDPLATFORM cgr.dev/chainguard/go:1.20 AS build

WORKDIR /work
COPY go.mod go.sum ./
RUN go mod download

COPY src ./src
COPY tailproxy.go ./
ARG TARGETOS TARGETARCH TARGETVARIANT
RUN CGO_ENABLED=0 go build tailproxy.go

FROM cgr.dev/chainguard/static:latest
ENV HOME /home/nonroot
COPY --from=build /work/tailproxy /tailproxy
ENTRYPOINT ["/tailproxy"]
