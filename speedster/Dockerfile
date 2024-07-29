FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go


FROM gcr.io/distroless/static:nonroot

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080

USER nonroot
CMD ["./main"]
