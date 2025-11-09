FROM golang:1.25.0 AS builder

WORKDIR /usr/src/app

COPY . .

RUN go mod download
RUN go build -v -o tododiplomas

FROM ubuntu:latest

WORKDIR /usr/src/app

COPY --from=builder /usr/src/app/tododiplomas /usr/src/app/
COPY --from=builder /usr/src/app/web /usr/src/app/web
COPY --from=builder /usr/src/app/scheduler.db /usr/src/app/
COPY --from=builder /usr/src/app/.env /usr/src/app/

CMD ["/usr/src/app/tododiplomas"] 