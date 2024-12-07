# Mini-Scan Data Processor

## Overview
The Mini-Scan project is designed to process internet scan results published to a Google Pub/Sub topic. It's part of a larger system intended to scale horizontally across thousands of machines, using distributed queues to manage data flow effectively. This component specifically pulls scan results, processes them, and maintains a database of unique `(ip, port, service)` combinations with their respective last scan time and service response.

## Architecture
The data processor subscribes to the `scan-sub` subscription, processes messages containing scan results, and stores them in a PostgreSQL database. The processor handles two data formats:
- `data_version: 1` with base64-encoded responses.
- `data_version: 2` with plain text responses.

## Technologies Used
- **Go**: Programming language for the processor.
- **Google Cloud Pub/Sub**: For message queue management.
- **PostgreSQL**: Database for storing processed data.
- **Docker**: Used to containerize and manage dependencies like the Google Pub/Sub emulator.

## Setup Instructions

### Prerequisites
- Go (1.15 or later)
- Docker and Docker Compose
- PostgreSQL database

### Configuration

## Set up Environment Variables:

```bash
export PUBSUB_EMULATOR_HOST=localhost:8085

## Configure the PostgreSQL Connection String in `processor.go`

```go
connStr := "postgres://username:password@localhost/dbname?sslmode=disable"

# PostgreSQL Configuration and Application Instructions

## Configure the PostgreSQL Connection String in `processor.go`

```go
connStr := "postgres://username:password@localhost/dbname?sslmode=disable"
```

Ensure this connection string is correctly configured for your local or production environment.

---

## Running the Application

### Start the Google Pub/Sub Emulator

```bash
docker-compose up
```

### Run the Processor

```bash
go run cmd/scanner/processor.go
```

---

## Testing Instructions

To test the application, follow these steps:

1. Ensure the Pub/Sub emulator is running.
2. Send a test message to the Pub/Sub topic:

    ```bash
    gcloud pubsub topics publish scan-topic --message='{"ip":"192.168.1.1","port":80,"service":"HTTP","timestamp":1609459200,"data_version":2,"data":{"response_str":"hello world"}}'
    ```

3. Observe the logs to ensure the message is processed and stored correctly.

---

## Data Storage

Data is stored in the `scan_records` table with the following schema:

```plaintext
Column       |  Type
-------------|---------
ip           | varchar
port         | integer
service      | varchar
last_scanned | timestamp
response     | text
```

Each record is uniquely identified by a combination of `ip`, `port`, and `service`.

---

## Contributing

Contributions to the Mini-Scan project are welcome. Please ensure you follow the existing code style and include tests for any new features or changes.

---

## License

Specify the license under which the project is made available.

---

## Contact Information

For help or issues, please submit a GitHub issue or contact [Your Name].

