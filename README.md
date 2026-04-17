# envlint
[![BU6IXwu.md.jpg](https://iili.io/BU6IXwu.md.jpg)](https://freeimage.host/i/BU6IXwu)

[![CI](https://img.shields.io/badge/ci-pending-lightgrey)](#)
[![Go Version](https://img.shields.io/badge/go-1.24+-blue)](#)
[![License](https://img.shields.io/badge/license-MIT-lightgrey)](#)

`envlint` is a lightweight Go CLI that audits `.env` files against your codebase. It helps you spot missing variables, unused values, `.env.example` drift, duplicate keys, and basic secret-leak risks before they turn into production surprises.

## Install

```bash
go install github.com/drawliin/envlint@latest
```

## Usage

| Command | What it does |
| --- | --- |
| `envlint` | Audit the current directory using `.env` and `.env.example`. |
| `envlint --path ./app` | Audit a specific project directory instead of the current one. |
| `envlint --env .env.local` | Use a different env file as the source of actual values. |
| `envlint --example-env example.env` | Use a different example env filename for comparison. |
| `envlint --fix` | Add keys found in the env file but missing from the example env file, writing them as empty values. |
| `envlint --json` | Output the audit result as JSON instead of terminal-friendly text. |
| `envlint --strict` | Exit with a non-zero status when issues are found. |

### Sample output

```text
❌ Missing vars
  - STRIPE_SECRET_KEY

⚠️ Unused vars
  - LEGACY_FLAG

❌ .env.example missing from .env
  - DATABASE_URL

⚠️ .env missing from .env.example
  - INTERNAL_API_TOKEN

⚠️ Secret leak detection
  - .env is not ignored by .gitignore

Summary: 5 issues (2 blocking, 3 non-blocking)
```

## What it checks

- Missing vars referenced in `.go`, `.js`, `.ts`, and `.py` source files but absent from `.env`
- Unused vars defined in `.env` but never referenced in the codebase
- Drift between `.env` and `.env.example`
- Duplicate keys inside `.env` files
- Whether `.env` is protected by `.gitignore`

## Contributing

Issues and pull requests are welcome. If you want to contribute, start by opening an issue or sharing the behavior you want to improve, then add tests and keep changes focused and idiomatic.
