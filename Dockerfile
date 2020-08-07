FROM golang:1.14.6

RUN mkdir /app 
ADD . /app/
WORKDIR /app 
RUN go build -o main .
CMD ["./main"]