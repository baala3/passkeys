# ---------- Stage 1: Build Frontend ----------
FROM node:20-alpine AS builder

WORKDIR /client
COPY client/ ./
RUN yarn install && yarn build

# ---------- Stage 2: Build Go Backend ----------
FROM golang:1.24-alpine AS backend

WORKDIR /server
COPY server/go.mod server/go.sum ./
RUN go mod download
COPY server/ ./
RUN go build -o main .
COPY start.sh ./
RUN chmod +x start.sh
CMD ["./start.sh"]
