#builder
FROM golang:alpine as builder
WORKDIR /home
COPY . .
RUN go mod tidy
RUN go build -o client main.go

#final image
ENV GRPC_SERVER=server:80
ENV CALL=1000
ENV GODEBUG=
FROM golang:alpine
COPY --from=builder /home/client .
EXPOSE 50051
CMD ./client -addr=${GRPC_SERVER} -call=${CALL}