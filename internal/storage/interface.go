package storage

import "github.com/blog-project/internal/models"

type DB interface {
	Open() (*DBBlog, error)
	GetPost(postID int) (*models.Post, error)
	GetPosts(limit, offset int, titleFilter string) (*[]models.Post, error)
	CreatePost(post models.Post) (*int, error)
	CreateComment(postId int, comment string) (*int, error)
}
