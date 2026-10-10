package tools

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/0xalexb/intervals-icu-mcp/src/app/clients/intervals"
)

func TestCoachTools_Requests(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		register func(*intervals.Client) ToolRegistration
		args     map[string]any
		method   string
		path     string
		query    map[string]string
		wantBody map[string]any
	}{
		{"list_sport_settings", NewListSportSettingsTool, map[string]any{}, http.MethodGet,
			"/api/v1/athlete/i12345/sport-settings", nil, nil},
		{"update_sport_settings", NewUpdateSportSettingsTool,
			map[string]any{"settings_id": "42", "ftp": 260, "recalc_hr_zones": true}, http.MethodPut,
			"/api/v1/athlete/i12345/sport-settings/42", map[string]string{"recalcHrZones": "true"},
			map[string]any{"ftp": float64(260)}},
		{"get_best_efforts", NewGetBestEffortsTool,
			map[string]any{"activity_id": "i9", "stream": "watts", "duration": 300, "count": 2}, http.MethodGet,
			"/api/v1/activity/i9/best-efforts", map[string]string{"stream": "watts", "duration": "300", "count": "2"}, nil},
		{"get_activity_histogram", NewGetActivityHistogramTool,
			map[string]any{"activity_id": "i9", "kind": "power", "bucket_size": 25}, http.MethodGet,
			"/api/v1/activity/i9/power-histogram", map[string]string{"bucketSize": "25"}, nil},
		{"get_activity_histogram", NewGetActivityHistogramTool,
			map[string]any{"activity_id": "i9", "kind": "gap"}, http.MethodGet,
			"/api/v1/activity/i9/gap-histogram", nil, nil},
		{"get_hr_curves", NewGetHRCurvesTool,
			map[string]any{"type": "Run", "newest": "2024-06-01", "curves": "42d,1y"}, http.MethodGet,
			"/api/v1/athlete/i12345/hr-curves.json",
			map[string]string{"type": "Run", "newest": "2024-06-01", "curves": "42d,1y"}, nil},
		{"get_pace_curves", NewGetPaceCurvesTool, map[string]any{"type": "Swim"}, http.MethodGet,
			"/api/v1/athlete/i12345/pace-curves.json", map[string]string{"type": "Swim"}, nil},
		{"get_time_at_hr", NewGetTimeAtHRTool, map[string]any{"activity_id": "i9"}, http.MethodGet,
			"/api/v1/activity/i9/time-at-hr", nil, nil},
		{"get_power_vs_hr", NewGetPowerVsHRTool, map[string]any{"activity_id": "i9"}, http.MethodGet,
			"/api/v1/activity/i9/power-vs-hr.json", nil, nil},
		{"list_workout_folders", NewListWorkoutFoldersTool, map[string]any{}, http.MethodGet,
			"/api/v1/athlete/i12345/folders", nil, nil},
		{"get_workout", NewGetWorkoutTool, map[string]any{"workout_id": "7"}, http.MethodGet,
			"/api/v1/athlete/i12345/workouts/7", nil, nil},
		{"update_workout", NewUpdateWorkoutTool, map[string]any{"workout_id": "7", "name": "VO2", "indoor": false},
			http.MethodPut, "/api/v1/athlete/i12345/workouts/7", nil,
			map[string]any{"name": "VO2", "indoor": false}},
		{"delete_workout", NewDeleteWorkoutTool, map[string]any{"workout_id": "7"}, http.MethodDelete,
			"/api/v1/athlete/i12345/workouts/7", nil, nil},
		{"search_activities", NewSearchActivitiesTool, map[string]any{"q": "#race", "limit": 5}, http.MethodGet,
			"/api/v1/athlete/i12345/activities/search", map[string]string{"q": "#race", "limit": "5"}, nil},
		{"search_intervals", NewSearchIntervalsTool,
			map[string]any{"min_secs": 240, "max_secs": 360, "min_intensity": 95, "max_intensity": 105, "type": "Ride"},
			http.MethodGet, "/api/v1/athlete/i12345/activities/interval-search",
			map[string]string{"minSecs": "240", "maxSecs": "360", "minIntensity": "95", "maxIntensity": "105", "type": "Ride"}, nil},
		{"list_activity_messages", NewListActivityMessagesTool, map[string]any{"activity_id": "i9", "limit": 3},
			http.MethodGet, "/api/v1/activity/i9/messages", map[string]string{"limit": "3"}, nil},
		{"add_activity_message", NewAddActivityMessageTool, map[string]any{"activity_id": "i9", "content": "Good job"},
			http.MethodPost, "/api/v1/activity/i9/messages", nil, map[string]any{"content": "Good job"}},
		{"list_gear", NewListGearTool, map[string]any{}, http.MethodGet, "/api/v1/athlete/i12345/gear.json", nil, nil},
		{"get_weather_forecast", NewGetWeatherForecastTool, map[string]any{}, http.MethodGet,
			"/api/v1/athlete/i12345/weather-forecast", nil, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name+"_"+tt.path, func(t *testing.T) {
			t.Parallel()

			var gotBody map[string]any

			result := callToolWithRegistration(t, tt.register, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != tt.method {
					t.Errorf("expected %s, got: %s", tt.method, r.Method)
				}

				if r.URL.Path != tt.path {
					t.Errorf("expected path %s, got: %s", tt.path, r.URL.Path)
				}

				for k, v := range tt.query {
					if got := r.URL.Query().Get(k); got != v {
						t.Errorf("expected query %s=%s, got: %s", k, v, got)
					}
				}

				if tt.wantBody != nil {
					raw, _ := io.ReadAll(r.Body)
					_ = json.Unmarshal(raw, &gotBody)
				}

				_, _ = w.Write([]byte(testResponse))
			}, tt.name, tt.args)

			if result.IsError {
				t.Fatalf("expected success, got error result: %+v", result.Content)
			}

			text, ok := result.Content[0].(*mcp.TextContent)
			if !ok || text.Text != testResponse {
				t.Fatalf("expected text %q, got: %+v", testResponse, result.Content[0])
			}

			if tt.wantBody != nil && !reflect.DeepEqual(gotBody, tt.wantBody) {
				t.Errorf("expected body %v, got: %v", tt.wantBody, gotBody)
			}
		})
	}
}

