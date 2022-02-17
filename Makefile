
main:
	go build -o bin/proxy ./cmd/eventsproxy

client:
	go build -o bin/client examples/app.go

