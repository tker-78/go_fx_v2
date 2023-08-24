package main

import (
	"fmt"
	"log"
	"time"

	"example.com/tker-78/fx2/models"
)

func spiner(delay time.Duration) {
	for {
		for _, r := range `-\|/` {
			fmt.Printf("\r%c", r)
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
	// 	models.CreateCandleWithDuration(key, "1m")
	// }

	go spiner(100 * time.Millisecond)

	models.CreateCandleWithDuration("5m", "1m")

	// controllers.StartServer()

	/*
		以降はデバッグ用temporary code
	*/

	timeTime := time.Date(2022, 10, 01, 00, 00, 00, 00, time.Local)
	df, err := models.GetCandlesAfterTime(timeTime)
	if err != nil {
		log.Println(err)
	}

	df.AddEma(5)
	df.AddEma(10)
	df.AddEma(25)

	for i := 0; i < len(df.Candles); i++ {
		ema1 := df.Emas[0].Value
		ema2 := df.Emas[1].Value

		candle := df.Candles[i]
		if ema1[i] <= ema2[i] && ema1[i+1] > ema2[i+1] {
			df.Signals.Buy(candle.Time, candle.Mid(), 1000, true)
		}
	}

}
