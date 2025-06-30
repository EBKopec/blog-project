package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/blog-project/internal/models"
	"github.com/blog-project/internal/storage"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	DB *storage.DBBlog
}

func NewHandler(db *storage.DBBlog) *Handler {
	return &Handler{
		DB: db,
	}
}

func (h *Handler) ListPosts(c *gin.Context) {
	posts, err := h.DB.GetPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, posts)
}

func (h *Handler) ListPost(c *gin.Context) {
	id := c.Param("id")
	postID, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("Invalid ID: %s\n", id)
		c.JSON(http.StatusBadGateway, gin.H{"error": "Invalid post ID"})
		return
	}
	post, err := h.DB.GetPost(postID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	fmt.Println(post)
	c.JSON(http.StatusOK, post)
}

func (h *Handler) SetPost(c *gin.Context) {
	var postRequest struct {
		Title   string `binding:"required" json:"title"`
		Content string `binding:"required" json:"content"`
	}

	if err := c.BindJSON(&postRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	post := models.Post{
		Title:   postRequest.Title,
		Content: postRequest.Content,
	}

	result, err := h.DB.CreatePost(post)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	log.Printf("Post %v has been created", result)

	c.JSON(http.StatusCreated, result)
}

func (h *Handler) AddComment(c *gin.Context) {
	postIDSTR := c.Param("id")
	postID, err := strconv.Atoi(postIDSTR)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var commentReq struct {
		Content string `binding:"required" json:"content"`
	}

	if err := c.BindJSON(&commentReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	commentID, err := h.DB.CreateComment(postID, commentReq.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusCreated, gin.H{"comment_id": commentID})
}
