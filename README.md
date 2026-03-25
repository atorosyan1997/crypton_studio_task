# Crypton Studio Task

**Author:** Arman Torosyan<br>
**Email:** armantorosyan1997@gmail.com<br>
**Telegram:** @atorosyan

## Description

Thread-safe map protected by a mutex with a concurrent test.

The map uses integer keys in the range from 1 to **1799** (birth year of Alexander Sergeyevich Pushkin).
Four goroutines concurrently access shuffled keys, incrementing each value.
After completion, every key holds the value **3**.

## Architecture

The project follows **Clean Architecture** principles — each layer has a single responsibility and depends only on inner layers through interfaces.

```
internal/
├── model/                        # Domain layer — data structures
│   └── safemap.go                # SafeMapStats (access/insert counters)
├── repository/
│   └── memory/                   # Data layer — in-memory storage
│       └── safemap.go            # Mutex-protected map implementation
└── service/                      # Business logic layer
    ├── safemap.go                # SafeMap service, counters, increment logic
    └── safemap_test.go           # Concurrent test + benchmarks
```

### Layers

| Layer | Package | Responsibility |
|-------|---------|----------------|
| **Model** | `internal/model` | Domain data structures shared across layers |
| **Repository** | `internal/repository/memory` | Data storage and synchronization (`sync.Mutex`) |
| **Service** | `internal/service` | Business logic: access/insert counting, increment operation |

The interface `SafeMapStorage` is defined in the **service** layer (consumer owns the interface) — a standard Go practice aligned with Clean Architecture.

## Technologies

| Technology | Purpose |
|------------|---------|
| **Go 1.25.5** | Programming language |
| **sync.Mutex** | Map synchronization in the repository layer |
| **sync/atomic** | Lock-free access/insert counters in the service layer |
| **testing** | Standard library for unit tests and benchmarks |

No external dependencies are used — the project relies solely on the Go standard library.

## Running Tests

```bash
# Unit test
go test ./internal/service/ -v

# Unit test with race detector
go test ./internal/service/ -v -race

# Benchmarks
go test ./internal/service/ -bench=. -benchmem
```

## Test Expectations

| Metric | Expected Value |
|--------|----------------|
| Value per key | 3 |
| Access count | 1799 * 3 = 5397 |
| Insert count | 1799 |
| Goroutines | 4 |
| Key access order | Shuffled (non-sequential) |
