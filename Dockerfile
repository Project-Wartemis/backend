FROM golang:1.22-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/backend ./cmd/backend

FROM alpine/curl:8.17.0
WORKDIR /app
ENV WARTEMIS_ENV=BUILD
COPY --from=builder /out/backend /app/backend

EXPOSE 80
ENTRYPOINT ["/app/backend"]
