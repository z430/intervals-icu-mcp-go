package tools

import (
	"net/http"
	"net/url"

	"github.com/0xalexb/intervals-icu-mcp/src/app/clients/intervals"
)

type updateSportSettingsArgs struct {
	SettingsID    string   `json:"settings_id"             jsonschema:"sport settings ID, from list_sport_settings"`
	FTP           *int     `json:"ftp,omitempty"           jsonschema:"functional threshold power in watts"`
	IndoorFTP     *int     `json:"indoor_ftp,omitempty"    jsonschema:"indoor FTP in watts"`
	WPrime        *int     `json:"w_prime,omitempty"       jsonschema:"W' (anaerobic work capacity) in joules"`
	LTHR          *int     `json:"lthr,omitempty"          jsonschema:"lactate threshold heart rate in bpm"`
	MaxHR         *int     `json:"max_hr,omitempty"        jsonschema:"maximum heart rate in bpm"`
	ThresholdPace *float64 `json:"threshold_pace,omitempty" jsonschema:"threshold pace in meters per second"`
	RecalcHRZones bool     `json:"recalc_hr_zones,omitempty" jsonschema:"recalculate HR zones from the new LTHR/max HR"`
}

type updateSportSettingsPayload struct {
	FTP           *int     `json:"ftp,omitempty"`
	IndoorFTP     *int     `json:"indoor_ftp,omitempty"`
	WPrime        *int     `json:"w_prime,omitempty"`
	LTHR          *int     `json:"lthr,omitempty"`
	MaxHR         *int     `json:"max_hr,omitempty"`
	ThresholdPace *float64 `json:"threshold_pace,omitempty"`
}

// NewUpdateSportSettingsTool registers update_sport_settings.
func NewUpdateSportSettingsTool(apiClient *intervals.Client) ToolRegistration {
	return writeTool(apiClient, "update_sport_settings",
		"Updates thresholds (FTP, W', LTHR, max HR, threshold pace) of a sport settings entry. Only the provided fields are changed.",
		http.MethodPut,
		func(c *intervals.Client, a updateSportSettingsArgs) (string, url.Values, any, error) {
			err := required("settings_id", a.SettingsID)
			if err != nil {
				return "", nil, nil, err
			}

			query := url.Values{}
			query.Set("recalcHrZones", boolString(a.RecalcHRZones))

			return athletePath(c, "/sport-settings/"+url.PathEscape(a.SettingsID)), query, updateSportSettingsPayload{
				FTP: a.FTP, IndoorFTP: a.IndoorFTP, WPrime: a.WPrime, LTHR: a.LTHR,
				MaxHR: a.MaxHR, ThresholdPace: a.ThresholdPace,
			}, nil
		})
}

func boolString(b bool) string {
	if b {
		return "true"
	}

	return "false"
}
