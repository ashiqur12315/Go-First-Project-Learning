# Go Backend Learning Journey & Context

> **Note for AI Assistants / Agents:**  
> This file is the primary context and learning roadmap for this project. When continuing this conversation or working with another AI agent, read this file to understand the current architecture, mastered concepts, completed milestones, and upcoming topics. Keep this file updated as new features or concepts are implemented.

---

## 📌 Project Overview

- **Project Path**: `d:\Go Mama\go-first-project`
- **Module Name**: `go-first-project`
- **Go Version**: `go1.26.0` (compatible with Go 1.22+ standard library enhancements)
- **Dependencies**: Pure Go Standard Library (`net/http`, `encoding/json`, `sync`, `strconv`) — no external frameworks (like Gin or Fiber) to master the fundamentals first.
- **Port**: `8080` (`http://localhost:8080`)

---

## 🏗️ Current Architecture & Files

| File | Purpose |
| :--- | :--- |
| `main.go` | Main application containing the HTTP server, models, in-memory store, REST handlers, CORS middleware, and route registrations. |
| `main_test.go` | Unit and integration test suite using Go's `net/http/httptest` package covering all endpoints and CORS behavior. |
| `go.mod` | Go module definition (`module go-first-project`). |
| `LEARNING_JOURNEY.md` | Persistent learning log and context tracker for the developer and AI agents. |

---

## 💡 Concepts Mastered So Far

### 1. Go 1.22+ Routing (`http.ServeMux`)
- Standard library routing syntax that supports HTTP methods and URL path parameters:
  - `"GET /products"`
  - `"GET /products/{id}"`
  - Retrieve path parameter with `r.PathValue("id")`.

### 2. `mux.Handle` vs `mux.HandleFunc`
- **`mux.Handle(pattern, http.Handler)`**: Expects a type implementing the `http.Handler` interface (`ServeHTTP(ResponseWriter, *Request)`).
- **`http.HandlerFunc` Adapter**: Converts ordinary functions `func(w http.ResponseWriter, r *http.Request)` into an `http.Handler`.
- All routes are registered using:
  ```go
  mux.Handle("GET /products", http.HandlerFunc(getAllProducts))
  ```

### 3. REST API Methods
- **GET**: Retrieve all items (`GET /products`) or single item by ID (`GET /products/{id}`).
- **POST**: Create a new item (`POST /products`), returns status `201 Created`.
- **PUT**: Full replacement of an existing item (`PUT /products/{id}`).
- **PATCH**: Partial update (`PATCH /products/{id}`).
  - **Pointer pattern**: Uses `*string`, `*float64`, and `*int` in `UpdateProductInput` to distinguish between an omitted field (`nil`) versus an explicit zero-value (`""`, `0`).
- **DELETE**: Remove an item by ID (`DELETE /products/{id}`).

### 4. In-Memory Data Storage & Concurrency
- `sync.RWMutex`:
  - `mu.RLock()` / `mu.RUnlock()` for safe concurrent reads (`GET`).
  - `mu.Lock()` / `mu.Unlock()` for exclusive writes (`POST`, `PUT`, `PATCH`, `DELETE`).

### 5. CORS & Middleware
- **What is CORS**: Security mechanism enforced by browsers when web clients (e.g., React on `localhost:3000`) communicate with backends on a different port/origin.
- **Preflight Check (`OPTIONS`)**: Browsers send an automatic `OPTIONS` request before complex calls (`POST`, `PUT`, `PATCH` with JSON headers). Handled by returning `200 OK` immediately.
- **Middleware Pattern**: A function taking `http.Handler` and returning `http.Handler`:
  ```go
  func corsMiddleware(next http.Handler) http.Handler {
      middlewareHandler := func(w http.ResponseWriter, r *http.Request) {
          w.Header().Set("Access-Control-Allow-Origin", "*")
          w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
          w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

          if r.Method == http.MethodOptions {
              w.WriteHeader(http.StatusOK)
              return
          }

          next.ServeHTTP(w, r)
      }

      wrappedHandler := http.HandlerFunc(middlewareHandler)
      return wrappedHandler
  }
  ```
