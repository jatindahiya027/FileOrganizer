# 🗂 File Organizer — Wails + Go + React

A native cross-platform desktop app that organizes files in a directory into categorized subfolders.
Built with **Go** backend, **React + TypeScript** frontend, and **Wails v2** for native OS integration.

---

## Features

- **15 file categories** — Images, Videos, Audio, Documents, Spreadsheets, Presentations, Ebooks, 3D Files, Executables, Archives, Disk Images, Code, Fonts, Database, Others
- **200+ extensions** recognized across all categories
- **Multithreaded** — goroutine worker pool, configurable thread count via UI slider
- **Platform-aware moves** — `robocopy` on Windows, `os.Rename` + `io.Copy` fallback on macOS/Linux
- **Real-time log feed** — every move, skip, and error appears live in the UI
- **Dry Run mode** — preview what would happen without moving anything
- **Conflict resolution** — `file_(1).jpg`, `file_(2).jpg`, etc. for name clashes
- **Flat-scan only** — only organizes files directly in the selected folder, never recurses into subdirectories
- **JSON log file** — `organizer.log` written after each run for audit trail
- **Undo-safe** — log path shown so you can review before re-running

---

## Requirements

| Tool | Version |
|------|---------|
| Go   | 1.21+   |
| Node | 18+     |
| Wails CLI | v2.9+ |

### Install Wails CLI

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

Verify:
```bash
wails doctor
```

---

## Project Structure

```
file-organizer/
├── main.go                     # Wails entrypoint
├── app.go                      # All exported Go→JS methods
├── wails.json                  # Wails project config
├── go.mod
├── organizer/
│   ├── types.go                # Shared structs (LogEntry, ProgressEvent, Summary…)
│   ├── categories.go           # Extension → Category map (200+ extensions)
│   ├── organizer.go            # Worker pool, scanning, task dispatch
│   ├── platform.go             # os.Rename / io.Copy move logic (macOS/Linux)
│   └── platform_windows.go     # robocopy move logic (Windows only, build tag)
└── frontend/
    ├── index.html
    ├── package.json
    ├── vite.config.ts
    ├── tailwind.config.js
    └── src/
        ├── main.tsx
        ├── App.tsx             # Full UI — single file, all components
        ├── types.ts            # TypeScript types mirroring Go structs
        └── index.css           # Tailwind + custom utilities
```

---

## Development

```bash
# Install frontend deps
cd frontend && npm install && cd ..

# Start dev server (hot-reload both Go and React)
wails dev
```

## Build (Production)

```bash
wails build
```

Output binary: `build/bin/FileOrganizer` (or `FileOrganizer.exe` on Windows)

---

## How It Works

### Scanning
Only the **top-level** files in the selected directory are scanned. Subdirectories and hidden files (`.`-prefixed) are always skipped.

### Worker Pool
```
Files → buffered channel → N goroutines → results channel → UI events
```
Thread count is user-configurable (1–32) via the UI slider. Defaults to `CPU cores × 2`, capped at 16.

### File Moving
- **Windows**: `robocopy <srcDir> <dstDir> <filename> /MOV /R:3 /W:1`
- **macOS/Linux**: `os.Rename()` (atomic, same volume) → falls back to `io.Copy` + `os.Remove` for cross-device

### Real-time Events
Go emits two Wails events:
- `"log"` — `LogEntry` for each file action (success / warn / error / info)
- `"progress"` — `ProgressEvent` after each file is processed (percent, counts, current filename)

React subscribes via `runtime.EventsOn('log', handler)` and updates state immediately.

---

## Safety Rules

Files that are **never** moved:
- Hidden files (`.gitignore`, `.DS_Store`, etc.)
- `organizer.log`
- `desktop.ini`, `thumbs.db`
- The app binary itself
- Files inside already-categorized subfolders (flat scan only)

---

## Customization

**Add a new category**: edit `organizer/categories.go` and add entries to `ExtensionMap`.

**Change colors**: edit `CATEGORY_COLORS` in `frontend/src/App.tsx` or `organizer/categories.go`.

**Exclude extensions**: use the Exclude Extensions field in the UI (e.g. `.tmp, .log, .bak`).

---

## License

MIT
