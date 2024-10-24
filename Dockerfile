FROM golang:1.23.2-alpine AS builder
RUN apk update && \
  apk add --no-cache build-base pkgconfig sqlite-dev
ENV CGO_ENABLED=1 \
  GOOS=linux \
  GO111MODULE=on
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o create-token ./cmd/create-token
RUN go build -o server ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates openssl
RUN addgroup -S donezo && adduser -S donezo -G donezo
WORKDIR /home/donezo/app
COPY --from=builder /app/data/sql ./data/sql
COPY --from=builder /app/create-token .
COPY --from=builder /app/server .
RUN mkdir -p /etc/donezo && \
  mkdir -p /data && \
  echo "Port: 8000" > /etc/donezo/config.yaml && \
  echo "Database: /data/db.sqlite" >> /etc/donezo/config.yaml && \
  echo "JWTSecret: $(openssl rand -base64 32)" >> /etc/donezo/config.yaml && \
  chown -R donezo:donezo /etc/donezo && \
  chown -R donezo:donezo /data && \
  chown -R donezo:donezo /home/donezo/app/
ENV PATH="/home/donezo/app/:${PATH}"
USER donezo
EXPOSE 8000
ENTRYPOINT [ "server" ]
