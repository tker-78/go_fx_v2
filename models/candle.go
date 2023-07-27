package models

import (
	"fmt"
	"time"

	"example.com/tker-78/fx2/config"
)

type Candle struct {
	Time  time.Time `json:"time"`
	Open  float64   `json:"open"`
	High  float64   `json:"high"`
	Low   float64   `json:"low"`
	Close float64   `json:"close"`
	Swap  int       `json:"swap"`
}

// ex) usd_jpy_1d
func GetTableName() string {
	currency_code := config.Config.CurrencyCode
	duration := config.Config.Duration
	return currency_code + "_" + duration
}

// 時刻を日付の形式に切り捨てて、RFC3339形式にフォーマットする
// データベースからの情報読み出し時に使用する

func TruncateTimeToDate(timeTime time.Time) string {
	return timeTime.Truncate(24 * time.Hour).Format(time.RFC3339)
}

// candleを一つ返す(日付で指定)
func GetCandle(timeTime time.Time) *Candle {
	tableName := GetTableName()
	truncatedTime := TruncateTimeToDate(timeTime)
	cmd := fmt.Sprintf(`
	SELECT * FROM %s WHERE time = $1
	`, tableName)

	row := DbConnection.QueryRow(cmd, truncatedTime)

	var candle Candle

	err := row.Scan(&candle.Time, &candle.Open, &candle.High, &candle.Low, &candle.Close, &candle.Swap)
	if err != nil {
		fmt.Println(err)
	}
	return &Candle{
		Time:  candle.Time,
		Open:  candle.Open,
		High:  candle.Open,
		Low:   candle.Low,
		Close: candle.Close,
		Swap:  candle.Swap,
	}
}

// 指定した期間のDataFrameCandleを返す
// limitで最新側の取得
func GetCandlesByLimit(limit int) (*DataFrameCandle, error) {
	tableName := GetTableName()
	cmd := fmt.Sprintf(`
	SELECT * FROM (
		SELECT time, open, high, low, close, swap FROM %s 
		ORDER BY time DESC limit $1
	) AS t1
	ORDER BY time ASC;
	`, tableName)

	rows, err := DbConnection.Query(cmd, limit)
	if err != nil {
		fmt.Println("error occured while query", err)
		return nil, err
	}
	defer rows.Close()

	dfCandle := &DataFrameCandle{}
	dfCandle.CurrencyCode = config.Config.CurrencyCode
	dfCandle.Duration = config.Config.Duration
	for rows.Next() {
		candle := Candle{}
		rows.Scan(&candle.Time, &candle.Open, &candle.High, &candle.Low, &candle.Close, &candle.Swap)
		dfCandle.Candles = append(dfCandle.Candles, candle)
	}
	return dfCandle, err
}

// 指定した期間のDataFrameCandleを返す
// between
func GetCandlesByBetween(start, end time.Time) (*DataFrameCandle, error) {
	tableName := GetTableName()
	startDate := TruncateTimeToDate(start)
	endDate := TruncateTimeToDate(end)

	cmd := fmt.Sprintf(`
	SELECT * FROM %s 
	WHERE time BETWEEN $1 AND $2
	`, tableName)

	rows, err := DbConnection.Query(cmd, startDate, endDate)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	dfCandle := &DataFrameCandle{}
	dfCandle.CurrencyCode = config.Config.CurrencyCode
	dfCandle.Duration = config.Config.Duration
	for rows.Next() {
		candle := Candle{}
		rows.Scan(&candle.Time, &candle.Open, &candle.High, &candle.Low, &candle.Close, &candle.Swap)
		dfCandle.Candles = append(dfCandle.Candles, candle)
	}
	return dfCandle, err
}
