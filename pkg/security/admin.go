package security

import (
	"log"
	"time"

	"github.com/StindCo/smart_ispt/internal/entities"
	repository "github.com/StindCo/smart_ispt/internal/pkg/identity/repository"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

var identityKey string = "id"

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func InitAdminSecurityMiddleware(repository repository.UserRepository) (*jwt.GinJWTMiddleware, error) {
	auth, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("VGhlIHNtYXJ0IGdyb3VwIGVzdCB1bmUgw6lxdWlwZSBkZSBnw6luaWUgZW4gaW5mb3JtYXRpcXVlCg=="),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,

		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*entities.User); ok {
				return jwt.MapClaims{
					identityKey: v.ID,
				}
			}
			return jwt.MapClaims{}
		},

		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &entities.User{
				ID: claims[identityKey].(string),
			}
		},

		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals Login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			username := loginVals.Username
			password := loginVals.Password

			user, err := repository.GetByUsername(username)
			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			}
			if user.IsAdmin == 0 {
				return nil, jwt.ErrFailedAuthentication
			}

			if (user.ValidatePassword(password)) != nil {
				return nil, jwt.ErrFailedAuthentication
			}

			return user, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			user1, _ := data.(*entities.User)

			_, err := repository.Get(user1.ID)

			return err == nil
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

	errInit := auth.MiddlewareInit()

	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}
	return auth, err
}
