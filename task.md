# Live Code Challenge: Task Management API with Golang, Fiber, and Advanced Concurrency  

**Objective**: Build a **Task Management API** using Golang and Fiber that demonstrates:

- CRUD operations
- Error handling with structured responses
- Concurrency patterns
- Context-aware request cancellation
- Middleware implementation

## **Core Requirements**

1. **Data Model**

   ```go
   type Task struct {
       ID          string    `json:"id" validate:"required,uuid4"`  
       Title       string    `json:"title" validate:"required,min=3"`  
       Description string    `json:"description,omitempty"`  
       Status      string    `json:"status" validate:"oneof=todo in_progress done"`  
       CreatedAt   time.Time `json:"created_at"`  
   }
   ```  

2. **In-Memory Storage**  
   - Use a thread-safe `map` with `sync.RWMutex` for concurrent access.  
   - Initialize a `Task` slice/map and a global `mutex`.

   ```go
   var (
       tasks   = make(map[string]Task)
       mutex   sync.RWMutex
   )
   ```  

3. **API Endpoints**  
   Implement the following routes with **Fiber**:  
   - `GET /tasks`: Fetch all tasks (use `mutex.RLock` for concurrent reads).  
   - `GET /tasks/:id`: Fetch a single task (return `404` if not found).  
   - `POST /tasks`: Create a task (validate input using `go-playground/validator`).  
   - `PUT /tasks/:id`: Update a task (validate input and handle concurrent writes with `mutex.Lock`).  
   - `DELETE /tasks/:id`: Delete a task.  

4. **Error Handling**  
   - Return **structured JSON errors** for all endpoints:  

     ```json
     { "error": "Task not found", "code": "TASK_404" }  
     ```

   - Use custom errors with `errors.Is` and `errors.As`:

     ```go
     var (  
         ErrTaskNotFound = errors.New("task not found")  
         ErrInvalidID    = errors.New("invalid ID format")  
     )  
     ```  

5. **Concurrency Patterns**  
   - **Goroutines & Channels**:  
     - Create a background goroutine to log task updates via a channel.

     ```go  
     func logTasks(taskChan <-chan Task) {  
         for task := range taskChan {  
             log.Printf("Task updated: ID=%s, Title=%s", task.ID, task.Title)  
         }  
     }  
     ```

   - **WaitGroup**:  
     - Ensure the logging goroutine completes before the app shuts down.

     ```go
     var wg sync.WaitGroup  
     defer func() {  
         close(taskChan)  
         wg.Wait()  
     }()  
     ```  

6. **Context Integration**  
   - Use `context.WithTimeout` in `GET /tasks` to handle request cancellation:

     ```go
     func getTasks(c *fiber.Ctx) error {  
         ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)  
         defer cancel()  
         
         // Simulate a long-running operation  
         select {  
         case <-ctx.Done():  
             return c.Status(fiber.StatusRequestTimeout).JSON(fiber.Map{"error": "request timed out"})  
         case <-time.After(3 * time.Second):  
             mutex.RLock()  
             defer mutex.RUnlock()  
             return c.JSON(tasks)  
         }  
     }  
     ```  

7. **Middleware**  
   - **Logging Middleware**: Log request method, path, and status code.  
   - **Mock Authentication**: Protect `POST`, `PUT`, and `DELETE` routes by checking the `X-API-Key` header.

   ```go  
   func authMiddleware(c *fiber.Ctx) error {  
       apiKey := c.Get("X-API-Key")  
       if apiKey != "12345" {  
           return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})  
       }  
       return c.Next()  
   }  
   ```  

---

## **Stretch Goals**

1. **Task Filtering**: Add query params to `GET /tasks` (e.g., `?status=done`).

2. **Worker Pool**: Process task updates using a pool of 3 workers.  
3. **Unit Tests**: Write tests for `GET /tasks` and `POST /tasks` using `testing` package.  

---

## **Starter Code**

```go  
package main  

import (  
    "context"  
    "errors"  
    "log"  
    "sync"  
    "time"  

    "github.com/go-playground/validator/v10"  
    "github.com/gofiber/fiber/v2"  
)  

type Task struct {  
    ID          string    `json:"id" validate:"required,uuid4"`  
    Title       string    `json:"title" validate:"required,min=3"`  
    Description string    `json:"description,omitempty"`  
    Status      string    `json:"status" validate:"oneof=todo in_progress done"`  
    CreatedAt   time.Time `json:"created_at"`  
}  

var (  
    tasks     = make(map[string]Task)  
    mutex     sync.RWMutex  
    validate  = validator.New()  
    taskChan  = make(chan Task)  
    wg        sync.WaitGroup  
)  

func main() {  
    app := fiber.New()  

    // Start logging goroutine  
    wg.Add(1)  
    go func() {  
        defer wg.Done()  
        logTasks(taskChan)  
    }()  

    // Add middleware  
    app.Use(loggingMiddleware)  
    app.Use("/tasks", authMiddleware)  

    // Implement routes here  

    defer func() {  
        close(taskChan)  
        wg.Wait()  
    }()  

    app.Listen(":3000")  
}  
```  

---

## **Evaluation Criteria**

1. **Correctness**: Endpoints behave as expected under concurrent access.  
2. **Concurrency**: Proper use of `mutex`, channels, and `WaitGroup`.  
3. **Error Handling**: Structured errors and context cancellation.  
4. **Code Quality**: Idiomatic Go, separation of concerns, and readability.  

This challenge tests core Go concepts while mimicking real-world scenarios like concurrent data access and request cancellation. Adjust the scope based on the developer’s pace! 🚀
