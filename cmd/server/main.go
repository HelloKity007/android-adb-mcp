package main

import (
	"context"
	"fmt"
	"os"

	"github.com/HelloKity007/android-adb-mcp/internal/adb"
	"github.com/HelloKity007/android-adb-mcp/internal/mcp"
	"github.com/HelloKity007/android-adb-mcp/internal/mcp/transport"
)

func main() {
	// Create ADB client
	client, err := adb.NewClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create ADB client: %v\n", err)
		os.Exit(1)
	}

	// Create MCP server
	server := mcp.NewServer("android-adb-mcp", "1.0.0")

	// Register all tools
	mcp.RegisterADBTools(server, client)

	// Start server with stdio transport
	fmt.Fprintln(os.Stderr, "Android ADB MCP Server starting...")
	t := &transport.StdioTransport{}
	if err := t.Serve(context.Background(), server); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
