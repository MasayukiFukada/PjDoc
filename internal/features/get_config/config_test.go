package get_config

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestConfigService_Default(t *testing.T) {
	service := NewConfigService("")
	cfg := service.GetConfig()

	if cfg.PlantUMLServer != "https://kroki.io" {
		t.Errorf("expected default PlantUMLServer to be 'https://kroki.io', got '%s'", cfg.PlantUMLServer)
	}
}

func TestConfigService_Custom(t *testing.T) {
	customServer := "http://localhost:8000"
	service := NewConfigService(customServer)
	cfg := service.GetConfig()

	if cfg.PlantUMLServer != customServer {
		t.Errorf("expected PlantUMLServer to be '%s', got '%s'", customServer, cfg.PlantUMLServer)
	}
}

func TestHandler_GetSuccess(t *testing.T) {
	service := NewConfigService("https://kroki.io")
	handler := NewHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/api/config", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var resp AppConfig
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response JSON: %v", err)
	}

	if resp.PlantUMLServer != "https://kroki.io" {
		t.Errorf("expected PlantUMLServer 'https://kroki.io', got '%s'", resp.PlantUMLServer)
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	service := NewConfigService("")
	handler := NewHandler(service)

	req := httptest.NewRequest(http.MethodPost, "/api/config", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405 Method Not Allowed, got %d", rec.Code)
	}
}
