package calendar

import (
	"context"
	"errors"
	"strings"
	"testing"

	model "github.com/dooray-go/dooray-sdk/openapi/model/calendar"
	"github.com/mark3labs/mcp-go/mcp"

	"dooray_mcp/internal/mcptest"
	"github.com/mark3labs/mcp-go/server"
)

const calendarSuccessBody = `{"header":{"isSuccessful":true,"resultCode":0,"resultMessage":""},"result":null}`

func TestCalendarUpdatePayload(t *testing.T) {
	event, err := calendarUpdatePayload(map[string]any{
		"subject":      "title",
		"content":      "<p>body</p>",
		"location":     "room A",
		"startedAt":    "2025-04-11T10:00:00+09:00",
		"endedAt":      "2025-04-11T11:00:00+09:00",
		"wholeDayFlag": false,
	})
	if err != nil {
		t.Fatalf("payload: %v", err)
	}
	if event.Subject == nil || *event.Subject != "title" || event.Location == nil || *event.Location != "room A" {
		t.Errorf("unexpected strings: %#v", event)
	}
	if event.Body == nil || event.Body.MimeType != "text/html" || event.Body.Content != "<p>body</p>" {
		t.Errorf("unexpected body: %#v", event.Body)
	}
	if event.WholeDayFlag == nil || *event.WholeDayFlag {
		t.Errorf("wholeDayFlag false must be set: %#v", event.WholeDayFlag)
	}
	if event.StartedAt == nil || event.StartedAt.String() != "2025-04-11T10:00:00+09:00" || event.EndedAt == nil || event.EndedAt.String() != "2025-04-11T11:00:00+09:00" {
		t.Errorf("unexpected times: started=%v ended=%v", event.StartedAt, event.EndedAt)
	}
}

func TestCalendarUpdatePayloadWholeDayFormatsDates(t *testing.T) {
	event, err := calendarUpdatePayload(map[string]any{
		"wholeDayFlag": true,
		"startedAt":    "2025-04-11T00:00:00+09:00",
		"endedAt":      "2025-04-12T00:00:00+09:00",
	})
	if err != nil {
		t.Fatalf("payload: %v", err)
	}
	if event.WholeDayFlag == nil || !*event.WholeDayFlag {
		t.Errorf("wholeDayFlag: %#v", event.WholeDayFlag)
	}
	if event.StartedAt == nil || event.StartedAt.String() != "2025-04-11+09:00" || event.EndedAt == nil || event.EndedAt.String() != "2025-04-12+09:00" {
		t.Errorf("unexpected whole-day times: started=%v ended=%v", event.StartedAt, event.EndedAt)
	}
}

func TestCalendarUpdatePayloadKeepsDateOnlyWithoutFlag(t *testing.T) {
	event, err := calendarUpdatePayload(map[string]any{"startedAt": "2025-04-11+09:00"})
	if err != nil {
		t.Fatalf("payload: %v", err)
	}
	if event.StartedAt == nil || event.StartedAt.String() != "2025-04-11+09:00" {
		t.Errorf("date-only start not preserved: %#v", event.StartedAt)
	}
}

func TestCalendarUpdateEventToolSubjectOnlyOmitsOtherFields(t *testing.T) {
	s := mcptest.NewServer()
	token := "test-token"
	called := 0
	calendarMutationTools(s, &token, func(_ context.Context, _, _, _ string, event model.UpdateEventRequest) (*model.EventResponse, error) {
		called++
		if event.Subject == nil || *event.Subject != "updated" {
			t.Errorf("subject: %#v", event.Subject)
		}
		if event.Body != nil || event.Location != nil || event.StartedAt != nil || event.EndedAt != nil || event.WholeDayFlag != nil || event.Users != nil || event.RecurrenceRule != nil {
			t.Errorf("partial update sent extra fields: %#v", event)
		}
		return successfulEventResponse(), nil
	}, refusingDelete(t))
	result := callCalendarTool(t, s, "dooray_calendar_update_event", updateEventArguments())
	if result.IsError || called != 1 {
		t.Fatalf("unexpected result: %+v calls=%d", result, called)
	}
}

func TestCalendarUpdateEventTool(t *testing.T) {
	s := mcptest.NewServer()
	token := "test-token"
	ctx := t.Context()
	called := 0
	calendarMutationTools(s, &token, func(gotCtx context.Context, gotToken, calendarID, eventID string, event model.UpdateEventRequest) (*model.EventResponse, error) {
		called++
		if gotCtx != ctx || gotToken != token || calendarID != "cal-123" || eventID != "evt-456" {
			t.Errorf("incorrect request context/ids")
		}
		if event.Subject == nil || *event.Subject != "updated" || event.Body == nil || event.Body.Content != "<p>body</p>" || event.Body.MimeType != "text/html" {
			t.Errorf("incorrect update payload: %#v", event)
		}
		return successfulEventResponse(), nil
	}, refusingDelete(t))
	args := updateEventArguments()
	args["content"] = "<p>body</p>"
	result := callCalendarTool(t, s, "dooray_calendar_update_event", args)
	if result.IsError || result.Content[0].(mcp.TextContent).Text != calendarSuccessBody || called != 1 {
		t.Fatalf("unexpected result: %+v calls=%d", result, called)
	}
}

