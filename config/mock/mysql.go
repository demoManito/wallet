package mock

import (
	"github.com/sirupsen/logrus"

	"wallet/config"
	"wallet/pkg/closer"
	goloader "wallet/pkg/core/loader"
	"wallet/pkg/mock"
)

func DBOption() mock.Option {
	return func(loader goloader.ILoader, m mock.IMockEnv) {
		conf, _ := config.GetMockConfig()
		db, err := config.InitDatabase(conf.DB)
		if err != nil {
			m.Close()
			m.Fatal(err)
			logrus.Fatalf("init test database error: %v", err)
		}
		m.AppendCloser(closer.NewCloserDelegate(func() error {
			sqldb, _ := db.DB()
			return sqldb.Close()
		}))
		loader.Register("wallet.db", db)
	}
}
