package service

import (
	"context"

	"github.com/gofrs/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/husanmus/bookshop/genproto"
	l "github.com/husanmus/bookshop/pkg/logger"
	"github.com/husanmus/bookshop/storage"
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

func (s *CatalogService) AuthorCreate(ctx context.Context, req *pb.Catalog) (*pb.Catalog, error) {
	id, err := uuid.NewV4()

	if err != nil {
		s.logger.Error("failed while generating uuid", l.Error(err))
		return nil, status.Error(codes.Internal, "failed generate uuid")
	}
	req.Id = id.String()

	task, err := s.storage.Catalog().AuthorCreate(*req)
	if err != nil {
		s.logger.Error("failed to create task", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to create task")
	}

	return &task, nil
}

func (s *CatalogService) AuthorGet(ctx context.Context, req *pb.ByIdReq) (*pb.Catalog, error) {
	task, err := s.storage.Catalog().AuthorGet(req.GetId())
	if err != nil {
		s.logger.Error("failed to get task", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to get task")
	}

	return &task, nil
}

func (s *CatalogService) AuthorList(ctx context.Context, req *pb.ListReq) (*pb.ListResp, error) {
	tasks, count, err := s.storage.Catalog().AuthorList(req.Page, req.Limit)
	if err != nil {
		s.logger.Error("failed to list tasks", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to list tasks")
	}

	return &pb.ListResp{
		Catalogs: tasks,
		Count:    count,
	}, nil
}

func (s *CatalogService) AuthorUpdate(ctx context.Context, req *pb.Catalog) (*pb.Catalog, error) {
	task, err := s.storage.Catalog().AuthorUpdate(*req)
	if err != nil {
		s.logger.Error("failed to update task", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to update task")
	}

	return &task, nil
}

func (s *CatalogService) AuthorDelete(ctx context.Context, req *pb.ByIdReq) (*pb.EmptyResp, error) {
	err := s.storage.Catalog().AuthorDelete(req.Id)
	if err != nil {
		s.logger.Error("failed to delete task", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to delete task")
	}

	return &pb.EmptyResp{}, nil
}
