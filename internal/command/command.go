package command

import (
	"context"

	"github.com/CodFrm/qqbot-official/internal/service"
	"github.com/CodFrm/qqbot-official/pkg/command"
	"github.com/robfig/cron/v3"
	"github.com/tencent-connect/botgo/openapi"
)

type Register interface {
	Register(ctx context.Context, command *command.Command)
}

func Registers(ctx context.Context, command *command.Command, reg ...Register) {
	for _, v := range reg {
		v.Register(ctx, command)
	}
}

func InitCommand(command *command.Command, api openapi.OpenAPI) {

	c := cron.New()

	punishSvc := service.NewPunish()
	userSvc := service.NewUser()
	clockinSvc := service.NewClockIn(c, userSvc, api)

	Registers(context.Background(), command,
		NewPunish(),
		newUser(userSvc, punishSvc),
		NewClockIn(api, c, clockinSvc, userSvc),
		NewUtils(),
	)

	c.Start()
}
