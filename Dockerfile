FROM golang AS builder

WORKDIR /home
COPY go.mod go.sum ./
RUN go mod download

COPY tailproxy.go ./
RUN go build -v -x -o tailproxy .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /home/tailproxy .
ENTRYPOINT ["./tailproxy"]
