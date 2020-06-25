package main

import (
	"CRUD/ginRoutes"
	"CRUD/middlewares"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	// r := gin.Default()

	userSr := r.Group("/user")
	fileSr := r.Group("/")

	userSr.Use(middlewares.UserAuthGin())

	ginRoutes.AddUserHandler(userSr)
	ginRoutes.AddFileHandler(fileSr)

	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{http.MethodGet, http.MethodDelete, http.MethodPost, http.MethodPut},
		AllowHeaders:  []string{"X-Requested-With", "Content-Type", "Accept", "Origin", "token"},
		ExposeHeaders: []string{"token"},
		// AllowCredentials: true,
		// AllowOriginFunc: func(origin string) bool {
		// 	return origin == "https://github.com"
		// },
		// MaxAge: 12 * time.Hour,
	}))
	// r.Use(cors.Default())

	r.Run(":8080")
}
