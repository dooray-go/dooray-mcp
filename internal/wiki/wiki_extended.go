package wiki

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"

	wikimodel "github.com/dooray-go/dooray-sdk/openapi/model/wiki"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func wikiCommentsTool(s *server.MCPServer, token *string, call func(context.Context, string, string, string, int, int) (*wikimodel.ListCommentsResponse, error)) {
	tool := mcp.NewTool("dooray_wiki_comments", mcp.WithDescription("List wiki page comments."), mcp.WithString("operation", mcp.Required(), mcp.Enum("find_comments")),
		mcp.WithString("wikiId", mcp.Required()),
		mcp.WithString("pageId", mcp.Required()),
		mcp.WithNumber("page", mcp.Description("Non-negative page number; default 0")), mcp.WithNumber("size", mcp.Description("Non-negative page size; 0 uses the API default")),
	)
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		if err := wikiRequireOperation(args, "find_comments"); err != nil {
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
		page, err := wikiNonnegativeInt(args, "page")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		size, err := wikiNonnegativeInt(args, "size")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		if err := wikiRequireToken(token); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		res, err := call(ctx, *token, wikiID, pageID, page, size)
		if res == nil {
			return wikiAPIResult(err, "failed to find comments", false, "", "", false)
		}
		return wikiAPIResult(err, "failed to find comments", true, res.RawJSON, res.Header.ResultMessage, res.Header.IsSuccessful)
	})
}

func wikiCommentTool(s *server.MCPServer, token *string, call func(context.Context, string, string, string, string) (*wikimodel.GetCommentResponse, error)) {
	tool := mcp.NewTool("dooray_wiki_comment", mcp.WithDescription("Get a wiki page comment."), mcp.WithString("operation", mcp.Required(), mcp.Enum("get_comment")),
		mcp.WithString("wikiId", mcp.Required()),
		mcp.WithString("pageId", mcp.Required()),
		mcp.WithString("commentId", mcp.Required()),
	)
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		if err := wikiRequireOperation(args, "get_comment"); err != nil {
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
		commentID, err := wikiRequiredID(args, "commentId")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		if err := wikiRequireToken(token); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		res, err := call(ctx, *token, wikiID, pageID, commentID)
		if res == nil {
			return wikiAPIResult(err, "failed to get comment", false, "", "", false)
		}
		return wikiAPIResult(err, "failed to get comment", true, res.RawJSON, res.Header.ResultMessage, res.Header.IsSuccessful)
	})
}

func wikiCommentPostTool(s *server.MCPServer, token *string, call func(context.Context, string, string, string, string) (*wikimodel.IDResponse, error)) {
	tool := mcp.NewTool("dooray_wiki_comment_post", mcp.WithDescription("Create a wiki page comment. Every call creates a new comment."), mcp.WithString("operation", mcp.Required(), mcp.Enum("create_comment")),
		mcp.WithString("wikiId", mcp.Required()),
		mcp.WithString("pageId", mcp.Required()),
		mcp.WithString("content", mcp.Required()),
	)
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		if err := wikiRequireOperation(args, "create_comment"); err != nil {
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
		content, err := wikiRequiredText(args, "content")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		if err := wikiRequireToken(token); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		res, err := call(ctx, *token, wikiID, pageID, content)
		if res == nil {
			return wikiAPIResult(err, "failed to create comment", false, "", "", false)
		}
		return wikiAPIResult(err, "failed to create comment", true, res.RawJSON, res.Header.ResultMessage, res.Header.IsSuccessful)
	})
}

