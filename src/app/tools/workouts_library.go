package tools

import (
	"net/http"
	"net/url"

	"github.com/0xalexb/intervals-icu-mcp/src/app/clients/intervals"
)

type workoutIDArgs struct {
	WorkoutID string `json:"workout_id" jsonschema:"the workout ID"`
}

type updateWorkoutArgs struct {
	WorkoutID   string `json:"workout_id"            jsonschema:"the workout ID to update"`
	Name        string `json:"name,omitempty"        jsonschema:"workout name"`
	Description string `json:"description,omitempty" jsonschema:"workout steps in Intervals.icu workout syntax plus optional notes"`
	Type        string `json:"type,omitempty"        jsonschema:"sport type (e.g. Ride or Run or Swim)"`
	FolderID    *int   `json:"folder_id,omitempty"   jsonschema:"folder ID to move the workout to"`
	Indoor      *bool  `json:"indoor,omitempty"      jsonschema:"whether this is an indoor workout"`
}

type updateWorkoutPayload struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type,omitempty"`
	FolderID    *int   `json:"folder_id,omitempty"`
	Indoor      *bool  `json:"indoor,omitempty"`
}

func workoutPath(c *intervals.Client, id string) (string, error) {
	err := required("workout_id", id)
	if err != nil {
		return "", err
	}

	return athletePath(c, "/workouts/"+url.PathEscape(id)), nil
}

// NewListWorkoutFoldersTool registers list_workout_folders.
func NewListWorkoutFoldersTool(apiClient *intervals.Client) ToolRegistration {
	return getTool(apiClient, "list_workout_folders",
		"Lists the athlete's workout folders and training plans, including their workouts.",
		func(c *intervals.Client, _ struct{}) (string, url.Values, error) {
			return athletePath(c, "/folders"), nil, nil
		})
}

// NewGetWorkoutTool registers get_workout.
func NewGetWorkoutTool(apiClient *intervals.Client) ToolRegistration {
	return getTool(apiClient, "get_workout", "Returns one workout from the workout library.",
		func(c *intervals.Client, a workoutIDArgs) (string, url.Values, error) {
			path, err := workoutPath(c, a.WorkoutID)

			return path, nil, err
		})
}

// NewUpdateWorkoutTool registers update_workout.
func NewUpdateWorkoutTool(apiClient *intervals.Client) ToolRegistration {
	return writeTool(apiClient, "update_workout",
		"Updates a workout in the library. Only the provided fields are changed.", http.MethodPut,
		func(c *intervals.Client, a updateWorkoutArgs) (string, url.Values, any, error) {
			path, err := workoutPath(c, a.WorkoutID)
			if err != nil {
				return "", nil, nil, err
			}

			return path, nil, updateWorkoutPayload{
				Name: a.Name, Description: a.Description, Type: a.Type, FolderID: a.FolderID, Indoor: a.Indoor,
			}, nil
		})
}

// NewDeleteWorkoutTool registers delete_workout.
func NewDeleteWorkoutTool(apiClient *intervals.Client) ToolRegistration {
	return writeTool(apiClient, "delete_workout",
		"Deletes a workout from the library. Calendar events created from it are not removed.", http.MethodDelete,
		func(c *intervals.Client, a workoutIDArgs) (string, url.Values, any, error) {
			path, err := workoutPath(c, a.WorkoutID)

			return path, nil, nil, err
		})
}
