.PHONY: bin
bin:
	go build -o bin/toml2docs main.go

clean:
	rm -r bin
