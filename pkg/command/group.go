package command

type Group struct {
	group   []*Group
	handler []HandlerFunc
}

func newGroup() *Group {
	return &Group{
		handler: make([]HandlerFunc, 0),
		group:   make([]*Group, 0),
	}
}

func (g *Group) handle(ctx *Context) {
	for _, v := range g.handler {
		v(ctx)
		if ctx.IsAborted() {
			break
		}
	}
	if ctx.IsAborted() {
		return
	}
	for _, v := range g.group {
		v.handle(ctx)
		if ctx.IsAborted() {
			break
		}
	}
}

func (g *Group) Use(handler ...HandlerFunc) {
	g.handler = append(g.handler, handler...)
}

func (g *Group) Group(handler ...HandlerFunc) *Group {
	ret := &Group{}
	ret.Use(handler...)
	g.group = append(g.group, ret)
	return ret
}

// Match 命令 [参数1] [参数2]
func (g *Group) Match(command string, handler ...HandlerFunc) {
	g.handler = append(append(g.handler, Match(command)), handler...)
}

// AtMeMatch Match 命令 [参数1] [参数2]
func (g *Group) AtMeMatch(command string, handler ...HandlerFunc) {
	g.handler = append(append(g.handler, AtMe()), handler...)
}

func (g *Group) AtMe(handler ...HandlerFunc) {
	g.handler = append(append(g.handler, AtMe()), handler...)
}
