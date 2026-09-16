package project

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	projectmodel "github.com/dooray-go/dooray-sdk/openapi/model/project"
	"github.com/mark3labs/mcp-go/mcp"

	"dooray_mcp/internal/mcptest"
)

func TestProjectToolsRegistration(t *testing.T) {
	s := mcptest.NewServer()
	token := "test-token"
	Tools(s, &token)

	tools := s.ListTools()

	names := make(map[string]bool)
	for name := range tools {
		names[name] = true
	}

	if !names["dooray_project"] {
		t.Error("dooray_project tool not registered")
	}
	if !names["dooray_posts"] {
		t.Error("dooray_posts tool not registered")
	}
	if !names["dooray_project_post"] {
		t.Error("dooray_project_post tool not registered")
	}
	if !names["dooray_post"] {
		t.Error("dooray_post tool not registered")
	}
}

func TestProjectToolArguments(t *testing.T) {
	s := mcptest.NewServer()
	token := "invalid-token"
	Tools(s, &token)

	tools := s.ListTools()
	tool, ok := tools["dooray_project"]
	if !ok {
		t.Fatal("dooray_project tool not registered")
	}

	// 핸들러 호출 - API 호출은 실패하지만 인자 파싱이 올바르게 동작하는지 확인
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "dooray_project",
			Arguments: map[string]any{
				"operation": "find_projects",
				"type":      "public",
				"state":     "active",
				"scope":     "private",
			},
		},
	}

	// API 호출 실패는 예상됨 - 인자 파싱 중 panic이 발생하지 않는지 확인
	_, err := tool.Handler(context.Background(), req)
	if err == nil {
		t.Log("handler succeeded (unexpected with invalid token, but argument parsing works)")
	} else {
		t.Logf("handler returned expected error (API call failed): %v", err)
	}
}

func TestPostsToolArguments(t *testing.T) {
	s := mcptest.NewServer()
	token := "invalid-token"
	Tools(s, &token)

	tools := s.ListTools()
	tool, ok := tools["dooray_posts"]
	if !ok {
		t.Fatal("dooray_posts tool not registered")
	}

	// 필수 인자만으로 호출
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "dooray_posts",
			Arguments: map[string]any{
				"operation": "find_posts",
				"projectId": "12345",
			},
		},
	}

	_, err := tool.Handler(context.Background(), req)
	if err == nil {
		t.Log("handler succeeded (unexpected with invalid token, but argument parsing works)")
	} else {
		t.Logf("handler returned expected error (API call failed): %v", err)
	}
}

func TestPostsToolWithAllOptions(t *testing.T) {
	s := mcptest.NewServer()
	token := "invalid-token"
	Tools(s, &token)

	tools := s.ListTools()
	tool, ok := tools["dooray_posts"]
	if !ok {
		t.Fatal("dooray_posts tool not registered")
	}

	// 모든 옵션 인자를 포함하여 호출 - panic 없이 파싱되는지 확인
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "dooray_posts",
			Arguments: map[string]any{
				"operation":           "find_posts",
				"projectId":           "12345",
				"page":                float64(0),
				"size":                float64(20),
				"fromEmailAddress":    "test@example.com",
				"fromMemberIds":       "member1,member2",
				"toMemberIds":         "member3",
				"toMemberSize":        float64(1),
				"ccMemberIds":         "member4",
				"tagIds":              "tag1,tag2",
				"parentPostId":        "parent1",
				"postNumber":          "100",
				"postWorkflowClasses": "registered,working",
				"postWorkflowIds":     "wf1",
				"milestoneIds":        "ms1",
				"subjects":            "test subject",
				"createdAt":           "today",
				"updatedAt":           "prev-7d",
				"dueAt":               "next-30d",
				"order":               "-createdAt",
			},
		},
	}

	_, err := tool.Handler(context.Background(), req)
	if err == nil {
		t.Log("handler succeeded")
	} else {
		t.Logf("handler returned expected error: %v", err)
	}
}

