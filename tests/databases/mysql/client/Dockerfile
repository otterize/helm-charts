FROM golang:1.23.1-alpine
WORKDIR /app
ADD tests/go.mod .
ADD tests/go.sum .
ADD tests/databases/mysql/client/client.go .
RUN go build -o main .

RUN ls -lh
#COPY ./main /app/main
CMD ["/app/main"]