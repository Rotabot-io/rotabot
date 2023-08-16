package db

import (
	"context"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/jackc/pgx/v5"
	. "github.com/onsi/ginkgo/v2"

	. "github.com/onsi/gomega"
)

var _ = Describe("Migration", func() {
	var ctx context.Context
	var connString string
	var q *Queries

	BeforeEach(func() {
		var err error
		ctx = context.Background()

		container, err := postgres.RunContainer(ctx,
			testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(5*time.Second)),
		)
		Expect(err).ToNot(HaveOccurred())

		connString, err = container.ConnectionString(ctx, "sslmode=disable")
		Expect(err).ToNot(HaveOccurred())

		conn, err := pgx.Connect(ctx, connString)
		Expect(err).ToNot(HaveOccurred())
		q = New(conn)

		DeferCleanup(func() {
			_ = container.Terminate(ctx)
			conn.Close(ctx)
		})
	})

	It("Queries should fail if table is not ready", func() {
		rotas, err := q.ListRotasByChannel(ctx, ListRotasByChannelParams{ChannelID: "foo", TeamID: "bar"})

		Expect(err.Error()).To(Equal("ERROR: relation \"rotas\" does not exist (SQLSTATE 42P01)"))
		Expect(len(rotas)).To(Equal(0))
	})

	It("Queries should not fail if table is ready", func() {
		err := Migrate(ctx, connString)
		Expect(err).ToNot(HaveOccurred())

		rotas, err := q.ListRotasByChannel(ctx, ListRotasByChannelParams{ChannelID: "foo", TeamID: "bar"})

		Expect(err).ToNot(HaveOccurred())
		Expect(len(rotas)).To(Equal(0))
	})
})
