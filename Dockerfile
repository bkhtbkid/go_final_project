FROM golang:1.25.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/scheduler .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/scheduler /app/scheduler

ENV TODO_PORT=7540
ENV TODO_DBFILE=scheduler.db

CMD [ "/app/scheduler" ]
