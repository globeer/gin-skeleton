package middleware

import (
	"log"
	"time"

	"bower.co.kr/c4bapi/models"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

var authMiddleware *jwt.GinJWTMiddleware
var identityKey = "email"

// Login struct
type Login struct {
	UserId string `form:"userId" json:"user_id" binding:"required"`
	UserPw string `form:"userPw" json:"user_pw" binding:"required,min=6,max=20"`
}

// Auth middleware
func Auth() *jwt.GinJWTMiddleware {
	return authMiddleware
}

func init() {
	var err error
	authMiddleware, err = jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "gin-skeleton",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		SendCookie:  true,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.User1); ok {
				return jwt.MapClaims{
					identityKey: v.UserId,
					"userNm":    v.UserNm,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &models.User1{
				UserId: claims[identityKey].(string),
				UserNm: claims["userNm"].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals Login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userId := loginVals.UserId
			userPw := loginVals.UserPw

			return models.LoginByEmailAndPassword(userId, userPw)
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(*models.User1); ok && v.UserId == "admin" {
				return true
			}
			// log.Default(c.Params.ByName())
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}
}
