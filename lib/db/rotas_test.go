package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rotabot-io/rotabot/internal"
)

var _ = Describe("Rotas", func() {
	var ctx context.Context
	var connString string
	var q *Queries

	BeforeEach(func() {
		ctx = context.Background()

		container, err := internal.RunContainer(ctx)
		Expect(err).ToNot(HaveOccurred())

		err = Migrate(ctx, container.ConnectionString())
		Expect(err).ToNot(HaveOccurred())

		connString = container.ConnectionString()

		conn, err := pgx.Connect(ctx, connString)
		Expect(err).ToNot(HaveOccurred())
		q = New(conn)

		DeferCleanup(func() {
			_ = container.Terminate(ctx)
			conn.Close(ctx)
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
})