func TestPostsToolMissingRequiredArg(t *testing.T) {
	s := mcptest.NewServer()
	token := "invalid-token"
	Tools(s, &token)

	tools := s.ListTools()
	tool, ok := tools["dooray_posts"]
	if !ok {
		t.Fatal("dooray_posts tool not registered")
	}

	// projectId 누락 시 panic 발생 여부 확인
	defer func() {
		if r := recover(); r != nil {
			t.Logf("recovered from panic (missing required arg): %v", r)
		}
	}()

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "dooray_posts",
			Arguments: map[string]any{
				"operation": "find_posts",
				// projectId 누락
			},
		},
	}

	tool.Handler(context.Background(), req)
}

func TestProjectToolMissingType(t *testing.T) {
	s := mcptest.NewServer()
	token := "invalid-token"
	Tools(s, &token)

	tool := s.ListTools()["dooray_project"]

	defer func() {
		if r := recover(); r != nil {
			t.Logf("recovered from panic (missing type): %v", r)
		}
	}()

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "dooray_project",
			Arguments: map[string]any{
				"operation": "find_projects",
				"state":     "active",
				"scope":     "private",
			},
		},
	}

	tool.Handler(context.Background(), req)
}

func TestProjectToolMissingState(t *testing.T) {
	s := mcptest.NewServer()
	token := "invalid-token"
	Tools(s, &token)

	tool := s.ListTools()["dooray_project"]

	defer func() {
		if r := recover(); r != nil {
			t.Logf("recovered from panic (missing state): %v", r)
		}
	}()

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "dooray_project",
			Arguments: map[string]any{
				"operation": "find_projects",
				"type":      "public",
				"scope":     "private",
			},
		},
	}

	tool.Handler(context.Background(), req)
}

func TestProjectToolMissingScope(t *testing.T) {
	s := mcptest.NewServer()
	token := "invalid-token"
	Tools(s, &token)

	tool := s.ListTools()["dooray_project"]

	defer func() {
		if r := recover(); r != nil {
			t.Logf("recovered from panic (missing scope): %v", r)
		}
	}()

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "dooray_project",
			Arguments: map[string]any{
				"operation": "find_projects",
				"type":      "public",
				"state":     "active",
			},
		},
	}

	tool.Handler(context.Background(), req)
}

func TestProjectToolCount(t *testing.T) {
	s := mcptest.NewServer()
	token := "test-token"
	Tools(s, &token)

	tools := s.ListTools()
	if len(tools) != 4 {
		t.Errorf("expected 4 project tools, got %d", len(tools))
	}
}

func TestPostsToolWithPagingOnly(t *testing.T) {
	s := mcptest.NewServer()
	token := "invalid-token"
	Tools(s, &token)

	tool := s.ListTools()["dooray_posts"]

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "dooray_posts",
			Arguments: map[string]any{
				"operation": "find_posts",
				"projectId": "12345",
				"page":      float64(2),
				"size":      float64(50),
			},
		},
	}

	_, err := tool.Handler(context.Background(), req)
	if err == nil {
		t.Log("handler succeeded with paging args")
	} else {
		t.Logf("handler returned expected error: %v", err)
	}
}

func TestPostsToolWithDateFilters(t *testing.T) {
	s := mcptest.NewServer()
	token := "invalid-token"
	Tools(s, &token)

	tool := s.ListTools()["dooray_posts"]

	testCases := []struct {
		name      string
		createdAt string
		updatedAt string
		dueAt     string
	}{
		{"today pattern", "today", "", ""},
		{"thisweek pattern", "", "thisweek", ""},
		{"prev-N pattern", "prev-30d", "", ""},
		{"next-N pattern", "", "", "next-7d"},
		{"ISO8601 range", "2025-01-01T00:00:00+09:00~2025-12-31T23:59:59+09:00", "", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			args := map[string]any{
				"operation": "find_posts",
				"projectId": "12345",
			}
			if tc.createdAt != "" {
				args["createdAt"] = tc.createdAt
			}
			if tc.updatedAt != "" {
				args["updatedAt"] = tc.updatedAt
			}
			if tc.dueAt != "" {
				args["dueAt"] = tc.dueAt
			}

			req := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name:      "dooray_posts",
					Arguments: args,
				},
			}

			_, err := tool.Handler(context.Background(), req)
			if err == nil {
				t.Log("handler succeeded")
			} else {
				t.Logf("handler returned expected error: %v", err)
			}
		})
	}
}

