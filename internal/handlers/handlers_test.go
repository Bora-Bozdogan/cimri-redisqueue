package handlers

import (
	"bytes"
	"cimrique-redis/internal/models"
	"encoding/json"
	"net/http"
	"testing"
)

func strPtr(s string) *string { return &s }
func intPtr(i int) *int       { return &i }

func TestHandleUpdate(t *testing.T) {
	MTest.TestHandleUpdate(t)
}

func (m *MainTest) TestHandleUpdate(t *testing.T) {
	app := m.app

	tests := []struct {
		name     string
		message  models.QueMessage
		rawBody  []byte
		apiKey   string
		expected int
		contains string
	}{
		{
			name: "invalid score",
			message: models.QueMessage{
				Message: models.Request{
					ApiKey:             strPtr("amazon-key"),
					ProductName:        strPtr("IPhone 16"),
					ProductDescription: strPtr("Latest Apple flagship phone"),
					ProductImage:       strPtr("https://example.com/iphone16.jpg"),
					StoreName:          strPtr("Amazon"),
					Price:              intPtr(1475),
					Stock:              intPtr(50),
					PopularityScore:    intPtr(5),
					UrgencyScore:       intPtr(5),
					},
				Score: 0,
			},
			apiKey:   "amazon-key",
			expected: http.StatusBadRequest,
			contains: "Invalid score",
		},
		{
			name: "low que, 300 score",
			message: models.QueMessage{
				Message: models.Request{
					ApiKey:             strPtr("amazon-key"),
					ProductName:        strPtr("IPhone 16"),
					ProductDescription: strPtr("Latest Apple flagship phone"),
					ProductImage:       strPtr("https://example.com/iphone16.jpg"),
					StoreName:          strPtr("Amazon"),
					Price:              intPtr(1475),
					Stock:              intPtr(50),
					PopularityScore:    intPtr(5),
					UrgencyScore:       intPtr(5),
					},
				Score: 300,
			},
			apiKey:   "amazon-key",
			expected: http.StatusOK,
			contains: "Enqueued request score 300 on queue low",
		},
		{
			name: "med que, 650 score",
			message: models.QueMessage{
				Message: models.Request{
					ApiKey:             strPtr("amazon-key"),
					ProductName:        strPtr("IPhone 16"),
					ProductDescription: strPtr("Latest Apple flagship phone"),
					ProductImage:       strPtr("https://example.com/iphone16.jpg"),
					StoreName:          strPtr("Amazon"),
					Price:              intPtr(1475),
					Stock:              intPtr(50),
					PopularityScore:    intPtr(5),
					UrgencyScore:       intPtr(5),
					},
				Score: 650,
			},
			apiKey:   "amazon-key",
			expected: http.StatusOK,
			contains: "Enqueued request score 650 on queue med",
		},
		{
			name: "high que, 900 score",
			message: models.QueMessage{
				Message: models.Request{
					ApiKey:             strPtr("amazon-key"),
					ProductName:        strPtr("IPhone 16"),
					ProductDescription: strPtr("Latest Apple flagship phone"),
					ProductImage:       strPtr("https://example.com/iphone16.jpg"),
					StoreName:          strPtr("Amazon"),
					Price:              intPtr(1475),
					Stock:              intPtr(50),
					PopularityScore:    intPtr(5),
					UrgencyScore:       intPtr(5),
					},
				Score: 900,
			},
			apiKey:   "amazon-key",
			expected: http.StatusOK,
			contains: "Enqueued request score 900 on queue high",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var body []byte
			var err error

			if tc.rawBody != nil {
				body = tc.rawBody
			} else {
				body, err = json.Marshal(tc.message)
				if err != nil {
					t.Fatalf("Failed to marshal JSON: %v", err)
				}
			}

			status, response := m.sendRequest(app, t, body, tc.apiKey)

			if status != tc.expected {
				t.Errorf("[%s] Expected status %d, got %d", tc.name, tc.expected, status)
			}
			if !bytes.Contains([]byte(response), []byte(tc.contains)) {
				t.Errorf("[%s] Expected response to contain %q, got %q", tc.name, tc.contains, response)
			}
		})
	}
}
