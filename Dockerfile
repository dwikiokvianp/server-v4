# Build Stage
FROM golang AS build

WORKDIR /app
COPY . .
RUN go run seed/seed.go && go run seed/seed_regions.go
RUN go mod tidy
RUN go build -o binary

# Runtime Stage
FROM ubuntu:latest AS runtime

WORKDIR /app
COPY --from=build /app/binary .
COPY --from=build /app/.env .
EXPOSE 8081

RUN cd /app

CMD ["./binary"]