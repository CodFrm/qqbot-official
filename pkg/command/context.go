package command

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/openapi"
)

type Context struct {
	ctx       context.Context
	isAborted bool

	Message *Message

	bot *dto.User
	api openapi.OpenAPI

	param map[string]string
}

func createContext(bot *dto.User, api openapi.OpenAPI, data *dto.Message) *Context {
	return &Context{
		ctx:       context.Background(),
		isAborted: false,
		Message:   &Message{Message: data},
		bot:       bot,
		api:       api,
		param:     map[string]string{},
	}
}

func (c *Context) Bot() *dto.User {
	return c.bot
}

func (c *Context) IsAborted() bool {
	return c.isAborted
}

func (c *Context) Abort() {
	c.isAborted = true
}

func (c *Context) ReplyText(content string) {
	if _, err := c.api.PostMessage(c.ctx, c.Message.Channel(), &dto.MessageToCreate{
		Content: "<@!" + c.Message.User() + ">" + content,
		MsgID:   c.Message.ID,
	}); err != nil {
		logrus.Errorf("post message: %v", err)
	}
}

func (c *Context) Error(err error) {
	logrus.Errorf("handle error: %+v", errors.WithStack(err))
	//if _, err := c.api.PostMessage(c.ctx, c.Message.Channel(), &dto.MessageToCreate{
	//	Content: "<@!" + c.Message.User() + ">" + err.Error(),
	//	MsgID:   c.Message.ID,
	//}); err != nil {
	//	logrus.Errorf("post message: %v", err)
	//}
	c.Abort()
}

func (c *Context) Param(k string) string {
	return c.param[k]
}
func (c *Context) setParam(k, v string) {
	c.param[k] = v
}

func (c *Context) OpenApi() openapi.OpenAPI {
	return c.api
}

func (c *Context) GuildMember() (*dto.Member, error) {
	return c.api.GuildMember(context.Background(), c.Message.Guild(), c.Message.User())
}

func (c *Context) Guild() (*dto.Guild, error) {
	return c.api.Guild(context.Background(), c.Message.Guild())
}

func (c *Context) IsAtMe() bool {
	for _, v := range c.Message.Mentions() {
		if v.ID == c.bot.ID {
			return true
		}
	}
	return false
}

func (c *Context) ReplyArk(ark *dto.Ark) {
	_, err := c.api.PostMessage(c.ctx, c.Message.Channel(), &dto.MessageToCreate{
		Ark:   ark,
		MsgID: c.Message.ID,
	})
	if err != nil {
		logrus.Errorf("post message: %v", err)
	}
}
