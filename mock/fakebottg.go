package mock

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type FakeBot struct {
	SentMessages []telegram.Chattable
}

func NewFakeBot() *FakeBot {
	return &FakeBot{}
}

func (f *FakeBot) Send(c telegram.Chattable) (telegram.Message, error) {
	f.SentMessages = append(f.SentMessages, c)
	return telegram.Message{}, nil
}

func (f *FakeBot) LastMessageText() string {
	if len(f.SentMessages) == 0 {
		return ""
	}
	msg, ok := f.SentMessages[len(f.SentMessages)-1].(telegram.MessageConfig)
	if !ok {
		return ""
	}

	return msg.Text
}
