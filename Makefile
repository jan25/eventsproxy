
main:
	go build -o bin/proxy main.go

client:
	go build -o bin/client examples/app.go

