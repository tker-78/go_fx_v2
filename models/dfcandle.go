package models

import "time"

type DataFrameCandle struct {
	CurrencyCode string        `json:"currency_code"`
	Duration     time.Duration `json:"duration"`
	Candles      []Candle      `json:"candles"`
}
