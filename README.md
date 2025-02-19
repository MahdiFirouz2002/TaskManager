# Task Manager 

## Overview
This project is a concurrent task management system with logging, designed to handle multiple tasks efficiently using Goroutines and channels in Go. It includes features for creating, updating, retrieving, and deleting tasks, with built-in validation and thread-safe access.

## Features
- **Concurrent Task Logging**: Uses worker Goroutines to log task updates.
- **Thread-Safe Task Management**: Ensures safe access to shared task data using sync.Mutex.
- **REST API with Gin**: Provides HTTP endpoints for task operations.
- **Validation with go-playground/validator**: Ensures data integrity.
- **Unit and Load Testing**: Uses testify for assertions and parallel load testing.

## Installation

```sh
# Clone the repository
git clone https://github.com/your-username/concurrent-task-logger.git
cd concurrent-task-logger

# Install dependencies
go mod tidy
```

## Usage

### Start the Server
```sh
go run main.go
```

### API Endpoints
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/tasks` | Retrieve all tasks |
| `POST` | `/tasks` | Create a new task |
| `GET` | `/tasks/:id` | Get a task by ID |
| `PUT` | `/tasks/:id` | Update a task |
| `DELETE` | `/tasks/:id` | Delete a task |

## Running Tests
Run unit and load tests:
```sh
go test ./...
```

## Project Structure
```
/concurrentlogger  # Concurrent logging logic
/server           # API handlers & tests
/structs/task     # Task model & validation
/utils           # Custom error handling
```