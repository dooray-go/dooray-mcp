package messenger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	model "github.com/dooray-go/dooray-sdk/openapi/model"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// GET /messenger/v1/channels/{id}/logs 는 공식 문서에 없는 경로라 SDK 가 감싸지 않는다.
// 최근 size 건만 주고(상한 1000) 페이징이 없어 그 이전 메시지는 가져올 수 없다.
const (
	logsEndpoint    = "https://api.dooray.com"
	defaultLogsSize = 50
	maxLogsSize     = 1000
)

type logsResponse struct {
	Header  model.ResponseHeader
	RawJSON string
}

func fetchChannelLogs(ctx context.Context, token, channelID string, size int) (*logsResponse, error) {
	u := fmt.Sprintf("%s/messenger/v1/channels/%s/logs?size=%d", logsEndpoint, url.PathEscape(channelID), size)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "dooray-api "+token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var parsed struct {
		Header model.ResponseHeader `json:"header"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("status %d: %s", res.StatusCode, string(body))
	}
	return &logsResponse{Header: parsed.Header, RawJSON: string(body)}, nil
}

func addLogs(s *server.MCPServer, t *string, c calls) {
	opts := []mcp.ToolOption{
		str("operation", true, "find_logs"),
		str("channelId", true),
		mcp.WithNumber("size", mcp.Description("1~1000, default 50. Only the most recent messages are returned; no paging.")),
	}
	add(s, "dooray_messenger_channel_logs", "Read recent messages of a channel (or thread channel)", opts, func(ctx context.Context, a map[string]any) (*mcp.CallToolResult, error) {
		if e := op(a, "find_logs"); e != nil {
			return bad(e)
		}
		ch, e := id(a, "channelId")
		if e != nil {
			return bad(e)
		}
		size, e := logsSize(a)
		if e != nil {
			return bad(e)
		}
		if e = tok(t); e != nil {
			return bad(e)
		}
		r, e := c.logs(ctx, *t, ch, size)
		if r == nil {
			return result(e, false, "", "", false)
		}
		return result(e, true, r.RawJSON, r.Header.ResultMessage, r.Header.IsSuccessful)
	})
}

func logsSize(a map[string]any) (int, error) {
	v, ok := a["size"]
	if !ok || v == nil {
		return defaultLogsSize, nil
	}
	var n int
	switch x := v.(type) {
	case float64:
		n = int(x)
		if float64(n) != x {
			return 0, fmt.Errorf("size must be an integer")
		}
	case string:
		p, err := strconv.Atoi(x)
		if err != nil {
			return 0, fmt.Errorf("size must be an integer")
		}
		n = p
	default:
		return 0, fmt.Errorf("size must be an integer")
	}
	if n < 1 || n > maxLogsSize {
		return 0, fmt.Errorf("size must be between 1 and %d", maxLogsSize)
	}
	return n, nil
}
