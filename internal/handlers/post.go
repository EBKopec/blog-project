package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/blog-project/internal/models"
	"github.com/blog-project/internal/storage"
	"github.com/gin-gonic/gin"
)

type handler struct {
	DB storage.DB
}

func NewHandler(db storage.DB) Handler {
	return &handler{
		DB: db,
	}
}

// ListPosts godoc
// @Summary List blog posts
// @Description Get all posts with comment count
// @Tags posts
// @Produce json
// @Success 200 {array} models.Post
// @Failure 500 {object} handlers.ErrorResponse
// @Router /api/posts [get]
func (h *handler) ListPosts(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		offset = 0
	}

	titleFilter := c.DefaultQuery("title", "")

	posts, err := h.DB.GetPosts(limit, offset, titleFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if posts == nil || len(*posts) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No posts found"})
		return
	}

	c.JSON(http.StatusOK, posts)
}

// ListPost godoc
// @Summary Get a blog post by ID
// @Description Get post details including full comments by ID
// @Tags posts
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} models.Post
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 404 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /api/posts/{id} [get]
func (h *handler) ListPost(c *gin.Context) {
	id := c.Param("id")
	postID, err := strconv.Atoi(id)
	if err != nil {
		respondWithError(c, http.StatusBadGateway, "Invalid post ID", err)
		return
	}
	post, err := h.DB.GetPost(postID)
	if err != nil {
		respondWithError(c, http.StatusNotFound, "Post not found", err)
		return
	}

	c.JSON(http.StatusOK, post)
}

// SetPost godoc
// @Summary Create a new blog post
// @Description Create post with title and content
// @Tags posts
// @Accept json
// @Produce json
// @Param post body models.Post true "Post data"
// @Success 201 {integer} int "New Post ID"
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /api/posts [post]
func (h *handler) SetPost(c *gin.Context) {
	var postRequest struct {
		Title   string `binding:"required" json:"title"`
		Content string `binding:"required" json:"content"`
	}

	if err := c.BindJSON(&postRequest); err != nil {
		respondWithError(c, http.StatusBadRequest, "Failed to create post", err)
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

type AddCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

// AddComment godoc
// @Summary Add a comment to a post
// @Description Create comment linked to a specific post by ID
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Param comment body handlers.AddCommentRequest true "Comment content"
// @Success 201 {integer} int "New Comment ID"
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /api/posts/{id}/comments [post]
func (h *handler) AddComment(c *gin.Context) {
	postIDSTR := c.Param("id")
	postID, err := strconv.Atoi(postIDSTR)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid post ID", err)
		return
	}

	var commentReq AddCommentRequest

	if err := c.BindJSON(&commentReq); err != nil {
		respondWithError(c, http.StatusBadRequest, "Failed to add comment", err)
		return
	}

	commentID, err := h.DB.CreateComment(postID, commentReq.Content)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to add comment", err)
	}

	c.JSON(http.StatusCreated, gin.H{"comment_id": commentID})
}
