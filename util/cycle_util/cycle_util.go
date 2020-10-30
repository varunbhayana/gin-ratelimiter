package cycle_util

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func DegdCall(to time.Duration, c *gin.Context, action func() (int, interface{})) {
	adone := make(chan *string, 1)
	var status int
	var res interface{}

	go func() {
		defer func() {
			var errorString string

			if e := recover(); e != nil {
				fmt.Println("e in defer", e)
				switch val := e.(type) {
				case error:
					errorString = val.Error()
				case string:
					errorString = val
				default:
					errorString = "error in routing"
				}
				adone <- &errorString
			} else {
				adone <- nil
			}
		}()

		status, res = action()
		fmt.Println("status", status, res)
	}()

	select {
	case e := <-adone:
		fmt.Println("e", status, res)

		if e != nil {
			c.String(500, *e)
			return
			///panic(errors.New("error route panic"))
		} else {
			if s, ok := res.(string); ok {
				c.String(status, s)
			} else {
				c.JSON(status, res)
			}
		}
	case <-time.After(to):
		c.String(500, "error route timeout")
		return

	}
}
