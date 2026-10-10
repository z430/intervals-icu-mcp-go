package tools

import (
	"context"
	"fmt"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/0xalexb/intervals-icu-mcp/src/app/clients/intervals"
)

type getFitnessSummaryArgs struct {
	Start string `json:"start" jsonschema:"start date in yyyy-MM-dd format"`
	End   string `json:"end"   jsonschema:"end date in yyyy-MM-dd format"`
}

// NewGetFitnessSummaryTool returns a ToolRegistration that registers the get_fitness_summary tool on an MCP server.
func NewGetFitnessSummaryTool(apiClient *intervals.Client) ToolRegistration {
	return func(server *mcp.Server) {
		mcp.AddTool(server,
			&mcp.Tool{
				Name:        "get_fitness_summary",
				Description: "Returns the athlete's training summary over a date range: weekly fitness (CTL), fatigue (ATL), form, time, distance and training load.",
			},
			func(ctx context.Context, _ *mcp.CallToolRequest, args getFitnessSummaryArgs) (*mcp.CallToolResult, any, error) {
				if args.Start == "" {
					return nil, nil, errMissingStart
				}

				if args.End == "" {
					return nil, nil, errMissingEnd
				}

				for _, d := range []string{args.Start, args.End} {
					err := validateDateFormat(d)
					if err != nil {
						return nil, nil, err
					}
				}

				query := url.Values{}
				query.Set("start", args.Start)
				query.Set("end", args.End)

				body, err := apiClient.Get(ctx, "/api/v1/athlete/"+apiClient.AthleteID()+"/athlete-summary.json", query)
				if err != nil {
					return nil, nil, fmt.Errorf("fetching fitness summary: %w", err)
				}

				return textResult(body), nil, nil
			},
		)
	}
}
