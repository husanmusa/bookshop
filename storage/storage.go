package storage

import (
	"github.com/jmoiron/sqlx"

	"github.com/husanmus/bookshop/storage/postgres"
	"github.com/husanmus/bookshop/storage/repo"
)

type IStorage interface {
	Catalog() repo.CatalogStorageI
}

type storagePg struct {
	db          *sqlx.DB
	catalogRepo repo.CatalogStorageI
}

func NewStoragePg(db *sqlx.DB) *storagePg {
	return &storagePg{
		db:          db,
		catalogRepo: postgres.NewCatalogRepo(db),
	}
}

func (s storagePg) Catalog() repo.CatalogStorageI {
	return s.catalogRepo
}
