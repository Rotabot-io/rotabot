package slack

import (
	"context"
	"time"

	"github.com/rotabot-io/rotabot/slack/slackclient"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/slack-go/slack"

	"github.com/rotabot-io/rotabot/lib/db"

	"github.com/slack-go/slack/slackevents"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gen "github.com/rotabot-io/rotabot/gen/slack"
	"github.com/rotabot-io/rotabot/slack/slackclient/mock_slackclient"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Service", func() {
	var (
		ctx        context.Context
		sc         *mock_slackclient.MockSlackClient
		svc        gen.Service
		connString string
		conn       *pgxpool.Pool
		container  *postgres.PostgresContainer
	)

	channelId := "C041Q6Z5FSP"
	teamId := "T042E16BURW"

	BeforeEach(func() {
		var err error
		ctx = context.Background()

		container, err = postgres.RunContainer(ctx,
			testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(5*time.Second)),
		)
		Expect(err).ToNot(HaveOccurred())

		connString, err = container.ConnectionString(ctx, "sslmode=disable")
		Expect(err).ToNot(HaveOccurred())

		err = db.Migrate(ctx, connString)
		Expect(err).ToNot(HaveOccurred())

		conn, err = pgxpool.New(ctx, connString)
		Expect(err).ToNot(HaveOccurred())

		svc = New(conn)

		DeferCleanup(func() {
			_ = container.Terminate(ctx)
			conn.Close()
		})
	})

	// Create a mock and assign it to the sc variable at the start of each test
	slackclient.MockSlackClient(&ctx, &sc, nil)

	Describe("Command", func() {
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
		It("Should save the rota that was requested by the user", func() {
			sc.EXPECT().UpdateViewContext(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)

			payload := "{\"type\":\"view_submission\",\"team\":{\"id\":\"T042E16BURW\",\"domain\":\"rotabot-workspace\"},\"user\":{\"id\":\"U0422UEJLD7\",\"username\":\"me1\",\"name\":\"me1\",\"team_id\":\"T042E16BURW\"},\"api_app_id\":\"A041MF0T137\",\"token\":\"XXXXXXXXXXXXXX\",\"trigger_id\":\"5885103319953.4082040402880.46a426db384926468ec6e193e1165f94\",\"view\":{\"id\":\"V05RNDTE2NN\",\"team_id\":\"T042E16BURW\",\"type\":\"modal\",\"blocks\":[{\"type\":\"input\",\"block_id\":\"ROTA_NAME\",\"label\":{\"type\":\"plain_text\",\"text\":\"Name:\",\"emoji\":true},\"optional\":false,\"dispatch_action\":false,\"element\":{\"type\":\"plain_text_input\",\"action_id\":\"ROTA_NAME\",\"placeholder\":{\"type\":\"plain_text\",\"text\":\"e.g. 'On Call'\",\"emoji\":true},\"dispatch_action_config\":{\"trigger_actions_on\":[\"on_enter_pressed\"]}}},{\"type\":\"section\",\"block_id\":\"ROTA_FREQUENCY\",\"text\":{\"type\":\"plain_text\",\"text\":\"Frequency:\",\"emoji\":true},\"accessory\":{\"type\":\"static_select\",\"action_id\":\"ROTA_FREQUENCY\",\"initial_option\":{\"text\":{\"type\":\"plain_text\",\"text\":\"Weekly\",\"emoji\":true},\"value\":\"Weekly\"},\"options\":[{\"text\":{\"type\":\"plain_text\",\"text\":\"Daily\",\"emoji\":true},\"value\":\"Daily\"},{\"text\":{\"type\":\"plain_text\",\"text\":\"Weekly\",\"emoji\":true},\"value\":\"Weekly\"},{\"text\":{\"type\":\"plain_text\",\"text\":\"Monthly\",\"emoji\":true},\"value\":\"Monthly\"}]}},{\"type\":\"section\",\"block_id\":\"ROTA_TYPE\",\"text\":{\"type\":\"plain_text\",\"text\":\"Scheduling Type:\",\"emoji\":true},\"accessory\":{\"type\":\"static_select\",\"action_id\":\"ROTA_TYPE\",\"initial_option\":{\"text\":{\"type\":\"plain_text\",\"text\":\"Created At\",\"emoji\":true},\"value\":\"Created At\"},\"options\":[{\"text\":{\"type\":\"plain_text\",\"text\":\"Created At\",\"emoji\":true},\"value\":\"Created At\"},{\"text\":{\"type\":\"plain_text\",\"text\":\"Randomly\",\"emoji\":true},\"value\":\"Randomly\"}]}}],\"private_metadata\":\"{\\\"rota_id\\\":\\\"\\\",\\\"channel_id\\\":\\\"C041Q6Z5FSP\\\"}\",\"callback_id\":\"SaveRota\",\"state\":{\"values\":{\"ROTA_NAME\":{\"ROTA_NAME\":{\"type\":\"plain_text_input\",\"value\":\"RotaName\"}},\"ROTA_FREQUENCY\":{\"ROTA_FREQUENCY\":{\"type\":\"static_select\",\"selected_option\":{\"text\":{\"type\":\"plain_text\",\"text\":\"Weekly\",\"emoji\":true},\"value\":\"Weekly\"}}},\"ROTA_TYPE\":{\"ROTA_TYPE\":{\"type\":\"static_select\",\"selected_option\":{\"text\":{\"type\":\"plain_text\",\"text\":\"Created At\",\"emoji\":true},\"value\":\"Created At\"}}}}},\"hash\":\"1694275539.EQSsFQL1\",\"title\":{\"type\":\"plain_text\",\"text\":\"Create Rota\",\"emoji\":true},\"clear_on_close\":true,\"notify_on_close\":true,\"close\":{\"type\":\"plain_text\",\"text\":\"Cancel\",\"emoji\":true},\"submit\":{\"type\":\"plain_text\",\"text\":\"Create\",\"emoji\":true},\"previous_view_id\":\"V05RNATD3QB\",\"root_view_id\":\"V05RNATD3QB\",\"app_id\":\"A041MF0T137\",\"external_id\":\"\",\"app_installed_team_id\":\"T042E16BURW\",\"bot_id\":\"B041A53D6ET\"},\"response_urls\":[],\"is_enterprise_install\":false,\"enterprise\":null}"
			res, err := svc.MessageActions(ctx, &gen.Action{
				Signature: "TEST",
				Timestamp: 1234567890,
				Payload:   []byte(payload),
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(res).ToNot(BeNil())

			Expect(res.ResponseAction).To(BeNil())

			rotas, err := db.New(conn).ListRotasByChannel(ctx, db.ListRotasByChannelParams{ChannelID: channelId, TeamID: teamId})
			Expect(err).ToNot(HaveOccurred())
			Expect(rotas).To(HaveLen(1))
		})

		It("Should return response without errors", func() {
			payload := "{\"type\":\"view_closed\",\"team\":{\"id\":\"T042E16BURW\",\"domain\":\"rotabot-workspace\"},\"user\":{\"id\":\"U0422UEJLD7\",\"username\":\"me1\",\"name\":\"me1\",\"team_id\":\"T042E16BURW\"},\"api_app_id\":\"A041MF0T137\",\"token\":\"XXXXXXXXXXXXXX\",\"view\":{\"id\":\"V05L4EGLTU0\",\"team_id\":\"T042E16BURW\",\"type\":\"modal\",\"blocks\":[{\"type\":\"actions\",\"block_id\":\"HOME_ACTIONS\",\"elements\":[{\"type\":\"button\",\"action_id\":\"HOME_ADD_ROTA\",\"text\":{\"type\":\"plain_text\",\"text\":\"Add Rota :heavy_plus_sign:\",\"emoji\":true}}]},{\"type\":\"header\",\"block_id\":\"=tQl\",\"text\":{\"type\":\"plain_text\",\"text\":\"Active Rotas:\",\"emoji\":true}},{\"type\":\"section\",\"block_id\":\"RTlBREta1IfYEawl\",\"text\":{\"type\":\"plain_text\",\"text\":\"Cool Rota\",\"emoji\":true},\"accessory\":{\"type\":\"overflow\",\"action_id\":\"ROTA_ELEMENT\",\"options\":[{\"text\":{\"type\":\"plain_text\",\"text\":\":spiral_note_pad: Edit Rota\",\"emoji\":true},\"value\":\"HOME_EDIT_ROTA\"}]}}],\"private_metadata\":\"{\\\"rota_id\\\":\\\"RTlBREta1IfYEawl\\\",\\\"channel_id\\\":\\\"C041Q6Z5FSP\\\"}\",\"callback_id\":\"Home\",\"state\":{\"values\":{}},\"hash\":\"1690750360.I4nd5wcR\",\"title\":{\"type\":\"plain_text\",\"text\":\"Rotabot Home\",\"emoji\":true},\"clear_on_close\":true,\"notify_on_close\":true,\"close\":null,\"submit\":null,\"previous_view_id\":null,\"root_view_id\":\"V05L4EGLTU0\",\"app_id\":\"A041MF0T137\",\"external_id\":\"\",\"app_installed_team_id\":\"T042E16BURW\",\"bot_id\":\"B041A53D6ET\"},\"is_cleared\":true,\"is_enterprise_install\":false,\"enterprise\":null}"
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
			payload := "{\"type\":\"view_ddddd\",\"team\":{\"id\":\"T042E16BURW\",\"domain\":\"rotabot-workspace\"},\"user\":{\"id\":\"U0422UEJLD7\",\"username\":\"me1\",\"name\":\"me1\",\"team_id\":\"T042E16BURW\"},\"api_app_id\":\"A041MF0T137\",\"token\":\"XXXXXXXXXXXXXX\",\"view\":{\"id\":\"V05L4EGLTU0\",\"team_id\":\"T042E16BURW\",\"type\":\"modal\",\"blocks\":[{\"type\":\"actions\",\"block_id\":\"HOME_ACTIONS\",\"elements\":[{\"type\":\"button\",\"action_id\":\"HOME_ADD_ROTA\",\"text\":{\"type\":\"plain_text\",\"text\":\"Add Rota :heavy_plus_sign:\",\"emoji\":true}}]},{\"type\":\"header\",\"block_id\":\"=tQl\",\"text\":{\"type\":\"plain_text\",\"text\":\"Active Rotas:\",\"emoji\":true}},{\"type\":\"section\",\"block_id\":\"RTlBREta1IfYEawl\",\"text\":{\"type\":\"plain_text\",\"text\":\"Cool Rota\",\"emoji\":true},\"accessory\":{\"type\":\"overflow\",\"action_id\":\"ROTA_ELEMENT\",\"options\":[{\"text\":{\"type\":\"plain_text\",\"text\":\":spiral_note_pad: Edit Rota\",\"emoji\":true},\"value\":\"HOME_EDIT_ROTA\"}]}}],\"private_metadata\":\"{\\\"rota_id\\\":\\\"RTlBREta1IfYEawl\\\",\\\"channel_id\\\":\\\"C041Q6Z5FSP\\\"}\",\"callback_id\":\"Home\",\"state\":{\"values\":{}},\"hash\":\"1690750360.I4nd5wcR\",\"title\":{\"type\":\"plain_text\",\"text\":\"Rotabot Home\",\"emoji\":true},\"clear_on_close\":true,\"notify_on_close\":true,\"close\":null,\"submit\":null,\"previous_view_id\":null,\"root_view_id\":\"V05L4EGLTU0\",\"app_id\":\"A041MF0T137\",\"external_id\":\"\",\"app_installed_team_id\":\"T042E16BURW\",\"bot_id\":\"B041A53D6ET\"},\"is_cleared\":true,\"is_enterprise_install\":false,\"enterprise\":null}"
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
