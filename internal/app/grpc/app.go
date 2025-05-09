package grpcapp

import (
	"fmt"
	appsgrpc "github.com/cuongpiger/sso/internal/grpc/apps"
	authgrpc "github.com/cuongpiger/sso/internal/grpc/auth"
	permgrpc "github.com/cuongpiger/sso/internal/grpc/permissions"
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	host       string
	port       string
}

func New(log *slog.Logger, authService authgrpc.Auth, permService permgrpc.Perm, appsService appsgrpc.Apps, host string, port string) *App {
	gRPCServer := grpc.NewServer()
	authgrpc.Register(gRPCServer, authService)
	permgrpc.Register(gRPCServer, permService)
	appsgrpc.Register(gRPCServer, appsService)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		host:       host,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.run(); err != nil {
		panic(err)
	}
}

func (a *App) run() error {
	const op = "app.grpc.app.Run"
	log := a.log.With(slog.String("op", op), slog.String("port", a.port))

	l, err := net.Listen("tcp", a.host+":"+a.port)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("gRPC server is running", slog.String("address", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (a *App) Stop() {
	const op = "app.grpc.app.Stop"
	log := a.log.With(slog.String("op", op))

	log.Info("stopping gRPC server")
	a.gRPCServer.GracefulStop()
}
