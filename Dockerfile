FROM golang:1.23.2-alpine AS builder
RUN apk update && \
  apk add --no-cache build-base pkgconfig sqlite-dev
ENV CGO_ENABLED=1 \
  GOOS=linux \
  GO111MODULE=on
WORKDIR /app
COPY . .
RUN go mod download && \
  go build -o create-token ./cmd/create-token && \
  go build -o server ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates openssl
RUN addgroup -S donezo && adduser -S donezo -G donezo
WORKDIR /home/donezo/app
COPY --from=builder /app/data/sql/migrations ./data/sql/migrations
COPY --from=builder /app/create-token .
COPY --from=builder /app/server .
RUN mkdir -p /data && \
  chown -R donezo:donezo /data && \
  chown -R donezo:donezo /home/donezo/app/
ENV PATH="/home/donezo/app/:${PATH}"
USER donezo
ENTRYPOINT [ "server" ]
