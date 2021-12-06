package http

import (
	"github.com/gin-gonic/gin"
)

func images(c *gin.Context) {
	file := c.Query("name")
	c.File("./assets/images/" + file)
}
