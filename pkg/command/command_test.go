package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tencent-connect/botgo/dto"
)

func TestCommand(t *testing.T) {
	cmd := NewCommand(&dto.User{
		ID: "1000",
	}, nil)
	atme := false
	cmd.AtMe(func(ctx *Context) {
		atme = true
	})
	match := false
	cmd.Match("测试", func(ctx *Context) {
		match = true
	})
	atMatch := false
	cmd.AtMeMatch("测试", func(ctx *Context) {
		atMatch = true
	})
	err := cmd.MessageHandler(&dto.Message{
		Content: "<@!1000>测试",
	})

	assert.Nil(t, err)
	assert.True(t, atme)
	assert.True(t, match)
	assert.True(t, atMatch)

	cmd = NewCommand(&dto.User{
		ID: "1000",
	}, nil)
	match = false
	cmd.Match("全命令", func(ctx *Context) {
		match = true
	})
	err = cmd.MessageHandler(&dto.Message{
		Content: "全命令",
	})
	assert.Nil(t, err)
	assert.True(t, match)
}
