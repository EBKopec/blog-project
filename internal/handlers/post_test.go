package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blog-project/internal/handlers"
	"github.com/blog-project/internal/models"
	"github.com/blog-project/internal/storage"
	_ "github.com/blog-project/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockDB struct{}

func (m *MockDB) Open() (*storage.DBBlog, error) {
	return nil, nil
}

func (m *MockDB) GetPost(postID int) (*models.Post, error) {
	return &models.Post{ID: postID, Title: "Mock Title", Content: "Mock Content"}, nil
}

func (m *MockDB) GetPosts(limit, offset int, titleFilter string) (*[]models.Post, error) {
	posts := []models.Post{
		{ID: 1, Title: "Test Post", Content: "Test content"},
	}
	return &posts, nil
}

func (m *MockDB) CreatePost(post models.Post) (*int, error) {
	id := 1
	return &id, nil
}

func (m *MockDB) CreateComment(postId int, comment string) (*int, error) {
	id := 1
	return &id, nil
}

func setupRouter() *gin.Engine {
	mockDB := &MockDB{}
	handler := handlers.NewHandler(mockDB)
	router := gin.Default()
	router.GET("/api/posts", handler.ListPosts)
	router.POST("/api/posts", handler.SetPost)
	return router
}

func TestListPosts(t *testing.T) {
	router := setupRouter()
	req, _ := http.NewRequest("GET", "/api/posts", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Test Post")
}

func TestSetPost(t *testing.T) {
	router := setupRouter()
	postData := map[string]string{"title": "New Post", "content": "New Content"}
	jsonData, _ := json.Marshal(postData)
	req, _ := http.NewRequest("POST", "/api/posts", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "1")
}
