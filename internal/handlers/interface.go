package handlers

import "github.com/gin-gonic/gin"

type Handler interface {
	ListPosts(c *gin.Context)
	ListPost(c *gin.Context)
	SetPost(c *gin.Context)
	AddComment(c *gin.Context)
}