func TestPostsToolWithSortOptions(t *testing.T) {
	s := mcptest.NewServer()
	token := "invalid-token"
	Tools(s, &token)

	tool := s.ListTools()["dooray_posts"]

	sortOrders := []string{"createdAt", "-createdAt", "postDueAt", "-postDueAt", "postUpdatedAt", "-postUpdatedAt"}

	for _, order := range sortOrders {
		t.Run("order_"+order, func(t *testing.T) {
			req := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "dooray_posts",
					Arguments: map[string]any{
						"operation": "find_posts",
						"projectId": "12345",
						"order":     order,
					},
				},
			}

			_, err := tool.Handler(context.Background(), req)
			if err == nil {
				t.Log("handler succeeded")
			} else {
				t.Logf("handler returned expected error: %v", err)
			}
		})
	}
}

func TestCreatePostPayload(t *testing.T) {
	for _, withUsers := range []bool{false, true} {
		t.Run(fmt.Sprint(withUsers), func(t *testing.T) {
			s := mcptest.NewServer()
			token := "test-token"
			called := 0
			ctx := t.Context()
			createPostTool(s, &token, func(gotCtx context.Context, gotToken, projectID string, post projectmodel.PostRequest) (*projectmodel.PostResponse, error) {
				called++
				if gotCtx != ctx || gotToken != token || projectID != "12345" {
					t.Errorf("incorrect request context/token/project")
				}
				if post.Subject != "Task title" || post.Body.Content != "**Task body**" {
					t.Errorf("incorrect post: %+v", post)
				}
				if withUsers {
					if post.Body.MimeType != "text/html" || post.Users == nil {
						t.Fatalf("incorrect optional payload: %+v", post)
					}
					if len(post.Users.To) != 2 || post.Users.To[0].Type != "member" || post.Users.To[0].Member.OrganizationMemberID != "member1" || post.Users.To[1].Member.OrganizationMemberID != "member2" || len(post.Users.Cc) != 1 || post.Users.Cc[0].Member.OrganizationMemberID != "member3" {
						t.Errorf("incorrect recipients: %+v", post.Users)
					}
				} else if post.Body.MimeType != "text/x-markdown" || post.Users != nil {
					t.Errorf("incorrect defaults: %+v", post)
				}
				res := &projectmodel.PostResponse{RawJSON: `{"header":{"isSuccessful":true},"result":{"id":"task1"}}`}
				res.Header.IsSuccessful = true
				res.Result.ID = "task1"
				return res, nil
			})
			args := createPostArguments()
			if withUsers {
				args["mimeType"] = "text/html"
				args["toMemberIds"] = "member1, member2"
				args["ccMemberIds"] = "member3"
			}
			result, err := s.ListTools()["dooray_project_post"].Handler(ctx, mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}})
			if err != nil || result == nil || result.IsError || called != 1 {
				t.Fatalf("creation failed: result=%+v err=%v calls=%d", result, err, called)
			}
			if result.Content[0].(mcp.TextContent).Text != `{"header":{"isSuccessful":true},"result":{"id":"task1"}}` {
				t.Errorf("response not preserved: %+v", result)
			}
		})
	}
}

func createPostArguments() map[string]any {
	return map[string]any{"operation": "create_post", "projectId": "12345", "subject": "Task title", "content": "**Task body**"}
}

