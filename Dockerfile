FROM golang:1.24-alpine AS builder
WORKDIR /src

COPY .env makefile go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd/
COPY docs/ ./docs/
COPY internal/ ./internal/
COPY pkg/ ./pkg/
COPY public/ ./public/

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/tasklybe ./cmd

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /out/tasklybe ./tasklybe
COPY --from=builder /src/.env ./.env
COPY --from=builder /src/makefile ./makefile
COPY --from=builder /src/go.mod ./go.mod
COPY --from=builder /src/go.sum ./go.sum
COPY --from=builder /src/docs ./docs
COPY --from=builder /src/public ./public
EXPOSE 4002
CMD ["./tasklybe"]
