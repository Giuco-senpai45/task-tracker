FROM golang:alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ./task-service main.go



FROM alpine

COPY --from=builder /app/task-service  /app/task-service

EXPOSE 8080

CMD ["/app/task-service"]
