package block

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/slack-go/slack"
)

var _ = Describe("Input", func() {
	Describe("NewTextInput", func() {
		It("Generates text input with all input", func() {
			i := NewTextInput(TextInput{
				BlockID: "blockId",
				Label:   "label",
				Hint:    "hint",
			})

			Expect(i.Type).To(Equal(slack.MBTInput))
			Expect(i.BlockID).To(Equal("blockId"))
			Expect(i.Label.Type).To(Equal(slack.PlainTextType))
			Expect(i.Label.Text).To(Equal("label"))
			Expect(i.Element.ElementType()).To(Equal(slack.METPlainTextInput))

			accessory, ok := i.Element.(*slack.PlainTextInputBlockElement)
			Expect(ok).To(BeTrue())

			Expect(accessory.ActionID).To(Equal("blockId"))
			Expect(accessory.Placeholder.Type).To(Equal(slack.PlainTextType))
			Expect(accessory.Placeholder.Text).To(Equal("hint"))
		})
	})

	Describe("NewStaticSelect", func() {
		It("generates static select without options", func() {
			s := NewStaticSelect(StaticSelect{
				BlockID: "blockId",
				Label:   "label",
				InitialOption: StaticSelectOption{
					Text: "initial",
				},
				Options: []StaticSelectOption{},
			})

			Expect(s.Type).To(Equal(slack.MBTSection))
			Expect(s.BlockID).To(Equal("blockId"))
			Expect(s.Text.Type).To(Equal(slack.PlainTextType))
			Expect(s.Text.Text).To(Equal("label"))

			Expect(s.Accessory).ToNot(BeNil())
			Expect(s.Accessory.SelectElement.Type).To(Equal(slack.OptTypeStatic))
			Expect(s.Accessory.SelectElement.ActionID).To(Equal("blockId"))
			Expect(s.Accessory.SelectElement.InitialOption.Text.Text).To(Equal("initial"))
			Expect(s.Accessory.SelectElement.Options).To(BeEmpty())
		})

		It("generates static select without options", func() {
			s := NewStaticSelect(StaticSelect{
				BlockID: "blockId",
				Label:   "label",
				InitialOption: StaticSelectOption{
					Text: "initial",
				},
				Options: []StaticSelectOption{
					{
						Text: "option1",
					},
					{
						Text: "option2",
					},
				},
			})

			Expect(s.Type).To(Equal(slack.MBTSection))
			Expect(s.BlockID).To(Equal("blockId"))
			Expect(s.Text.Type).To(Equal(slack.PlainTextType))
			Expect(s.Text.Text).To(Equal("label"))

			Expect(s.Accessory).ToNot(BeNil())
			Expect(s.Accessory.SelectElement.Type).To(Equal(slack.OptTypeStatic))
			Expect(s.Accessory.SelectElement.ActionID).To(Equal("blockId"))
			Expect(s.Accessory.SelectElement.InitialOption.Text.Text).To(Equal("initial"))
			Expect(s.Accessory.SelectElement.Options).To(HaveLen(2))

			Expect(s.Accessory.SelectElement.Options[0].Text.Text).To(Equal("option1"))
			Expect(s.Accessory.SelectElement.Options[0].Value).To(Equal("option1"))
		})
	})
})
