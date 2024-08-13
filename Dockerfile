FROM golang:1.22.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o db_tools ./scripts/db_tools.go

FROM alpine:3.18

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/db_tools .
COPY --from=builder /app/scripts/migrations ./scripts/migrations

EXPOSE 8080

CMD ["./main"]
