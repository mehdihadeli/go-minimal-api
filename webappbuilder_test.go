package minimalapi

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestWebApplication(t *testing.T) {
	// Create a new WebApplication
	app, err := NewWebApplicationBuilder().
		WithPort(8080).
		Build()

	// Start the app
	go app.Run()

	// Send a request to the app
	res, err := http.Get("http://localhost:8080/")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}
