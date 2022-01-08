package main

import (
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/husanmus/bookshop/config"
	pb "github.com/husanmus/bookshop/genproto/catalogService"
	"github.com/husanmus/bookshop/pkg/db"
	"github.com/husanmus/bookshop/pkg/logger"
	"github.com/husanmus/bookshop/service"
	"github.com/husanmus/bookshop/storage"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.LogLevel, "Catalog")
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
		log.Fatal("sqlx connetion to postgres error", logger.Error(err))
	}

	pgStorage := storage.NewStoragePg(connDB)

	catalogService := service.NewCatalogService(pgStorage, log)

	lis, err := net.Listen("tcp", cfg.RPCPort)
	if err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}

	s := grpc.NewServer()
	pb.RegisterCatalogServiceServer(s, catalogService)
	reflection.Register(s)
	log.Info("main: server running", logger.String("port", cfg.RPCPort))

	if err := s.Serve(lis); err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}
}
