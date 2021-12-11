package command

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/CodFrm/qqbot-official/internal/db"
	"github.com/CodFrm/qqbot-official/internal/middleware"
	utils2 "github.com/CodFrm/qqbot-official/internal/utils"
	"github.com/CodFrm/qqbot-official/internal/utils/api"
	"github.com/CodFrm/qqbot-official/pkg/command"
	"github.com/sirupsen/logrus"
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
		ctx.ReplyText(err.Error())
		return
	}
	if num > 40 {
		ctx.ReplyText("今天已经处理够多人了")
		return
	}
	for _, user := range ctx.Message.Mentions() {
		if user.ID == ctx.Bot().ID || user.ID == ctx.Message.User() {
			continue
		}
		num, err := db.Incr(fmt.Sprintf("guild:punish:%v:%v", ctx.Message.Guild(), user.ID), 1, 604800)
		member += utils2.At(user.ID)
		if err != nil {
			member += "错误:" + err.Error() + "\n"
			continue
		}
		m, err := api.NewGuildApi(ctx.OpenApi()).GuildMember(ctx.Message.Guild(), user.ID)
		if err != nil {
			member += "错误:" + err.Error() + "\n"
			continue
		}
		flag := false
		for _, v := range m.Roles {
			if v == "4" || v == "2" || v == "5" {
				flag = true
				break
			}
		}
		if flag {
			member += "管理员间无法警告,请反馈给频道主\n"
			continue
		}
		switch num {
		case 1, 2:
			member += "警告一次"
			var punishRole *dto.Role
			roles, err := api.NewGuildApi(ctx.OpenApi()).UserGroup(ctx.Message.Guild())
			if err != nil {
				member += " " + err.Error()
			} else {
				for _, role := range roles {
					if strings.Index(role.Name, "警告") != -1 {
						punishRole = role
						break
					}
				}
				if err := api.NewGuildApi(ctx.OpenApi()).SetSignalRole(ctx.Message.Guild(), user.ID, punishRole.ID); err != nil {
					member += " " + err.Error()
				} else {
					member += "并设置" + punishRole.Name + "用户组"
				}
				break
			}
			if err := api.NewGuildApi(ctx.OpenApi()).SetSignalRole(ctx.Message.Guild(), user.ID, punishRole.ID); err != nil {
				member += " " + err.Error()
			} else {
				member += "并移除所有用户组"
			}
		case 3:
			member += "踢出频道"
			if err := ctx.OpenApi().DeleteGuildMember(context.Background(), ctx.Message.Guild(), user.ID); err != nil {
				member += " " + err.Error()
			}
		case 4:
			member += "拉黑此人(暂时需要管理员手动操作拉黑)"
			g, err := ctx.Guild()
			if err != nil {
				member += " " + err.Error()
			} else {
				member += utils2.At(g.OwnerID)
			}
		default:
			g, err := ctx.Guild()
			if err != nil {
				member += err.Error()
			} else {
				if g.OwnerID == user.ID {
					member += "这人咋还在？请求最高权限:"
				} else {
					member += "这人咋还在？请求最高权限:"
				}
				member += utils2.At(g.OwnerID)
			}
		}
		member += "\n"
	}
	ctx.ReplyText("对以下成员做出处理:\n" + member)
}

func (p *punish) Register(ctx context.Context, cmd *command.Command) {
	cg := cmd.Group(middleware.Member(func(m *dto.Member) (bool, error) {
		logrus.Infof("user: %v role: %v", m.Nick, m.Roles)
		for _, v := range m.Roles {
			if v == "4" || v == "2" || v == "5" {
				return true, nil
			}
		}
		return false, nil
	}))
	cg.AtMeMatch("警告", p.punish)
}
