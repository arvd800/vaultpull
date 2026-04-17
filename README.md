# vaultpull

> CLI tool to sync HashiCorp Vault secrets into local `.env` files safely

---

## Installation

```bash
go install github.com/youruser/vaultpull@latest
```

Or build from source:

```bash
git clone https://github.com/youruser/vaultpull.git
cd vaultpull
go build -o vaultpull .
```

---

## Usage

Authenticate with Vault and pull secrets into a local `.env` file:

```bash
export VAULT_ADDR="https://vault.example.com"
export VAULT_TOKEN="s.yourtoken"

vaultpull --path secret/data/myapp --output .env
```

This will fetch all key-value pairs stored at the given Vault path and write them to `.env` in the format:

```
KEY=value
ANOTHER_KEY=another_value
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--path` | Vault secret path to pull from | *(required)* |
| `--output` | Output file path | `.env` |
| `--overwrite` | Overwrite existing file | `false` |

---

## Requirements

- Go 1.21+
- A running HashiCorp Vault instance
- Valid `VAULT_ADDR` and `VAULT_TOKEN` environment variables

---

## License

[MIT](LICENSE) © 2024 youruser