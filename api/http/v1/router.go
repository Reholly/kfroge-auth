package v1

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"sso-service/api/http/v1/handler"
	"sso-service/api/http/v1/middleware"
	"sso-service/internal/service"
	"time"
)

func InitRouter(service service.ServiceManager) *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "*"
		},
		MaxAge: 12 * time.Hour,
	}))

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	authMiddleware := middleware.AuthMiddleware(
		service.TokenService,
		service.PermissionService,
		service.AccountService,
	)

	api := handler.NewApi(service)
	base := router.Group("/api")

	auth := base.Group("/auth")
	{
		auth.POST("/register", api.SignUp)
		auth.POST("/login", api.SignIn)
		auth.POST("/refresh", api.Refresh)
	}

	account := base.Group("/auth/account")
	{
		account.PUT("/confirmreset", api.ConfirmResetPassword)
		account.PUT("/sendresetcode", api.SendResetPasswordCode)
		account.GET("/confirm", api.ConfirmEmail)
	}

	protectedAccount := base.Group("/auth/account/info").
		Use(authMiddleware)
	{
		protectedAccount.GET("/", api.GetProfile)
		protectedAccount.PUT("/", api.ChangeMainInfo)
	}

	admin := base.Group("/auth/admin").
		Use(authMiddleware)
	{
		admin.PUT("/createmoder", api.CreateModerator)
		admin.PUT("/deletemoder", api.DeleteModerator)
		admin.PUT("/ban", api.BanUser)
		admin.PUT("/unban", api.UnbanUser)
	}

	return router
}
