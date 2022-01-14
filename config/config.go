package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Http  *Http     `yaml:"http"`
	DB    *Database `yaml:"db"`
	Redis *Redis    `yaml:"redis"`
}

type Http struct {
	Port string `yaml:"port"`
}

type Redis struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
	PoolSize int    `yaml:"pool_size"`
}

type Database struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	Name         string `yaml:"name"`
	User         string `yaml:"user"`
	Pass         string `yaml:"pass"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
}

func LoadConfig(confPath string) (Config, error) {
	var conf Config
	configFile, err := ioutil.ReadFile(confPath)
	if err != nil {
		return conf, err
	}
	err = yaml.Unmarshal(configFile, &conf)
	if err != nil {
		return conf, err
	}
	return conf, nil
}

func GetMockConfig() (*Config, error) {
	conf := new(Config)
	configFile, err := ioutil.ReadFile("../config.yml")
	if err != nil {
		return conf, err
	}
	err = yaml.Unmarshal(configFile, &conf)
	if err != nil {
		return conf, err
	}
	return conf, nil
}
