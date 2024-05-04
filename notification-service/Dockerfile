FROM golang:alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ./notification-service main.go


FROM alpine

COPY --from=builder /app/notification-service  /app/notification-service

CMD ["/app/notification-service"]
