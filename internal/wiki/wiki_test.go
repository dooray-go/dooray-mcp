package wiki

import (
	"context"
	"errors"
	"strings"
	"testing"

	wikimodel "github.com/dooray-go/dooray-sdk/openapi/model/wiki"
	"github.com/mark3labs/mcp-go/mcp"

	"dooray_mcp/internal/mcptest"
)

func TestWikiToolsRegistration(t *testing.T) {
	s := mcptest.NewServer()
	token := "test-token"
	Tools(s, &token)
	tools := s.ListTools()
	for _, name := range []string{"dooray_wikis", "dooray_wiki_pages", "dooray_wiki_page", "dooray_wiki_page_post", "dooray_wiki_page_update"} {
		if _, ok := tools[name]; !ok {
			t.Errorf("%s tool not registered", name)
		}
	}
	if len(tools) != 19 {
		t.Errorf("expected 19 wiki tools, got %d", len(tools))
	}
}

func TestWikiListReturnsRawJSON(t *testing.T) {
	s := mcptest.NewServer()
	token := "test-token"
	ctx := t.Context()
	called := 0
	raw := `{"header":{"isSuccessful":true},"result":[{"id":"w1"}]}`
	wikiListTool(s, &token, func(gotCtx context.Context, gotToken string, page, size int) (*wikimodel.ListWikisResponse, error) {
		called++
		if gotCtx != ctx || gotToken != token || page != 1 || size != 20 {
			t.Errorf("unexpected list args page=%d size=%d", page, size)
		}
		res := &wikimodel.ListWikisResponse{RawJSON: raw}
		res.Header.IsSuccessful = true
		return res, nil
	})
	result, err := s.ListTools()["dooray_wikis"].Handler(ctx, mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{
		"operation": "find_wikis", "page": float64(1), "size": float64(20),
	}}})
	if err != nil || result == nil || result.IsError || called != 1 {
		t.Fatalf("list failed: %+v %v calls=%d", result, err, called)
	}
	if result.Content[0].(mcp.TextContent).Text != raw {
		t.Errorf("raw JSON not preserved: %+v", result)
	}
}

func TestWikiPagesAndGetPage(t *testing.T) {
	s := mcptest.NewServer()
	token := "test-token"
	ctx := t.Context()
	wikiPagesTool(s, &token, func(_ context.Context, _, wikiID, parent string) (*wikimodel.ListPagesResponse, error) {
		if wikiID != "wiki-1" || parent != "page-10" {
			t.Errorf("ids wiki=%s parent=%s", wikiID, parent)
		}
		res := &wikimodel.ListPagesResponse{RawJSON: `{"pages":true}`}
		res.Header.IsSuccessful = true
		return res, nil
	})
	wikiGetPageTool(s, &token,
		func(context.Context, string, string) (*wikimodel.GetPageResponse, error) {
			t.Fatal("GetPage should not be used when wikiId is set")
			return nil, nil
		},
		func(_ context.Context, _, wikiID, pageID string) (*wikimodel.GetPageResponse, error) {
			if wikiID != "wiki-1" || pageID != "page-10" {
				t.Errorf("get wiki=%s page=%s", wikiID, pageID)
			}
			res := &wikimodel.GetPageResponse{RawJSON: `{"page":true}`}
			res.Header.IsSuccessful = true
			return res, nil
		},
	)
	listRes, err := s.ListTools()["dooray_wiki_pages"].Handler(ctx, mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{
		"operation": "find_pages", "wikiId": "wiki-1", "parentPageId": "page-10",
	}}})
	if err != nil || listRes == nil || listRes.IsError || listRes.Content[0].(mcp.TextContent).Text != `{"pages":true}` {
		t.Fatalf("find_pages: %+v %v", listRes, err)
	}
	getRes, err := s.ListTools()["dooray_wiki_page"].Handler(ctx, mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{
		"operation": "get_page", "wikiId": "wiki-1", "pageId": "page-10",
	}}})
	if err != nil || getRes == nil || getRes.IsError || getRes.Content[0].(mcp.TextContent).Text != `{"page":true}` {
		t.Fatalf("get_page: %+v %v", getRes, err)
	}
}

