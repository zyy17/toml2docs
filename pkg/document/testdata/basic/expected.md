| Key | Type | Default | Descriptions |
| --- | -----| ------- | ----------- |
| `mode` | String | `distributed` | The running mode of the server, can be `standalone` or `distributed`. |
| `default_timezone` | String | `UTC` | The default timezone of the server. |
| `accept_methods` | Array | `None` | Accept methods |
| `heartbeat` | -- | -- | The heartbeat options for server. |
| `heartbeat.interval` | String | `5s` | Interval for sending heartbeat task. |
| `heartbeat.retry_interval` | String | `5s` | Interval for retrying to send heartbeat task. |
| `http` | -- | -- | The http server options. |
| `http.addr` | String | `127.0.0.1:4000` | The address to bind the http server. |
| `http.timeout` | String | `30s` | The timeout for the http server. |
| `http.body_limit` | String | `64MB` | The body size limit for the http server. |
| `grpc` | -- | -- | gRPC server options, see `standalone.example.toml`. |
| `grpc.addr` | String | `127.0.0.1:4001` | -- |
| `grpc.runtime_size` | Integer | `8` | -- |
| `mysql` | -- | -- | The MySQL server options. |
| `mysql.enable` | Bool | `true` | Whether to enable the MySQL server. |
| `mysql.addr` | String | `127.0.0.1:4002` | The address to bind the MySQL server. |
| `mysql.runtime_size` | Integer | `2` | The runtime size of the MySQL server. |
| `mysql.tls` | -- | -- | The TLS options for MySQL server. |
| `mysql.tls.mode` | String | `disable` | The mode of the MySQL server TLS. |
| `mysql.tls.cert_path` | String | `""` | The certificate path of the MySQL server TLS. |
| `mysql.tls.key_path` | String | `""` | The key path of the MySQL server TLS. |
| `mysql.tls.watch` | Bool | `false` | Whether to watch the certificate changes of the MySQL server TLS. |
| `postgres` | -- | -- | The PostgresSQL server options. |
| `postgres.enable` | Bool | `true` | Whether to enable the PostgresSQL server. |
| `postgres.addr` | String | `127.0.0.1:4003` | The address to bind the PostgresSQL server. |
| `postgres.runtime_size` | Integer | `2` | The runtime size of the PostgresSQL server. |
| `postgres.tls` | -- | -- | The PostgresSQL server TLS options, see `standalone.example.toml`. |
| `postgres.tls.mode` | String | `disable` | The mode of the PostgresSQL server TLS. |
| `postgres.tls.cert_path` | String | `""` | The certificate path of the PostgresSQL server TLS. |
| `postgres.tls.key_path` | String | `""` | The key path of the PostgresSQL server TLS. |
| `postgres.tls.watch` | Bool | `false` | Whether to watch the certificate changes of the PostgresSQL server TLS. |
| `prom_store` | -- | -- | The prometheus service options. |
| `prom_store.enable` | Bool | `true` | Whether to enable the prometheus service. |
| `prom_store.with_metric_engine` | Bool | `true` | Whether to use the metric engine for the prometheus service. |
| `meta_client` | -- | -- | The meta client options. |
| `meta_client.addrs` | Array | -- | The address of the meta servers. |
| `meta_client.timeout` | String | `3s` | The timeout for the meta client. |
| `datanode` | -- | -- | The datanode options. |
| `datanode.client` | -- | -- | The client options for the datanode. |
| `datanode.client.timeout` | String | `10s` | The timeout for the datanode client. |
| `datanode.client.connect_timeout` | String | `10s` | The connect timeout for the datanode client. |
| `datanode.client.tcp_nodelay` | Bool | `true` | Whether to enable `tcp_nodelay` for the datanode client. |
| `wal` | -- | -- | The wal options. |
| `wal.provider` | String | `local` | The provider of the wal.<br/>- `local`: use the local wal.<br/>- `remote`: use the remote wal. |
