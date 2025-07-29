# Contributing to Berth

Welcome, and thanks for your interest in contributing to **Berth**! 🎉  
We appreciate contributions of all kinds — code, documentation, bug reports, and ideas.

## 🛠 Project Setup

### Requirements

- Go `1.24+` installed and in your `PATH`
- [`make`](https://www.gnu.org/software/make/) installed
- [`golangci-lint`](https://golangci-lint.run/) installed (optional for local linting)

### Setup

```bash
git clone https://github.com/rluders/berth.git
cd berth
go mod tidy
````

## 🧪 Running Tests

```bash
make test
```

Runs all unit tests in the project using `go test ./...`.

## ⚙️ Building

```bash
make build
```

Compiles the binary into the `bin/` directory.

## 🧼 Linting

Berth uses [golangci-lint](https://golangci-lint.run/) in CI. You can run it locally with:

```bash
golangci-lint run
```

To install:

```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

## 🔄 Making a Pull Request

1. Fork the repository and create a branch from `main`
2. Make your changes
3. Run tests and linters locally
4. Push your branch and open a pull request

Please ensure that:

* Code is well-commented
* You’ve run `go mod tidy` to keep dependencies clean
* Lint passes: `golangci-lint run`
* Tests pass: `make test`

## 🚀 Releasing a New Version

Releases are automated using [GoReleaser](https://goreleaser.com) and GitHub Actions.

### Steps:

1. Make sure all changes are merged into `main`
2. Bump the version using [semantic versioning](https://semver.org/), e.g.:

```bash
git tag v1.0.0
git push origin v1.0.0
```

3. GitHub Actions will:

    * Build the binaries for multiple platforms
    * Generate checksums
    * Create a GitHub Release with attached binaries

You can view release results in the **Actions** tab and under **Releases** in GitHub.

## 💬 Need Help?

Open an issue or start a discussion. We’re happy to help!

