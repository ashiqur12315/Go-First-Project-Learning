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

| File / Package | Purpose |
| :--- | :--- |
| `main.go` | Ultra-clean entry point (`package main`): 5 lines of code that simply calls `cmd.Serve()`. |
| `cmd/serve.go` | `package cmd`: Configures the `ServeMux`, registers routes from `handlers`, wraps with `middleware.CorsMiddleware`, and runs `http.ListenAndServe`. |
| `models/product.go` | `package models`: Exported data models (`Product`, `UpdateProductInput`) and thread-safe in-memory store (`Products`, `Mu`, `NextID`). |
| `handlers/helpers.go` | `package handlers`: Package-private helper `writeJSON` (lowercase 'w') used by handlers. |
| `handlers/get_product.go` | `package handlers`: Exported query handlers `GetAllProducts` and `GetProductByID`. |
| `handlers/create_product.go` | `package handlers`: Exported handler `CreateProduct`. |
| `handlers/update_product.go` | `package handlers`: Exported handlers `UpdateProduct` and `PatchProduct`. |
| `handlers/delete_product.go` | `package handlers`: Exported handler `DeleteProduct`. |
| `middleware/cors.go` | `package middleware`: Exported `CorsMiddleware` with preflight `OPTIONS` and caching. |
| `middleware/logger.go` | `package middleware`: Exported `LoggerMiddleware` with `responseRecorder` capturing method, path, status, and latency. |
| `main_test.go` | Integration and unit tests covering all handlers and middleware. |
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
- **Why Wrapping `mux` Works (The Key Insight)**:
  - In Go, `http.ServeMux` **implements `http.Handler`**!
  - Therefore, `corsMiddleware(mux)` passes `mux` in as the `next` argument.
  - When a request arrives, Go executes `corsMiddleware`, which executes `w.Header().Set(...)` and then calls `next.ServeHTTP(w, r)`.
  - Since `next == mux`, it calls `mux.ServeHTTP(w, r)`! `mux` then matches the URL pattern (e.g. `GET /products`) and forwards `w` and `r` to `getAllProducts(w, r)`.

### 6. Automated Testing with `httptest`
- Uses `net/http/httptest.NewRecorder()` and `httptest.NewRequest(...)` to test HTTP handlers in isolation without spinning up a live network socket.
- Run tests anytime with `go test -v .`.

### 7. Go Packages & Visibility Rules (Exported vs. Unexported)
- **Directory = Package**: In Go, each directory represents a package (e.g., `models/` is `package models`, `handlers/` is `package handlers`).
- **Capitalization Rule**:
  - **Uppercase (PascalCase)**: **Exported** (public, accessible by other packages). E.g., `handlers.GetAllProducts`, `models.Product`, `middleware.CorsMiddleware`.
  - **Lowercase (camelCase)**: **Unexported** (package-private). E.g., `handlers.writeJSON` is only usable by files inside the `handlers` folder.
- **Importing Local Packages**:
  - Imported using module name prefix from `go.mod`:
    ```go
    import (
        "go-first-project/handlers"
        "go-first-project/middleware"
        "go-first-project/models"
    )
    ```

### 8. Logger Middleware & Middleware Chaining
- **Response Recorder Wrapper**: In Go, `http.ResponseWriter` doesn't expose the status code after writing it. We create a `responseRecorder` struct embedding `http.ResponseWriter` that intercepts `WriteHeader(code)` to record `statusCode`.
- **Latency Measurement**: Captures `start := time.Now()` before calling `next.ServeHTTP(rec, r)` and `duration := time.Since(start)` after it finishes.
- **Middleware Chaining**: Middlewares wrap each other like onion layers:
  ```go
  rootHandler := middleware.LoggerMiddleware(middleware.CorsMiddleware(mux))
  ```
  Request pipeline:
  $$\text{Request} \longrightarrow \mathbf{Logger} \longrightarrow \mathbf{CORS} \longrightarrow \mathbf{ServeMux} \longrightarrow \mathbf{Handler}$$
  Response pipeline:
  $$\text{Response} \longleftarrow \mathbf{Logger\ (logs\ latency\ \&\ status)} \longleftarrow \mathbf{CORS} \longleftarrow \mathbf{Handler}$$

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
- **Session 2**:
  - Deep-dive into Global Middleware mechanics and Go's `http.Handler` abstraction.
  - Demystified `mux` as an `http.Handler`: discovered that `http.ServeMux` implements `ServeHTTP(w, r)`.
  - Traced the complete line-by-line runtime lifecycle of a request (`GET /products`) through `corsMiddleware` into `mux` and finally `getAllProducts`.
- **Session 3**:
  - Verified and explored the `OPTIONS` preflight request lifecycle in `corsMiddleware`.
  - Confirmed why `next.ServeHTTP(w, r)` is intentionally skipped for `OPTIONS` (answering the browser preflight security gate immediately).
  - Added `Access-Control-Max-Age: 86400` to cache preflight permission for 24 hours, saving unnecessary extra network roundtrips.
- **Session 4**:
  - Cleaned up demonstration endpoints (`/example-direct-cors` and `/example-inline-cors`) from `main.go` and `main_test.go`.
  - Verified that all core product REST endpoints and CORS preflight tests continue to pass with 100% success.
- **Session 5**:
  - Renamed `handlerWithCORS` to the idiomatic `rootHandler` in `main.go`.
  - Discussed standard Go naming conventions for root HTTP handlers.
- **Session 6**:
  - Destructured and modularized the monolithic `main.go` into dedicated files: `models.go`, `helpers.go`, `middleware.go`, `get_product.go`, `create_product.go`, `update_product.go`, `delete_product.go`, and a clean, concise `main.go`.
  - Highlighted the critical Go CLI tip: use `go run .` instead of `go run main.go` when working with multiple files in `package main`.
- **Session 7**:
  - Migrated from a single flat package to multi-package architecture: `models/`, `handlers/`, and `middleware/`.
  - Learned Go's visibility rules (Uppercase = Exported/Public, Lowercase = Unexported/Package-private).
  - Learned how to import subpackages using the module prefix (`"go-first-project/..."`).
  - Updated all unit tests in `main_test.go` to test across the new packages.
- **Session 8**:
  - Introduced the standard `cmd/` package pattern with `cmd/serve.go`.
  - Moved server initialization and routing logic into `cmd.Serve()`.
  - Reduced `main.go` to an ultra-clean 5-line entry point that delegates execution to `cmd.Serve()`.
- **Session 9**:
  - Implemented `middleware/logger.go` with a custom `responseRecorder` to capture HTTP status codes and measure request latency.
  - Chained multiple middlewares (`LoggerMiddleware` wrapping `CorsMiddleware` wrapping `mux`).
  - Verified terminal logging output across all unit tests and endpoints.
