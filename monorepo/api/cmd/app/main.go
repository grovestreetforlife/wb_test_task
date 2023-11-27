package main

import (
	"context"
	"flag"
	logger "github.com/dany-ykl/logger"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	_ "wb_test_task/api/docs"
	"wb_test_task/api/internal/application"
)

func main() {
	if err := logger.InitLogger(logger.Config{
		Namespace:   "wb_test_task",
		Development: false,
		Level:       logger.InfoLevel,
	}); err != nil {
		log.Fatalln(err)
	}

	configFile := flag.String("config", "configs/config.yml", "Path to config file")
	shutdownTimeout := flag.Int("shutdown-timeout", 1, "Time to shutdown second")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app, err := application.New(ctx, *configFile)
	if err != nil {
		logger.Fatal("fail to init application", zap.Error(err))
	}

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig

		logger.Info("--- shutdown application ---")
		cancel()

		time.Sleep(time.Duration(*shutdownTimeout))

		if err := app.Shutdown(context.Background()); err != nil {
			logger.Fatal("graceful shutdown is not complete: ", zap.String("error", err.Error()))
		}
	}()

	if err := app.Start(ctx); err != nil {
		log.Fatalln(err)
	}
}
