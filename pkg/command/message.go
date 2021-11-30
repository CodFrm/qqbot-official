package command

import "github.com/tencent-connect/botgo/dto"

type Message struct {
	*dto.Message
}

func (m *Message) Context() string {
	return m.Content
}

func (m *Message) User() string {
	return m.Author.ID
}

func (m *Message) Guild() string {
	return m.GuildID
}

func (m *Message) Channel() string {
	return m.ChannelID
}
