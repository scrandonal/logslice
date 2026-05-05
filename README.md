# logslice

Fast log file parser that extracts time-range segments from large structured logs.

## Installation

```bash
go install github.com/yourusername/logslice@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/logslice.git
cd logslice && go build ./...
```

## Usage

Extract log entries between two timestamps:

```bash
logslice --from "2024-01-15T08:00:00Z" --to "2024-01-15T09:00:00Z" --file app.log
```

Pipe from stdin:

```bash
cat app.log | logslice --from "2024-01-15T08:00:00Z" --to "2024-01-15T09:00:00Z"
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--file` | Path to log file | stdin |
| `--from` | Start timestamp (RFC3339) | required |
| `--to` | End timestamp (RFC3339) | required |
| `--format` | Timestamp layout | RFC3339 |
| `--field` | Timestamp field name | `time` |

### Example Output

```
{"time":"2024-01-15T08:12:44Z","level":"info","msg":"server started"}
{"time":"2024-01-15T08:45:01Z","level":"error","msg":"connection timeout"}
```

## Requirements

- Go 1.21+

## License

MIT © 2024 yourusername