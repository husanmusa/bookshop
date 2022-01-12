package storage

import (
	"github.com/jmoiron/sqlx"

	"github.com/husanmusa/bookshop/storage/postgres"
	"github.com/husanmusa/bookshop/storage/repo"
)

type IStorage interface {
	Category() repo.CategoryStorageI
	Author() repo.AuthorStorageI
	Book() repo.BookStorageI
}

type storagePg struct {
	db           *sqlx.DB
	authorRepo   repo.AuthorStorageI
	categoryRepo repo.CategoryStorageI
	bookRepo     repo.BookStorageI
}

func NewStoragePg(db *sqlx.DB) *storagePg {
	return &storagePg{
		db:           db,
		authorRepo:   postgres.NewAuthorRepo(db),
		categoryRepo: postgres.NewCategoryRepo(db),
		bookRepo:     postgres.NewBookRepo(db),
	}
}

func (s storagePg) Author() repo.AuthorStorageI {
	return s.authorRepo
}

func (s storagePg) Category() repo.CategoryStorageI {
	return s.categoryRepo
}

func (s storagePg) Book() repo.BookStorageI {
	return s.bookRepo
}
