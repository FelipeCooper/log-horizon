package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/felipecooper/log-horizon/app/domain/mlogapp"
	protomlog "github.com/felipecooper/log-horizon/app/sdk/proto/mlog"
	"github.com/felipecooper/log-horizon/business/domain/mlog"
	"github.com/felipecooper/log-horizon/business/domain/mlog/mongodb"
	"github.com/felipecooper/log-horizon/foundation/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	logger := newLogger()
	logger.Info(context.Background(), "starting server", "version", "0.1.0")

	mongoURI := getEnv("MONGODB_URI", "mongodb://localhost:27017")
	mongoDBName := getEnv("MONGODB_DBNAME", "loghorizon")
	mongoCollection := getEnv("MONGODB_COLLECTION", "logs")
	exportPath := getEnv("EXPORT_PATH", "./exports")

	grpcPort := getEnv("GRPC_PORT", "50051")

	mongoConfig := mongodb.Config{
		URI:              mongoURI,
		DatabaseName:     mongoDBName,
		CollectionName:   mongoCollection,
		ExportPath:       exportPath,
		CompressionLevel: 9,
	}

	ctx := context.Background()
	store, err := mongodb.NewStore(ctx, logger, mongoConfig)
	if err != nil {
		logger.Error(context.Background(), "failed to create MongoDB store", "error", err)
		os.Exit(1)
	}

	mlogBusiness := mlog.NewMlog(logger, store)
	app := mlogapp.NewApp(logger, mlogBusiness)
	server := grpc.NewServer()
	protomlog.RegisterLogWriterServer(server, app)
	protomlog.RegisterLogReaderServer(server, app)
	reflection.Register(server)
	addr := fmt.Sprintf(":%s", grpcPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error(context.Background(), "failed to listen", "error", err)
		os.Exit(1)
	}

	logger.Info(context.Background(), "server started", "port", grpcPort)

	go func() {
		if err := server.Serve(listener); err != nil {
			logger.Error(context.Background(), "failed to serve", "error", err)
			os.Exit(1)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown

	logger.Info(context.Background(), "shutting down server")
	server.GracefulStop()
	logger.Info(context.Background(), "server stopped")
}

func newLogger() logger.Logger {
	return &simpleLogger{}
}

type simpleLogger struct{}

func (l *simpleLogger) Info(ctx context.Context, msg string, keyValues ...interface{}) {
	log.Printf("INFO: %s %v\n", msg, keyValues)
}

func (l *simpleLogger) Error(ctx context.Context, msg string, keyValues ...interface{}) {
	log.Printf("ERROR: %s %v\n", msg, keyValues)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
