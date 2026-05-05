# pk-proxy Development Guide

## Project Overview
This is a learning project building an HTTP reverse proxy in Go. Two services run independently:
- **Proxy** (`cmd/proxy/main.go`) - Port 8090: Forwards requests to backend server
- **Server** (`cmd/server/main.go`) - Port 8091: Backend server responding to proxied requests

## Architecture

### Request Flow
```
Client → Proxy:8090/hello → Server:8091/hello → Response → Client
```

The proxy intercepts requests, forwards them to the backend server, and streams responses back to the client while preserving headers and status codes.

### Key Components

**`cmd/proxy/main.go`** - Proxy service that:
- Accepts incoming HTTP requests
- Creates new requests to backend server at `http://127.0.0.1:8091`
- Copies response headers (Content-Type) and body back to client using `io.Copy(w, res.Body)`
- Currently hardcoded to forward `/hello` to server's `/hello` endpoint

**`cmd/server/main.go`** - Backend server with:
- `/hello` - Returns JSON error responses using custom `writeJSONError` helper
- `/headers` - Echoes request headers as JSON response
- Consistent JSON response format with proper Content-Type headers

## Development Patterns

### JSON Response Handling
Always set Content-Type and status code before writing body:
```go
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusNotFound)
json.NewEncoder(w).Encode(response)
```

Use the `writeJSONError` helper in server for consistent error responses:
```go
writeJSONError(w, http.StatusNotFound, "error", "user not found")
```

### Proxy Response Forwarding
Copy headers from backend response before streaming body:
```go
w.Header().Set("Content-Type", res.Header.Get("Content-Type"))
io.Copy(w, res.Body)  // Efficient streaming without buffering
```

## Running the Services

Start both services in separate terminals:
```bash
# Terminal 1 - Backend server
go run cmd/server/main.go  # Listens on :8091

# Terminal 2 - Proxy
go run cmd/proxy/main.go   # Listens on :8090
```

Test the proxy: `curl http://localhost:8090/hello`

## Important Notes

- Status codes are transmitted in the response itself, not as headers (HTTP protocol handles this)
- Use `io.Copy` for efficient response streaming without loading entire body into memory
- Both services must be running for the proxy to work
- Future work: Dynamic server mapping (see TODO comment in proxy)
