package wiki

import (
	"context"
	"dooray_mcp/internal/mcptest"
	"errors"
	wikimodel "github.com/dooray-go/dooray-sdk/openapi/model/wiki"
	"github.com/mark3labs/mcp-go/mcp"
	"reflect"
	"strings"
	"testing"
)

func TestWikiCommentsRoutingAndFailures(t *testing.T) {
	for _, mode := range []string{"success", "transport", "nil", "rejected"} {
		t.Run(mode, func(t *testing.T) {
			s := mcptest.NewServer()
			token := "token"
			ctx := t.Context()
			calls := 0
			wikiCommentsTool(s, &token, func(ctx context.Context, token string, wikiID string, pageID string, page, size int) (*wikimodel.ListCommentsResponse, error) {
				calls++
				if ctx != t.Context() || token != "token" {
					t.Error("context or token not forwarded")
				}
				if wikiID != "wikiId-value" {
					t.Errorf("unexpected wikiId: %s", wikiID)
				}
				if pageID != "pageId-value" {
					t.Errorf("unexpected pageId: %s", pageID)
				}
				if page != 2 || size != 15 {
					t.Errorf("pagination %d %d", page, size)
				}
				if mode == "transport" {
					return nil, errors.New("transport failed")
				}
				if mode == "nil" {
					return nil, nil
				}
				res := &wikimodel.ListCommentsResponse{RawJSON: `{"ok":true}`}
				res.Header.IsSuccessful = mode != "rejected"
				res.Header.ResultMessage = "denied"
				return res, nil
			})
			args := map[string]any{"wikiId": "wikiId-value", "pageId": "pageId-value", "operation": "find_comments", "page": 2, "size": 15}
			res, err := s.ListTools()["dooray_wiki_comments"].Handler(ctx, mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}})
			if err != nil || res == nil || calls != 1 {
				t.Fatalf("result=%+v err=%v calls=%d", res, err, calls)
			}
			if res.IsError != (mode != "success") {
				t.Fatalf("unexpected error state: %+v", res)
			}
			if mode == "success" && res.Content[0].(mcp.TextContent).Text != `{"ok":true}` {
				t.Error("raw JSON lost")
			}
		})
	}
}
func TestWikiCommentRoutingAndFailures(t *testing.T) {
	for _, mode := range []string{"success", "transport", "nil", "rejected"} {
		t.Run(mode, func(t *testing.T) {
			s := mcptest.NewServer()
			token := "token"
			ctx := t.Context()
			calls := 0
			wikiCommentTool(s, &token, func(ctx context.Context, token string, wikiID string, pageID string, commentID string) (*wikimodel.GetCommentResponse, error) {
				calls++
				if ctx != t.Context() || token != "token" {
					t.Error("context or token not forwarded")
				}
				if wikiID != "wikiId-value" {
					t.Errorf("unexpected wikiId: %s", wikiID)
				}
				if pageID != "pageId-value" {
					t.Errorf("unexpected pageId: %s", pageID)
				}
				if commentID != "commentId-value" {
					t.Errorf("unexpected commentId: %s", commentID)
				}
				if mode == "transport" {
					return nil, errors.New("transport failed")
				}
				if mode == "nil" {
					return nil, nil
				}
				res := &wikimodel.GetCommentResponse{RawJSON: `{"ok":true}`}
				res.Header.IsSuccessful = mode != "rejected"
				res.Header.ResultMessage = "denied"
				return res, nil
			})
			args := map[string]any{"wikiId": "wikiId-value", "pageId": "pageId-value", "commentId": "commentId-value", "operation": "get_comment"}
			res, err := s.ListTools()["dooray_wiki_comment"].Handler(ctx, mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}})
			if err != nil || res == nil || calls != 1 {
				t.Fatalf("result=%+v err=%v calls=%d", res, err, calls)
			}
			if res.IsError != (mode != "success") {
				t.Fatalf("unexpected error state: %+v", res)
			}
			if mode == "success" && res.Content[0].(mcp.TextContent).Text != `{"ok":true}` {
				t.Error("raw JSON lost")
			}
		})
	}
}
func TestWikiCommentPostRoutingAndFailures(t *testing.T) {
	for _, mode := range []string{"success", "transport", "nil", "rejected"} {
		t.Run(mode, func(t *testing.T) {
			s := mcptest.NewServer()
			token := "token"
			ctx := t.Context()
			calls := 0
			wikiCommentPostTool(s, &token, func(ctx context.Context, token string, wikiID string, pageID string, content string) (*wikimodel.IDResponse, error) {
				calls++
				if ctx != t.Context() || token != "token" {
					t.Error("context or token not forwarded")
				}
				if wikiID != "wikiId-value" {
					t.Errorf("unexpected wikiId: %s", wikiID)
				}
				if pageID != "pageId-value" {
					t.Errorf("unexpected pageId: %s", pageID)
				}
				if content != "comment body" {
					t.Error("content lost")
				}
				if mode == "transport" {
					return nil, errors.New("transport failed")
				}
				if mode == "nil" {
					return nil, nil
				}
				res := &wikimodel.IDResponse{RawJSON: `{"ok":true}`}
				res.Header.IsSuccessful = mode != "rejected"
				res.Header.ResultMessage = "denied"
				return res, nil
			})
			args := map[string]any{"wikiId": "wikiId-value", "pageId": "pageId-value", "operation": "create_comment", "content": "comment body"}
			res, err := s.ListTools()["dooray_wiki_comment_post"].Handler(ctx, mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}})
			if err != nil || res == nil || calls != 1 {
				t.Fatalf("result=%+v err=%v calls=%d", res, err, calls)
			}
			if res.IsError != (mode != "success") {
				t.Fatalf("unexpected error state: %+v", res)
			}
			if mode == "success" && res.Content[0].(mcp.TextContent).Text != `{"ok":true}` {
				t.Error("raw JSON lost")
			}
		})
	}
}
func TestWikiCommentUpdateRoutingAndFailures(t *testing.T) {
	for _, mode := range []string{"success", "transport", "nil", "rejected"} {
		t.Run(mode, func(t *testing.T) {
			s := mcptest.NewServer()
			token := "token"
			ctx := t.Context()
			calls := 0
			wikiCommentUpdateTool(s, &token, func(ctx context.Context, token string, wikiID string, pageID string, commentID string, content string) (*wikimodel.HeaderOnlyResponse, error) {
				calls++
				if ctx != t.Context() || token != "token" {
					t.Error("context or token not forwarded")
				}
				if wikiID != "wikiId-value" {
					t.Errorf("unexpected wikiId: %s", wikiID)
				}
				if pageID != "pageId-value" {
					t.Errorf("unexpected pageId: %s", pageID)
				}
				if commentID != "commentId-value" {
					t.Errorf("unexpected commentId: %s", commentID)
				}
				if content != "comment body" {
					t.Error("content lost")
				}
				if mode == "transport" {
					return nil, errors.New("transport failed")
				}
				if mode == "nil" {
					return nil, nil
				}
				res := &wikimodel.HeaderOnlyResponse{RawJSON: `{"ok":true}`}
				res.Header.IsSuccessful = mode != "rejected"
				res.Header.ResultMessage = "denied"
				return res, nil
			})
			args := map[string]any{"wikiId": "wikiId-value", "pageId": "pageId-value", "commentId": "commentId-value", "operation": "update_comment", "content": "comment body"}
			res, err := s.ListTools()["dooray_wiki_comment_update"].Handler(ctx, mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}})
			if err != nil || res == nil || calls != 1 {
				t.Fatalf("result=%+v err=%v calls=%d", res, err, calls)
			}
			if res.IsError != (mode != "success") {
				t.Fatalf("unexpected error state: %+v", res)
			}
			if mode == "success" && res.Content[0].(mcp.TextContent).Text != `{"ok":true}` {
				t.Error("raw JSON lost")
			}
		})
	}
}
func TestWikiCommentDeleteRoutingAndFailures(t *testing.T) {
	for _, mode := range []string{"success", "transport", "nil", "rejected"} {
		t.Run(mode, func(t *testing.T) {
			s := mcptest.NewServer()
			token := "token"
			ctx := t.Context()
			calls := 0
			wikiCommentDeleteTool(s, &token, func(ctx context.Context, token string, wikiID string, pageID string, commentID string) (*wikimodel.HeaderOnlyResponse, error) {
				calls++
				if ctx != t.Context() || token != "token" {
					t.Error("context or token not forwarded")
				}
				if wikiID != "wikiId-value" {
					t.Errorf("unexpected wikiId: %s", wikiID)
				}
				if pageID != "pageId-value" {
					t.Errorf("unexpected pageId: %s", pageID)
				}
				if commentID != "commentId-value" {
					t.Errorf("unexpected commentId: %s", commentID)
				}
				if mode == "transport" {
					return nil, errors.New("transport failed")
				}
				if mode == "nil" {
					return nil, nil
				}
				res := &wikimodel.HeaderOnlyResponse{RawJSON: `{"ok":true}`}
				res.Header.IsSuccessful = mode != "rejected"
				res.Header.ResultMessage = "denied"
				return res, nil
			})
			args := map[string]any{"wikiId": "wikiId-value", "pageId": "pageId-value", "commentId": "commentId-value", "operation": "delete_comment"}
			res, err := s.ListTools()["dooray_wiki_comment_delete"].Handler(ctx, mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}})
			if err != nil || res == nil || calls != 1 {
				t.Fatalf("result=%+v err=%v calls=%d", res, err, calls)
			}
			if res.IsError != (mode != "success") {
				t.Fatalf("unexpected error state: %+v", res)
			}
			if mode == "success" && res.Content[0].(mcp.TextContent).Text != `{"ok":true}` {
				t.Error("raw JSON lost")
			}
		})
	}
}
func TestWikiSharedLinksRoutingAndFailures(t *testing.T) {
	for _, mode := range []string{"success", "transport", "nil", "rejected"} {
		t.Run(mode, func(t *testing.T) {
			s := mcptest.NewServer()
			token := "token"
			ctx := t.Context()
			calls := 0
			wikiSharedLinksTool(s, &token, func(ctx context.Context, token string, wikiID string, pageID string, page, size int, valid *bool) (*wikimodel.ListSharedLinksResponse, error) {
				calls++
				if ctx != t.Context() || token != "token" {
					t.Error("context or token not forwarded")
				}
				if wikiID != "wikiId-value" {
					t.Errorf("unexpected wikiId: %s", wikiID)
				}
				if pageID != "pageId-value" {
					t.Errorf("unexpected pageId: %s", pageID)
				}
				if page != 2 || size != 15 {
					t.Errorf("pagination %d %d", page, size)
				}
				if valid == nil || *valid {
					t.Error("explicit false valid lost")
				}
				if mode == "transport" {
					return nil, errors.New("transport failed")
				}
				if mode == "nil" {
					return nil, nil
				}
				res := &wikimodel.ListSharedLinksResponse{RawJSON: `{"ok":true}`}
				res.Header.IsSuccessful = mode != "rejected"
				res.Header.ResultMessage = "denied"
				return res, nil
			})
			args := map[string]any{"wikiId": "wikiId-value", "pageId": "pageId-value", "operation": "find_shared_links", "page": 2, "size": 15, "valid": false}
			res, err := s.ListTools()["dooray_wiki_shared_links"].Handler(ctx, mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}})
			if err != nil || res == nil || calls != 1 {
				t.Fatalf("result=%+v err=%v calls=%d", res, err, calls)
			}
			if res.IsError != (mode != "success") {
				t.Fatalf("unexpected error state: %+v", res)
			}
			if mode == "success" && res.Content[0].(mcp.TextContent).Text != `{"ok":true}` {
				t.Error("raw JSON lost")
			}
		})
	}
}
func TestWikiPageDeleteRoutingAndFailures(t *testing.T) {
	for _, mode := range []string{"success", "transport", "nil", "rejected"} {
		t.Run(mode, func(t *testing.T) {
			s := mcptest.NewServer()
			token := "token"
			ctx := t.Context()
			calls := 0
			wikiPageDeleteTool(s, &token, func(ctx context.Context, token string, wikiID string, pageID string) (*wikimodel.HeaderOnlyResponse, error) {
				calls++
				if ctx != t.Context() || token != "token" {
					t.Error("context or token not forwarded")
				}
				if wikiID != "wikiId-value" {
					t.Errorf("unexpected wikiId: %s", wikiID)
				}
				if pageID != "pageId-value" {
					t.Errorf("unexpected pageId: %s", pageID)
				}
				if mode == "transport" {
					return nil, errors.New("transport failed")
				}
				if mode == "nil" {
					return nil, nil
				}
				res := &wikimodel.HeaderOnlyResponse{RawJSON: `{"ok":true}`}
				res.Header.IsSuccessful = mode != "rejected"
				res.Header.ResultMessage = "denied"
				return res, nil
			})
			args := map[string]any{"wikiId": "wikiId-value", "pageId": "pageId-value", "operation": "delete_page"}
			res, err := s.ListTools()["dooray_wiki_page_delete"].Handler(ctx, mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}})
			if err != nil || res == nil || calls != 1 {
				t.Fatalf("result=%+v err=%v calls=%d", res, err, calls)
			}
			if res.IsError != (mode != "success") {
				t.Fatalf("unexpected error state: %+v", res)
			}
			if mode == "success" && res.Content[0].(mcp.TextContent).Text != `{"ok":true}` {
				t.Error("raw JSON lost")
			}
		})
	}
}
func TestWikiPageMoveRoutingAndFailures(t *testing.T) {
	for _, mode := range []string{"success", "transport", "nil", "rejected"} {
		t.Run(mode, func(t *testing.T) {
			s := mcptest.NewServer()
			token := "token"
			ctx := t.Context()
			calls := 0
			wikiPageMoveTool(s, &token, func(ctx context.Context, token string, wikiID string, pageID string, req wikimodel.MovePageRequest) (*wikimodel.MovePageResponse, error) {
				calls++
				if ctx != t.Context() || token != "token" {
					t.Error("context or token not forwarded")
				}
				if wikiID != "wikiId-value" {
					t.Errorf("unexpected wikiId: %s", wikiID)
				}
				if pageID != "pageId-value" {
					t.Errorf("unexpected pageId: %s", pageID)
				}
				if req.TargetParentPageID != "targetParentPageId-value" || req.TargetWikiID != "target-wiki" || req.BeforePageID != "before-page" || req.WithChildren == nil || *req.WithChildren {
					t.Errorf("move request: %+v", req)
				}
				if mode == "transport" {
					return nil, errors.New("transport failed")
				}
				if mode == "nil" {
					return nil, nil
				}
				res := &wikimodel.MovePageResponse{RawJSON: `{"ok":true}`}
				res.Header.IsSuccessful = mode != "rejected"
				res.Header.ResultMessage = "denied"
				return res, nil
			})
			args := map[string]any{"wikiId": "wikiId-value", "pageId": "pageId-value", "targetParentPageId": "targetParentPageId-value", "operation": "move_page", "targetWikiId": "target-wiki", "beforePageId": "before-page", "withChildren": false}
			res, err := s.ListTools()["dooray_wiki_page_move"].Handler(ctx, mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}})
			if err != nil || res == nil || calls != 1 {
				t.Fatalf("result=%+v err=%v calls=%d", res, err, calls)
			}
			if res.IsError != (mode != "success") {
				t.Fatalf("unexpected error state: %+v", res)
			}
			if mode == "success" && res.Content[0].(mcp.TextContent).Text != `{"ok":true}` {
				t.Error("raw JSON lost")
			}
		})
	}
}
func TestWikiPageReferrersRoutingAndFailures(t *testing.T) {
	for _, mode := range []string{"success", "transport", "nil", "rejected"} {
		t.Run(mode, func(t *testing.T) {
			s := mcptest.NewServer()
			token := "token"
			ctx := t.Context()
			calls := 0
			wikiPageReferrersTool(s, &token, func(ctx context.Context, token string, wikiID string, pageID string, referrers []wikimodel.Actor) (*wikimodel.HeaderOnlyResponse, error) {
				calls++
				if ctx != t.Context() || token != "token" {
					t.Error("context or token not forwarded")
				}
				if wikiID != "wikiId-value" {
					t.Errorf("unexpected wikiId: %s", wikiID)
				}
				if pageID != "pageId-value" {
					t.Errorf("unexpected pageId: %s", pageID)
				}
				if !reflect.DeepEqual(referrers, []wikimodel.Actor{{Type: "member", Member: wikimodel.Member{OrganizationMemberID: "member-1"}}}) {
					t.Errorf("referrers: %+v", referrers)
				}
				if mode == "transport" {
					return nil, errors.New("transport failed")
				}
				if mode == "nil" {
					return nil, nil
				}
				res := &wikimodel.HeaderOnlyResponse{RawJSON: `{"ok":true}`}
				res.Header.IsSuccessful = mode != "rejected"
				res.Header.ResultMessage = "denied"
				return res, nil
			})
			args := map[string]any{"wikiId": "wikiId-value", "pageId": "pageId-value", "operation": "update_referrers", "referrerMemberIds": []any{"member-1"}}
			res, err := s.ListTools()["dooray_wiki_page_referrers_update"].Handler(ctx, mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}})
			if err != nil || res == nil || calls != 1 {
				t.Fatalf("result=%+v err=%v calls=%d", res, err, calls)
			}
			if res.IsError != (mode != "success") {
				t.Fatalf("unexpected error state: %+v", res)
			}
			if mode == "success" && res.Content[0].(mcp.TextContent).Text != `{"ok":true}` {
				t.Error("raw JSON lost")
			}
		})
	}
}
func TestWikiFileUploadRoutingAndFailures(t *testing.T) {
	for _, mode := range []string{"success", "transport", "nil", "rejected"} {
		t.Run(mode, func(t *testing.T) {
			s := mcptest.NewServer()
			token := "token"
			ctx := t.Context()
			calls := 0
			wikiFileUploadTool(s, &token, func(ctx context.Context, token string, wikiID string, fileType, filename string, content []byte) (*wikimodel.UploadFileResponse, error) {
				calls++
				if ctx != t.Context() || token != "token" {
					t.Error("context or token not forwarded")
				}
				if wikiID != "wikiId-value" {
					t.Errorf("unexpected wikiId: %s", wikiID)
				}
				if filename != "file.bin" || fileType != "inline_image" || !reflect.DeepEqual(content, []byte{0, 255}) {
					t.Error("upload data lost")
				}
				if mode == "transport" {
					return nil, errors.New("transport failed")
				}
				if mode == "nil" {
					return nil, nil
				}
				res := &wikimodel.UploadFileResponse{RawJSON: `{"ok":true}`}
				res.Header.IsSuccessful = mode != "rejected"
				res.Header.ResultMessage = "denied"
				return res, nil
			})
			args := map[string]any{"wikiId": "wikiId-value", "operation": "upload_wiki_file", "filename": "file.bin", "contentBase64": "AP8=", "fileType": "inline_image"}
			res, err := s.ListTools()["dooray_wiki_file_upload"].Handler(ctx, mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}})
			if err != nil || res == nil || calls != 1 {
				t.Fatalf("result=%+v err=%v calls=%d", res, err, calls)
			}
			if res.IsError != (mode != "success") {
				t.Fatalf("unexpected error state: %+v", res)
			}
			if mode == "success" && res.Content[0].(mcp.TextContent).Text != `{"ok":true}` {
				t.Error("raw JSON lost")
			}
		})
	}
}
func TestWikiPageFileUploadRoutingAndFailures(t *testing.T) {
	for _, mode := range []string{"success", "transport", "nil", "rejected"} {
		t.Run(mode, func(t *testing.T) {
			s := mcptest.NewServer()
			token := "token"
			ctx := t.Context()
			calls := 0
			wikiPageFileUploadTool(s, &token, func(ctx context.Context, token string, wikiID string, pageID string, fileType, filename string, content []byte) (*wikimodel.UploadFileResponse, error) {
				calls++
				if ctx != t.Context() || token != "token" {
					t.Error("context or token not forwarded")
				}
				if wikiID != "wikiId-value" {
					t.Errorf("unexpected wikiId: %s", wikiID)
				}
				if pageID != "pageId-value" {
					t.Errorf("unexpected pageId: %s", pageID)
				}
				if filename != "file.bin" || fileType != "inline_image" || !reflect.DeepEqual(content, []byte{0, 255}) {
					t.Error("upload data lost")
				}
				if mode == "transport" {
					return nil, errors.New("transport failed")
				}
				if mode == "nil" {
					return nil, nil
				}
				res := &wikimodel.UploadFileResponse{RawJSON: `{"ok":true}`}
				res.Header.IsSuccessful = mode != "rejected"
				res.Header.ResultMessage = "denied"
				return res, nil
			})
			args := map[string]any{"wikiId": "wikiId-value", "pageId": "pageId-value", "operation": "upload_page_file", "filename": "file.bin", "contentBase64": "AP8=", "fileType": "inline_image"}
			res, err := s.ListTools()["dooray_wiki_page_file_upload"].Handler(ctx, mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}})
			if err != nil || res == nil || calls != 1 {
				t.Fatalf("result=%+v err=%v calls=%d", res, err, calls)
			}
			if res.IsError != (mode != "success") {
				t.Fatalf("unexpected error state: %+v", res)
			}
			if mode == "success" && res.Content[0].(mcp.TextContent).Text != `{"ok":true}` {
				t.Error("raw JSON lost")
			}
		})
	}
}
func TestWikiAttachFileDownloadRoutingAndFailures(t *testing.T) {
	for _, mode := range []string{"success", "transport", "nil", "rejected"} {
		t.Run(mode, func(t *testing.T) {
			s := mcptest.NewServer()
			token := "token"
			ctx := t.Context()
			calls := 0
			wikiAttachFileDownloadTool(s, &token, func(ctx context.Context, token string, wikiID string, attachFileID string) (*wikimodel.Download, error) {
				calls++
				if ctx != t.Context() || token != "token" {
					t.Error("context or token not forwarded")
				}
				if wikiID != "wikiId-value" {
					t.Errorf("unexpected wikiId: %s", wikiID)
				}
				if attachFileID != "attachFileId-value" {
					t.Errorf("unexpected attachFileId: %s", attachFileID)
				}
				if mode == "transport" {
					return nil, errors.New("transport failed")
				}
				if mode == "nil" {
					return nil, nil
				}
				status := 200
				if mode == "rejected" {
					status = 307
				}
				return &wikimodel.Download{Content: []byte{0, 255}, ContentType: "application/octet-stream", StatusCode: status}, nil
			})
			args := map[string]any{"wikiId": "wikiId-value", "attachFileId": "attachFileId-value", "operation": "download_attach_file"}
			res, err := s.ListTools()["dooray_wiki_attach_file_download"].Handler(ctx, mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}})
			if err != nil || res == nil || calls != 1 {
				t.Fatalf("result=%+v err=%v calls=%d", res, err, calls)
			}
			if res.IsError != (mode != "success") {
				t.Fatalf("unexpected error state: %+v", res)
			}
			if mode == "success" && !strings.Contains(res.Content[0].(mcp.TextContent).Text, `"contentBase64":"AP8="`) {
				t.Errorf("download bytes: %+v", res)
			}
		})
	}
}
func TestWikiPageFileDownloadRoutingAndFailures(t *testing.T) {
	for _, mode := range []string{"success", "transport", "nil", "rejected"} {
		t.Run(mode, func(t *testing.T) {
			s := mcptest.NewServer()
			token := "token"
			ctx := t.Context()
			calls := 0
			wikiPageFileDownloadTool(s, &token, func(ctx context.Context, token string, wikiID string, pageID string, fileID string) (*wikimodel.Download, error) {
				calls++
				if ctx != t.Context() || token != "token" {
					t.Error("context or token not forwarded")
				}
				if wikiID != "wikiId-value" {
					t.Errorf("unexpected wikiId: %s", wikiID)
				}
				if pageID != "pageId-value" {
					t.Errorf("unexpected pageId: %s", pageID)
				}
				if fileID != "fileId-value" {
					t.Errorf("unexpected fileId: %s", fileID)
				}
				if mode == "transport" {
					return nil, errors.New("transport failed")
				}
				if mode == "nil" {
					return nil, nil
				}
				status := 200
				if mode == "rejected" {
					status = 307
				}
				return &wikimodel.Download{Content: []byte{0, 255}, ContentType: "application/octet-stream", StatusCode: status}, nil
			})
			args := map[string]any{"wikiId": "wikiId-value", "pageId": "pageId-value", "fileId": "fileId-value", "operation": "download_page_file"}
			res, err := s.ListTools()["dooray_wiki_page_file_download"].Handler(ctx, mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}})
			if err != nil || res == nil || calls != 1 {
				t.Fatalf("result=%+v err=%v calls=%d", res, err, calls)
			}
			if res.IsError != (mode != "success") {
				t.Fatalf("unexpected error state: %+v", res)
			}
			if mode == "success" && !strings.Contains(res.Content[0].(mcp.TextContent).Text, `"contentBase64":"AP8="`) {
				t.Errorf("download bytes: %+v", res)
			}
		})
	}
}
func TestWikiPageFileDeleteRoutingAndFailures(t *testing.T) {
	for _, mode := range []string{"success", "transport", "nil", "rejected"} {
		t.Run(mode, func(t *testing.T) {
			s := mcptest.NewServer()
			token := "token"
			ctx := t.Context()
			calls := 0
			wikiPageFileDeleteTool(s, &token, func(ctx context.Context, token string, wikiID string, pageID string, fileID string) (*wikimodel.HeaderOnlyResponse, error) {
				calls++
				if ctx != t.Context() || token != "token" {
					t.Error("context or token not forwarded")
				}
				if wikiID != "wikiId-value" {
					t.Errorf("unexpected wikiId: %s", wikiID)
				}
				if pageID != "pageId-value" {
					t.Errorf("unexpected pageId: %s", pageID)
				}
				if fileID != "fileId-value" {
					t.Errorf("unexpected fileId: %s", fileID)
				}
				if mode == "transport" {
					return nil, errors.New("transport failed")
				}
				if mode == "nil" {
					return nil, nil
				}
				res := &wikimodel.HeaderOnlyResponse{RawJSON: `{"ok":true}`}
				res.Header.IsSuccessful = mode != "rejected"
				res.Header.ResultMessage = "denied"
				return res, nil
			})
			args := map[string]any{"wikiId": "wikiId-value", "pageId": "pageId-value", "fileId": "fileId-value", "operation": "delete_page_file"}
			res, err := s.ListTools()["dooray_wiki_page_file_delete"].Handler(ctx, mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}})
			if err != nil || res == nil || calls != 1 {
				t.Fatalf("result=%+v err=%v calls=%d", res, err, calls)
			}
			if res.IsError != (mode != "success") {
				t.Fatalf("unexpected error state: %+v", res)
			}
			if mode == "success" && res.Content[0].(mcp.TextContent).Text != `{"ok":true}` {
				t.Error("raw JSON lost")
			}
		})
	}
}