func TestCreatePostRejectsInvalidArgumentsWithoutCallingAPI(t *testing.T) {
	cases := []struct {
		name, key string
		value     any
		remove    bool
	}{
		{"missing operation", "operation", nil, true}, {"missing project", "projectId", nil, true},
		{"missing subject", "subject", nil, true}, {"missing content", "content", nil, true},
		{"wrong type", "subject", 12, false}, {"blank title", "subject", "  ", false},
		{"unknown operation", "operation", "delete_post", false}, {"multiple projects", "projectId", "1,2", false},
		{"path injection", "projectId", "1/../2", false}, {"encoded path", "projectId", "%2f", false},
		{"bad mime", "mimeType", "application/json", false}, {"bad optional type", "toMemberIds", true, false},
		{"empty member", "ccMemberIds", "1,,2", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := mcptest.NewServer()
			token := "test-token"
			createPostTool(s, &token, func(context.Context, string, string, projectmodel.PostRequest) (*projectmodel.PostResponse, error) {
				t.Fatal("API must not be called")
				return nil, nil
			})
			args := createPostArguments()
			if tc.remove {
				delete(args, tc.key)
			} else {
				args[tc.key] = tc.value
			}
			result, err := s.ListTools()["dooray_project_post"].Handler(t.Context(), mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}})
			if err != nil || result == nil || !result.IsError {
				t.Fatalf("expected tool error: %+v %v", result, err)
			}
			if !strings.Contains(result.Content[0].(mcp.TextContent).Text, tc.key) {
				t.Errorf("error does not identify argument: %+v", result)
			}
		})
	}
}

func TestCreatePostReportsAPIFailures(t *testing.T) {
	rejected := &projectmodel.PostResponse{}
	rejected.Header.ResultMessage = "permission denied"
	cases := []struct {
		name     string
		response *projectmodel.PostResponse
		err      error
		message  string
	}{
		{"transport", nil, errors.New("network unavailable"), "network unavailable"},
		{"rejected", rejected, nil, "permission denied"},
		{"empty response", nil, nil, "empty response"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := mcptest.NewServer()
			token := "test-token"
			createPostTool(s, &token, func(context.Context, string, string, projectmodel.PostRequest) (*projectmodel.PostResponse, error) {
				return tc.response, tc.err
			})
			result, err := s.ListTools()["dooray_project_post"].Handler(t.Context(), mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: createPostArguments()}})
			if err != nil || result == nil || !result.IsError || !strings.Contains(result.Content[0].(mcp.TextContent).Text, tc.message) {
				t.Fatalf("expected %q failure: %+v %v", tc.message, result, err)
			}
		})
	}
}

func TestCreatePostRequiresToken(t *testing.T) {
	for _, token := range []*string{nil, new(""), new(" ")} {
		s := mcptest.NewServer()
		createPostTool(s, token, func(context.Context, string, string, projectmodel.PostRequest) (*projectmodel.PostResponse, error) {
			t.Fatal("API must not be called without a token")
			return nil, nil
		})
		result, err := s.ListTools()["dooray_project_post"].Handler(t.Context(), mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: createPostArguments()}})
		if err != nil || result == nil || !result.IsError || !strings.Contains(result.Content[0].(mcp.TextContent).Text, "token") {
			t.Fatalf("expected missing token error: %+v %v", result, err)
		}
	}
}

func getPostArguments() map[string]any {
	return map[string]any{"operation": "get_post", "projectId": "proj-1", "postId": "post-1"}
}

