// Package tools contains MCP tool registrations for the intervals.icu MCP server.
package tools

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var (
	errMissingEventID    = errors.New("event_id is required")
	errMissingDate       = errors.New("date is required")
	errMissingOldest     = errors.New("oldest is required")
	errMissingNewest     = errors.New("newest is required")
	errMissingActivityID = errors.New("activity_id is required")
	errMissingType       = errors.New("type is required")
	errMissingName       = errors.New("name is required")
	errMissingStartDate  = errors.New("start_date_local is required")
	errMissingCategory   = errors.New("category is required")
	errMissingTypes      = errors.New("types is required")
	errMissingStart      = errors.New("start is required")
	errMissingEnd        = errors.New("end is required")
)

func validateDateFormat(s string) error {
	_, err := time.Parse(time.DateOnly, s)
	if err != nil {
		return fmt.Errorf("invalid date format %q, expected yyyy-MM-dd: %w", s, err)
	}

	return nil
}

// ToolRegistration is a function that registers an MCP tool on the given server.
type ToolRegistration func(server *mcp.Server)

const localDateTimeLayout = "2006-01-02T15:04:05"

var errInvalidCategory = errors.New("invalid category")

var eventCategories = []string{ //nolint:gochecknoglobals // static enum of API-accepted categories.
	"WORKOUT", "RACE_A", "RACE_B", "RACE_C", "NOTE", "PLAN", "HOLIDAY", "SICK",
	"INJURED", "SET_EFTP", "FITNESS_DAYS", "SEASON_START", "TARGET", "SET_FITNESS",
}

// normalizeLocalDateTime accepts yyyy-MM-dd, yyyy-MM-ddTHH:mm or yyyy-MM-ddTHH:mm:ss and returns
// the yyyy-MM-ddTHH:mm:ss form the Intervals.icu API requires. A bare date means midnight.
func normalizeLocalDateTime(s string) (string, error) {
	for _, layout := range []string{time.DateOnly, "2006-01-02T15:04", localDateTimeLayout} {
		parsed, err := time.Parse(layout, s)
		if err == nil {
			return parsed.Format(localDateTimeLayout), nil
		}
	}

	return "", fmt.Errorf("invalid date %q, expected yyyy-MM-dd or yyyy-MM-ddTHH:mm:ss", s)
}

// normalizeEventCategory upper-cases and validates an event category against the API enum.
func normalizeEventCategory(s string) (string, error) {
	category := strings.ToUpper(strings.TrimSpace(s))
	if !slices.Contains(eventCategories, category) {
		return "", fmt.Errorf("%w %q, must be one of: %s", errInvalidCategory, s, strings.Join(eventCategories, ", "))
	}

	return category, nil
}
