package src

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/gin-gonic/gin"
	"github.com/golang-grpc-proxy/config"
	"github.com/golang-grpc-proxy/src/server"
	"github.com/golang-grpc-proxy/src/services"
	"github.com/golang-grpc-proxy/src/storage"
	"google.golang.org/grpc"
)

type App struct {
	logger *slog.Logger
	api services.Api
}

func NewApp(config *config.Config, logger *slog.Logger) App {
	storage, err := storage.New(config, logger)
	
	if err != nil {
		panic("database connection error: " + err.Error())
	}

	api := services.NewApi(&storage, config)

	return App{logger: logger, api: api}
}

func (a *App) Run(grpcPort int, httpPort int) {
	go a.serveGrpc(grpcPort)
	a.serverHttp(httpPort)
}

func (a *App) serveGrpc(grpcPort int) {
	grpcServer := grpc.NewServer()
	server.RegisterGrpcServer(grpcServer, a.api)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
    if err != nil {
        panic("Starting grpc server error: " + err.Error())
    }

	a.logger.Info("grpc server started", slog.String("addr", l.Addr().String()))

	err = grpcServer.Serve(l)
    if err != nil {
        panic("Running grpc server error: " + err.Error())
    }
}

func (a *App) serverHttp(httpPort int) {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	server.RegisterHttpServer(engine, a.api)

	address := fmt.Sprintf("localhost:%d", httpPort)
	a.logger.Info("http server started", slog.String("addr", address))

	engine.Run(address)
}