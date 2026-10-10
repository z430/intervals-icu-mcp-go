package tools

import "github.com/modelcontextprotocol/go-sdk/mcp"

func textResult(body []byte) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: string(body),
			},
		},
	}
}
