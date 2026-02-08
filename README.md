# pen2pdf

[![Build](https://github.com/vpedrosa/pen2pdf/actions/workflows/build.yml/badge.svg)](https://github.com/vpedrosa/pen2pdf/actions/workflows/build.yml)
[![Test](https://github.com/vpedrosa/pen2pdf/actions/workflows/test.yml/badge.svg)](https://github.com/vpedrosa/pen2pdf/actions/workflows/test.yml)
[![Lint](https://github.com/vpedrosa/pen2pdf/actions/workflows/lint.yml/badge.svg)](https://github.com/vpedrosa/pen2pdf/actions/workflows/lint.yml)
[![Release](https://img.shields.io/github/v/release/vpedrosa/pen2pdf)](https://github.com/vpedrosa/pen2pdf/releases/latest)
[![Go](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![Cobra](https://img.shields.io/badge/Cobra-v1.10-blue)](https://cobra.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

CLI tool to generate PDF documents from `.pen` files — a JSON-based design format.

## Features

- **Full layout engine** — Flexbox-like layout with `vertical`/`horizontal` stacking, `gap`, `padding`, `justifyContent`, `alignItems`, and `fill_container` responsive sizing
- **Typography** — Font embedding with `fontFamily`, `fontSize`, `fontWeight`, `fontStyle`, `letterSpacing`, `lineHeight`, and `textAlign`
- **Auto-sizing** — Frames without explicit dimensions automatically size to fit their content
- **Design variables** — Reusable `$variable` tokens for colors, fonts, spacing, and sizes
- **Image fills** — Background images with cover mode, clipping, and configurable opacity
- **Rounded corners** — Frames with `cornerRadius` and solid or image backgrounds
- **Multi-page** — Each top-level frame becomes a separate PDF page
- **Auto font download** — Missing fonts are detected and downloaded from Google Fonts with a single prompt
- **Fallback fonts** — Embedded Go fonts as fallback when fonts are unavailable

## Installation

### From GitHub Releases (recommended)

Download the latest pre-compiled binary for your platform from the [Releases](https://github.com/vpedrosa/pen2pdf/releases) page.

```bash
# Linux (amd64)
curl -Lo pen2pdf https://github.com/vpedrosa/pen2pdf/releases/latest/download/pen2pdf-linux-amd64
chmod +x pen2pdf
sudo mv pen2pdf /usr/local/bin/

# macOS (Apple Silicon)
curl -Lo pen2pdf https://github.com/vpedrosa/pen2pdf/releases/latest/download/pen2pdf-darwin-arm64
chmod +x pen2pdf
sudo mv pen2pdf /usr/local/bin/

# macOS (Intel)
curl -Lo pen2pdf https://github.com/vpedrosa/pen2pdf/releases/latest/download/pen2pdf-darwin-amd64
chmod +x pen2pdf
sudo mv pen2pdf /usr/local/bin/
```

### From source

Requires Go 1.21+.

```bash
go install github.com/vpedrosa/pen2pdf@latest
```

### Build from source

```bash
git clone https://github.com/vpedrosa/pen2pdf.git
cd pen2pdf
go build -o pen2pdf .
```

## Usage

### Render a `.pen` file to PDF

```bash
pen2pdf render input.pen -o output.pdf
```

On first run, `pen2pdf` detects missing fonts and offers to download them from Google Fonts:

```
Missing 15 font(s):
  - Inter 700
  - Montserrat 300
  ...

Download from Google Fonts to ./fonts? [Y/n] y
  downloaded: Inter-Bold.ttf
  downloaded: Montserrat-Light.ttf
  ...
15 font(s) downloaded to ./fonts

PDF written to output.pdf (2 pages)
```

Downloaded fonts are saved next to the `.pen` file in a `fonts/` directory and reused on subsequent runs.

### Render specific pages

```bash
pen2pdf render input.pen --pages "Travel Flyer"
```

### Non-interactive mode

```bash
pen2pdf render input.pen --no-prompt
```

Skips the font download prompt and uses fallback fonts. Useful for CI/CD pipelines.

### Validate a `.pen` file

```bash
pen2pdf validate input.pen
```

Checks that the file parses and variables resolve correctly, without rendering.

### Show document info

```bash
pen2pdf info input.pen
```

Displays pages, variables, and fonts used in the document.

## The `.pen` Format

A `.pen` file is a JSON document describing a tree of visual nodes:

- **`frame`** — Container with optional fill (solid color or image), corner radius, clipping, and layout properties
- **`text`** — Text node with full typography control

```json
{
  "version": "2.7",
  "children": [
    {
      "type": "frame",
      "id": "page1",
      "name": "My Page",
      "width": 800,
      "height": 1000,
      "fill": "#FFFFFF",
      "layout": "vertical",
      "gap": 20,
      "padding": 40,
      "justifyContent": "center",
      "alignItems": "center",
      "children": [
        {
          "type": "text",
          "id": "title",
          "name": "title",
          "content": "Hello World",
          "fill": "$primary-color",
          "fontFamily": "Inter",
          "fontSize": 48,
          "fontWeight": "700"
        }
      ]
    }
  ],
  "variables": {
    "primary-color": { "type": "color", "value": "#FF6B35" }
  }
}
```

## Development

```bash
task          # lint + test + build
task test     # go test with race detector and coverage
task build    # build to ./bin/pen2pdf
task lint     # golangci-lint
task fmt      # format code
task dev      # start Air hot-reload
task clean    # remove build artifacts
```

## License

MIT License. See [LICENSE](LICENSE) for details.
