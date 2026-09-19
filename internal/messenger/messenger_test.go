package messenger

import (
	"context"
	"errors"
	"testing"

	ms "github.com/dooray-go/dooray-sdk/openapi/messenger"
	model "github.com/dooray-go/dooray-sdk/openapi/model"
	mm "github.com/dooray-go/dooray-sdk/openapi/model/messenger"
	"github.com/mark3labs/mcp-go/mcp"

	"dooray_mcp/internal/mcptest"
)

func TestDirectMessageCallbackFailuresAndRawJSON(t *testing.T) {
	for _, mode := range []string{"success", "transport", "nil", "rejected"} {
		t.Run(mode, func(t *testing.T) {
			s := mcptest.NewServer()
			token := "token"
			callCount := 0
			register(s, &token, calls{direct: func(ctx context.Context, gotToken string, req *ms.DirectSendRequest) (*mm.DirectSendResponse, error) {
				callCount++
				if ctx != t.Context() || gotToken != token || req.OrganizationMemberId != "member-1" || req.Text != "hello" {
					t.Error("arguments not forwarded")
				}
				if mode == "transport" {
					return nil, errors.New("down")
				}
				if mode == "nil" {
					return nil, nil
				}
				return &mm.DirectSendResponse{Header: model.ResponseHeader{IsSuccessful: mode != "rejected", ResultMessage: "denied"}, RawJSON: `{"id":1}`}, nil
			}})
			res, err := s.ListTools()["dooray_messenger"].Handler(t.Context(), mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{"operation": "send", "to": "member-1", "message": "hello"}}})
			if err != nil || callCount != 1 || res.IsError != (mode != "success") {
				t.Fatalf("res=%+v err=%v calls=%d", res, err, callCount)
			}
			if mode == "success" && res.Content[0].(mcp.TextContent).Text != `{"id":1}` {
				t.Fatal("raw JSON lost")
			}
		})
	}
}

func TestChannelCreateForwardsArguments(t *testing.T) {
	s := mcptest.NewServer()
	token := "token"
	called := false
	register(s, &token, calls{create: func(ctx context.Context, gotToken, idType string, req mm.CreateChannelRequest) (*mm.IDResponse, error) {
		called = true
		if ctx != t.Context() || gotToken != "token" || idType != "email" || req.Type != "private" || req.Capacity != "8" || req.Title != "Team" || len(req.MemberIDs) != 2 {
			t.Fatalf("unexpected request: %+v", req)
		}
		return &mm.IDResponse{Header: model.ResponseHeader{IsSuccessful: true}, RawJSON: `{"result":{"id":"c1"}}`}, nil
	}})
	args := map[string]any{"operation": "create_channel", "idType": "email", "type": "private", "capacity": "8", "memberIds": []any{"a@example.com", "b@example.com"}, "title": "Team"}
	res, err := s.ListTools()["dooray_messenger_channel"].Handler(t.Context(), mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}})
	if err != nil || res.IsError || !called {
		t.Fatalf("res=%+v err=%v called=%v", res, err, called)
	}
}

func TestMessengerValidationDoesNotCallSDK(t *testing.T) {
	s := mcptest.NewServer()
	token := "token"
	callCount := 0
	register(s, &token, calls{join: func(context.Context, string, string, []string) (*mm.HeaderOnlyResponse, error) {
		callCount++
		return nil, nil
	}, leave: func(context.Context, string, string, []string) (*mm.HeaderOnlyResponse, error) {
		callCount++
		return nil, nil
	}})
	cases := []map[string]any{{"operation": "join_members", "channelId": "bad/id", "memberIds": []any{"m1"}}, {"operation": "join_members", "channelId": "c1", "memberIds": []any{}}, {"operation": "unknown", "channelId": "c1", "memberIds": []any{"m1"}}}
	for _, args := range cases {
		res, err := s.ListTools()["dooray_messenger_channel_members"].Handler(t.Context(), mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}})
		if err != nil || !res.IsError {
			t.Fatalf("expected tool error: %+v %v", res, err)
		}
	}
	if callCount != 0 {
		t.Fatalf("SDK called %d times", callCount)
	}
}

func TestMessengerToolCount(t *testing.T) {
	s := mcptest.NewServer()
	token := "test-token"
	Tools(s, &token)

	tools := s.ListTools()
	for _, name := range []string{"dooray_messenger", "dooray_messenger_channels", "dooray_messenger_channel", "dooray_messenger_channel_members", "dooray_messenger_channel_message", "dooray_messenger_channel_log", "dooray_messenger_channel_thread"} {
		if _, ok := tools[name]; !ok {
			t.Errorf("missing tool %s", name)
		}
	}
}

