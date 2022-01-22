package command

import (
	"context"
	"math"
	"strconv"

	"github.com/CodFrm/qqbot-official/internal/config"
	"github.com/CodFrm/qqbot-official/internal/db"
	"github.com/CodFrm/qqbot-official/internal/service"
	"github.com/CodFrm/qqbot-official/internal/utils/api"
	"github.com/CodFrm/qqbot-official/pkg/command"
	utils2 "github.com/CodFrm/qqbot/utils"
	"github.com/tencent-connect/botgo/dto"
)

type user struct {
	svc    service.User
	punish service.Punish
}

func newUser(svc service.User, punish service.Punish) *user {
	return &user{svc: svc, punish: punish}
}

func (i *user) info(c *command.Context) {
	punishLevel, err := i.punish.PunishLevel(c.Message.Guild(), c.Message.User())
	if err != nil {
		c.ReplyText("警告:" + err.Error())
		return
	}
	integral, err := i.svc.Integral(c.Message.Guild(), c.Message.User())
	if err != nil {
		c.ReplyText("积分:" + err.Error())
		return
	}
	arkList := []*dto.ArkObj{
		{
			ObjKV: []*dto.ArkObjKV{
				{Key: "desc", Value: c.Message.Author.Username + " 在本频道的信息:"},
			},
		}, {
			ObjKV: []*dto.ArkObjKV{
				{Key: "desc", Value: "用户积分:" + strconv.FormatInt(integral, 10) + " 用户id:" + c.Message.User()},
			},
		}, {
			ObjKV: []*dto.ArkObjKV{
				{Key: "desc", Value: "警告等级:" + strconv.FormatInt(punishLevel, 10)},
			},
		},
	}
	if punishLevel != 0 {
		arkList = append(arkList, &dto.ArkObj{ObjKV: []*dto.ArkObjKV{
			{Key: "desc", Value: "警告人员不可设置用户组"},
		}})
		c.ReplyArk(&dto.Ark{
			TemplateID: 23,
			KV: []*dto.ArkKV{
				{Key: "#DESC#", Value: "用户信息"},
				{Key: "#PROMPT#", Value: "用户信息"},
				{
					Key: "#LIST#",
					Obj: arkList,
				},
			},
		})
		return
	}
	next := &dto.ArkObjKV{Key: "desc", Value: "下面为可申请的用户组,点击相应链接可切换用户组,下一晋级用户组为:"}
	arkList = append(arkList, &dto.ArkObj{ObjKV: []*dto.ArkObjKV{
		next,
	}})
	session := utils2.RandStringRunes(10)
	if err := db.Put("session:"+session, []byte(c.Message.User()), 600); err != nil {
		c.ReplyText("session:" + err.Error())
		return
	}
	list, err := api.NewGuildApi(c.OpenApi()).UserGroup(c.Message.Guild())
	n := -1
	var userList []*dto.ArkObj
	for i := 0; i < len(list); i++ {
		v := list[len(list)-i-1]
		if v.Hoist == 1 {
			continue
		}
		n++
		if integral < 10*int64(math.Pow(2, float64(n)))-10 {
			next.Value += v.Name + " 还差:" + strconv.FormatInt(10*int64(math.Pow(2, float64(n)))-10-integral, 10) + "积分"
			break
		}
		userList = append([]*dto.ArkObj{{
			ObjKV: []*dto.ArkObjKV{
				{Key: "desc", Value: v.Name},
				{Key: "link", Value: config.AppConfig.MsgUrl + "?action=setUserRole&session=" +
					session + "&guild=" + c.Message.Guild() + "&user=" + c.Message.User() + "&role=" + string(v.ID)},
			},
		}}, userList...)
	}
	arkList = append(arkList, userList...)
	c.ReplyArk(&dto.Ark{
		TemplateID: 23,
		KV: []*dto.ArkKV{
			{Key: "#DESC#", Value: "用户信息"},
			{Key: "#PROMPT#", Value: "用户信息"},
			{
				Key: "#LIST#",
				Obj: arkList,
			},
		},
	})
}

func (i *user) Register(ctx context.Context, cmd *command.Command) {
	cg := cmd.Group()
	cg.Match("我的信息", i.info)
}
