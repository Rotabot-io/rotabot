package block

import "github.com/slack-go/slack"

type Button struct {
	ActionID string
	Text     string
}

func NewButton(b Button) *slack.ButtonBlockElement {
	return &slack.ButtonBlockElement{
		Type:     slack.METButton,
		ActionID: b.ActionID,
		Text:     NewDefaultText(b.Text),
	}
}
