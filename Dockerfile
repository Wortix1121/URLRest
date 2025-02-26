FROM golang:1.23.4-alpine AS builder

WORKDIR /app

RUN apk --no-cache add bash git make gcc gettext musl-dev

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 go build -o app .

FROM alpine:latest

COPY --from=builder /app .
COPY config/local.yaml /local.yaml

CMD ["./app"]