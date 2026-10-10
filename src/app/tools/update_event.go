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

type updateEventArgs struct {
	EventID        string   `json:"event_id"                   jsonschema:"the event ID to update"`
	Name           string   `json:"name,omitempty"             jsonschema:"event name"`
	Description    string   `json:"description,omitempty"      jsonschema:"event description or notes"`
	StartDateLocal string   `json:"start_date_local,omitempty" jsonschema:"start date as yyyy-MM-dd or yyyy-MM-ddTHH:mm:ss (local time); a bare date means midnight"`
	Category       string   `json:"category,omitempty"         jsonschema:"event category: WORKOUT, RACE_A, RACE_B, RACE_C, NOTE, PLAN, HOLIDAY, SICK, INJURED, SET_EFTP, FITNESS_DAYS, SEASON_START, TARGET or SET_FITNESS"`
	Type           string   `json:"type,omitempty"             jsonschema:"sport type (e.g. Ride or Run or Swim)"`
	MovingTime     *float64 `json:"moving_time,omitempty"      jsonschema:"planned moving time in seconds"`
	Distance       *float64 `json:"distance,omitempty"         jsonschema:"planned distance in meters"`
	TrainingLoad   *float64 `json:"training_load,omitempty"    jsonschema:"planned training load (TSS or similar)"`
}

// NewUpdateEventTool returns a ToolRegistration that registers the update_event tool on an MCP server.
func NewUpdateEventTool(apiClient *intervals.Client) ToolRegistration {
	return func(server *mcp.Server) {
		mcp.AddTool(server,
			&mcp.Tool{
				Name:        "update_event",
				Description: "Updates an existing event for the athlete. Only the provided fields will be updated.",
			},
			func(ctx context.Context, _ *mcp.CallToolRequest, args updateEventArgs) (*mcp.CallToolResult, any, error) {
				if args.EventID == "" {
					return nil, nil, errMissingEventID
				}

				payload := eventPayload{
					Name:         args.Name,
					Description:  args.Description,
					Type:         args.Type,
					MovingTime:   roundedInt(args.MovingTime),
					Distance:     args.Distance,
					TrainingLoad: roundedInt(args.TrainingLoad),
				}

				if args.StartDateLocal != "" {
					startDate, err := normalizeLocalDateTime(args.StartDateLocal)
					if err != nil {
						return nil, nil, err
					}

					payload.StartDate = startDate
				}

				if args.Category != "" {
					category, err := normalizeEventCategory(args.Category)
					if err != nil {
						return nil, nil, err
					}

					payload.Category = category
				}

				payloadBytes, err := json.Marshal(payload)
				if err != nil {
					return nil, nil, fmt.Errorf("marshaling event body: %w", err)
				}

				path := "/api/v1/athlete/" + apiClient.AthleteID() + "/events/" + url.PathEscape(args.EventID)

				body, err := apiClient.Put(ctx, path, bytes.NewReader(payloadBytes))
				if err != nil {
					return nil, nil, fmt.Errorf("updating event: %w", err)
				}

				return &mcp.CallToolResult{
					Content: []mcp.Content{
						&mcp.TextContent{
							Text: string(body),
						},
					},
				}, nil, nil
			},
		)
	}
}
