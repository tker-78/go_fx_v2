package main

import (
	"fmt"
	"log"
	"time"

	"example.com/tker-78/fx2/models"
)

func spiner(delay time.Duration, text string) {
	for {
		for _, r := range `-\|/` {
			fmt.Printf("\r%c %s", r, text)
			time.Sleep(delay)
		}
	}
}

func main() {
	// 初回のみの読み込み
	// models.LoadM1CSV()

	// for key, _ := range config.Config.Durations {
	// 	if key == "1m" {
	// 		continue
	// 	}
	// 	go models.CreateCandleWithDuration(key)
	// }

	// controllers.StartServer()

	/*
		以降はデバッグ用temporary code
	*/

	// 30分足でシグナルに基づいて売買して、その結果をデータベースに保存する
	// 売買結果はデータベースで確認する

	// backtest ema
	timeTime := time.Date(2010, 01, 01, 00, 00, 00, 00, time.Local)
	df, err := models.GetCandlesAfterTime(timeTime, "30m")
	if err != nil {
		log.Println(err)
	}

	df.AddEma(5)
	df.AddEma(10)
	df.AddEma(25)
	df.AddSignals()

	ema1 := df.Emas[0].Value
	ema2 := df.Emas[1].Value

	for i := 1; i < len(df.Candles); i++ {
		candle := df.Candles[i-1]
		if ema1[i-1] <= ema2[i-1] && ema1[i] > ema2[i] {
			df.Signals.Buy(candle.Time, candle.Mid(), 1000, true)
		}

		if df.Signals.Profit(candle.Mid()) > 5000 {
			fmt.Println(df.Signals.Profit(candle.Mid()))
			df.Signals.Sell(candle.Time, candle.Mid(), true)
		}

		if df.Signals.Profit(candle.Mid()) < -5000 {
			fmt.Println(df.Signals.Profit(candle.Mid()))
			df.Signals.Sell(candle.Time, candle.Mid(), true)
		}
	}

	fmt.Println(df.Signals)

}
