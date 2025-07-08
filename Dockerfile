# ---------- Stage 1: Build Frontend ----------
FROM node:20-alpine AS builder

WORKDIR /client
COPY client/ ./
RUN yarn install && yarn build


