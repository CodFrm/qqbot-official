package command

import (
	"context"

	"github.com/CodFrm/qqbot-official/internal/config"
	"github.com/CodFrm/qqbot-official/internal/service"
	utils2 "github.com/CodFrm/qqbot-official/internal/utils"
	"github.com/CodFrm/qqbot-official/pkg/command"
	"github.com/tencent-connect/botgo/dto"
)

type clockIn struct {
	svc  service.ClockIn
	user service.User
}

func NewClockIn(svc service.ClockIn, user service.User) *clockIn {
	return &clockIn{svc: svc, user: user}
}

func (c *clockIn) sleep(ctx *command.Context) {
	msg, err := c.svc.SleepClockIn(ctx.Message.Guild(), ctx.Message.User())
	if err != nil {
		ctx.ReplyText(err.Error())
	} else {
		ctx.ReplyArk(&dto.Ark{
			TemplateID: 37,
			KV: []*dto.ArkKV{
				{Key: "#PROMPT#", Value: "打卡成功"},
				{Key: "#METATITLE#", Value: "早睡打卡成功"},
				{Key: "#METASUBTITLE#", Value: msg},
				{Key: "#METACOVER#", Value: config.AppConfig.MsgUrl + "?action=images&name=sleep.jpg"},
			},
		})
	}
}

func (c *clockIn) getUp(ctx *command.Context) {
	msg, err := c.svc.GetUpClockIn(ctx.Message.Guild(), ctx.Message.User())
	if err != nil {
		ctx.ReplyText(err.Error())
		return
	}
	ctx.ReplyArk(&dto.Ark{
		TemplateID: 37,
		KV: []*dto.ArkKV{
			{Key: "#PROMPT#", Value: "打卡成功"},
			{Key: "#METATITLE#", Value: "早起打卡成功"},
			//{Key: "#METASUBTITLE#", Value: fmt.Sprintf("新的一天开始啦~您是第%d位起床的,昨晚睡了%d小时,奖励%d积分", n, hour/3600)},
			{Key: "#METASUBTITLE#", Value: msg},
			{Key: "#METACOVER#", Value: config.AppConfig.MsgUrl + "?action=images&name=getup.jpg"},
		},
	})
}

func (c *clockIn) clockIn(ctx *command.Context) {
	msg, err := c.svc.ClockIn(ctx.Message.Guild(), ctx.Message.User())
	if err != nil {
		ctx.ReplyText(err.Error())
		return
	}
	ctx.ReplyText(msg)
}

func (c *clockIn) getUpList(ctx *command.Context) {
	list, err := c.svc.GetUpList(ctx.Message.Guild())
	if err != nil {
		ctx.ReplyText(err.Error())
		return
	}
	msg := "以下人员未成功早起打卡:"
	for _, v := range list {
		msg += utils2.At(v)
	}
	ctx.ReplyText(msg)
}

func (c *clockIn) Register(ctx context.Context, cmd *command.Command) {
	cg := cmd.Group(command.AtMe())
	cg.Match("早睡打卡", c.sleep)
	cg.Match("早起打卡", c.getUp)
	cg.Match("打卡耻辱榜", c.getUpList)

	cg.Match("\\s打卡[\\s]?$", c.clockIn)
}
