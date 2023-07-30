package views

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gen "github.com/rotabot-io/rotabot/gen/slack"
	"github.com/rotabot-io/rotabot/internal"
	"github.com/rotabot-io/rotabot/lib/db"
	"github.com/rotabot-io/rotabot/slack/slackclient"
	"github.com/rotabot-io/rotabot/slack/slackclient/mock_slackclient"
	"github.com/slack-go/slack"
)

var _ = Describe("Home", func() {
	var (
		ctx     context.Context
		sc      *mock_slackclient.MockSlackClient
		queries *db.Queries
		home    *Home
	)

	BeforeEach(func() {
		ctx = context.Background()
	})

	// Create a mock and assign it to the sc variable at the start of each test
	slackclient.MockSlackClient(&ctx, &sc, nil)

	BeforeEach(func() {
		container, err := internal.RunContainer(ctx)
		Expect(err).ToNot(HaveOccurred())

		err = db.Migrate(ctx, container.ConnectionString())
		Expect(err).ToNot(HaveOccurred())

		conn, err := pgx.Connect(ctx, container.ConnectionString())
		Expect(err).ToNot(HaveOccurred())

		queries = db.New(conn)

		DeferCleanup(func() {
			conn.Close(ctx)
		})

		home = &Home{
			Queries: queries,
		}
	})

	Describe("Callback", func() {
		It("resolves a home view without actions", func() {
			Expect(home.CallbackID()).To(Equal(VTHome))
		})
	})

	Describe("DefaultState", func() {
		It("returns a default state", func() {
			expectedState := &HomeState{}

			Expect(home.DefaultState()).To(Equal(expectedState))
		})
	})

	Describe("BuildProps", func() {
		BeforeEach(func() {
			home.State = &HomeState{
				TriggerID: "T123",
				ChannelID: "C123",
				TeamID:    "TM123",
			}
		})

		It("returns props when no rota exists", func() {
			p, err := home.BuildProps(ctx)
			Expect(err).ToNot(HaveOccurred())

			Expect(p).To(BeAssignableToTypeOf(&HomeProps{}))
			props := p.(*HomeProps)
			Expect(props.title.Text).To(Equal("Rotabot Home"))

			Expect(props.blocks.BlockSet).To(HaveLen(2))

			Expect(props.blocks.BlockSet[0]).To(BeAssignableToTypeOf(&slack.ActionBlock{}))
			Expect(props.blocks.BlockSet[1]).To(BeAssignableToTypeOf(&slack.SectionBlock{}))

			actionBlock := props.blocks.BlockSet[0].(*slack.ActionBlock)
			Expect(actionBlock.BlockID).To(Equal("HOME_ACTIONS"))
			Expect(actionBlock.Elements.ElementSet).To(HaveLen(1))
			Expect(actionBlock.Elements.ElementSet[0]).To(BeAssignableToTypeOf(&slack.ButtonBlockElement{}))
			button := actionBlock.Elements.ElementSet[0].(*slack.ButtonBlockElement)
			Expect(button.Text.Text).To(Equal("Add Rota :heavy_plus_sign:"))

			sectionBlock := props.blocks.BlockSet[1].(*slack.SectionBlock)
			Expect(sectionBlock.Text.Text).To(Equal("Active Rotas:"))
		})

		It("returns home props when rotas exist", func() {
			id, err := home.Queries.SaveRota(ctx, db.SaveRotaParams{
				ChannelID: home.State.ChannelID,
				TeamID:    home.State.TeamID,
				Name:      "Test Rota",
			})
			Expect(err).ToNot(HaveOccurred())

			p, err := home.BuildProps(ctx)
			Expect(err).ToNot(HaveOccurred())

			Expect(p).To(BeAssignableToTypeOf(&HomeProps{}))
			props := p.(*HomeProps)
			Expect(props.title.Text).To(Equal("Rotabot Home"))

			Expect(props.blocks.BlockSet).To(HaveLen(3))

			Expect(props.blocks.BlockSet[0]).To(BeAssignableToTypeOf(&slack.ActionBlock{}))
			Expect(props.blocks.BlockSet[1]).To(BeAssignableToTypeOf(&slack.SectionBlock{}))
			Expect(props.blocks.BlockSet[2]).To(BeAssignableToTypeOf(&slack.SectionBlock{}))

			actionBlock := props.blocks.BlockSet[0].(*slack.ActionBlock)
			Expect(actionBlock.BlockID).To(Equal("HOME_ACTIONS"))
			Expect(actionBlock.Elements.ElementSet).To(HaveLen(1))
			Expect(actionBlock.Elements.ElementSet[0]).To(BeAssignableToTypeOf(&slack.ButtonBlockElement{}))
			button := actionBlock.Elements.ElementSet[0].(*slack.ButtonBlockElement)
			Expect(button.Text.Text).To(Equal("Add Rota :heavy_plus_sign:"))

			sectionBlock := props.blocks.BlockSet[1].(*slack.SectionBlock)
			Expect(sectionBlock.Text.Text).To(Equal("Active Rotas:"))

			sectionBlock = props.blocks.BlockSet[2].(*slack.SectionBlock)
			Expect(sectionBlock.Text.Text).To(Equal("Test Rota"))
			Expect(strings.Contains(sectionBlock.BlockID, "ROTA_ELEMENT")).To(BeTrue())
			Expect(strings.Contains(sectionBlock.BlockID, id)).To(BeTrue())
		})
	})

	Describe("OnAction", func() {
		var (
			channelID string
			teamID    string
			triggerID string
		)
		BeforeEach(func() {
			channelID = "CH123"
			teamID = "TM123"
			triggerID = "TR123"

			home.State = &HomeState{
				ChannelID: channelID,
				TeamID:    teamID,
				TriggerID: triggerID,
			}
		})

		It("returns an error when actions is unknown", func() {
			home.State.Action = "unknown"

			_, err := home.OnAction(ctx)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("unknown_action"))
		})

		It("calls slack api to push add_rota modal", func() {
			home.State.Action = HomeActionAddRota

			addRota := &AddRota{
				Queries: queries,
				State: &AddRotaState{
					TriggerID:      triggerID,
					ChannelID:      channelID,
					TeamID:         teamID,
					frequency:      db.RFWeekly,
					schedulingType: db.RSCreated,
				},
			}
			p, err := addRota.BuildProps(ctx)
			Expect(err).ToNot(HaveOccurred())
			props := p.(*AddRotaProps)

			expectedModal := slack.ModalViewRequest{
				Type:            slack.VTModal,
				Title:           props.title,
				Blocks:          props.blocks,
				Close:           props.close,
				Submit:          props.submit,
				CallbackID:      string(VTAddRota),
				NotifyOnClose:   true,
				ClearOnClose:    true,
				PrivateMetadata: channelID,
			}
			sc.EXPECT().PushViewContext(ctx, triggerID, expectedModal).Return(nil, nil).Times(1)

			_, err = home.OnAction(ctx)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("OnClose", func() {
		It("returns without doing anything", func() {
			res, err := home.OnClose(ctx)
			Expect(err).ToNot(HaveOccurred())

			expectedRes := &gen.ActionResponse{}
			Expect(res).To(Equal(expectedRes))
		})
	})

	Describe("OnSubmit", func() {
		It("fails", func() {
			_, err := home.OnSubmit(ctx)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("Render", func() {
		BeforeEach(func() {
			home.State = &HomeState{
				TriggerID: "T123",
				ChannelID: "C123",
				TeamID:    "TM123",
			}
		})

		It("calls slack to open modal", func() {
			p, err := home.BuildProps(ctx)
			Expect(err).ToNot(HaveOccurred())
			Expect(p).To(BeAssignableToTypeOf(&HomeProps{}))
			props := p.(*HomeProps)

			expectedView := slack.ModalViewRequest{
				Type:            slack.VTModal,
				Title:           props.title,
				Blocks:          props.blocks,
				CallbackID:      string(VTHome),
				NotifyOnClose:   true,
				ClearOnClose:    true,
				PrivateMetadata: home.State.ChannelID,
			}
			sc.EXPECT().
				OpenViewContext(ctx, home.State.TriggerID, expectedView).
				Return(nil, nil).Times(1)

			err = home.Render(ctx, p)
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
