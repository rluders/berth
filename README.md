<div align="center">

# 🚢 Berth – Terminal UI for Containers

[![Go Version](https://img.shields.io/badge/go-1.24-blue?logo=go)](https://golang.org)
[![License](https://img.shields.io/github/license/rluders/berth)](LICENSE)
[![Build](https://img.shields.io/badge/build-passing-brightgreen)]()
[![Status](https://img.shields.io/badge/status-alpha-orange)]()
[![Made with Bubbletea](https://img.shields.io/badge/made%20with-bubbletea-ff69b4?logo=github)](https://github.com/charmbracelet/bubbletea)

> **Berth** is a terminal-based UI to manage your containers, images, volumes, networks, and system usage — with support for **Docker** and **Podman**.
> 🧠 Name origin: In maritime terms, a **berth** is a designated place where a ship is docked — just like containers in your stack. Clean, organized, and under control.

</div>

## 📖 Table of Contents

- [✨ Overview](#-overview)
- [🚀 Installation](#-installation)
- [🧭 Usage](#-usage)
- [🛠️ Technology Stack](#-technology-stack)
- [📂 Project Structure](#-project-structure)
- [🤝 Contributing](#-contributing)
- [📜 License](#-license)

## ✨ Overview

Berth is a comprehensive terminal user interface (TUI) application built in Go, designed to simplify the management of Docker and Podman container environments. It provides a real-time, interactive experience for listing, inspecting, and controlling containers, images, volumes, and networks directly from your terminal. Berth aims to offer a `k9s`-like experience for container orchestration, focusing on usability, visual consistency, and efficient workflow.

## 🚀 Installation

### Prerequisites

-   [Go](https://golang.org/doc/install) (version 1.24 or higher recommended)
-   [Docker](https://docs.docker.com/get-docker/) or [Podman](https://podman.io/docs/installation) installed and running

### Steps

```bash
# 1. Clone the repository
git clone https://github.com/rluders/berth.git

# 2. Enter the project directory
cd berth

# 3. Build the binary
make build

# 4. Run it!
make run
```

## 🧭 Usage

Berth provides an intuitive keyboard-driven interface.

### 🎹 Navigation

*   `1` — Containers View
*   `2` — Images View
*   `3` — Volumes View
*   `4` — Networks View
*   `5` — System View

### 🛠️ Container Actions

| Key | Action                   |
| --- | ------------------------ |
| `s` | Start selected container |
| `x` | Stop selected container  |
| `d` | Remove container         |
| `l` | View logs                |
| `i` | Inspect container        |

### 📦 Image & Volume Actions

*   `d` — Remove selected image or volume

### 🧼 System Cleanup

| Key | Action           |
| --- | ---------------- |
| `b` | Basic Cleanup    |
| `a` | Advanced Cleanup |
| `t` | Total Cleanup    |

### 🔙 Back / Exit

*   `q` or `esc` — Return to the previous view or quit the application from the main views.

## 🛠️ Technology Stack

-   **Language**: [Go](https://golang.org/)
-   **TUI Framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea)
-   **Styling**: [Lipgloss](https://github.com/charmbracelet/lipgloss)
-   **Reusable TUI Components**: [Bubbles](https://github.com/charmbracelet/bubbles)

## 📂 Project Structure

```
.
├── cmd/                 # CLI entry point (e.g., main.go)
├── internal/
│   ├── tui/             # All Bubbletea models/views/components
│   ├── engine/          # Docker/Podman abstraction layer
│   ├── controller/      # Logic for container/image/volume actions
│   ├── state/           # Global state models
│   └── utils/           # Helpers: formatting, exec wrappers, etc.
├── assets/              # Logo, themes, maybe future plugins
├── go.mod
└── README.md
```

## 🤝 Contributing

We welcome contributions to Berth! If you're interested in improving the project, please consider:

-   Reporting bugs or suggesting features via [GitHub Issues](https://github.com/rluders/berth/issues).
-   Submitting pull requests for bug fixes or new features. Please ensure your code adheres to the existing style and includes appropriate tests.

## 📜 License

This project is licensed under the [MIT License](LICENSE).