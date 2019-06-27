FROM golang:alpine

RUN apk update; apk upgrade
RUN apk add git gcc sqlite sqlite-dev g++ make

# Only re-fetch the deps if the go.mod file changes
RUN mkdir /dep 
ADD go.mod /dep/go.mod 
RUN cd /dep; export CGO_ENABLED=1; export CC=gcc; go get; rm -r /dep

RUN mkdir /app 
ADD . /app/
WORKDIR /app 

RUN export CGO_ENABLED=1; export CC=gcc; go get
RUN export CGO_ENABLED=1; export CC=gcc; go test
RUN export CGO_ENABLED=1; export CC=gcc; go build -o main .

RUN rm -r /go

CMD ["./main"]
