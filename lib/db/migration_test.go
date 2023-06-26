package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	. "github.com/onsi/ginkgo/v2"

	. "github.com/onsi/gomega"
	"github.com/rotabot-io/rotabot/internal"
)

var _ = Describe("Migration", func() {
	var ctx context.Context
	var connString string
	var q *Queries

	BeforeEach(func() {
		ctx = context.Background()

		container, err := internal.RunContainer(ctx)
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

	It("Queries should fail if table is not ready", func() {
		rotas, err := q.ListRotas(ctx)

		Expect(err.Error()).To(Equal("ERROR: relation \"rotas\" does not exist (SQLSTATE 42P01)"))
		Expect(len(rotas)).To(Equal(0))
	})

	It("Queries should not fail if table is ready", func() {
		err := Migrate(ctx, connString)
		Expect(err).ToNot(HaveOccurred())

		rotas, err := q.ListRotas(ctx)

		Expect(err).ToNot(HaveOccurred())
		Expect(len(rotas)).To(Equal(0))
	})
})
