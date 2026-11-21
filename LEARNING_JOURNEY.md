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
| `cmd/serve.go` | `package cmd`: Orchestrates application startup: loads `config.Load()`, prints startup banner, and passes `cfg` to `rest.Start(cfg)`. |
| `config/config.go` | `package config`: Pure standard-library configuration manager — reads `.env` via `bufio.Scanner`, sets `os.Setenv`, and populates `Config` struct (`Version`, `ServiceName`, `Port`) with fallback defaults. |
| `config/config_test.go` | `package config`: Unit test suite verifying `.env` parsing, comments, quotes, and fallback defaults. |
| `.env` | Environment configuration file containing `VERSION`, `SERVICE_NAME`, and `PORT`. |
| `models/product.go` | `package models`: Pure domain entities and schema definitions (`Product`, `UpdateProductInput`), strictly decoupled from storage logic. |
| `database/database.go` | `package database`: Dedicated data layer with in-memory thread-safe repository CRUD functions (`GetAllProducts`, `GetProductByID`, `CreateProduct`, `UpdateProduct`, `PatchProduct`, `DeleteProduct`, and `Reset`). |
| `database/database_test.go` | `package database`: Complete unit test suite verifying CRUD operations, sentinel errors, and test isolation. |
| `utils/response.go` | `package utils`: Reusable cross-package utility functions including `SendData(w, status, data)` for standardized JSON responses. |
| `utils/response_test.go` | `package utils`: Unit test suite verifying `SendData` headers, status code, and JSON payload encoding. |
| `rest/server.go` | `package rest`: Encapsulates REST server setup — accepts `cfg *config.Config`, initializes `mux` and `manager`, calls `RegisterRoutes(mux, manager)`, wraps with `manager.Wrap(mux)`, and starts `http.ListenAndServe`. |
| `rest/routes.go` | `package rest`: Single source of truth for routing — defines `RegisterRoutes(mux, manager)` and maps URL paths/methods to handler functions. |
| `rest/handlers/get_product.go` | `package handlers`: Exported query handlers `GetAllProducts` and `GetProductByID` using `database.*` and `utils.SendData`. |
| `rest/handlers/create_product.go` | `package handlers`: Exported handler `CreateProduct` with validation, database insertion, and `utils.SendData`. |
| `rest/handlers/update_product.go` | `package handlers`: Exported handlers `UpdateProduct` and `PatchProduct` using `database.*` and `utils.SendData`. |
| `rest/handlers/delete_product.go` | `package handlers`: Exported handler `DeleteProduct` using `database.*` and `utils.SendData`. |
| `rest/middleware/cors.go` | `package middleware`: Exported `CorsMiddleware` with preflight `OPTIONS` and caching. |
| `rest/middleware/logger.go` | `package middleware`: Exported `LoggerMiddleware` with `responseRecorder` capturing method, path, status, and latency. |
| `rest/middleware/manager.go` | `package middleware`: Exported `Manager` and `Middleware` type to chain middlewares cleanly without nesting. |
| `main_test.go` | Integration and unit tests covering all handlers and middleware via `rest.RegisterRoutes`. |
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
  mux.Handle("GET /products", http.HandlerFunc(handlers.GetAllProducts))
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
  func CorsMiddleware(next http.Handler) http.Handler {
      middlewareHandler := func(w http.ResponseWriter, r *http.Request) {
          w.Header().Set("Access-Control-Allow-Origin", "*")
          w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
          w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
          w.Header().Set("Access-Control-Max-Age", "86400")

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
  1. **Directly inside a single handler**: Setting headers directly on `w`.
  2. **Per-route inline**: `mux.Handle("/url", CorsMiddleware(handler))`.
  3. **Globally across all routes**: `rootHandler := CorsMiddleware(mux)`.
- **Why Wrapping `mux` Works (The Key Insight)**:
  - In Go, `http.ServeMux` **implements `http.Handler`**!
  - Therefore, `CorsMiddleware(mux)` passes `mux` in as the `next` argument.
  - When a request arrives, Go executes `CorsMiddleware`, which sets headers and calls `next.ServeHTTP(w, r)`.
  - Since `next == mux`, it calls `mux.ServeHTTP(w, r)`! `mux` matches the route (e.g. `GET /products`) and forwards `w` and `r` to `handlers.GetAllProducts(w, r)`.

### 6. Automated Testing with `httptest`
- Uses `net/http/httptest.NewRecorder()` and `httptest.NewRequest(...)` to test HTTP handlers in isolation without spinning up a live network socket.
- Run tests across all packages anytime with:
  ```powershell
  go test -v ./...
  ```

### 7. Go Packages & Visibility Rules (Exported vs. Unexported)
- **Directory = Package**: In Go, each directory represents a package (e.g., `models/` is `package models`, `handlers/` is `package handlers`).
- **Capitalization Rule**:
  - **Uppercase (PascalCase)**: **Exported** (public, accessible by other packages). E.g., `handlers.GetAllProducts`, `models.Product`, `middleware.CorsMiddleware`.
  - **Lowercase (camelCase)**: **Unexported** (package-private). E.g., `handlers.sendData` is only usable by files inside the `handlers` folder.
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

### 9. Middleware Manager & Clean Chaining (Solving the Onion Pyramid)
- **The Problem**: Wrapping multiple middlewares creates deeply nested code:
  ```go
  rootHandler := middleware.Logger(middleware.Cors(middleware.Auth(middleware.RateLimit(mux))))
  ```
- **The Type Definition**: Every standard Go middleware matches `func(http.Handler) http.Handler`, so we define:
  ```go
  type Middleware func(http.Handler) http.Handler
  ```
- **The Manager Pattern**: A `Manager` struct holds a slice `[]Middleware`.
- **The Reverse Loop in `Wrap()`**:
  ```go
  func (m *Manager) Wrap(handler http.Handler) http.Handler {
      for i := len(m.middlewares) - 1; i >= 0; i-- {
          handler = m.middlewares[i](handler)
      }
      return handler
  }
  ```
  By iterating backwards, the first middleware passed to `manager.Use(m1, m2)` becomes the outermost wrapper, meaning incoming requests hit `m1` first, then `m2`, then the router.
- **Clean Configuration**:
  ```go
  manager := middleware.NewManager()
  manager.Use(
      middleware.LoggerMiddleware,
      middleware.CorsMiddleware,
  )
  rootHandler := manager.Wrap(mux)
  ```

### 10. Dedicated Route Registry (`cmd/routes.go`)
- **Separation of Concerns**: `cmd/serve.go` manages server lifecycle, port listening, and global middleware configuration. `cmd/routes.go` is the single source of truth for URL path declarations and handler mappings.
- **Passing Mux and Manager as Input**:
  ```go
  func RegisterRoutes(mux *http.ServeMux, manager *middleware.Manager) {
      mux.Handle("GET /products", http.HandlerFunc(handlers.GetAllProducts))
      // ...
  }
  ```
  Passing both `mux` and `manager` gives route definitions full access to the router while also allowing route-specific middlewares to be applied via `manager` when needed.
- **Single Source of Truth for Tests**: Integration tests in `main_test.go` call `cmd.RegisterRoutes(mux, manager)` directly, eliminating duplicated endpoint definitions.

### 11. Pure Standard-Library Config & `.env` Parser (Zero 3rd-Party Dependencies)
- **Why No External Packages?**: Go's `os`, `bufio`, and `strings` packages are sufficient to parse `.env` files and manage environment configuration without adding external dependencies (like `godotenv`).
- **Parsing Mechanics (`bufio.Scanner`)**:
  - `os.Open(".env")` opens the file safely; if missing, returns `nil` so production/container environments relying on system environment variables continue working.
  - `scanner.Scan()` reads line-by-line, stripping whitespace and skipping blank lines and comments (`#`).
  - `strings.SplitN(line, "=", 2)` cleanly separates keys and values, stripping optional enclosing quotes (`"` or `'`).
  - `os.Setenv(key, value)` sets keys into the runtime environment only if not already provided by the host OS.
- **Strongly Typed `Config` Struct**:
  - `Config` exposes `Version`, `ServiceName`, and `Port`.
  - Uses `getEnv(key, fallback)` helper to guarantee safe default values if `.env` keys are omitted.

### 12. REST Layer Encapsulation (`rest/` package)
- **Architectural Motivation**: Separates transport protocols (REST/HTTP) from business logic, CLI commands, and database models. This prepares the application to support other transports in the future (like gRPC or WebSocket) without modifying core domain logic.
- **The `rest/server.go` Entry Point**:
  ```go
  func Start(cfg *config.Config) { ... }
  ```
  Encapsulates mux creation, middleware configuration, route binding, and HTTP listener binding.
- **Passing Configuration Down**:
  `cmd/serve.go` simply calls `rest.Start(cfg)`, keeping CLI/cmd code decoupled from HTTP server implementation details.

### 13. Reusable Utility Functions (`utils/` package & `SendData`)
- **Centralizing Cross-Cutting Helpers**: Utility functions like sending JSON responses (`utils.SendData`) are separated into a dedicated `utils/` package.
- **Exported `utils.SendData`**:
  ```go
  func SendData(w http.ResponseWriter, status int, data any)
  ```
  Sets `Content-Type: application/json`, writes status code, and encodes JSON payload cleanly.
- **Handler Integration**: Handlers import `"go-first-project/utils"` and call `utils.SendData(w, status, data)` directly, eliminating unnecessary wrapper functions.

### 14. Data Layer Abstraction & In-Memory Repository Pattern (`database/` package)
- **Separation of Concerns (SoC)**:
  - `models/`: Strictly defines schema entities (`Product`, `UpdateProductInput`). No storage variables, no mutexes, no database imports.
  - `database/`: Manages storage state, concurrency synchronization (`sync.RWMutex`), query filtering, and data mutations.
  - `handlers/`: Handles HTTP request decoding, validation, delegating to `database.*`, and returning JSON responses via `sendData`.
- **Clean CRUD Operations**:
  ```go
  GetAllProducts() []models.Product
  GetProductByID(id int) (*models.Product, error)
  CreateProduct(p models.Product) models.Product
  UpdateProduct(id int, p models.Product) (*models.Product, error)
  PatchProduct(id int, input models.UpdateProductInput) (*models.Product, error)
  DeleteProduct(id int) error
  ```
- **Sentinel Errors**:
  - `var ErrNotFound = errors.New("product not found")` exported by `database`.
  - Handlers use `errors.Is(err, database.ErrNotFound)` to cleanly distinguish 404 Not Found from 500 Internal Server Error.
- **Automated Test Isolation**:
  - `database.Reset()` restores the initial seed data.
  - Calling `database.Reset()` in test setup functions prevents cross-test state leakage.

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
   - Replace in-memory slice with SQLite or PostgreSQL using Go's standard `database/sql` driver (or GORM / sqlc / pgx).
2. **Structured Request Validation**:
   - Validating incoming payloads (e.g., checking minimum name length, positive price, valid email/data formats) and returning standardized error JSON structures.
3. **Authentication & Authorization**:
   - Writing an Auth Middleware to inspect `Authorization: Bearer <token>` headers using JWT.
4. **Graceful Shutdown**:
   - Listening for OS interrupt signals (`os.Interrupt`, `syscall.SIGTERM`) to cleanly close server listeners and database connections using `server.Shutdown(ctx)`.
5. **Dockerizing the Application**:
   - Writing a multi-stage `Dockerfile` to create a tiny, production-ready Go binary container.

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
- **Session 10**:
  - Solved the nested middleware "pyramid of doom" by implementing a Middleware Manager pattern in `middleware/manager.go`.
  - Defined `type Middleware func(http.Handler) http.Handler`.
  - Implemented `Manager` struct with `Use(middlewares...)` and `Wrap(handler)`.
  - Learned why iterating backwards through the middleware slice ensures the first registered middleware executes first for incoming requests.
  - Refactored `cmd/serve.go` and `main_test.go` to use `manager.Wrap(mux)`.
  - Added dedicated unit test `TestMiddlewareManagerExecutionOrder` verifying exact start/end lifecycle order.
- **Session 11**:
  - Extracted all route declarations from `cmd/serve.go` into a dedicated `cmd/routes.go` file.
  - Defined `RegisterRoutes(mux *http.ServeMux, manager *middleware.Manager)` accepting both router and middleware manager as input.
  - Refactored `cmd/serve.go` to instantiate `mux` and `manager`, configure global middlewares, and pass them to `RegisterRoutes(mux, manager)`.
  - Refactored `main_test.go` to use `cmd.RegisterRoutes(mux, manager)`, keeping routes and tests DRY (Don't Repeat Yourself).
- **Session 12**:
  - Implemented pure Go `.env` file parsing and environment management in a new `config/` package without any 3rd-party dependencies.
  - Defined `Config` struct (`Version`, `ServiceName`, `Port`) with fallback defaults via `config.Load()`.
  - Added startup configuration banner in `cmd/serve.go` printing service name, version, and dynamically binding server to configured `PORT`.
  - Added automated test suite in `config/config_test.go` covering `.env` loading, comment handling, quote stripping, and fallback defaults (total 12 unit tests passing).
- **Session 13**:
  - Re-architected project under REST domain philosophy: created dedicated `rest/` directory.
  - Moved `handlers/` and `middleware/` into `rest/handlers/` and `rest/middleware/`.
  - Created `rest/routes.go` and `rest/server.go` with `rest.Start(cfg *config.Config)`.
  - Updated `cmd/serve.go` to import `rest` and call `rest.Start(cfg)`.
  - Updated `main_test.go` imports and verified all 12 tests continue to pass with 100% success.
- **Session 14**:
  - Created dedicated `utils/` package with `utils.SendData(w, status, data)` for standardized JSON response generation across all layers.
  - Added unit test suite `utils/response_test.go` verifying Content-Type header, status code, and JSON encoding.
  - Refactored `rest/handlers/helpers.go` to bridge `sendData` to `utils.SendData`.
  - Refactored all handlers (`get_product`, `create_product`, `update_product`, `delete_product`) to use `sendData`.
  - Verified all 13 automated tests across `main`, `config`, and `utils` pass with 100% success.
- **Session 15**:
  - Decoupled storage from domain models by creating a dedicated `database/` package implementing Approach 2 (Repository CRUD functions).
  - Cleaned `models/product.go` to strictly hold domain entity definitions (`Product`, `UpdateProductInput`), eliminating global storage variables and mutexes.
  - Implemented thread-safe CRUD operations in `database/database.go` (`GetAllProducts`, `GetProductByID`, `CreateProduct`, `UpdateProduct`, `PatchProduct`, `DeleteProduct`, and `Reset`).
  - Defined sentinel error `database.ErrNotFound` and refactored handlers (`rest/handlers/`) to use `errors.Is(err, database.ErrNotFound)` for clean 404 responses.
  - Added comprehensive test suite in `database/database_test.go` covering full CRUD lifecycle, not-found errors, and test state reset.
  - Integrated `database.Reset()` into `main_test.go` test setup for complete test isolation.
  - Verified all 19 automated tests across `main`, `config`, `database`, and `utils` pass with 100% success.
- **Session 16**:
  - Removed redundant `rest/handlers/helpers.go` file.
  - Refactored all handlers (`get_product.go`, `create_product.go`, `update_product.go`, `delete_product.go`) to import and directly call `utils.SendData(w, status, data)`.
  - Re-ran complete test suite; all 19 tests continue passing with 100% success.
