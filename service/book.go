package service

import (
	"context"

	"github.com/gofrs/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/husanmusa/bookshop/genproto/catalog"
	l "github.com/husanmusa/bookshop/pkg/logger"
)

func (s *CatalogService) BookCreate(ctx context.Context, req *pb.Book) (*pb.Book, error) {
	id, err := uuid.NewV4()
	if err != nil {
		s.logger.Error("failed while generating uuid", l.Error(err))
		return nil, status.Error(codes.Internal, "failed generate uuid")
	}
	req.Id = id.String()

	book, err := s.storage.Book().BookCreate(*req)
	if err != nil {
		s.logger.Error("failed to create book", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to create book")
	}

	return &book, nil
}

func (s *CatalogService) BookGet(ctx context.Context, req *pb.ByIdReq) (*pb.Book, error) {
	book, err := s.storage.Book().BookGet(req.Id)
	if err != nil {
		s.logger.Error("failed to get book", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to get book")
	}

	return &book, nil
}

func (s *CatalogService) BookList(ctx context.Context, req *pb.ListBookReq) (*pb.BookListResp, error) {
	books, count, err := s.storage.Book().BookList(req.Page, req.Limit, req.Filters)
	if err != nil {
		s.logger.Error("failed to list books", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to list books")
	}

	return &pb.BookListResp{
		Books: books,
		Count: count,
	}, nil
}

func (s *CatalogService) BookUpdate(ctx context.Context, req *pb.Book) (*pb.Book, error) {
	book, err := s.storage.Book().BookUpdate(*req)
	if err != nil {
		s.logger.Error("failed to update book", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to update book")
	}

	return &book, nil
}

func (s *CatalogService) BookDelete(ctx context.Context, req *pb.ByIdReq) (*pb.EmptyRes, error) {
	err := s.storage.Book().BookDelete(req.Id)
	if err != nil {
		s.logger.Error("failed to delete book", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to delete book")
	}

	return &pb.EmptyRes{}, nil
}
