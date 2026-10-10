package tools

import (
	"fmt"
	"net/url"

	"github.com/0xalexb/intervals-icu-mcp/src/app/clients/intervals"
)

type searchActivitiesArgs struct {
	Query string `json:"q"               jsonschema:"case-insensitive name search, or an exact tag search if it starts with #"`
	Limit *int   `json:"limit,omitempty" jsonschema:"maximum number of results"`
}

type searchIntervalsArgs struct {
	MinSecs      int    `json:"min_secs"      jsonschema:"minimum interval length in seconds"`
	MaxSecs      int    `json:"max_secs"      jsonschema:"maximum interval length in seconds"`
	MinIntensity int    `json:"min_intensity" jsonschema:"minimum intensity in percent of threshold"`
	MaxIntensity int    `json:"max_intensity" jsonschema:"maximum intensity in percent of threshold"`
	Type         string `json:"type,omitempty" jsonschema:"activity type, e.g. Ride or Run"`
	MinReps      *int   `json:"min_reps,omitempty" jsonschema:"minimum number of matching intervals in an activity"`
	MaxReps      *int   `json:"max_reps,omitempty" jsonschema:"maximum number of matching intervals in an activity"`
	Limit        *int   `json:"limit,omitempty"    jsonschema:"maximum number of results"`
}

// NewSearchActivitiesTool registers search_activities.
func NewSearchActivitiesTool(apiClient *intervals.Client) ToolRegistration {
	return getTool(apiClient, "search_activities",
		"Searches activities by name or tag (prefix with #) across the whole history.",
		func(c *intervals.Client, a searchActivitiesArgs) (string, url.Values, error) {
			err := required("q", a.Query)
			if err != nil {
				return "", nil, err
			}

			q := url.Values{}
			q.Set("q", a.Query)
			setInt(q, "limit", a.Limit)

			return athletePath(c, "/activities/search"), q, nil
		})
}

// NewSearchIntervalsTool registers search_intervals.
func NewSearchIntervalsTool(apiClient *intervals.Client) ToolRegistration {
	return getTool(apiClient, "search_intervals",
		"Finds past activities containing intervals of a given length and intensity, e.g. 4-6 minute efforts at 95-105% of threshold, to compare similar sessions.",
		func(c *intervals.Client, a searchIntervalsArgs) (string, url.Values, error) {
			if a.MaxSecs <= 0 || a.MinSecs > a.MaxSecs || a.MaxIntensity <= 0 || a.MinIntensity > a.MaxIntensity {
				return "", nil, fmt.Errorf("%w: need 0 < min <= max for both secs and intensity", errInvalidValue)
			}

			q := url.Values{}
			q.Set("minSecs", fmt.Sprint(a.MinSecs))
			q.Set("maxSecs", fmt.Sprint(a.MaxSecs))
			q.Set("minIntensity", fmt.Sprint(a.MinIntensity))
			q.Set("maxIntensity", fmt.Sprint(a.MaxIntensity))

			if a.Type != "" {
				q.Set("type", a.Type)
			}

			setInt(q, "minReps", a.MinReps)
			setInt(q, "maxReps", a.MaxReps)
			setInt(q, "limit", a.Limit)

			return athletePath(c, "/activities/interval-search"), q, nil
		})
}
