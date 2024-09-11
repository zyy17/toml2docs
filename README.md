# toml2docs

`toml2docs` is a tool that generates a Markdown document from toml.

The tool will discover the comments from the toml file and generate a description for each field.  For example:

```toml
# The running mode of the server, can be `standalone` or `distributed`.
mode = "distributed"

# The default timezone of the server.
# The value should be a valid timezone name, such as `UTC`, `Local`, `Asia/Shanghai`, etc.
# @toml2docs:none-default
default_timezone = "UTC"

# The heartbeat options for server.
[heartbeat]
# Interval for sending heartbeat task.
interval = "5s"
# Interval for retrying to send heartbeat task.
retry_interval = "5s"

# The PostgresSQL server options.
[postgres]
# Whether to enable the PostgresSQL server.
enable = true
# The address to bind the PostgresSQL server.
addr = "127.0.0.1:4003"
# The runtime size of the PostgresSQL server.
runtime_size = 2

# The PostgresSQL server TLS options.
[postgres.tls]
# The mode of the PostgresSQL server TLS.
mode = "disable"
# The certificate path of the PostgresSQL server TLS.
cert_path = ""
# The key path of the PostgresSQL server TLS.
key_path = ""
# Whether to watch the certificate changes of the PostgresSQL server TLS.
watch = false

# The meta client options.
[meta_client]
# The address of the meta servers.
addrs = ["127.0.0.1:3002"]
# The timeout for the meta client.
timeout = "3s"

# The x86_64 processor options.
[[processors]]
thread_num = 2
arch = "x86_64"
```

And it will output the following Markdown file:

| Key | Type | Default | Descriptions |
| --- | -----| ------- | ----------- |
| `mode` | String | `distributed` | The running mode of the server, can be `standalone` or `distributed`. |
| `default_timezone` | String | `None` | The default timezone of the server.<br/>The value should be a valid timezone name, such as `UTC`, `Local`, `Asia/Shanghai`, etc. |
| `heartbeat` | -- | -- | The heartbeat options for server. |
| `heartbeat.interval` | String | `5s` | Interval for sending heartbeat task. |
| `heartbeat.retry_interval` | String | `5s` | Interval for retrying to send heartbeat task. |
| `postgres` | -- | -- | The PostgresSQL server options. |
| `postgres.enable` | Bool | `true` | Whether to enable the PostgresSQL server. |
| `postgres.addr` | String | `127.0.0.1:4003` | The address to bind the PostgresSQL server. |
| `postgres.runtime_size` | Integer | `2` | The runtime size of the PostgresSQL server. |
| `postgres.tls` | -- | -- | The PostgresSQL server TLS options. |
| `postgres.tls.mode` | String | `disable` | The mode of the PostgresSQL server TLS. |
| `postgres.tls.cert_path` | String | `""` | The certificate path of the PostgresSQL server TLS. |
| `postgres.tls.key_path` | String | `""` | The key path of the PostgresSQL server TLS. |
| `postgres.tls.watch` | Bool | `false` | Whether to watch the certificate changes of the PostgresSQL server TLS. |
| `meta_client` | -- | -- | The meta client options. |
| `meta_client.addrs` | Array | -- | The address of the meta servers. |
| `meta_client.timeout` | String | `3s` | The timeout for the meta client. |
| `[[processors]]` | -- | -- | The x86_64 processor options. |
| `processors.thread_num` | Integer | `2` | The thread number of the processor. |
| `processors.arch` | String | `x86_64` | The arch of the processor. |

## 🚀 Quick Start

Use the docker image to experience the `toml2docs`:

```console
docker run --rm \
  -v $(pwd)/pkg/document/testdata/basic:/data \
  toml2docs/toml2docs:latest \
  -i /data/input.toml
```

- `-i/--input-file`: Specifies the input toml file.
- `-o/--output-file`: Specifies the output markdown file. If not specified, the output will be printed to stdout.
- `-t/--template-file`: Specifies the template file.
- `-p/--docs-comment-prefix`: Specifies the prefix of the comments that will be used as the description of the field. Default is `#`.

## Template

toml2docs supports generating docs from a template file. For example, you can provide a template file like this:

```markdown
# Configurations

## Datanode

{{ toml2docs "./testdata/templates/datanode.toml" }}

## Frontend

{{ toml2docs "./testdata/templates/frontend.toml" }}
```

Then run with the `-t/--template-file` option and the `{{ toml2docs <file> }}` section will be replaced with the generated docs.

## Metadata in Comments

You can add some special metadata in the comments to control the output of the field:

- `@toml2docs:none-default=<custom-default-value>`

  There is no default value for the field. It's useful for the scenario that supports `None` value, for example, `None` in the Rust Option structure.

  By default, the default value is `Unset` without code reference. You can also specify a custom default value by using `@toml2docs:none-default=<custom-default-value>`.

- `#+`

  `toml2docs` will treat the comment as a toml field and use it to generate the document.  

## Development

### Build and Test

Run the following commands to build the `toml2docs`:

```console
make
```

The `toml2docs` will be generated in `bin/`.

- `make test`: Run all the unit tests.
- `make lint`: Run the lint.
