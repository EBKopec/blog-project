package router

import (
	"net/http"

	"github.com/blog-project/internal/handlers"
	"github.com/gin-gonic/gin"
)

func Router(db *handlers.Handler) *gin.Engine {
	server := gin.Default()

	server.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusText(http.StatusOK),
		})
	})

	posts := server.Group("/api/posts")
	{
		posts.GET("", db.ListPosts)
		posts.GET("/:id", db.ListPost)
		posts.POST("", db.SetPost)
		posts.POST("/:id/comments", db.AddComment)
	}

	return server
}
