package models

import (
	"fmt"
	"log"
	"math"
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
	Results      []Result      `json:"results,omitempty"`
	Smas         []Sma         `json:"smas,omitempty"`
	Emas         []Ema         `json:"emas,omitempty"`
	BBands       *BBands       `json:"bbands,omitempty"`
	StochRSIs    *StochRSI     `json:"stoch_rsi,omitempty"`
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

type StochRSI struct {
	FastPeriod  []float64 `json:"fast_period,omitempty"`
	FastDPeriod []float64 `json:"fast_d_period,omitempty"`
	MaType      int       `json:"ma_type,omitempty"`
}

// シミュレーション結果を格納する
type Result struct {
	Entry         time.Time `json:"entry,omitempty"`
	Exit          time.Time `json:"exit,omitempty"`
	CapitalProfit float64   `json:"capital_profit,omitempty"`
	SwapProfit    float64   `json:"swap_profit,omitempty"`
	Duration      float64   `json:"duration,omitempty"`
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

// Stochastic
func (df *DataFrameCandle) AddStochastic(Period, FastPeriod, FastDPeriod int) bool {
	if Period <= len(df.Closes()) {
		outFastK, outFastD := talib.StochRsi(df.Closes(), Period, FastPeriod, FastDPeriod, 0)
		df.StochRSIs = &StochRSI{
			FastPeriod:  outFastK,
			FastDPeriod: outFastD,
			MaType:      0,
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

func (df *DataFrameCandle) CheckSell(currentPrice, swapProfit float64) bool {
	if df.Signals.Profit(currentPrice)+swapProfit < -50000 || df.Signals.Profit(currentPrice)+swapProfit > 30000 {
		return true
	} else {
		return false
	}
}

// Technicalを用いてエントリーポイントを選択する todo: 動作が変。最初の立ち上がりを回避したいけど回避できていない.
func (df *DataFrameCandle) ChooseStartCandleNumWithStochRSI() int {

	for i := 20; i < len(df.Candles); i++ {
		if df.StochRSIs.FastDPeriod[i] < 50 && df.StochRSIs.FastDPeriod[i+1] >= 50 {
			return i + 1
		}
	}
	return 0
}

// Todo: Start ~ Endの期間のdataframeを切り出す
func (df *DataFrameCandle) ExtractDataFrame() *DataFrameCandle {

	return nil
}

// メインのシミュレーション処理
func (df *DataFrameCandle) ExeSimWithStartDate() bool {
	DeleteSignals()
	startCandleNum := df.ChooseStartCandleNumWithStochRSI() // Todo: このエントリーポイントを、technicalを使って抽出できるようにする
	df.Signals.Buy(df.Candles[startCandleNum].Time, df.Candles[startCandleNum].Mid(), 1000, true)

	var total_swap_profit float64
	var lastCandleTime time.Time
	for i := 1; i < len(df.Candles)-startCandleNum; i++ {
		swap_profit := 0.0

		if df.Signals.LastSignal().Side == "SELL" {
			break
		}
		currentCandle := df.Candles[startCandleNum+i]

		swap_profit = df.Signals.TempTotalSize() * float64(currentCandle.Swap) / 10000

		total_swap_profit += swap_profit

		currentDuration, err := time.ParseDuration(currentCandle.Time.Sub(df.Signals.LastSignal().Time).String())
		if err != nil {
			log.Println("error occured while parsing duration from last signal", err)
		}
		currentDurationFromLastSignal := currentDuration.Hours() / 24

		if currentCandle.Low < df.Signals.LastSignal().Price-1 && len(df.Signals.Signals) < 10 {
			df.Signals.Buy(currentCandle.Time, df.Signals.LastSignal().Price-1, 1000, true)
		} else if currentDurationFromLastSignal > 15 && len(df.Signals.Signals) < 10 {
			df.Signals.Buy(currentCandle.Time, currentCandle.Mid(), 1000, true)
		} else if df.CheckSell(currentCandle.High, total_swap_profit) { // 利益が出る側での売却
			df.Signals.Sell(currentCandle.Time, df.Signals.SellPrice(30000, total_swap_profit), true)
			lastCandleTime = currentCandle.Time
		} else if df.CheckSell(currentCandle.Low, total_swap_profit) { //損失が出る側での売却
			df.Signals.Sell(currentCandle.Time, df.Signals.SellPrice(-50000, total_swap_profit), true)
			lastCandleTime = currentCandle.Time
		}

	}

	// resultに結果を格納
	d, _ := time.ParseDuration(lastCandleTime.Sub(df.Candles[startCandleNum].Time).String())
	days := d.Hours() / 24

	result := Result{
		Entry:         df.Candles[startCandleNum].Time,
		Exit:          lastCandleTime,
		CapitalProfit: df.Signals.FinalProfit(),
		SwapProfit:    total_swap_profit,
		Duration:      days,
	}

	result.Save()

	df.AddResults()

	// CandlesをEntry~Exitの期間で書き換え
	redf, err := GetCandlesByBetween(df.Results[len(df.Results)-1].Entry.AddDate(0, 0, -30), df.Results[len(df.Results)-1].Exit)
	if err != nil {
		log.Println("290", err)

	}

	df.Candles = redf.Candles

	return true

}

// 全てのSignalsをデータベースから読み込んで、dfに与える
func (df *DataFrameCandle) AddSignals() bool {
	// Todo
	signals, err := GetAllSignals()
	if err != nil {
		return false
	}
	df.Signals = signals
	return true
}

// Signalsの内容からシミュレーションの結果を分析して、Resultsに入れる
func (result *Result) Save() bool {
	cmd := fmt.Sprintf(`
	INSERT INTO %s (entry, exit, capital_profit, swap_profit, duration) VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (entry) DO UPDATE
	SET exit = $2, capital_profit = $3, swap_profit = $4, duration = $5;
	`, simulationResultsTableName)

	_, err := DbConnection.Exec(cmd, result.Entry, result.Exit, math.Round(result.CapitalProfit), math.Round(result.SwapProfit), result.Duration)
	if err != nil {
		log.Println("error occured while saving result:", err)
		return false
	}
	return true
}

// 全てのResultsをデータベースから読み込んで、dfに与える
func (df *DataFrameCandle) AddResults() bool {

	cmd := fmt.Sprintf(`SELECT * FROM %s`, simulationResultsTableName)

	rows, err := DbConnection.Query(cmd)
	if err != nil {
		log.Println("error occured while querying for simuration results:", err)
		return false
	}
	defer rows.Close()

	for rows.Next() {
		r := Result{}
		rows.Scan(&r.Entry, &r.Exit, &r.CapitalProfit, &r.SwapProfit, &r.Duration)
		df.Results = append(df.Results, r)
	}

	return true
}

// Resultsをリセットする
func (df *DataFrameCandle) DeleteResults() bool {

	cmd := fmt.Sprintf(`DELETE FROM %s`, simulationResultsTableName)
	_, err := DbConnection.Exec(cmd)
	if err != nil {
		log.Println("error occured while deleting results:", err)
		return false
	}
	return true
}
