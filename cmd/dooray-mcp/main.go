package main

import (
	"flag"
	"fmt"

	"github.com/mark3labs/mcp-go/server"

	"dooray_mcp/internal/account"
	"dooray_mcp/internal/calendar"
	"dooray_mcp/internal/messenger"
	"dooray_mcp/internal/ostool"
	"dooray_mcp/internal/project"
	"dooray_mcp/internal/wiki"
)

const mcpVersion = "1.5.0-beta.2"

func main() {
	token := flag.String("token", "", "개인설정 > API > 개인 인증 토큰 메뉴에서 생성할 수 있습니다.")
	flag.Parse()

	if *token == "" {
		fmt.Printf("token must be set!!")
		return
	}

	s := server.NewMCPServer(
		"dooray",
		mcpVersion,
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)

	ostool.Tools(s, token)
	messenger.Tools(s, token)
	account.Tools(s, token)
	calendar.Tools(s, token)
	project.Tools(s, token)
	wiki.Tools(s, token)

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
