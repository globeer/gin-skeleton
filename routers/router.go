// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"bower.co.kr/c4bapi/controllers"
	"bower.co.kr/c4bapi/models"

	beego "github.com/beego/beego/v2/server/web"

	"bower.co.kr/c4bapi/middleware"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func init() {
	ns := beego.NewNamespace("/v1",

		beego.NSNamespace("/board",
			beego.NSInclude(
				&controllers.BoardController{},
			),
		),

		beego.NSNamespace("/board_comment",
			beego.NSInclude(
				&controllers.BoardCommentController{},
			),
		),

		beego.NSNamespace("/board_file",
			beego.NSInclude(
				&controllers.BoardFileController{},
			),
		),

		beego.NSNamespace("/board_reader",
			beego.NSInclude(
				&controllers.BoardReaderController{},
			),
		),

		beego.NSNamespace("/category",
			beego.NSInclude(
				&controllers.CategoryController{},
			),
		),

		beego.NSNamespace("/code",
			beego.NSInclude(
				&controllers.CodeController{},
			),
		),

		beego.NSNamespace("/code_detail",
			beego.NSInclude(
				&controllers.CodeDetailController{},
			),
		),

		beego.NSNamespace("/customer",
			beego.NSInclude(
				&controllers.CustomerController{},
			),
		),

		beego.NSNamespace("/deposit",
			beego.NSInclude(
				&controllers.DepositController{},
			),
		),

		beego.NSNamespace("/goods",
			beego.NSInclude(
				&controllers.GoodsController{},
			),
		),

		beego.NSNamespace("/group",
			beego.NSInclude(
				&controllers.GroupController{},
			),
		),

		beego.NSNamespace("/group_user",
			beego.NSInclude(
				&controllers.GroupUserController{},
			),
		),

		beego.NSNamespace("/letter",
			beego.NSInclude(
				&controllers.LetterController{},
			),
		),

		beego.NSNamespace("/letter_aligo",
			beego.NSInclude(
				&controllers.LetterAligoController{},
			),
		),

		beego.NSNamespace("/login_trial",
			beego.NSInclude(
				&controllers.LoginTrialController{},
			),
		),

		beego.NSNamespace("/order",
			beego.NSInclude(
				&controllers.OrderController{},
			),
		),

		beego.NSNamespace("/order_detail",
			beego.NSInclude(
				&controllers.OrderDetailController{},
			),
		),

		beego.NSNamespace("/reserve",
			beego.NSInclude(
				&controllers.ReserveController{},
			),
		),

		beego.NSNamespace("/set_letter",
			beego.NSInclude(
				&controllers.SetLetterController{},
			),
		),

		beego.NSNamespace("/shop",
			beego.NSInclude(
				&controllers.ShopController{},
			),
		),

		beego.NSNamespace("/staff",
			beego.NSInclude(
				&controllers.StaffController{},
			),
		),

		beego.NSNamespace("/stat_basic",
			beego.NSInclude(
				&controllers.StatBasicController{},
			),
		),

		beego.NSNamespace("/stat_hist",
			beego.NSInclude(
				&controllers.StatHistController{},
			),
		),

		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
	)
	beego.AddNamespace(ns)
}

func Route(router *gin.Engine) {
	indexController := new(controllers.IndexController)
	router.GET(
		"/", indexController.GetIndex,
	)

	auth := router.Group("/auth")
	authMiddleware := middleware.Auth()
	auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/hello", func(c *gin.Context) {
			claims := jwt.ExtractClaims(c)
			user, _ := c.Get("userId")
			c.JSON(200, gin.H{
				"userId": claims["userId"],
				"userNm": user.(*models.User1).UserNm,
				"text":   "Hello World.",
			})
		})
	}

	userController := new(controllers.User1Controller)
	router.GET(
		"/user/:id", userController.GetUser,
	).GET(
		"/signup", userController.SignupForm,
	).POST(
		"/signup", userController.Signup,
	).GET(
		"/login", userController.LoginForm,
	).POST(
		"/login", authMiddleware.LoginHandler,
	)

	api := router.Group("/api")
	{
		api.GET("/version", indexController.GetVersion)
	}
}
