package block

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/slack-go/slack"
)

var _ = Describe("Text", func() {
	It("Generates a default text", func() {
		t := NewDefaultText("Hello World")

		Expect(t.Type).To(Equal(slack.PlainTextType))
		Expect(t.Text).To(Equal("Hello World"))
		Expect(t.Emoji).To(BeFalse())
		Expect(t.Verbatim).To(BeFalse())
	})
})
