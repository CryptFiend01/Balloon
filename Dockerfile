FROM golang:1.20.6 AS builder

WORKDIR /app

COPY . /app/

ENV GO111MODULE=on
ENV GOPROXY="https://goproxy.cn"

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o Balloon

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/Balloon /app/
COPY config.yaml /app/
COPY logger.json /app/
RUN chmod +x /app/Balloon

CMD [ "/app/Balloon" ]
