package design

import . "goa.design/goa/v3/dsl"

var SlackService = Service("Slack", func() {
	Description("Slack api for interacting with slack commands, actions, events etc.")

	HTTP(func() {
		Path("/slack")
	})

	Method("Commands", func() {
		Payload(commandPayload)

		HTTP(func() {
			POST("commands")
			Header("signature:X-Slack-Signature")
			Header("timestamp:X-Slack-Request-Timestamp")
			Response(StatusOK)
		})
	})

	Method("Events", func() {
		Payload(eventPayload)
		Result(eventResponse)

		HTTP(func() {
			POST("events")
			Header("signature:X-Slack-Signature")
			Header("timestamp:X-Slack-Request-Timestamp")
			Response(StatusOK)
		})
	})

	Method("MessageActions", func() {
		Payload(actionPayload)
		Result(actionResponse)

		HTTP(func() {
			POST("message_actions")
			Header("signature:X-Slack-Signature")
			Header("timestamp:X-Slack-Request-Timestamp")
			Response(StatusOK)
		})
	})
})

var commandPayload = Type("Command", func() {
	Description("https://api.slack.com/interactivity/slash-commands")
	Attribute("signature", String)
	Attribute("timestamp", Int64)
	Attribute("token", String)
	Attribute("command", String)
	Attribute("text", String)
	Attribute("response_url", String)
	Attribute("trigger_id", String)
	Attribute("user_id", String)
	Attribute("user_name", String)
	Attribute("team_id", String)
	Attribute("team_domain", String)
	Attribute("enterprise_id", String)
	Attribute("enterprise_name", String)
	Attribute("is_enterprise_install", Boolean)
	Attribute("channel_id", String)
	Attribute("channel_name", String)
	Attribute("api_app_id", String)
	Required(
		"signature",
		"timestamp",
		"token",
		"command",
		"trigger_id",
		"user_id",
		"team_id",
		"channel_id",
	)
})

var eventPayload = Type("Event", func() {
	Description("https://api.slack.com/apis/connections/events-api")
	Attribute("signature", String)
	Attribute("timestamp", Int64)
	Attribute("token", String)
	Attribute("team_id", String)
	Attribute("challenge", String)
	Attribute("type", String)
	Attribute("api_app_id", String)
	Attribute("event", func() {
		Description("The actual event information")
		Attribute("type", String)
	})
	Required(
		"signature",
		"timestamp",
		"token",
		"team_id",
		"type",
		"api_app_id",
	)
})

var eventResponse = Type("EventResponse", func() {
	Description("https://api.slack.com/apis/connections/events-api")
	Attribute("challenge", String, func() {
		Example("randomstring")
	})
})

var actionPayload = Type("Action", func() {
	Description("https://api.slack.com/reference/interaction-payloads")
	Attribute("signature", String)
	Attribute("timestamp", Int64)
	Attribute("payload", Bytes)
	Required("signature", "timestamp", "payload")
})

var actionResponse = Type("ActionResponse", func() {
	Description("https://api.slack.com/surfaces/modals#displaying_errors")
	Attribute("response_action", String, func() {
		Example("errors")
	})
	Attribute("view", Any)
	Attribute("errors", MapOf(String, String), func() {
		Example(map[string]string{"foo": "bar"})
	})
})
