package tools

import (
	"net/url"

	"github.com/0xalexb/intervals-icu-mcp/src/app/clients/intervals"
)

// NewListGearTool registers list_gear.
func NewListGearTool(apiClient *intervals.Client) ToolRegistration {
	return getTool(apiClient, "list_gear",
		"Lists the athlete's gear (bikes, shoes) with accumulated distance and time, and maintenance reminders.",
		func(c *intervals.Client, _ struct{}) (string, url.Values, error) {
			return athletePath(c, "/gear.json"), nil, nil
		})
}

// NewGetWeatherForecastTool registers get_weather_forecast.
func NewGetWeatherForecastTool(apiClient *intervals.Client) ToolRegistration {
	return getTool(apiClient, "get_weather_forecast",
		"Returns the weather forecast for the athlete's configured location, for planning outdoor sessions.",
		func(c *intervals.Client, _ struct{}) (string, url.Values, error) {
			return athletePath(c, "/weather-forecast"), nil, nil
		})
}
