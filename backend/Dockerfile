FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git make

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make build

FROM alpine:latest AS api

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /app/bin/api .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./api"]

FROM alpine:latest AS worker

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /app/bin/worker .

CMD ["./worker"]

FROM alpine:latest AS migrator

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/bin/migrator .
COPY --from=builder /app/migrations ./migrations

ENTRYPOINT ["./migrator"]
CMD ["up"]