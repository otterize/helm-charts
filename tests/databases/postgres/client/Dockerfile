FROM golang:1.22.1-alpine AS build
WORKDIR /app
ADD ../../../go.mod .
ADD ../../../go.sum .
ADD client.go .
RUN go build -o main .

COPY ./main /app/main
CMD ["/app/main"]