package tools

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/0xalexb/intervals-icu-mcp/src/app/clients/intervals"
)

const testResponse = `{"ok":true}`

func callToolWithRegistration(
	t *testing.T,
	register func(*intervals.Client) ToolRegistration,
	apiHandler http.HandlerFunc,
	name string,
	args map[string]any,
) *mcp.CallToolResult {
	t.Helper()

	apiServer := httptest.NewServer(apiHandler)
	t.Cleanup(apiServer.Close)

	c := intervals.NewTestClient(apiServer.URL, "test-key", "i12345", apiServer.Client())

	mcpServer := mcp.NewServer(&mcp.Implementation{Name: "test-server", Version: "v0.0.1"}, nil)
	register(c)(mcpServer)

	serverTransport, clientTransport := mcp.NewInMemoryTransports()

	serverCtx, serverCancel := context.WithCancel(context.Background())
	serverDone := make(chan struct{})

	go func() {
		defer close(serverDone)

		_ = mcpServer.Run(serverCtx, serverTransport)
	}()

	t.Cleanup(func() {
		serverCancel()
		<-serverDone
	})

	mcpClient := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "v0.0.1"}, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	session, err := mcpClient.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatalf("expected client to connect, got: %v", err)
	}

	t.Cleanup(func() { _ = session.Close() })

	result, err := session.CallTool(ctx, &mcp.CallToolParams{Name: name, Arguments: args})
	if err != nil {
		t.Fatalf("expected CallTool to succeed, got: %v", err)
	}

	return result
}

func TestNewTools_Requests(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		register  func(*intervals.Client) ToolRegistration
		args      map[string]any
		method    string
		path      string
		query     map[string]string
		wantBody  map[string]any
		checkBody bool
	}{
		{
			name:     "update_wellness",
			register: NewUpdateWellnessTool,
			args:     map[string]any{"date": "2024-06-01", "hrv": 62.5, "restingHR": 48, "comments": "fresh"},
			method:   http.MethodPut,
			path:     "/api/v1/athlete/i12345/wellness/2024-06-01",
			wantBody: map[string]any{"hrv": 62.5, "restingHR": float64(48), "comments": "fresh"},
		},
		{
			name:     "get_activity_streams",
			register: NewGetActivityStreamsTool,
			args:     map[string]any{"activity_id": "i999", "types": "watts,heartrate"},
			method:   http.MethodGet,
			path:     "/api/v1/activity/i999/streams.json",
			query:    map[string]string{"types": "watts,heartrate"},
		},
		{
			name:     "get_activity_intervals",
			register: NewGetActivityIntervalsTool,
			args:     map[string]any{"activity_id": "i999"},
			method:   http.MethodGet,
			path:     "/api/v1/activity/i999/intervals",
		},
		{
			name:     "update_activity",
			register: NewUpdateActivityTool,
			args:     map[string]any{"activity_id": "i999", "name": "Tempo", "icu_rpe": 7, "race": false},
			method:   http.MethodPut,
			path:     "/api/v1/activity/i999",
			wantBody: map[string]any{"name": "Tempo", "icu_rpe": float64(7), "race": false},
		},
		{
			name:     "get_fitness_summary",
			register: NewGetFitnessSummaryTool,
			args:     map[string]any{"start": "2024-05-01", "end": "2024-06-01"},
			method:   http.MethodGet,
			path:     "/api/v1/athlete/i12345/athlete-summary.json",
			query:    map[string]string{"start": "2024-05-01", "end": "2024-06-01"},
		},
		{
			name:     "list_workouts",
			register: NewListWorkoutsTool,
			args:     map[string]any{},
			method:   http.MethodGet,
			path:     "/api/v1/athlete/i12345/workouts",
		},
		{
			name:     "create_workout",
			register: NewCreateWorkoutTool,
			args:     map[string]any{"name": "Sweet spot", "description": "- 20m 90%", "folder_id": 5, "indoor": true},
			method:   http.MethodPost,
			path:     "/api/v1/athlete/i12345/workouts",
			wantBody: map[string]any{"name": "Sweet spot", "description": "- 20m 90%", "folder_id": float64(5), "indoor": true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

				if user, pass, ok := r.BasicAuth(); !ok || user != "API_KEY" || pass != "test-key" {
					t.Errorf("unexpected basic auth: %s:%s (ok=%v)", user, pass, ok)
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

func TestNewTools_Validation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		register func(*intervals.Client) ToolRegistration
		args     map[string]any
	}{
		{"update_wellness", NewUpdateWellnessTool, map[string]any{"date": "06-01-2024"}},
		{"get_activity_streams", NewGetActivityStreamsTool, map[string]any{"activity_id": "i999", "types": ""}},
		{"get_activity_intervals", NewGetActivityIntervalsTool, map[string]any{"activity_id": ""}},
		{"update_activity", NewUpdateActivityTool, map[string]any{"activity_id": ""}},
		{"get_fitness_summary", NewGetFitnessSummaryTool, map[string]any{"start": "2024-05-01", "end": "bad"}},
		{"create_workout", NewCreateWorkoutTool, map[string]any{"name": ""}},
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

func TestNewTools_APIError(t *testing.T) {
	t.Parallel()

	result := callToolWithRegistration(t, NewGetActivityIntervalsTool, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}, "get_activity_intervals", map[string]any{"activity_id": "i999"})

	if !result.IsError {
		t.Fatal("expected IsError to be true on API 404")
	}
}
