package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/0xalexb/intervals-icu-mcp/src/app/clients/intervals"
)

var (
	errRequiredField = errors.New("missing required argument")
	errInvalidValue  = errors.New("invalid argument value")
)

func required(name, value string) error {
	if value == "" {
		return fmt.Errorf("%w: %s", errRequiredField, name)
	}

	return nil
}

func optionalDate(value string) error {
	if value == "" {
		return nil
	}

	return validateDateFormat(value)
}

// requestBuilder turns tool arguments into an API path and query string.
type requestBuilder[A any] func(c *intervals.Client, args A) (string, url.Values, error)

// writeBuilder additionally returns the JSON payload for PUT/POST requests.
type writeBuilder[A any] func(c *intervals.Client, args A) (string, url.Values, any, error)

func athletePath(c *intervals.Client, suffix string) string {
	return "/api/v1/athlete/" + c.AthleteID() + suffix
}

func activityPath(activityID, suffix string) string {
	return "/api/v1/activity/" + url.PathEscape(activityID) + suffix
}

// getTool registers a read-only tool backed by a single GET request.
func getTool[A any](c *intervals.Client, name, description string, build requestBuilder[A]) ToolRegistration {
	return func(server *mcp.Server) {
		mcp.AddTool(server,
			&mcp.Tool{Name: name, Description: description},
			func(ctx context.Context, _ *mcp.CallToolRequest, args A) (*mcp.CallToolResult, any, error) {
				path, query, err := build(c, args)
				if err != nil {
					return nil, nil, err
				}

				body, err := c.Get(ctx, path, query)
				if err != nil {
					return nil, nil, fmt.Errorf("%s: %w", name, err)
				}

				return textResult(body), nil, nil
			},
		)
	}
}

// writeTool registers a tool backed by a PUT, POST or DELETE request.
func writeTool[A any](
	c *intervals.Client, name, description, method string, build writeBuilder[A],
) ToolRegistration {
	return func(server *mcp.Server) {
		mcp.AddTool(server,
			&mcp.Tool{Name: name, Description: description},
			func(ctx context.Context, _ *mcp.CallToolRequest, args A) (*mcp.CallToolResult, any, error) {
				path, query, payload, err := build(c, args)
				if err != nil {
					return nil, nil, err
				}

				var body []byte

				switch method {
				case http.MethodDelete:
					body, err = c.Delete(ctx, path, query)
				default:
					body, err = sendJSON(ctx, c, method, path, query, payload)
				}

				if err != nil {
					return nil, nil, fmt.Errorf("%s: %w", name, err)
				}

				return textResult(body), nil, nil
			},
		)
	}
}

func sendJSON(
	ctx context.Context, c *intervals.Client, method, path string, query url.Values, payload any,
) ([]byte, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshaling request body: %w", err)
	}

	if len(query) > 0 {
		path += "?" + query.Encode()
	}

	if method == http.MethodPost {
		return c.Post(ctx, path, bytes.NewReader(payloadBytes))
	}

	return c.Put(ctx, path, bytes.NewReader(payloadBytes))
}

func setInt(q url.Values, key string, v *int) {
	if v != nil {
		q.Set(key, fmt.Sprint(*v))
	}
}
