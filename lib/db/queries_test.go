package db

import (
	"context"
	"errors"
	"path/filepath"

	"github.com/testcontainers/testcontainers-go"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rotabot-io/rotabot/internal"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var _ = Describe("Rotas", func() {
	var ctx context.Context
	var q *Queries

	BeforeEach(func() {
		ctx = context.Background()

		container, err := internal.RunContainer(ctx,
			postgres.WithInitScripts(filepath.Join("..", "..", "assets", "structure.sql")),
			testcontainers.WithWaitStrategy(internal.DefaultWaitStrategy()),
		)
		Expect(err).ToNot(HaveOccurred())

		dbUrl, err := container.ConnectionString(ctx, "sslmode=disable")
		Expect(err).ToNot(HaveOccurred())

		conn, err := pgx.Connect(ctx, dbUrl)
		Expect(err).ToNot(HaveOccurred())
		q = New(conn)

		DeferCleanup(func() {
			_ = container.Terminate(ctx)
			conn.Close(ctx)
		})
	})

	Describe("CreateOrUpdateRota", func() {
		It("Should create rota if id is null", func() {
			id, err := q.CreateOrUpdateRota(ctx, CreateOrUpdateRotaParams{
				ChannelID: "foo",
				TeamID:    "bar",
				Name:      "baz",
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(id).ToNot(BeEmpty())
		})

		It("Should fail to create two identical rotas", func() {
			req := CreateOrUpdateRotaParams{
				ChannelID: "foo",
				TeamID:    "bar",
				Name:      "baz",
			}
			_, err := q.CreateOrUpdateRota(ctx, req)
			Expect(err).ToNot(HaveOccurred())

			_, err = q.CreateOrUpdateRota(ctx, req)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(ErrAlreadyExists))
		})

		It("Should update rota if id is null", func() {
			id, err := q.CreateOrUpdateRota(ctx, CreateOrUpdateRotaParams{
				ChannelID: "foo",
				TeamID:    "bar",
				Name:      "baz",
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(id).ToNot(BeEmpty())

			updated, err := q.CreateOrUpdateRota(ctx, CreateOrUpdateRotaParams{
				RotaID: id,
				Name:   "bazbaz",
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(id).To(Equal(updated))

			rota, err := q.FindRotaByID(ctx, id)
			Expect(err).ToNot(HaveOccurred())
			Expect(rota.Name).To(Equal("bazbaz"))
		})

		It("Should fail to update something it does not exist", func() {
			_, err := q.CreateOrUpdateRota(ctx, CreateOrUpdateRotaParams{
				RotaID: "not_found",
				Name:   "bazbaz",
			})
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(ErrNotFound))
		})
	})

	Describe("FindRotaByID", func() {
		When("rota does not exist", func() {
			It("should return ErrNotFound", func() {
				_, err := q.FindRotaByID(ctx, "not_found")
				Expect(err).To(HaveOccurred())
			})
		})
		When("rota exists", func() {
			It("should return rota", func() {
				id, err := q.saveRota(ctx, saveRotaParams{
					ChannelID: "C123",
					TeamID:    "T123",
					Name:      "test",
					Metadata: RotaMetadata{
						Frequency:      RFDaily,
						SchedulingType: RSRandom,
					},
				})
				Expect(err).ToNot(HaveOccurred())

				rota, err := q.FindRotaByID(ctx, id)
				Expect(err).ToNot(HaveOccurred())
				Expect(rota.ID).To(Equal(id))
				Expect(rota.Name).To(Equal("test"))
				Expect(rota.Metadata.Frequency).To(Equal(RFDaily))
				Expect(rota.Metadata.SchedulingType).To(Equal(RSRandom))
			})
		})
	})

	Describe("ListRotasByChannel", func() {
		var (
			channelID string
			teamID    string
		)
		BeforeEach(func() {
			channelID = "C123"
			teamID = "T123"
		})

		It("should return empty array if no rotas", func() {
			rotas, err := q.ListRotasByChannel(ctx, ListRotasByChannelParams{ChannelID: channelID, TeamID: teamID})
			Expect(err).ToNot(HaveOccurred())
			Expect(rotas).To(HaveLen(0))
		})

		It("should return rotas when one exist", func() {
			_, err := q.saveRota(ctx, saveRotaParams{
				ChannelID: channelID,
				TeamID:    teamID,
				Name:      "test",
				Metadata: RotaMetadata{
					Frequency:      RFDaily,
					SchedulingType: RSRandom,
				},
			})
			Expect(err).ToNot(HaveOccurred())

			rotas, err := q.ListRotasByChannel(ctx, ListRotasByChannelParams{ChannelID: channelID, TeamID: teamID})
			Expect(err).ToNot(HaveOccurred())
			Expect(rotas).To(HaveLen(1))
		})

		It("should return return nothing rota is on another team", func() {
			_, err := q.saveRota(ctx, saveRotaParams{
				ChannelID: channelID,
				TeamID:    "another_team",
				Name:      "test",
				Metadata: RotaMetadata{
					Frequency:      RFDaily,
					SchedulingType: RSRandom,
				},
			})
			Expect(err).ToNot(HaveOccurred())

			rotas, err := q.ListRotasByChannel(ctx, ListRotasByChannelParams{ChannelID: channelID, TeamID: teamID})
			Expect(err).ToNot(HaveOccurred())
			Expect(rotas).To(HaveLen(0))
		})
	})

	Describe("SaveRota", func() {
		var (
			channelID string
			teamID    string
			name      string
		)
		BeforeEach(func() {
			channelID = "C123"
			teamID = "T123"
			name = "test"
		})

		It("should create rota", func() {
			id, err := q.saveRota(ctx, saveRotaParams{
				ChannelID: channelID,
				TeamID:    teamID,
				Name:      name,
				Metadata: RotaMetadata{
					Frequency:      RFDaily,
					SchedulingType: RSRandom,
				},
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(id).ToNot(BeEmpty())
		})

		It("should fail when rota already exist", func() {
			p := saveRotaParams{
				ChannelID: channelID,
				TeamID:    teamID,
				Name:      name,
				Metadata: RotaMetadata{
					Frequency:      RFDaily,
					SchedulingType: RSRandom,
				},
			}
			id, err := q.saveRota(ctx, p)
			Expect(err).ToNot(HaveOccurred())
			Expect(id).ToNot(BeEmpty())

			_, err = q.saveRota(ctx, p)
			Expect(err).To(HaveOccurred())

			var pgError *pgconn.PgError
			Expect(errors.As(err, &pgError)).To(BeTrue())
			Expect(pgError.Code).To(Equal("23505"))

			rotas, err := q.ListRotasByChannel(ctx, ListRotasByChannelParams{ChannelID: channelID, TeamID: teamID})
			Expect(err).ToNot(HaveOccurred())
			Expect(rotas).To(HaveLen(1))
		})
	})

	Describe("UpdateRota", func() {
		When("rota does not exist", func() {
			It("should return ErrNotFound", func() {
				id, err := q.updateRota(
					ctx,
					updateRotaParams{
						ID:   "not_found",
						Name: "test",
						Metadata: RotaMetadata{
							Frequency:      RFDaily,
							SchedulingType: RSRandom,
						},
					},
				)
				Expect(err).To(HaveOccurred())
				Expect(id).To(BeEmpty())
			})
		})
		When("rota exists", func() {
			It("should return rota", func() {
				id, err := q.saveRota(ctx, saveRotaParams{
					ChannelID: "C123",
					TeamID:    "T123",
					Name:      "test",
					Metadata: RotaMetadata{
						Frequency:      RFDaily,
						SchedulingType: RSRandom,
					},
				})
				Expect(err).ToNot(HaveOccurred())

				id, err = q.updateRota(
					ctx,
					updateRotaParams{
						ID:   id,
						Name: "test test",
						Metadata: RotaMetadata{
							Frequency:      RFWeekly,
							SchedulingType: RSCreated,
						},
					},
				)
				Expect(err).ToNot(HaveOccurred())

				rota, err := q.FindRotaByID(ctx, id)
				Expect(err).ToNot(HaveOccurred())
				Expect(rota.ID).To(Equal(id))
				Expect(rota.Name).To(Equal("test test"))
				Expect(rota.Metadata.Frequency).To(Equal(RFWeekly))
				Expect(rota.Metadata.SchedulingType).To(Equal(RSCreated))
			})
		})
	})
})