func TestWikiExtendedValidation(t *testing.T) {
	cases := []struct {
		name string
		args map[string]any
	}{
		{"dooray_wiki_comments", map[string]any{"operation": "find_comments", "wikiId": "w", "pageId": "p", "page": -1}},
		{"dooray_wiki_comments", map[string]any{"operation": "find_comments", "wikiId": "w", "pageId": "p", "size": 1.5}},
		{"dooray_wiki_comments", map[string]any{"operation": "find_comments", "wikiId": "w", "pageId": "p", "page": 1e30}},
		{"dooray_wiki_shared_links", map[string]any{"operation": "find_shared_links", "wikiId": "w", "pageId": "p", "valid": "false"}},
		{"dooray_wiki_comment_post", map[string]any{"operation": "create_comment", "wikiId": "w", "pageId": "p", "content": " "}},
		{"dooray_wiki_comment_update", map[string]any{"operation": "update_comment", "wikiId": "w", "pageId": "p", "commentId": "../x", "content": "a"}},
		{"dooray_wiki_page_move", map[string]any{"operation": "move_page", "wikiId": "w", "pageId": "p", "targetParentPageId": "x", "withChildren": "false"}},
		{"dooray_wiki_page_referrers_update", map[string]any{"operation": "update_referrers", "wikiId": "w", "pageId": "p", "referrerMemberIds": []any{1}}},
		{"dooray_wiki_page_referrers_update", map[string]any{"operation": "update_referrers", "wikiId": "w", "pageId": "p", "referrerMemberIds": []any{"a/b"}}},
		{"dooray_wiki_file_upload", map[string]any{"operation": "upload_wiki_file", "wikiId": "w", "filename": "x", "contentBase64": "!"}},
		{"dooray_wiki_page_file_upload", map[string]any{"operation": "upload_page_file", "wikiId": "w", "pageId": "p", "filename": "x", "contentBase64": "YQ==", "fileType": "other"}},
	}
	s := mcptest.NewServer()
	Tools(s, nil)
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := s.ListTools()[tc.name].Handler(t.Context(), mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: tc.args}})
			if err != nil || res == nil || !res.IsError {
				t.Fatalf("expected error: %+v %v", res, err)
			}
			if strings.Contains(res.Content[0].(mcp.TextContent).Text, "token") {
				t.Fatalf("invalid input reached token validation: %+v", res)
			}
		})
	}
	for name, tool := range s.ListTools() {
		if strings.HasPrefix(name, "dooray_wiki_") {
			t.Run(name+" operation", func(t *testing.T) {
				res, err := tool.Handler(t.Context(), mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{"operation": "invalid"}}})
				if err != nil || res == nil || !res.IsError {
					t.Fatalf("invalid operation accepted: %+v %v", res, err)
				}
			})
		}
	}
}

func TestWikiEmptyReferrersAndUploadDefaults(t *testing.T) {
	referrers, err := wikiReferrers(map[string]any{"referrerMemberIds": []any{}}, true)
	if err != nil || referrers == nil || len(referrers) != 0 {
		t.Fatalf("empty referrers must clear: %+v %v", referrers, err)
	}
	fileType, filename, content, err := wikiUploadArgs(map[string]any{"filename": "empty.txt", "contentBase64": ""})
	if err != nil || fileType != "general" || filename != "empty.txt" || len(content) != 0 {
		t.Fatalf("defaults: %s %s %v %v", fileType, filename, content, err)
	}
}
