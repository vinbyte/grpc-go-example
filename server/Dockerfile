#builder
FROM golang:alpine as builder
WORKDIR /home
COPY . .
RUN go mod tidy
RUN go build -o server main.go

#final image
FROM alpine
COPY --from=builder /home/server .
EXPOSE 50051
CMD ./server