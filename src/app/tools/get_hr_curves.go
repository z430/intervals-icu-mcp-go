package tools

import (
	"net/url"

	"github.com/0xalexb/intervals-icu-mcp/src/app/clients/intervals"
)

type getCurvesArgs struct {
	Type   string `json:"type"             jsonschema:"activity type, e.g. Ride, Run, Swim"`
	Newest string `json:"newest,omitempty" jsonschema:"newest date in yyyy-MM-dd format (inclusive), defaults to today"`
	Curves string `json:"curves,omitempty" jsonschema:"comma-separated curve periods to compare, e.g. 42d,1y; defaults to the last year"`
}

func curvesRequest(path string) requestBuilder[getCurvesArgs] {
	return func(c *intervals.Client, a getCurvesArgs) (string, url.Values, error) {
		err := required("type", a.Type)
		if err != nil {
			return "", nil, err
		}

		err = optionalDate(a.Newest)
		if err != nil {
			return "", nil, err
		}

		q := url.Values{}
		q.Set("type", a.Type)

		if a.Newest != "" {
			q.Set("newest", a.Newest)
		}

		if a.Curves != "" {
			q.Set("curves", a.Curves)
		}

		return athletePath(c, path), q, nil
	}
}

// NewGetHRCurvesTool registers get_hr_curves.
func NewGetHRCurvesTool(apiClient *intervals.Client) ToolRegistration {
	return getTool(apiClient, "get_hr_curves",
		"Returns the athlete's heart rate curve (best sustained heart rate per duration) filtered by activity type.",
		curvesRequest("/hr-curves.json"))
}

// NewGetPaceCurvesTool registers get_pace_curves.
func NewGetPaceCurvesTool(apiClient *intervals.Client) ToolRegistration {
	return getTool(apiClient, "get_pace_curves",
		"Returns the athlete's pace curve (best pace per distance) filtered by activity type, e.g. Run or Swim.",
		curvesRequest("/pace-curves.json"))
}
