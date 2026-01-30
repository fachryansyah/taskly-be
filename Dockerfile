FROM golang:1.24-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/ cmd/
COPY pkg/ pkg/
COPY internal/ internal/
COPY go.mod go.mod
COPY go.sum go.sum
RUN go build -o /tasklybe ./cmd/tasklybe

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /
COPY --from=builder /tasklybe /tasklybe
EXPOSE 4002
CMD ["/tasklybe"]