# Configurations

## Datanode

| Key | Type | Default | Descriptions |
| --- | -----| ------- | ----------- |
| `datanode` | -- | -- | The datanode options.<br/>Datanode is the data service of the cluster. |
| `datanode.engine` | String | `mito` | The engine to use for the datanode.<br/>- `m1`: The M1 engine.<br/>- `m2`: The M2 engine. |


## Frontend

| Key | Type | Default | Descriptions |
| --- | -----| ------- | ----------- |
| `frontend` | -- | -- | The frontend options.<br/>Frontend options are used to configure the frontend of the application. |
| `frontend.addr` | String | `127.0.0.1:8080` | The server address to bind to. |
