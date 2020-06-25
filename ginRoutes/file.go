package ginRoutes

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func AddFileHandler(r *gin.RouterGroup) {
	// r.Static("/static", "../images")
	r.Use(static.Serve("/static/", static.LocalFile("../yourdirectory", false)))
}
