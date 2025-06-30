package models

import "time"

type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Comments  []Comment `json:"comments"`
	CreatedAt time.Time `json:"created_at"`
}
