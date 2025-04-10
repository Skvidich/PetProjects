# UpOrDownDetector - Service Status Monitoring System

A distributed system for monitoring service statuses, detecting incidents, and analyzing service health metrics in real-time.

## Project Structure

```
UpOrDownDetector/
├── data_sql/                     # SQL initialization scripts  
├── dataCollector/                # Data collection microservice  
├── dataProcessor/                # Data processing microservice  
├── collector.Dockerfile          # Dockerfile for data collector  
├── processor.Dockerfile          # Dockerfile for data processor  
└── docker-compose.yaml           # Docker Compose configuration  
```

## System Components

### 1. [Data Collector Microservice](./dataCollector)
- Collects service status data via external HTTP APIs
- Provides a gRPC API for service management and control
- Sends raw service data to Kafka
- Stores data in PostgreSQL (raw dataset)

### 2. [Data Processor Microservice](./dataProcessor)
- Consumes raw data from Kafka
- Detects service incidents
- Aggregates health metrics
- Stores analyzed data in PostgreSQL (statistics DB)
- Sends alerts (incident start/end) via Kafka

### 3. Infrastructure Services
- Apache Kafka (3-node cluster) — distributed message brokering
- Zookeeper — Kafka coordination
- PostgreSQL — raw + processed data storage
- Kafka UI — real-time monitoring of Kafka topics
- PGAdmin — management UI for PostgreSQL

## Getting Started

### Prerequisites
- Docker
- Docker Compose

### Running the System

1. Clone the repository
   ```shell
   git clone https://github.com/your-repo/UpOrDownDetector.git
   cd UpOrDownDetector
   ```

2. Start all services
   ```shell
   docker-compose up -d
   ```

3. Access services
   - Kafka UI: http://localhost:8081
   - PGAdmin: http://localhost:8080
     - Email: admin@example.com
     - Password: admin
   - Data Collector API: gRPC on port 8888

4. Stop all services
   ```shell
   docker-compose down
   ```

## Monitoring & Logging

### Log Files

| Service         | Location           |
|----------------|--------------------|
| Data Collector | collector-logs/    |
| Data Processor | processor-logs/    |

### Kafka Topics

| Topic              | Description                 |
|--------------------|-----------------------------|
| Statuses           | Raw service status data     |
| Started_incidents  | Incident start alerts       |
| Ended_incidents    | Incident end alerts         |

### SQL Tables

#### Raw DB (postgres-raw)
- raw_reports — Raw status reports
- raw_component_metrics — Statuses of individual components

#### Statistic DB (postgres-statistic)
- reports — Aggregated service-level metrics
- component_metrics — Aggregated component metrics
- incidents — Incident log with duration and state

## Technologies

### Core Stack
- Golang — main language for all services
- PostgreSQL — relational DB for storing both raw and processed data
- Apache Kafka — high-throughput messaging pipeline
- Docker & Docker Compose — for containerization and orchestration

### Go Libraries

| Library                    | Purpose                          |
|----------------------------|----------------------------------|
| github.com/IBM/sarama      | Kafka client (producer/consumer) |
| github.com/jackc/pgx/v5    | PostgreSQL driver                |
| google.golang.org/grpc     | gRPC framework                   |
| gopkg.in/ini.v1            | INI file parsing                 |

## Usage Examples

```bash
# Check Kafka topics
docker exec -it kafka1 kafka-topics --bootstrap-server kafka1:9092 --list

# View collector logs
tail -f collector-logs/collector.log

# View processor logs
tail -f processor-logs/processor.log
```

## Environment Variables

| Variable                   | Description                   | Default              |
|----------------------------|-------------------------------|----------------------|
| POSTGRES_USER              | Username for both DBs         | user                 |
| POSTGRES_PASSWORD          | Password for both DBs         | password             |
| PGADMIN_DEFAULT_EMAIL      | Email for PGAdmin login       | admin@example.com    |
| PGADMIN_DEFAULT_PASSWORD   | Password for PGAdmin login    | admin                |


## Feedback

Open an issue or submit a PR if you'd like to contribute or suggest changes.

## License

This project is licensed under the [MIT License](./LICENSE) © 2025 [skvidich](https://github.com/skvidich)

