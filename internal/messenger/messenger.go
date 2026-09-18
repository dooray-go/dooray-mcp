package messenger

import (
	"context"
	"fmt"
	"strings"

	ms "github.com/dooray-go/dooray-sdk/openapi/messenger"
	mm "github.com/dooray-go/dooray-sdk/openapi/model/messenger"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type calls struct {
	direct        func(context.Context, string, *ms.DirectSendRequest) (*mm.DirectSendResponse, error)
	list          func(context.Context, string) (*mm.ListChannelsResponse, error)
	create        func(context.Context, string, string, mm.CreateChannelRequest) (*mm.IDResponse, error)
	join, leave   func(context.Context, string, string, []string) (*mm.HeaderOnlyResponse, error)
	send          func(context.Context, string, string, *ms.SendMessageRequest) (*mm.SendMessageResponse, error)
	update        func(context.Context, string, string, string, string) (*mm.HeaderOnlyResponse, error)
	delete        func(context.Context, string, string, string) (*mm.HeaderOnlyResponse, error)
	reply         func(context.Context, string, string, string, string) (*mm.SendMessageResponse, error)
	thread        func(context.Context, string, string, mm.CreateAndSendThreadRequest) (*mm.SendMessageResponse, error)
	threadFromLog func(context.Context, string, string, string, string) (*mm.SendMessageResponse, error)
}

func Tools(s *server.MCPServer, token *string) {
	c := ms.NewDefaultMessenger()
	register(s, token, calls{c.DirectSendContext, c.GetChannelsContext, c.CreateChannelContext, c.JoinChannelMembersContext, c.LeaveChannelMembersContext, c.SendMessageContext, c.UpdateChannelLogContext, c.DeleteChannelLogContext, c.ReplyChannelLogContext, c.CreateAndSendThreadContext, c.CreateAndSendThreadFromLogContext})
}

func register(s *server.MCPServer, token *string, c calls) {
	add(s, "dooray_messenger", "Send a direct message", []mcp.ToolOption{str("operation", true, "send"), str("to", true), str("message", true)}, func(ctx context.Context, a map[string]any) (*mcp.CallToolResult, error) {
		if e := op(a, "send"); e != nil {
			return bad(e)
		}
		to, e := id(a, "to")
		if e != nil {
			return bad(e)
		}
		text, e := required(a, "message")
		if e != nil {
			return bad(e)
		}
		if e = tok(token); e != nil {
			return bad(e)
		}
		r, e := c.direct(ctx, *token, &ms.DirectSendRequest{Text: text, OrganizationMemberId: to})
		if r == nil {
			return result(e, false, "", "", false)
		}
		return result(e, true, r.RawJSON, r.Header.ResultMessage, r.Header.IsSuccessful)
	})
	add(s, "dooray_messenger_channels", "List channels", []mcp.ToolOption{str("operation", true, "find_channels")}, func(ctx context.Context, a map[string]any) (*mcp.CallToolResult, error) {
		if e := op(a, "find_channels"); e != nil {
			return bad(e)
		}
		if e := tok(token); e != nil {
			return bad(e)
		}
		r, e := c.list(ctx, *token)
		if r == nil {
			return result(e, false, "", "", false)
		}
		return result(e, true, r.RawJSON, r.Header.ResultMessage, r.Header.IsSuccessful)
	})
	add(s, "dooray_messenger_channel", "Create a channel", []mcp.ToolOption{str("operation", true, "create_channel"), str("idType", true, "email", "member-id"), str("type", true, "direct", "private"), str("capacity", false), arr("memberIds"), str("title", false)}, func(ctx context.Context, a map[string]any) (*mcp.CallToolResult, error) {
		if e := op(a, "create_channel"); e != nil {
			return bad(e)
		}
		it, e := enm(a, "idType", "email", "member-id")
		if e != nil {
			return bad(e)
		}
		ct, e := enm(a, "type", "direct", "private")
		if e != nil {
			return bad(e)
		}
		members, e := array(a, "memberIds")
		if e != nil {
			return bad(e)
		}
		cap, e := optional(a, "capacity")
		if e != nil {
			return bad(e)
		}
		title, e := optional(a, "title")
		if e != nil {
			return bad(e)
		}
		if e = tok(token); e != nil {
			return bad(e)
		}
		r, e := c.create(ctx, *token, it, mm.CreateChannelRequest{Type: ct, Capacity: cap, MemberIDs: members, Title: title})
		if r == nil {
			return result(e, false, "", "", false)
		}
		return result(e, true, r.RawJSON, r.Header.ResultMessage, r.Header.IsSuccessful)
	})
	add(s, "dooray_messenger_channel_members", "Join or leave channel members", []mcp.ToolOption{str("operation", true, "join_members", "leave_members"), str("channelId", true), arr("memberIds")}, func(ctx context.Context, a map[string]any) (*mcp.CallToolResult, error) {
		o, e := enm(a, "operation", "join_members", "leave_members")
		if e != nil {
			return bad(e)
		}
		ch, e := id(a, "channelId")
		if e != nil {
			return bad(e)
		}
		members, e := array(a, "memberIds")
		if e != nil {
			return bad(e)
		}
		if e = tok(token); e != nil {
			return bad(e)
		}
		var r *mm.HeaderOnlyResponse
		if o == "join_members" {
			r, e = c.join(ctx, *token, ch, members)
		} else {
			r, e = c.leave(ctx, *token, ch, members)
		}
		if r == nil {
			return result(e, false, "", "", false)
		}
		return result(e, true, r.RawJSON, r.Header.ResultMessage, r.Header.IsSuccessful)
	})
	add(s, "dooray_messenger_channel_message", "Send a channel message", []mcp.ToolOption{str("operation", true, "send_message"), str("channelId", true), str("text", true)}, func(ctx context.Context, a map[string]any) (*mcp.CallToolResult, error) {
		if e := op(a, "send_message"); e != nil {
			return bad(e)
		}
		ch, e := id(a, "channelId")
		if e != nil {
			return bad(e)
		}
		text, e := required(a, "text")
		if e != nil {
			return bad(e)
		}
		if e = tok(token); e != nil {
			return bad(e)
		}
		r, e := c.send(ctx, *token, ch, &ms.SendMessageRequest{Text: text})
		if r == nil {
			return result(e, false, "", "", false)
		}
		return result(e, true, r.RawJSON, r.Header.ResultMessage, r.Header.IsSuccessful)
	})
	addLog(s, token, c)
	addThread(s, token, c)
}

func addLog(s *server.MCPServer, t *string, c calls) {
	add(s, "dooray_messenger_channel_log", "Update, delete, or reply to a channel log", []mcp.ToolOption{str("operation", true, "update_log", "delete_log", "reply_log"), str("channelId", true), str("logId", true), str("text", false)}, func(ctx context.Context, a map[string]any) (*mcp.CallToolResult, error) {
		o, e := enm(a, "operation", "update_log", "delete_log", "reply_log")
		if e != nil {
			return bad(e)
		}
		ch, e := id(a, "channelId")
		if e != nil {
			return bad(e)
		}
		log, e := id(a, "logId")
		if e != nil {
			return bad(e)
		}
		if e = tok(t); e != nil {
			return bad(e)
		}
		if o == "delete_log" {
			r, x := c.delete(ctx, *t, ch, log)
			if r == nil {
				return result(x, false, "", "", false)
			}
			return result(x, true, r.RawJSON, r.Header.ResultMessage, r.Header.IsSuccessful)
		}
		text, e := required(a, "text")
		if e != nil {
			return bad(e)
		}
		if o == "update_log" {
			r, x := c.update(ctx, *t, ch, log, text)
			if r == nil {
				return result(x, false, "", "", false)
			}
			return result(x, true, r.RawJSON, r.Header.ResultMessage, r.Header.IsSuccessful)
		}
		r, x := c.reply(ctx, *t, ch, log, text)
		if r == nil {
			return result(x, false, "", "", false)
		}
		return result(x, true, r.RawJSON, r.Header.ResultMessage, r.Header.IsSuccessful)
	})
}
func addThread(s *server.MCPServer, t *string, c calls) {
	add(s, "dooray_messenger_channel_thread", "Create a thread or create one from a log", []mcp.ToolOption{str("operation", true, "create_thread", "create_thread_from_log"), str("channelId", true), str("logId", false), str("text", true), str("threadText", false)}, func(ctx context.Context, a map[string]any) (*mcp.CallToolResult, error) {
		o, e := enm(a, "operation", "create_thread", "create_thread_from_log")
		if e != nil {
			return bad(e)
		}
		ch, e := id(a, "channelId")
		if e != nil {
			return bad(e)
		}
		text, e := required(a, "text")
		if e != nil {
			return bad(e)
		}
		if e = tok(t); e != nil {
			return bad(e)
		}
		var r *mm.SendMessageResponse
		if o == "create_thread" {
			tt, x := optional(a, "threadText")
			if x != nil {
				return bad(x)
			}
			r, e = c.thread(ctx, *t, ch, mm.CreateAndSendThreadRequest{Text: text, ThreadText: tt})
		} else {
			log, x := id(a, "logId")
			if x != nil {
				return bad(x)
			}
			r, e = c.threadFromLog(ctx, *t, ch, log, text)
		}
		if r == nil {
			return result(e, false, "", "", false)
		}
		return result(e, true, r.RawJSON, r.Header.ResultMessage, r.Header.IsSuccessful)
	})
}

func add(s *server.MCPServer, name, desc string, opts []mcp.ToolOption, h func(context.Context, map[string]any) (*mcp.CallToolResult, error)) {
	tool := mcp.NewTool(name, append([]mcp.ToolOption{mcp.WithDescription(desc)}, opts...)...)
	s.AddTool(tool, func(ctx context.Context, r mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return h(ctx, r.GetArguments())
	})
}
func str(k string, req bool, values ...string) mcp.ToolOption {
	o := []mcp.PropertyOption{}
	if req {
		o = append(o, mcp.Required())
	}
	if len(values) > 0 {
		o = append(o, mcp.Enum(values...))
	}
	return mcp.WithString(k, o...)
}
func arr(k string) mcp.ToolOption { return mcp.WithArray(k, mcp.Required(), mcp.WithStringItems()) }
func tok(t *string) error {
	if t == nil || strings.TrimSpace(*t) == "" {
		return fmt.Errorf("Dooray token is required")
	}
	return nil
}
func op(a map[string]any, w string) error { _, e := enm(a, "operation", w); return e }
func required(a map[string]any, k string) (string, error) {
	v, ok := a[k].(string)
	if !ok || strings.TrimSpace(v) == "" {
		return "", fmt.Errorf("%s must be a non-empty string", k)
	}
	return v, nil
}
func optional(a map[string]any, k string) (string, error) {
	v, ok := a[k]
	if !ok || v == nil {
		return "", nil
	}
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("%s must be a string", k)
	}
	return s, nil
}
func enm(a map[string]any, k string, values ...string) (string, error) {
	v, e := required(a, k)
	if e != nil {
		return "", e
	}
	for _, x := range values {
		if v == x {
			return v, nil
		}
	}
	return "", fmt.Errorf("%s must be one of %s", k, strings.Join(values, ", "))
}
func id(a map[string]any, k string) (string, error) {
	v, e := required(a, k)
	if e != nil {
		return "", e
	}
	if strings.ContainsAny(v, ",/\\?#% \t\r\n") || v == "." || v == ".." {
		return "", fmt.Errorf("%s must be a single ID", k)
	}
	return v, nil
}
func array(a map[string]any, k string) ([]string, error) {
	raw, ok := a[k]
	if !ok {
		return nil, fmt.Errorf("%s must be a non-empty string array", k)
	}
	var out []string
	switch v := raw.(type) {
	case []string:
		out = v
	case []any:
		for _, x := range v {
			s, ok := x.(string)
			if !ok {
				return nil, fmt.Errorf("%s must contain only strings", k)
			}
			out = append(out, s)
		}
	default:
		return nil, fmt.Errorf("%s must be a non-empty string array", k)
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("%s must be a non-empty string array", k)
	}
	for _, v := range out {
		if strings.TrimSpace(v) == "" {
			return nil, fmt.Errorf("%s must not contain empty values", k)
		}
	}
	return out, nil
}
func bad(e error) (*mcp.CallToolResult, error) { return mcp.NewToolResultError(e.Error()), nil }
func result(e error, p bool, raw, msg string, ok bool) (*mcp.CallToolResult, error) {
	if e != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Dooray request failed: %v", e)), nil
	}
	if !p {
		return mcp.NewToolResultError("Dooray returned an empty response"), nil
	}
	if !ok {
		return mcp.NewToolResultError(fmt.Sprintf("Dooray rejected the request: %s", msg)), nil
	}
	return mcp.NewToolResultText(raw), nil
}
