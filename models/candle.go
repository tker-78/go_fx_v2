package models

import (
	"fmt"
	"log"
	"math"
	"time"

	"example.com/tker-78/fx2/config"
)

type Candle struct {
	Duration    time.Duration `json:"duration"`
	DurationKey string        `json:"duration_key"`
	Time        time.Time     `json:"time"`
	Open        float64       `json:"open"`
	High        float64       `json:"high"`
	Low         float64       `json:"low"`
	Close       float64       `json:"close"`
	Swap        float64       `json:"swap"`
}

func NewCandle(durationName string, timeTime time.Time, open, high, low, close, swap float64) *Candle {
	duration := config.Config.Durations[durationName]
	return &Candle{
		duration,
		durationName,
		timeTime,
		open,
		high,
		low,
		close,
		swap,
	}
}

func (candle Candle) Mid() float64 {
	return (candle.High + candle.Low) / 2
}

// 時刻を日付の形式に切り捨てて、RFC3339形式にフォーマットする
// データベースからの情報読み出し時に使用する

func TruncateTimeToDate(timeTime time.Time) string {
	return timeTime.Truncate(24 * time.Hour).Format(time.RFC3339)
}

func TruncateTimeToDuration(timeTime time.Time, durationName string) string {
	duration := config.Config.Durations[durationName]
	return timeTime.Truncate(duration).Format(time.RFC3339)
}

// candleを保存
func (candle *Candle) Save() bool {

	tableName := GetTableName(candle.DurationKey)

	cmd := fmt.Sprintf(`
	INSERT INTO %s (time, open, high, low, close, swap) VALUES($1, $2, $3, $4, $5, $6)
	ON CONFLICT (time) DO UPDATE SET open = $2, high = $3, low = $4, close = $5, swap = $6
	`, tableName)

	_, err := DbConnection.Exec(cmd, candle.Time, candle.Open, candle.High, candle.Low, candle.Close, candle.Swap)
	if err != nil {
		log.Println("error occured while inserting candle:", err)
		return false
	}
	return true
}

// candleを一つ返す(日付で指定)
func GetCandle(timeTime time.Time, durationName string) *Candle {
	tableName := GetTableName(durationName)
	fmt.Println(tableName)
	truncatedTime := TruncateTimeToDuration(timeTime, durationName)
	cmd := fmt.Sprintf(`
	SELECT * FROM %s WHERE time = $1
	`, tableName)

	row := DbConnection.QueryRow(cmd, truncatedTime)

	var candle Candle

	err := row.Scan(&candle.Time, &candle.Open, &candle.High, &candle.Low, &candle.Close, &candle.Swap)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &Candle{
		Duration:    config.Config.Durations[durationName],
		DurationKey: durationName,
		Time:        candle.Time,
		Open:        candle.Open,
		High:        candle.Open,
		Low:         candle.Low,
		Close:       candle.Close,
		Swap:        candle.Swap,
	}
}

// 指定した期間のDataFrameCandleを返す
// limitで最新側の取得
func GetCandlesByLimit(limit int) (*DataFrameCandle, error) {
	tableName := GetTableName("1m")
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
	dfCandle.Duration = config.Config.TradeDuration
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
	tableName := GetTableName("1m")
	startDate := TruncateTimeToDate(start)
	endDate := TruncateTimeToDate(end)

	cmd := fmt.Sprintf(`
	SELECT * FROM %s 
	WHERE time BETWEEN $1 AND $2
	`, tableName)

	rows, err := DbConnection.Query(cmd, startDate, endDate)
	if err != nil {
		fmt.Println("error occured while calling GetCandlesByBetween", err)
		return nil, err
	}
	defer rows.Close()

	dfCandle := &DataFrameCandle{}
	dfCandle.CurrencyCode = config.Config.CurrencyCode
	dfCandle.Duration = config.Config.TradeDuration
	for rows.Next() {
		candle := Candle{}
		rows.Scan(&candle.Time, &candle.Open, &candle.High, &candle.Low, &candle.Close, &candle.Swap)
		dfCandle.Candles = append(dfCandle.Candles, candle)
	}
	return dfCandle, err
}

// 指定した日付以降のDataFrameCandleを返す
func GetCandlesAfterTime(dateTime time.Time) (*DataFrameCandle, error) {
	tableName := GetTableName("1m")
	startDate := TruncateTimeToDate(dateTime)

	cmd := fmt.Sprintf(`
		SELECT time, open, high, low, close, swap FROM %s 
		WHERE time >= $1
		ORDER BY time 
	`, tableName)

	rows, err := DbConnection.Query(cmd, startDate)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	dfCandle := &DataFrameCandle{}
	dfCandle.CurrencyCode = config.Config.CurrencyCode
	dfCandle.Duration = config.Config.TradeDuration

	for rows.Next() {
		var candle Candle
		rows.Scan(&candle.Time, &candle.Open, &candle.High, &candle.Low, &candle.Close, &candle.Swap)
		dfCandle.Candles = append(dfCandle.Candles, candle)
	}
	return dfCandle, err
}

func CreateCandleWithDuration(durationName string) bool {
	// 1mのデータベースから値を読み出して、所定のデータベースに格納する
	tableName := GetTableName("1m")
	duration := config.Config.Durations[durationName]

	cmd := fmt.Sprintf(`
	SELECT * FROM %s	
	`, tableName)

	rows, err := DbConnection.Query(cmd)
	if err != nil {
		log.Println("error occured while querying to 1m table:", err)
		return false
	}
	defer rows.Close()

	for rows.Next() {
		var candle Candle
		rows.Scan(&candle.Time, &candle.Open, &candle.High, &candle.Low, &candle.Close, &candle.Swap)

		currentCandle := GetCandle(candle.Time, durationName)
		price := candle.Mid()

		if currentCandle == nil {
			currentCandle = NewCandle(durationName, candle.Time.Truncate(duration), price, price, price, price, candle.Swap)
		}

		if price > currentCandle.High {
			currentCandle.High = math.Round(price*100) / 100
		} else if price < currentCandle.Low {
			currentCandle.Low = math.Round(price*100) / 100
		}
		currentCandle.Close = math.Round(price*100) / 100
		currentCandle.Save()
	}

	return true

}
