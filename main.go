package main

import "example.com/tker-78/fx2/controllers"

func main() {
	// 初回のみの読み込み
	// models.LoadCSV()

	controllers.StartServer()

	/*
		以降はデバッグ用temporary code
	*/

	// dateTime, _ := time.Parse("2006-01-02 03:04:05", "2011-03-25 01:02:12")

	// start := time.Date(2011, 03, 25, 0, 0, 0, 0, time.UTC)
	// end := time.Date(2011, 05, 25, 0, 0, 0, 0, time.UTC)
	// fmt.Println(models.GetCandlesByBetween(start, end))

	// df, _ := models.GetCandlesByLimit(10)
	// df.Signals = models.NewSignalEvents()

	// var startTime time.Time = time.Date(2022, 12, 02, 00, 00, 00, 00, time.UTC)
	// for i := 0; i < 10; i++ {
	// 	currentDate := startTime.AddDate(0, 0, i)
	// 	df.Signals.Buy(currentDate, df.Candles[i].High, 1000, true)
	// }
	// fmt.Println(df.Signals)
}
