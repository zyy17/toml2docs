.PHONY: toml2docs
toml2docs:
	GO111MODULE=on CGO_ENABLED=0 go build -o bin/toml2docs cmd/toml2docs/main.go

.PHONY: docker
docker:
	docker build -t toml2docs:latest .

.PHONY: test
test:
	go test -v ./pkg/... -cover

.PHONY: clean
clean:
	rm -r bin

.PHONY: lint
lint: golangci-lint toml2docs ## Run lint.
	golangci-lint run -v ./...

.PHONY: golangci-lint
golangci-lint:
	@which golangci-lint || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.55.2
