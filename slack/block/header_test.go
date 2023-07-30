package block

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/slack-go/slack"
)

var _ = Describe("Header", func() {
	It("Generates Slack Header text", func() {
		h := NewHeader("Hello World")

		Expect(h.Type).To(Equal(slack.MBTHeader))
		Expect(h.Text.Type).To(Equal(slack.PlainTextType))
		Expect(h.Text.Text).To(Equal("Hello World"))
	})
})
