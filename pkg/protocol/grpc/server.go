package grpc

import (
	"context"
	"net"
	"os"
	"os/signal"

	"github.com/Jooho/integration-framework-server/pkg/logger"
	"github.com/Jooho/integration-framework-server/pkg/protocol/grpc/middleware"
	userapi "github.com/Jooho/integration-framework-server/pkg/apis/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// RunServer runs gRPC service to publish ToDo service
func RunServer(ctx context.Context, port string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	// gRPC server statup options
	opts := []grpc.ServerOption{}

	// add middleware
	opts = middleware.AddLogging(logger.Log, opts)

	// register service
	server := grpc.NewServer(opts...)
	
	reflection.Register(server)
	userapi.NewUserServer(*server)

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			logger.Log.Warn("shutting down gRPC server...")

			server.GracefulStop()

			<-ctx.Done()
		}
	}()

	// start gRPC server
	logger.Log.Info("starting gRPC server...Port: "+port)
	return server.Serve(listen)
}