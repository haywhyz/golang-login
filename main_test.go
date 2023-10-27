// main_test.go

package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheckHandler(t *testing.T) {
	// Set the router to be in release mode for testing
	gin.SetMode(gin.ReleaseMode)

	// Create a temporary router
	router := gin.New()
	router.GET("/api/health", HealthCheckHandler)

	// Create a request to the health check endpoint
	req, err := http.NewRequest("GET", "/api/health", nil)
	assert.Nil(t, err)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Serve the request to the router
	router.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Check the response body
	var response HealthCheckResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "OK", response.Status)
}
