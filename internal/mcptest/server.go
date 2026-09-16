package mcptest

import "github.com/mark3labs/mcp-go/server"

func NewServer() *server.MCPServer {
	return server.NewMCPServer("dooray-test", "1.0.0", server.WithToolCapabilities(true))
}
