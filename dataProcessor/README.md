# Data Processor Microservice

The **Data Processor** is a core component of the UpOrDownDetector system. It consumes real-time service status updates from Kafka, detects incidents, and performs time-window aggregation to produce structured health metrics. This data is then persisted in PostgreSQL and alerts are published back to Kafka for notification systems.

---

## Overview

**Purpose**: Analyze service behavior over time and produce structured, meaningful metrics.

**Key Functions**:
- Stream consumption via Kafka
- Stateful detection of incidents (start and end)
- Time-windowed aggregation of status data per service and component
- Persistent storage of metrics and incident history
- Asynchronous alerting via Kafka topics

---

## Architecture

### Data Flow

1. **Input**: Stream of raw `ServiceStatus` messages (JSON) from Kafka
2. **Processing**:
   - Track active incidents per service/component
   - Aggregate per-interval statistics for service health
3. **Output**:
   - Store incident history and metrics in PostgreSQL
   - Emit `incident started` and `incident ended` messages back to Kafka

---

## Technologies

### Core Components

| Component        | Purpose                             |
|------------------|-------------------------------------|
| **Go (Golang)**  | Main implementation language        |
| **Apache Kafka** | Event streaming & alert publishing  |
| **PostgreSQL**   | Persistent metrics and incident DB  |
| **Custom Queue** | Concurrent message queue buffering  |
| **gopkg.in/ini.v1** | Config management via INI files |

### Go Libraries

| Library                   | Description                      |
|---------------------------|----------------------------------|
| `github.com/IBM/sarama`   | Kafka producer/consumer client   |
| `github.com/jackc/pgx/v5` | PostgreSQL driver and toolkit    |
| `gopkg.in/ini.v1`         | INI config parsing               |

---

## Code Structure

```
dataProcessor/
├── cmd/                  # Entry points (main, cli)
├── configs/              # INI-based configuration
├── internal/
│   ├── config/           # Typed configuration structs
│   ├── core/
│   │   ├── alert/        # Kafka-based alerting
│   │   ├── processor/    # Metrics + incident engine
│   │   ├── reader/       # Kafka consumer and buffering
│   │   └── storage/      # PostgreSQL repository
│   └── logger/           # File-based error logging
├── pkg/
│   ├── models/           # ServiceStatus, Metrics, Incident
│   └── errors/           # Tagged error classification
└── logs/                 # Runtime log output
```

---

## Configuration

Located in `configs/app.ini`. Sample format:

```ini
[KafkaConsumer]
Brokers = kafka1:9092,kafka2:9093,kafka3:9094
Topic = Statuses

[KafkaReporter]
Brokers = kafka1:9092,kafka2:9093,kafka3:9094
StartTopic = Started_incidents
EndTopic = Ended_incidents

[PostgresStorage]
DSN = postgres://user:password@host:5432/statistic
ReportTable = reports
ComponentTable = component_metrics
IncidentTable = incidents

[ProcessEngine]
AggregationInterval = 5m
ReadTimeout = 3s

[]
ErrLog = logs/processor.log
RetryMax = 5
```

---

## Run the Service

Ensure all dependencies (Kafka, PostgreSQL) are running and accessible.

```bash
go run ./cmd/main.go ./configs/app.ini
```

Logs will be written to the path specified in `ErrLog`.

---

## Data Model Summary

### Input: `ServiceStatus` (Kafka)

```json
{
  "name": "Discord",
  "time": "2025-04-10T12:00:00Z",
  "components": [
    { "name": "API", "status": "operational" },
    { "name": "Voice", "status": "down" }
  ]
}
```

### Output: Stored Data

- `reports` — per-service aggregation over time window
- `component_metrics` — frequency of status per component
- `incidents` — incident start and end tracking per component

### Output: Alerts (Kafka)

- `Started_incidents`
- `Ended_incidents`

---

## License

This project is licensed under the [MIT License](../LICENSE) © 2025 [skvidich](https://github.com/skvidich)
