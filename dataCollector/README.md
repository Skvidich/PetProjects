# Data Collector Microservice

The **Data Collector** is a core component of the UpOrDownDetector system. It is responsible for periodically fetching status data from external services, formatting it as `ServiceStatus` messages, and publishing this data to Kafka for downstream processing by other services.

---

## Overview

**Purpose**: Act as the system's ingestion layer, collecting live service status data for processing.

**Key Functions**:
- Scheduled polling of external service health endpoints (e.g. REST APIs)
- Transformation into structured `ServiceStatus` messages
- Publication to Kafka for further analysis and aggregation
- Runtime control and config reloading via gRPC API
- Optional storage of raw reports in PostgreSQL for traceability

---

## Architecture

### Data Flow

1. **Input**: List of monitored services and their endpoints (from config)
2. **Processing**:
   - Concurrent HTTP requests to fetch status data
   - Transformation into `ServiceStatus` format (defined via proto)
3. **Output**:
   - Kafka publication (`Statuses` topic)
   - PostgreSQL write (optional for raw storage)

---

## Technologies

### Core Components

| Component        | Purpose                             |
|------------------|-------------------------------------|
| **Go (Golang)**  | Main implementation language        |
| **Apache Kafka** | Event streaming target              |
| **PostgreSQL**   | Optional persistence for raw reports|
| **gRPC**         | Control interface (reload/start/stop) |
| **INI Files**    | Declarative configuration           |

### Go Libraries

| Library                       | Description                      |
|-------------------------------|----------------------------------|
| `github.com/IBM/sarama`       | Kafka producer client            |
| `github.com/jackc/pgx/v5`     | PostgreSQL driver                |
| `google.golang.org/grpc`      | gRPC server and API              |
| `gopkg.in/ini.v1`             | INI config parsing               |

---

## Code Structure

```
dataCollector/
├── cmd/                      # Entry point (main.go)
├── internal/
│   ├── core/
│   │   ├── coordinator       # Orchestrates fetch-storage-publish cycle
│   │   ├── getters           # HTTP-based data fetchers
│   │   ├── relay             # Kafka producer logic
│   │   └── storage           # PostgreSQL interactions
│   ├── app/server            # gRPC server setup
│   ├── config                # INI-based configuration loading
│   ├── grpcAPI               # gRPC API logic (handlers, storage)
│   └── logger                # Abstracted logging
├── pkg/types                 # Shared application-level types
├── go.mod / go.sum           # Dependency management
```

---

## Protobuf Interface

The service relies on shared `.proto` files from an external repository to maintain compatibility across services.

- **Repository**: [Skvidich/CollectorProto](https://github.com/Skvidich/CollectorProto)
- **Schema File**: `collector.proto`
- **Generated Code**: `gen/go/collector.pb.go`, `collector_grpc.pb.go`

---

## Configuration

Located in `app.ini`. Example format:

```ini
[Coordinator]
ReqDelay = 30s
Getters = Github,Discord,Cloudflare

[Relay]
Save = true
Resend = true

[Logger]
ErrLog = /app/code/logs/error.log

[Producer]
Brokers = kafka1:9092,kafka2:9093,kafka3:9094
Topic = Statuses

[Storage]
DSN = postgres://user:password@postgres-raw:5432/raw?sslmode=disable
ReportTable = raw_reports
ComponentTable = raw_component_metrics

[Server]
Address = data-collector:8888
```
## Run the Service

Ensure all dependencies (Kafka, PostgreSQL) are running and accessible.

```bash
go run ./cmd/main.go ./configs/app.ini
```

Logs will be written to the path specified in `ErrLog`.

---

## Data Model Summary

### Output: `ServiceStatus` (Kafka)

```json
{
  "name": "Discord",
  "time": "2025-04-10T12:00:00Z",
  "components": [
    { "name": "API", "status": "operational" },
    { "name": "Voice", "status": "degraded_performance" }
  ]
}
```

- Each message represents a snapshot of a service’s overall and component-level health.
- Messages are published to the Kafka topic `Statuses`.

### Output: Stored Data (PostgreSQL)

If enabled, the following tables are populated:

- `raw_reports` — JSON-encoded full status reports
- `raw_component_metrics` — Breakdown of component statuses per service

---

## License

This project is licensed under the [MIT License](../LICENSE) © 2025 [skvidich](https://github.com/skvidich)
