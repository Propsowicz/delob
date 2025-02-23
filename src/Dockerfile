# build delob main app 
FROM golang:1.22.4-bullseye AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/cmd
RUN go build -o delob

# final stage 
FROM debian:buster-slim

ENV BUILD_ENV=docker

RUN groupadd -r delobgroup && useradd -r -s /bin/false -g delobgroup delobuser

WORKDIR /home/delobuser

COPY --from=builder /app/cmd/delob .
RUN chmod +x delob && chown delobuser:delobgroup .

EXPOSE 5678

USER delobuser

CMD ["./delob"]
