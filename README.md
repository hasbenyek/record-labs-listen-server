# Record Labs Listen Together server

The WebSocket server behind **Listen Together** in Record Labs for Android and Windows.
Clients connect to `wss://<host>/ws`; messages use protobuf with optional gzip compression.

## Run with Docker

```bash
docker build -t recordlabs-listen .
docker run -p 8080:8080 recordlabs-listen
curl http://localhost:8080/health   # {"status":"ok"}
```

On Render, create a **Web Service** from this repository (runtime: Docker). Render sets `PORT` itself.

## Run locally

```bash
go build -o recordlabs-listen ./cmd/recordlabs-listen
PORT=8080 UA_POLICY_FILE=./ua_policy.json ./recordlabs-listen
```

## Configuration

| Variable | Default | Meaning |
| --- | --- | --- |
| `PORT` | `8080` | Port to listen on |
| `UA_POLICY_FILE` | `/app/ua_policy.json` in Docker | Which apps may host rooms (`allow`), and which are blocked |
| `DATABASE_FILE` | `recordlabs-listen.db` | Saved rooms and connection counts |
| `UA_ADMIN_TOKEN` | unset | Enables `GET /uas` (send `Authorization: Bearer <token>`) |

`ua_policy.json` allows every client whose User-Agent contains `recordlabs`
(Android sends `com.recordlabs.music`, Windows sends `RecordLabs-Windows/<version>`).

## Tests

```bash
go test ./...
```

## License

GPL-3.0 — see [LICENSE](LICENSE) and [NOTICE](NOTICE).