- **Three Ways to Apply CORS**:
  1. **Directly inside a single handler**: Setting headers on `w` (see `exampleDirectCORS`).
  2. **Per-route inline**: `mux.Handle("/url", corsMiddleware(handler))`.
  3. **Globally across all routes**: `http.ListenAndServe(":8080", corsMiddleware(mux))`.

### 6. Automated Testing with `httptest`
- Uses `net/http/httptest.NewRecorder()` and `httptest.NewRequest(...)` to test HTTP handlers in isolation without spinning up a live network socket.
- Run tests anytime with `go test -v .`.

---

## 📡 Registered Endpoints

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/hello` | Basic test endpoint returning `"Hello, World!"` |
| `GET` | `/products` | Returns all products as JSON array |
| `GET` | `/products/{id}` | Returns single product by integer ID (or 404) |
| `POST` | `/products` | Creates a new product (validates name & positive price) |
| `PUT` | `/products/{id}` | Replaces all fields of an existing product |
| `PATCH` | `/products/{id}` | Updates only provided fields (price, stock, etc.) |
| `DELETE` | `/products/{id}` | Deletes product by ID |
| `GET` | `/example-direct-cors` | Demonstrates CORS written directly inside handler |
| `GET` | `/example-inline-cors` | Demonstrates wrapping a single route with middleware |

---

## 🧪 Postman Quick Reference

- **Base URL**: `http://localhost:8080`
- **Headers**: `Content-Type: application/json`

### Example Request Bodies

#### POST `/products`
```json
{
  "name": "USB-C Hub",
  "description": "7-in-1 multi-port adapter with 4K HDMI",
  "price": 34.50,
  "stock": 40
}
```

#### PUT `/products/1` (Full Replace)
```json
{
  "name": "Mechanical Keyboard Pro",
  "description": "Hot-swappable tactile switches with sound dampening",
  "price": 89.99,
  "stock": 18
}
```

#### PATCH `/products/1` (Partial Update)
```json
{
  "price": 74.99,
  "stock": 12
}
```

---

## 🗺️ Next Steps Roadmap (Learning Suggestions)

When you're ready to proceed in your video tutorial or with any AI assistant, recommended next topics are:

1. **Persistent Database Integration**:
   - Replace in-memory slice with SQLite or PostgreSQL using Go's standard `database/sql` driver.
2. **Project Structure Refactoring**:
   - Separate code into standard Go packages (`cmd/server/main.go`, `internal/handlers`, `internal/models`, `internal/middleware`).
3. **Structured Request Validation**:
   - Input sanitization, checking string lengths, required fields, and returning consistent error JSON.
4. **Environment Variables**:
   - Reading PORT and configuration from `.env` or system environment (`os.Getenv`).
5. **Authentication & Authorization**:
   - Writing an Auth Middleware to check JWT bearer tokens or API keys.
6. **Logging Middleware**:
   - Creating a middleware that logs request method, path, status code, and latency (e.g. using `log/slog`).

---

## 📝 Learning Log / Session History

- **Session 1**:
  - Implemented Product REST API (`GET`, `POST`, `PUT`, `PATCH`, `DELETE`) with in-memory storage.
  - Learned difference between `PUT` (full replacement) and `PATCH` (pointer-based partial update).
  - Converted router from `mux.HandleFunc` to `mux.Handle` using `http.HandlerFunc`.
  - Added comprehensive unit test suite in `main_test.go`.
  - Implemented CORS middleware with preflight `OPTIONS` handling.
  - Refactored `corsMiddleware` with explicit step-by-step variables and clear final `return`.
  - Added demonstration endpoints for direct in-handler CORS vs. inline route-wrapped CORS.
  - Created `LEARNING_JOURNEY.md` to track progress and share state across AI agents.
