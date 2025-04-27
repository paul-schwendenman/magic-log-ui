package jqfilter_test

import (
	"testing"

	"github.com/paul-schwendenman/magic-log-ui/internal/jqfilter"
)

func TestSimpleJQApply(t *testing.T) {
	jqfilter.Init("{id: .trace_id, text: .message}")

	input := map[string]string{
		"trace_id": "abc123",
		"message":  "hello world",
	}

	out := jqfilter.Apply(input)

	if out["id"] != "abc123" {
		t.Errorf("Expected id abc123, got %s", out["id"])
	}
	if out["text"] != "hello world" {
		t.Errorf("Expected text hello world, got %s", out["text"])
	}
}
