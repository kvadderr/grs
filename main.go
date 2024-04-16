package main

import (
	"log/slog"
	"os"

	"github.com/golang-grpc-proxy/config"
	"github.com/golang-grpc-proxy/src"
)

//	@title		Http proxy server
//	@version	1.0

func main() {
	cfg := config.Load()

	logger := setupLogger(cfg.Env)

	app := src.NewApp(cfg, logger)
	app.Run(cfg.GRPC.Port, cfg.HTTP.Port)

	// loggingOpts := []logging.Option{
    //     logging.WithLogOnEvents(
    //         logging.PayloadReceived, logging.PayloadSent,
    //     ),
    // }

    // recoveryOpts := []recovery.Option{
    //     recovery.WithRecoveryHandler(func(p interface{}) (err error) {
    //         logger.Error("Recovered from panic", slog.Any("panic", p))

    //         return status.Errorf(codes.Internal, "internal error")
    //     }),
    // }
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case "local":
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "prod":
		logger = slog.New(
            slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
        )
	}

	return logger
}