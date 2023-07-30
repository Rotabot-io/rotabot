package block

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/slack-go/slack"
)

var _ = Describe("Button", func() {
	It("Generates Slack Button", func() {
		b := NewButton(Button{
			ActionID: "awesome_button",
			Text:     "Click me!",
		})

		Expect(b.Type).To(Equal(slack.METButton))
		Expect(b.Text.Text).To(Equal("Click me!"))
		Expect(b.ActionID).To(Equal("awesome_button"))
	})
})
