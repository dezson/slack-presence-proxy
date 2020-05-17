.PHONY: build clean deploy test

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/handler handler/main.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose

test:
	sls invoke -f getPresence