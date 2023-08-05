package models

import (
	"fmt"
	"log"
	"strings"
	"time"

	"example.com/tker-78/fx2/config"
)

/*
DataFrameCandleにSignalEventsのフィールドを持たせて、
SignalEventsのメソッドとして定義した、
#Buy, #Sellメソッドを呼び出して売買を記録する.

売買成立時には#Save()メソッドでデータベースにsignalEventを保存する.

売買のルールは、DataFrameCandleに保存された情報を起点にするので、
DataFrameCandleのメソッドとして定義する
(signalsに保存された最新のpriceから1円下がったら,
SMAやEMAのようなtechnicalsが条件を満たしたら、売買するのような条件を起点にするので。)

*/

// 売買情報を扱う型
// 売買条件に一致した場合に、signalEventをデータベースに保存する
type SignalEvent struct {
	Time         time.Time `json:"time"`
	CurrencyCode string    `json:"currency_code"`
	Side         string    `json:"side"`
	Price        float64   `json:"price"`
	Size         float64   `json:"size"`
}

type SignalEvents struct {
	Signals []SignalEvent `json:"signals,omitempty"`
}

// Newメソッド
func NewSignalEvents() *SignalEvents {
	return &SignalEvents{}
}

// signalをデータベースに保存する
func (s *SignalEvent) Save() bool {
	cmd := fmt.Sprintf(`
	INSERT INTO %s (time, currency_code, side, price, size) VALUES ($1, $2, $3, $4, $5)
	`, signalEventsTableName)
	_, err := DbConnection.Exec(cmd, s.Time.Format(time.RFC3339), s.CurrencyCode, s.Side, s.Price, s.Size)

	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			// すでにその日付のデータがある場合は、エラーを表示してtrueを返す
			log.Println(err)
			return true
		}
	}
	return true
}

// 最新のsignalを取得する
func (signals *SignalEvents) LastSignal() *SignalEvent {
	cmd := fmt.Sprintf(`
		SELECT * FROM %s
		ORDER BY time DESC LIMIT 1
	`, signalEventsTableName)

	row := DbConnection.QueryRow(cmd)
	s := &SignalEvent{}

	err := row.Scan(&s.Time, &s.CurrencyCode, &s.Side, &s.Price, &s.Size)
	if err != nil {
		log.Println(err)
	}
	return s
}

// 全てのsignalEventをデータベースから読み取る
func GetAllSignals() (*SignalEvents, error) {
	cmd := fmt.Sprintf(`
	SELECT * FROM %s
	`, signalEventsTableName)

	rows, err := DbConnection.Query(cmd)
	if err != nil {
		log.Println("error occured while querying Signals:", err)
		return nil, err
	}

	signalEvents := &SignalEvents{}
	for rows.Next() {
		s := SignalEvent{}
		err = rows.Scan(&s.Time, &s.CurrencyCode, &s.Side, &s.Price, &s.Size)
		if err != nil {
			log.Println("error occured while scanning:", err)
			return nil, err
		}

		signalEvents.Signals = append(signalEvents.Signals, s)
	}
	return signalEvents, err
}

// Signalsの利益/損失を返す
// Todo: df.CheckSellと連携する
func (signals *SignalEvents) Profit(currentPrice float64) float64 {

	profit := 0.0

	tmp_amount := 0.0

	for _, v := range signals.Signals {
		if v.Side == "BUY" {
			tmp_amount += v.Size * v.Price
		} else if v.Side == "SELL" {
			continue
		}
	}

	tmp_total_size := signals.TempTotalSize()
	profit = currentPrice*tmp_total_size - tmp_amount

	return profit
}

// 現時点の建玉のサイズを返す(BUYで計算)
func (signals *SignalEvents) TempTotalSize() float64 {
	tmp_total_size := 0.0
	for _, v := range signals.Signals {
		if v.Side == "BUY" {
			tmp_total_size += v.Size
		}
	}
	return tmp_total_size
}

func (signals *SignalEvents) FinalProfit() float64 {

	avg_price := 0.0
	sell_price := 0.0

	tmp := 0.0
	for _, v := range signals.Signals {
		if v.Side == "BUY" {
			tmp += v.Price * v.Size
		}
	}

	avg_price = tmp / signals.TempTotalSize()

	sell_price = signals.LastSignal().Price

	final_profit := (sell_price - avg_price) * signals.TotalBuySize()
	return final_profit
}

// BUYの建玉数を返す
func (signals *SignalEvents) TotalBuySize() float64 {
	total_size := 0.0

	for _, v := range signals.Signals {
		if v.Side == "BUY" {
			total_size += v.Size
		}
	}
	return total_size
}

// signalsを削除する
func DeleteSignals() bool {
	cmd := fmt.Sprintf(`
		DELETE FROM %s
	`, signalEventsTableName)
	_, err := DbConnection.Exec(cmd)
	if err != nil {
		log.Println("error occured while deleting SignalEvents:", err)
		return false
	}
	return true
}

// Buy関連メソッド

// signalsに関する購入条件制約
// 例えば、signalsの数が上限を超えていたら購入できない、など
func (signals *SignalEvents) CanBuy() bool {

	return true //temporary
}

func (signals *SignalEvents) Buy(dateTime time.Time, price, size float64, save bool) bool {
	if !signals.CanBuy() {
		return false
	}

	signalEvent := SignalEvent{
		Time:         dateTime,
		CurrencyCode: config.Config.CurrencyCode,
		Side:         "BUY",
		Price:        price,
		Size:         size,
	}

	if save {
		signalEvent.Save()
	}

	signals.Signals = append(signals.Signals, signalEvent)

	return true // temporary
}

// Sell関連メソッド

// signalsに関する売却条件制約
// 例えば、BUYのsignalがない場合は、売却できない、など
func (signals *SignalEvents) CanSell() bool {

	return true //temporary
}

func (signals *SignalEvents) Sell(dateTime time.Time, price float64, save bool) bool {
	if !signals.CanSell() {
		return false
	}

	signalEvent := SignalEvent{
		Time:         dateTime,
		CurrencyCode: config.Config.CurrencyCode,
		Side:         "SELL",
		Price:        price,
		Size:         signals.TotalBuySize(),
	}

	if save {
		signalEvent.Save()
	}

	signals.Signals = append(signals.Signals, signalEvent)

	return true // temporary
}
