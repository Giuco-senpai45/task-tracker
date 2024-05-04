FROM golang:alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ./auth-service main.go



FROM alpine

COPY --from=builder /app/auth-service  /app/auth-service

EXPOSE 8081

CMD ["/app/auth-service"]
