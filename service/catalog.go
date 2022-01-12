package service

import (
	l "github.com/husanmusa/bookshop/pkg/logger"
	"github.com/husanmusa/bookshop/storage"
)

type CatalogService struct {
	storage storage.IStorage
	logger  l.Logger
}

func NewCatalogService(storage storage.IStorage, log l.Logger) *CatalogService {
	return &CatalogService{
		storage: storage,
		logger:  log,
	}
}
