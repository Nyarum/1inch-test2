
build/service:
	cd cmd/service && go build -o service main.go

run/service: build/service
	cd cmd/service && ./service 

build/client:
	cd cmd/client && go build -o client main.go

run/client: build/client
	cd cmd/client && ./client 