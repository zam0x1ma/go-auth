package routes

import (
    "go-auth/controllers"

    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
    router := gin.Default()

    config := cors.DefaultConfig()
    config.AllowCredentials = true
    config.AllowAllOrigins = true

    router.Use(cors.New(config))

    router.POST("/api/register", controllers.Register)
    router.POST("/api/login", controllers.Login)
    router.GET("/api/user", controllers.User)
    router.POST("/api/logout", controllers.Logout)

    return router
}
