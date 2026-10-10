package tools

import (
	"net/url"

	"github.com/0xalexb/intervals-icu-mcp/src/app/clients/intervals"
)

type activityIDArgs struct {
	ActivityID string `json:"activity_id" jsonschema:"the activity ID"`
}

func activityRequest(suffix string) requestBuilder[activityIDArgs] {
	return func(_ *intervals.Client, a activityIDArgs) (string, url.Values, error) {
		err := required("activity_id", a.ActivityID)
		if err != nil {
			return "", nil, err
		}

		return activityPath(a.ActivityID, suffix), nil, nil
	}
}

// NewGetTimeAtHRTool registers get_time_at_hr.
func NewGetTimeAtHRTool(apiClient *intervals.Client) ToolRegistration {
	return getTool(apiClient, "get_time_at_hr",
		"Returns time spent at each heart rate value for an activity.",
		activityRequest("/time-at-hr"))
}

// NewGetPowerVsHRTool registers get_power_vs_hr.
func NewGetPowerVsHRTool(apiClient *intervals.Client) ToolRegistration {
	return getTool(apiClient, "get_power_vs_hr",
		"Returns the power versus heart rate relationship of an activity, including aerobic decoupling, to judge aerobic fitness and fatigue.",
		activityRequest("/power-vs-hr.json"))
}
