FROM golang AS builder

WORKDIR /home
COPY go.mod go.sum ./
RUN go mod download

COPY tailproxy.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o tailproxy .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /home/tailproxy .
ENTRYPOINT ["./tailproxy"]
