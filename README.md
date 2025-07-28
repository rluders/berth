# 🚢 Berth – Terminal UI for Containers

[![Go Version](https://img.shields.io/badge/go-1.24-blue?logo=go)](https://golang.org)
[![License](https://img.shields.io/github/license/rluders/berth)](LICENSE)
[![Build](https://img.shields.io/badge/build-passing-brightgreen)]()
[![Status](https://img.shields.io/badge/status-alpha-orange)]()
[![Made with Bubbletea](https://img.shields.io/badge/made%20with-bubbletea-ff69b4?logo=github)](https://github.com/charmbracelet/bubbletea)

> **Berth** is a terminal-based UI to manage your containers, images, volumes, networks, and system usage — with support for **Docker** and **Podman**.
> 🧠 Name origin: In maritime terms, a **berth** is a designated place where a ship is docked — just like containers in your stack. Clean, organized, and under control.
 
## 🚀 Installation

```bash
# 1. Clone the repository
git clone https://github.com/rluders/berth.git

# 2. Enter the project directory
cd berth

# 3. Build the binary
go build -o berth ./cmd/berth/main.go

# 4. Run it!
./berth
````

## 🧭 Usage

### 🎹 Navigation

* `1` — Containers View
* `2` — Images View
* `3` — Volumes View
* `4` — Networks View
* `5` — System View

### 🛠️ Container Actions

| Key | Action                   |
| --- | ------------------------ |
| `s` | Start selected container |
| `x` | Stop selected container  |
| `d` | Remove container         |
| `l` | View logs                |
| `i` | Inspect container        |

### 📦 Image & Volume Actions

* `d` — Remove selected image or volume

### 🧼 System Cleanup

| Key | Action           |
| --- | ---------------- |
| `b` | Basic Cleanup    |
| `a` | Advanced Cleanup |
| `t` | Total Cleanup    |

### 🔙 Back / Exit

* `q` or `esc` — Return or quit

## 📁 Project Structure

.
├── cmd/               # Entry point
├── internal/
│   ├── tui/           # Bubbletea models/views
│   ├── engine/        # Docker/Podman abstraction
│   ├── controller/    # Business logic
│   ├── state/         # Shared UI/application state
│   └── utils/         # Helpers and runners
└── assets/            # Logos, themes, future plugins
```

## 📜 License

Licensed under the [MIT License](LICENSE).
