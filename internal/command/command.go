package command

import (
	"context"

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

	punish := NewPunish()

	Registers(context.Background(), command,
		punish,
		NewUtils(),
	)
}

func atReplay(c *command.Context, content string) {
	if c.IsAtMe() {
		c.ReplayText(content)
	}
}