func wikiCommentUpdateTool(s *server.MCPServer, token *string, call func(context.Context, string, string, string, string, string) (*wikimodel.HeaderOnlyResponse, error)) {
	tool := mcp.NewTool("dooray_wiki_comment_update", mcp.WithDescription("Update a wiki page comment."), mcp.WithString("operation", mcp.Required(), mcp.Enum("update_comment")),
		mcp.WithString("wikiId", mcp.Required()),
		mcp.WithString("pageId", mcp.Required()),
		mcp.WithString("commentId", mcp.Required()),
		mcp.WithString("content", mcp.Required()),
	)
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		if err := wikiRequireOperation(args, "update_comment"); err != nil {
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
		commentID, err := wikiRequiredID(args, "commentId")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		content, err := wikiRequiredText(args, "content")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		if err := wikiRequireToken(token); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		res, err := call(ctx, *token, wikiID, pageID, commentID, content)
		if res == nil {
			return wikiAPIResult(err, "failed to update comment", false, "", "", false)
		}
		return wikiAPIResult(err, "failed to update comment", true, res.RawJSON, res.Header.ResultMessage, res.Header.IsSuccessful)
	})
}

func wikiCommentDeleteTool(s *server.MCPServer, token *string, call func(context.Context, string, string, string, string) (*wikimodel.HeaderOnlyResponse, error)) {
	tool := mcp.NewTool("dooray_wiki_comment_delete", mcp.WithDescription("Delete a wiki page comment."), mcp.WithString("operation", mcp.Required(), mcp.Enum("delete_comment")),
		mcp.WithString("wikiId", mcp.Required()),
		mcp.WithString("pageId", mcp.Required()),
		mcp.WithString("commentId", mcp.Required()),
	)
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		if err := wikiRequireOperation(args, "delete_comment"); err != nil {
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
		commentID, err := wikiRequiredID(args, "commentId")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		if err := wikiRequireToken(token); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		res, err := call(ctx, *token, wikiID, pageID, commentID)
		if res == nil {
			return wikiAPIResult(err, "failed to delete comment", false, "", "", false)
		}
		return wikiAPIResult(err, "failed to delete comment", true, res.RawJSON, res.Header.ResultMessage, res.Header.IsSuccessful)
	})
}

func wikiSharedLinksTool(s *server.MCPServer, token *string, call func(context.Context, string, string, string, int, int, *bool) (*wikimodel.ListSharedLinksResponse, error)) {
	tool := mcp.NewTool("dooray_wiki_shared_links", mcp.WithDescription("List wiki page shared links; omit valid for all links."), mcp.WithString("operation", mcp.Required(), mcp.Enum("find_shared_links")),
		mcp.WithString("wikiId", mcp.Required()),
		mcp.WithString("pageId", mcp.Required()),
		mcp.WithNumber("page", mcp.Description("Non-negative page number; default 0")), mcp.WithNumber("size", mcp.Description("Non-negative page size; 0 uses the API default")),
		mcp.WithBoolean("valid"),
	)
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		if err := wikiRequireOperation(args, "find_shared_links"); err != nil {
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
		page, err := wikiNonnegativeInt(args, "page")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		size, err := wikiNonnegativeInt(args, "size")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		valid, err := wikiOptionalBool(args, "valid")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		if err := wikiRequireToken(token); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		res, err := call(ctx, *token, wikiID, pageID, page, size, valid)
		if res == nil {
			return wikiAPIResult(err, "failed to find shared links", false, "", "", false)
		}
		return wikiAPIResult(err, "failed to find shared links", true, res.RawJSON, res.Header.ResultMessage, res.Header.IsSuccessful)
	})
}

func wikiPageDeleteTool(s *server.MCPServer, token *string, call func(context.Context, string, string, string) (*wikimodel.HeaderOnlyResponse, error)) {
	tool := mcp.NewTool("dooray_wiki_page_delete", mcp.WithDescription("Delete a wiki page."), mcp.WithString("operation", mcp.Required(), mcp.Enum("delete_page")),
		mcp.WithString("wikiId", mcp.Required()),
		mcp.WithString("pageId", mcp.Required()),
	)
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		if err := wikiRequireOperation(args, "delete_page"); err != nil {
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
		if err := wikiRequireToken(token); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		res, err := call(ctx, *token, wikiID, pageID)
		if res == nil {
			return wikiAPIResult(err, "failed to delete page", false, "", "", false)
		}
		return wikiAPIResult(err, "failed to delete page", true, res.RawJSON, res.Header.ResultMessage, res.Header.IsSuccessful)
	})
}

