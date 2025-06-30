package storage_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/blog-project/internal/models"
	"github.com/blog-project/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestCreatePost(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("INSERT INTO blog").
		WithArgs("title1", "content1").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	dbBlog := &storage.DBBlog{DB: db}
	post := models.Post{Title: "title1", Content: "content1"}
	pk, err := dbBlog.CreatePost(post)

	assert.NoError(t, err)
	assert.NotNil(t, pk)
	assert.Equal(t, 1, *pk)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCreateComment(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("INSERT INTO comments").
		WithArgs(1, "comment1").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(42))

	dbBlog := &storage.DBBlog{DB: db}
	pk, err := dbBlog.CreateComment(1, "comment1")

	assert.NoError(t, err)
	assert.NotNil(t, pk)
	assert.Equal(t, 42, *pk)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetPosts(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	postRows := sqlmock.NewRows([]string{"id", "title", "content"}).
		AddRow(1, "post1", "content1").
		AddRow(2, "post2", "content2")

	mock.ExpectQuery("SELECT id, title, content FROM blog").
		WillReturnRows(postRows)

	commentRows1 := sqlmock.NewRows([]string{"id", "post_id", "content"}).
		AddRow(101, 1, "comment A").
		AddRow(102, 1, "comment B")
	mock.ExpectQuery("SELECT id, post_id, content FROM comments WHERE post_id = \\$1").
		WithArgs(1).
		WillReturnRows(commentRows1)

	commentRows2 := sqlmock.NewRows([]string{"id", "post_id", "content"})
	mock.ExpectQuery("SELECT id, post_id, content FROM comments WHERE post_id = \\$1").
		WithArgs(2).
		WillReturnRows(commentRows2)

	dbBlog := &storage.DBBlog{DB: db}
	posts, err := dbBlog.GetPosts()

	assert.NoError(t, err)
	assert.Len(t, *posts, 2)

	assert.Equal(t, 1, (*posts)[0].ID)
	assert.Equal(t, "post1", (*posts)[0].Title)
	assert.Len(t, (*posts)[0].Comments, 2)

	assert.Equal(t, 2, (*posts)[1].ID)
	assert.Equal(t, "post2", (*posts)[1].Title)
	assert.Len(t, (*posts)[1].Comments, 0)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetPost(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	createdAt := time.Now()

	mock.ExpectQuery("SELECT \\* FROM blog WHERE ID = \\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "content", "created_at"}).
			AddRow(1, "post1", "content1", createdAt))

	commentRows := sqlmock.NewRows([]string{"id", "post_id", "content"}).
		AddRow(101, 1, "comment A").
		AddRow(102, 1, "comment B")

	mock.ExpectQuery("SELECT id, post_id, content FROM comments WHERE post_id = \\$1").
		WithArgs(1).
		WillReturnRows(commentRows)

	dbBlog := &storage.DBBlog{DB: db}
	post, err := dbBlog.GetPost(1)

	assert.NoError(t, err)
	assert.Equal(t, 1, post.ID)
	assert.Equal(t, "post1", post.Title)
	assert.Equal(t, "content1", post.Content)
	assert.Equal(t, createdAt, post.CreatedAt)
	assert.Len(t, post.Comments, 2)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetPostNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT \\* FROM blog WHERE ID = \\$1").
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	dbBlog := &storage.DBBlog{DB: db}
	post, err := dbBlog.GetPost(999)

	assert.Nil(t, post)
	assert.Error(t, err)
	assert.Equal(t, "post not found", err.Error())

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
