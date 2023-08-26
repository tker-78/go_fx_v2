package main

import (
	"fmt"
	"time"

	"example.com/tker-78/fx2/controllers"
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

	controllers.StartServer()

	/*
		以降はデバッグ用temporary code
	*/

	// 30分足でシグナルに基づいて売買して、その結果をデータベースに保存する
	// 売買結果はデータベースで確認する

	// backtest ema

	// startTime := time.Date(2020, 01, 01, 00, 00, 00, 00, time.Local)
	// endTime := time.Date(2020, 01, 31, 00, 00, 00, 00, time.Local)
	// df, err := models.GetCandlesByBetween(startTime, endTime, "30m")
	// if err != nil {
	// 	log.Println(err)
	// }

	// signals := df.BacktestEma(5, 12)
	// fmt.Println(signals)
	// fmt.Println("total profit:", signals.ParseProfit())

}
