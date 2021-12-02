package command

import (
	"context"
	"fmt"
	"time"

	"github.com/CodFrm/qqbot-official/internal/db"
	"github.com/CodFrm/qqbot-official/internal/middleware"
	"github.com/CodFrm/qqbot-official/pkg/command"
	"github.com/CodFrm/qqbot-official/pkg/guild"
	"github.com/tencent-connect/botgo/dto"
)

type punish struct {
}

func NewPunish() *punish {
	return &punish{}
}

func (p *punish) punish(ctx *command.Context) {
	member := ""
	num, err := db.Incr("guild:punish:user:num:"+time.Now().Format("2006:01:02")+":"+ctx.Message.User(), int64(len(ctx.Message.Mentions())), 3600*24)
	if err != nil {
		ctx.ReplayText(err.Error())
		return
	}
	if num > 40 {
		ctx.ReplayText("今天已经处理够多人了")
		return
	}
	for _, v := range ctx.Message.Mentions() {
		if v.ID == ctx.Bot().ID || v.ID == ctx.Message.User() {
			continue
		}
		num, err := db.Incr(fmt.Sprintf("guild:punish:%v:%v", ctx.Message.Guild(), v.ID), 1, 604800)
		member += guild.At(v.ID)
		if err != nil {
			member += "错误:" + err.Error() + "\n"
			continue
		}
		switch num {
		case 1:
			member += "警告一次"
		case 2:
			member += "踢出频道"
			if err := ctx.OpenApi().DeleteGuildMember(context.Background(), ctx.Message.Guild(), v.ID); err != nil {
				member += " " + err.Error()
			}
		case 3:
			member += "拉黑此人(暂时需要管理员手动操作拉黑)"
			g, err := ctx.Guild()
			if err != nil {
				member += " " + err.Error()
			} else {
				member += guild.At(g.OwnerID)
			}
		default:
			g, err := ctx.Guild()
			if err != nil {
				member += err.Error()
			} else {
				if g.OwnerID == v.ID {
					member += "这人咋还在？请求最高权限:"
				} else {
					member += "这人咋还在？请求最高权限:"
				}
				member += guild.At(g.OwnerID)
			}
		}
		member += "\n"
	}
	atReplay(ctx, guild.At(ctx.Message.User())+"对以下成员做出处理:\n"+member)
}

func (p *punish) Register(ctx context.Context, cmd *command.Command) {
	cg := cmd.Group(middleware.Member(func(m *dto.Member) (bool, error) {
		if len(m.Roles) == 0 {
			return false, nil
		}
		fmt.Println(m.Roles)
		return true, nil
	}))
	cg.Match("警告", p.punish)
}
