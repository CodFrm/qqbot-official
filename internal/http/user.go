package http

import (
	"context"

	"github.com/CodFrm/qqbot-official/internal/utils/api"
	"github.com/gin-gonic/gin"
	"github.com/tencent-connect/botgo/dto"
)

type user struct {
	api *api.GuildApi
}

func NewUser(api *api.GuildApi) *user {
	return &user{api: api}
}

func (u *user) setRole(c *gin.Context) {
	handle(c, func() interface{} {
		guild := c.Query("guild")
		user := c.Query("user")
		role := c.Query("role")
		if err := u.api.SetSignalRole(guild, user, dto.RoleID(role)); err != nil {
			return err
		}
		return "用户组切换成功"
	})
}

func (u *user) Registry(ctx context.Context, r *gin.Engine) {
	actionMap["setUserRole"] = u.setRole
}
