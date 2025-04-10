# UpOrDownDetector - Service Status Monitoring System

A distributed system for monitoring service statuses, detecting incidents, and analyzing service health metrics.

## Project Structure

<pre>
UpOrDownDetector/  
├── collector-logs/               # Log files for data collector  
├── data_sql/ # SQL initialization scripts  
├── dataCollector/ # Data collection microservice  
├── dataProcessor/ # Data processing microservice  
├── processor-logs/ # Log files for data processor  
├── collector.Dockerfile # Dockerfile for data collector  
├── processor.Dockerfile # Dockerfile for data processor  
└── docker-compose.yaml # Docker Compose configuration  
</pre>

## System Components

### 1. Data Collector Microservice ([details](./dataCollector))
- Collects service status data via HTTP APIs
- Provides gRPC API for system control
- Sends raw data to Kafka
- Stores data in PostgreSQL

### 2. Data Processor Microservice ([details](./dataProcessor))
- Consumes data from Kafka
- Detects service incidents
- Generates service health metrics
- Stores analyzed data in PostgreSQL
- Sends alerts via Kafka

### 3. Infrastructure Services
- **Kafka Cluster**: 3-node cluster for message brokering
- **PostgreSQL**: 
  - Raw data storage ([config](./data_sql/raw.sql))
  - Processed statistics storage ([config](./data_sql/stat.sql))
- **Zookeeper**: For Kafka coordination
- **Kafka UI**: Web interface for cluster monitoring
- **PGAdmin**: Web interface for database management

## Getting Started

### Running the System
1. Clone the repository:
    ```bash
   git clone https://github.com/your-repo/UpOrDownDetector.git
   cd UpOrDownDetector

2. Start all services:
    ```bash
    docker-compose up -d
3. Access services:
    - Kafka UI: http://localhost:8081
    - PGAdmin: http://localhost:8080 (admin@example.com/admin)
    - Data Collector API: gRPC on port 8888
      
4. Stopping the System:
    ```bash
    docker-compose down
    
## Monitoring and Logging

### Log files:
  - Data Collector: collector-logs/
  - Data Processor: processor-logs/
    
### Kafka Topics:
  - Statuses: Raw service status data
  - Started_incidents: Incident start notifications
  - Ended_incidents: Incident end notifications
    
### SQL tables
  - raw_reports - Raw service status data
  - raw_component_metrics - Raw component-level status data
  - reports - Aggregated metrics
  - component_metrics - Component-level status metrics
  - incidents - Tracked incidents

## Technologies and Libraries

### Core Technologies
  - Golang (Primary programming language)
  - PostgreSQL (Relational database for both raw and processed data)
  - Apache Kafka  (Distributed event streaming platform)
  - Docker (Containerization)
  - Docker Compose (Orchestration)
### Go Libraries
  - github.com/IBM/sarama	Kafka client library	
  - github.com/jackc/pgx/v5	PostgreSQL driver	
  - google.golang.org/grpc	gRPC framework	
  - gopkg.in/ini.v1	INI configuration parsing	
