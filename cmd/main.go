package main

import (
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/husanmusa/bookshop/config"
	pb "github.com/husanmusa/bookshop/genproto/catalog"
	"github.com/husanmusa/bookshop/pkg/db"
	"github.com/husanmusa/bookshop/pkg/logger"
	"github.com/husanmusa/bookshop/service"
	"github.com/husanmusa/bookshop/storage"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.LogLevel, "CatalogService")
	defer func(l logger.Logger) {
		err := logger.Cleanup(l)
		if err != nil {
			log.Fatal("failed cleanup logger", logger.Error(err))
		}
	}(log)

	log.Info("main: sqlxConfig",
		logger.String("host", cfg.PostgresHost),
		logger.Int("port", cfg.PostgresPort),
		logger.String("database", cfg.PostgresDatabase),
	)

	connDB, err := db.ConnectToDB(cfg)
	if err != nil {
		log.Fatal("sqlx connection to postgres error", logger.Error(err))
	}

	pgStorage := storage.NewStoragePg(connDB)

	CatalogService := service.NewCatalogService(pgStorage, log)

	lis, err := net.Listen("tcp", cfg.RPCPort)
	if err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}

	s := grpc.NewServer()
	pb.RegisterCatalogServiceServer(s, CatalogService)

	reflection.Register(s)
	log.Info("main: server running", logger.String("port", cfg.RPCPort))

	if err := s.Serve(lis); err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}
}
