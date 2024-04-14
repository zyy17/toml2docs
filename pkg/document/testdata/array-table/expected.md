| Key | Type | Default | Descriptions |
| --- | -----| ------- | ----------- |
| `[[engine]]` | -- | -- | The engine options. |
| `engine[0].name` | String | `sqlite` | -- |
| `engine[0].path` | String | `db.sqlite` | -- |
| `engine[1].name` | String | `postgres` | Use postgres. |
| `engine[1].path` | String | `/var/lib/postgresql/data` | -- |
| `engine[2].name` | String | `mysql` | -- |
| `engine[2].path` | String | `/var/lib/mysql/data` | -- |
| `[[proxy]]` | -- | -- | -- |
| `proxy[0].name` | String | `haproxy` | -- |
| `proxy[1].name` | String | `nginx` | -- |
