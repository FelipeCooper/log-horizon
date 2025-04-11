# Log Horizon

An efficient log storage and retrieval service with support for compression and time-range queries.

## Overview

Log Horizon is an application that provides a gRPC API for:

1. **Storing logs** - Stores logs with automatic compression for large messages
2. **Querying logs** - Retrieves logs by time range and level
3. **Exporting logs** - Exports logs to files or provides log streaming

## Project Structure

The project follows a well-defined layered architecture:

```
.
├── app/                     # API Layer (gRPC)
│   ├── domain/              # APIs for specific domains
│   │   └── mlogapp/         # API for the logs domain
│   └── sdk/                 # Utilities for the API layer
│       ├── errs/            # Error handling
│       └── proto/           # Protobuf definitions
├── business/                # Business Layer
│   └── domain/              # Business domains
│       └── mlog/            # Logs domain
│           └── stores/      # Persistence interfaces
│               └── mongodb/ # MongoDB implementation
├── deploy/                  # Deployment configurations
│   └── proto/               # API definitions
└── foundation/              # Generic utilities
    ├── compress/            # Data compression
    ├── logger/              # Logging
    └── transaction/         # Transaction support
```

## Key Features

- **Efficient Storage**: Uses gzip compression to reduce storage space
- **Time Range Queries**: API to search logs in specific time periods
- **Log Streaming**: Support for retrieving logs in chunks for large datasets
- **File Export**: Capability to export logs to compressed files
- **Scalable Design**: Architecture based on domains and well-defined interfaces

## Technologies Used

- **Go**: Programming language
- **gRPC**: High-performance API framework
- **Protocol Buffers**: For data serialization
- **MongoDB**: Log storage
- **Gzip**: Data compression

## Getting Started

### Prerequisites

- Go 1.19+
- MongoDB
- Protoc (Protocol Buffers compiler)
- Docker & Docker Compose (optional, for containerized deployment)

### Installation

1. Clone the repository:

```bash
git clone https://github.com/felipecooper/log-horizon.git
cd log-horizon
```

2. Install dependencies:

```bash
go mod download
```

3. Generate Go files from protobuf:

```bash
protoc --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative app/sdk/proto/mlog/logs.proto
```

4. Run the server:

```bash
go run cmd/server/main.go
```

### Docker Deployment

For easy deployment, you can use Docker Compose:

```bash
docker-compose up -d
```

This will start both the MongoDB and Log Horizon application containers.

## MongoDB Configuration

### Local Setup

1. Install MongoDB:

   ```bash
   # Ubuntu
   sudo apt-get install mongodb

   # macOS
   brew install mongodb-community
   ```

2. Start MongoDB:

   ```bash
   # Ubuntu
   sudo systemctl start mongodb

   # macOS
   brew services start mongodb-community
   ```

3. Create a user (optional but recommended):

   ```bash
   mongosh
   > use admin
   > db.createUser({
       user: "root",
       pwd: "example",
       roles: [ { role: "root", db: "admin" } ]
     })
   ```

4. Configure environment variables for the application:
   ```bash
   export MONGODB_URI="mongodb://root:example@localhost:27017/admin"
   export MONGODB_DBNAME="loghorizon"
   export MONGODB_COLLECTION="logs"
   ```

### Docker Configuration

The Docker Compose file includes a pre-configured MongoDB instance. You can customize it by modifying the `docker-compose.yml` file:

```yaml
mongo:
  image: mongo:latest
  environment:
    MONGO_INITDB_ROOT_USERNAME: root
    MONGO_INITDB_ROOT_PASSWORD: example
  ports:
    - "27017:27017"
  volumes:
    - mongo-data:/data/db
```

## API Usage

For detailed API documentation, see [API.md](docs/API.md).

### Using the Client

The repository includes a client application that demonstrates how to use the API:

```bash
# Run the client to register a log
./client register "This is a test log message" info

# Search for logs in the last hour
./client search --start=-1h

# Export logs from the last week to a file
./client export --start=-168h --as-file
```

### Programmatic Usage

#### Initializing the Client

```go
package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	protomlog "github.com/felipecooper/log-horizon/app/sdk/proto/mlog"
)

func main() {
	// Connect to the gRPC server
	addr := "localhost:50051"
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

client := protomlog.NewLogWriterClient(conn)

    // Register a log
    resp, err := client.Register(context.Background(), &protomlog.NewLog{
        Message:   "Test log message",
        Level:     "info",
        Timestamp: time.Now().Unix(),
        Metadata:  map[string]string{"service": "example-service"},
    })
    if err != nil {
        log.Fatalf("Failed to register log: %v", err)
    }
    log.Printf("Log registered with ID: %s", resp.Id)
}
```

#### Registering a Log

```go
// Register a log with metadata
resp, err := writerClient.Register(context.Background(), &protomlog.NewLog{
	Message:   "Example log message",
	Level:     "info",
	Timestamp: time.Now().Unix(),
	Metadata: map[string]string{
		"service": "auth-service",
		"user_id": "12345",
		"request_id": "abc-123-xyz",
		"host": "api-server-1",
	},
})

if err != nil {
	log.Fatalf("Failed to register log: %v", err)
}

log.Printf("Log registered with ID: %s", resp.Id)
```

#### Searching Logs

