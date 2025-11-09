[<-- Back to Main README](../README.md)

# Buffman Convert 🤖✨

Welcome to the `convert` command. This is your tool for transforming `.proto` schema files into intermediate buffer schemas (`.fbs` for FlatBuffers and `.nbf` for NanoBuffers). You can run conversions directly through the CLI or automate them using your `buffman.yml`.

Both `flatbuffers` and `nanobuffers` are subcommands of `buffman convert`. More formats may be supported in the future.

## 🔧 Quick Command Reference

| Command                                                   | Description                                           |
| --------------------------------------------------------- | ----------------------------------------------------- |
| `buffman convert flatbuffers`                             | Converts `.proto` files to `.fbs` using `buffman.yml` |
| `buffman convert flatbuffers -I ./my-protos -o ./my-fbs`  | Converts using direct CLI flags                       |
| `buffman convert nanobuffers`                             | Converts `.proto` files to `.nbf` using `buffman.yml` |
| `buffman convert nanobuffers -I ./my-protos -o ./my-nbf`  | Converts using direct CLI flags                       |

## 🧠 Full Command

### FlatBuffers

```bash
buffman convert flatbuffers --proto_dir ./my-protos --output_dir ./my-fbs
```

Or using short flags:

```bash
buffman convert flatbuffers -I ./my-protos -o ./my-fbs
```

### NanoBuffers

```bash
buffman convert nanobuffers --proto_dir ./my-protos --output_dir ./my-nbf
```

Or using short flags:

```bash
buffman convert nanobuffers -I ./my-protos -o ./my-nbf
```

These commands read all `.proto` files from the specified input directory and write the converted schema files to the specified output directory.  
If `--output_dir` is not provided, Buffman defaults to the current working directory.

## 🚀 Usage Modes

You can run the convert command in two ways — CLI flags or using a configuration file.

### CLI Mode

Use this for quick conversions when you do not want to set up a config file.

```bash
buffman convert flatbuffers --proto_dir ./my-protos --output_dir ./my-fbs
buffman convert nanobuffers --proto_dir ./my-protos --output_dir ./my-nbf
```

### Docker CLI Mode

If you're using Buffman via Docker:

```bash
docker run --rm \
    -v $(pwd):/buffman \
    -w /buffman \
    ghcr.io/tarran-sidhaarth/buffman convert flatbuffers \
    --proto_dir /buffman/protos \
    --output_dir /buffman/flatbuffers

docker run --rm \
    -v $(pwd):/buffman \
    -w /buffman \
    ghcr.io/tarran-sidhaarth/buffman convert nanobuffers \
    --proto_dir /buffman/protos \
    --output_dir /buffman/nanobuffers
```

> 📌 Make sure paths like `/buffman/protos` and `/buffman/flatbuffers` exist inside your mounted directory. Paths must be **relative to `/buffman`** when using Docker.


## 🚩 Flags

| Flag           | Shorthand | Description                                                                     | Required |
| -------------- | --------- | ------------------------------------------------------------------------------- | -------- |
| `--proto_dir`  | `-I`      | The directory containing `.proto` files                                         | Yes (CLI) |
| `--output_dir` | `-o`      | The directory where `.fbs` or `.nbf` files will be written. Defaults to current directory | No       |

[<-- Back to Main README](../README.md)
