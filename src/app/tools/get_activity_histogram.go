package tools

import (
	"fmt"
	"net/url"

	"github.com/0xalexb/intervals-icu-mcp/src/app/clients/intervals"
)

type getActivityHistogramArgs struct {
	ActivityID string `json:"activity_id"           jsonschema:"the activity ID"`
	Kind       string `json:"kind"                  jsonschema:"histogram kind: power or hr or pace or gap"`
	BucketSize *int   `json:"bucket_size,omitempty" jsonschema:"bucket size (power in watts, hr in bpm); ignored for pace and gap"`
}

// NewGetActivityHistogramTool registers get_activity_histogram.
func NewGetActivityHistogramTool(apiClient *intervals.Client) ToolRegistration {
	return getTool(apiClient, "get_activity_histogram",
		"Returns the distribution of time spent at each power, heart rate, pace or grade-adjusted pace within an activity.",
		func(_ *intervals.Client, a getActivityHistogramArgs) (string, url.Values, error) {
			err := required("activity_id", a.ActivityID)
			if err != nil {
				return "", nil, err
			}

			switch a.Kind {
			case "power", "hr":
				q := url.Values{}
				setInt(q, "bucketSize", a.BucketSize)

				return activityPath(a.ActivityID, "/"+a.Kind+"-histogram"), q, nil
			case "pace", "gap":
				return activityPath(a.ActivityID, "/"+a.Kind+"-histogram"), nil, nil
			default:
				return "", nil, fmt.Errorf("%w: kind must be power, hr, pace or gap", errInvalidValue)
			}
		})
}
