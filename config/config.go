package config

import (
	"log"

	"gopkg.in/ini.v1"
)

type ConfigList struct {
	SQLDriver    string
	DbName       string
	CurrencyCode string
	Duration     string
	Port         int
}

var Config ConfigList

func init() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalln(err)
	}

	Config = ConfigList{
		SQLDriver:    cfg.Section("db").Key("SQLDriver").String(),
		DbName:       cfg.Section("db").Key("DbName").String(),
		CurrencyCode: cfg.Section("fxtrading").Key("currency_code").String(),
		Duration:     cfg.Section("fxtrading").Key("duration").String(),
		Port:         cfg.Section("web").Key("port").MustInt(),
	}

}
