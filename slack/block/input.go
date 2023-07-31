package block

import "github.com/slack-go/slack"

type TextInput struct {
	BlockID string
	Label   string
	Hint    string
	Value   string
}

func NewTextInput(input TextInput) *slack.InputBlock {
	return &slack.InputBlock{
		Type:    slack.MBTInput,
		BlockID: input.BlockID,
		Element: &slack.PlainTextInputBlockElement{
			Type:         slack.METPlainTextInput,
			ActionID:     input.BlockID,
			Placeholder:  NewDefaultText(input.Hint),
			InitialValue: input.Value,
		},
		Label: NewDefaultText(input.Label),
	}
}

type StaticSelect struct {
	BlockID       string
	Label         string
	InitialOption StaticSelectOption
	Options       []StaticSelectOption
}

type StaticSelectOption struct {
	Text string
}

func NewStaticSelect(input StaticSelect) *slack.SectionBlock {
	options := []*slack.OptionBlockObject{}
	for _, option := range input.Options {
		options = append(options, staticSelectOption(option))
	}

	return &slack.SectionBlock{
		Type:    slack.MBTSection,
		BlockID: input.BlockID,
		Text:    NewDefaultText(input.Label),
		Accessory: slack.NewAccessory(
			&slack.SelectBlockElement{
				Type:          slack.OptTypeStatic,
				InitialOption: staticSelectOption(input.InitialOption),
				ActionID:      input.BlockID,
				Options:       options,
			},
		),
	}
}

func staticSelectOption(option StaticSelectOption) *slack.OptionBlockObject {
	return &slack.OptionBlockObject{
		Text:  NewDefaultText(option.Text),
		Value: option.Text,
	}
}
