package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	. "github.com/onsi/ginkgo/v2"

	. "github.com/onsi/gomega"
	"github.com/rotabot-io/rotabot/internal"
)

var _ = Describe("Repository", func() {
	var ctx context.Context
	var connString string
	var conn *pgx.Conn

	BeforeEach(func() {
		var err error
		ctx = context.Background()

		container, err := internal.RunContainer(ctx)
		Expect(err).ToNot(HaveOccurred())

		err = Migrate(ctx, container.ConnectionString())
		Expect(err).ToNot(HaveOccurred())

		connString = container.ConnectionString()

		conn, err = pgx.Connect(ctx, connString)
		Expect(err).ToNot(HaveOccurred())

		DeferCleanup(func() {
			_ = container.Terminate(ctx)
			conn.Close(ctx)
		})
	})

	It("Should create rota if id is null", func() {
		tx, err := conn.Begin(ctx)
		Expect(err).ToNot(HaveOccurred())

		id, err := CreateOrUpdateRota(ctx, tx, CreateOrUpdateRotaParams{
			ChannelID: "foo",
			TeamID:    "bar",
			Name:      "baz",
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(id).ToNot(BeEmpty())
	})

	It("Should fail to create two identical rotas", func() {
		tx, err := conn.Begin(ctx)
		Expect(err).ToNot(HaveOccurred())

		req := CreateOrUpdateRotaParams{
			ChannelID: "foo",
			TeamID:    "bar",
			Name:      "baz",
		}
		_, err = CreateOrUpdateRota(ctx, tx, req)
		Expect(err).ToNot(HaveOccurred())

		_, err = CreateOrUpdateRota(ctx, tx, req)
		Expect(err).To(HaveOccurred())
		Expect(err).To(Equal(ErrAlreadyExists))
	})

	It("Should update rota if id is null", func() {
		tx, err := conn.Begin(ctx)
		Expect(err).ToNot(HaveOccurred())

		id, err := CreateOrUpdateRota(ctx, tx, CreateOrUpdateRotaParams{
			ChannelID: "foo",
			TeamID:    "bar",
			Name:      "baz",
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(id).ToNot(BeEmpty())

		updated, err := CreateOrUpdateRota(ctx, tx, CreateOrUpdateRotaParams{
			RotaID: id,
			Name:   "bazbaz",
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(id).To(Equal(updated))

		rota, err := New(tx).FindRotaByID(ctx, id)
		Expect(err).ToNot(HaveOccurred())
		Expect(rota.Name).To(Equal("bazbaz"))
	})

	It("Should fail to update something it does not exist", func() {
		tx, err := conn.Begin(ctx)
		Expect(err).ToNot(HaveOccurred())

		_, err = CreateOrUpdateRota(ctx, tx, CreateOrUpdateRotaParams{
			RotaID: "not_found",
			Name:   "bazbaz",
		})
		Expect(err).To(HaveOccurred())
		Expect(err).To(Equal(ErrNotFound))
	})
})
