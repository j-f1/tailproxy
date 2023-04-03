FROM cgr.dev/chainguard/go:1.20 AS build

# deps
WORKDIR /work
COPY go.mod go.sum ./
RUN go mod download

# build
COPY src ./src
COPY tailproxy.go ./
ARG TARGETOS TARGETARCH TARGETVARIANT
RUN CGO_ENABLED=0 go build -v tailproxy.go

# certs
FROM alpine AS certs
RUN apk --update add ca-certificates

# run
FROM scratch
ENV HOME /home/nonroot
COPY --from=build /work/tailproxy /tailproxy
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENTRYPOINT ["/tailproxy"]
