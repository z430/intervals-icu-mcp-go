package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/0xalexb/intervals-icu-mcp/src/app/clients/intervals"
)

type createEventArgs struct {
	Name           string   `json:"name"                    jsonschema:"event name"`
	StartDateLocal string   `json:"start_date_local"        jsonschema:"start date as yyyy-MM-dd or yyyy-MM-ddTHH:mm:ss (local time); a bare date means midnight"`
	Category       string   `json:"category"                jsonschema:"event category: WORKOUT, RACE_A, RACE_B, RACE_C, NOTE, PLAN, HOLIDAY, SICK, INJURED, SET_EFTP, FITNESS_DAYS, SEASON_START, TARGET or SET_FITNESS"`
	Type           string   `json:"type,omitempty"          jsonschema:"sport type (e.g. Ride or Run or Swim); required for WORKOUT and RACE_* events"`
	Description    string   `json:"description,omitempty"   jsonschema:"event description or notes"`
	MovingTime     *float64 `json:"moving_time,omitempty"   jsonschema:"planned moving time in seconds"`
	Distance       *float64 `json:"distance,omitempty"      jsonschema:"planned distance in meters"`
	TrainingLoad   *float64 `json:"training_load,omitempty" jsonschema:"planned training load (TSS or similar)"`
}

// eventPayload is the request body the Intervals.icu API accepts for events.
type eventPayload struct {
	Name         string   `json:"name,omitempty"`
	StartDate    string   `json:"start_date_local,omitempty"`
	Category     string   `json:"category,omitempty"`
	Type         string   `json:"type,omitempty"`
	Description  string   `json:"description,omitempty"`
	MovingTime   *int     `json:"moving_time,omitempty"`
	Distance     *float64 `json:"distance,omitempty"`
	TrainingLoad *int     `json:"icu_training_load,omitempty"`
}

func roundedInt(v *float64) *int {
	if v == nil {
		return nil
	}

	rounded := int(math.Round(*v))

	return &rounded
}

// NewCreateEventTool returns a ToolRegistration that registers the create_event tool on an MCP server.
func NewCreateEventTool(apiClient *intervals.Client) ToolRegistration {
	return func(server *mcp.Server) {
		mcp.AddTool(server,
			&mcp.Tool{
				Name:        "create_event",
				Description: "Creates a new event (workout, race, note, etc.) for the athlete on the specified date.",
			},
			func(ctx context.Context, _ *mcp.CallToolRequest, args createEventArgs) (*mcp.CallToolResult, any, error) {
				if args.Name == "" {
					return nil, nil, errMissingName
				}

				if args.StartDateLocal == "" {
					return nil, nil, errMissingStartDate
				}

				if args.Category == "" {
					return nil, nil, errMissingCategory
				}

				startDate, err := normalizeLocalDateTime(args.StartDateLocal)
				if err != nil {
					return nil, nil, err
				}

				category, err := normalizeEventCategory(args.Category)
				if err != nil {
					return nil, nil, err
				}

				if args.Type == "" && strings.HasPrefix(category, "RACE_") {
					return nil, nil, fmt.Errorf("%w: type (sport) is required for category %s", errRequiredField, category)
				}

				payload, err := json.Marshal(eventPayload{
					Name:         args.Name,
					StartDate:    startDate,
					Category:     category,
					Type:         args.Type,
					Description:  args.Description,
					MovingTime:   roundedInt(args.MovingTime),
					Distance:     args.Distance,
					TrainingLoad: roundedInt(args.TrainingLoad),
				})
				if err != nil {
					return nil, nil, fmt.Errorf("marshaling event body: %w", err)
				}

				path := "/api/v1/athlete/" + apiClient.AthleteID() + "/events"

				body, err := apiClient.Post(ctx, path, bytes.NewReader(payload))
				if err != nil {
					return nil, nil, fmt.Errorf("creating event: %w", err)
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
