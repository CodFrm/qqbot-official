package middleware

import (
	"github.com/CodFrm/qqbot-official/internal/db"
	"github.com/CodFrm/qqbot-official/pkg/command"
	"github.com/tencent-connect/botgo/dto"
)

func Member(h func(m *dto.Member) (bool, error)) command.HandlerFunc {
	return func(ctx *command.Context) {
		m := &dto.Member{}
		err := db.GetOrSet("guild:userinfo:"+ctx.Message.Guild()+":"+ctx.Message.User(), m, func() (interface{}, error) {
			m, err := ctx.GuildMember()
			if err != nil {
				return nil, err
			}
			return m, nil
		}, 3600)
		if err != nil {
			ctx.Error(err)
			return
		}
		if ok, err := h(m); err != nil {
			ctx.Error(err)
			return
		} else if !ok {
			ctx.Abort()
		}
	}
}
