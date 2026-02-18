# 🚀 Concurrent TCP Chat Engine

A high-performance, multi-room capable TCP chat server built with **Golang**. This project demonstrates low-level networking, thread-safe memory management, and the power of Go's concurrency model (CSP).

## 🛠 Features

* **High Concurrency:** Handles hundreds of simultaneous connections using lightweight Goroutines.
* **Real-time Broadcasting:** Instant message delivery via asynchronous channels.
* **Thread-Safety:** Implements `sync.Mutex` to prevent race conditions during client registration.
* **Connection Lifecycle:** Graceful handling of client joins, nickname registration, and unexpected disconnects.
* **Timestamping:** All messages are logged with precise server-side time tracking.

## 🏗 System Architecture

The server operates on a **Single-Broadcaster, Multi-Worker** model:

1. **Listener:** Accepts raw TCP connections on port `8080`.
2. **Handler:** Each client gets a dedicated Goroutine for reading inputs.
3. **Broadcaster:** A central "hub" that manages the state and fans out messages to every client's private message channel.

## 🚀 Getting Started

### Prerequisites

* Go 1.18 or higher
* `telnet` or `netcat` (for testing)

### Installation

1. **Clone the repository:**
```bash
git clone https://github.com
cd tcp_chat-server
