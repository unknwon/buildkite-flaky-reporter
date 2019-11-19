NAME = buildkite-flaky-reporter
LDFLAGS += -X "main.Version=$(shell git rev-parse HEAD)"

build:
	go build -v -o $(NAME)

web: build
	./$(NAME)

release:
	env GOOS=linux GOARCH=amd64 go build -ldflags '$(LDFLAGS)' -o $(NAME); tar czf linux_amd64.tar.gz $(NAME)

clean:
	go clean
	rm -f *.tar.gz
