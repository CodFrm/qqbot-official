package http

import (
	"context"
	"net/http"

	"github.com/CodFrm/qqbot-official/internal/config"
	"github.com/CodFrm/qqbot-official/internal/db"
	api2 "github.com/CodFrm/qqbot-official/internal/utils/api"
	"github.com/CodFrm/qqbot-official/internal/utils/errs"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tencent-connect/botgo/openapi"
)

var actionMap = map[string]func(c *gin.Context){}

func msgUrl(openapi openapi.OpenAPI) func(c *gin.Context) {
	return func(c *gin.Context) {
		action := c.Query("action")
		if action == "images" {
			images(c)
			return
		}
		session := c.Query("session")
		user := c.Query("user")
		if v, err := db.Get("session:" + session); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "该链接已经点过啦"})
			return
		} else if v == nil || string(v) != user {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "错误链接"})
			return
		}
		_ = db.Del("session:" + session)
		f, ok := actionMap[action]
		if ok {
			f(c)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "错误方法"})
		}
	}
}

type Service interface {
	Registry(ctx context.Context, r *gin.Engine)
}

func Registry(ctx context.Context, r *gin.Engine, registry ...Service) {
	for _, v := range registry {
		v.Registry(ctx, r)
	}
}

func handle(ctx *gin.Context, f func() interface{}) {
	resp := f()
	if resp == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0, "msg": "ok",
		})
		return
	}
	handleResp(ctx, resp)
}

func handleResp(ctx *gin.Context, resp interface{}) {
	switch resp.(type) {
	case *errs.JsonRespondError:
		err := resp.(*errs.JsonRespondError)
		ctx.JSON(err.Status, err)
	case error:
		err := resp.(error)
		logrus.Errorf("%s - %s: %v", ctx.Request.RequestURI, ctx.ClientIP(), err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": -1, "msg": "系统错误",
		})
	case string:

	default:
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0, "msg": "ok", "data": resp,
		})
	}
}

func StartWeb(openapi openapi.OpenAPI) error {
	api := api2.NewGuildApi(openapi)
	r := gin.Default()

	r.GET("/redirect", msgUrl(openapi))

	Registry(context.Background(), r,
		NewUser(api),
	)

	return r.Run(config.AppConfig.WebPort)
}
