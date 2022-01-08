package repo

import (
	pb "github.com/husanmus/bookshop/genproto/catalogService"
)

type CatalogStorageI interface {
	BookCreate(pb.Book) (pb.Book, error)
	BookGet(id string) (pb.Book, error)
	BookList(page, limit int64) ([]*pb.Book, int64, error)
	BookUpdate(pb.Book) (pb.Book, error)
	BookDelete(id string) error

	AuthorCreate(pb.Author) (pb.Author, error)
	AuthorGet(id string) (pb.Author, error)
	AuthorList(page, limit int64) ([]*pb.Author, int64, error)
	AuthorUpdate(pb.Author) (pb.Author, error)
	AuthorDelete(id string) error

	CategoryCreate(pb.Category) (pb.Category, error)
	CategoryGet(id string) (pb.Category, error)
	CategoryList(page, limit int64) ([]*pb.Category, int64, error)
	CategoryUpdate(pb.Category) (pb.Category, error)
	CategoryDelete(id string) error
}
