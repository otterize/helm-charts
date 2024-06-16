FROM golang:1.22.1-alpine
WORKDIR /app
ADD tests/go.mod .
ADD tests/go.sum .
ADD tests/databases/postgres/client/client.go .
RUN go build -o main .

RUN ls -lh
#COPY ./main /app/main
CMD ["/app/main"]