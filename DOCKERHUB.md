# Phoenix Domain Reputation

A fast, in-memory domain reputation API backed by the
[IPFire Domain Blocklist (DBL)](https://www.ipfire.org/dbl/how-to-use).
It downloads five IPFire category lists — gambling, malware, phishing,
pornography, and violence — on a schedule, builds an in-memory snapshot,
and serves reputation lookups fast enough for a mail/spam-processing hot
path. No database required.

**Source & documentation:** https://github.com/Yukthi-Systems/Phoenix-Domain-Reputation-API

## Quick Start

```bash
docker run --rm -p 8080:8080 \
  -e IPFIRE_UPDATE_INTERVAL=1h \
  -e IPFIRE_MALWARE_SCORE=5 \
  rjyspl/phoenix-domain-reputation-api
```

```bash
curl http://localhost:8080/health
curl http://localhost:8080/v1/reputation/domain/example.com
```

### With Docker Compose

```yaml
services:
  phoenix-domain-reputation:
    image: rjyspl/Phoenix-Domain-Reputation-API-api
    container_name: phoenix-domain-reputation
    env_file:
      - .env
    ports:
      - "127.0.0.1:8080:8080"
    restart: always
```

## Supported Tags

| Tag      | Description                              |
|----------|-------------------------------------------|
| `latest` | Latest build from the `main` branch       |
| `x.y.z`  | Immutable release build (semantic version) |

## Configuration

All configuration is via environment variables — no `.env` file is baked
into the image, so everything below must be supplied at container runtime
(`-e` flags, an `env_file`, or your orchestrator's secret/config
mechanism).

| Variable                    | Default | Description                                          |
|-------------------------------|---------|---------------------------------------------------------|
| `SERVER_PORT`                | `8080`  | HTTP listen port                                        |
| `IPFIRE_UPDATE_INTERVAL`     | `1h`    | Interval between background IPFire list downloads       |
| `IPFIRE_HTTP_TIMEOUT`        | `30s`   | Per-request timeout for IPFire downloads                |
| `IPFIRE_GAMBLING_SCORE`      | `2`     | Score added when a domain is on the gambling list       |
| `IPFIRE_MALWARE_SCORE`       | `5`     | Score added when a domain is on the malware list        |
| `IPFIRE_PHISHING_SCORE`      | `5`     | Score added when a domain is on the phishing list       |
| `IPFIRE_PORNOGRAPHY_SCORE`   | `2`     | Score added when a domain is on the pornography list    |
| `IPFIRE_VIOLENCE_SCORE`      | `2`     | Score added when a domain is on the violence list       |
| `LOG_LEVEL`                  | `INFO`  | `DEBUG`, `INFO`, `WARN`, or `ERROR`                     |

## Health Checks

- `GET /health` — liveness; always `200` once the process is up.
- `GET /ready` — readiness; `200` once the first IPFire snapshot has been
  published, `503` before that. Use this for your orchestrator's readiness
  probe, not `/health`.

## API

```bash
GET  /v1/reputation/domain/{domain}
GET  /v1/reputation/domain/{domain}/score
POST /v1/reputation/domains
```

Full API documentation, request examples, and architecture notes live in
the project [README](https://github.com/Yukthi-Systems/Phoenix-Domain-Reputation-API#readme).

## License

GNU General Public License v3.0. See the
[LICENSE](https://github.com/Yukthi-Systems/Phoenix-Domain-Reputation-API/blob/main/LICENSE)
file in the source repository.

## Support

- Issues & source: https://github.com/Yukthi-Systems/Phoenix-Domain-Reputation-API
- Discord: https://discord.com/invite/2BS7Z4FhJ
- Email: connect@yukthi.com
