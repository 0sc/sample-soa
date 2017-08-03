FROM golang:alpine

ADD . .
RUN go install 

CMD ["go", "run", "*.go"]