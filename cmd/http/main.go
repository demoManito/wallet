package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"

	"wallet/config"
	"wallet/handler"
	"wallet/models"
	goloader "wallet/pkg/core/loader"
	_ "wallet/service"
)

// GOARCH=amd64 CGO_ENABLED=0 GOOS=linux go build -o bin/http ./cmd/http/main.go

var configFilePath = flag.String("config", "./config.yml", "config file path")

func init() {
	// init ctx
	ctx, cancel := context.WithCancel(context.Background())
	goloader.Register(func(loader goloader.ILoader) error {
		err := loader.Register("wallet.Context", ctx)
		if err != nil {
			return err
		}
		return loader.Register("wallet.Cancel", cancel)
	})
	// init config
	conf, err := config.LoadConfig(*configFilePath)
	if err != nil {
		logrus.Fatalf("error: %v", err)
	}
	goloader.Register(func(loader goloader.ILoader) (err error) {
		return loader.Register("wallet.config", &conf)
	})
	// init db
	goloader.Register(func(loader goloader.ILoader) error {
		db, err := config.InitDatabase(conf.DB)
		db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&models.User{})
		db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&models.Order{})
		db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&models.Wallet{})
		if err != nil {
			logrus.Fatalf("init database error: %v", err)
			return err
		}
		return loader.Register("wallet.db", db)
	})
	// init redis
	goloader.Register(func(loader goloader.ILoader) error {
		rdb := config.InitRedis(conf.Redis)
		return loader.Register("wallet.redis", rdb)
	})
}

func main() {
	srv1 := handler.NewServer()
	srv1.Run()

	quit := make(chan os.Signal)
	signal.Reset(os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	<-quit
	logrus.Info("server stopped...")
}
