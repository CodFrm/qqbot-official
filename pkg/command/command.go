package command

import (
	"github.com/sirupsen/logrus"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/openapi"
)

type Command struct {
	bot          *dto.User
	api          openapi.OpenAPI
	defaultGroup *Group
}

func NewCommand(bot *dto.User, api openapi.OpenAPI) *Command {
	return &Command{
		bot:          bot,
		api:          api,
		defaultGroup: newGroup(),
	}
}

func (c *Command) Use(handler ...HandlerFunc) {
	c.defaultGroup.Use(handler...)
}

func (c *Command) Group(handler ...HandlerFunc) *Group {
	return c.defaultGroup.Group(handler...)
}

// Match 命令 [参数1] [参数2]
func (c *Command) Match(command string, handler ...HandlerFunc) {
	c.defaultGroup.Match(command, handler...)
}

// AtMeMatch Match 命令 [参数1] [参数2]
func (c *Command) AtMeMatch(command string, handler ...HandlerFunc) {
	c.defaultGroup.AtMeMatch(command, handler...)
}

func (c *Command) AtMe(handler ...HandlerFunc) {
	c.defaultGroup.AtMe(handler...)
}

func (c *Command) MessageHandler(data *dto.Message) error {
	if data.DirectMessage {
		return nil
	}
	// 判断是否为艾特我
	ctx := createContext(c.bot, c.api, data)
	c.defaultGroup.handle(ctx)
	logrus.Infof("message: guild(%v) channel(%v) user(%+v): %v", data.GuildID, data.ChannelID, data.Author, data.Content)
	return nil
}
