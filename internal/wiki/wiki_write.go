package wiki

import (
	"context"
	"fmt"
	"strings"

	wikimodel "github.com/dooray-go/dooray-sdk/openapi/model/wiki"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func wikiCreatePageTool(s *server.MCPServer, token *string, create func(context.Context, string, string, wikimodel.CreatePageRequest) (*wikimodel.CreatePageResponse, error)) {
	tool := mcp.NewTool("dooray_wiki_page_post",
		mcp.WithDescription("Create a Dooray wiki page. Creates a new page on every call. Use dooray_wikis and dooray_wiki_pages to choose wikiId and parentPageId."),
		mcp.WithString("operation", mcp.Required(), mcp.Enum("create_page")),
		mcp.WithString("wikiId", mcp.Required(), mcp.Description("Wiki ID from dooray_wikis")),
		mcp.WithString("subject", mcp.Required(), mcp.Description("Page title")),
		mcp.WithString("content", mcp.Required(), mcp.Description("Page body")),
		mcp.WithString("parentPageId", mcp.Description("Parent page ID. Omit to create under the wiki root")),
		mcp.WithArray("attachFileIds", mcp.WithStringItems(), mcp.Description("Attachment IDs returned by dooray_wiki_file_upload")),
		mcp.WithArray("referrerMemberIds", mcp.WithStringItems(), mcp.Description("Organization member IDs to add as referrers")),
		mcp.WithString("mimeType", mcp.Enum("text/x-markdown", "text/html"), mcp.Description("Body format, default text/x-markdown")),
	)
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		if err := wikiRequireOperation(args, "create_page"); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		wikiID, err := wikiRequiredID(args, "wikiId")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		subject, err := wikiRequiredText(args, "subject")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		content, err := wikiRequiredText(args, "content")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		parentPageID, err := wikiOptionalID(args, "parentPageId")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		mimeType, err := wikiOptionalString(args, "mimeType")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		if mimeType == "" {
			mimeType = "text/x-markdown"
		}
		if mimeType != "text/x-markdown" && mimeType != "text/html" {
			return mcp.NewToolResultError("mimeType must be text/x-markdown or text/html"), nil
		}
		if err := wikiRequireToken(token); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		attachFileIDs, err := wikiIDArray(args, "attachFileIds", false)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		referrers, err := wikiReferrers(args, false)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		req := wikimodel.CreatePageRequest{
			AttachFileIDs: attachFileIDs, Referrers: referrers,
			ParentPageID: parentPageID,
			Subject:      subject,
			Body:         wikimodel.Body{MimeType: mimeType, Content: content},
		}
		res, err := create(ctx, *token, wikiID, req)
		if res == nil {
			return wikiAPIResult(err, "failed to create wiki page", false, "", "", false)
		}
		return wikiAPIResult(err, "failed to create wiki page", true, res.RawJSON, res.Header.ResultMessage, res.Header.IsSuccessful)
	})
}

func wikiUpdatePageTool(s *server.MCPServer, token *string, update func(context.Context, string, string, string, wikimodel.UpdatePageRequest) (*wikimodel.HeaderOnlyResponse, error), updateTitle func(context.Context, string, string, string, string) (*wikimodel.HeaderOnlyResponse, error), updateContent func(context.Context, string, string, string, wikimodel.Body) (*wikimodel.HeaderOnlyResponse, error)) {
	tool := mcp.NewTool("dooray_wiki_page_update",
		mcp.WithDescription("Update a Dooray wiki page title and/or body. Provide subject, content, or both. Omitted fields stay unchanged."),
		mcp.WithString("operation", mcp.Required(), mcp.Enum("update_page")),
		mcp.WithString("wikiId", mcp.Required(), mcp.Description("Wiki ID")),
		mcp.WithString("pageId", mcp.Required(), mcp.Description("Page ID from dooray_wiki_pages or dooray_wiki_page")),
		mcp.WithString("subject", mcp.Description("New page title")),
		mcp.WithString("content", mcp.Description("New page body")),
		mcp.WithString("mimeType", mcp.Enum("text/x-markdown", "text/html"), mcp.Description("Body format when content is set, default text/x-markdown")),
	)
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		if err := wikiRequireOperation(args, "update_page"); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		wikiID, err := wikiRequiredID(args, "wikiId")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		pageID, err := wikiRequiredID(args, "pageId")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		subject, err := wikiOptionalString(args, "subject")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		content, err := wikiOptionalString(args, "content")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		mimeType, err := wikiOptionalString(args, "mimeType")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		hasSubject := strings.TrimSpace(subject) != ""
		hasContent := content != ""
		if !hasSubject && !hasContent {
			return mcp.NewToolResultError("provide subject, content, or both"), nil
		}
		if mimeType == "" {
			mimeType = "text/x-markdown"
		}
		if hasContent && mimeType != "text/x-markdown" && mimeType != "text/html" {
			return mcp.NewToolResultError("mimeType must be text/x-markdown or text/html"), nil
		}
		if err := wikiRequireToken(token); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		var res *wikimodel.HeaderOnlyResponse
		switch {
		case hasSubject && hasContent:
			res, err = update(ctx, *token, wikiID, pageID, wikimodel.UpdatePageRequest{
				Subject: subject,
				Body:    wikimodel.Body{MimeType: mimeType, Content: content},
			})
		case hasSubject:
			res, err = updateTitle(ctx, *token, wikiID, pageID, subject)
		default:
			res, err = updateContent(ctx, *token, wikiID, pageID, wikimodel.Body{MimeType: mimeType, Content: content})
		}
		if res == nil {
			return wikiAPIResult(err, "failed to update wiki page", false, "", "", false)
		}
		return wikiAPIResult(err, "failed to update wiki page", true, res.RawJSON, res.Header.ResultMessage, res.Header.IsSuccessful)
	})
}

func wikiRequiredText(args map[string]any, key string) (string, error) {
	v, ok := args[key]
	if !ok {
		return "", fmt.Errorf("%s must be a non-empty string", key)
	}
	s, ok := v.(string)
	if !ok || strings.TrimSpace(s) == "" {
		return "", fmt.Errorf("%s must be a non-empty string", key)
	}
	return s, nil
}
