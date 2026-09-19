package main

import "testing"

func TestAdvertisedMCPVersion(t *testing.T) {
	if mcpVersion == "1.0.0" {
		t.Fatal("MCP initialize version is still 1.0.0; bump it with the release")
	}
	if mcpVersion != "1.5.0" {
		t.Errorf("mcpVersion = %q, want 1.5.0", mcpVersion)
	}
}
