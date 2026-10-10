# envlint

![envlint CLI](https://iili.io/BU6IXwu.jpg)

[![CI](https://github.com/drawliin/envlint/actions/workflows/ci.yml/badge.svg)](https://github.com/drawliin/envlint/actions/workflows/ci.yml)
[![Go version](https://img.shields.io/badge/Go-1.23.2%2B-00ADD8?logo=go)](go.mod)
[![MIT license](https://img.shields.io/badge/license-MIT-blue)](LICENSE)

`envlint` is a Go command-line tool for auditing environment-variable **references** and `.env`/`.env.example` drift. It reports missing or unused keys, duplicate definitions, and basic `.gitignore` hygiene warnings.

It uses simple source-pattern matching. **It is not a general-purpose secret scanner or a guarantee that a secret has never been committed.**

## Install

Requires **Go 1.23.2 or newer** (see [`go.mod`](go.mod)).

```sh
go install github.com/drawliin/envlint@latest
```

Or build from source:

```sh
git clone https://github.com/drawliin/envlint.git
cd envlint
go build -o envlint .
```

## Quick start

Run the command from the root of a Go, JavaScript, TypeScript, or Python project:

```sh
envlint
envlint --path ./my-app
envlint --strict --json
```

By default, it reads `.env` and `.env.example` in the selected directory and scans supported source files recursively. It skips common dependency, cache, and build directories, including `node_modules`, `vendor`, `.venv`, and `dist`.

| Flag | Behavior |
| --- | --- |
| `--path ./app` | Scan a different project directory (default: working directory). |
| `--env .env.local` | Select the actual environment file. |
| `--example-env example.env` | Select the example/documentation file. |
| `--fix` | Append **empty** entries to the example file for keys found in the actual env file but absent from the example. |
| `--json` | Emit the audit result as JSON. |
| `--strict` | Exit with status `1` when the audit finds any issue; otherwise exit `0`. |

`--fix` does not remove unused variables, fill in missing secrets, or modify source files. Review the generated example file before committing it.

## Checks

- Variables referenced in supported `.go`, `.js`, `.ts`, and `.py` files but missing from the actual env file.
- Variables defined in the actual env file but not matched by the scanner.
- Keys present in one env file but absent from the other.
- Duplicate key definitions in either env file.
- A **basic** warning when the actual env file exists but the root `.gitignore` lacks a matching simple ignore rule.

The scanner recognizes common literal access patterns such as `os.Getenv("KEY")`, `process.env.KEY`, `process.env["KEY"]`, `os.getenv("KEY")`, `os.environ.get("KEY")`, and `os.environ["KEY"]`. Dynamically computed variable names and arbitrary syntax are not fully analyzed. An "unused" result can therefore be a false positive.

## Example output

```text
❌  Missing vars
  - STRIPE_SECRET_KEY

⚠️  Unused vars
  - LEGACY_FLAG

❌  .env.example missing from .env
  - DATABASE_URL

⚠️  .env missing from .env.example
  - INTERNAL_API_TOKEN

⚠️  Secret leak detection
  - .env is not ignored by .gitignore

Summary: 5 issues (2 blocking, 3 non-blocking)
```

The `.gitignore` check is heuristic; confirm secret handling and repository history independently.

## Development

```sh
go test ./...
go vet ./...
go build ./...
```

The [CI workflow](.github/workflows/ci.yml) runs these checks (including race-enabled tests) on pushes and pull requests.

## Contributing

Focused issues and pull requests are welcome. Include a regression test for behavior changes, and describe any syntax the scanner deliberately does not support.

## License

[MIT](LICENSE).
