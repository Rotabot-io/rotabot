package block

import (
	"github.com/slack-go/slack"
)

type OverflowSection struct {
	ElementID   string
	ElementName string
	SectionName string
	Actions     []OverflowAction
}

func NewOverflowSectionElement(input OverflowSection) *slack.SectionBlock {
	options := []*slack.OptionBlockObject{}
	for _, action := range input.Actions {
		options = append(options, buildOverflowAction(action))
	}
	return &slack.SectionBlock{
		Type:      slack.MBTSection,
		BlockID:   input.ElementID,
		Text:      NewDefaultText(input.ElementName),
		Fields:    nil,
		Accessory: slack.NewAccessory(slack.NewOverflowBlockElement(input.SectionName, options...)),
	}
}

type OverflowAction struct {
	Action string
	Name   string
}

// buildOverflowAction creates a clickable option for the overflow section.
// Please refer to https://api.slack.com/reference/block-kit/composition-objects#option for more context
func buildOverflowAction(input OverflowAction) *slack.OptionBlockObject {
	return &slack.OptionBlockObject{
		Text:  NewDefaultText(input.Name),
		Value: input.Action,
	}
}
