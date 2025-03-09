# build delob main app 
FROM golang:1.23-bullseye AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/cmd
RUN go build -o delob

WORKDIR /app/cli
RUN go build -o elo

# final stage 
FROM debian:buster-slim

ARG USERNAME
ARG PASSWORD

ENV BUILD_ENV=docker

RUN groupadd -r delobgroup && useradd -r -s /bin/false -g delobgroup delobuser

WORKDIR /home/delobuser

COPY --from=builder /app/cmd/delob .
COPY --from=builder /app/cli/elo .
RUN chmod +x delob && chown delobuser:delobgroup .

EXPOSE 5678

USER delobuser
RUN ./elo --add-user ${USERNAME} ${PASSWORD}

CMD ["./delob"]
