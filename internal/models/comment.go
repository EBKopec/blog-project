package models

type Comment struct {
	ID      int    `json:"id"`
	PostId  int    `json:"postId"`
	Content string `json:"content"`
}
