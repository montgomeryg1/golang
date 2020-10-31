FROM golang:alpine AS builder

ENV GO111MODULE=auto \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64  

WORKDIR /app
ADD . /app
RUN go mod download
RUN go build -o blobreader .

FROM alpine:latest  

RUN apk --no-cache add ca-certificates
ENTRYPOINT ["/blobreader"]
COPY --from=builder /app .