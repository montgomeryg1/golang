FROM golang:alpine AS builder

ENV GO111MODULE=auto \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64  

WORKDIR /app
ADD . /app
RUN go mod download
RUN go build -o grafana-dashboards ./cmd/web

FROM alpine:latest  

ENTRYPOINT ["/grafana-dashboards"]
COPY --from=builder /app/grafana-dashboards .