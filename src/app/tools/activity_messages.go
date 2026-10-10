package tools

import (
	"net/http"
	"net/url"

	"github.com/0xalexb/intervals-icu-mcp/src/app/clients/intervals"
)

type listActivityMessagesArgs struct {
	ActivityID string `json:"activity_id"     jsonschema:"the activity ID"`
	Limit      *int   `json:"limit,omitempty" jsonschema:"maximum number of messages"`
}

type addActivityMessageArgs struct {
	ActivityID string `json:"activity_id" jsonschema:"the activity ID"`
	Content    string `json:"content"     jsonschema:"message text"`
}

type activityMessagePayload struct {
	Content string `json:"content"`
}

// NewListActivityMessagesTool registers list_activity_messages.
func NewListActivityMessagesTool(apiClient *intervals.Client) ToolRegistration {
	return getTool(apiClient, "list_activity_messages",
		"Lists the comments on an activity (athlete notes, coach feedback).",
		func(_ *intervals.Client, a listActivityMessagesArgs) (string, url.Values, error) {
			err := required("activity_id", a.ActivityID)
			if err != nil {
				return "", nil, err
			}

			q := url.Values{}
			setInt(q, "limit", a.Limit)

			return activityPath(a.ActivityID, "/messages"), q, nil
		})
}

// NewAddActivityMessageTool registers add_activity_message.
func NewAddActivityMessageTool(apiClient *intervals.Client) ToolRegistration {
	return writeTool(apiClient, "add_activity_message",
		"Adds a comment to an activity, e.g. coaching feedback on the session.", http.MethodPost,
		func(_ *intervals.Client, a addActivityMessageArgs) (string, url.Values, any, error) {
			err := required("activity_id", a.ActivityID)
			if err != nil {
				return "", nil, nil, err
			}

			err = required("content", a.Content)
			if err != nil {
				return "", nil, nil, err
			}

			return activityPath(a.ActivityID, "/messages"), nil, activityMessagePayload{Content: a.Content}, nil
		})
}
