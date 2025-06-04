# JSON Edit

A web-based JSON editor application.

## Features

- Edit JSON files via a web interface
- Compare JSON files
- Read-only mode option
- Configurable indentation

## Usage

```bash
# Run with default settings
jsonedit

# Run with custom settings
jsonedit --port 3000 --host 0.0.0.0 --indent "    " --read-only
```

## Configuration

The application can be configured using command-line flags or environment variables:

| Flag | Environment Variable | Default | Description |
|------|---------------------|---------|-------------|
| `--port` | `JSON_EDIT_PORT` | 8080 | Port to listen on |
| `--host` | `JSON_EDIT_HOST` | localhost | Host to listen on |
| `--indent` | `JSON_EDIT_INDENT` | "  " | Indentation level |
| `--read-only` | `JSON_EDIT_READ_ONLY` | false | Read-only mode |
| `--log-level` | `JSON_EDIT_LOG_LEVEL` | info | Log level (debug, info, warn, error) |

## Building and Releasing

This project uses [GoReleaser](https://goreleaser.com/) to build and release binaries for multiple platforms:

- macOS (arm64, amd64)
- Linux (arm64, amd64)
- Windows (amd64)

### Local Testing

To test the build process locally:

```bash
# Install GoReleaser if you haven't already
go install github.com/goreleaser/goreleaser@latest

# Test the build without releasing
goreleaser build --snapshot --clean
```

The built binaries will be available in the `dist/` directory.

### Creating a Release

To create a new release:

1. Tag the commit you want to release:
   ```bash
   git tag -a v0.1.0 -m "First release"
   git push origin v0.1.0
   ```

2. The GitHub Actions workflow will automatically build and release the binaries.

3. The release will be available on the GitHub Releases page.

## License

[Add license information here]