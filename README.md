## QQ频道机器人
> 有很多自用的功能,也有一个简单的机器人命令框架,自用功能可以直接编译然后配置好配置文件就可使用




### 命令框架
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



### 机器人

#### 机器人命令

> 以下是本仓库的机器人命令,语料配置可以查看:[语料配置模版.csv](./docs/语料配置模版.csv)

| 命令           | 返回                               | 备注                             |
| -------------- | ---------------------------------- | -------------------------------- |
| 警告           | 管理员警告用户                     | 需要艾特并且是管理员权限才会触发 |
| 我的信息       | 返回用户信息、警告等级、用户组选择 | 警告成员无法获得用户信息         |
| 摇色子\|摇骰子 | 回复1-6                            |                                  |


#### 配置和启动

将目录下的`config.yaml.example`复制粘贴一份`config.yaml`并修改其中的配置,然后直接启动编译好的二进制程序即可


### 机器人申请

机器人的申请方式请看:[QQ机器人](https://bot.q.qq.com/wiki/)



#### 常见问题

> 很多申请的东西和新更新的咨询需要加入`QQ机器人官方频道`获取

* [完善信息时测试频道ID如何获取](https://github.com/CodFrm/qqbot-official/issues/1)
* 主动推送消息限制:公域每日2条,私域每日100条(可能变化)

