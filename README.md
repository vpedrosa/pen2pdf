# pen2pdf

CLI tool to generate PDF documents from `.pen` files (Pencil design format).

## Overview

`pen2pdf` parses `.pen` files — a JSON-based design document format — and renders them into high-fidelity PDF output. Each top-level frame in the `.pen` file becomes a separate page in the resulting PDF.

### The `.pen` Format

A `.pen` file describes a tree of visual nodes with two primitive types:

- **`frame`** — Container with optional background (solid color or image), border radius, clipping, and Flexbox-like layout (`vertical`/`horizontal`, `gap`, `padding`, `justifyContent`, `alignItems`).
- **`text`** — Text node with full typography control (`fontFamily`, `fontSize`, `fontWeight`, `fontStyle`, `letterSpacing`, `lineHeight`, `textAlign`).

It also supports:

- **Design variables** — Reusable tokens (`$variable-name`) for colors, fonts, spacing, and sizes.
- **Responsive sizing** — Fixed pixel dimensions or `fill_container` to expand to the parent's available space.
- **Image fills** — Background images with `fill` mode and configurable opacity.

## Architecture

```
.pen file
   │
   ▼
┌──────────────────┐
│   1. Parser      │  JSON → typed Go AST (Node tree)
└────────┬─────────┘
         ▼
┌──────────────────┐
│   2. Resolver    │  Replace $variable references with concrete values
└────────┬─────────┘
         ▼
┌──────────────────┐
│   3. Assets      │  Load images, resolve relative paths, register fonts
└────────┬─────────┘
         ▼
┌──────────────────┐
│   4. Layout      │  Calculate absolute (x, y, w, h) for every node
│                  │  Implements a Flexbox subset:
│                  │    - fill_container resolution
│                  │    - padding, gap distribution
│                  │    - justifyContent / alignItems
│                  │    - intrinsic text measurement
└────────┬─────────┘
         ▼
┌──────────────────┐
│   5. Renderer    │  Walk the positioned tree and draw to PDF:
│                  │    - Filled rectangles with corner radius
│                  │    - Images (clipped, scaled to fill)
│                  │    - Text with font embedding
│                  │  Each root-level frame → 1 PDF page
└──────────────────┘
```

### Project Structure

```
pen2pdf/
├── cmd/
│   └── root.go              # Cobra root command
│   └── render.go            # `pen2pdf render` subcommand
│   └── info.go              # `pen2pdf info` subcommand
│   └── validate.go          # `pen2pdf validate` subcommand
├── internal/
│   ├── parser/
│   │   ├── parser.go        # JSON deserialization
│   │   └── types.go         # Node, Frame, Text, Fill, Variable structs
│   ├── resolver/
│   │   └── resolver.go      # $variable substitution
│   ├── assets/
│   │   └── loader.go        # Image and font loading
│   ├── layout/
│   │   ├── engine.go        # Layout algorithm (measure → resolve → position)
│   │   ├── measure.go       # Text intrinsic sizing
│   │   └── types.go         # LayoutBox with computed geometry
│   └── render/
│       ├── pdf.go           # PDF document assembly
│       └── draw.go          # Drawing primitives (rect, text, image)
├── example/
│   ├── test.pen             # Sample .pen file
│   ├── images/              # Referenced images
│   └── export/              # Expected output (PNG reference)
├── main.go                  # Entry point
├── go.mod
├── go.sum
├── LICENSE
├── README.md
└── CONTRIBUTING.md
```

## Usage

```bash
# Render a .pen file to PDF
pen2pdf render input.pen -o output.pdf

# Render specific pages only
pen2pdf render input.pen -o output.pdf --pages "Travel Flyer"

# Show document info (pages, variables, fonts used)
pen2pdf info input.pen

# Validate a .pen file without rendering
pen2pdf validate input.pen
```

## Building

```bash
go build -o pen2pdf .
```

## Dependencies

Key Go libraries:

| Library | Purpose |
|---------|---------|
| [cobra](https://github.com/spf13/cobra) | CLI framework |
| [gopdf](https://github.com/signintech/gopdf) | PDF generation with TTF/OTF support |

## License

MIT License. See [LICENSE](LICENSE) for details.
