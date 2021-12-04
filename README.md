## QQ频道官方机器人
> 有很多自用的功能,也有一个简单的机器人命令框架,自用功能可以直接编译然后配置好配置文件就可使用

### command
参考gin的设计,实现命令框架

```go
// init
cmd := command.NewCommand(me, api)
// 只会响应管理员艾特me的消息
cg := cmd.Group(command.AtMe(), middleware.Admin)
cg.Match("警告", p.punish)
cg.Match("移除警告", p.remove)

// 匹配艾特命令+摇色子
cmd.AtMeMatch("摇色子|摇骰子", func(ctx *command.Context) {
    ctx.ReplyText(fmt.Sprintf("%v", rand.Int31n(6)+1))
})

```
