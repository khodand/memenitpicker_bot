FROM golang:1.22 AS builder

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 go build -installsuffix 'static' -o meme_bot ./cmd/main.go

FROM alpine:latest AS runtime

COPY --from=builder /app/meme_bot .

ENV CONFIG_PATH="/config/secret.yaml"
VOLUME [ "/config" ]

CMD ["./meme_bot"]
