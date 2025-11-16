# CLAUDE.md - AI Assistant Guide

This document helps AI assistants understand and work with the Applause Button codebase effectively.

## Project Overview

**Applause Button** is a self-hosted, privacy-focused applause/clap button service for websites and blogs. Think of it as a "like" button where readers can show varying levels of appreciation by clapping multiple times.

**Tech Stack:**
- Language: Go 1.21+
- Web Framework: Fiber v2
- Database: BadgerDB (embedded key-value store)
- Deployment: Docker, Kubernetes

## Architecture

```
applause-btn/
├── api/                    # HTTP request handlers
│   ├── get_claps.go       # Get clap count for a URL
│   ├── get_clappers.go    # Get clapper information
│   ├── get_multiple.go    # Batch query multiple URLs
│   └── update_claps.go    # Add claps to a URL
├── config/                 # Configuration management
│   └── config.go          # Env var parsing, validation, defaults
├── utils/                  # Shared utilities
│   ├── db.go              # Database operations (BadgerDB wrapper)
│   ├── structs.go         # Data structures (Item, ClapperInfo)
│   ├── url.go             # URL validation and normalization
│   └── number.go          # Number utilities (Clamp)
├── public/                 # Static files served to clients
│   ├── applause-button.js
│   ├── applause-button.css
│   └── test.html          # Demo page
├── server.go              # Main entry point, server setup
├── server_test.go         # Integration tests
└── go.mod                 # Go dependencies
```

## Critical Security Considerations

### 1. IP-Based Duplicate Prevention (MOST IMPORTANT!)

**The Core Security Mechanism:**
- Each IP address can only clap **once** per URL
- Implemented using IP sets in `utils/structs.go`
- Stored as `map[string]bool` in the `Item.SourceIPs` field

**Implementation Details:**
```go
// utils/structs.go
type Item struct {
    SourceIP  string          // DEPRECATED - kept for backward compatibility
    SourceIPs map[string]bool // The actual IP set (current implementation)
    Claps     int
    Clappers  []ClapperInfo
}

// Methods to check and add IPs
func (i *Item) HasClappedFrom(ip string) bool
func (i *Item) AddClapFrom(ip string, claps int)
```

**Why This Matters:**
- Previous version stored only ONE IP per URL
- Bug: When person B clapped, it overwrote person A's IP
- Result: Person A could clap again → broken security
- **Fixed in v2.0.0** by using IP sets

**Testing This:**
- See `api/update_claps.go:62-67` for duplicate check
- Returns HTTP 429 (Too Many Requests) on duplicate
- Critical to never break this logic!

### 2. Input Validation

**Always validate:**
- URLs must be valid (using `utils.IsURL()`)
- Clap count clamped 1-10 (using `utils.Clamp()`)
- Batch requests limited to 100 URLs (see `api/get_multiple.go:11`)
- All inputs sanitized before database storage

### 3. No Panic in Production

**Rule:** Never use `panic()` in production code
- All panics removed in v2.0.0
- Use proper error handling with context
- Return errors to caller, log with details
- See `utils/db.go` for examples

## Code Conventions

### Error Handling

**Pattern:**
```go
// Good: Consistent JSON error responses
if err != nil {
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
        "error": "Clear error message",
    })
}

// Bad: Don't use errors.New() or fmt.Errorf() directly
// (Fiber's error handler won't format it consistently)
```

### Database Operations

**Two Patterns Available:**

1. **Legacy Global DB (deprecated but supported):**
```go
item, err := utils.GetItem(sourceURL)
utils.PutItem(sourceURL, item)
```

2. **Modern Approach (preferred for new code):**
```go
db, err := utils.NewDatabase("/path/to/db")
defer db.Close()
item, err := db.GetItem(sourceURL)
```

**Important:**
- Current handlers use legacy pattern (backward compatible)
- Global `utils.DB` is initialized in `server.go:43`
- Future refactor: Pass DB as dependency injection

### Configuration

**All configuration via environment variables:**
- Parsed in `config/config.go`
- Validated on startup
- No hardcoded values
- See `config.Config` struct for all options

### Logging

**Pattern:**
```go
log.Printf("Action: %s (context: %v)", action, context)
```

- Use structured logging with context
- Include relevant info: URL, IP, error details
- Don't log sensitive data (full IPs are okay, emails be careful)

## Data Structures

### Item (utils/structs.go)

```go
type Item struct {
    SourceIP  string            // LEGACY - kept for backward compat
    SourceIPs map[string]bool   // Current: Set of IPs that clapped
    Claps     int               // Total clap count
    Clappers  []ClapperInfo     // Who clapped (if authenticated)
}
```

**Key Methods:**
- `HasClappedFrom(ip string) bool` - Check if IP already clapped
- `AddClapFrom(ip, claps)` - Record clap from IP