func wikiPageMoveTool(s *server.MCPServer, token *string, call func(context.Context, string, string, string, wikimodel.MovePageRequest) (*wikimodel.MovePageResponse, error)) {
	tool := mcp.NewTool("dooray_wiki_page_move", mcp.WithDescription("Move a wiki page. targetParentPageId identifies the destination parent."), mcp.WithString("operation", mcp.Required(), mcp.Enum("move_page")),
		mcp.WithString("wikiId", mcp.Required()),
		mcp.WithString("pageId", mcp.Required()),
		mcp.WithString("targetParentPageId", mcp.Required()),
		mcp.WithString("targetWikiId"), mcp.WithString("beforePageId"), mcp.WithBoolean("withChildren"),
	)
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		if err := wikiRequireOperation(args, "move_page"); err != nil {
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
		targetParentPageID, err := wikiRequiredID(args, "targetParentPageId")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		targetWikiID, err := wikiOptionalID(args, "targetWikiId")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		beforePageID, err := wikiOptionalID(args, "beforePageId")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		withChildren, err := wikiOptionalBool(args, "withChildren")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		req := wikimodel.MovePageRequest{TargetWikiID: targetWikiID, TargetParentPageID: targetParentPageID, BeforePageID: beforePageID, WithChildren: withChildren}
		if err := wikiRequireToken(token); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		res, err := call(ctx, *token, wikiID, pageID, req)
		if res == nil {
			return wikiAPIResult(err, "failed to move page", false, "", "", false)
		}
		return wikiAPIResult(err, "failed to move page", true, res.RawJSON, res.Header.ResultMessage, res.Header.IsSuccessful)
	})
}

func wikiPageReferrersTool(s *server.MCPServer, token *string, call func(context.Context, string, string, string, []wikimodel.Actor) (*wikimodel.HeaderOnlyResponse, error)) {
	tool := mcp.NewTool("dooray_wiki_page_referrers_update", mcp.WithDescription("Replace wiki page referrers. An empty referrerMemberIds array clears them."), mcp.WithString("operation", mcp.Required(), mcp.Enum("update_referrers")),
		mcp.WithString("wikiId", mcp.Required()),
		mcp.WithString("pageId", mcp.Required()),
		mcp.WithArray("referrerMemberIds", mcp.Required(), mcp.WithStringItems()),
	)
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		if err := wikiRequireOperation(args, "update_referrers"); err != nil {
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
		referrers, err := wikiReferrers(args, true)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		if err := wikiRequireToken(token); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		res, err := call(ctx, *token, wikiID, pageID, referrers)
		if res == nil {
			return wikiAPIResult(err, "failed to update referrers", false, "", "", false)
		}
		return wikiAPIResult(err, "failed to update referrers", true, res.RawJSON, res.Header.ResultMessage, res.Header.IsSuccessful)
	})
}

func wikiFileUploadTool(s *server.MCPServer, token *string, call func(context.Context, string, string, string, string, []byte) (*wikimodel.UploadFileResponse, error)) {
	tool := mcp.NewTool("dooray_wiki_file_upload", mcp.WithDescription("Upload a wiki attachment before creating a page; pass the returned attachment ID in attachFileIds."), mcp.WithString("operation", mcp.Required(), mcp.Enum("upload_wiki_file")),
		mcp.WithString("wikiId", mcp.Required()),
		mcp.WithString("filename", mcp.Required()), mcp.WithString("contentBase64", mcp.Required(), mcp.Description("Standard base64 encoded file bytes")), mcp.WithString("fileType", mcp.Enum("general", "inline_image"), mcp.Description("Default general")),
	)
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		if err := wikiRequireOperation(args, "upload_wiki_file"); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		wikiID, err := wikiRequiredID(args, "wikiId")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		fileType, filename, content, err := wikiUploadArgs(args)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		if err := wikiRequireToken(token); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		res, err := call(ctx, *token, wikiID, fileType, filename, content)
		if res == nil {
			return wikiAPIResult(err, "failed to upload wiki file", false, "", "", false)
		}
		return wikiAPIResult(err, "failed to upload wiki file", true, res.RawJSON, res.Header.ResultMessage, res.Header.IsSuccessful)
	})
}

