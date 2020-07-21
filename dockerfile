FROM golang:1.8.3 as builder
WORKDIR /app
COPY upload.go  .
RUN go get -u github.com/gorilla/mux
RUN go get -u github.com/rs/cors
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o upload .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/upload .
CMD ["./upload"]