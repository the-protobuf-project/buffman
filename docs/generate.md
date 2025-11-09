[<-- Back to Main README](../README.md)

# Buffman Generate 🚀🔥

The `generate` command turns `.fbs` schema files into language-specific source code using the `flatc` compiler. Each backend format is handled through a subcommand.

Currently supported:

- `flatbuffers` — Generates code from `.fbs` using FlatBuffers

This document covers the `flatbuffers` subcommand.

## 📚 Table of Contents

- [What is `flatbuffers`](#what-is-flatbuffers)
- [Quick Command Reference](#quick-command-reference)
- [Full Command](#full-command)
- [Usage Modes](#usage-modes)
  - [CLI Mode](#cli-mode)
  - [Config Mode (Recommended)](#config-mode-recommended)
- [Flags](#flags)
- [Note on `-f` for configuration](#note-on--f-for-configuration)

## 🧾 What is `flatbuffers`

`buffman generate flatbuffers` uses FlatBuffers to generate code in multiple languages from `.fbs` files.

It is a subcommand under `generate`, allowing Buffman to remain modular. Additional subcommands like `nanobuffers` will be supported in the future.

To use a configuration-based setup, use the root `generate` command with the `-f` flag.

## 🔧 Quick Command Reference

| Command                                                                                                                      | Description                                     |
|------------------------------------------------------------------------------------------------------------------------------|-------------------------------------------------|
| `buffman generate -f buffman.yml`                                                                                            | Generates all code defined in the config file   |
| `buffman generate flatbuffers -I ./my-fbs -l go -o ./gen/go -m "go_package=github.com/me/project/fb"`                       | Generates Go code via CLI flags                 |

## 🧠 Full Command

```bash
buffman generate flatbuffers --flatbuffers_dir ./my-fbs --language go --target_dir ./gen/go --module_options "go_package=github.com/me/project/fb"
```

Or with shorthand:

```bash
buffman generate flatbuffers -I ./my-fbs -l go -o ./gen/go -m "go_package=github.com/me/project/fb"
```

If `--target_dir` is omitted, Buffman will write the generated code to the current working directory.

## 🚀 Usage Modes

### CLI Mode

Use this mode for quick generation when targeting a single language:

```bash
buffman generate flatbuffers -I ./my-fbs -l cpp -o ./gen/cpp
```

**Run With Docker**

```bash
docker run --rm \
    -v $(pwd):/buffman \
    -w /buffman \
    ghcr.io/tarran-sidhaarth/buffman generate flatbuffers \
    -I /buffman/fbs \
    -l cpp \
    -o /buffman/gen/cpp
```

> 📌 When using Docker, paths passed to flags must be relative to `/buffman`, as that’s where the host directory is mounted.

### Config Mode (Recommended)

Use `buffman.yml` to define your schemas and language targets, then run:

```bash
buffman generate -f ./buffman.yml
```

**Run With Docker**

```bash
docker run --rm \
    -v $(pwd):/buffman \
    -w /buffman \
    ghcr.io/tarran-sidhaarth/buffman generate \
    -f /buffman/buffman.yml
```

> 📌 When using Docker, all paths inside `buffman.yml` must be relative to `/buffman`.

This approach is ideal for multi-language projects and CI automation.

## 🚩 Flags (for `flatbuffers` subcommand)

| Flag                | Shorthand | Description                                                             | Required |
|---------------------|-----------|-------------------------------------------------------------------------|----------|
| `--flatbuffers_dir` | `-I`      | Directory containing source `.fbs` files                                | Yes      |
| `--language`        | `-l`      | Target language (`go`, `cpp`, `java`, `ts`, etc.)                       | Yes      |
| `--target_dir`      | `-o`      | Directory to write generated code. Defaults to current directory        | No       |
| `--module_options`  | `-m`      | Language-specific options (e.g., `go_package`, `java_package_prefix`)   | No       |

## 🧾 Note on `-f` for configuration

Only the **root `generate` command** supports the `-f` flag:

```bash
buffman generate -f ./buffman.yml
```

This lets you run all conversions and code generation in one step based on your config file.  
Subcommands like `flatbuffers` do **not** support `-f` and rely entirely on CLI flags.

[<-- Back to Main README](../README.md)
