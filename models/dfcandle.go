package models

import (
	"time"

	"github.com/markcheno/go-talib"
)

/*
データベースからのcandleデータの抽出と、
売買ルールの指定はDataDrameCandleで定義する

*/

// データベースから指定した条件のcandleを格納するための型
type DataFrameCandle struct {
	CurrencyCode string        `json:"currency_code"`
	Duration     string        `json:"duration"`
	Candles      []Candle      `json:"candles"`
	Signals      *SignalEvents `json:"signals,omitempty"`
	Smas         []Sma         `json:"smas,omitempty"`
	Emas         []Ema         `json:"emas,omitempty"`
	BBands       *BBands       `json:"bbands,omitempty"`
}

// technical用の型
type Sma struct {
	Period int       `json:"period,omitempty"`
	Value  []float64 `json:"value,omitempty"`
}

type Ema struct {
	Period int       `json:"period,omitempty"`
	Value  []float64 `json:"value,omitempty"`
}

type BBands struct {
	N    int       `json:"n,omitempty"`
	K    float64   `json:"k,omitempty"`
	Up   []float64 `json:"up,omitempty"`
	Mid  []float64 `json:"mid,omitempty"`
	Down []float64 `json:"down,omitempty"`
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

func (df *DataFrameCandle) MidPrices() []float64 {
	s := make([]float64, len(df.Candles))
	for i, v := range df.Candles {
		s[i] = (v.High + v.Low) / 2
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

// EMA
// SMAと同じ実装方法
func (df *DataFrameCandle) AddEma(period int) bool {
	if len(df.Candles) < period {
		return false
	}
	emaVal := talib.Ema(df.Closes(), period)
	df.Emas = append(df.Emas, Ema{
		Period: period,
		Value:  emaVal,
	})
	return true
}

// BBand
func (df *DataFrameCandle) AddBBands(n int, k float64) bool {
	if n <= len(df.Closes()) {
		up, mid, down := talib.BBands(df.Closes(), n, k, k, 0)
		df.BBands = &BBands{
			N:    n,
			K:    k,
			Up:   up,
			Mid:  mid,
			Down: down,
		}
		return true
	}
	return false
}

// 売買ルールの指定
// この引数のtimeTimeを1日ごとに送ることで、シミュレーションができる
func (df *DataFrameCandle) BuyRule(timeTime time.Time) bool {
	if !df.Signals.CanBuy() {
		return false
	}

	candle := GetCandle(timeTime)
	/*
		前回の購入金額よりも1円下がったら追加購入する LastSignal()メソッドで判定
		購入回数は10回まで

	*/
	if candle.Low < df.Signals.LastSignal().Price-1 && len(df.Signals.Signals) < 10 {
		return true
	} else {
		return false
	}

}

func (df *DataFrameCandle) SellRule() bool {
	return true // temporary
}

// Signalsの追加
func (df *DataFrameCandle) AddSignals() {
	// Todo
}
