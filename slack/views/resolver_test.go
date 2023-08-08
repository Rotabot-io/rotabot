package views

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rotabot-io/rotabot/lib/db"
	"github.com/slack-go/slack"
)

var _ = Describe("Resolver", func() {
	var ctx context.Context

	BeforeEach(func() {
		ctx = context.Background()
	})

	Describe("Home", func() {
		It("resolves a home view without actions", func() {
			params := ResolverParams{
				Action: slack.InteractionCallback{
					View: slack.View{
						CallbackID:      string(VTHome),
						PrivateMetadata: "{\"rota_id\":\"\",\"channel_id\":\"C123\"}",
					},
					TriggerID: "T123",
					Team: slack.Team{
						ID: "TE123",
					},
				},
			}

			view, err := Resolve(ctx, params)
			Expect(err).ToNot(HaveOccurred())

			homeView, ok := view.(*Home)
			Expect(ok).To(BeTrue())

			Expect(homeView.State.TriggerID).To(Equal("T123"))
			Expect(homeView.State.ChannelID).To(Equal("C123"))
			Expect(homeView.State.TeamID).To(Equal("TE123"))
			Expect(homeView.State.action).To(BeEmpty())
		})

		It("resolves save view action without rota_id to trigger a create", func() {
			params := ResolverParams{
				Action: slack.InteractionCallback{
					View: slack.View{
						PrivateMetadata: "{\"rota_id\":\"ROTA_ID\",\"channel_id\":\"C123\"}",
						CallbackID:      string(VTHome),
					},
					ActionCallback: slack.ActionCallbacks{
						BlockActions: []*slack.BlockAction{
							{
								ActionID: string(HASaveRota),
							},
						},
					},
				},
			}

			view, err := Resolve(ctx, params)
			Expect(err).ToNot(HaveOccurred())

			homeView, ok := view.(*Home)
			Expect(ok).To(BeTrue())

			Expect(homeView.State.action).To(Equal(HASaveRota))
			Expect(homeView.State.rotaID).To(BeEmpty())
		})

		It("resolves save view action with the rota_id to trigger an update", func() {
			params := ResolverParams{
				Action: slack.InteractionCallback{
					View: slack.View{
						PrivateMetadata: "{\"rota_id\":\"ROTA_ID\",\"channel_id\":\"C123\"}",
						CallbackID:      string(VTHome),
					},
					ActionCallback: slack.ActionCallbacks{
						BlockActions: []*slack.BlockAction{
							{
								ActionID: string(HSRota),
								BlockID:  "ROTA_ID",
								SelectedOption: slack.OptionBlockObject{
									Value: string(HASaveRota),
								},
							},
						},
					},
				},
			}

			view, err := Resolve(ctx, params)
			Expect(err).ToNot(HaveOccurred())

			homeView, ok := view.(*Home)
			Expect(ok).To(BeTrue())

			Expect(homeView.State.action).To(Equal(HASaveRota))
			Expect(homeView.State.rotaID).To(Equal("ROTA_ID"))
		})
	})
	Describe("SaveRota", func() {
		It("resolves a add rota view with default state", func() {
			params := ResolverParams{
				Action: slack.InteractionCallback{
					View: slack.View{
						CallbackID:      string(VTSaveRota),
						PrivateMetadata: "{\"rota_id\":\"\",\"channel_id\":\"C123\"}",
						PreviousViewID:  "PREV123",
						ExternalID:      "EXT123",
						State:           &slack.ViewState{},
					},
					TriggerID: "T123",
					Team: slack.Team{
						ID: "TE123",
					},
				},
			}

			view, err := Resolve(ctx, params)
			Expect(err).ToNot(HaveOccurred())

			addView, ok := view.(*SaveRota)
			Expect(ok).To(BeTrue())

			Expect(addView.State.TriggerID).To(Equal("T123"))
			Expect(addView.State.ChannelID).To(Equal("C123"))
			Expect(addView.State.TeamID).To(Equal("TE123"))
			Expect(addView.State.previousViewID).To(Equal("PREV123"))
			Expect(addView.State.externalID).To(Equal("EXT123"))

			Expect(addView.State.rotaName).To(BeEmpty())
			Expect(addView.State.rotaID).To(BeEmpty())
			Expect(addView.State.frequency).To(Equal(db.RFWeekly))
			Expect(addView.State.schedulingType).To(Equal(db.RSCreated))
		})

		It("resolves a add rota view with the state given on the action", func() {
			params := ResolverParams{
				Action: slack.InteractionCallback{
					View: slack.View{
						CallbackID:      string(VTSaveRota),
						PrivateMetadata: "{\"rota_id\":\"ROTA_ID\",\"channel_id\":\"C123\"}",
						State: &slack.ViewState{
							Values: map[string]map[string]slack.BlockAction{
								"ROTA_NAME": {
									"ROTA_NAME": {
										Value: "Test Rota",
									},
								},
								"ROTA_FREQUENCY": {
									"ROTA_FREQUENCY": {
										SelectedOption: slack.OptionBlockObject{
											Value: string(db.RFMonthly),
										},
									},
								},
								"ROTA_TYPE": {
									"ROTA_TYPE": {
										SelectedOption: slack.OptionBlockObject{
											Value: string(db.RSRandom),
										},
									},
								},
							},
						},
					},
				},
			}

			view, err := Resolve(ctx, params)
			Expect(err).ToNot(HaveOccurred())

			addView, ok := view.(*SaveRota)
			Expect(ok).To(BeTrue())
			Expect(addView.State.rotaID).To(Equal("ROTA_ID"))

			Expect(addView.State.rotaName).To(Equal("Test Rota"))
			Expect(addView.State.frequency).To(Equal(db.RFMonthly))
			Expect(addView.State.schedulingType).To(Equal(db.RSRandom))
		})
	})
})
