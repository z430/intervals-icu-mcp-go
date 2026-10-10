package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/0xalexb/intervals-icu-mcp/src/app/clients/intervals"
)

type createWorkoutArgs struct {
	Name        string `json:"name"                  jsonschema:"workout name"`
	Description string `json:"description,omitempty" jsonschema:"workout steps in Intervals.icu workout syntax (e.g. '- 10m 55%' per line) plus optional notes"`
	Type        string `json:"type,omitempty"        jsonschema:"sport type (e.g. Ride or Run or Swim)"`
	FolderID    *int   `json:"folder_id,omitempty"   jsonschema:"folder ID to place the workout in"`
	Indoor      *bool  `json:"indoor,omitempty"      jsonschema:"whether this is an indoor workout"`
}

type createWorkoutPayload struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type,omitempty"`
	FolderID    *int   `json:"folder_id,omitempty"`
	Indoor      *bool  `json:"indoor,omitempty"`
}

// NewCreateWorkoutTool returns a ToolRegistration that registers the create_workout tool on an MCP server.
func NewCreateWorkoutTool(apiClient *intervals.Client) ToolRegistration {
	return func(server *mcp.Server) {
		mcp.AddTool(server,
			&mcp.Tool{
				Name:        "create_workout",
				Description: "Creates a structured workout in the athlete's workout library. The description uses Intervals.icu workout syntax.",
			},
			func(ctx context.Context, _ *mcp.CallToolRequest, args createWorkoutArgs) (*mcp.CallToolResult, any, error) {
				if args.Name == "" {
					return nil, nil, errMissingName
				}

				payloadBytes, err := json.Marshal(createWorkoutPayload{
					Name:        args.Name,
					Description: args.Description,
					Type:        args.Type,
					FolderID:    args.FolderID,
					Indoor:      args.Indoor,
				})
				if err != nil {
					return nil, nil, fmt.Errorf("marshaling workout body: %w", err)
				}

				body, err := apiClient.Post(ctx, "/api/v1/athlete/"+apiClient.AthleteID()+"/workouts", bytes.NewReader(payloadBytes))
				if err != nil {
					return nil, nil, fmt.Errorf("creating workout: %w", err)
				}

				return textResult(body), nil, nil
			},
		)
	}
}
