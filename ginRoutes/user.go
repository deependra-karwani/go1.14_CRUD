package ginRoutes

import (
	"CRUD/controllers"
	"CRUD/workers"

	"github.com/gin-gonic/gin"
)

func AddUserHandler(r *gin.RouterGroup) {
	r.POST("/register", func(c *gin.Context) {
		done := make(chan bool)
		defer close(done)
		workers.AddToQueue(controllers.Register, c.Writer, c.Request, done)
		<-done
	})

	r.PUT("/login", func(c *gin.Context) {
		done := make(chan bool)
		defer close(done)
		workers.AddToQueue(controllers.Login, c.Writer, c.Request, done)
		<-done
	})

	r.PUT("/forgot", func(c *gin.Context) {
		done := make(chan bool)
		defer close(done)
		workers.AddToQueue(controllers.ForgotPassword, c.Writer, c.Request, done)
		<-done
	})

	r.GET("/logout", func(c *gin.Context) {
		done := make(chan bool)
		defer close(done)
		workers.AddToQueue(controllers.Logout, c.Writer, c.Request, done)
		<-done
	})

	r.GET("/getAll", func(c *gin.Context) {
		done := make(chan bool)
		defer close(done)
		workers.AddToQueue(controllers.GetAllUsers, c.Writer, c.Request, done)
		<-done
	})

	r.GET("/getDetails", func(c *gin.Context) {
		done := make(chan bool)
		defer close(done)
		workers.AddToQueue(controllers.GetUserDetails, c.Writer, c.Request, done)
		<-done
	})

	r.PUT("/updProf", func(c *gin.Context) {
		done := make(chan bool)
		defer close(done)
		workers.AddToQueue(controllers.UpdateUserProfile, c.Writer, c.Request, done)
		<-done
	})

	r.DELETE("/delAcc", func(c *gin.Context) {
		done := make(chan bool)
		defer close(done)
		workers.AddToQueue(controllers.DeleteUserAccount, c.Writer, c.Request, done)
		<-done
	})
}