### ClapperInfo

```go
type ClapperInfo struct {
    Email     string
    UID       string
    CreatedAt time.Time
}
```

Populated from headers if reverse proxy provides auth info.

## API Endpoints

All endpoints documented in README.md, but key points:

### GET /get-claps
- Requires `Referer` header (source URL)
- Returns integer (clap count)
- Returns 0 if URL not found

### POST /update-claps
- Requires `Referer` header
- Body: `"count,..."` (comma-separated, first is clap count)
- **Security Check:** Returns 429 if IP already clapped
- Clamps count to 1-10

### POST /get-multiple
- Body: JSON array of URLs
- Max 100 URLs (see `api.MaxURLsPerRequest`)
- Returns array of Item objects

## Testing

### Running Tests

```bash
go test -v ./...
```

### Test Structure

```go
func getTestConfig() *config.Config {
    // Returns config for testing
}

func TestSomething(t *testing.T) {
    app := Setup(getTestConfig())
    req := httptest.NewRequest("GET", "/endpoint", nil)
    resp, err := app.Test(req, -1)
    // assertions...
}
```

**Important:**
- Tests need valid config (use `getTestConfig()`)
- Use `httptest` for HTTP testing
- DB is in-memory for tests (BadgerDB handles this)

## Common Tasks

### Adding a New Endpoint

1. Create handler in `api/` directory
2. Add route in `server.go` Setup function
3. Follow error handling pattern
4. Add tests in `server_test.go`
5. Document in README.md

### Modifying Database Schema

**Critical:** BadgerDB stores serialized structs
- Changes to `utils.Item` must be backward compatible
- Add new fields, don't remove old ones
- Use struct tags for encoding control
- Test migration from old data

### Updating Dependencies

```bash
go get -u github.com/package/name
go mod tidy
go mod verify
```

Run tests after updates!

### Building for Production

```bash
# Local build
go build -ldflags='-w -s' -o applause-btn .

# Docker build
docker build -t applause-btn .
```

## Deployment

### Environment Variables

**Required:**
- `PORT` - Server port (default: 3000)
- `DB_LOCATION` - Database path (default: /tmp/badger)

**Security:**
- `ALLOWED_ORIGINS` - CORS allowed origins (comma-separated)

**Optional:**
- `HEADER_USER_EMAIL` - Auth header name
- `HEADER_USER_ID` - Auth header name
- `MAX_URLS_PER_REQUEST` - Batch query limit
- `MAX_CLAPS_PER_UPDATE` - Max claps per request

### Docker

**Multi-stage build:**
- Builder: Go 1.21-alpine
- Runtime: scratch (minimal)
- Health checks built-in

**Volumes:**
- Mount `/data` for persistent BadgerDB storage

### Kubernetes

Use Helm chart: https://github.com/duyet/charts/tree/master/applause-btn

## Known Issues & TODOs

### Current Limitations

1. **No rate limiting yet** - Configured but not implemented
2. **No metrics endpoint** - Planned for future
3. **No distributed deployment** - BadgerDB is single-node
4. **No data export tool** - Manual extraction only

### Future Improvements

1. Add Prometheus metrics endpoint
2. Implement rate limiting middleware
3. Add admin API for management
4. Database backup/restore utilities
5. Migration to distributed DB option

## Troubleshooting

### "database not initialized" error
- Check `utils.DB` is set in main()
- Ensure `NewDatabase()` is called before handler setup

### "multiple claps from same IP" not working
- Verify IP extraction: `c.IP()` gets correct IP
- Check reverse proxy forwards real IP
- Ensure `Item.SourceIPs` map is initialized

### Tests failing with "file not found"
- BadgerDB creates temp directories
- Cleanup may fail - check `/tmp/badger-test`
- Use unique paths per test if parallel

## Performance Notes

### Database

- BadgerDB is fast for reads/writes
- GC runs periodically (see `utils.Database.RunGC()`)
- Keep database on SSD for best performance
- Size grows with number of unique URLs

### Scaling

**Current:** Single-instance deployment
**Bottleneck:** BadgerDB (embedded, not distributed)
**Recommendation:**
- Use caching layer (Redis) for reads if needed
- Consider CockroachDB/PostgreSQL for multi-region

## Version History

### v2.0.0 (Current)
- **CRITICAL FIX:** IP tracking now actually works
- Graceful shutdown
- Health endpoints
- Modernized to Go 1.21
- Comprehensive documentation

### v1.0.0
- Initial release
- Basic functionality (but IP tracking was broken!)

## Questions?

Check these files:
- `README.md` - User documentation
- `CHANGELOG.md` - Detailed version history
- `server.go` - Server setup and initialization
- `config/config.go` - Configuration options

---

**Last Updated:** 2025-11-16
**Maintained By:** Community
**License:** MIT
