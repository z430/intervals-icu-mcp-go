package tools

import (
	"context"
	"fmt"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/0xalexb/intervals-icu-mcp/src/app/clients/intervals"
)

type getActivityIntervalsArgs struct {
	ActivityID string `json:"activity_id" jsonschema:"the activity ID"`
}

// NewGetActivityIntervalsTool returns a ToolRegistration that registers the get_activity_intervals tool on an MCP server.
func NewGetActivityIntervalsTool(apiClient *intervals.Client) ToolRegistration {
	return func(server *mcp.Server) {
		mcp.AddTool(server,
			&mcp.Tool{
				Name:        "get_activity_intervals",
				Description: "Returns the detected and manual intervals (laps, efforts) of an activity with their metrics.",
			},
			func(ctx context.Context, _ *mcp.CallToolRequest, args getActivityIntervalsArgs) (*mcp.CallToolResult, any, error) {
				if args.ActivityID == "" {
					return nil, nil, errMissingActivityID
				}

				body, err := apiClient.Get(ctx, "/api/v1/activity/"+url.PathEscape(args.ActivityID)+"/intervals", nil)
				if err != nil {
					return nil, nil, fmt.Errorf("fetching activity intervals: %w", err)
				}

				return textResult(body), nil, nil
			},
		)
	}
}
