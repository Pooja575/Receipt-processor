
FROM golang:1.17-alpine as builder


WORKDIR /app


COPY go.mod .
COPY go.sum .


RUN go mod download


COPY . .


RUN go build -o main .


FROM alpine:latest


WORKDIR /root/


COPY --from=builder /app/main .


EXPOSE 8080

# Command to run the executable
CMD ["./main"]