```go
// Search logs from the last 24 hours with error level
startTime := time.Now().Add(-24 * time.Hour).Unix()
endTime := time.Now().Unix()

logs, err := readerClient.Search(context.Background(), &protomlog.SearchQuery{
	StartTime: startTime,
	EndTime:   endTime,
	Level:     "error",
	Page:      0,
	PageSize:  100,
})

if err != nil {
	log.Fatalf("Failed to search logs: %v", err)
}

log.Printf("Found %d logs (total: %d, has more: %v)", len(logs.Logs), logs.Total, logs.HasMore)

// Process the logs
for i, log := range logs.Logs {
	log.Printf("%d. [%s] %s (ID: %s)",
		i+1,
		time.Unix(log.Timestamp, 0).Format(time.RFC3339),
		log.Message,
		log.Id)

	// Print metadata if available
	if len(log.Metadata) > 0 {
		log.Println("   Metadata:")
		for k, v := range log.Metadata {
			log.Printf("   - %s: %s", k, v)
		}
	}
}
```

#### Exporting Logs to a File

```go
// Export logs from the last week to a file
startTime := time.Now().Add(-7 * 24 * time.Hour).Unix()
endTime := time.Now().Unix()

fileResp, err := readerClient.ExportToFile(context.Background(), &protomlog.SearchQuery{
	StartTime: startTime,
	EndTime:   endTime,
	AsFile:    true,
})

if err != nil {
	log.Fatalf("Failed to export logs: %v", err)
}

log.Printf("Logs exported to file: %s (size: %d bytes, compression: %s)",
	fileResp.FileUrl, fileResp.FileSize, fileResp.Compression)
```

#### Streaming Logs

```go
// Stream logs with a specific level
ctx := context.Background()
stream, err := readerClient.StreamFile(ctx, &protomlog.SearchQuery{
	StartTime: time.Now().Add(-24 * time.Hour).Unix(),
	EndTime:   time.Now().Unix(),
	Level:     "info",
	PageSize:  50,  // Smaller chunks for streaming
})

if err != nil {
	log.Fatalf("Failed to start streaming: %v", err)
}

// Process the stream of logs
for {
	batch, err := stream.Recv()
	if err == io.EOF {
		break  // End of stream
	}
	if err != nil {
		log.Fatalf("Error while streaming: %v", err)
	}

	// Process this batch of logs
	log.Printf("Received batch with %d logs", len(batch.Logs))

	// ... process logs ...
}
```

## Troubleshooting

### Common Issues

1. **Connection Refused**

   - Make sure the server is running and listening on the expected port
   - Check if there's a firewall blocking the connection

2. **Authentication Failed**

   - Verify MongoDB credentials in environment variables
   - Check that the MongoDB URI includes the authentication database (e.g., `/admin`)

3. **Empty Search Results**
   - Use broad time ranges for testing
   - Make sure you're connecting to the correct database and collection

### Debugging with gRPCurl

You can use [gRPCurl](https://github.com/fullstorydev/grpcurl) to test the gRPC service:

```bash
# List available services
grpcurl -plaintext localhost:50051 list

# List methods for a service
grpcurl -plaintext localhost:50051 list logs.LogReader

# Call the Search method
grpcurl -plaintext -d '{"start_time": 0, "end_time": 9999999999, "page_size": 10}' localhost:50051 logs.LogReader/Search
```

#Error Codes

Error Codes
The gRPC API uses standard gRPC error codes to indicate the status of operations. Below is a list of common error codes and their meanings:

Error Code Description

INVALID_ARGUMENT: The client provided invalid input, such as an unrecognized log level or invalid time range.

NOT_FOUND: The requested resource (e.g., log entry) was not found.
ALREADY_EXISTS The resource being created already exists.

PERMISSION_DENIED: The client does not have permission to perform the operation.

UNAUTHENTICATED: Authentication failed or was not provided.

RESOURCE_EXHAUSTED: The server has exhausted its resources (e.g., rate limits or storage).

INTERNAL: An internal server error occurred.

UNAVAILABLE: The service is currently unavailable (e.g., due to maintenance or overload).

DEADLINE_EXCEEDED: The operation took too long to complete and timed out.

##Service-Specific Errors

###LogWriter Service

Error Code: Scenario
INVALID_ARGUMENT: Log level is invalid or metadata is malformed.
INTERNAL: Failed to register the log due to a server-side issue.

###LogReader Service
Error Code: Scenario
INVALID_ARGUMENT: Time range is invalid or page size exceeds the limit.
NOT_FOUND: No logs found for the given query.
INTERNAL: Failed to retrieve logs due to a server-side issue.

###ExportToFile
Error Code: Scenario
INVALID_ARGUMENT: Time range or log level is invalid.
INTERNAL: Failed to export logs to a file due to a server-side issue.

###StreamFile
Error Code: Scenario
INVALID_ARGUMENT: Time range or log level is invalid.
INTERNAL: Failed to stream logs due to a server-side issue.
UNAVAILABLE: Streaming was interrupted due to server unavailability.

##How to Handle Errors
Check the Error Code: Use the error code to determine the type of issue.
Retry on Transient Errors: For errors like UNAVAILABLE or DEADLINE_EXCEEDED, implement retry logic with exponential backoff.
Fix Client-Side Issues: For errors like INVALID_ARGUMENT, ensure the request parameters are valid.
Contact Support: For persistent INTERNAL errors, contact the API support team.

## Quick-Start Guide

Provide a simple guide to help new clients integrate quickly.

Example:
Install Dependencies:

go get google.golang.org/grpc
go get google.golang.org/protobuf

Generate Protobuf Files:

protoc --go_out=. --go-grpc_out=. app/sdk/proto/mlog/logs.proto

Connect to the API:

conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))

Use the API:

client := protomlog.NewLogWriterClient(conn)

By implementing these suggestions, you can make your API more accessible and easier to integrate for new clients.

## Contributing

Contributions are welcome! Please read the contribution guidelines before submitting a pull request.

## License

This project is licensed under the [MIT License](LICENSE).
