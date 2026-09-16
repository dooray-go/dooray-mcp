package wiki

import (
	"context"
	"fmt"
	"strings"

	wikimodel "github.com/dooray-go/dooray-sdk/openapi/model/wiki"
	"github.com/dooray-go/dooray-sdk/openapi/wiki"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func Tools(s *server.MCPServer, token *string) {
	client := wiki.NewDefaultWiki()
	wikiListTool(s, token, func(ctx context.Context, apikey string, page, size int) (*wikimodel.ListWikisResponse, error) {
		return client.GetWikisCustomHTTPContext(ctx, apikey, nil, page, size)
	})
	wikiPagesTool(s, token, client.GetPagesContext)
	wikiGetPageTool(s, token, client.GetPageContext, client.GetWikiPageContext)
	wikiCreatePageTool(s, token, client.CreatePageContext)
	wikiUpdatePageTool(s, token, client.UpdatePageContext, client.UpdatePageTitleContext, client.UpdatePageContentContext)
}

func wikiListTool(s *server.MCPServer, token *string, list func(context.Context, string, int, int) (*wikimodel.ListWikisResponse, error)) {
	tool := mcp.NewTool("dooray_wikis",
		mcp.WithDescription("List Dooray wikis the token can access. Home page IDs are in result[].home.pageId."),
		mcp.WithString("operation", mcp.Required(), mcp.Enum("find_wikis")),
		mcp.WithNumber("page", mcp.Description("page number, default 0")),
		mcp.WithNumber("size", mcp.Description("page size, default is the Dooray API default")),
	)
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		if err := wikiRequireOperation(args, "find_wikis"); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		if err := wikiRequireToken(token); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		page := wikiIntArg(args, "page")
		size := wikiIntArg(args, "size")
		res, err := list(ctx, *token, page, size)
		if res == nil {
			return wikiAPIResult(err, "failed to list wikis", false, "", "", false)
		}
		return wikiAPIResult(err, "failed to list wikis", true, res.RawJSON, res.Header.ResultMessage, res.Header.IsSuccessful)
	})
}

func wikiPagesTool(s *server.MCPServer, token *string, list func(context.Context, string, string, string) (*wikimodel.ListPagesResponse, error)) {
	tool := mcp.NewTool("dooray_wiki_pages",
		mcp.WithDescription("List one level of Dooray wiki pages. Omit parentPageId for root children; pass a page id to list its children."),
		mcp.WithString("operation", mcp.Required(), mcp.Enum("find_pages")),
		mcp.WithString("wikiId", mcp.Required(), mcp.Description("Wiki ID from dooray_wikis")),
		mcp.WithString("parentPageId", mcp.Description("Parent page ID. Omit to list the wiki root's children")),
	)
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		if err := wikiRequireOperation(args, "find_pages"); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		wikiID, err := wikiRequiredID(args, "wikiId")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		parentPageID, err := wikiOptionalID(args, "parentPageId")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		if err := wikiRequireToken(token); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		res, err := list(ctx, *token, wikiID, parentPageID)
		if res == nil {
			return wikiAPIResult(err, "failed to list wiki pages", false, "", "", false)
		}
		return wikiAPIResult(err, "failed to list wiki pages", true, res.RawJSON, res.Header.ResultMessage, res.Header.IsSuccessful)
	})
}

func wikiGetPageTool(s *server.MCPServer, token *string, getPage func(context.Context, string, string) (*wikimodel.GetPageResponse, error), getWikiPage func(context.Context, string, string, string) (*wikimodel.GetPageResponse, error)) {
	tool := mcp.NewTool("dooray_wiki_page",
		mcp.WithDescription("Get a Dooray wiki page including body, referrers, files, and images. Prefer pageId from dooray_wiki_pages or dooray_wikis home.pageId."),
		mcp.WithString("operation", mcp.Required(), mcp.Enum("get_page")),
		mcp.WithString("pageId", mcp.Required(), mcp.Description("Page ID")),
		mcp.WithString("wikiId", mcp.Description("Wiki ID. When set, calls the wiki-scoped page API")),
	)
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		if err := wikiRequireOperation(args, "get_page"); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		pageID, err := wikiRequiredID(args, "pageId")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		wikiID, err := wikiOptionalID(args, "wikiId")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		if err := wikiRequireToken(token); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		var res *wikimodel.GetPageResponse
		if wikiID != "" {
			res, err = getWikiPage(ctx, *token, wikiID, pageID)
		} else {
			res, err = getPage(ctx, *token, pageID)
		}
		if res == nil {
			return wikiAPIResult(err, "failed to get wiki page", false, "", "", false)
		}
		return wikiAPIResult(err, "failed to get wiki page", true, res.RawJSON, res.Header.ResultMessage, res.Header.IsSuccessful)
	})
}

func wikiRequireOperation(args map[string]any, want string) error {
	v, ok := args["operation"].(string)
	if !ok || strings.TrimSpace(v) == "" {
		return fmt.Errorf("operation must be a non-empty string")
	}
	if v != want {
		return fmt.Errorf("operation must be %s", want)
	}
	return nil
}

func wikiRequireToken(token *string) error {
	if token == nil || strings.TrimSpace(*token) == "" {
		return fmt.Errorf("Dooray token is required")
	}
	return nil
}

func wikiRequiredID(args map[string]any, key string) (string, error) {
	v, ok := args[key]
	if !ok {
		return "", fmt.Errorf("%s must be a non-empty string", key)
	}
	id, ok := v.(string)
	if !ok || strings.TrimSpace(id) == "" {
		return "", fmt.Errorf("%s must be a non-empty string", key)
	}
	if err := wikiSingleID(key, id); err != nil {
		return "", err
	}
	return id, nil
}

func wikiOptionalID(args map[string]any, key string) (string, error) {
	v, ok := args[key]
	if !ok || v == nil {
		return "", nil
	}
	id, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("%s must be a non-empty string", key)
	}
	if strings.TrimSpace(id) == "" {
		return "", nil
	}
	if err := wikiSingleID(key, id); err != nil {
		return "", err
	}
	return id, nil
}

func wikiSingleID(key, id string) error {
	if strings.ContainsAny(id, ",/\\?#% \t\r\n") || id == "." || id == ".." {
		return fmt.Errorf("%s must be a single ID", key)
	}
	return nil
}

func wikiIntArg(args map[string]any, key string) int {
	v, ok := args[key]
	if !ok || v == nil {
		return 0
	}
	switch n := v.(type) {
	case float64:
		return int(n)
	case int:
		return n
	default:
		return 0
	}
}

func wikiOptionalString(args map[string]any, key string) (string, error) {
	v, ok := args[key]
	if !ok || v == nil {
		return "", nil
	}
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("%s must be a string", key)
	}
	return s, nil
}

func wikiAPIResult(err error, failPrefix string, present bool, raw, msg string, successful bool) (*mcp.CallToolResult, error) {
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("%s: %v", failPrefix, err)), nil
	}
	if !present {
		return mcp.NewToolResultError("Dooray returned an empty response"), nil
	}
	if !successful {
		return mcp.NewToolResultError(fmt.Sprintf("Dooray rejected the request: %s", msg)), nil
	}
	return mcp.NewToolResultText(raw), nil
}
