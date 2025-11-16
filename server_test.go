package main

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/duyet/applause-btn/config"
	"github.com/stretchr/testify/assert"
)

func getTestConfig() *config.Config {
	return &config.Config{
		Port:              "3000",
		DBLocation:        "/tmp/badger-test",
		HeaderUserEmail:   "x-authenticated-user-email",
		HeaderUserID:      "x-authenticated-uid",
		AllowedOrigins:    []string{"*"},
		RateLimitEnabled:  false, // Disable for tests
		MaxURLsPerRequest: 100,
		MaxClapsPerUpdate: 10,
	}
}

func TestIndexRoute(t *testing.T) {
	app := Setup(getTestConfig())

	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req, -1)

	assert.NoError(t, err, "Index route should not error")
	assert.Equal(t, 200, resp.StatusCode, "statusCode should be 200")

	var body map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	assert.NoError(t, err, "Should decode JSON response")
	assert.Equal(t, "Applause Button", body["service"])
	assert.Equal(t, "running", body["status"])
}

func TestHealthRoute(t *testing.T) {
	app := Setup(getTestConfig())

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var body map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	assert.NoError(t, err)
	assert.Equal(t, "healthy", body["status"])
}

func TestGetClapsNoReferer(t *testing.T) {
	app := Setup(getTestConfig())
	req := httptest.NewRequest("GET", "/get-claps", nil)

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode, "Should return 400 for missing referer")

	var body map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	assert.NoError(t, err)
	assert.Contains(t, body["error"], "Referer")
}

func TestGetClapsWrongReferer(t *testing.T) {
	app := Setup(getTestConfig())
	req := httptest.NewRequest("GET", "/get-claps", nil)

	referer := "not-a-url"
	req.Header.Add("Referer", referer)

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode, "Should return 400 for invalid URL")

	var body map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	assert.NoError(t, err)
	assert.Contains(t, body["error"], "URL")
}

func TestGetClapsValidReferer(t *testing.T) {
	app := Setup(getTestConfig())
	req := httptest.NewRequest("GET", "/get-claps", nil)

	referer := "https://example.com/page"
	req.Header.Add("Referer", referer)

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestStaticFiles(t *testing.T) {
	app := Setup(getTestConfig())

	req := httptest.NewRequest("GET", "/public/test.html", nil)
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	// Should either find the file (200) or not (404), but not error
	assert.True(t, resp.StatusCode == 200 || resp.StatusCode == 404)
}
