package command

import (
	"regexp"

	"github.com/tencent-connect/botgo/dto"
)

type Message struct {
	*dto.Message
}

func (m *Message) Mentions() []*dto.User {
	// 处理艾特信息,现在官方不知道为什么只能收到一个
	reg := regexp.MustCompile("<@!(\\d+)>")
	all := reg.FindAllStringSubmatch(m.Content, -1)
	m.Message.Mentions = make([]*dto.User, len(all))
	for n, v := range all {
		m.Message.Mentions[n] = &dto.User{
			ID: v[1],
		}
	}
	return m.Message.Mentions
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
