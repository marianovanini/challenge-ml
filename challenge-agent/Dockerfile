FROM golang:1.18
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go mod vendor
RUN go build -o main /app/main.go
CMD ["/app/main"]