#build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN go build -o /go/bin/app

#final stage
FROM alpine:latest
#RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/app /bin/app
ENV SERVER_PORT=80
ENTRYPOINT /bin/app
LABEL Name=restserver Version=0.0.1
EXPOSE 80 
