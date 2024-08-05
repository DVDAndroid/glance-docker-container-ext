FROM golang:1.22-alpine AS builder

WORKDIR /build
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build --trimpath -o glance-docker-container-ext .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /build/glance-docker-container-ext .
COPY widget.gohtml /app/widget.gohtml

CMD ["./glance-docker-container-ext"]

