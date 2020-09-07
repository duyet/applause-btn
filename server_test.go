package main

import (
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	app := Setup()

	expectedBody := "Hello, World!"

	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req, -1)

	assert.Equalf(t, false, err != nil, "Index route")
	assert.Equalf(t, 200, resp.StatusCode, "statusCode should be 200")

	body, err := ioutil.ReadAll(resp.Body)
	assert.Nilf(t, err, "err should be nil")
	assert.Equalf(t, string(body), expectedBody, "body should ok")
}

func TestGetClapsNoReferer(t *testing.T) {
	app := Setup()
	req := httptest.NewRequest("GET", "/get-claps", nil)

	resp, err := app.Test(req, -1)
	assert.Equalf(t, false, err != nil, "err should be nil")
	assert.NotEqual(t, 200, resp.StatusCode, "statusCode should not be 200")

	body, err := ioutil.ReadAll(resp.Body)
	assert.Nilf(t, err, "err should be nil")
	assert.Equalf(t, string(body), "no referer set", "body should be 'no referer set'")
}

func TestGetClapsWrongReferer(t *testing.T) {
	app := Setup()
	req := httptest.NewRequest("GET", "/get-claps", nil)

	referer := "ahihi"
	req.Header.Add("Referer", referer)

	// test StatusCode
	resp, err := app.Test(req, -1)
	assert.Equalf(t, false, err != nil, "err should be nil")
	assert.NotEqual(t, 200, resp.StatusCode, "statusCode should not be 200")

	// test body
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nilf(t, err, "err should be nil")
	expectedBody := fmt.Sprintf("Referer is not a URL [%s]", referer)
	assert.Equalf(t, string(body), expectedBody, fmt.Sprintf("body should be: %s", expectedBody))
}

func TestGetClapsSuccessReferer(t *testing.T) {
	app := Setup()
	req := httptest.NewRequest("GET", "/get-claps", nil)

	referer := "https://duyet.net"
	req.Header.Add("Referer", referer)

	// test StatusCode
	resp, err := app.Test(req, -1)
	assert.Equalf(t, false, err != nil, "err should be nil")
	assert.Equalf(t, 200, resp.StatusCode, "statusCode should be 200")

	// TODO: mock DB and test for body
}
