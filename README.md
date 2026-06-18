# new

A fast, single-binary CLI for scaffolding files and projects.

```sh
new index.html
new my-app --template react
```

---

## Installation

```sh
go install github.com/DevShedLabs/New@latest
```

Or build from source:

```sh
git clone https://github.com/DevShedLabs/New
cd New
go build -o new .
```

---

## Usage

### Create a file

Detects the file type by extension and generates smart boilerplate automatically.

```sh
new index.html       # HTML5 skeleton
new main.go          # Go package main
new styles.css       # CSS reset + box-sizing
new app.ts           # TypeScript module stub
new script.sh        # Bash shebang + set -euo pipefail
new data.json        # Empty JSON object
new notes.md         # Markdown heading
```

Any extension without a built-in template creates an empty file.

### Scaffold a project

Use `--template` (or `-t`) to scaffold a full project from a built-in template or user blueprint.

```sh
new my-site --template html
new my-app --template react
new my-app --template my-blueprint   # user blueprint from ~/.new/blueprints/
```

The project is created as a subdirectory of the current working directory.

### Flags

| Flag | Short | Description |
|---|---|---|
| `--template` | `-t` | Template or blueprint name |
| `--output` | `-o` | Output directory (default: current directory) |
| `--var key=value` | `-v` | Pass variables to the template (repeatable) |

**Example with variables:**

```sh
new my-app --template react --var Author="Jeffrey"
```

---

## Built-in Templates

| Template | Description |
|---|---|
| `html` | HTML5 project with `index.html`, `css/main.css`, `js/main.js` |
| `react` | Vite + React + TypeScript starter |

---

## Blueprints

Blueprints are user-defined templates that live in `~/.new/blueprints/`. Any blueprint directory placed there is immediately available as a `--template` value.

**Blueprint resolution order:**
1. `~/.new/blueprints/<name>/` — user blueprints take precedence
2. Built-in embedded templates

This means you can override any built-in template by creating a blueprint with the same name.

### Creating a blueprint

A blueprint is just a folder of files. Template variables are rendered using Go's `text/template` syntax — `{{.Name}}`, `{{.Author}}`, etc.

```
~/.new/blueprints/
└── my-blueprint/
    ├── blueprint.yaml       # optional manifest
    ├── index.html
    └── src/
        └── main.js
```

Variable names in file paths are also rendered:

```
src/{{.Name}}.go   →   src/my-app.go
```

### blueprint.yaml

The manifest is optional but recommended. It documents the blueprint and declares expected variables.

```yaml
name: my-blueprint
description: My custom project starter.
vars:
  - Name
  - Author
defaults:
  Author: "Your Name"
```

### Available variables

| Variable | Value |
|---|---|
| `{{.Name}}` | The project name passed as the first argument |
| `{{.Template}}` | The template/blueprint name |
| Any `--var` flag | `--var Foo=bar` → `{{.Foo}}` |
| Blueprint defaults | Defined in `blueprint.yaml` under `defaults` |

---

## Roadmap

- `new list` — show available built-in templates and user blueprints
- `new blueprint capture` — snapshot an existing project into a blueprint
- Icon and favicon generation
- Variable prompting for blueprints that declare `vars`
- Remote blueprint fetching

---

## License

MIT
