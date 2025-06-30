package models_test

import (
	"testing"

	"github.com/blog-project/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestPostModel(t *testing.T) {
	post := models.Post{
		ID:      1,
		Title:   "Sample Title",
		Content: "Sample Content",
	}

	assert.Equal(t, 1, post.ID)
	assert.Equal(t, "Sample Title", post.Title)
	assert.Equal(t, "Sample Content", post.Content)
}

func TestCommentModel(t *testing.T) {
	comment := models.Comment{
		ID:      1,
		PostId:  1,
		Content: "Sample comment content",
	}

	assert.Equal(t, 1, comment.ID)
	assert.Equal(t, 1, comment.PostId)
	assert.Equal(t, "Sample comment content", comment.Content)
}
