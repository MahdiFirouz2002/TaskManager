# Task Manager 

## Overview
This project is a concurrent task management system with logging, designed to handle multiple tasks efficiently using Goroutines and channels in Go. It includes features for creating, updating, retrieving, and deleting tasks, with built-in validation and thread-safe access.

## Features
- **Concurrent Task Logging**: Uses worker Goroutines to log task updates.
- **Thread-Safe Task Management**: Ensures safe access to shared task data using sync.Mutex.
- **REST API with Gin**: Provides HTTP endpoints for task operations.
- **Authentication Middleware**: Secures sensitive API endpoints using API key authentication.
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
| Method | Endpoint | Description | Authentication Required |
|--------|----------|-------------|-----------------------|
| `GET` | `/tasks` | Retrieve all tasks | No |
| `POST` | `/tasks` | Create a new task | Yes |
| `GET` | `/tasks/:id` | Get a task by ID | No |
| `PUT` | `/tasks/:id` | Update a task | Yes |
| `DELETE` | `/tasks/:id` | Delete a task | Yes |

### API Descriptions

#### Get All Tasks
- **Endpoint**: `GET /tasks`
- **Description**: Retrieves all tasks from the system. Supports optional filtering by status.
- **Authentication**: Not required.
- **Response**:
  - `200 OK`: Returns a list of tasks.
  - `204 No Content`: No tasks available.
  - `408 Request Timeout`: Request took too long.

#### Get Task by ID
- **Endpoint**: `GET /tasks/:id`
- **Description**: Retrieves a single task by its unique identifier.
- **Authentication**: Not required.
- **Response**:
  - `200 OK`: Returns the task details.
  - `404 Not Found`: Task does not exist.
  - `500 Internal Server Error`: Unexpected server error.

#### Create Task
- **Endpoint**: `POST /tasks`
- **Description**: Creates a new task.
- **Authentication**: Required (`X-API-Key: 12345` header).
- **Request Body**:
  ```json
  {
    "title": "Task Title",
    "description": "Task Description",
    "status": "pending"
  }
  ```
- **Response**:
  - `201 Created`: Task successfully created.
  - `400 Bad Request`: Invalid input.
  - `500 Internal Server Error`: Failed to create task.

#### Update Task
- **Endpoint**: `PUT /tasks/:id`
- **Description**: Updates an existing task.
- **Authentication**: Required (`X-API-Key: 12345` header).
- **Request Body**:
  ```json
  {
    "title": "Updated Title",
    "description": "Updated Description",
    "status": "completed"
  }
  ```
- **Response**:
  - `200 OK`: Task successfully updated.
  - `400 Bad Request`: Invalid task format.
  - `404 Not Found`: Task does not exist.
  - `500 Internal Server Error`: Failed to update task.

#### Delete Task
- **Endpoint**: `DELETE /tasks/:id`
- **Description**: Deletes a task by its ID.
- **Authentication**: Required (`X-API-Key: 12345` header).
- **Response**:
  - `200 OK`: Task deleted.
  - `404 Not Found`: Task does not exist.
  - `500 Internal Server Error`: Failed to delete task.

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

