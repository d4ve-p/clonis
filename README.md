# 🚀 Clonis
**Clonis is a lightweight, self-hosted backup automation tool designed for transient VPS infrastructure.**


It streamlines the lifecycle of "VPS Hopping" (migrating between providers for cost optimization) by automating directory syncing to Google Drive via a minimal HTMX web dashboard. Built for reliability, low resource footprint, and zero-dependency deployments.

---

## 🚀 Quick Start

Get Clonis running in under 30 seconds.

1.  **Download the Recipe**
    ```bash
    curl -o docker-compose.yml [https://raw.githubusercontent.com/d4ve-p/clonis/main/docker-compose.yml](https://raw.githubusercontent.com/d4ve-p/clonis/main/docker-compose.yml)
    ```

2.  **Define Environment**
    Open the file and update the `SERVER_NAME` to identify this machine.
    ```bash
    nano docker-compose.yml
    ```

3.  **Ignition**
    Start the container in detached mode.
    ```bash
    docker compose up -d
    ```

4.  **Launch**
    Visit [http://localhost:8080](http://localhost:8080).
    > 💡 **Note:** Your backup configurations and database are safely persisted in the local `./config` directory.

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

1. **Initialization**: On startup, Clonis authenticates and establishes a Clonis/ root in the connected Drive (If user has made a connection previously).
2. **Server Namespace**: Creates a sub-directory based on the `SERVER_NAME` environment variable (e.g., `Clonis/DigitalOcean-SG1/`).
3. **Path Mapping**: Registered local paths are mirrored inside the namespace, preserving the directory structure.
    - *Local*: `/var/www/html`
    - *Remote*: `Clonis/DigitalOcean-SG1/var/www/html`
4. **Compressing**: Clonis compresses the mirrored paths into a single archive before uploading to Google Drive.
5. **Retention**: Clonis retains backups for a configurable period, deleting older versions automatically.

---
# 🔧 Tech Stack
- **Core**: Go (Golang)
- **Frontend**: HTMX + Go Templates (Server Side Rendering)
- **Containerization**: Docker (Multi-stage build based on Alpine Linux)
- **Storage**: Google Drive API v3