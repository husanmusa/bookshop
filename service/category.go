package service

import (
	"context"

	"github.com/gofrs/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/husanmusa/bookshop/genproto/catalog"
	l "github.com/husanmusa/bookshop/pkg/logger"
)

func (s *CatalogService) CategoryCreate(ctx context.Context, req *pb.Category) (*pb.Category, error) {
	id, err := uuid.NewV4()
	if err != nil {
		s.logger.Error("failed while generating uuid", l.Error(err))
		return nil, status.Error(codes.Internal, "failed generate uuid")
	}
	req.Id = id.String()

	category, err := s.storage.Category().CategoryCreate(*req)
	if err != nil {
		s.logger.Error("failed to create category", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to create category")
	}

	return &category, nil
}

func (s *CatalogService) CategoryGet(ctx context.Context, req *pb.ByIdReq) (*pb.Category, error) {
	category, err := s.storage.Category().CategoryGet(req.GetId())
	if err != nil {
		s.logger.Error("failed to get category", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to get category")
	}

	return &category, nil
}

func (s *CatalogService) CategoryList(ctx context.Context, req *pb.ListReq) (*pb.CategoryListResp, error) {
	categories, count, err := s.storage.Category().CategoryList(req.Page, req.Limit)
	if err != nil {
		s.logger.Error("failed to list Categories", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to list Categories")
	}

	return &pb.CategoryListResp{
		Categories: categories,
		Count:      count,
	}, nil
}

func (s *CatalogService) CategoryUpdate(ctx context.Context, req *pb.Category) (*pb.Category, error) {
	category, err := s.storage.Category().CategoryUpdate(*req)
	if err != nil {
		s.logger.Error("failed to update category", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to update category")
	}

	return &category, nil
}

func (s *CatalogService) CategoryDelete(ctx context.Context, req *pb.ByIdReq) (*pb.EmptyRes, error) {
	err := s.storage.Category().CategoryDelete(req.Id)
	if err != nil {
		s.logger.Error("failed to delete category", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to delete category")
	}

	return &pb.EmptyRes{}, nil
}
