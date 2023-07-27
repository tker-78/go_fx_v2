package models

// データベースから指定した条件のcandleを格納するための型
type DataFrameCandle struct {
	CurrencyCode string   `json:"currency_code"`
	Duration     string   `json:"duration"`
	Candles      []Candle `json:"candles"`
}
