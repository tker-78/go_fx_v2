package config

import (
	"log"
	"time"

	"gopkg.in/ini.v1"
)

type ConfigList struct {
	SQLDriver     string
	DbName        string
	CurrencyCode  string
	Durations     map[string]time.Duration
	TradeDuration time.Duration
	Port          int
	BaseURL       string
}

var Config ConfigList

func init() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalln(err)
	}

	durations := map[string]time.Duration{
		"1m":  time.Minute,
		"5m":  5 * time.Minute,
		"30m": 30 * time.Minute,
		"1h":  time.Hour,
		"4h":  4 * time.Hour,
		"1d":  24 * time.Hour,
	}

	Config = ConfigList{
		SQLDriver:     cfg.Section("db").Key("SQLDriver").String(),
		DbName:        cfg.Section("db").Key("DbName").String(),
		CurrencyCode:  cfg.Section("fxtrading").Key("currency_code").String(),
		Durations:     durations,
		TradeDuration: durations[cfg.Section("fxtrading").Key("duration").String()],
		Port:          cfg.Section("web").Key("port").MustInt(),
		BaseURL:       cfg.Section("web").Key("baseURL").String(),
	}

}
