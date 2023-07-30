package block

import "github.com/slack-go/slack"

func NewHeader(text string) *slack.SectionBlock {
	return &slack.SectionBlock{
		Type: slack.MBTHeader,
		Text: NewDefaultText(text),
	}
}
