package block

import (
	"github.com/slack-go/slack"
)

//	 This file tries to hide the complexities of the slack api regarding texts, i
//	 The following is the documentation copied from the slack website
//
//		return &slack.TextBlockObject{
//	   The formatting to use for this text object. Can be one of plain_textor mrkdwn.
//			Type:  slack.MarkdownType || slack.PlainTextType,
//
//	   The text for the block. This field accepts any of the standard text formatting markup when type is mrkdwn. The minimum length is 1 and maximum length is 3000 characters.
//			Text:  text,
//
//	   Indicates whether emojis in a text field should be escaped into the colon emoji format. This field is only usable when type is plain_text.
//			Emoji: boolean,
//
//			When set to false (as is default) URLs will be auto-converted into links, conversation names will be link-ified, and certain mentions will be automatically parsed.
//			Using a value of true will skip any preprocessing of this nature, although you can still include manual parsing strings. This field is only usable when type is mrkdwn.
//			Verbatim: boolean,
//		}
//
// source: https://api.slack.com/reference/block-kit/composition-objects#text

func NewDefaultText(text string) *slack.TextBlockObject {
	return &slack.TextBlockObject{
		Type: slack.PlainTextType,
		Text: text,
	}
}
