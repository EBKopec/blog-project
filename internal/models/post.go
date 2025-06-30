package models

import "time"

type Post struct {
	ID            int       `json:"id"`
	Title         string    `json:"title"`
	Content       string    `json:"content,omitempty"`
	CommentsCount int       `json:"comments_count"`
	Comments      []Comment `json:"comments,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}
