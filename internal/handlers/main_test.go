package handlers

import (
	"bytes"
	"cimrique-redis/internal/models"
	"cimrique-redis/internal/service"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
)

type MainTest struct {
	app     *fiber.App
	client  mockClientInterface
	service service.ServicesFuncs
	handler Handler
}

type mockClient struct{}

type mockClientInterface interface {
	UnpackRequest(body []byte) (models.Request, int)
	EnqueueHigh(req models.Request) error
	EnqueueMed(req models.Request) error
	EnqueueLow(req models.Request) error
}

var MTest MainTest

func TestMain(m *testing.M) {
	MTest.setupApp()
	code := m.Run()
	os.Exit(code)
}

func (m *MainTest) setupApp() {
	m.app = fiber.New()
	m.client = new(mockClient) //create mock client here
	m.service = service.NewServicesFuncs(m.client)
	m.handler = NewHandler(&m.service)
	m.app.Post("/enqueue", m.handler.HandleEnqueue)
}

func (m *MainTest) sendRequest(app *fiber.App, t *testing.T, body []byte, apiKey string) (int, string) {
	req := httptest.NewRequest(http.MethodPost, "/enqueue", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", apiKey)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Error sending request: %v", err)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	return resp.StatusCode, buf.String()
}

// mock functions
func (m mockClient) UnpackRequest(body []byte) (models.Request, int) {
	//que logic here
	var msg models.QueMessage
	json.Unmarshal(body, &msg)
	req := msg.Message
	score := msg.Score
	return req, score
}
func (m mockClient) EnqueueHigh(req models.Request) error {
	return nil
}
func (m mockClient) EnqueueMed(req models.Request) error {
	return nil
}
func (m mockClient) EnqueueLow(req models.Request) error {
	return nil
}
