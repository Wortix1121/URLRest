dockerbuild:
	docker build -t dev:local .

start:
	go run main.go

