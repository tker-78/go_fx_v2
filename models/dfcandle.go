package models

type DataFrameCandle struct {
	CurrencyCode string   `json:"currency_code"`
	Duration     string   `json:"duration"`
	Candles      []Candle `json:"candles"`
}
