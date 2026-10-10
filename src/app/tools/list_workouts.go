package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/0xalexb/intervals-icu-mcp/src/app/clients/intervals"
)

type listWorkoutsArgs struct{}

// NewListWorkoutsTool returns a ToolRegistration that registers the list_workouts tool on an MCP server.
func NewListWorkoutsTool(apiClient *intervals.Client) ToolRegistration {
	return func(server *mcp.Server) {
		mcp.AddTool(server,
			&mcp.Tool{
				Name:        "list_workouts",
				Description: "Lists the structured workouts in the athlete's workout library.",
			},
			func(ctx context.Context, _ *mcp.CallToolRequest, _ listWorkoutsArgs) (*mcp.CallToolResult, any, error) {
				body, err := apiClient.Get(ctx, "/api/v1/athlete/"+apiClient.AthleteID()+"/workouts", nil)
				if err != nil {
					return nil, nil, fmt.Errorf("listing workouts: %w", err)
				}

				return textResult(body), nil, nil
			},
		)
	}
}
