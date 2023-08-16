package db

import (
	"context"
	"errors"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Rotas", func() {
	var ctx context.Context
	var q *Queries

	BeforeEach(func() {
		ctx = context.Background()

		container, err := postgres.RunContainer(ctx,
			testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(5*time.Second)),
		)
		Expect(err).ToNot(HaveOccurred())

		dbUrl, err := container.ConnectionString(ctx, "sslmode=disable")
		Expect(err).ToNot(HaveOccurred())

		err = Migrate(ctx, dbUrl)
		Expect(err).ToNot(HaveOccurred())

		conn, err := pgx.Connect(ctx, dbUrl)
		Expect(err).ToNot(HaveOccurred())
		q = New(conn)

		DeferCleanup(func() {
			_ = container.Terminate(ctx)
			conn.Close(ctx)
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
				id, err := q.SaveRota(ctx, SaveRotaParams{
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
			_, err := q.SaveRota(ctx, SaveRotaParams{
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
			_, err := q.SaveRota(ctx, SaveRotaParams{
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
			id, err := q.SaveRota(ctx, SaveRotaParams{
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
			p := SaveRotaParams{
				ChannelID: channelID,
				TeamID:    teamID,
				Name:      name,
				Metadata: RotaMetadata{
					Frequency:      RFDaily,
					SchedulingType: RSRandom,
				},
			}
			id, err := q.SaveRota(ctx, p)
			Expect(err).ToNot(HaveOccurred())
			Expect(id).ToNot(BeEmpty())

			_, err = q.SaveRota(ctx, p)
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
				id, err := q.UpdateRota(
					ctx,
					UpdateRotaParams{
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
				id, err := q.SaveRota(ctx, SaveRotaParams{
					ChannelID: "C123",
					TeamID:    "T123",
					Name:      "test",
					Metadata: RotaMetadata{
						Frequency:      RFDaily,
						SchedulingType: RSRandom,
					},
				})
				Expect(err).ToNot(HaveOccurred())

				id, err = q.UpdateRota(
					ctx,
					UpdateRotaParams{
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
