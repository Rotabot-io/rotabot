package block

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/slack-go/slack"
)

var _ = DescribeTable("OverflowSection",
	func(i OverflowSection, blockId string, opts []slack.OptionBlockObject) {
		o := NewOverflowSectionElement(i)

		Expect(o.Type).To(Equal(slack.MBTSection))
		Expect(o.BlockID).To(Equal(i.ElementID))
		Expect(o.Fields).To(BeNil())

		Expect(o.Text.Text).To(Equal(i.ElementName))
		Expect(o.Text.Type).To(Equal(slack.PlainTextType))

		Expect(o.Accessory.OverflowElement.Type).To(Equal(slack.METOverflow))
		Expect(o.Accessory.OverflowElement.ActionID).To(Equal(i.SectionName))
		Expect(len(o.Accessory.OverflowElement.Options)).To(Equal(len(opts)))

		for inx, option := range opts {
			Expect(o.Accessory.OverflowElement.Options[inx].Description).To(BeNil())
			Expect(o.Accessory.OverflowElement.Options[inx].URL).To(BeEmpty())
			Expect(o.Accessory.OverflowElement.Options[inx].Text.Type).To(Equal(slack.PlainTextType))
			Expect(o.Accessory.OverflowElement.Options[inx].Text.Text).To(Equal(option.Text.Text))
			Expect(o.Accessory.OverflowElement.Options[inx].Value).To(Equal(option.Value))
		}
	},
	Entry(
		"No Actions",
		OverflowSection{
			ElementID:   "ElementID",
			ElementName: "ElementName",
			SectionName: "SectionName",
		},
		"BlockID_ElementID",
		[]slack.OptionBlockObject{},
	),
	Entry(
		"One action",
		OverflowSection{
			ElementID:   "ElementID",
			ElementName: "ElementName",
			SectionName: "SectionName",
			Actions: []OverflowAction{
				{
					Action: "action",
					Name:   "Name",
				},
			},
		},
		"BlockID_ElementID",
		[]slack.OptionBlockObject{
			{
				Text:  NewDefaultText("Name"),
				Value: "action",
			},
		}),
	Entry(
		"Many Actions, actions are ordered",
		OverflowSection{
			ElementID:   "ElementID",
			ElementName: "ElementName",
			SectionName: "SectionName",
			Actions: []OverflowAction{
				{
					Action: "action#1",
					Name:   "Name#1",
				},
				{
					Action: "action#2",
					Name:   "Name#2",
				},
				{
					Action: "action#3",
					Name:   "Name#3",
				},
			},
		},
		"BlockID_ElementID",
		[]slack.OptionBlockObject{
			{
				Text:  NewDefaultText("Name#1"),
				Value: "action#1",
			},
			{
				Text:  NewDefaultText("Name#2"),
				Value: "action#2",
			},
			{
				Text:  NewDefaultText("Name#3"),
				Value: "action#3",
			},
		}),
)
