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

type updateWellnessArgs struct {
	Date          string   `json:"date"                    jsonschema:"date in yyyy-MM-dd format"`
	Weight        *float64 `json:"weight,omitempty"        jsonschema:"body weight in kg"`
	RestingHR     *int     `json:"restingHR,omitempty"     jsonschema:"resting heart rate in bpm"`
	HRV           *float64 `json:"hrv,omitempty"           jsonschema:"heart rate variability (rMSSD)"`
	SleepSecs     *int     `json:"sleepSecs,omitempty"     jsonschema:"sleep duration in seconds"`
	SleepQuality  *int     `json:"sleepQuality,omitempty"  jsonschema:"sleep quality 1-4 (1 great, 4 poor)"`
	AvgSleepingHR *float64 `json:"avgSleepingHR,omitempty" jsonschema:"average sleeping heart rate in bpm"`
	Soreness      *int     `json:"soreness,omitempty"      jsonschema:"soreness 1-4 (1 low, 4 extreme)"`
	Fatigue       *int     `json:"fatigue,omitempty"       jsonschema:"fatigue 1-4 (1 low, 4 extreme)"`
	Stress        *int     `json:"stress,omitempty"        jsonschema:"stress 1-4 (1 low, 4 extreme)"`
	Mood          *int     `json:"mood,omitempty"          jsonschema:"mood 1-4 (1 great, 4 poor)"`
	Motivation    *int     `json:"motivation,omitempty"    jsonschema:"motivation 1-4 (1 high, 4 low)"`
	Readiness     *float64 `json:"readiness,omitempty"     jsonschema:"readiness score"`
	SpO2          *float64 `json:"spO2,omitempty"          jsonschema:"blood oxygen saturation in percent"`
	BodyFat       *float64 `json:"bodyFat,omitempty"       jsonschema:"body fat in percent"`
	Steps         *int     `json:"steps,omitempty"         jsonschema:"step count"`
	KcalConsumed  *int     `json:"kcalConsumed,omitempty"  jsonschema:"calories consumed"`
	Comments      string   `json:"comments,omitempty"      jsonschema:"free-text comments"`
}

type updateWellnessPayload struct {
	Weight        *float64 `json:"weight,omitempty"`
	RestingHR     *int     `json:"restingHR,omitempty"`
	HRV           *float64 `json:"hrv,omitempty"`
	SleepSecs     *int     `json:"sleepSecs,omitempty"`
	SleepQuality  *int     `json:"sleepQuality,omitempty"`
	AvgSleepingHR *float64 `json:"avgSleepingHR,omitempty"`
	Soreness      *int     `json:"soreness,omitempty"`
	Fatigue       *int     `json:"fatigue,omitempty"`
	Stress        *int     `json:"stress,omitempty"`
	Mood          *int     `json:"mood,omitempty"`
	Motivation    *int     `json:"motivation,omitempty"`
	Readiness     *float64 `json:"readiness,omitempty"`
	SpO2          *float64 `json:"spO2,omitempty"`
	BodyFat       *float64 `json:"bodyFat,omitempty"`
	Steps         *int     `json:"steps,omitempty"`
	KcalConsumed  *int     `json:"kcalConsumed,omitempty"`
	Comments      string   `json:"comments,omitempty"`
}

// NewUpdateWellnessTool returns a ToolRegistration that registers the update_wellness tool on an MCP server.
func NewUpdateWellnessTool(apiClient *intervals.Client) ToolRegistration {
	return func(server *mcp.Server) {
		mcp.AddTool(server,
			&mcp.Tool{
				Name:        "update_wellness",
				Description: "Creates or updates wellness data (weight, HRV, sleep, mood, etc.) for the athlete on a date. Only the provided fields are changed.",
			},
			func(ctx context.Context, _ *mcp.CallToolRequest, args updateWellnessArgs) (*mcp.CallToolResult, any, error) {
				if args.Date == "" {
					return nil, nil, errMissingDate
				}

				err := validateDateFormat(args.Date)
				if err != nil {
					return nil, nil, err
				}

				payloadBytes, err := json.Marshal(updateWellnessPayload{
					Weight:        args.Weight,
					RestingHR:     args.RestingHR,
					HRV:           args.HRV,
					SleepSecs:     args.SleepSecs,
					SleepQuality:  args.SleepQuality,
					AvgSleepingHR: args.AvgSleepingHR,
					Soreness:      args.Soreness,
					Fatigue:       args.Fatigue,
					Stress:        args.Stress,
					Mood:          args.Mood,
					Motivation:    args.Motivation,
					Readiness:     args.Readiness,
					SpO2:          args.SpO2,
					BodyFat:       args.BodyFat,
					Steps:         args.Steps,
					KcalConsumed:  args.KcalConsumed,
					Comments:      args.Comments,
				})
				if err != nil {
					return nil, nil, fmt.Errorf("marshaling wellness body: %w", err)
				}

				path := "/api/v1/athlete/" + apiClient.AthleteID() + "/wellness/" + url.PathEscape(args.Date)

				body, err := apiClient.Put(ctx, path, bytes.NewReader(payloadBytes))
				if err != nil {
					return nil, nil, fmt.Errorf("updating wellness: %w", err)
				}

				return textResult(body), nil, nil
			},
		)
	}
}