func TestCalendarDeleteEventTool(t *testing.T) {
	s := mcptest.NewServer()
	token := "test-token"
	ctx := t.Context()
	called := 0
	occurrenceID := "3748560866668749754-20240228T013000Z"
	calendarMutationTools(s, &token, refusingUpdate(t), func(gotCtx context.Context, gotToken, calendarID, eventID string, req model.DeleteEventRequest) (*model.EventResponse, error) {
		called++
		if gotCtx != ctx || gotToken != token || calendarID != "cal-123" || eventID != occurrenceID {
			t.Errorf("incorrect delete request")
		}
		if req.DeleteType != model.DeleteTypeWhole {
			t.Errorf("deleteType: %#v", req.DeleteType)
		}
		return successfulEventResponse(), nil
	})
	args := deleteEventArguments()
	args["eventId"] = occurrenceID
	args["deleteType"] = "whole"
	result := callCalendarTool(t, s, "dooray_calendar_delete_event", args)
	if result.IsError || result.Content[0].(mcp.TextContent).Text != calendarSuccessBody || called != 1 {
		t.Fatalf("unexpected result: %+v calls=%d", result, called)
	}
}

func TestCalendarUpdateEventRejectsInvalidArgumentsWithoutCallingAPI(t *testing.T) {
	cases := []struct {
		name, key, message string
		value              any
		remove             bool
	}{
		{name: "missing operation", key: "operation", message: "operation", remove: true},
		{name: "missing calendar", key: "calendarId", message: "calendarId", remove: true},
		{name: "missing event", key: "eventId", message: "eventId", remove: true},
		{name: "wrong operation", key: "operation", value: "create_event", message: "operation"},
		{name: "blank subject", key: "subject", value: "  ", message: "subject"},
		{name: "empty content", key: "content", value: "", message: "content"},
		{name: "bad type", key: "location", value: 12, message: "location"},
		{name: "path injection", key: "calendarId", value: "1/../2", message: "calendarId"},
		{name: "comma calendar", key: "calendarId", value: "1,2", message: "calendarId"},
		{name: "wildcard event", key: "eventId", value: "*", message: "eventId"},
		{name: "bad time", key: "startedAt", value: "not-a-date", message: "startedAt"},
		{name: "bad flag", key: "wholeDayFlag", value: "true", message: "wholeDayFlag"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := mcptest.NewServer()
			token := "test-token"
			calendarMutationTools(s, &token, refusingUpdate(t), refusingDelete(t))
			args := updateEventArguments()
			if tc.remove {
				delete(args, tc.key)
			} else {
				args[tc.key] = tc.value
			}
			result := callCalendarTool(t, s, "dooray_calendar_update_event", args)
			if !result.IsError || !strings.Contains(toolErrorText(t, result), tc.message) {
				t.Fatalf("expected %q error: %+v", tc.message, result)
			}
		})
	}
}

func TestCalendarUpdateEventRequiresAField(t *testing.T) {
	s := mcptest.NewServer()
	token := "test-token"
	calendarMutationTools(s, &token, refusingUpdate(t), refusingDelete(t))
	args := map[string]any{"operation": "update_event", "calendarId": "cal-123", "eventId": "evt-456"}
	result := callCalendarTool(t, s, "dooray_calendar_update_event", args)
	if !strings.Contains(toolErrorText(t, result), "at least one field") {
		t.Fatalf("expected missing-field error: %+v", result)
	}
}

func TestCalendarDeleteEventRejectsInvalidArgumentsWithoutCallingAPI(t *testing.T) {
	cases := []struct {
		name, key, message string
		value              any
		remove             bool
	}{
		{name: "missing deleteType", key: "deleteType", message: "deleteType", remove: true},
		{name: "unknown deleteType", key: "deleteType", value: "all", message: "deleteType"},
		{name: "path injection", key: "eventId", value: "evt/../x", message: "eventId"},
		{name: "wrong operation", key: "operation", value: "update_event", message: "operation"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := mcptest.NewServer()
			token := "test-token"
			calendarMutationTools(s, &token, refusingUpdate(t), refusingDelete(t))
			args := deleteEventArguments()
			if tc.remove {
				delete(args, tc.key)
			} else {
				args[tc.key] = tc.value
			}
			result := callCalendarTool(t, s, "dooray_calendar_delete_event", args)
			if !result.IsError || !strings.Contains(toolErrorText(t, result), tc.message) {
				t.Fatalf("expected %q error: %+v", tc.message, result)
			}
		})
	}
}

