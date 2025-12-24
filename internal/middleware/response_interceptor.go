package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ResponseWriter is a wrapper around gin.ResponseWriter to capture the response body
type ResponseWriter struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

// Write captures the response body
func (w *ResponseWriter) Write(b []byte) (int, error) {
	return w.Body.Write(b)
}

// Meta holds the response metadata
type Meta struct {
	Code       int    `json:"code"`
	StatusCode string `json:"statusCode"`
}

// APIResponse is the simple standardized response structure requested
type APIResponse struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

// ResponseInterceptor middleware to standardize API responses
func ResponseInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Replace output writer with our custom wrapper
		w := &ResponseWriter{Body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = w

		// Process request
		c.Next()

		// Determine status
		status := c.Writer.Status()

		// Parse the original body
		var originalBody interface{}
		// Try to parse as JSON, otherwise use as string
		if err := json.Unmarshal(w.Body.Bytes(), &originalBody); err != nil {
			// If not JSON, use the raw string, or nil if empty
			if w.Body.Len() > 0 {
				originalBody = w.Body.String()
			} else {
				originalBody = nil
			}
		}

		// Helper to check if map is likely our APIResponse to prevent double wrapping
		// This can happen if an error handler wraps it, or we re-enter middleware
		isResponseStructure := false
		if bodyMap, ok := originalBody.(map[string]interface{}); ok {
			if _, hasMeta := bodyMap["meta"]; hasMeta {
				if _, hasData := bodyMap["data"]; hasData {
					isResponseStructure = true
				}
			}
		}

		if isResponseStructure {
			// Already wrapped, just write original bytes to the underlying writer
			w.ResponseWriter.Write(w.Body.Bytes())
			return
		}

		// Construct the new response
		response := APIResponse{
			Meta: Meta{
				Code:       status,
				StatusCode: http.StatusText(status),
			},
			Data: originalBody,
		}

		// Marshal the new response
		newBody, err := json.Marshal(response)
		if err != nil {
			// Fallback in case of marshal error (unlikely)
			w.ResponseWriter.Write(w.Body.Bytes())
			return
		}

		// Write to the actual response writer
		w.ResponseWriter.Write(newBody)
	}
}
