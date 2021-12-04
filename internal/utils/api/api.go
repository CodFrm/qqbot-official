package api

import (
	"context"

	"github.com/CodFrm/qqbot-official/internal/db"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/openapi"
)

// GuildApi 对openapi封装,增加了缓存
type GuildApi struct {
	api openapi.OpenAPI
}

func NewGuildApi(api openapi.OpenAPI) *GuildApi {
	return &GuildApi{api: api}
}

func (g *GuildApi) OpenApi() openapi.OpenAPI {
	return g.api
}

func (g *GuildApi) UserGroup(guild string) ([]*dto.Role, error) {
	var list []*dto.Role
	if err := db.GetOrSet("guild:cache:role:"+guild, &list, func() (interface{}, error) {
		l, err := g.api.Roles(context.Background(), guild)
		if err != nil {
			return nil, err
		}
		for _, v := range l.Roles {
			list = append(list, v)
		}
		return list, nil
	}, 3600); err != nil {
		return nil, err
	}
	return list, nil
}

func (g *GuildApi) SetSignalRole(guild, user string, role dto.RoleID) error {
	m, err := g.api.GuildMember(context.Background(), guild, user)
	if err != nil {
		return err
	}
	list, err := g.UserGroup(guild)
	if err != nil {
		return err
	}
	roleMap := map[string]*dto.Role{}
	for _, v := range list {
		roleMap[string(v.ID)] = v
	}
	for _, v := range m.Roles {
		r, ok := roleMap[v]
		if ok && r.Hoist != 1 && role != dto.RoleID(v) {
			g.api.MemberDeleteRole(context.Background(), guild, dto.RoleID(v), user, nil)
		}
	}
	if role == "" {
		return nil
	}
	return g.api.MemberAddRole(context.Background(), guild, role, user, nil)
}
