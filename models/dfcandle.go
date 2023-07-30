package models

import (
	"time"

	"github.com/markcheno/go-talib"
)

// データベースから指定した条件のcandleを格納するための型
type DataFrameCandle struct {
	CurrencyCode string   `json:"currency_code"`
	Duration     string   `json:"duration"`
	Candles      []Candle `json:"candles"`
	Smas         []Sma    `json:"smas,omitempty"`
}

// technical用の型
type Sma struct {
	Period int       `json:"period,omitempty"`
	Value  []float64 `json:"value,omitempty"`
}

// テクニカル分析用のデータの準備
func (df *DataFrameCandle) Times() []time.Time {
	s := make([]time.Time, len(df.Candles))
	for i, v := range df.Candles {
		s[i] = v.Time
	}
	return s
}

func (df *DataFrameCandle) Opens() []float64 {
	s := make([]float64, len(df.Candles))
	for i, v := range df.Candles {
		s[i] = v.Open
	}
	return s
}

func (df *DataFrameCandle) Highs() []float64 {
	s := make([]float64, len(df.Candles))
	for i, v := range df.Candles {
		s[i] = v.High
	}
	return s
}

func (df *DataFrameCandle) Lows() []float64 {
	s := make([]float64, len(df.Candles))
	for i, v := range df.Candles {
		s[i] = v.Low
	}
	return s
}

func (df *DataFrameCandle) Closes() []float64 {
	s := make([]float64, len(df.Candles))
	for i, v := range df.Candles {
		s[i] = v.Close
	}
	return s
}

func (df *DataFrameCandle) Swaps() []int {
	s := make([]int, len(df.Candles))
	for i, v := range df.Candles {
		s[i] = v.Swap
	}
	return s
}

// technicalの定義

// SMA
// dfにSMAの配列をappendする
// AddSmaはwebserverから呼び出す
func (df *DataFrameCandle) AddSma(period int) bool {

	if len(df.Candles) < period {
		return false
	}

	smaVal := talib.Sma(df.Closes(), period)
	df.Smas = append(df.Smas, Sma{
		Period: period,
		Value:  smaVal,
	})
	return true
}
