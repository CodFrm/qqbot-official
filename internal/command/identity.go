package command

import (
	"context"

	"github.com/CodFrm/qqbot-official/pkg/command"
)

type identity struct {
}

func newIdentity() *identity {
	return &identity{}
}

func (i *identity) Register(ctx context.Context, cmd *command.Command) {

}
