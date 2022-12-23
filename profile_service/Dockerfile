FROM golang:latest as builder
WORKDIR /app

COPY ./profile_service/go.mod .
COPY ./profile_service/go.sum .
RUN go mod download
COPY ./profile_service/ .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine
WORKDIR /root/
COPY --from=builder /app/main .
COPY  ./auth_service/certificates ./auth_service/certificates
COPY ./profile_service/security/password_blacklist.txt ./profile_service/security/password_blacklist.txt
EXPOSE 8002
CMD ["./main"]