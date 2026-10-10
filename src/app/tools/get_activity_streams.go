package tools

import (
	"context"
	"fmt"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/0xalexb/intervals-icu-mcp/src/app/clients/intervals"
)

type getActivityStreamsArgs struct {
	ActivityID string `json:"activity_id" jsonschema:"the activity ID"`
	Types      string `json:"types"       jsonschema:"comma-separated stream types to return (e.g. watts or heartrate or cadence or velocity_smooth or altitude or time)"`
}

// NewGetActivityStreamsTool returns a ToolRegistration that registers the get_activity_streams tool on an MCP server.
func NewGetActivityStreamsTool(apiClient *intervals.Client) ToolRegistration {
	return func(server *mcp.Server) {
		mcp.AddTool(server,
			&mcp.Tool{
				Name:        "get_activity_streams",
				Description: "Returns second-by-second data streams (power, heart rate, cadence, speed, altitude, ...) for an activity. Always request specific types; full streams can be very large.",
			},
			func(ctx context.Context, _ *mcp.CallToolRequest, args getActivityStreamsArgs) (*mcp.CallToolResult, any, error) {
				if args.ActivityID == "" {
					return nil, nil, errMissingActivityID
				}

				if args.Types == "" {
					return nil, nil, errMissingTypes
				}

				query := url.Values{}
				query.Set("types", args.Types)

				body, err := apiClient.Get(ctx, "/api/v1/activity/"+url.PathEscape(args.ActivityID)+"/streams.json", query)
				if err != nil {
					return nil, nil, fmt.Errorf("fetching activity streams: %w", err)
				}

				return textResult(body), nil, nil
			},
		)
	}
}
