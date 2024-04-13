# 🚧 toml2docs

`toml2docs` is a tool to generate a Markdown document from toml.

The tool will discover the comments from the toml file and generate a description for each field.  For example:

```toml
# The running mode of the server, can be `standalone` or `distributed`.
mode = "distributed"

# The default timezone of the server.
# The value should be a valid timezone name, such as `UTC`, `Local`, `Asia/Shanghai`, etc.
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

# The aarch64 processor options.
[[processors]]
thread_num = 4
arch = "aarch64"
```

And it will output the following Markdown file:

| Key | Type | Default | Descriptions |
| --- | -----| ------- | ----------- |
| `mode` | String | `distributed` | The running mode of the server, can be `standalone` or `distributed`. |
| `default_timezone` | String | `UTC` | The default timezone of the server. The value should be a valid timezone name, such as `UTC`, `Local`, `Asia/Shanghai`, etc. |
| `heartbeat` | -- | -- | The heartbeat options for server. |
| `heartbeat.interval` | String | `5s` | Interval for sending heartbeat task. |
| `heartbeat.retry_interval` | String | `5s` | Interval for retrying to send heartbeat task. |
| `postgres` | -- | -- | The PostgresSQL server options. |
| `postgres.enable` | Bool | `true` | Whether to enable the PostgresSQL server. |
| `postgres.addr` | String | `127.0.0.1:4003` | The address to bind the PostgresSQL server. |
| `postgres.runtime_size` | Integer | `2` | The runtime size of the PostgresSQL server. |
| `postgres.tls` | -- | -- | The PostgresSQL server TLS options. |
| `postgres.tls.mode` | String | `disable` | The mode of the PostgresSQL server TLS. |
| `postgres.tls.cert_path` | String | -- | The certificate path of the PostgresSQL server TLS. |
| `postgres.tls.key_path` | String | -- | The key path of the PostgresSQL server TLS. |
| `postgres.tls.watch` | Bool | `false` | Whether to watch the certificate changes of the PostgresSQL server TLS. |
| `meta_client` | -- | -- | The meta client options. |
| `meta_client.addrs` | Array | -- | The address of the meta servers. |
| `meta_client.timeout` | String | `3s` | The timeout for the meta client. |
| `processors[0].thread_num` | Integer | `2` | The x86_64 processor options. |
| `processors[0].arch` | String | `x86_64` | -- |
| `processors[1].thread_num` | Integer | `4` | The aarch64 processor options. |
| `processors[1].arch` | String | `aarch64` | -- |

## How to use

1. Build the project:

   ```
   make
   ```

   The `toml2docs` will be generated in `bin/`.

2. Provide the input file and generate the Markdown file(output to stdout by default)

   ```
   ./bin/toml2docs -i <input-file>
   ```
   
   You can also specify the output file:
   
   ```
   ./bin/toml2docs -i <input-file> -o <output-file>
   ```
