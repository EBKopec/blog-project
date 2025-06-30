package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/blog-project/internal/models"
	_ "github.com/lib/pq"
)

type DBBlog struct {
	User   string `json:"user"`
	Passwd string `json:"passwd"`
	DBName string `json:"db_name"`
	DBHost string `json:"db_host"`
	DB     *sql.DB
}

func NewDBBlog(data DBBlog) *DBBlog {
	return &DBBlog{
		User:   data.User,
		Passwd: data.Passwd,
		DBName: data.DBName,
		DBHost: data.DBHost,
	}
}

func (db *DBBlog) Open() (*DBBlog, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable",
		db.User, db.Passwd, db.DBHost, db.DBName)

	sqlDB, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err = sqlDB.Ping(); err != nil {
		return nil, err
	}

	db.DB = sqlDB
	return db, nil
}

func (db *DBBlog) CreatePost(post models.Post) (*int, error) {
	query := "INSERT INTO blog (title, content) VALUES ($1, $2) RETURNING ID"
	var pk int
	err := db.DB.QueryRow(query, post.Title, post.Content).Scan(&pk)
	if err != nil {
		return nil, err
	}
	return &pk, nil
}

func (db *DBBlog) CreateComment(postId int, comment string) (*int, error) {
	query := "INSERT INTO comments (post_id, content) VALUES ($1, $2) RETURNING ID"
	var pk int
	err := db.DB.QueryRow(query, postId, comment).Scan(&pk)
	if err != nil {
		return nil, err
	}
	return &pk, nil
}

func (db *DBBlog) GetPosts(limit, offset int, titleFilter string) (*[]models.Post, error) {
	var posts []models.Post

	queryPosts := "SELECT blog.id, blog.title, blog.content, blog.created_at, COUNT(comments.id) AS comments_count " +
		"FROM blog " +
		"LEFT JOIN comments ON blog.id = comments.post_id " +
		"WHERE title ILIKE '%' || $1 || '%' " +
		"GROUP BY blog.id, blog.title, blog.content, blog.created_at " +
		"ORDER BY blog.id " +
		"LIMIT $2 OFFSET $3"
	rows, err := db.DB.Query(queryPosts, titleFilter, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.CommentsCount); err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}
	return &posts, nil
}

func (db *DBBlog) GetPost(postID int) (*models.Post, error) {
	post := &models.Post{}
	query := "SELECT blog.id, blog.title, blog.content, blog.created_at, COUNT(comments.id) AS comments_count " +
		"FROM blog " +
		"LEFT JOIN comments ON blog.id = comments.post_id " +
		"WHERE blog.id = $1 " +
		"GROUP BY blog.id, blog.title, blog.content, blog.created_at " +
		"ORDER BY blog.id"
	err := db.DB.QueryRow(query, postID).Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.CommentsCount)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("post not found")
		}
		return nil, err
	}

	queryComments := "SELECT id, post_id, content FROM comments WHERE post_id = $1"
	rows, err := db.DB.Query(queryComments, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var c models.Comment
		if err := rows.Scan(&c.ID, &c.PostId, &c.Content); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	post.Comments = comments

	return post, nil
}
