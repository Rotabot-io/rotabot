package views

import (
	"context"
	"encoding/json"
	"path/filepath"

	"go.uber.org/mock/gomock"

	"github.com/testcontainers/testcontainers-go"

	"github.com/jackc/pgx/v5"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gen "github.com/rotabot-io/rotabot/gen/slack"
	"github.com/rotabot-io/rotabot/internal"
	"github.com/rotabot-io/rotabot/lib/db"
	"github.com/rotabot-io/rotabot/slack/slackclient"
	"github.com/rotabot-io/rotabot/slack/slackclient/mock_slackclient"
	"github.com/slack-go/slack"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var _ = Describe("Home", func() {
	var (
		ctx  context.Context
		sc   *mock_slackclient.MockSlackClient
		home *Home
		repo db.Repository
		conn *pgx.Conn
	)

	BeforeEach(func() {
		ctx = context.Background()

		container, err := internal.RunContainer(ctx,
			postgres.WithInitScripts(filepath.Join("..", "..", "assets", "structure.sql")),
			testcontainers.WithWaitStrategy(internal.DefaultWaitStrategy()),
		)
		Expect(err).ToNot(HaveOccurred())

		connString, err := container.ConnectionString(ctx, "sslmode=disable")
		Expect(err).ToNot(HaveOccurred())

		conn, err = pgx.Connect(ctx, connString)
		Expect(err).ToNot(HaveOccurred())

		tx, err := conn.Begin(ctx)
		Expect(err).ToNot(HaveOccurred())

		repo = db.New(tx)
		home = &Home{
			Repository: repo,
		}

		DeferCleanup(func() {
			_ = container.Terminate(ctx)
			_ = conn.Close(ctx)
			_ = tx.Rollback(ctx)
		})
	})

	// Create a mock and assign it to the sc variable at the start of each test
	slackclient.MockSlackClient(&ctx, &sc, nil)

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
			id, err := repo.CreateOrUpdateRota(ctx, db.CreateOrUpdateRotaParams{
				Name:      "Test Rota",
				ChannelID: home.State.ChannelID,
				TeamID:    home.State.TeamID,
				Metadata: db.RotaMetadata{
					Frequency:      db.RFMonthly,
					SchedulingType: db.RSRandom,
				},
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
			Expect(sectionBlock.BlockID).To(Equal(id))
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
			home.State.action = "unknown"

			_, err := home.OnAction(ctx)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("unknown_action"))
		})

		It("calls slack api to push add_rota modal", func() {
			home.State.action = HASaveRota

			sc.EXPECT().PushViewContext(ctx, triggerID, gomock.Cond(func(x any) bool {
				view := x.(slack.ModalViewRequest)

				var m Metadata
				err := json.Unmarshal([]byte(view.PrivateMetadata), &m)
				Expect(err).ToNot(HaveOccurred())

				Expect(view.CallbackID).To(Equal(string(VTSaveRota)))
				Expect(m.ChannelID).To(Equal(channelID))
				return Expect(m.RotaID).To(BeEmpty())
			})).Return(nil, nil).Times(1)

			_, err := home.OnAction(ctx)
			Expect(err).ToNot(HaveOccurred())
		})

		It("calls slack api to push update modal", func() {
			home.State.action = HASaveRota
			id, err := home.Repository.CreateOrUpdateRota(ctx, db.CreateOrUpdateRotaParams{
				Name:      "Rota",
				ChannelID: channelID,
				TeamID:    teamID,
			})
			Expect(err).ToNot(HaveOccurred())

			home.State.rotaID = id
			sc.EXPECT().PushViewContext(ctx, triggerID, gomock.Cond(func(x any) bool {
				view := x.(slack.ModalViewRequest)

				var m Metadata
				err := json.Unmarshal([]byte(view.PrivateMetadata), &m)
				Expect(err).ToNot(HaveOccurred())

				Expect(view.CallbackID).To(Equal(string(VTSaveRota)))
				Expect(m.ChannelID).To(Equal(channelID))
				return Expect(m.RotaID).To(Equal(id))
			})).Return(nil, nil).Times(1)

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

			sc.EXPECT().OpenViewContext(ctx, home.State.TriggerID, gomock.Cond(func(x any) bool {
				view := x.(slack.ModalViewRequest)

				var m Metadata
				err := json.Unmarshal([]byte(view.PrivateMetadata), &m)
				Expect(err).ToNot(HaveOccurred())

				Expect(view.CallbackID).To(Equal(string(VTHome)))
				Expect(m.ChannelID).To(Equal("C123"))
				return Expect(m.RotaID).To(BeEmpty())
			})).Return(nil, nil).Times(1)

			err = home.Render(ctx, p)
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
