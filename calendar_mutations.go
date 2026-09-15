package main

import (
	"context"
	"fmt"
	"slices"
	"strings"

	model "github.com/dooray-go/dooray-sdk/openapi/model/calendar"
	"github.com/dooray-go/dooray-sdk/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

var calendarDeleteTypes = []model.DeleteType{model.DeleteTypeThis, model.DeleteTypeWholeFromThis, model.DeleteTypeWhole}

type calendarUpdateFunc func(context.Context, string, string, string, model.UpdateEventRequest) (*model.EventResponse, error)

type calendarDeleteFunc func(context.Context, string, string, string, model.DeleteEventRequest) (*model.EventResponse, error)

func calendarMutationTools(s *server.MCPServer, token *string, update calendarUpdateFunc, deleteEvent calendarDeleteFunc) {
	calendarUpdateEventTool(s, token, update)
	calendarDeleteEventTool(s, token, deleteEvent)
}

func calendarUpdateEventTool(s *server.MCPServer, token *string, update calendarUpdateFunc) {
	tool := mcp.NewTool("dooray_calendar_update_event",
		mcp.WithDescription("Update fields of an existing Dooray calendar event. Provide at least one of subject, content, startedAt, endedAt, wholeDayFlag, or location. Omitted fields stay unchanged; empty strings cannot clear a field. Recurrence and attendees cannot be updated. When changing wholeDayFlag, send both startedAt and endedAt."),
		mcp.WithString("operation", mcp.Required(), mcp.Description("The operation to perform"), mcp.Enum("update_event")),
		mcp.WithString("calendarId", mcp.Required(), mcp.Description("Target calendar ID from dooray_calendar_calendars or dooray_calendar_events")),
		mcp.WithString("eventId", mcp.Required(), mcp.Description("Target event ID from dooray_calendar_events")),
		mcp.WithString("subject", mcp.Description("New event title")),
		mcp.WithString("content", mcp.Description("New event body (text/html)")),
		mcp.WithString("startedAt", mcp.Description("New start time in ISO 8601, e.g. 2025-04-11T10:00:00+09:00")),
		mcp.WithString("endedAt", mcp.Description("New end time in ISO 8601, e.g. 2025-04-11T11:00:00+09:00")),
		mcp.WithBoolean("wholeDayFlag", mcp.Description("Set true for a whole-day event. When true, startedAt/endedAt use date-only format such as 2025-04-11+09:00")),
		mcp.WithString("location", mcp.Description("New event location")),
	)
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		if err := calendarRequireOperation(args, "update_event"); err != nil {
			return calendarToolError(err)
		}
		calendarID, eventID, err := calendarEventIDs(args)
		if err != nil {
			return calendarToolError(err)
		}
		event, err := calendarUpdatePayload(args)
		if err != nil {
			return calendarToolError(err)
		}
		apiToken, err := calendarAPIToken(token)
		if err != nil {
			return calendarToolError(err)
		}
		res, err := update(ctx, apiToken, calendarID, eventID, event)
		return calendarMutationResult("update", res, err)
	})
}

func calendarDeleteEventTool(s *server.MCPServer, token *string, deleteEvent calendarDeleteFunc) {
	tool := mcp.NewTool("dooray_calendar_delete_event",
		mcp.WithDescription("Delete a Dooray calendar event. For a recurring series, eventId is the occurrence ID returned by dooray_calendar_events."),
		mcp.WithString("operation", mcp.Required(), mcp.Description("The operation to perform"), mcp.Enum("delete_event")),
		mcp.WithString("calendarId", mcp.Required(), mcp.Description("Target calendar ID from dooray_calendar_calendars or dooray_calendar_events")),
		mcp.WithString("eventId", mcp.Required(), mcp.Description("Target event ID. Recurring occurrences may include a date suffix")),
		mcp.WithString("deleteType", mcp.Required(), mcp.Description("this: this occurrence only; wholeFromThis: this and following occurrences; whole: the entire series"), mcp.Enum("this", "wholeFromThis", "whole")),
	)
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		if err := calendarRequireOperation(args, "delete_event"); err != nil {
			return calendarToolError(err)
		}
		calendarID, eventID, err := calendarEventIDs(args)
		if err != nil {
			return calendarToolError(err)
		}
		deleteType, err := calendarArgumentString(args, "deleteType", true)
		if err != nil {
			return calendarToolError(err)
		}
		parsed := model.DeleteType(deleteType)
		if !slices.Contains(calendarDeleteTypes, parsed) {
			return mcp.NewToolResultError("deleteType must be this, wholeFromThis, or whole"), nil
		}
		apiToken, err := calendarAPIToken(token)
		if err != nil {
			return calendarToolError(err)
		}
		res, err := deleteEvent(ctx, apiToken, calendarID, eventID, model.DeleteEventRequest{DeleteType: parsed})
		return calendarMutationResult("delete", res, err)
	})
}

