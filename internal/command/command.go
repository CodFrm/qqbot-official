package command

import (
	"context"

	"github.com/CodFrm/qqbot-official/internal/service"
	"github.com/CodFrm/qqbot-official/pkg/command"
)

type Register interface {
	Register(ctx context.Context, command *command.Command)
}

func Registers(ctx context.Context, command *command.Command, reg ...Register) {
	for _, v := range reg {
		v.Register(ctx, command)
	}
}

func InitCommand(command *command.Command) {

	punishSvc := service.NewPunish()

	Registers(context.Background(), command,
		NewPunish(),
		newIdentity(punishSvc),
		NewUtils(),
	)
}

func atReplay(c *command.Context, content string) {
	if c.IsAtMe() {
		c.ReplyText(content)
	}
}
