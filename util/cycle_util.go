package cycle_util

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func DegdCall(to time.Duration, c *gin.Context, action func() (int, interface{})) {
	adone := make(chan bool, 1)
	var status int
	var res interface{}

	go func() {
		defer func() {
			if e := recover(); e != nil {
				fmt.Println(e)
				adone <- true
			} else {
				adone <- false
			}
		}()

		status, res = action()
	}()

	select {
	case e := <-adone:
		if e {
			panic(errors.New(enum.ERR_ROUTE_PANIC))
		} else {
			if s, ok := res.(string); ok {
				c.String(status, s)
			} else {
				c.JSON(status, res)
			}
		}
	case <-time.After(to):
		panic(errors.New(enum.ERR_ROUTE_TIMEOUT))
	}
}
