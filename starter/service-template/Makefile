.PHONY: build run test clean docker

build:
	go build -o bin/service .

run:
	go run .

test:
	go test ./... -v -race

clean:
	rm -rf bin/

docker:
	docker build -t service-template .
