package main

import "example.com/tker-78/fx2/controllers"

func main() {
	// 初回のみの読み込み
	// models.LoadM1CSV()

	// models.CreateCandleWithDuration("1m")

	// for key, _ := range config.Config.Durations {
	// 	if key == "1m" {
	// 		continue
	// 	}
	// 	models.CreateCandleWithDuration(key)
	// }

	controllers.StartServer()

	/*
		以降はデバッグ用temporary code
	*/

}
