BINARY := kubernetes-live-images

.PHONY: build format test clean

build:
	go build -o $(BINARY) .

format:
	gofmt -w .

test:
	go test ./...

clean:
	rm -f $(BINARY)
