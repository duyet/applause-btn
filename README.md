# Applause Button 👏

A high-performance, self-hosted applause button service for your blog or website. Track and display appreciation from your readers with a simple, elegant button.

<center><img src=".github/demo.png" width="100" /></center>

## ✨ Features

- **Fast & Lightweight**: Built with Go and Fiber framework
- **Self-Hosted**: Full control over your data
- **Secure**: IP-based duplicate prevention, rate limiting, CORS configuration
- **Scalable**: Embedded BadgerDB for efficient storage
- **Production-Ready**: Graceful shutdown, health checks, structured logging
- **Docker Support**: Easy deployment with Docker and Kubernetes
- **Modern**: Go 1.21+, updated dependencies, best practices

## 🚀 Quick Start

### Using Docker (Recommended)

```bash
docker run -it -p 3000:3000 -v $(pwd)/data:/data duyetdev/applause-btn
```

### Using Docker Compose

```yaml
version: '3.8'
services:
  applause:
    image: duyetdev/applause-btn
    ports:
      - "3000:3000"
    volumes:
      - ./data:/data
    environment:
      - PORT=3000
      - DB_LOCATION=/data/badger
      - ALLOWED_ORIGINS=https://yourblog.com,https://www.yourblog.com
```

### Using Kubernetes Helm

```bash
helm repo add duyet https://duyet.github.io/charts
helm install applause duyet/applause-btn
```

More details at https://github.com/duyet/charts/tree/master/applause-btn

### Build from Source

```bash
# Clone the repository
git clone https://github.com/duyet/applause-btn.git
cd applause-btn

# Build
go build -o applause .

# Run
./applause
```

## 📖 Usage

Once deployed, integrate the applause button into your website:

```html
<!DOCTYPE html>
<html>
<head>
    <link rel="stylesheet" href="https://applause.yourdomain.com/public/applause-button.css" />
    <script src="https://applause.yourdomain.com/public/applause-button.js"></script>
</head>
<body>
    <h1>My Blog Post</h1>
    <p>Great content here...</p>

    <!-- Add the applause button -->
    <applause-button
        style="width: 58px; height: 58px"
        multiclap="true"
        api="https://applause.yourdomain.com"
    />
</body>
</html>
```

## 🔧 Configuration

Configure the service using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `3000` |
| `DB_LOCATION` | BadgerDB storage path | `/tmp/badger` |
| `ALLOWED_ORIGINS` | CORS allowed origins (comma-separated) | `*` (all) |
| `HEADER_USER_EMAIL` | Header name for authenticated user email | `x-authenticated-user-email` |
| `HEADER_USER_ID` | Header name for authenticated user ID | `x-authenticated-uid` |
| `MAX_URLS_PER_REQUEST` | Max URLs in batch query | `100` |
| `MAX_CLAPS_PER_UPDATE` | Max claps per update | `10` |

### Example Configuration

```bash
export PORT=8080
export DB_LOCATION=/var/lib/applause/data
export ALLOWED_ORIGINS=https://blog.example.com,https://www.example.com
./applause
```

## 📡 API Endpoints

### `GET /health`

Health check endpoint for monitoring.

**Response:**
```json
{
  "status": "healthy",
  "time": "2024-01-15T10:30:00Z"
}
```

### `GET /`

Service information and available endpoints.

**Response:**
```json
{
  "service": "Applause Button",
  "version": "2.0.0",
  "status": "running",
  "endpoints": {
    "health": "/health",
    "get_claps": "/get-claps",
    "get_clappers": "/get-clappers",
    "get_multiple": "/get-multiple (POST)",
    "update_claps": "/update-claps (POST)"
  }
}
```

### `GET /get-claps`

Get the number of claps for a URL (from Referer header).

**Headers:**
- `Referer`: The URL to get claps for (required)

**Response:**
```json
42
```

### `GET /get-clappers`

Get information about who clapped (if authenticated headers are provided).

**Headers:**
- `Referer`: The URL to get clappers for (required)

**Response:**
```json
[
  {
    "Email": "user@example.com",
    "UID": "user123",
    "CreatedAt": "2024-01-15T10:30:00Z"
  }
]
```

### `POST /get-multiple`

Get claps for multiple URLs in a single request.

**Request Body:**
```json
[
  "https://blog.example.com/post1",
  "https://blog.example.com/post2",
  "https://blog.example.com/post3"
]
```

**Response:**
```json
[
  { "Claps": 42, "Clappers": [...] },
  { "Claps": 15, "Clappers": [...] },
  { "Claps": 7, "Clappers": [...] }
]
```

**Limits:**
- Maximum 100 URLs per request (configurable via `MAX_URLS_PER_REQUEST`)

### `POST /update-claps`

Add claps to a URL.

**Headers:**
- `Referer`: The URL to add claps to (required)

**Request Body:**
```
1,10,<other-data>
```
First number is the clap count (1-10). Default is 1 if not provided.

**Response:**
```json
43
```

**Security:**
- Each IP address can only clap once per URL
- Duplicate claps return HTTP 429 (Too Many Requests)
- Clap count is clamped between 1-10

## 🔒 Security Features

### IP-Based Duplicate Prevention

Each IP address can only clap once per URL. The system maintains a set of IP addresses per URL to prevent abuse.

### CORS Configuration

Configure allowed origins to prevent unauthorized domains from using your service:

```bash
export ALLOWED_ORIGINS=https://yourblog.com,https://www.yourblog.com
```

### Input Validation

- URL validation for all endpoints
- Request size limits
- Array size limits for batch operations

### Error Handling

- Consistent JSON error responses
- Proper HTTP status codes
- No panic crashes - all errors are handled gracefully

## 🧪 Development

### Running Tests

```bash
go test -v ./...
```

### Running Locally

```bash
go run server.go
```

The server will start on `http://localhost:3000`

### Building

```bash
# Local build
go build -o applause .

# Docker build
docker build -t applause-btn .
```

## 🏗️ Architecture

```
applause-btn/
├── api/                    # API handlers
│   ├── get_claps.go       # Get clap count
│   ├── get_clappers.go    # Get clapper info
│   ├── get_multiple.go    # Batch query
│   └── update_claps.go    # Add claps
├── config/                 # Configuration management
│   └── config.go
├── utils/                  # Utilities
│   ├── db.go              # Database operations
│   ├── structs.go         # Data structures
│   ├── url.go             # URL utilities
│   └── number.go          # Number utilities
├── public/                 # Static files (JS, CSS)
├── server.go              # Main server
├── server_test.go         # Tests
├── Dockerfile             # Docker configuration
└── go.mod                 # Go dependencies
```

## 📊 Monitoring

### Health Checks

The `/health` endpoint provides service health status:

```bash
curl http://localhost:3000/health
```

### Logs

The service provides structured logging with:
- Request/response logging
- Error tracking with context
- Performance metrics (latency, status codes)

### Docker Health Checks

Built-in Docker health checks monitor service availability.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Development Guidelines

1. Write tests for new features
2. Update documentation
3. Follow Go best practices
4. Run `go fmt` before committing
5. Ensure all tests pass

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Credits

- Original applause-button.com concept
- Built with [Fiber](https://gofiber.io/) web framework
- Storage powered by [BadgerDB](https://github.com/dgraph-io/badger)

## 📮 Support

- Issues: https://github.com/duyet/applause-btn/issues
- Docker Hub: https://hub.docker.com/r/duyetdev/applause-btn
- GitHub: https://github.com/duyet/applause-btn

---

Made with ❤️ by the community
