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

The project follows **hexagonal architecture** (ports & adapters) with **vertical slicing**. Each functional group is a self-contained slice with its own domain, application, and infrastructure layers.

### Pipeline

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

### Hexagonal Layers

Each vertical slice contains three layers:

- **`domain/`** — Entities, value objects, and port interfaces. Zero external dependencies.
- **`application/`** — Use cases and services that orchestrate domain logic.
- **`infrastructure/`** — Adapters: concrete implementations of domain ports (JSON parser, filesystem loader, gopdf renderer).

### Project Structure

```
pen2pdf/
├── cmd/                                        # Cobra CLI commands (driving adapters)
│   ├── root.go                                 # Root command, version flag
│   ├── render.go                               # pen2pdf render
│   ├── validate.go                             # pen2pdf validate
│   └── info.go                                 # pen2pdf info
├── internal/
│   ├── shared/                                 # Cross-cutting shared types
│   │   ├── domain/
│   │   │   ├── node.go                         # Node interface, Frame, Text structs
│   │   │   ├── document.go                     # Document top-level container
│   │   │   ├── fill.go                         # Fill types (solid color, image)
│   │   │   └── variable.go                     # Variable type system
│   │   ├── application/
│   │   └── infrastructure/
│   ├── parser/                                 # Parse vertical slice
│   │   ├── domain/
│   │   │   └── port.go                         # Parser port interface
│   │   ├── application/
│   │   │   └── service.go                      # Parse orchestration
│   │   └── infrastructure/
│   │       └── json/
│   │           └── parser.go                   # JSON .pen file adapter
│   ├── resolver/                               # Resolve vertical slice
│   │   ├── domain/
│   │   │   └── port.go                         # Resolver port interface
│   │   ├── application/
│   │   │   └── service.go                      # Resolution orchestration
│   │   └── infrastructure/
│   │       └── engine/
│   │           └── resolver.go                 # $variable substitution engine
│   ├── asset/                                  # Asset loading vertical slice
│   │   ├── domain/
│   │   │   └── port.go                         # AssetLoader port, ImageData, FontData
│   │   ├── application/
│   │   │   └── service.go                      # Asset orchestration
│   │   └── infrastructure/
│   │       └── fs/
│   │           └── loader.go                   # Filesystem image/font loader
│   ├── layout/                                 # Layout vertical slice
│   │   ├── domain/
│   │   │   ├── port.go                         # LayoutEngine port, TextMeasurer
│   │   │   └── types.go                        # LayoutBox, Page
│   │   ├── application/
│   │   │   └── service.go                      # Layout orchestration
│   │   └── infrastructure/
│   │       └── flexbox/
│   │           ├── engine.go                   # Flexbox layout algorithm
│   │           └── measure.go                  # Text intrinsic measurement
│   └── renderer/                               # Render vertical slice
│       ├── domain/
│       │   └── port.go                         # Renderer port interface
│       ├── application/
│       │   └── service.go                      # Render orchestration
│       └── infrastructure/
│           └── pdf/
│               ├── renderer.go                 # gopdf document/page setup
│               └── draw.go                     # Drawing primitives
├── main.go
├── .air.toml
├── Makefile
├── go.mod
└── go.sum
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
