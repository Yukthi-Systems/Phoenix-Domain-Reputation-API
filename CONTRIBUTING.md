# Contributing to Phoenix Domain Reputation

Thanks for taking the time to contribute! This document covers how to get
set up, the standards we hold changes to, and how to submit a pull
request.

## Code of Conduct

This project follows a [Code of Conduct](CODE_OF_CONDUCT.md). By
participating you agree to uphold it.

## Getting Started

```bash
git clone https://github.com/Yukthi-Systems/Phoenix-Domain-Reputation-API.git
cd Phoenix-Domain-Reputation-API
cp .env.example .env
go run ./cmd/server
```

Requirements: Go 1.25+.

## Reporting Bugs

Search [existing issues](https://github.com/Yukthi-Systems/Phoenix-Domain-Reputation-API/issues)
first. If nothing matches, open a new issue using the **Bug report**
template and include steps to reproduce, what you expected, and what
happened instead.

For security vulnerabilities, do **not** open a public issue — follow
[`SECURITY.md`](SECURITY.md) instead.

## Suggesting Features

Open an issue using the **Feature request** template. Explain the
problem you're trying to solve, not just the solution you have in mind —
it makes it much easier to evaluate alternatives.

## Development Workflow

1. Fork the repository and create a branch off `main`:
   `git checkout -b feat/short-description`.
2. Make your change.
3. Run the checks below and make sure they pass.
4. Commit with a clear, descriptive message.
5. Open a pull request using the PR template, describing what changed and
   why.

## Code Standards

This codebase intentionally stays close to the Go standard library and
avoids unnecessary abstraction. Please keep that spirit:

- **Standard library first.** `net/http`, `log/slog`, `sync/atomic` are
  preferred over third-party frameworks. A new dependency needs a clear,
  concrete reason.
- **No business logic in HTTP handlers.** Handlers translate HTTP ⇄ the
  domain packages (`internal/reputation`, `internal/ipfire`); they don't
  decide anything.
- **Doc comments on every exported identifier**, in standard Go style
  (comment starts with the identifier's name). Non-exported helpers should
  have a doc comment when their behavior isn't obvious from the name.
- **Errors are wrapped with context** (`fmt.Errorf("...: %w", err)`), not
  swallowed or panicked on for ordinary operational failures.
- **No premature abstraction.** Don't introduce an interface, generic, or
  config option for a single call site or a hypothetical future need.
- **Never mutate a published `reputation.Snapshot`.** The store is
  read-lock-free by design — new data always means building a new
  snapshot and calling `Store.Replace`, never editing the current one in
  place.

Before opening a PR, run:

```bash
gofmt -l .          # must print nothing
go vet ./...
go build ./...
go test ./... -race -cover
```

## Tests

New behavior should come with tests. Table-driven tests are preferred
where they fit. Tests must not depend on live network access or the real
IPFire service — use `httptest.Server` for HTTP interactions and fakes
(e.g. `internal/updater`'s `Downloader` interface) for dependencies that
would otherwise reach the network.

## Pull Requests

- Keep PRs focused — one logical change per PR is easier to review than a
  bundle of unrelated fixes.
- Update `README.md` and `.env.example` if you add or change
  configuration.
- Link the issue your PR addresses, if any.
- A maintainer will review and may ask for changes before merging.

## License

By contributing, you agree that your contributions will be licensed under
the project's [GNU General Public License v3.0](LICENSE).

## Questions?

Ask in the [Discord](https://discord.com/invite/2BS7Z4FhJ) or email
[connect@yukthi.com](mailto:connect@yukthi.com).
