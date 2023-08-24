package main

import (
	"context"
	"fmt"
	"monapi"
	"os"
	"os/signal"
	"syscall"

	"monapi/internal/config"
	"monapi/internal/handler"
	"monapi/internal/repository"
	"monapi/internal/service"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	cfg, err := config.InitConfig("")
	if err != nil {
		panic(fmt.Sprintf("main(): error initializing config %s", err))
	}

	// dsn := "sqlserver://login:pass@host:port?database=name"
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s",
		cfg.DB.Username,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.DBname)

	db, err := repository.NewMSSQLDB(dsn)

	if err != nil {
		logrus.Fatalf("main(): failed to initialize db: %s", err.Error())
	}

	// create repository
	repos := repository.NewRepository(db)
	//create services
	services := service.NewServices(repos)
	// create hanlers
	handlers := handler.NewHandler(services)

	// create server
	srv := new(monapi.Server)
	go func() {
		if err := srv.Run(cfg.App.Port, handlers.InitRoutes()); err != nil {
			logrus.Fatalf("main(): error occured while running http server: %s", err.Error())
		}

	}()

	logrus.Printf("main(): Service %s started on port = %d ", cfg.App.ServiceName, cfg.App.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Printf("main(): Service %s shutting down", cfg.App.ServiceName)

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("main(): error occured on server shutting down: %s", err.Error())
	}
}
