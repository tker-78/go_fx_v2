package main

import (
	"fmt"
	"time"

	"example.com/tker-78/fx2/models"
)

func main() {
	// 初回のみの読み込み
	// models.LoadCSV()

	dateTime, _ := time.Parse("2006-01-02 03:04:05", "2011-03-25 01:02:12")
	fmt.Println(models.GetCandle(dateTime))
	fmt.Println(models.GetCandlesByLimit(10))
}
