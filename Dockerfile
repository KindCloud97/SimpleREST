# Compile stage
FROM golang:1.16.3 AS build-env

WORKDIR /SimpleREST
ADD . .

RUN CGO_ENABLED=0 GOOS=linux go build -o bin/simple Simple.go

# Final stage
FROM alpine:latest AS deploy

EXPOSE 8000
EXPOSE 8080

WORKDIR /SimpleREST
COPY --from=build-env /SimpleREST/bin/simple .

CMD ["./simple"]