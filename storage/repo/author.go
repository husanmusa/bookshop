package repo

import pb "github.com/husanmusa/bookshop/genproto/catalog"

type AuthorStorageI interface {
	AuthorCreate(pb.Author) (pb.Author, error)
	AuthorGet(id string) (pb.Author, error)
	AuthorList(page, limit int64) ([]*pb.Author, int64, error)
	AuthorUpdate(pb.Author) (pb.Author, error)
	AuthorDelete(id string) error
}
