# Buffman

<p align="center">
  <img src="docs/buffman.png" alt="Buffman Logo" width="400" />
</p>

**Buffman** is a CLI tool that wraps around the `flatc` compiler. It simplifies converting `.proto` files to `.fbs`, and generates code in multiple languages using a declarative YAML config (`buffman.yml`).

It currently supports two plugin types:

- `flatbuffers` — multi-language code generation
- `nanobuffers` — minimal and ultra-fast C-only serialization

> [!IMPORTANT]  
> This project is under active development. APIs, configurations, and features may change without notice. Use with caution in production environments.

- [Buffman](#buffman)
  - [Installation](#installation)
  - [Quickstart](#quickstart)
  - [Commands](#commands)
  - [Configuration](#configuration)
  - [Examples](#examples)
    - [Minimal example](#minimal-example)
    - [Multi-language production example](#multi-language-production-example)
  - [License](#license)

## Installation

You can install Buffman in four ways:

1. **Homebrew** (macOS / Linux)

   ```bash
   brew install the-protobuf-project/tap/buffman
   ```

2. **curl installer**

   ```bash
   curl -sSL https://raw.githubusercontent.com/the-protobuf-project/buffman/main/scripts/install.sh | bash
   ```

   Visit the [Releases page](https://github.com/the-protobuf-project/buffman/releases) to download a specific version.

3. **Build from Source**

   ```bash
   git clone https://github.com/the-protobuf-project/buffman.git
   cd buffman
   go build -o buffman .
   ```

4. **Docker Image**

   ```bash
   docker pull ghcr.io/the-protobuf-project/buffman:latest
   ```

> [!TIP]
> Add the binary to your `PATH` for convenient use from anywhere.

## Quickstart

Buffman requires a YAML configuration file and **does not** auto-detect it.  
You **must specify the file explicitly** using the `-f` flag.

Here's a minimal example config (`buffman.yml`):

```yaml
version: v1

inputs:
  - name: source
    path: "./proto"

  - name: googleprotobuf
    remote: https://github.com/protocolbuffers/protobuf
    commit: <commit-hash>

plugins:
  - name: flatbuffers
    out: "./fbs"
    languages:
      - language: go
        out: "./generated/go"
        opt:
          - go_package=github.com/username/project/fb

  - name: nanobuffers
    out: "./nano"
```

Then run:

```bash
buffman generate -f ./buffman.yml
```

Or use Docker:

```bash
docker run --rm \
    -v $(pwd):/buffman \
    -w /buffman \
    ghcr.io/the-protobuf-project/buffman:latest generate -f /buffman/buffman.yml
```

> [!NOTE]
> When using Docker, all paths in your `buffman.yml` must be **relative to `/buffman`**, since that's where your local project is mounted in the container.

You can use any filename and location for the config—just update the path with `-f`.

## Commands

| Command             | Description                                                                                               |
|---------------------|-----------------------------------------------------------------------------------------------------------|
| `buffman generate`  | Generates code as defined in your config file. Use the `-f` flag to specify the config path.              |
| `buffman convert`   | Converts `.proto` files to `.fbs` files using your config. [Learn more](docs/convert.md)                  |

## Configuration

Buffman uses a YAML configuration file (`buffman.yml`) to define your input sources, output directories, plugins, and language-specific options.

### Structure

```yaml
version: v1

inputs:
  - name: source
    path: "./proto"

  # Optional external repositories
  # - name: googleprotobuf
  #   remote: https://github.com/protocolbuffers/protobuf
  #   commit: <commit-hash>

plugins:
  - name: flatbuffers
    out: "./fbs"
    languages:
      - language: cpp
        out: "./generated/cpp"

      - language: go
        out: "./generated/go"
        opt:
          - go_package=github.com/username/project/fb

      - language: java
        out: "./generated/java"
        opt:
          - java_package_prefix=com.fb

      - language: kotlin
        out: "./generated/kotlin"

      - language: php
        out: "./generated/php"

      - language: swift
        out: "./generated/swift"

      - language: dart
        out: "./generated/dart"

      - language: csharp
        out: "./generated/csharp"

      - language: python
        out: "./generated/python"

      - language: rust
        out: "./generated/rust"

      - language: ts
        out: "./generated/ts"

  - name: nanobuffers
    out: "./nano"
```

- `inputs` define your schema sources.
- `plugins` define how `.proto` files are converted and which language targets to generate.
- `flatbuffers` supports multiple `languages` with optional config per target.
- `nanobuffers` is **C-only**, so it does not require a `languages` field.
- `opt` is required only for `go` (`go_package`) and `java` (`java_package_prefix`).

## Examples

The [examples/](examples/) directory contains ready-to-run `.proto` schemas and pre-generated output for every supported language.

### Minimal example

Single language (Go) — the simplest possible config:

```yaml
version: v1
inputs:
  - name: source
    path: "./proto"

plugins:
  - name: flatbuffers
    out: "./fbs"
    languages:
      - language: go
        out: "./generated/go"
        opt:
          - go_package=github.com/username/project/fb
```

### Multi-language production example

```yaml
version: v1
inputs:
  - name: source
    path: "./schemas"

plugins:
  - name: flatbuffers
    out: "./build/fbs"
    languages:
      - language: go
        out: "./services/go/generated"
        opt:
          - go_package=github.com/company/project/fb

      - language: cpp
        out: "./native/cpp/generated"

      - language: java
        out: "./services/java/generated"
        opt:
          - java_package_prefix=com.company.project.fb

      - language: ts
        out: "./web/src/generated"

      - language: python
        out: "./analytics/generated"

  - name: nanobuffers
    out: "./build/nano"
```

See [examples/configs/](examples/configs/) for all-languages and other ready-to-use configs.

## License

Copyright © 2026 The Protobuf Project

Licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for details.
