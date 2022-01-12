package service

import (
	"context"

	"github.com/gofrs/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/husanmusa/bookshop/genproto/catalog"
	l "github.com/husanmusa/bookshop/pkg/logger"
)

func (s *CatalogService) AuthorCreate(ctx context.Context, req *pb.Author) (*pb.Author, error) {
	id, err := uuid.NewV4()
	if err != nil {
		s.logger.Error("failed while generating uuid", l.Error(err))
		return nil, status.Error(codes.Internal, "failed generate uuid")
	}
	req.Id = id.String()

	author, err := s.storage.Author().AuthorCreate(*req)
	if err != nil {
		s.logger.Error("failed to create author", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to create author")
	}

	return &author, nil
}

func (s *CatalogService) AuthorGet(ctx context.Context, req *pb.ByIdReq) (*pb.Author, error) {
	author, err := s.storage.Author().AuthorGet(req.GetId())
	if err != nil {
		s.logger.Error("failed to get author", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to get author")
	}

	return &author, nil
}

func (s *CatalogService) AuthorList(ctx context.Context, req *pb.ListReq) (*pb.AuthorListResp, error) {
	authors, count, err := s.storage.Author().AuthorList(req.Page, req.Limit)
	if err != nil {
		s.logger.Error("failed to list authors", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to list authors")
	}

	return &pb.AuthorListResp{
		Authors: authors,
		Count:   count,
	}, nil
}

func (s *CatalogService) AuthorUpdate(ctx context.Context, req *pb.Author) (*pb.Author, error) {
	author, err := s.storage.Author().AuthorUpdate(*req)
	if err != nil {
		s.logger.Error("failed to update author", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to update author")
	}

	return &author, nil
}

func (s *CatalogService) AuthorDelete(ctx context.Context, req *pb.ByIdReq) (*pb.EmptyRes, error) {
	err := s.storage.Author().AuthorDelete(req.Id)
	if err != nil {
		s.logger.Error("failed to delete author", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to delete author")
	}

	return &pb.EmptyRes{}, nil
}
