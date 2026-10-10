package tools

import (
	"fmt"
	"net/url"

	"github.com/0xalexb/intervals-icu-mcp/src/app/clients/intervals"
)

type getBestEffortsArgs struct {
	ActivityID string   `json:"activity_id"         jsonschema:"the activity ID"`
	Stream     string   `json:"stream"              jsonschema:"stream to search, e.g. watts or heartrate or velocity_smooth"`
	Duration   *int     `json:"duration,omitempty"  jsonschema:"effort length in seconds (use this or distance)"`
	Distance   *float64 `json:"distance,omitempty"  jsonschema:"effort length in meters (use this or duration)"`
	Count      *int     `json:"count,omitempty"     jsonschema:"number of efforts to return"`
}

// NewGetBestEffortsTool registers get_best_efforts.
func NewGetBestEffortsTool(apiClient *intervals.Client) ToolRegistration {
	return getTool(apiClient, "get_best_efforts",
		"Finds the best efforts within one activity, e.g. best 5 minute power or fastest 1 km. Provide either duration or distance.",
		func(_ *intervals.Client, a getBestEffortsArgs) (string, url.Values, error) {
			err := required("activity_id", a.ActivityID)
			if err != nil {
				return "", nil, err
			}

			err = required("stream", a.Stream)
			if err != nil {
				return "", nil, err
			}

			if a.Duration == nil && a.Distance == nil {
				return "", nil, fmt.Errorf("%w: duration or distance", errRequiredField)
			}

			q := url.Values{}
			q.Set("stream", a.Stream)
			setInt(q, "duration", a.Duration)
			setInt(q, "count", a.Count)

			if a.Distance != nil {
				q.Set("distance", fmt.Sprint(*a.Distance))
			}

			return activityPath(a.ActivityID, "/best-efforts"), q, nil
		})
}
