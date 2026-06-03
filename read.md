# Task Flow 🚀

`TaskFlow` is a concurrent, CLI-driven task pipeline engine built in Go. It implements a **Clean Architecture / DDD-lite** structure to process tasks through multiple lifecycles using a high-performance worker pool.

This project serves as a practical, production-grade sandbox for mastering Go's native concurrency primitives, memory management, and decoupled software design.

---

## 🏗️ Architecture

The project strictly follows decoupled layer boundaries to isolate core business rules from infrastructure details:

- **Domain Layer (`/internal/domain`)**: Pure Go enterprise logic. Contains the `Task` entity, `Status/Priority` custom types, and implementation-agnostic repository interfaces.
- **Application Layer (`/internal/application`)**: Use case orchestrator. Coordinates moving tasks through states and publishes events.
- **Infrastructure Layer (`/internal/infrastructure`)**:
  - `inmemory`: Thread-safe state storage utilizing a map wrapped in a `sync.RWMutex`.
  - `worker`: High-performance concurrency hub utilizing channels and worker goroutines.
- **Interface Layer (`/cmd`)**: CLI entrypoint handling configuration and standard I/O streams.

---

## 🛠️ Concurrency Patterns Demonstrated

- **Worker Pool**: Spawns $N$ configurable worker goroutines consuming jobs from a shared, bounded buffer.
- **Fan-Out / Fan-In**: A central dispatcher fans out jobs to the worker pool, while a collector goroutine fans in results from a single results channel.
- **Graceful Shutdown**: Utilizes `context.Context` propagation to signals workers to complete active jobs and stop processing cleanly on interrupt.
- **Data Race Prevention**: Protects volatile in-memory repositories using fine-grained read/write locking mechanics (`sync.RWMutex`).

---

## 🚀 Getting Started

### Prerequisites

- Go 1.21 or higher
- `make` installed on your machine

### 1. Installation & Setup

Clone the repository and install the development tools (including `air` for hot-reloading and linters):

```bash
make setup
```