func TestWikiGetPageWithoutWikiID(t *testing.T) {
	s := mcptest.NewServer()
	token := "test-token"
	wikiGetPageTool(s, &token,
		func(_ context.Context, _, pageID string) (*wikimodel.GetPageResponse, error) {
			if pageID != "page-10" {
				t.Errorf("pageID=%s", pageID)
			}
			res := &wikimodel.GetPageResponse{RawJSON: `{"byPageId":true}`}
			res.Header.IsSuccessful = true
			return res, nil
		},
		func(context.Context, string, string, string) (*wikimodel.GetPageResponse, error) {
			t.Fatal("GetWikiPage should not be used without wikiId")
			return nil, nil
		},
	)
	result, err := s.ListTools()["dooray_wiki_page"].Handler(t.Context(), mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{
		"operation": "get_page", "pageId": "page-10",
	}}})
	if err != nil || result == nil || result.IsError || result.Content[0].(mcp.TextContent).Text != `{"byPageId":true}` {
		t.Fatalf("get by pageId: %+v %v", result, err)
	}
}

func TestWikiCreatePage(t *testing.T) {
	s := mcptest.NewServer()
	token := "test-token"
	ctx := t.Context()
	wikiCreatePageTool(s, &token, func(_ context.Context, _, wikiID string, req wikimodel.CreatePageRequest) (*wikimodel.CreatePageResponse, error) {
		if wikiID != "wiki-1" || req.ParentPageID != "page-10" || req.Subject != "제목" || req.Body.Content != "본문" || req.Body.MimeType != "text/html" || len(req.AttachFileIDs) != 1 || req.AttachFileIDs[0] != "attachment-1" || len(req.Referrers) != 1 || req.Referrers[0].Member.OrganizationMemberID != "member-1" {
			t.Errorf("create payload wiki=%s req=%+v", wikiID, req)
		}
		res := &wikimodel.CreatePageResponse{RawJSON: `{"id":"new"}`}
		res.Header.IsSuccessful = true
		res.Result.ID = "new"
		return res, nil
	})
	result, err := s.ListTools()["dooray_wiki_page_post"].Handler(ctx, mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{
		"operation": "create_page", "wikiId": "wiki-1", "parentPageId": "page-10",
		"subject": "제목", "content": "본문", "mimeType": "text/html",
		"attachFileIds": []any{"attachment-1"}, "referrerMemberIds": []any{"member-1"},
	}}})
	if err != nil || result == nil || result.IsError || result.Content[0].(mcp.TextContent).Text != `{"id":"new"}` {
		t.Fatalf("create: %+v %v", result, err)
	}
}

func TestWikiUpdatePageRoutes(t *testing.T) {
	cases := []struct {
		name        string
		args        map[string]any
		wantUpdate  bool
		wantTitle   bool
		wantContent bool
		wantSubject string
		wantBody    string
	}{
		{"both", map[string]any{"subject": "새 제목", "content": "새 본문"}, true, false, false, "새 제목", "새 본문"},
		{"title only", map[string]any{"subject": "새 제목"}, false, true, false, "새 제목", ""},
		{"content only", map[string]any{"content": "새 본문"}, false, false, true, "", "새 본문"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := mcptest.NewServer()
			token := "test-token"
			wikiUpdatePageTool(s, &token,
				func(_ context.Context, _, wikiID, pageID string, req wikimodel.UpdatePageRequest) (*wikimodel.HeaderOnlyResponse, error) {
					if !tc.wantUpdate {
						t.Fatal("UpdatePage should not be called")
					}
					if wikiID != "wiki-1" || pageID != "page-10" || req.Subject != tc.wantSubject || req.Body.Content != tc.wantBody {
						t.Errorf("update %+v", req)
					}
					res := &wikimodel.HeaderOnlyResponse{RawJSON: `{"updated":true}`}
					res.Header.IsSuccessful = true
					return res, nil
				},
				func(_ context.Context, _, wikiID, pageID, subject string) (*wikimodel.HeaderOnlyResponse, error) {
					if !tc.wantTitle {
						t.Fatal("UpdatePageTitle should not be called")
					}
					if subject != tc.wantSubject {
						t.Errorf("title %s", subject)
					}
					res := &wikimodel.HeaderOnlyResponse{RawJSON: `{"updated":true}`}
					res.Header.IsSuccessful = true
					return res, nil
				},
				func(_ context.Context, _, wikiID, pageID string, body wikimodel.Body) (*wikimodel.HeaderOnlyResponse, error) {
					if !tc.wantContent {
						t.Fatal("UpdatePageContent should not be called")
					}
					if body.Content != tc.wantBody || body.MimeType != "text/x-markdown" {
						t.Errorf("body %+v", body)
					}
					res := &wikimodel.HeaderOnlyResponse{RawJSON: `{"updated":true}`}
					res.Header.IsSuccessful = true
					return res, nil
				},
			)
			args := map[string]any{"operation": "update_page", "wikiId": "wiki-1", "pageId": "page-10"}
			for k, v := range tc.args {
				args[k] = v
			}
			result, err := s.ListTools()["dooray_wiki_page_update"].Handler(t.Context(), mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}})
			if err != nil || result == nil || result.IsError {
				t.Fatalf("update: %+v %v", result, err)
			}
		})
	}
}

