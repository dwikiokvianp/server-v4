# Build Stage
FROM golang AS build

WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o binary

# Runtime Stage
FROM golang AS runtime

WORKDIR /app
COPY --from=build /app/binary .
COPY --from=build /app/.env .
EXPOSE 8080

CMD ["./binary"]