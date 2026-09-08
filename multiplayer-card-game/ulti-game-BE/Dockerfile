FROM golang:1.24-alpine AS builder

WORKDIR /backend

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /ulti-game-BE

FROM alpine:latest

WORKDIR /app

COPY --from=builder /ulti-game-BE /app/ulti-game-BE

EXPOSE 3000

CMD [ "/app/ulti-game-BE" ] 