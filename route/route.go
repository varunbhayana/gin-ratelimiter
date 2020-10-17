package route

import (
	"fmt"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/varunbhayana/gin-ratelimiter/db"
	handler "github.com/varunbhayana/gin-ratelimiter/route/handlers"
)

//CORSMiddleware ...
//CORS (Cross-Origin Resource Sharing)
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Origin, Authorization, Accept, Client-Security-Token, Accept-Encoding, x-access-token")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			fmt.Println("OPTIONS")
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}
func Settle(e *gin.Engine) {
	// disableMiddle4xx := func(origin func(*gin.Context)) func(*gin.Context) {
	// 	return func(c *gin.Context) {
	// 		origin(c)
	// 		st := c.Writer.Status()
	// 		if st >= 400 && st < 500 {
	// 			component.GincGetCtx(c).IgnoreError4xx = true
	// 		}
	// 	}
	// }

	// e.Use(
	// 	middleLog,
	// 	newrelic.Middle(),
	// 	validator.MiddleCheckPostReq,
	// 	//	validator.MiddleEncryptResp,
	// 	gzip.Gzip(gzip.DefaultCompression),
	// 	validator.MiddleSignResp,
	// 	middleRecov,
	// 	middlePrepareBody,
	// 	middleRespServerTs,
	// 	middleRespAddCommonHeader,
	// )

	e.Use(CORSMiddleware())
	//e.Use(RequestIDMiddleware())
	e.Use(gzip.Gzip(gzip.DefaultCompression))

	//Start PostgreSQL database
	//Example: db.GetDB() - More info in the models folder
	//db.Init()

	//Start Redis on database 1 - it's used to store the JWT but you can use it for anythig else
	//Example: db.GetRedis().Set(KEY, VALUE, at.Sub(now)).Err()
	db.InitRedis("1")

	//)
	// v1 := r.Group("/v1")
	// {
	// 	/*** START USER ***/
	// 	user := new(controllers.UserController)

	// 	v1.POST("/user/login", user.Login)
	// 	v1.POST("/user/register", user.Register)
	// 	v1.GET("/user/logout", user.Logout)

	// 	/*** START AUTH ***/
	// 	auth := new(controllers.AuthController)

	// 	//Rerfresh the token when needed to generate new access_token and refresh_token for the user
	// 	v1.POST("/token/refresh", auth.Refresh)

	// 	/*** START Article ***/
	// 	article := new(controllers.ArticleController)

	// 	v1.POST("/article", TokenAuthMiddleware(), article.Create)
	// 	v1.GET("/articles", TokenAuthMiddleware(), article.All)
	// 	v1.GET("/article/:id", TokenAuthMiddleware(), article.One)
	// 	v1.PUT("/article/:id", TokenAuthMiddleware(), article.Update)
	// 	v1.DELETE("/article/:id", TokenAuthMiddleware(), article.Delete)
	// }

	//r.LoadHTMLGlob("./public/html/*")

	//r.Static("/public", "./public")
	e.GET("/rate", handler.RateLimit())
	e.GET("/", func(c *gin.Context) {
		c.String(200, "check")
		// c.HTML(http.StatusOK, "index.html", gin.H{
		// 	"ginBoilerplateVersion": "v0.03",
		// 	"goVersion":             runtime.Version(),
		// })
	})
}
