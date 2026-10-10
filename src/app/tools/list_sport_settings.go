package tools

import (
	"net/url"

	"github.com/0xalexb/intervals-icu-mcp/src/app/clients/intervals"
)

// NewListSportSettingsTool registers list_sport_settings: FTP, LTHR, max HR, thresholds and zones per sport.
func NewListSportSettingsTool(apiClient *intervals.Client) ToolRegistration {
	return getTool(apiClient, "list_sport_settings",
		"Returns the athlete's sport settings: FTP, W', LTHR, max HR, threshold pace and the power/HR/pace zones for each sport. Use this to interpret any workout intensity.",
		func(c *intervals.Client, _ struct{}) (string, url.Values, error) {
			return athletePath(c, "/sport-settings"), nil, nil
		})
}
