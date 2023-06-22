FROM golang:1.18-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o ./go-demo ./main.go


FROM alpine:latest AS runner
WORKDIR /app
COPY --from=builder /app/go-demo .
EXPOSE 8080
ENTRYPOINT ["./go-demo"]