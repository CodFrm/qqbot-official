package command

import (
	"regexp"
)

type HandlerFunc func(ctx *Context)

func Match(command string) HandlerFunc {
	m := regexp.MustCompile("\\[(.+?)\\]")
	var paramName []string
	regexCmd := regexp.MustCompile(m.ReplaceAllStringFunc(command, func(s string) string {
		if len(s) == 0 {
			return ""
		}
		if s[0] == '?' {
			paramName = append(paramName, s[1:])
			return "(.*)"
		}
		paramName = append(paramName, s)
		return "(.+)"
	}))

	return func(ctx *Context) {
		if param := regexCmd.FindStringSubmatch(ctx.Message.Context()); len(param)-1 != len(paramName) {
			ctx.Abort()
			return
		} else {
			for k, v := range param[1:] {
				ctx.setParam(paramName[k], v)
			}
		}
	}
}

func AtMe() HandlerFunc {
	return func(ctx *Context) {
		if !ctx.IsAtMe() {
			ctx.Abort()
		}
	}
}