func TestChannelOperationRouting(t *testing.T) {
	s := mcptest.NewServer()
	token := "token"
	got := ""
	header := func() *mm.HeaderOnlyResponse {
		return &mm.HeaderOnlyResponse{Header: model.ResponseHeader{IsSuccessful: true}, RawJSON: `{"ok":true}`}
	}
	message := func() *mm.SendMessageResponse {
		return &mm.SendMessageResponse{Header: model.ResponseHeader{IsSuccessful: true}, RawJSON: `{"ok":true}`}
	}
	register(s, &token, calls{
		list: func(ctx context.Context, tok string) (*mm.ListChannelsResponse, error) {
			if ctx == nil || tok != "token" {
				t.Error("list context/token")
			}
			got = "list"
			return &mm.ListChannelsResponse{Header: model.ResponseHeader{IsSuccessful: true}, RawJSON: `{"ok":true}`}, nil
		},
		join: func(ctx context.Context, tok, ch string, ids []string) (*mm.HeaderOnlyResponse, error) {
			if ctx == nil || tok != "token" || ch != "c1" || len(ids) != 1 || ids[0] != "m1" {
				t.Error("join args")
			}
			got = "join"
			return header(), nil
		},
		leave: func(ctx context.Context, tok, ch string, ids []string) (*mm.HeaderOnlyResponse, error) {
			if ch != "c1" || ids[0] != "m1" {
				t.Error("leave args")
			}
			got = "leave"
			return header(), nil
		},
		send: func(ctx context.Context, tok, ch string, r *ms.SendMessageRequest) (*mm.SendMessageResponse, error) {
			if ch != "c1" || r.Text != "hello" {
				t.Error("send args")
			}
			got = "send"
			return message(), nil
		},
		update: func(ctx context.Context, tok, ch, log, text string) (*mm.HeaderOnlyResponse, error) {
			if ch != "c1" || log != "l1" || text != "hello" {
				t.Error("update args")
			}
			got = "update"
			return header(), nil
		},
		delete: func(ctx context.Context, tok, ch, log string) (*mm.HeaderOnlyResponse, error) {
			if ch != "c1" || log != "l1" {
				t.Error("delete args")
			}
			got = "delete"
			return header(), nil
		},
		reply: func(ctx context.Context, tok, ch, log, text string) (*mm.SendMessageResponse, error) {
			if ch != "c1" || log != "l1" || text != "hello" {
				t.Error("reply args")
			}
			got = "reply"
			return message(), nil
		},
		thread: func(ctx context.Context, tok, ch string, r mm.CreateAndSendThreadRequest) (*mm.SendMessageResponse, error) {
			if ch != "c1" || r.Text != "hello" || r.ThreadText != "first reply" {
				t.Error("thread args")
			}
			got = "thread"
			return message(), nil
		},
		threadFromLog: func(ctx context.Context, tok, ch, log, text string) (*mm.SendMessageResponse, error) {
			if ch != "c1" || log != "l1" || text != "hello" {
				t.Error("thread log args")
			}
			got = "thread-log"
			return message(), nil
		},
	})
	tests := []struct {
		name, want string
		args       map[string]any
	}{
		{"dooray_messenger_channels", "list", map[string]any{"operation": "find_channels"}},
		{"dooray_messenger_channel_members", "join", map[string]any{"operation": "join_members", "channelId": "c1", "memberIds": []any{"m1"}}},
		{"dooray_messenger_channel_members", "leave", map[string]any{"operation": "leave_members", "channelId": "c1", "memberIds": []any{"m1"}}},
		{"dooray_messenger_channel_message", "send", map[string]any{"operation": "send_message", "channelId": "c1", "text": "hello"}},
		{"dooray_messenger_channel_log", "update", map[string]any{"operation": "update_log", "channelId": "c1", "logId": "l1", "text": "hello"}},
		{"dooray_messenger_channel_log", "delete", map[string]any{"operation": "delete_log", "channelId": "c1", "logId": "l1"}},
		{"dooray_messenger_channel_log", "reply", map[string]any{"operation": "reply_log", "channelId": "c1", "logId": "l1", "text": "hello"}},
		{"dooray_messenger_channel_thread", "thread", map[string]any{"operation": "create_thread", "channelId": "c1", "text": "hello", "threadText": "first reply"}},
		{"dooray_messenger_channel_thread", "thread-log", map[string]any{"operation": "create_thread_from_log", "channelId": "c1", "logId": "l1", "text": "hello"}},
	}
	for _, tc := range tests {
		t.Run(tc.want, func(t *testing.T) {
			got = ""
			res, err := s.ListTools()[tc.name].Handler(t.Context(), mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: tc.args}})
			if err != nil || res.IsError || got != tc.want || res.Content[0].(mcp.TextContent).Text != `{"ok":true}` {
				t.Fatalf("got=%s res=%+v err=%v", got, res, err)
			}
		})
	}
}

func TestConditionalRequiredArguments(t *testing.T) {
	s := mcptest.NewServer()
	token := "token"
	register(s, &token, calls{})
	for _, tc := range []struct {
		name string
		args map[string]any
	}{
		{"dooray_messenger_channel_log", map[string]any{"operation": "update_log", "channelId": "c1", "logId": "l1"}},
		{"dooray_messenger_channel_log", map[string]any{"operation": "reply_log", "channelId": "c1", "logId": "l1"}},
		{"dooray_messenger_channel_thread", map[string]any{"operation": "create_thread_from_log", "channelId": "c1", "text": "hello"}},
	} {
		res, err := s.ListTools()[tc.name].Handler(t.Context(), mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: tc.args}})
		if err != nil || !res.IsError {
			t.Fatalf("expected validation error: %+v %v", res, err)
		}
	}
}

func TestDirectMessageMissingArgumentsReturnErrors(t *testing.T) {
	s := mcptest.NewServer()
	token := "token"
	register(s, &token, calls{})
	for _, args := range []map[string]any{{"operation": "send", "message": "hello"}, {"operation": "send", "to": "m1"}, {"operation": "send"}} {
		res, err := s.ListTools()["dooray_messenger"].Handler(t.Context(), mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}})
		if err != nil || res == nil || !res.IsError {
			t.Fatalf("expected non-panicking tool error: %+v %v", res, err)
		}
	}
}
