package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"

	"wallet/config"
	goloader "wallet/pkg/core/loader"
	"wallet/pkg/work"
)

var configFilePath = flag.String("config", "./config.yml", "config file path")

func init() {
	// init config
	conf, err := config.LoadConfig(*configFilePath)
	if err != nil {
		logrus.Fatalf("error: %v", err)
	}
	goloader.Register(func(loader goloader.ILoader) (err error) {
		return loader.Register("wallet.config", &conf)
	})
	// init ctx
	goloader.Register(func(loader goloader.ILoader) error {
		ctx, cancel := context.WithCancel(context.Background())
		err := loader.Register("wallet.Context", ctx)
		if err != nil {
			return err
		}
		return loader.Register("wallet.Cancel", cancel)
	})
	// init db
	goloader.Register(func(loader goloader.ILoader) error {
		db, err := config.InitDatabase(conf.DB)
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
	h := &handler{}
	goloader.LoadingAll(goloader.DefaultLoader())
	goloader.DefaultLoader().Loading(&h)

	h.WalletMangerWorker.WalletMangerSubscribe()

	quit := make(chan os.Signal)
	signal.Reset(os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	<-quit
	logrus.Info("work stopped...")
}

type handler struct {
	WalletMangerWorker *work.WalletMangerWorker `load:"wallet.worker.WalletMangerWorker"`
}
