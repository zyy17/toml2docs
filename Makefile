.PHONY: bin
bin:
	go build -o bin/toml2docs cmd/toml2docs/main.go

.PHONY: test
test:
	go test -v ./pkg/... -cover

.PHONY: clean
clean:
	rm -r bin
