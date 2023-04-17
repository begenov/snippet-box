package repository

import (
	"database/sql"
	"errors"
	"snippet/model"
)

type Database struct {
	db *sql.DB
}

var ErrNoRecord = errors.New("error")

type IAuthSnippetModel interface {
	Insert(title, content, expires string) (int, error)
	Get(id int) (*model.Snippet, error)
	Latest() ([]*model.Snippet, error)
	InsertUser(name, email, password string) error
	Authenticate(email, password string) (int, error)
}

func NewDB(db *sql.DB) IAuthSnippetModel {
	return &Database{
		db: db,
	}
}
