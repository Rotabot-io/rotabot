// Code generated by goa v3.14.0, DO NOT EDIT.
//
// Slack HTTP client types
//
// Command:
// $ goa gen github.com/rotabot-io/rotabot/design

package client

import (
	slack "github.com/rotabot-io/rotabot/gen/slack"
)

// CommandsRequestBody is the type of the "Slack" service "Commands" endpoint
// HTTP request body.
type CommandsRequestBody struct {
	Token               string  `form:"token" json:"token" xml:"token"`
	Command             string  `form:"command" json:"command" xml:"command"`
	Text                *string `form:"text,omitempty" json:"text,omitempty" xml:"text,omitempty"`
	ResponseURL         *string `form:"response_url,omitempty" json:"response_url,omitempty" xml:"response_url,omitempty"`
	TriggerID           string  `form:"trigger_id" json:"trigger_id" xml:"trigger_id"`
	UserID              string  `form:"user_id" json:"user_id" xml:"user_id"`
	UserName            *string `form:"user_name,omitempty" json:"user_name,omitempty" xml:"user_name,omitempty"`
	TeamID              string  `form:"team_id" json:"team_id" xml:"team_id"`
	TeamDomain          *string `form:"team_domain,omitempty" json:"team_domain,omitempty" xml:"team_domain,omitempty"`
	EnterpriseID        *string `form:"enterprise_id,omitempty" json:"enterprise_id,omitempty" xml:"enterprise_id,omitempty"`
	EnterpriseName      *string `form:"enterprise_name,omitempty" json:"enterprise_name,omitempty" xml:"enterprise_name,omitempty"`
	IsEnterpriseInstall *bool   `form:"is_enterprise_install,omitempty" json:"is_enterprise_install,omitempty" xml:"is_enterprise_install,omitempty"`
	ChannelID           string  `form:"channel_id" json:"channel_id" xml:"channel_id"`
	ChannelName         *string `form:"channel_name,omitempty" json:"channel_name,omitempty" xml:"channel_name,omitempty"`
	APIAppID            *string `form:"api_app_id,omitempty" json:"api_app_id,omitempty" xml:"api_app_id,omitempty"`
}

// EventsRequestBody is the type of the "Slack" service "Events" endpoint HTTP
// request body.
type EventsRequestBody struct {
	Token     string  `form:"token" json:"token" xml:"token"`
	TeamID    string  `form:"team_id" json:"team_id" xml:"team_id"`
	Challenge *string `form:"challenge,omitempty" json:"challenge,omitempty" xml:"challenge,omitempty"`
	Type      string  `form:"type" json:"type" xml:"type"`
	APIAppID  string  `form:"api_app_id" json:"api_app_id" xml:"api_app_id"`
	// The actual event information
	Event *struct {
		Type *string `form:"type" json:"type" xml:"type"`
	} `form:"event,omitempty" json:"event,omitempty" xml:"event,omitempty"`
}

// MessageActionsRequestBody is the type of the "Slack" service
// "MessageActions" endpoint HTTP request body.
type MessageActionsRequestBody struct {
	Payload []byte `form:"payload" json:"payload" xml:"payload"`
}

// EventsResponseBody is the type of the "Slack" service "Events" endpoint HTTP
// response body.
type EventsResponseBody struct {
	Challenge *string `form:"challenge,omitempty" json:"challenge,omitempty" xml:"challenge,omitempty"`
}

// MessageActionsResponseBody is the type of the "Slack" service
// "MessageActions" endpoint HTTP response body.
type MessageActionsResponseBody struct {
	ResponseAction *string           `form:"response_action,omitempty" json:"response_action,omitempty" xml:"response_action,omitempty"`
	View           any               `form:"view,omitempty" json:"view,omitempty" xml:"view,omitempty"`
	Errors         map[string]string `form:"errors,omitempty" json:"errors,omitempty" xml:"errors,omitempty"`
}

// NewCommandsRequestBody builds the HTTP request body from the payload of the
// "Commands" endpoint of the "Slack" service.
func NewCommandsRequestBody(p *slack.Command) *CommandsRequestBody {
	body := &CommandsRequestBody{
		Token:               p.Token,
		Command:             p.Command,
		Text:                p.Text,
		ResponseURL:         p.ResponseURL,
		TriggerID:           p.TriggerID,
		UserID:              p.UserID,
		UserName:            p.UserName,
		TeamID:              p.TeamID,
		TeamDomain:          p.TeamDomain,
		EnterpriseID:        p.EnterpriseID,
		EnterpriseName:      p.EnterpriseName,
		IsEnterpriseInstall: p.IsEnterpriseInstall,
		ChannelID:           p.ChannelID,
		ChannelName:         p.ChannelName,
		APIAppID:            p.APIAppID,
	}
	return body
}

// NewEventsRequestBody builds the HTTP request body from the payload of the
// "Events" endpoint of the "Slack" service.
func NewEventsRequestBody(p *slack.Event) *EventsRequestBody {
	body := &EventsRequestBody{
		Token:     p.Token,
		TeamID:    p.TeamID,
		Challenge: p.Challenge,
		Type:      p.Type,
		APIAppID:  p.APIAppID,
	}
	if p.Event != nil {
		body.Event = &struct {
			Type *string `form:"type" json:"type" xml:"type"`
		}{
			Type: p.Event.Type,
		}
	}
	return body
}

// NewMessageActionsRequestBody builds the HTTP request body from the payload
// of the "MessageActions" endpoint of the "Slack" service.
func NewMessageActionsRequestBody(p *slack.Action) *MessageActionsRequestBody {
	body := &MessageActionsRequestBody{
		Payload: p.Payload,
	}
	return body
}

// NewEventsEventResponseOK builds a "Slack" service "Events" endpoint result
// from a HTTP "OK" response.
func NewEventsEventResponseOK(body *EventsResponseBody) *slack.EventResponse {
	v := &slack.EventResponse{
		Challenge: body.Challenge,
	}

	return v
}

// NewMessageActionsActionResponseOK builds a "Slack" service "MessageActions"
// endpoint result from a HTTP "OK" response.
func NewMessageActionsActionResponseOK(body *MessageActionsResponseBody) *slack.ActionResponse {
	v := &slack.ActionResponse{
		ResponseAction: body.ResponseAction,
		View:           body.View,
	}
	if body.Errors != nil {
		v.Errors = make(map[string]string, len(body.Errors))
		for key, val := range body.Errors {
			tk := key
			tv := val
			v.Errors[tk] = tv
		}
	}

	return v
}
