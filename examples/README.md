# Examples

This directory contains example `.proto` schemas and the generated output for all supported languages.

## Structure

```text
examples/
├── configs/           # Ready-to-use buffman.yml configs
│   ├── minimal.yml        — single language (Go) quickstart
│   ├── production.yml     — common multi-language setup with nanobuffers
│   └── all-languages.yml  — every supported language + nanobuffers
├── protos/            # Source .proto schemas used by all examples
│   ├── common/            — shared types (address, status, timestamp)
│   └── services/          — service definitions (user, notification, analytics)
├── flatbuffers/       — converted .fbs schemas
├── go/                — FlatBuffers Go bindings
├── cpp/               — FlatBuffers C++ headers
├── java/              — FlatBuffers Java code
├── kotlin/            — FlatBuffers Kotlin code
├── php/               — FlatBuffers PHP code
├── swift/             — FlatBuffers Swift code
├── dart/              — FlatBuffers Dart code
├── csharp/            — FlatBuffers C# code
├── python/            — FlatBuffers Python code
├── rust/              — FlatBuffers Rust code
├── ts/                — FlatBuffers TypeScript code
├── nanobuffers/       — NanoBuffers C output
└── justfile           — build automation (requires just)
```

## Running the examples

Build the binary and regenerate all output from the `examples/` directory:

```bash
just run
```

Or run individual steps:

```bash
just build                 # compile buffman
just convert               # proto → fbs + nanobuf
just generate              # fbs → all language bindings
```

## Using a config

From the repo root, run any of the example configs with:

```bash
buffman generate -f examples/configs/minimal.yml
buffman generate -f examples/configs/production.yml
buffman generate -f examples/configs/all-languages.yml
```
