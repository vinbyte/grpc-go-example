compile-client:
	@mkdir -p ./client/helloworld
	@protoc --proto_path=proto/$* \
		--go_out=client/helloworld \
		--go_opt=paths=source_relative proto/helloworld.proto \
		--go-grpc_out=client/helloworld \
		--go-grpc_opt=paths=source_relative proto/helloworld.proto

compile-server:
	@mkdir -p ./server/helloworld
	@protoc --proto_path=proto/$* \
		--go_out=server/helloworld \
		--go_opt=paths=source_relative proto/helloworld.proto \
		--go-grpc_out=server/helloworld \
		--go-grpc_opt=paths=source_relative proto/helloworld.proto

compile: compile-client compile-server

run-server: compile-server
	@go run server/main.go

run-client: compile-client
	@GODEBUG=http2debug=2 go run client/main.go

compose-up: compile
	@docker-compose up

compose-down:
	@docker-compose down

compose-build:
	@docker rmi grpc-example-client grpc-example-server &> /dev/null || true
	@docker-compose build

k8s-up:
	@kubectl apply -f k8s/server-deployment.yaml,k8s/server-service.yaml
	@sleep 10
	@kubectl apply -f k8s/client-deployment.yaml
	@sleep 5
	@kubectl logs -l app=client -f 

k8s-up-linkerd:
	@linkerd inject ./k8s | kubectl apply -f -

k8s-down:
	@kubectl delete -f ./k8s


.PHONY: compile-client compile-server compile run-server run-client compose-up compose-down compose-build k8s-prepare k8s-up k8s-down k8s-up-linkerd