package main

import controller "example.com/tker-78/fx2/controllers"

func main() {
	// 初回のみの読み込み
	// models.LoadCSV()

	controller.StartServer()

	/*
		以降はデバッグ用temporary code
	*/

	// dateTime, _ := time.Parse("2006-01-02 03:04:05", "2011-03-25 01:02:12")

	// start := time.Date(2011, 03, 25, 0, 0, 0, 0, time.UTC)
	// end := time.Date(2011, 05, 25, 0, 0, 0, 0, time.UTC)
	// fmt.Println(models.GetCandlesByBetween(start, end))

	// var startTime time.Time = time.Date(2022, 12, 01, 00, 00, 00, 00, time.UTC)
	// df, _ := models.GetCandlesByLimit(10)
	// df.Signals = models.NewSignalEvents()
	// df.Signals.Buy(startTime, df.Candles[0].High, 1, true)
	// fmt.Println(df)
}
