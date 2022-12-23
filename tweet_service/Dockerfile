FROM golang:latest as builder
WORKDIR /app

COPY ./tweet_service/go.mod .
COPY ./tweet_service/go.sum .
RUN go mod download
COPY ./tweet_service/ .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine
WORKDIR /root/
COPY --from=builder /app/main .
COPY  ./auth_service/certificates ./auth_service/certificates
EXPOSE 8003
CMD ["./main"]