package shared_test

import (
	"encoding/json"
	"testing"

	"github.com/paul-schwendenman/magic-log-ui/internal/shared"
)

func TestMustJson(t *testing.T) {
	in := map[string]any{"hello": "world"}
	out := shared.MustJson(in)

	var decoded map[string]any
	if err := json.Unmarshal(out, &decoded); err != nil {
		t.Fatalf("Expected valid JSON, got error: %v", err)
	}
	if decoded["hello"] != "world" {
		t.Errorf("Expected 'world', got %v", decoded["hello"])
	}
}
