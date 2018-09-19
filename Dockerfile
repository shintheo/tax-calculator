FROM golang:1.10
RUN mkdir -p /go/src/app
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...
CMD ["go", "run", "main.go"]