func wikiPageFileUploadTool(s *server.MCPServer, token *string, call func(context.Context, string, string, string, string, string, []byte) (*wikimodel.UploadFileResponse, error)) {
	tool := mcp.NewTool("dooray_wiki_page_file_upload", mcp.WithDescription("Upload a file or inline image to an existing wiki page."), mcp.WithString("operation", mcp.Required(), mcp.Enum("upload_page_file")),
		mcp.WithString("wikiId", mcp.Required()),
		mcp.WithString("pageId", mcp.Required()),
		mcp.WithString("filename", mcp.Required()), mcp.WithString("contentBase64", mcp.Required(), mcp.Description("Standard base64 encoded file bytes")), mcp.WithString("fileType", mcp.Enum("general", "inline_image"), mcp.Description("Default general")),
	)
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		if err := wikiRequireOperation(args, "upload_page_file"); err != nil {
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
		fileType, filename, content, err := wikiUploadArgs(args)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		if err := wikiRequireToken(token); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		res, err := call(ctx, *token, wikiID, pageID, fileType, filename, content)
		if res == nil {
			return wikiAPIResult(err, "failed to upload page file", false, "", "", false)
		}
		return wikiAPIResult(err, "failed to upload page file", true, res.RawJSON, res.Header.ResultMessage, res.Header.IsSuccessful)
	})
}

func wikiAttachFileDownloadTool(s *server.MCPServer, token *string, call func(context.Context, string, string, string) (*wikimodel.Download, error)) {
	tool := mcp.NewTool("dooray_wiki_attach_file_download", mcp.WithDescription("Download a wiki attachment as JSON containing contentBase64, contentType and statusCode."), mcp.WithString("operation", mcp.Required(), mcp.Enum("download_attach_file")),
		mcp.WithString("wikiId", mcp.Required()),
		mcp.WithString("attachFileId", mcp.Required()),
	)
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		if err := wikiRequireOperation(args, "download_attach_file"); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		wikiID, err := wikiRequiredID(args, "wikiId")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		attachFileID, err := wikiRequiredID(args, "attachFileId")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		if err := wikiRequireToken(token); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		res, err := call(ctx, *token, wikiID, attachFileID)
		return wikiDownloadResult(res, err)
	})
}

func wikiPageFileDownloadTool(s *server.MCPServer, token *string, call func(context.Context, string, string, string, string) (*wikimodel.Download, error)) {
	tool := mcp.NewTool("dooray_wiki_page_file_download", mcp.WithDescription("Download a wiki page file as JSON containing contentBase64, contentType and statusCode."), mcp.WithString("operation", mcp.Required(), mcp.Enum("download_page_file")),
		mcp.WithString("wikiId", mcp.Required()),
		mcp.WithString("pageId", mcp.Required()),
		mcp.WithString("fileId", mcp.Required()),
	)
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		if err := wikiRequireOperation(args, "download_page_file"); err != nil {
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
		fileID, err := wikiRequiredID(args, "fileId")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		if err := wikiRequireToken(token); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		res, err := call(ctx, *token, wikiID, pageID, fileID)
		return wikiDownloadResult(res, err)
	})
}

func wikiPageFileDeleteTool(s *server.MCPServer, token *string, call func(context.Context, string, string, string, string) (*wikimodel.HeaderOnlyResponse, error)) {
	tool := mcp.NewTool("dooray_wiki_page_file_delete", mcp.WithDescription("Delete an attached wiki page file."), mcp.WithString("operation", mcp.Required(), mcp.Enum("delete_page_file")),
		mcp.WithString("wikiId", mcp.Required()),
		mcp.WithString("pageId", mcp.Required()),
		mcp.WithString("fileId", mcp.Required()),
	)
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		if err := wikiRequireOperation(args, "delete_page_file"); err != nil {
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
		fileID, err := wikiRequiredID(args, "fileId")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		if err := wikiRequireToken(token); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		res, err := call(ctx, *token, wikiID, pageID, fileID)
		if res == nil {
			return wikiAPIResult(err, "failed to delete page file", false, "", "", false)
		}
		return wikiAPIResult(err, "failed to delete page file", true, res.RawJSON, res.Header.ResultMessage, res.Header.IsSuccessful)
	})
}

