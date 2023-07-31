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
	SELECT * FROM (
		SELECT * FROM %s
		ORDER BY time LIMIT 1
	)`, signalEventsTableName)

	row := DbConnection.QueryRow(cmd)
	s := &SignalEvent{}

	err := row.Scan(&s.Time, &s.CurrencyCode, &s.Side, &s.Price, &s.Size)
	if err != nil {
		log.Println(err)
	}
	return s
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

func (signals *SignalEvents) Sell(dateTime time.Time, price, size float64, save bool) bool {
	if !signals.CanSell() {
		return false
	}

	signalEvent := SignalEvent{
		Time:         dateTime,
		CurrencyCode: config.Config.CurrencyCode,
		Side:         "SELL",
		Price:        price,
		Size:         size,
	}

	if save {
		signalEvent.Save()
	}

	signals.Signals = append(signals.Signals, signalEvent)

	return true // temporary
}
