FROM golang:1.25-alpine AS builder

ENV CGO_ENABLED=1

RUN apk add --no-cache \
    # Important: required for go-sqlite3
    gcc \
    # Required for Alpine
    musl-dev

RUN mkdir /data && \
    chown -R 65532:65532 /data

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY main.go main.go
COPY cmd/ cmd/
COPY internal/ internal/
COPY pkg/ pkg/

RUN go build -ldflags='-s -w -extldflags "-static"' -o /quotobot main.go


FROM scratch

ENV USER=quotobot
ENV GROUP=quotobot

WORKDIR /

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /data /data

COPY --from=builder /quotobot /quotobot

USER 65532:65532

EXPOSE 8080

ENTRYPOINT ["/quotobot"]