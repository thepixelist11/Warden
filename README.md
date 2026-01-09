# Warden

Warden is a Go-based command-line tool for querying and monitoring Minecraft servers, providing detailed server information and an optional live terminal-based monitor with player lists and server icon visualization. Supports both Java and Bedrock editions.

---

## Features

- Query Minecraft server status via [mcsrvstat.us API](https://mcsrvstat.us/).
- Supports Java and Bedrock servers.
- Displays server information in JSON format.
- Optional live monitoring mode with:
  - Server icon display (if available)
  - Hostname, version, IP:Port, and online status
  - MOTD (Message of the Day)
  - Player list with online count
- Optional output to a JSON file.
- Silent mode to suppress console output.
- Debug mode using pre-saved JSON files.
- Calculates the primary color of server icons for UI styling.

---

## Installation

Ensure you have Go installed (version 1.20+ recommended).

```bash
git clone https://github.com/thepixelist11/Warden.git
cd Warden
go mod tidy
go build -o warden
