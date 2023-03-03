FROM --platform=$BUILDPLATFORM cgr.dev/chainguard/go:1.20 AS build

WORKDIR /work
COPY go.mod go.sum ./
RUN go mod download

COPY src ./src
ARG TARGETOS TARGETARCH TARGETVARIANT
RUN \
    if [ "${TARGETARCH}" = "arm" ] && [ -n "${TARGETVARIANT}" ]; then \
      export GOARM="${TARGETVARIANT#v}"; \
    fi; \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} CGO_ENABLED=0 go build -v -o ./tailproxy ./src

FROM cgr.dev/chainguard/static:latest
ENV HOME /home/nonroot
COPY --from=build /work/tailproxy /tailproxy
ENTRYPOINT ["/tailproxy"]
