package command

import (
	"context"
	"fmt"

	"github.com/CodFrm/qqbot-official/internal/middleware"
	"github.com/CodFrm/qqbot-official/pkg/command"
	"github.com/tencent-connect/botgo/dto"
)

type punish struct {
}

func NewPunish() *punish {
	return &punish{}
}

func (p *punish) punish(ctx *command.Context) {
	ctx.ReplayText("警告！")
}

func (p *punish) Register(ctx context.Context, cmd *command.Command) {
	cg := cmd.Group(command.AtMe(), middleware.Member(func(m *dto.Member) (bool, error) {
		if len(m.Roles) == 0 {
			return false, nil
		}
		fmt.Println(m.Roles)
		return true, nil
	}))
	cg.Match("警告", p.punish)
}
