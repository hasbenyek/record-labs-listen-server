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
| `CLIENT_KEY` | unset | Secret shared with the official app builds (see below) |
| `REQUIRE_CLIENT_KEY` | `false` | `true`: refuse every client without a valid key |

`ua_policy.json` allows every client whose User-Agent contains `recordlabs`
(Android sends `com.recordlabs.music`, Windows sends `RecordLabs-Windows/<version>`).

## Official apps only

The source code of the apps is public, so a User-Agent alone cannot tell an official
build from a fork. Official builds therefore carry a secret key that is not in any
repository (the apps read it from `local.properties` or a CI secret at build time).
On every connection they send

```
X-RecordLabs-Auth: v1.<unix seconds>.<hex HMAC-SHA256(key, "recordlabs-listen.v1.<unix seconds>")>
```

and the server checks it against `CLIENT_KEY` (tokens are valid for 12 hours either side
of the server clock).

- `CLIENT_KEY` unset: nothing changes; the User-Agent policy decides.
- `CLIENT_KEY` set: a valid token is always allowed; other clients still follow the User-Agent policy.
- `CLIENT_KEY` set and `REQUIRE_CLIENT_KEY=true`: clients without a valid token are refused with
  "This Listen Together server only works with the official Record Labs app."

Turn `REQUIRE_CLIENT_KEY` on once the app versions that send the key are out.
A key inside a distributed app can still be dug out by someone determined; if that happens,
change `CLIENT_KEY` (and the apps' key) and release new builds.

## Tests

```bash
go test ./...
```

## License

GPL-3.0 — see [LICENSE](LICENSE) and [NOTICE](NOTICE).
