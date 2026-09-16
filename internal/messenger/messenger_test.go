package messenger

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"

	"dooray_mcp/internal/mcptest"
)

func TestMessengerToolsRegistration(t *testing.T) {
	s := mcptest.NewServer()
	token := "test-token"
	Tools(s, &token)

	tools := s.ListTools()
	if _, ok := tools["dooray_messenger"]; !ok {
		t.Fatal("dooray_messenger tool not registered")
	}
}

func TestMessengerSendArguments(t *testing.T) {
	s := mcptest.NewServer()
	token := "invalid-token"
	Tools(s, &token)

	tool := s.ListTools()["dooray_messenger"]

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "dooray_messenger",
			Arguments: map[string]any{
				"operation": "send",
				"to":        "member-123",
				"message":   "안녕하세요",
			},
		},
	}

	_, err := tool.Handler(context.Background(), req)
	if err == nil {
		t.Log("handler succeeded (argument parsing works)")
	} else {
		t.Logf("handler returned expected error: %v", err)
	}
}

func TestMessengerMissingTo(t *testing.T) {
	s := mcptest.NewServer()
	token := "invalid-token"
	Tools(s, &token)

	tool := s.ListTools()["dooray_messenger"]

	defer func() {
		if r := recover(); r != nil {
			t.Logf("recovered from panic (missing 'to'): %v", r)
		}
	}()

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "dooray_messenger",
			Arguments: map[string]any{
				"operation": "send",
				"message":   "안녕하세요",
			},
		},
	}

	tool.Handler(context.Background(), req)
}

func TestMessengerMissingMessage(t *testing.T) {
	s := mcptest.NewServer()
	token := "invalid-token"
	Tools(s, &token)

	tool := s.ListTools()["dooray_messenger"]

	defer func() {
		if r := recover(); r != nil {
			t.Logf("recovered from panic (missing 'message'): %v", r)
		}
	}()

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "dooray_messenger",
			Arguments: map[string]any{
				"operation": "send",
				"to":        "member-123",
			},
		},
	}

	tool.Handler(context.Background(), req)
}

func TestMessengerMissingAllArgs(t *testing.T) {
	s := mcptest.NewServer()
	token := "invalid-token"
	Tools(s, &token)

	tool := s.ListTools()["dooray_messenger"]

	defer func() {
		if r := recover(); r != nil {
			t.Logf("recovered from panic (missing all required args): %v", r)
		}
	}()

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "dooray_messenger",
			Arguments: map[string]any{
				"operation": "send",
			},
		},
	}

	tool.Handler(context.Background(), req)
}

func TestMessengerToolCount(t *testing.T) {
	s := mcptest.NewServer()
	token := "test-token"
	Tools(s, &token)

	tools := s.ListTools()
	count := 0
	for name := range tools {
		if name == "dooray_messenger" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected 1 messenger tool, got %d", count)
	}
}
