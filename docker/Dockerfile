FROM golang:1.21.8-alpine3.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -C cmd/gophermart -a -installsuffix cgo -o main .

FROM alpine:latest  
WORKDIR /root/
COPY --from=builder /app/cmd/gophermart/main .
EXPOSE 8080
CMD ["./main"] 