func TestCalendarMutationToolsRequireToken(t *testing.T) {
	for _, token := range []*string{nil, new(""), new(" ")} {
		s := mcptest.NewServer()
		calendarMutationTools(s, token, refusingUpdate(t), refusingDelete(t))
		for _, tc := range []struct {
			name string
			args map[string]any
		}{
			{"update", updateEventArguments()},
			{"delete", deleteEventArguments()},
		} {
			result := callCalendarTool(t, s, "dooray_calendar_"+tc.name+"_event", tc.args)
			if !strings.Contains(toolErrorText(t, result), "token") {
				t.Fatalf("%s expected missing token error: %+v", tc.name, result)
			}
		}
	}
}

func TestCalendarMutationToolsReportAPIFailures(t *testing.T) {
	rejected := &model.EventResponse{}
	rejected.Header.ResultMessage = "permission denied"
	cases := []struct {
		name     string
		response *model.EventResponse
		err      error
		message  string
	}{
		{"transport", nil, errors.New("network unavailable"), "network unavailable"},
		{"rejected", rejected, nil, "permission denied"},
		{"empty response", nil, nil, "empty response"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := mcptest.NewServer()
			token := "test-token"
			fail := func(context.Context, string, string, string, model.UpdateEventRequest) (*model.EventResponse, error) {
				return tc.response, tc.err
			}
			failDelete := func(context.Context, string, string, string, model.DeleteEventRequest) (*model.EventResponse, error) {
				return tc.response, tc.err
			}
			calendarMutationTools(s, &token, fail, failDelete)
			for _, name := range []string{"dooray_calendar_update_event", "dooray_calendar_delete_event"} {
				args := updateEventArguments()
				if name == "dooray_calendar_delete_event" {
					args = deleteEventArguments()
				}
				result := callCalendarTool(t, s, name, args)
				if !strings.Contains(toolErrorText(t, result), tc.message) {
					t.Fatalf("%s expected %q: %+v", name, tc.message, result)
				}
			}
		})
	}
}

func TestCalendarIDValid(t *testing.T) {
	valid := []string{"cal-123", "evt_2026-09-15", "3748560866668749754-20240228T013000Z"}
	invalid := []string{"", " ", "1,2", "1/2", "1\\2", "a?b", "a#b", "%2f", ".", "..", "*", "id with space"}
	for _, id := range valid {
		if !calendarIDValid(id) {
			t.Errorf("%q should be valid", id)
		}
	}
	for _, id := range invalid {
		if calendarIDValid(id) {
			t.Errorf("%q should be invalid", id)
		}
	}
}

func updateEventArguments() map[string]any {
	return map[string]any{"operation": "update_event", "calendarId": "cal-123", "eventId": "evt-456", "subject": "updated"}
}

func deleteEventArguments() map[string]any {
	return map[string]any{"operation": "delete_event", "calendarId": "cal-123", "eventId": "evt-456", "deleteType": "this"}
}

func callCalendarTool(t *testing.T, s *server.MCPServer, name string, args map[string]any) *mcp.CallToolResult {
	t.Helper()
	tool, ok := s.ListTools()[name]
	if !ok {
		t.Fatalf("%s not registered", name)
	}
	result, err := tool.Handler(t.Context(), mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}})
	if err != nil || result == nil {
		t.Fatalf("%s handler failed: %+v %v", name, result, err)
	}
	return result
}

func toolErrorText(t *testing.T, result *mcp.CallToolResult) string {
	t.Helper()
	if !result.IsError {
		t.Fatalf("expected tool error, got %+v", result)
	}
	return result.Content[0].(mcp.TextContent).Text
}

func successfulEventResponse() *model.EventResponse {
	res := &model.EventResponse{RawJSON: calendarSuccessBody}
	res.Header.IsSuccessful = true
	return res
}

func refusingUpdate(t *testing.T) calendarUpdateFunc {
	t.Helper()
	return func(context.Context, string, string, string, model.UpdateEventRequest) (*model.EventResponse, error) {
		t.Fatal("update API must not be called")
		return nil, nil
	}
}

func refusingDelete(t *testing.T) calendarDeleteFunc {
	t.Helper()
	return func(context.Context, string, string, string, model.DeleteEventRequest) (*model.EventResponse, error) {
		t.Fatal("delete API must not be called")
		return nil, nil
	}
}
