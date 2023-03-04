FROM --platform=$BUILDPLATFORM cgr.dev/chainguard/go:1.20 AS build

WORKDIR /work
COPY go.mod go.sum ./
RUN go mod download

COPY src ./src
COPY tailproxy.go ./
ARG TARGETOS TARGETARCH TARGETVARIANT
RUN CGO_ENABLED=0 go build tailproxy.go

FROM cgr.dev/chainguard/static:latest
FROM scratch
ENV HOME /home/nonroot
ENV TAILPROXY_DATA_DIR /home/nonroot/data
COPY --from=build /work/tailproxy /tailproxy
ENTRYPOINT ["/tailproxy"]
