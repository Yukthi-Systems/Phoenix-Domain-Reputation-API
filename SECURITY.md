# Security Policy

## Reporting a Vulnerability

If you discover a security vulnerability in Phoenix Domain Reputation,
please **do not** open a public GitHub issue. Instead, report it
privately by emailing:

**[connect@yukthi.com](mailto:connect@yukthi.com)**

Please include:

- A description of the vulnerability and its potential impact
- Steps to reproduce it (a minimal proof of concept helps a lot)
- The version/commit you tested against

We will acknowledge your report as soon as possible and keep you updated
as we investigate and address the issue. Once a fix is available, we will
coordinate disclosure with you and credit you in the release notes, unless
you'd prefer to remain anonymous.

## Supported Versions

This project is pre-1.0 and moves quickly. Security fixes are applied to
the `main` branch and the latest published Docker image tag. If you are
running an older build, please upgrade before reporting an issue that may
already be fixed.

## Scope

This service downloads plaintext domain lists from
[IPFire's DBL](https://www.ipfire.org/dbl/how-to-use) and serves reputation
lookups over HTTP. Areas of particular interest for security review
include:

- Input handling for the `{domain}` path parameter and the batch request
  body
- The `X-API-Key` authentication middleware (`internal/httpapi/middleware.go`)
  guarding the `/v1/*` routes
- HTTP client behavior when downloading IPFire lists (timeouts, response
  size limits, TLS)
- Anything that could allow a malicious IPFire response to affect more
  than the reputation data itself
