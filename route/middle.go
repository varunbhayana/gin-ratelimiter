package route

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func middleRecov(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			c.String(500, "")
		}
	}()
	c.Next()
}
