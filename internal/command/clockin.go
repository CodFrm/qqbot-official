package command

import (
	"context"
	"strings"
	"time"

	"github.com/CodFrm/qqbot-official/internal/config"
	"github.com/CodFrm/qqbot-official/internal/middleware"
	"github.com/CodFrm/qqbot-official/internal/service"
	utils2 "github.com/CodFrm/qqbot-official/internal/utils"
	"github.com/CodFrm/qqbot-official/pkg/command"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/openapi"
)

type clockIn struct {
	svc        service.ClockIn
	channelSvc map[string]service.ChannelClockIn
	user       service.User
}

func NewClockIn(api openapi.OpenAPI, c *cron.Cron, svc service.ClockIn, user service.User) *clockIn {
	return &clockIn{
		svc:  svc,
		user: user,
		channelSvc: map[string]service.ChannelClockIn{
			"学习": service.NewChannelClockIn(c, user, api, "learn", &service.ChannelClockInOptions{Integral: 4}),
			"早起": service.NewChannelClockIn(c, user, api, "struggle", &service.ChannelClockInOptions{Integral: 3}),
		},
	}
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

func (c *clockIn) setNoticeChannel(ctx *command.Context) {
	svc, ok := c.channelSvc[ctx.Param("[功能]")]
	if !ok {
		ctx.ReplyText("没有对应的功能")
		return
	}
	if err := svc.SetNotice(ctx.Message.Guild(), ctx.Message.Channel(), strings.ReplaceAll(ctx.Param("[时间]"), "-", " "), ctx.Param("[标题]"), ctx.Param("[文案]")); err != nil {
		ctx.ReplyText(err.Error())
		return
	}
	ctx.ReplyText("设置成功")
}

func (c *clockIn) setClockChannel(ctx *command.Context) {
	svc, ok := c.channelSvc[ctx.Param("[功能]")]
	if !ok {
		ctx.ReplyText("没有对应的功能")
		return
	}
	if err := svc.SetClock(ctx.Message.Guild(), ctx.Message.Channel()); err != nil {
		ctx.ReplyText(err.Error())
		return
	}
	ctx.ReplyText("设置成功")
}

func (c *clockIn) struggle(ctx *command.Context) {
	svc := c.channelSvc["早起"]
	now := time.Now()
	if !(now.Hour() == 7 && now.Minute() > 30) {
		ctx.ReplyText("请在7:30-8:00之间打卡")
		return
	}
	if err := svc.ClockIn(ctx.Message.Guild(), ctx.Message.User()); err != nil {
		ctx.ReplyText(err.Error())
		return
	}
	ctx.ReplyText("早八人打卡成功,+3积分")
}

func (c *clockIn) learn(ctx *command.Context) {
	svc := c.channelSvc["学习"]
	if err := svc.ClockIn(ctx.Message.Guild(), ctx.Message.User()); err != nil {
		ctx.ReplyText(err.Error())
		return
	}
	ctx.ReplyText("学习打卡成功,+4积分")
}

func (c *clockIn) isLearn(ctx *command.Context) {
	svc := c.channelSvc["学习"]
	ok, err := svc.IsClock(ctx.Message.Guild(), ctx.Message.User())
	if err != nil {
		ctx.ReplyText(err.Error())
		return
	}
	if ok {
		ctx.ReplyText("已经打开过卡了,+4积分")
		return
	}
	ctx.ReplyText("暂未打卡,请分享学习软件进行打卡")
}

func (c *clockIn) Register(ctx context.Context, cmd *command.Command) {
	cg := cmd.Group()
	cg.Match("早睡打卡", c.sleep)
	cg.Match("早起打卡", c.getUp)
	cg.Match("打卡耻辱榜", c.getUpList)
	cg.Match("早八人打卡", c.struggle)
	cg.Match("学习打卡", c.isLearn)
	cg.Match("打卡", c.clockIn)
	cmd.Group(func(ctx *command.Context) {
		// 检测分享打卡
		if !strings.HasPrefix(ctx.Message.Content, "[分享]") {
			ctx.Abort()
			return
		}
		if !(strings.HasSuffix(ctx.Message.Content, "来自: 百词斩") ||
			strings.HasSuffix(ctx.Message.Content, "来自: 墨墨背单词")) {
			ctx.Abort()
		}
	}, c.learn)

	//NOTE: 后续可能调整为只允许频道主
	var admin = middleware.Member(func(m *dto.Member) (bool, error) {
		logrus.Infof("user: %v role: %v", m.Nick, m.Roles)
		for _, v := range m.Roles {
			if v == "4" || v == "2" || v == "5" {
				return true, nil
			}
		}
		return false, nil
	})

	cg.Match("设置通知频道 [功能] [时间] [标题] [文案]", admin, c.setNoticeChannel)
	cg.Match("设置打卡频道 [功能]", admin, c.setClockChannel)
}
