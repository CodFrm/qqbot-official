package command

import (
	"context"
	"strconv"

	"github.com/CodFrm/qqbot-official/internal/config"
	"github.com/CodFrm/qqbot-official/internal/db"
	"github.com/CodFrm/qqbot-official/internal/service"
	"github.com/CodFrm/qqbot-official/internal/utils/api"
	"github.com/CodFrm/qqbot-official/pkg/command"
	utils2 "github.com/CodFrm/qqbot/utils"
	"github.com/tencent-connect/botgo/dto"
)

type identity struct {
	punish service.Punish
}

func newIdentity(punish service.Punish) *identity {
	return &identity{punish: punish}
}

func (i *identity) info(c *command.Context) {
	punishLevel, err := i.punish.PunishLevel(c.Message.Guild(), c.Message.User())
	if err != nil {
		c.ReplyText(err.Error())
		return
	}
	arkList := []*dto.ArkObj{
		{
			ObjKV: []*dto.ArkObjKV{
				{Key: "desc", Value: c.Message.Author.Username + " 在本频道的信息:"},
			},
		}, {
			ObjKV: []*dto.ArkObjKV{
				{Key: "desc", Value: "警告等级:" + strconv.FormatInt(punishLevel, 10) + " 用户id:" + c.Message.User()},
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
	arkList = append(arkList, &dto.ArkObj{ObjKV: []*dto.ArkObjKV{
		{Key: "desc", Value: "下面为可申请的用户组,点击相应链接可切换用户组"},
	}})
	session := utils2.RandStringRunes(10)
	if err := db.Put("session:"+session, []byte(c.Message.User()), 600); err != nil {
		c.ReplyText(err.Error())
		return
	}
	list, err := api.NewGuildApi(c.OpenApi()).UserGroup(c.Message.Guild())
	for _, v := range list {
		if v.Hoist == 1 {
			continue
		}
		arkList = append(arkList, &dto.ArkObj{
			ObjKV: []*dto.ArkObjKV{
				{Key: "desc", Value: v.Name},
				{Key: "link", Value: config.AppConfig.MsgUrl + "?action=setUserRole&session=" +
					session + "&guild=" + c.Message.Guild() + "&user=" + c.Message.User() + "&role=" + string(v.ID)},
			},
		})
	}
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

func (i *identity) Register(ctx context.Context, cmd *command.Command) {
	cg := cmd.Group(command.AtMe())
	cg.Match("我的信息", i.info)
}
