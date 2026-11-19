# Agent Monitor: Codex Session Reader

A small Go service that repeatedly reads a Codex session file.

It can watch either:

- one specific session file; or
- a directory, in which case it reads the most recently modified regular file.

The service reads the whole selected file on each poll and logs the path, size, and modification time. It does not parse or store the session contents yet.

## Run it

From the repository root:

```bash
make run
```

That is all you need for local development. If there is no `.env` configuration, the service uses the included sample:

```text
test/data/sample-session.jsonl
```

Stop it with `Ctrl+C`.

## Point it at your Codex sessions

Create your local configuration:

```bash
cp .env.example .env
```

Edit `.env`:

```dotenv
CODEX_SESSION_PATH=/Users/you/.codex/sessions
POLL_INTERVAL=5s
```

`CODEX_SESSION_PATH` may be a single file or a directory. When it is a directory, the newest file is selected on every poll. Process environment variables override values in `.env`.

For a one-time override, without editing `.env`:

```bash
CODEX_SESSION_PATH="$HOME/.codex/sessions" make run
```

## Useful commands

```bash
make help    # List commands
make run     # Run the service
make test    # Run tests
make build   # Build bin/session-reader
```

You can also run the Go command directly:

```bash
go run ./cmd/session-reader
```

Command-line flags override `.env`:

```bash
go run ./cmd/session-reader \
  -session-path "$HOME/.codex/sessions" \
  -poll-interval 10s
```

## Configuration order

Settings are applied in this order, from highest to lowest priority:

1. command-line flags;
2. process environment variables;
3. values from `.env`;
4. built-in defaults.

| Setting | Default | Description |
|---|---|---|
| `CODEX_SESSION_PATH` | `test/data/sample-session.jsonl` | Session file or directory |
| `POLL_INTERVAL` | `5s` | Time between reads |

`CODEX_SESSION_DIR` is also accepted as a backwards-compatible alias when `CODEX_SESSION_PATH` is unset.

## Requirements

- Go 1.23 or newer
- `make` for the short commands above (optional)