func TestCoachTools_Validation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		register func(*intervals.Client) ToolRegistration
		args     map[string]any
	}{
		{"update_sport_settings", NewUpdateSportSettingsTool, map[string]any{"settings_id": ""}},
		{"get_best_efforts", NewGetBestEffortsTool, map[string]any{"activity_id": "i9", "stream": "watts"}},
		{"get_best_efforts", NewGetBestEffortsTool, map[string]any{"activity_id": "i9", "duration": 60}},
		{"get_activity_histogram", NewGetActivityHistogramTool, map[string]any{"activity_id": "i9", "kind": "bogus"}},
		{"get_hr_curves", NewGetHRCurvesTool, map[string]any{"type": ""}},
		{"get_pace_curves", NewGetPaceCurvesTool, map[string]any{"type": "Run", "newest": "bad"}},
		{"get_time_at_hr", NewGetTimeAtHRTool, map[string]any{"activity_id": ""}},
		{"get_workout", NewGetWorkoutTool, map[string]any{"workout_id": ""}},
		{"delete_workout", NewDeleteWorkoutTool, map[string]any{"workout_id": ""}},
		{"search_activities", NewSearchActivitiesTool, map[string]any{"q": ""}},
		{"search_intervals", NewSearchIntervalsTool,
			map[string]any{"min_secs": 300, "max_secs": 60, "min_intensity": 90, "max_intensity": 100}},
		{"add_activity_message", NewAddActivityMessageTool, map[string]any{"activity_id": "i9", "content": ""}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := callToolWithRegistration(t, tt.register, func(_ http.ResponseWriter, _ *http.Request) {
				t.Error("API should not be called for validation errors")
			}, tt.name, tt.args)

			if !result.IsError {
				t.Fatal("expected IsError to be true")
			}
		})
	}
}

func TestNormalizeLocalDateTime(t *testing.T) {
	t.Parallel()

	tests := map[string]string{
		"2024-06-01":          "2024-06-01T00:00:00",
		"2024-06-01T07:30":    "2024-06-01T07:30:00",
		"2024-06-01T07:30:15": "2024-06-01T07:30:15",
	}

	for in, want := range tests {
		got, err := normalizeLocalDateTime(in)
		if err != nil || got != want {
			t.Errorf("normalizeLocalDateTime(%q) = %q, %v; want %q", in, got, err, want)
		}
	}

	if _, err := normalizeLocalDateTime("06/01/2024"); err == nil {
		t.Error("expected error for unsupported date format")
	}
}

func TestEventCategoryAndPayload(t *testing.T) {
	t.Parallel()

	if got, err := normalizeEventCategory(" race_a "); err != nil || got != "RACE_A" {
		t.Errorf("expected RACE_A, got %q, %v", got, err)
	}

	if _, err := normalizeEventCategory("RACE"); err == nil {
		t.Error("expected RACE to be rejected")
	}

	var body map[string]any

	result := callToolWithRegistration(t, NewCreateEventTool, func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &body)
		_, _ = w.Write([]byte(testResponse))
	}, "create_event", map[string]any{
		"name": "Ride", "start_date_local": "2030-01-01", "category": "workout", "type": "Ride",
		"moving_time": 3600, "training_load": 49.6,
	})

	if result.IsError {
		t.Fatalf("expected success, got: %+v", result.Content)
	}

	want := map[string]any{
		"name": "Ride", "start_date_local": "2030-01-01T00:00:00", "category": "WORKOUT", "type": "Ride",
		"moving_time": float64(3600), "icu_training_load": float64(50),
	}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("expected body %v, got: %v", want, body)
	}

	for _, args := range []map[string]any{
		{"name": "Race", "start_date_local": "2030-01-01", "category": "RACE_A"},
		{"name": "Race", "start_date_local": "2030-01-01", "category": "RACE", "type": "Run"},
	} {
		res := callToolWithRegistration(t, NewCreateEventTool, func(_ http.ResponseWriter, _ *http.Request) {
			t.Error("API should not be called for invalid events")
		}, "create_event", args)
		if !res.IsError {
			t.Errorf("expected validation error for %v", args)
		}
	}
}
