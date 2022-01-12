package repo

import pb "github.com/husanmusa/bookshop/genproto/catalog"

type BookStorageI interface {
	BookCreate(pb.Book) (pb.Book, error)
	BookGet(id string) (pb.Book, error)
	BookList(page, limit int64, filters map[string]string) ([]*pb.Book, int64, error)
	BookUpdate(pb.Book) (pb.Book, error)
	BookDelete(id string) error
}
