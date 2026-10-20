package messenger

import (
	"context"
	"testing"

	model "github.com/dooray-go/dooray-sdk/openapi/model"
	"github.com/mark3labs/mcp-go/mcp"

	"dooray_mcp/internal/mcptest"
)

func TestChannelLogsForwardsArgumentsAndSize(t *testing.T) {
	cases := []struct {
		name    string
		size    any
		want    int
		isError bool
	}{
		{"default", nil, defaultLogsSize, false},
		{"number", float64(200), 200, false},
		{"string", "3", 3, false},
		{"zero", float64(0), 0, true},
		{"over", float64(1001), 0, true},
		{"fraction", 1.5, 0, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := mcptest.NewServer()
			token := "token"
			got := 0
			register(s, &token, calls{logs: func(ctx context.Context, gotToken, ch string, size int) (*logsResponse, error) {
				if gotToken != "token" || ch != "c1" {
					t.Fatalf("token=%s ch=%s", gotToken, ch)
				}
				got = size
				return &logsResponse{Header: model.ResponseHeader{IsSuccessful: true}, RawJSON: `{"result":[]}`}, nil
			}})
			args := map[string]any{"operation": "find_logs", "channelId": "c1"}
			if tc.size != nil {
				args["size"] = tc.size
			}
			res, err := s.ListTools()["dooray_messenger_channel_logs"].Handler(t.Context(), mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}})
			if err != nil || res.IsError != tc.isError || got != tc.want {
				t.Fatalf("res=%+v err=%v size=%d", res, err, got)
			}
		})
	}
}
