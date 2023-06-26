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
})

var commandPayload = Type("Command", func() {
	Description("https://api.slack.com/interactivity/slash-commands")
	Attribute("signature", String, func() {
		Example("v0=somethingveryimportant")
	})
	Attribute("timestamp", Int64, func() {
		Example(1688832649)
	})
	Attribute("token", String, func() {
		Example("gIkuvaNzQIHg97ATvDxqgjtO")
	})
	Attribute("command", String, func() {
		Example("/weather")
	})
	Attribute("text", String, func() {
		Example("94070")
	})
	Attribute("response_url", String, func() {
		Example("https://hooks.slack.com/commands/1234/5678")
	})
	Attribute("trigger_id", String, func() {
		Example("13345224609.738474920.8088930838d88f008e0")
	})
	Attribute("user_id", String, func() {
		Example("U2147483697")
	})
	Attribute("user_name", String, func() {
		Example("Steve")
	})
	Attribute("team_id", String, func() {
		Example("T0001")
	})
	Attribute("team_domain", String, func() {
		Example("example")
	})
	Attribute("enterprise_id", String, func() {
		Example("E0001")
	})
	Attribute("enterprise_name", String, func() {
		Example("Globular Construct Inc")
	})
	Attribute("is_enterprise_install", Boolean, func() {
		Example(true)
	})
	Attribute("channel_id", String, func() {
		Example("C2147483705")
	})
	Attribute("channel_name", String, func() {
		Example("test")
	})
	Attribute("api_app_id", String, func() {
		Example("A123456")
	})
})

var eventPayload = Type("Event", func() {
	Description("https://api.slack.com/apis/connections/events-api")
	Attribute("signature", String, func() {
		Example("v0=somethingveryimportant")
	})
	Attribute("timestamp", Int64, func() {
		Example(1688832649)
	})
	Attribute("token", String, func() {
		Example("gIkuvaNzQIHg97ATvDxqgjtO")
	})
	Attribute("team_id", String, func() {
		Example("T0001")
	})
	Attribute("challenge", String, func() {
		Example("randomstring")
	})
	Attribute("type", String, func() {
		Example("event_callback")
	})
	Attribute("api_app_id", String, func() {
		Example("A123456")
	})
	Attribute("event", func() {
		Description("The actual event information")
		Attribute("type", String)
	})
})

var eventResponse = Type("EventResponse", func() {
	Description("https://api.slack.com/apis/connections/events-api")
	Attribute("challenge", String, func() {
		Example("randomstring")
	})
})
