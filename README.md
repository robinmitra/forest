[![Build Status](https://travis-ci.com/robinmitra/forest.svg?branch=master)](https://travis-ci.com/robinmitra/forest)

# Forest

*For the forest on your computer*

Forest is a command-line tool for analysing and browsing files and directories on the filesystem.

## Features

### Analyse files

The `analyse` command analyses files and directories at a given path, and summarises the following
metrics:
* Total number of files and directories
* Total disk space usage
* Top 5 file types (by occurrence and disk usage)
* Ability to ignore certain files and/or directories

#### Usage

```bash
forest analyse [path]
```

* `[path]` - Optional path from where to start analysing (defaults to current working directory).

##### Options

* `--include-dot-files`: Include hidden dot files in the analysis. These are excluded by default.
* `--format`: The output format of the summary. Options include `normal` (default) and `rainbow`.

### Browse files

The `browse` command presents a traversable tree of files and directories, with some metadata.

#### Usage

```bash
forest browse [path]
```

* `[path]` - Optional path from where to start browsing (defaults to current working directory).

## Development

### Building

The project requires Go version 0.11 and up, since it uses Go Modules for dependency management.

#### Building natively

Use standard Go toolchain to compile, run and install. As a refresher:

Compile and execute the application by running `go run main.go` - this does not leave the compiled 
binary behind, so its useful while developing.

Compile the application by running `go build` - you can then execute the generated binary directly.

Compile and install the application by running `go install` - this builds the binary and installs it
(in other words, moves it to the `bin` directory within your `GOPATH`), so that you can run it from
anywhere in your system.

#### Building using Docker

Compile the application in a Docker container by running: `make build`.

Install the binary by running: `make install`.

### Testing

#### Testing natively

Run tests for all packages by running `go test ./...`.

Run tests for a particular package (e.g. the `analyse` package): `go test ./cmd/analyse`.

#### Testing using Docker

Run tests for all package in a Docker container by running `make docker_test`.

### Releasing

Bump version by running `make bump` - this bumps the version in files, commits them, and creates a
Git tag.

Bump version and push changes by running `make release` - this bumps version like above and pushes
these changes to upstream.
