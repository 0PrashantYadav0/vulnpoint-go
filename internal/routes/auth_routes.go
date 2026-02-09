package routes

import (
	"github.com/datmedevil17/go-vuln/internal/handlers"
	"github.com/gin-gonic/gin"
)

// RegisterAuthRoutes registers public authentication routes
func RegisterAuthRoutes(rg *gin.RouterGroup, authHandler *handlers.AuthHandler) {
	auth := rg.Group("/auth")
	{
		auth.GET("/github", authHandler.GetAuthURL)
		auth.GET("/github/callback", authHandler.HandleCallback)
		auth.POST("/logout", authHandler.Logout)
	}
}