func TestWikiToolsRejectBadInput(t *testing.T) {
	s := mcptest.NewServer()
	token := "test-token"
	wikiListTool(s, &token, func(context.Context, string, int, int) (*wikimodel.ListWikisResponse, error) {
		t.Fatal("API must not be called")
		return nil, nil
	})
	wikiPagesTool(s, &token, func(context.Context, string, string, string) (*wikimodel.ListPagesResponse, error) {
		t.Fatal("API must not be called")
		return nil, nil
	})
	wikiGetPageTool(s, &token, func(context.Context, string, string) (*wikimodel.GetPageResponse, error) {
		t.Fatal("API must not be called")
		return nil, nil
	}, func(context.Context, string, string, string) (*wikimodel.GetPageResponse, error) {
		t.Fatal("API must not be called")
		return nil, nil
	})
	wikiCreatePageTool(s, &token, func(context.Context, string, string, wikimodel.CreatePageRequest) (*wikimodel.CreatePageResponse, error) {
		t.Fatal("API must not be called")
		return nil, nil
	})
	wikiUpdatePageTool(s, &token,
		func(context.Context, string, string, string, wikimodel.UpdatePageRequest) (*wikimodel.HeaderOnlyResponse, error) {
			t.Fatal("API must not be called")
			return nil, nil
		},
		func(context.Context, string, string, string, string) (*wikimodel.HeaderOnlyResponse, error) {
			t.Fatal("API must not be called")
			return nil, nil
		},
		func(context.Context, string, string, string, wikimodel.Body) (*wikimodel.HeaderOnlyResponse, error) {
			t.Fatal("API must not be called")
			return nil, nil
		},
	)
	cases := []struct {
		tool, key string
		args      map[string]any
	}{
		{"dooray_wikis", "operation", map[string]any{"operation": "get_page"}},
		{"dooray_wiki_pages", "wikiId", map[string]any{"operation": "find_pages", "wikiId": "1,2"}},
		{"dooray_wiki_page", "pageId", map[string]any{"operation": "get_page"}},
		{"dooray_wiki_page_post", "subject", map[string]any{"operation": "create_page", "wikiId": "w", "content": "c"}},
		{"dooray_wiki_page_update", "subject", map[string]any{"operation": "update_page", "wikiId": "w", "pageId": "p"}},
	}
	for _, tc := range cases {
		t.Run(tc.tool, func(t *testing.T) {
			result, err := s.ListTools()[tc.tool].Handler(t.Context(), mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: tc.args}})
			if err != nil || result == nil || !result.IsError || !strings.Contains(result.Content[0].(mcp.TextContent).Text, tc.key) {
				t.Fatalf("expected %s in error: %+v %v", tc.key, result, err)
			}
		})
	}
}

func TestWikiToolsRequireTokenAndReportAPIErrors(t *testing.T) {
	s := mcptest.NewServer()
	wikiListTool(s, nil, func(context.Context, string, int, int) (*wikimodel.ListWikisResponse, error) {
		t.Fatal("API must not be called")
		return nil, nil
	})
	result, err := s.ListTools()["dooray_wikis"].Handler(t.Context(), mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{"operation": "find_wikis"}}})
	if err != nil || result == nil || !result.IsError || !strings.Contains(result.Content[0].(mcp.TextContent).Text, "token") {
		t.Fatalf("token: %+v %v", result, err)
	}

	token := "test-token"
	s = mcptest.NewServer()
	wikiGetPageTool(s, &token, func(context.Context, string, string) (*wikimodel.GetPageResponse, error) {
		return nil, errors.New("not found")
	}, func(context.Context, string, string, string) (*wikimodel.GetPageResponse, error) {
		t.Fatal("unused")
		return nil, nil
	})
	result, err = s.ListTools()["dooray_wiki_page"].Handler(t.Context(), mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{"operation": "get_page", "pageId": "p"}}})
	if err != nil || result == nil || !result.IsError || !strings.Contains(result.Content[0].(mcp.TextContent).Text, "not found") {
		t.Fatalf("api error: %+v %v", result, err)
	}
}
