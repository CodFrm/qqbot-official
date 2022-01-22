package command

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/CodFrm/qqbot-official/pkg/command"
)

type utils struct {
}

func NewUtils() *utils {
	return &utils{}
}

func (u *utils) Register(ctx context.Context, cmd *command.Command) {
	cmd.Match("摇色子|摇骰子", func(ctx *command.Context) {
		ctx.ReplyText(fmt.Sprintf("%v", rand.Int31n(6)+1))
	})
}
