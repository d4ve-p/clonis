# 🚀 Clonis
**Clonis is a lightweight, self-hosted backup automation tool designed for transient VPS infrastructure.**

It streamlines the lifecycle of "VPS Hopping" (migrating between providers for cost optimization) by automating directory syncing to Google Drive via a minimal HTMX web dashboard. Built for reliability, low resource footprint, and zero-dependency deployments.
---

# ✨ Key Features
- **Infrastructure Agnostic**: Designed for ephemeral servers. Spin up, backup, destroy, restore.
- **Google Drive Integration**: Native API integration creates a structured, browsable backup hierarchy (`/Clonis/{Server_name}/{Path}`)
- **Resource Efficient**: Written in **Go** for minimal memory overhead (<50 MB RAM).
- **Web Dashboard**: Simple management UI built with **HTMX** (No heavy React/Node bundles)
- **Docker Ready**: One-click deployment via docker compose
---

# 🛠 Architecture & Logic
Clonis enforces a strict "Server-Root" hierarchy to prevent data collisions when managing multiple VPS instances.

1. **Initialization**: On startup, Clonis authenticates and establishes a Clonis/ root in the connected Drive.
2. **Server Namespace**: Creates a sub-directory based on the `SERVER_NAME` environment variable (e.g., `Clonis/DigitalOcean-SG1/`).
3. **Path Mapping**: Registered local paths are mirrored inside the namespace, preserving the directory structure.
    - *Local*: `/var/www/html`
    - *Remote*: `Clonis/DigitalOcean-SG1/var/www/html`
---
# 🚀 Quick Start
**Prerequisites**
  - Docker & Docker Compose
  - Google Cloud Service Account Credentials (`credentials.json`)

**Installation**
1. **Clone the repository**
    ```Bash
    git clone https://github.com/d4ve-p/clonis.git
    cd clonis
    ```
2. **Configure Environment**
    ```Bash
    cp .env.example .env
    # Edit .env to set SERVER_NAME and PORT
    ```
3. **Run with Docker**
    ```Bash
    docker-compose up -d

    Access the Dashboard Visit http://localhost:3130 to manage backup paths.
    ```

# 🔧 Tech Stack
- **Core**: Go (Golang)
- **Frontend**: HTMX + Go Templates (Server Side Rendering)
- **Containerization**: Docker (Multi-stage build based on Alpine Linux)
- **Storage**: Google Drive API v3