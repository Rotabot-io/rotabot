package design

import (
	slack "github.com/rotabot-io/rotabot/slack/design"
	. "goa.design/goa/v3/dsl"
)

var _ = slack.SlackService

var _ = API("Rotabot", func() {
	Title("Rotabot - Making rotas dead simple")
	Description("A service for working with rotas across multiples tools i.e Slack, Teams, etc")
	HTTP(func() {
		Produces("application/json")
		Consumes("application/json", "application/x-www-form-urlencoded")
	})

	Server("Rotabot", func() {
		Description("Backend for the rotabot application.")

		// List the services hosted by this server.
		Services("Slack")

		// List the Hosts and their transport URLs.
		Host("development", func() {
			Description("Development hosts.")
			// Transport specific URLs, supported schemes are:
			// 'http', 'https', 'grpc' and 'grpcs' with the respective default
			// ports: 80, 443, 8080, 8443.
			URI("http://localhost:8080/")
		})
	})
})
