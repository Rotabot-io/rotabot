package slack

import (
	"context"

	"github.com/slack-go/slack"

	"github.com/jackc/pgx/v5"
	"github.com/rotabot-io/rotabot/internal"

	"github.com/rotabot-io/rotabot/lib/db"

	"github.com/slack-go/slack/slackevents"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gen "github.com/rotabot-io/rotabot/gen/slack"
	"github.com/rotabot-io/rotabot/slack/slackclient"
	"github.com/rotabot-io/rotabot/slack/slackclient/mock_slackclient"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Service", func() {
	var (
		ctx context.Context
		sc  *mock_slackclient.MockSlackClient
		svc gen.Service
	)

	BeforeEach(func() {
		ctx = context.Background()
	})

	// Create a mock and assign it to the sc variable at the start of each test
	slackclient.MockSlackClient(&ctx, &sc, nil)

	Describe("Command", func() {
		BeforeEach(func() {
			container, err := internal.RunContainer(ctx)
			Expect(err).ToNot(HaveOccurred())

			err = db.Migrate(ctx, container.ConnectionString())
			Expect(err).ToNot(HaveOccurred())

			conn, err := pgx.Connect(ctx, container.ConnectionString())
			Expect(err).ToNot(HaveOccurred())

			DeferCleanup(func() {
				conn.Close(ctx)
			})

			svc = New(db.New(conn))
		})
		It("should open view when command api is called", func() {
			sc.EXPECT().OpenViewContext(gomock.Any(), "T123", gomock.Any()).Return(nil, nil).Times(1)

			err := svc.Commands(ctx, &gen.Command{
				Signature: "TEST",
				Timestamp: 1234567890,
				TriggerID: "T123",
				Command:   "/rotabot",
				Token:     "TEST",
				TeamID:    "TE123",
				ChannelID: "CH123",
			})

			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("Events", func() {
		BeforeEach(func() {
			svc = New(&db.Queries{})
		})
		It("Should return the challenge", func() {
			challenge := "TEST"
			res, err := svc.Events(ctx, &gen.Event{
				Signature: "TEST",
				Timestamp: 1234567890,
				Token:     "TEST",
				Type:      slackevents.URLVerification,
				TeamID:    "T123",
				Challenge: &challenge,
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(res).ToNot(BeNil())

			responseChallenge := *res.Challenge
			Expect(responseChallenge).ToNot(BeEmpty())
			Expect(responseChallenge).To(Equal(challenge))
		})

		It("Should return the challenge", func() {
			res, err := svc.Events(ctx, &gen.Event{
				Signature: "TEST",
				Timestamp: 1234567890,
				Token:     "TEST",
				Type:      string(slackevents.AppMention),
				TeamID:    "T123",
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(res).ToNot(BeNil())
			Expect(res.Challenge).To(BeNil())
		})
	})

	Describe("Actions", func() {
		BeforeEach(func() {
			svc = New(&db.Queries{})
		})

		It("Should return response without errors", func() {
			payload := "{\"type\":\"view_closed\",\"team\":{\"id\":\"T042E16BURW\",\"domain\":\"rotabot-workspace\"},\"user\":{\"id\":\"U0422UEJLD7\",\"username\":\"me1\",\"name\":\"me1\",\"team_id\":\"T042E16BURW\"},\"api_app_id\":\"A041MF0T137\",\"token\":\"dasdasdasdasdas\",\"view\":{\"id\":\"V05KE4NJCDR\",\"team_id\":\"T042E16BURW\",\"type\":\"modal\",\"blocks\":[{\"type\":\"actions\",\"block_id\":\"HOME_ACTIONS\",\"elements\":[{\"type\":\"button\",\"action_id\":\"HOME_ADD_ROTA\",\"text\":{\"type\":\"plain_text\",\"text\":\"Add Rota :heavy_plus_sign:\",\"emoji\":true}}]},{\"type\":\"header\",\"block_id\":\"R4Ay\",\"text\":{\"type\":\"plain_text\",\"text\":\"Active Rotas:\",\"emoji\":true}},{\"type\":\"section\",\"block_id\":\"ROTA_ELEMENT_RTGZBrAeq3dDWYvK\",\"text\":{\"type\":\"plain_text\",\"text\":\"gdgdfgdf\",\"emoji\":true},\"accessory\":{\"type\":\"overflow\",\"action_id\":\"ROTA_ELEMENT_RTGZBrAeq3dDWYvK\",\"options\":[{\"text\":{\"type\":\"plain_text\",\"text\":\"Example1\",\"emoji\":true},\"value\":\"ExampleAction1_RTGZBrAeq3dDWYvK\"},{\"text\":{\"type\":\"plain_text\",\"text\":\"Example2\",\"emoji\":true},\"value\":\"ExampleAction2_RTGZBrAeq3dDWYvK\"},{\"text\":{\"type\":\"plain_text\",\"text\":\"Example3\",\"emoji\":true},\"value\":\"ExampleAction3_RTGZBrAeq3dDWYvK\"},{\"text\":{\"type\":\"plain_text\",\"text\":\"Example4\",\"emoji\":true},\"value\":\"ExampleAction4_RTGZBrAeq3dDWYvK\"},{\"text\":{\"type\":\"plain_text\",\"text\":\"Example5\",\"emoji\":true},\"value\":\"ExampleAction5_RTGZBrAeq3dDWYvK\"}]}},{\"type\":\"section\",\"block_id\":\"ROTA_ELEMENT_RT9fTe6m7wITGEvK\",\"text\":{\"type\":\"plain_text\",\"text\":\"dasfafs\",\"emoji\":true},\"accessory\":{\"type\":\"overflow\",\"action_id\":\"ROTA_ELEMENT_RT9fTe6m7wITGEvK\",\"options\":[{\"text\":{\"type\":\"plain_text\",\"text\":\"Example1\",\"emoji\":true},\"value\":\"ExampleAction1_RT9fTe6m7wITGEvK\"},{\"text\":{\"type\":\"plain_text\",\"text\":\"Example2\",\"emoji\":true},\"value\":\"ExampleAction2_RT9fTe6m7wITGEvK\"},{\"text\":{\"type\":\"plain_text\",\"text\":\"Example3\",\"emoji\":true},\"value\":\"ExampleAction3_RT9fTe6m7wITGEvK\"},{\"text\":{\"type\":\"plain_text\",\"text\":\"Example4\",\"emoji\":true},\"value\":\"ExampleAction4_RT9fTe6m7wITGEvK\"},{\"text\":{\"type\":\"plain_text\",\"text\":\"Example5\",\"emoji\":true},\"value\":\"ExampleAction5_RT9fTe6m7wITGEvK\"}]}}],\"private_metadata\":\"C041Q6Z5FSP\",\"callback_id\":\"Home\",\"state\":{\"values\":{}},\"hash\":\"1690721765.rDNKI2V9\",\"title\":{\"type\":\"plain_text\",\"text\":\"Rotabot Home\",\"emoji\":true},\"clear_on_close\":true,\"notify_on_close\":true,\"close\":null,\"submit\":null,\"previous_view_id\":null,\"root_view_id\":\"V05KE4NJCDR\",\"app_id\":\"A041MF0T137\",\"external_id\":\"\",\"app_installed_team_id\":\"T042E16BURW\",\"bot_id\":\"B041A53D6ET\"},\"is_cleared\":true,\"is_enterprise_install\":false,\"enterprise\":null}"
			res, err := svc.MessageActions(ctx, &gen.Action{
				Signature: "TEST",
				Timestamp: 1234567890,
				Payload:   []byte(payload),
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(res).ToNot(BeNil())

			Expect(res.ResponseAction).To(BeNil())
		})

		It("Should return clear when the view is unknown", func() {
			payload := "{\"type\":\"ddddddd\",\"team\":{\"id\":\"T042E16BURW\",\"domain\":\"rotabot-workspace\"},\"user\":{\"id\":\"U0422UEJLD7\",\"username\":\"me1\",\"name\":\"me1\",\"team_id\":\"T042E16BURW\"},\"api_app_id\":\"A041MF0T137\",\"token\":\"dasdasdasdasdas\",\"view\":{\"id\":\"V05KE4NJCDR\",\"team_id\":\"T042E16BURW\",\"type\":\"modal\",\"blocks\":[{\"type\":\"actions\",\"block_id\":\"HOME_ACTIONS\",\"elements\":[{\"type\":\"button\",\"action_id\":\"HOME_ADD_ROTA\",\"text\":{\"type\":\"plain_text\",\"text\":\"Add Rota :heavy_plus_sign:\",\"emoji\":true}}]},{\"type\":\"header\",\"block_id\":\"R4Ay\",\"text\":{\"type\":\"plain_text\",\"text\":\"Active Rotas:\",\"emoji\":true}},{\"type\":\"section\",\"block_id\":\"ROTA_ELEMENT_RTGZBrAeq3dDWYvK\",\"text\":{\"type\":\"plain_text\",\"text\":\"gdgdfgdf\",\"emoji\":true},\"accessory\":{\"type\":\"overflow\",\"action_id\":\"ROTA_ELEMENT_RTGZBrAeq3dDWYvK\",\"options\":[{\"text\":{\"type\":\"plain_text\",\"text\":\"Example1\",\"emoji\":true},\"value\":\"ExampleAction1_RTGZBrAeq3dDWYvK\"},{\"text\":{\"type\":\"plain_text\",\"text\":\"Example2\",\"emoji\":true},\"value\":\"ExampleAction2_RTGZBrAeq3dDWYvK\"},{\"text\":{\"type\":\"plain_text\",\"text\":\"Example3\",\"emoji\":true},\"value\":\"ExampleAction3_RTGZBrAeq3dDWYvK\"},{\"text\":{\"type\":\"plain_text\",\"text\":\"Example4\",\"emoji\":true},\"value\":\"ExampleAction4_RTGZBrAeq3dDWYvK\"},{\"text\":{\"type\":\"plain_text\",\"text\":\"Example5\",\"emoji\":true},\"value\":\"ExampleAction5_RTGZBrAeq3dDWYvK\"}]}},{\"type\":\"section\",\"block_id\":\"ROTA_ELEMENT_RT9fTe6m7wITGEvK\",\"text\":{\"type\":\"plain_text\",\"text\":\"dasfafs\",\"emoji\":true},\"accessory\":{\"type\":\"overflow\",\"action_id\":\"ROTA_ELEMENT_RT9fTe6m7wITGEvK\",\"options\":[{\"text\":{\"type\":\"plain_text\",\"text\":\"Example1\",\"emoji\":true},\"value\":\"ExampleAction1_RT9fTe6m7wITGEvK\"},{\"text\":{\"type\":\"plain_text\",\"text\":\"Example2\",\"emoji\":true},\"value\":\"ExampleAction2_RT9fTe6m7wITGEvK\"},{\"text\":{\"type\":\"plain_text\",\"text\":\"Example3\",\"emoji\":true},\"value\":\"ExampleAction3_RT9fTe6m7wITGEvK\"},{\"text\":{\"type\":\"plain_text\",\"text\":\"Example4\",\"emoji\":true},\"value\":\"ExampleAction4_RT9fTe6m7wITGEvK\"},{\"text\":{\"type\":\"plain_text\",\"text\":\"Example5\",\"emoji\":true},\"value\":\"ExampleAction5_RT9fTe6m7wITGEvK\"}]}}],\"private_metadata\":\"C041Q6Z5FSP\",\"callback_id\":\"Home\",\"state\":{\"values\":{}},\"hash\":\"1690721765.rDNKI2V9\",\"title\":{\"type\":\"plain_text\",\"text\":\"Rotabot Home\",\"emoji\":true},\"clear_on_close\":true,\"notify_on_close\":true,\"close\":null,\"submit\":null,\"previous_view_id\":null,\"root_view_id\":\"V05KE4NJCDR\",\"app_id\":\"A041MF0T137\",\"external_id\":\"\",\"app_installed_team_id\":\"T042E16BURW\",\"bot_id\":\"B041A53D6ET\"},\"is_cleared\":true,\"is_enterprise_install\":false,\"enterprise\":null}"
			res, err := svc.MessageActions(ctx, &gen.Action{
				Signature: "TEST",
				Timestamp: 1234567890,
				Payload:   []byte(payload),
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(res).ToNot(BeNil())

			clear := string(slack.RAClear)
			Expect(res.ResponseAction).To(Equal(&clear))
		})
	})
})