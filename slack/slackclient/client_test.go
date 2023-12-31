package slackclient

import (
	"context"

	"github.com/rotabot-io/rotabot/slack/slackclient/mock_slackclient"

	"github.com/slack-go/slack"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ClientFor", func() {
	var (
		ctx context.Context
		sc  *mock_slackclient.MockSlackClient
	)

	BeforeEach(func() {
		ctx = context.Background()
	})

	// Create a mock and assign it to the sc variable at the start of each test
	MockSlackClient(&ctx, &sc, nil)

	Describe("mocking a channel call", func() {
		// Apply an expectation in the BeforeEach, before the test runs
		BeforeEach(func() {
			sc.EXPECT().GetConversationInfoContext(ctx, &slack.GetConversationInfoInput{
				ChannelID:         "CH123",
				IncludeLocale:     false,
				IncludeNumMembers: false,
			}).
				Return(&slack.Channel{
					GroupConversation: slack.GroupConversation{
						Name: "my-channel",
						Conversation: slack.Conversation{
							NameNormalized: "my-channel",
							ID:             "CH123",
						},
					},
				}, nil).Times(1)
		})

		Specify("returns a client that responds with the mock", func() {
			client, err := ClientFor(ctx, "OR123")
			Expect(err).NotTo(HaveOccurred(), "Slack client should have built with no error")

			channel, err := client.GetConversationInfoContext(ctx, &slack.GetConversationInfoInput{
				ChannelID:         "CH123",
				IncludeLocale:     false,
				IncludeNumMembers: false,
			})
			Expect(err).NotTo(HaveOccurred())

			// We'll only receive this if the client generated by ClientFor is the mock we
			// configured with a fake response in our BeforeEach.
			Expect(channel.NameNormalized).To(Equal("my-channel"))
		})
	})
})
