.PHONY: build clean deploy

build:
#	env GOOS=linux go build -ldflags="-s -w" -o bin/bettercommit cmd/bettercommit/main.go
	go build -ldflags="-s -w" -o bin/bettercommit cmd/bettercommit/main.go

clean:
	rm -rf ./bin ./vendor

deploy: clean build
	sls deploy --verbose

