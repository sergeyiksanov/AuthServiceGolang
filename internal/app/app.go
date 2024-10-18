package app

import (
	"AuthService/internal/config"
	proto "AuthService/pkg/api/v1"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
)

type App struct {
	grpcServer      *grpc.Server
	ServiceProvider *serviceProvider
}

const (
	portName = "GRPC_PORT"
)

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initGRPC,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) Run() error {
	return a.runGRPC()
}

func (a *App) initConfig(_ context.Context) error {
	err := config.Load(".env")
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.ServiceProvider = newServiceProvider()
	return nil
}

func (a *App) initGRPC(_ context.Context) error {
	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))

	reflection.Register(a.grpcServer)

	authServerImpl := a.ServiceProvider.AuthServerImpl()

	proto.RegisterAuthServer(a.grpcServer, authServerImpl)

	return nil
}

func (a *App) runGRPC() error {
	log.Printf("Starting gRPC Server")
	port := os.Getenv(portName)

	list, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	err = a.grpcServer.Serve(list)
	if err != nil {
		return err
	}

	return nil
}
