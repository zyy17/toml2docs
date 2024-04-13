FROM golang:1.21 as builder

ENV LANG en_US.utf8
WORKDIR /toml2docs

# Build the project.
COPY . .

RUN make toml2docs

# Export the binary to the clean image.
FROM alpine:3.14 as base
COPY --from=builder /toml2docs/bin/toml2docs /bin/toml2docs
ENTRYPOINT [ "toml2docs" ]
