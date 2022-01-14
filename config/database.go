package config

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	defaultParseTime = true
	defaultCharset   = "utf8mb4,utf8"
	defaultLoc       = "Local"
)

// InitDatabase init database connection
func InitDatabase(conf *Database) (*gorm.DB, error) {
	dsn, err := getDsn(conf)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(mysql.New(mysql.Config{DSN: dsn}), &gorm.Config{})
	if err != nil {
		return db, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return db, err
	}
	if err = sqlDB.Ping(); err != nil {
		return db, err
	}
	sqlDB.SetMaxIdleConns(conf.MaxIdleConns)
	// for db invalid connection after EOF
	sqlDB.SetConnMaxLifetime(time.Second)

	// connect success
	return db, nil
}

// getDsn 获取 DSN
func getDsn(conf *Database) (dsn string, err error) {
	if conf.Host == "" || conf.Port == 0 || conf.Name == "" || conf.User == "" || conf.Pass == "" {
		err = errors.New("db config should not be empty")
		return
	}

	dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s", conf.User, conf.Pass, conf.Host,
		conf.Port, conf.Name, defaultCharset, defaultParseTime, defaultLoc)
	return
}
