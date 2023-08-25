package main

import (
	"fmt"
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

	go spiner(100*time.Millisecond, "loading candle stick data...")

	models.CreateCandleWithDuration("5m", "1m")

	models.CreateCandleWithDuration("30m", "5m")
	models.CreateCandleWithDuration("1h", "30m")
	models.CreateCandleWithDuration("4h", "1h")
	models.CreateCandleWithDuration("24h", "4h")
	// controllers.StartServer()

	/*
		以降はデバッグ用temporary code
	*/

	// timeTime := time.Date(2022, 12, 30, 16, 00, 00, 00, time.Local)
	// df, err := models.GetCandlesAfterTime(timeTime)
	// if err != nil {
	// 	log.Println(err)
	// }

	// df.AddEma(5)
	// df.AddEma(10)
	// df.AddEma(25)
	// df.AddSignals()

	// ema1 := df.Emas[0].Value
	// ema2 := df.Emas[1].Value

	// for i := 1; i < len(df.Candles); i++ {
	// 	candle := df.Candles[i-1]
	// 	if ema1[i-1] <= ema2[i-1] && ema1[i] > ema2[i] {
	// 		df.Signals.Buy(candle.Time, candle.Mid(), 1000, true)
	// 	}
	// }
	// fmt.Println(df.Signals)

}
