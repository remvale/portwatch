# portwatch

A lightweight CLI daemon that monitors port usage and alerts on unexpected bindings or conflicts.

---

## Installation

```bash
go install github.com/youruser/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/youruser/portwatch.git && cd portwatch && go build -o portwatch .
```

---

## Usage

Start the daemon to watch all active ports:

```bash
portwatch start
```

Watch specific ports and get alerted on unexpected bindings:

```bash
portwatch watch --ports 8080,5432,6379
```

Run a one-time snapshot of current port bindings:

```bash
portwatch scan
```

Example output:

```
[INFO]  Watching ports: 8080, 5432, 6379
[ALERT] Unexpected binding detected on :8081 by process nginx (pid 3821)
[WARN]  Port conflict: :5432 claimed by unknown process (pid 4102)
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--ports` | Comma-separated list of ports to watch | all |
| `--interval` | Poll interval in seconds | `5` |
| `--log` | Log output file path | stdout |
| `--quiet` | Suppress info messages | false |

---

## License

MIT © 2024 youruser