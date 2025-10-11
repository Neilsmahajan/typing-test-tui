# typing-test-tui

[![CI](https://github.com/Neilsmahajan/typing-test-tui/actions/workflows/ci.yml/badge.svg)](https://github.com/Neilsmahajan/typing-test-tui/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/neilsmahajan/typing-test-tui.svg)](https://pkg.go.dev/github.com/neilsmahajan/typing-test-tui)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

`typing-test-tui` is a terminal-native typing trainer built in Go with [Cobra](https://github.com/spf13/cobra), [Bubble Tea](https://github.com/charmbracelet/bubbletea), and [Lip Gloss](https://github.com/charmbracelet/lipgloss). It ships with multiple practice modes, curated word lists in several natural languages, and code snippets across popular programming languages. Whether you want to benchmark your WPM or practice typing syntax-heavy snippets, this TUI keeps you focused on your keyboard.

Built and maintained by **Neil Mahajan** (<neilsmahajan@gmail.com> · [links.neilsmahajan.com](https://links.neilsmahajan.com/)).

## Table of contents

- [Features](#features)
- [Requirements](#requirements)
- [Installation](#installation)
- [Quick start](#quick-start)
- [Usage](#usage)
  - [Modes](#modes)
  - [Flags](#flags)
  - [Languages](#languages)
- [Development](#development)
  - [Project layout](#project-layout)
  - [Makefile tasks](#makefile-tasks)
  - [Testing & linting](#testing--linting)
- [Continuous integration](#continuous-integration)
- [Contributing](#contributing)
- [License](#license)
- [Contact](#contact)

## Features

- **Three focused practice modes** – quote, words, and time-based sessions tailored to different training goals.
- **Rich language dataset** – curated word banks for English, Spanish, French, and Simplified Chinese, plus code samples for C, Go, JavaScript, Rust, and more.
- **Adaptive configuration** – toggle punctuation, numbers, durations, and word counts with ergonomic CLI flags.
- **Accessible TUI design** – responsive layouts, multi-line quotes, and consistent theming courtesy of Bubble Tea and Lip Gloss.
- **Offline-friendly** – ships with all required datasets; no external network calls once installed.

## Requirements

- Go **1.25** or newer (module minimum specified in `go.mod`).
- macOS, Linux, or Windows terminal with ANSI color support.

Optional tooling used by the project:

- [`golangci-lint`](https://golangci-lint.run/) for static analysis.
- [`make`](https://www.gnu.org/software/make/) for development shortcuts (recommended).

## Installation

### Go install

```sh
go install github.com/neilsmahajan/typing-test-tui@latest
```

The binary will be placed in your `$GOBIN` (default: `$HOME/go/bin`).

### Build from source

```sh
git clone https://github.com/Neilsmahajan/typing-test-tui.git
cd typing-test-tui
make build
```

This produces `./bin/typing-test-tui` using the default `Makefile` target.

## Quick start

```sh
typing-test-tui --mode quote
```

Use <kbd>Ctrl+C</kbd> at any time to exit the session. At the end of each test you’ll see your words-per-minute, accuracy, and mistake breakdown in the terminal.

## Usage

Run `typing-test-tui --help` (or `make run ARGS="--help"`) to view all options directly in your terminal.

### Modes

| Mode    | Description                                                 | Key options                                                  |
| ------- | ----------------------------------------------------------- | ------------------------------------------------------------ |
| `quote` | Type through inspirational quotes or programming aphorisms. | Language only. Duration & word-count flags are ignored.      |
| `words` | Timed practice over a fixed set of words.                   | `--word-count`, `--include-punctuation`, `--include-numbers` |
| `time`  | Open-ended stream of words for a chosen duration.           | `--duration`, `--include-punctuation`, `--include-numbers`   |

### Flags

| Flag                          | Default   | Modes           | Description                                                           |
| ----------------------------- | --------- | --------------- | --------------------------------------------------------------------- |
| `-m`, `--mode`                | `quote`   | all             | Select the practice mode (`quote`, `words`, `time`).                  |
| `-l`, `--language`            | `english` | all             | Choose the content language or code corpus (see list below).          |
| `-d`, `--duration`            | `60`      | `time`          | Session length in seconds. Must be one of `15`, `30`, `60`, or `120`. |
| `-w`, `--word-count`          | `50`      | `words`         | Total words in the session. Pick from `10`, `25`, `50`, or `100`.     |
| `-p`, `--include-punctuation` | `false`   | `words`, `time` | Adds punctuation symbols to the text stream.                          |
| `-n`, `--include-numbers`     | `false`   | `words`, `time` | Adds numbers to the text stream.                                      |

Invalid combinations return actionable error messages before the TUI launches, preventing accidental misuse.

### Languages

Natural languages:

- `english`
- `spanish`
- `french`
- `chinese_simplified`

Programming corpora:

- `code_assembly`, `code_c`, `code_c++`, `code_csharp`, `code_css`, `code_go`, `code_java`, `code_javascript`, `code_kotlin`, `code_lua`, `code_php`, `code_python`, `code_r`, `code_ruby`, `code_rust`, `code_typescript`

Aliases such as `en`, `es`, `rust`, or `typescript` are automatically normalized; see `internal/models/config.go` for the full mapping. Each dataset lives under `internal/data/quotes` and `internal/data/words` and can be extended with your own JSON files.

## Development

### Project layout

- `main.go` – entry point that simply delegates to `cmd/`.
- `cmd/` – Cobra commands and CLI flag wiring.
- `internal/app/` – orchestrates session state and transitions.
- `internal/modes/` – mode-specific services for quotes, timed tests, and word lists.
- `internal/ui/` – Bubble Tea models, views, and input components.
- `internal/data/` – JSON corpora for quotes and word lists across languages and code stacks.

### Makefile tasks

Common development shortcuts:

| Target                        | Description                                      |
| ----------------------------- | ------------------------------------------------ |
| `make help`                   | Print all available tasks and default variables. |
| `make build`                  | Compile the binary into `./bin/typing-test-tui`. |
| `make run ARGS="--mode time"` | Launch the TUI with additional CLI flags.        |
| `make test`                   | Execute `go test ./...`.                         |
| `make lint`                   | Run `go vet` and `golangci-lint` (if installed). |
| `make fmt`                    | Apply `gofmt` to all Go files.                   |
| `make tidy`                   | Sync module dependencies.                        |
| `make clean`                  | Remove the `bin/` output directory.              |

### Testing & linting

Unit tests live alongside their packages inside `internal/`. Run them with:

```sh
make test
# or
go test ./...
```

Static analysis combines `go vet` and `golangci-lint`:

```sh
make lint
# ensures gofmt, vet, and golangci-lint all succeed
```

## Continuous integration

The repository ships with a GitHub Actions workflow (`.github/workflows/ci.yml`) that runs on push and pull requests targeting `main`. The pipeline performs:

1. Go toolchain setup (Go 1.25 from `go.mod`).
2. Formatting verification via `gofmt -l`.
3. `go vet` static analysis.
4. `golangci-lint` with a 5-minute timeout, built using Go 1.25.
5. `go test ./...` across all packages.

Badges for the CI status appear at the top of this README.

## Contributing

Contributions are welcome! If you’d like to propose improvements:

1. Fork the repository and create a topic branch.
2. Run `make fmt && make lint && make test` before submitting your pull request.
3. Open a PR with context describing what problem you’re solving or which feature you’re adding.

Bug reports and feature requests are also appreciated—feel free to open an issue.

## License

Distributed under the [MIT License](LICENSE). See the license file for full text.

## Contact

- Author: **Neil Mahajan**
- Email: <neilsmahajan@gmail.com>
- Links: [links.neilsmahajan.com](https://links.neilsmahajan.com/)

If you build something cool with `typing-test-tui`, let me know—I’d love to hear about it!
