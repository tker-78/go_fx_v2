package main

import (
	"fmt"
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

}
