| Key | Type | Default | Descriptions |
| --- | -----| ------- | ----------- |
| `mysql` | -- | -- | The MySQL configuration file is used to configure the MySQL connection. |
| `mysql.user` | String | `root` | -- |
| `mysql.tls` | -- | -- | Configure the MySQL connection to use TLS/SSL. |
| `mysql.tls.cert` | String | `/etc/mysql/certs/cert.pem` | -- |
| `mysql.tls.key` | String | `/etc/mysql/certs/key.pem` | -- |
| `mysql.tls.ca` | String | `/etc/mysql/certs/ca.pem` | -- |
| `mysql.tls.options` | -- | -- | The TLS options section is used to configure the TLS/SSL options for the MySQL connection. |
| `mysql.tls.options.mode` | String | `REQUIRED` | -- |
