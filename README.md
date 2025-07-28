<div align="center">
  <img src="assets/logo.png" />
  <div>
    <img alt="Version" src="https://img.shields.io/github/package-json/v/diogopereiradev/httpzen/latest?style=for-the-badge">
    <img alt="Top language" src="https://img.shields.io/github/languages/top/diogopereiradev/httpzen?style=for-the-badge">
    <img alt="Closed pull requests" src="https://img.shields.io/github/issues-pr-closed-raw/diogopereiradev/httpzen?style=for-the-badge">
    <img alt="GitHub" src="https://img.shields.io/github/license/diogopereiradev/httpzen?style=for-the-badge&cacheSeconds=60">
  </div>
</div>

<br />
<div align="center">
  <p>HTTPZen is a modern, terminal-based HTTP client designed for developers who want a fast, beautiful, and scriptable way to interact with APIs and web services. Built with Go, HTTPZen offers a rich TUI (Text User Interface) experience, advanced request/response inspection, and flexible configuration.</p>

  [Downloads](#downloads) •
  [Building from source](#building-from-source) •
  [Contribute](#contribute)
  
</div>

## Features

- **Interactive TUI**: Navigate requests, responses, and configuration with a clean, keyboard-driven interface.
- **Request Builder**: Easily craft HTTP requests with support for all methods, custom headers, and body types (JSON, form, file, etc).
- **Response Viewer**: Pretty-print JSON, HTML, and plain text responses. Inspect headers, status, and timings.
- **Cross-platform**: Runs on Linux, Windows, and macOS.
- **Benchmarking**: Stress test your API routes.
- **Scriptable**: Integrate with shell scripts and automate API testing.

<br />

## Downloads [](#downloads)
You can downloads the project release on [Downloads](https://github.com/diogopereiradev/httpzen/releases/latest)

<br />

## Prerequisites
- Go 1.20 or newer
- (Optional) Flatpak, dpkg, or rpm for package installation

<br />

## Installation

#### From Source
```sh
git clone https://github.com/diogopereiradev/httpzen.git
cd httpzen
go mod download
make build
./build/httpzen
```

#### Debian/Ubuntu
```sh
sudo dpkg -i build/debian/httpzen.deb
```

#### RPM (Fedora, CentOS, etc)
```sh
sudo rpm -i build/httpzen.rpm
```

#### Flatpak
```sh
flatpak install --user ./build/flatpak/httpzen.flatpak
```

<br />

## Usage

To start HTTPZen:
```sh
go run main.go
# or if installed
httpzen
```

### Example: Sending a GET Request
1. Launch HTTPZen `httpzen [METHOD] [URL] [FLAGS...]`

### Command Line Options
- `httpzen version` — Show version and build info
- `httpzen help` — Show help and usage

<br />

## Development

1. Fork and clone the repository.
2. Create a new branch for your feature or bugfix:
```sh
git checkout -b feat/my-feature
```

3. Run tests:
```sh
make test
```

4. Lint the code:
```sh
make lint
```

5. Build the project:
```sh
make build
```

See [CONTRIBUTING.md](CONTRIBUTING.md) for more details.

<br />

## Contributing
Contributions are welcome! Please read the [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

<br />

## Acknowledgements
- [Charmbracelet lipgloss](https://github.com/charmbracelet/lipgloss) for estilization
- [Charmbracelet bubbletea](https://github.com/charmbracelet/bubbletea) for interactive menus
- [Testify](https://github.com/stretchr/testify) for unit testing
- [Cobra](https://github.com/spf13/cobra) for CLI framework

<br />

## License
HTTPZen is licensed under the MIT License. See [LICENSE](LICENSE) for details.
