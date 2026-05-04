# Contributing

Contributions are welcome. Please open an issue before starting significant work so we can discuss the approach first.

## Getting Started

```
git clone https://github.com/smythg4/go-ynab
cd go-ynab
go mod download
```

Run the test suite:

```
go test ./...
```

## Guidelines

**Code style** — All code must be formatted with `gofmt`. The CI pipeline will reject unformatted code. Run `gofmt -w .` before pushing.

**Tests** — New endpoints and bug fixes must include tests. The existing tests in `ynab/*_test.go` show the patterns to follow. Write table-driven tests where a function has multiple interesting cases.

**Doc comments** — Exported types and functions require a doc comment. Unexported helpers and test functions do not.

**Commits** — Keep commits focused. One logical change per commit makes review and bisection easier.

**No generated code** — This library is hand-written to stay idiomatic. Do not use OpenAPI generators or similar tools.

## Adding an Endpoint

1. Add the method to the appropriate `ynab/*.go` file, following the existing pattern for the HTTP verb
2. Add the corresponding test in `ynab/*_test.go`
3. Add a row to the API Coverage table in `README.md`

## Reporting Bugs

Open a GitHub issue with:
- The method you called
- The response or error you received
- What you expected to happen

If the bug involves sensitive account data, describe the shape of the response rather than the values.

## Versioning

This project follows [Semantic Versioning](https://semver.org). Breaking changes to public types or method signatures require a major version bump.

## License

By contributing you agree that your contributions will be licensed under the same [MIT License](LICENSE) as the project.