func wikiNonnegativeInt(args map[string]any, key string) (int, error) {
	v, ok := args[key]
	if !ok {
		return 0, nil
	}
	switch n := v.(type) {
	case int:
		if n >= 0 {
			return n, nil
		}
	case float64:
		if n >= 0 && n < float64(math.MaxInt) && math.Trunc(n) == n {
			return int(n), nil
		}
	}
	return 0, fmt.Errorf("%s must be a non-negative integer", key)
}

func wikiOptionalBool(args map[string]any, key string) (*bool, error) {
	v, ok := args[key]
	if !ok {
		return nil, nil
	}
	b, ok := v.(bool)
	if !ok {
		return nil, fmt.Errorf("%s must be a boolean", key)
	}
	return new(b), nil
}

func wikiIDArray(args map[string]any, key string, required bool) ([]string, error) {
	v, ok := args[key]
	if !ok && !required {
		return nil, nil
	}
	var values []string
	switch a := v.(type) {
	case []string:
		values = a
	case []any:
		values = make([]string, 0, len(a))
		for _, item := range a {
			s, ok := item.(string)
			if !ok {
				return nil, fmt.Errorf("%s must be an array of IDs", key)
			}
			values = append(values, s)
		}
	default:
		return nil, fmt.Errorf("%s must be an array of IDs", key)
	}
	for _, id := range values {
		if id == "" {
			return nil, fmt.Errorf("%s must contain non-empty IDs", key)
		}
		if err := wikiSingleID(key, id); err != nil {
			return nil, err
		}
	}
	return values, nil
}

func wikiReferrers(args map[string]any, required bool) ([]wikimodel.Actor, error) {
	ids, err := wikiIDArray(args, "referrerMemberIds", required)
	if err != nil {
		return nil, err
	}
	actors := make([]wikimodel.Actor, 0, len(ids))
	for _, id := range ids {
		actors = append(actors, wikimodel.Actor{Type: "member", Member: wikimodel.Member{OrganizationMemberID: id}})
	}
	return actors, nil
}

func wikiUploadArgs(args map[string]any) (string, string, []byte, error) {
	filename, err := wikiRequiredText(args, "filename")
	if err != nil {
		return "", "", nil, err
	}
	encoded, ok := args["contentBase64"].(string)
	if !ok {
		return "", "", nil, fmt.Errorf("contentBase64 must be a base64 string")
	}
	content, err := base64.StdEncoding.Strict().DecodeString(encoded)
	if err != nil {
		return "", "", nil, fmt.Errorf("contentBase64 must be valid standard base64")
	}
	fileType, err := wikiOptionalString(args, "fileType")
	if err != nil {
		return "", "", nil, err
	}
	if fileType == "" {
		fileType = "general"
	}
	if fileType != "general" && fileType != "inline_image" {
		return "", "", nil, fmt.Errorf("fileType must be general or inline_image")
	}
	return fileType, filename, content, nil
}

func wikiDownloadResult(res *wikimodel.Download, err error) (*mcp.CallToolResult, error) {
	if err != nil || res == nil {
		return wikiAPIResult(err, "failed to download wiki file", false, "", "", false)
	}
	if res.StatusCode != 200 {
		return mcp.NewToolResultError(fmt.Sprintf("wiki file download returned status %d", res.StatusCode)), nil
	}
	raw, err := json.Marshal(struct {
		ContentBase64 string `json:"contentBase64"`
		ContentType   string `json:"contentType"`
		StatusCode    int    `json:"statusCode"`
	}{base64.StdEncoding.EncodeToString(res.Content), res.ContentType, res.StatusCode})
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(string(raw)), nil
}