func TestGetPostReturnsRawJSON(t *testing.T) {
	s := mcptest.NewServer()
	token := "test-token"
	called := 0
	ctx := t.Context()
	raw := `{"header":{"isSuccessful":true},"result":{"id":"post-1","subject":"업무 상세"}}`
	getPostTool(s, &token, func(gotCtx context.Context, gotToken, projectID, postID string) (*projectmodel.GetPostResponse, error) {
		called++
		if gotCtx != ctx || gotToken != token || projectID != "proj-1" || postID != "post-1" {
			t.Errorf("incorrect request context/token/ids")
		}
		res := &projectmodel.GetPostResponse{RawJSON: raw}
		res.Header.IsSuccessful = true
		res.Result.ID = "post-1"
		return res, nil
	})
	result, err := s.ListTools()["dooray_post"].Handler(ctx, mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: getPostArguments()}})
	if err != nil || result == nil || result.IsError || called != 1 {
		t.Fatalf("lookup failed: result=%+v err=%v calls=%d", result, err, called)
	}
	if result.Content[0].(mcp.TextContent).Text != raw {
		t.Errorf("response not preserved: %+v", result)
	}
}

func TestGetPostRejectsInvalidArgumentsWithoutCallingAPI(t *testing.T) {
	cases := []struct {
		name, key string
		value     any
		remove    bool
	}{
		{"missing operation", "operation", nil, true},
		{"missing project", "projectId", nil, true},
		{"missing post", "postId", nil, true},
		{"wrong type", "postId", 12, false},
		{"blank post", "postId", "  ", false},
		{"unknown operation", "operation", "find_posts", false},
		{"multiple projects", "projectId", "1,2", false},
		{"path injection", "postId", "1/../2", false},
		{"encoded path", "projectId", "%2f", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := mcptest.NewServer()
			token := "test-token"
			getPostTool(s, &token, func(context.Context, string, string, string) (*projectmodel.GetPostResponse, error) {
				t.Fatal("API must not be called")
				return nil, nil
			})
			args := getPostArguments()
			if tc.remove {
				delete(args, tc.key)
			} else {
				args[tc.key] = tc.value
			}
			result, err := s.ListTools()["dooray_post"].Handler(t.Context(), mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}})
			if err != nil || result == nil || !result.IsError {
				t.Fatalf("expected tool error: %+v %v", result, err)
			}
			if !strings.Contains(result.Content[0].(mcp.TextContent).Text, tc.key) {
				t.Errorf("error does not identify argument: %+v", result)
			}
		})
	}
}

func TestGetPostReportsAPIFailures(t *testing.T) {
	rejected := &projectmodel.GetPostResponse{}
	rejected.Header.ResultMessage = "not found"
	cases := []struct {
		name     string
		response *projectmodel.GetPostResponse
		err      error
		message  string
	}{
		{"transport", nil, errors.New("network unavailable"), "network unavailable"},
		{"rejected", rejected, nil, "not found"},
		{"empty response", nil, nil, "empty response"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := mcptest.NewServer()
			token := "test-token"
			getPostTool(s, &token, func(context.Context, string, string, string) (*projectmodel.GetPostResponse, error) {
				return tc.response, tc.err
			})
			result, err := s.ListTools()["dooray_post"].Handler(t.Context(), mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: getPostArguments()}})
			if err != nil || result == nil || !result.IsError || !strings.Contains(result.Content[0].(mcp.TextContent).Text, tc.message) {
				t.Fatalf("expected %q failure: %+v %v", tc.message, result, err)
			}
		})
	}
}

func TestGetPostRequiresToken(t *testing.T) {
	for _, token := range []*string{nil, new(""), new(" ")} {
		s := mcptest.NewServer()
		getPostTool(s, token, func(context.Context, string, string, string) (*projectmodel.GetPostResponse, error) {
			t.Fatal("API must not be called without a token")
			return nil, nil
		})
		result, err := s.ListTools()["dooray_post"].Handler(t.Context(), mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: getPostArguments()}})
		if err != nil || result == nil || !result.IsError || !strings.Contains(result.Content[0].(mcp.TextContent).Text, "token") {
			t.Fatalf("expected missing token error: %+v %v", result, err)
		}
	}
}
