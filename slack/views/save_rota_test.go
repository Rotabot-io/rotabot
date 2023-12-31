package views

import (
	"context"
	"path/filepath"

	"github.com/testcontainers/testcontainers-go"

	"github.com/rotabot-io/rotabot/internal"

	"github.com/rotabot-io/rotabot/slack/slackclient"

	"github.com/jackc/pgx/v5"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gen "github.com/rotabot-io/rotabot/gen/slack"
	"github.com/rotabot-io/rotabot/lib/db"
	"github.com/rotabot-io/rotabot/slack/slackclient/mock_slackclient"
	"github.com/slack-go/slack"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"go.uber.org/mock/gomock"
)

var _ = Describe("SaveRota", func() {
	var (
		ctx       context.Context
		sc        *mock_slackclient.MockSlackClient
		repo      db.Repository
		addRota   *SaveRota
		conn      *pgx.Conn
		channelID string
		teamID    string
		triggerID string
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
		addRota = &SaveRota{
			Repository: repo,
		}
		channelID = "CH123"
		teamID = "TM123"
		triggerID = "TR123"

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
			Expect(addRota.CallbackID()).To(Equal(VTSaveRota))
		})
	})

	Describe("DefaultState", func() {
		It("returns a default state", func() {
			expectedState := &SaveRotaState{
				frequency:      db.RFWeekly,
				schedulingType: db.RSCreated,
			}

			Expect(addRota.DefaultState()).To(Equal(expectedState))
		})
	})

	Describe("BuildProps", func() {
		When("Rota does not exist", func() {
			BeforeEach(func() {
				addRota.State = &SaveRotaState{
					TriggerID: triggerID,
					ChannelID: channelID,
					TeamID:    teamID,
				}
			})
			It("builds props with the correct values in the state", func() {
				addRota.State = &SaveRotaState{
					TriggerID:      triggerID,
					ChannelID:      channelID,
					TeamID:         teamID,
					frequency:      db.RFMonthly,
					schedulingType: db.RSRandom,
				}

				p, err := addRota.BuildProps(ctx)
				Expect(err).ToNot(HaveOccurred())
				Expect(p).To(BeAssignableToTypeOf(&SaveRotaProps{}))

				props := p.(*SaveRotaProps)
				Expect(props.title.Text).To(Equal("Create Rota"))
				Expect(props.close.Text).To(Equal("Cancel"))
				Expect(props.submit.Text).To(Equal("Create"))

				Expect(props.blocks.BlockSet).To(HaveLen(4))
				Expect(props.blocks.BlockSet[0]).To(BeAssignableToTypeOf(&slack.InputBlock{}))
				Expect(props.blocks.BlockSet[1]).To(BeAssignableToTypeOf(&slack.SectionBlock{}))
				Expect(props.blocks.BlockSet[2]).To(BeAssignableToTypeOf(&slack.SectionBlock{}))
				Expect(props.blocks.BlockSet[3]).To(BeAssignableToTypeOf(&slack.SectionBlock{}))

				inputBlock := props.blocks.BlockSet[0].(*slack.InputBlock)
				Expect(inputBlock.BlockID).To(Equal("ROTA_NAME"))

				frequencySelect := props.blocks.BlockSet[1].(*slack.SectionBlock)
				Expect(frequencySelect.BlockID).To(Equal("ROTA_FREQUENCY"))

				schedulingType := props.blocks.BlockSet[2].(*slack.SectionBlock)
				Expect(schedulingType.BlockID).To(Equal("ROTA_TYPE"))

				userSelect := props.blocks.BlockSet[3].(*slack.SectionBlock)
				Expect(userSelect.BlockID).To(Equal("ROTA_MEMBERS"))
			})
		})

		When("Rota does exist", func() {
			BeforeEach(func() {
				id, err := repo.CreateOrUpdateRota(ctx, db.CreateOrUpdateRotaParams{
					Name:      "Test Rota",
					TeamID:    teamID,
					ChannelID: channelID,
					Metadata: db.RotaMetadata{
						Frequency:      db.RFMonthly,
						SchedulingType: db.RSRandom,
					},
				})
				Expect(err).ToNot(HaveOccurred())

				addRota.State = &SaveRotaState{
					rotaID:    id,
					TriggerID: triggerID,
					ChannelID: channelID,
					TeamID:    teamID,
				}
			})
			It("builds the props with the values from the rota", func() {
				p, err := addRota.BuildProps(ctx)
				Expect(err).ToNot(HaveOccurred())
				Expect(p).To(BeAssignableToTypeOf(&SaveRotaProps{}))

				props := p.(*SaveRotaProps)
				Expect(props.title.Text).To(Equal("Update Rota"))
				Expect(props.close.Text).To(Equal("Cancel"))
				Expect(props.submit.Text).To(Equal("Update"))

				Expect(props.blocks.BlockSet).To(HaveLen(4))
				Expect(props.blocks.BlockSet[0]).To(BeAssignableToTypeOf(&slack.InputBlock{}))
				Expect(props.blocks.BlockSet[1]).To(BeAssignableToTypeOf(&slack.SectionBlock{}))
				Expect(props.blocks.BlockSet[2]).To(BeAssignableToTypeOf(&slack.SectionBlock{}))
				Expect(props.blocks.BlockSet[3]).To(BeAssignableToTypeOf(&slack.SectionBlock{}))

				inputBlock := props.blocks.BlockSet[0].(*slack.InputBlock)
				Expect(inputBlock.BlockID).To(Equal("ROTA_NAME"))
				Expect(inputBlock.Element.(*slack.PlainTextInputBlockElement).InitialValue).To(Equal("Test Rota"))

				frequencySelect := props.blocks.BlockSet[1].(*slack.SectionBlock)
				Expect(frequencySelect.BlockID).To(Equal("ROTA_FREQUENCY"))

				schedulingType := props.blocks.BlockSet[2].(*slack.SectionBlock)
				Expect(schedulingType.BlockID).To(Equal("ROTA_TYPE"))

				userSelect := props.blocks.BlockSet[3].(*slack.SectionBlock)
				Expect(userSelect.BlockID).To(Equal("ROTA_MEMBERS"))
			})
		})
	})

	Describe("OnAction", func() {
		It("returns without doing anything", func() {
			res, err := addRota.OnAction(ctx)
			Expect(err).ToNot(HaveOccurred())

			expectedRes := &gen.ActionResponse{}
			Expect(res).To(Equal(expectedRes))
		})
	})

	Describe("OnClose", func() {
		It("returns without doing anything", func() {
			res, err := addRota.OnClose(ctx)
			Expect(err).ToNot(HaveOccurred())

			expectedRes := &gen.ActionResponse{}
			Expect(res).To(Equal(expectedRes))
		})
	})

	Describe("OnSubmit", func() {
		When("the user creates a rota that already exists", func() {
			It("returns an error", func() {
				_, err := repo.CreateOrUpdateRota(ctx, db.CreateOrUpdateRotaParams{
					Name:      "test",
					TeamID:    teamID,
					ChannelID: channelID,
					Metadata: db.RotaMetadata{
						Frequency:      db.RFMonthly,
						SchedulingType: db.RSRandom,
					},
				})
				Expect(err).ToNot(HaveOccurred())

				addRota.State = &SaveRotaState{
					TriggerID:      triggerID,
					ChannelID:      channelID,
					TeamID:         teamID,
					rotaName:       "test",
					frequency:      db.RFWeekly,
					schedulingType: db.RSCreated,
				}

				res, err := addRota.OnSubmit(ctx)
				Expect(err).ToNot(HaveOccurred())
				Expect(res).ToNot(BeNil())

				expectedResAction := string(slack.RAErrors)
				expectedRes := &gen.ActionResponse{
					ResponseAction: &expectedResAction,
					Errors: map[string]string{
						"ROTA_NAME": "A rota with this name already exists in this channel.",
					},
				}
				Expect(res).To(Equal(expectedRes))
			})
		})
		When("the user creates a rota that does not exist", func() {
			It("creates the rota and updates the home view", func() {
				addRota.State = &SaveRotaState{
					TriggerID:      triggerID,
					ChannelID:      channelID,
					TeamID:         teamID,
					rotaName:       "test",
					frequency:      db.RFWeekly,
					schedulingType: db.RSCreated,
					externalID:     "E123",
					previousViewID: "PV123",
				}
				sc.EXPECT().
					UpdateViewContext(ctx, gomock.Any(), "E123", "", "PV123").
					Return(nil, nil).Times(1)

				res, err := addRota.OnSubmit(ctx)
				Expect(err).ToNot(HaveOccurred())
				Expect(res).ToNot(BeNil())

				expectedRes := &gen.ActionResponse{}
				Expect(res).To(Equal(expectedRes))
			})
		})
	})

	Describe("Render", func() {
		BeforeEach(func() {
			addRota.State = &SaveRotaState{
				TriggerID: triggerID,
				ChannelID: channelID,
				TeamID:    teamID,
			}
		})

		It("calls slack to open modal", func() {
			p, err := addRota.BuildProps(ctx)
			Expect(err).ToNot(HaveOccurred())
			Expect(p).To(BeAssignableToTypeOf(&SaveRotaProps{}))
			props := p.(*SaveRotaProps)

			expectedView := slack.ModalViewRequest{
				Type:            slack.VTModal,
				Title:           props.title,
				Blocks:          props.blocks,
				Close:           props.close,
				Submit:          props.submit,
				CallbackID:      string(VTSaveRota),
				NotifyOnClose:   true,
				ClearOnClose:    true,
				PrivateMetadata: "{\"rota_id\":\"\",\"channel_id\":\"CH123\"}",
			}
			sc.EXPECT().
				OpenViewContext(ctx, addRota.State.TriggerID, expectedView).
				Return(nil, nil).Times(1)

			err = addRota.Render(ctx, p)
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
