FROM golang:1.19 AS builder
WORKDIR /usr/src/app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o auth ./cmd/main.go

FROM scratch
COPY --from=builder /usr/src/app/auth ./auth
CMD ["./auth"]
