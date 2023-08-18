// Code generated by goa v3.12.4, DO NOT EDIT.
//
// Slack service
//
// Command:
// $ goa gen github.com/rotabot-io/rotabot/design

package slack

import (
	"context"
)

// Slack api for interacting with slack commands, actions, events etc.
type Service interface {
	// Commands implements Commands.
	Commands(context.Context, *Command) (err error)
	// Events implements Events.
	Events(context.Context, *Event) (res *EventResponse, err error)
	// MessageActions implements MessageActions.
	MessageActions(context.Context, *Action) (res *ActionResponse, err error)
}

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "Slack"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [3]string{"Commands", "Events", "MessageActions"}

// Action is the payload type of the Slack service MessageActions method.
type Action struct {
	Signature string
	Timestamp int64
	Payload   []byte
}

// ActionResponse is the result type of the Slack service MessageActions method.
type ActionResponse struct {
	ResponseAction *string
	View           any
	Errors         map[string]string
}

// Command is the payload type of the Slack service Commands method.
type Command struct {
	Signature           string
	Timestamp           int64
	Token               string
	Command             string
	Text                *string
	ResponseURL         *string
	TriggerID           string
	UserID              string
	UserName            *string
	TeamID              string
	TeamDomain          *string
	EnterpriseID        *string
	EnterpriseName      *string
	IsEnterpriseInstall *bool
	ChannelID           string
	ChannelName         *string
	APIAppID            *string
}

// Event is the payload type of the Slack service Events method.
type Event struct {
	Signature string
	Timestamp int64
	Token     string
	TeamID    string
	Challenge *string
	Type      string
	APIAppID  string
	// The actual event information
	Event *struct {
		Type *string
	}
}

// EventResponse is the result type of the Slack service Events method.
type EventResponse struct {
	Challenge *string
}
