package main

import (
	"context"
	"log"
	"time"

	command2 "github.com/CodFrm/qqbot-official/internal/command"
	"github.com/CodFrm/qqbot-official/internal/config"
	"github.com/CodFrm/qqbot-official/internal/db"
	"github.com/CodFrm/qqbot-official/pkg/command"
	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/token"
	"github.com/tencent-connect/botgo/websocket"
)

func main() {
	if err := config.Init("config.yaml"); err != nil {
		log.Fatalf("init config: %v", err)
	}
	if err := db.InitDatabase(); err != nil {
		log.Fatalf("init database: %v", err)
	}

	token := token.BotToken(config.AppConfig.AppID, config.AppConfig.AccessToken)
	api := botgo.NewOpenAPI(token).WithTimeout(3 * time.Second)
	ctx := context.Background()

	ws, err := api.WS(ctx, nil, "")
	if err != nil {
		log.Fatalf("%+v, err:%v", ws, err)
	}

	me, err := api.Me(ctx)
	if err != nil {
		log.Fatalf("%+v, err:%v", me, err)
	}

	command := command.NewCommand(me, api)

	command2.InitCommand(command)

	// 监听哪类事件就需要实现哪类的 handler，定义：websocket/event_handler.go
	var atMessage websocket.MessageEventHandler = func(event *dto.WSPayload, data *dto.WSMessageData) error {
		return command.MessageHandler((*dto.Message)(data))
	}
	intent := websocket.RegisterHandlers(atMessage)
	// 启动 session manager 进行 ws 连接的管理，如果接口返回需要启动多个 shard 的连接，这里也会自动启动多个
	if err := botgo.NewSessionManager().Start(ws, token, &intent); err != nil {
		log.Fatalf("start: %v", err)
	}
}
