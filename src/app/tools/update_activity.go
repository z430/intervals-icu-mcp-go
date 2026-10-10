package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/0xalexb/intervals-icu-mcp/src/app/clients/intervals"
)

type updateActivityArgs struct {
	ActivityID  string `json:"activity_id"           jsonschema:"the activity ID to update"`
	Name        string `json:"name,omitempty"        jsonschema:"activity name"`
	Description string `json:"description,omitempty" jsonschema:"activity description or notes"`
	Type        string `json:"type,omitempty"        jsonschema:"sport type (e.g. Ride or Run or Swim)"`
	Feel        *int   `json:"feel,omitempty"        jsonschema:"how it felt, 1 (great) to 5 (poor)"`
	RPE         *int   `json:"icu_rpe,omitempty"     jsonschema:"rate of perceived exertion 1-10"`
	Commute     *bool  `json:"commute,omitempty"     jsonschema:"mark as a commute"`
	Trainer     *bool  `json:"trainer,omitempty"     jsonschema:"mark as indoor trainer"`
	Race        *bool  `json:"race,omitempty"        jsonschema:"mark as a race"`
}

type updateActivityPayload struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type,omitempty"`
	Feel        *int   `json:"feel,omitempty"`
	RPE         *int   `json:"icu_rpe,omitempty"`
	Commute     *bool  `json:"commute,omitempty"`
	Trainer     *bool  `json:"trainer,omitempty"`
	Race        *bool  `json:"race,omitempty"`
}

// NewUpdateActivityTool returns a ToolRegistration that registers the update_activity tool on an MCP server.
func NewUpdateActivityTool(apiClient *intervals.Client) ToolRegistration {
	return func(server *mcp.Server) {
		mcp.AddTool(server,
			&mcp.Tool{
				Name:        "update_activity",
				Description: "Updates an existing activity (name, description, type, feel, RPE, flags). Only the provided fields are changed.",
			},
			func(ctx context.Context, _ *mcp.CallToolRequest, args updateActivityArgs) (*mcp.CallToolResult, any, error) {
				if args.ActivityID == "" {
					return nil, nil, errMissingActivityID
				}

				payloadBytes, err := json.Marshal(updateActivityPayload{
					Name:        args.Name,
					Description: args.Description,
					Type:        args.Type,
					Feel:        args.Feel,
					RPE:         args.RPE,
					Commute:     args.Commute,
					Trainer:     args.Trainer,
					Race:        args.Race,
				})
				if err != nil {
					return nil, nil, fmt.Errorf("marshaling activity body: %w", err)
				}

				body, err := apiClient.Put(ctx, "/api/v1/activity/"+url.PathEscape(args.ActivityID), bytes.NewReader(payloadBytes))
				if err != nil {
					return nil, nil, fmt.Errorf("updating activity: %w", err)
				}

				return textResult(body), nil, nil
			},
		)
	}
}