func calendarMutationResult(action string, res *model.EventResponse, err error) (*mcp.CallToolResult, error) {
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to %s event: %v", action, err)), nil
	}
	if res == nil {
		return mcp.NewToolResultError("Dooray returned an empty response"), nil
	}
	if !res.Header.IsSuccessful {
		return mcp.NewToolResultError(fmt.Sprintf("Dooray rejected event %s: %s", action, res.Header.ResultMessage)), nil
	}
	return mcp.NewToolResultText(res.RawJSON), nil
}

func calendarRequireOperation(args map[string]any, want string) error {
	op, err := calendarArgumentString(args, "operation", true)
	if err != nil {
		return err
	}
	if op != want {
		return fmt.Errorf("operation must be %s", want)
	}
	return nil
}

func calendarEventIDs(args map[string]any) (string, string, error) {
	calendarID, err := calendarRequiredID(args, "calendarId")
	if err != nil {
		return "", "", err
	}
	eventID, err := calendarRequiredID(args, "eventId")
	if err != nil {
		return "", "", err
	}
	return calendarID, eventID, nil
}

func calendarRequiredID(args map[string]any, key string) (string, error) {
	id, err := calendarArgumentString(args, key, true)
	if err != nil {
		return "", err
	}
	if !calendarIDValid(id) {
		return "", fmt.Errorf("%s must be a single %s", key, strings.TrimSuffix(key, "Id")+" ID")
	}
	return id, nil
}

func calendarUpdatePayload(args map[string]any) (model.UpdateEventRequest, error) {
	var event model.UpdateEventRequest
	for _, key := range []string{"subject", "content", "location"} {
		value, err := calendarArgumentString(args, key, false)
		if err != nil {
			return event, err
		}
		if _, present := args[key]; !present {
			continue
		}
		if strings.TrimSpace(value) == "" {
			return event, fmt.Errorf("%s must be a non-empty string", key)
		}
		switch key {
		case "content":
			event.Body = new(model.Body{MimeType: "text/html", Content: value})
		case "subject":
			event.Subject = new(value)
		default:
			event.Location = new(value)
		}
	}
	wholeDay, wholeDayPresent, err := calendarArgumentBool(args, "wholeDayFlag")
	if err != nil {
		return event, err
	}
	if wholeDayPresent {
		event.WholeDayFlag = new(wholeDay)
	}
	for _, key := range []string{"startedAt", "endedAt"} {
		value, err := calendarArgumentString(args, key, false)
		if err != nil {
			return event, err
		}
		if _, present := args[key]; !present {
			continue
		}
		if strings.TrimSpace(value) == "" {
			return event, fmt.Errorf("%s must be a non-empty string", key)
		}
		parsed, err := utils.ConvertISO8601ToTime(value)
		if err != nil {
			return event, fmt.Errorf("%s must be ISO 8601: %w", key, err)
		}
		dateOnly := (wholeDayPresent && wholeDay) || (!wholeDayPresent && utils.IsDateOnlyFormat(value))
		var when utils.JsonTime
		if dateOnly {
			when = utils.NewJsonDate(parsed)
		} else {
			when = utils.NewJsonTime(parsed)
		}
		if key == "startedAt" {
			event.StartedAt = new(when)
		} else {
			event.EndedAt = new(when)
		}
	}
	if event == (model.UpdateEventRequest{}) {
		return event, fmt.Errorf("at least one field to update is required")
	}
	return event, nil
}

func calendarArgumentString(args map[string]any, key string, required bool) (string, error) {
	value, present := args[key]
	if !present {
		if !required {
			return "", nil
		}
		return "", fmt.Errorf("%s must be a non-empty string", key)
	}
	text, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("%s must be a string", key)
	}
	if required && strings.TrimSpace(text) == "" {
		return "", fmt.Errorf("%s must be a non-empty string", key)
	}
	return text, nil
}

func calendarArgumentBool(args map[string]any, key string) (bool, bool, error) {
	value, present := args[key]
	if !present {
		return false, false, nil
	}
	flag, ok := value.(bool)
	if !ok {
		return false, true, fmt.Errorf("%s must be a boolean", key)
	}
	return flag, true, nil
}

func calendarIDValid(id string) bool {
	return strings.TrimSpace(id) != "" && !strings.ContainsAny(id, ",/\\?#% \t\r\n") && id != "." && id != ".." && id != "*"
}

func calendarAPIToken(token *string) (string, error) {
	if token == nil || strings.TrimSpace(*token) == "" {
		return "", fmt.Errorf("Dooray token is required")
	}
	return *token, nil
}

func calendarToolError(err error) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultError(err.Error()), nil
}
