package repo

import (
	pb "github.com/husanmusa/bookshop/genproto/catalog"
)

type CategoryStorageI interface {
	CategoryCreate(pb.Category) (pb.Category, error)
	CategoryGet(id string) (pb.Category, error)
	CategoryList(page, limit int64) ([]*pb.Category, int64, error)
	CategoryUpdate(pb.Category) (pb.Category, error)
	CategoryDelete(id string) error
}
