# Changelog

All notable changes to this project will be documented in this file.

## [2.0.0] - 2025-11-16

### 🚨 BREAKING CHANGES

- Updated minimum Go version to 1.21
- Database schema updated to support IP sets (backward compatible)
- Server setup now requires config parameter

### 🔒 Security Fixes

- **CRITICAL**: Fixed IP tracking bug that allowed users to bypass single-clap-per-IP restriction
  - Changed from storing single IP to storing IP set per URL
  - Each IP can now only clap once per URL (previously broken)
  - Added proper HTTP 429 response for duplicate clap attempts

- Removed all `panic()` calls from production code
  - Replaced with proper error handling and logging
  - Server no longer crashes on encoding/decoding errors

- Added input validation and limits:
  - Maximum 100 URLs per batch request (configurable)
  - Maximum 10 claps per update (configurable)
  - Proper URL validation on all endpoints

- Improved CORS configuration:
  - Now configurable via `ALLOWED_ORIGINS` environment variable
  - Supports comma-separated list of allowed origins
  - More restrictive default behavior

### ✨ New Features

- **Configuration Management**:
  - New `config` package for centralized configuration
  - Environment variable validation on startup
  - Sensible defaults for all settings

- **Graceful Shutdown**:
  - Proper signal handling (SIGTERM, SIGINT)
  - 30-second grace period for in-flight requests
  - Clean database closure

- **Health Monitoring**:
  - New `/health` endpoint for service monitoring
  - Kubernetes-ready health checks
  - Docker HEALTHCHECK support

- **Enhanced Logging**:
  - Structured request/response logging
  - Error tracking with context (path, IP)
  - Performance metrics (latency, status codes)

- **Better Error Handling**:
  - Consistent JSON error responses across all endpoints
  - Proper HTTP status codes
  - Custom error handler with detailed logging
  - No more panic-induced crashes

### 🏗️ Architecture Improvements

- **Database Layer**:
  - New `Database` struct for proper encapsulation
  - Added Close() method for clean shutdown
  - Removed global mutable state (backward compatible legacy support)
  - Better error messages with context

- **API Handlers**:
  - Consistent error response format across all endpoints
  - Proper status codes (400, 429, 500, etc.)
  - Better input validation
  - Improved documentation

- **Server Structure**:
  - Separated configuration from setup
  - Added panic recovery middleware
  - Better middleware organization
  - Improved static file serving with caching

### 🐛 Bug Fixes

- Fixed deprecated `ioutil` usage (replaced with `io` package)
- Fixed potential nil pointer issues in clapper list handling
- Fixed inconsistent error handling in API routes
- Improved URL normalization and validation

### 🧪 Testing

- Updated test suite for new configuration pattern
- Added comprehensive tests for:
  - Index and health endpoints
  - Error handling
  - Referer validation
  - Static file serving
- Better test structure with helper functions

### 📦 Docker & Deployment

- **Modern Dockerfile**:
  - Multi-stage build with Go 1.21-alpine
  - Minimal scratch-based runtime image
  - Better layer caching
  - Security improvements
  - Built-in health checks

- **New .dockerignore**:
  - Optimized build context
  - Smaller image sizes
  - Faster builds

### 📚 Documentation

- **Comprehensive README**:
  - Complete API documentation
  - Configuration guide
  - Security features explained
  - Architecture overview
  - Deployment examples
  - Monitoring guide

- **Improved Demo Page**:
  - Beautiful, modern UI
  - Responsive design
  - Feature showcase
  - Usage instructions
  - API endpoint documentation

### 🔧 Developer Experience

- **GitHub Actions CI/CD**:
  - Automated testing
  - Linting
  - Docker build and push
  - Code coverage reporting

- **Better Code Organization**:
  - New `config` package
  - Clearer separation of concerns
  - Improved code documentation
  - Consistent error handling patterns

### 📈 Performance

- Disabled verbose BadgerDB logging
- Improved compression configuration
- Static file caching (24 hours)
- Better connection reuse

### 🔄 Migration Guide

#### For Developers

1. Update Go to version 1.21 or higher
2. The `Setup()` function now requires a config parameter:
   ```go
   // Old
   app := Setup()

   // New
   cfg, _ := config.Load()
   app := Setup(cfg)
   ```

3. If using the database directly, consider migrating to the new Database struct:
   ```go
   // Old
   utils.DB.View(...)

   // New
   db, _ := utils.NewDatabase("/path/to/db")
   defer db.Close()
   db.GetItem(...)
   ```

#### For Users

- Existing data is automatically compatible
- Old IP tracking data will be migrated on first read
- No manual migration required
- Set `ALLOWED_ORIGINS` environment variable for better security

### Dependencies

- Updated to Go 1.21 (was 1.14)
- Kept compatible dependency versions due to network constraints
- All deprecated packages replaced with modern equivalents

---

## [1.0.0] - Previous

- Initial release
- Basic applause button functionality
- BadgerDB storage
- Docker support
