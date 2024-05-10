FROM golang:1.22.1-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o main cmd/main.go

FROM busybox

WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 3000

CMD ["./main"]
