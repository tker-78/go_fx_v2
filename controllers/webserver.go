package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"example.com/tker-78/fx2/config"
	"example.com/tker-78/fx2/models"
)

func StartServer() error {
	http.HandleFunc("/api/candle/", apiCandleHandler)
	return http.ListenAndServe(fmt.Sprintf(":%d", config.Config.Port), nil)
}

func apiCandleHandler(w http.ResponseWriter, r *http.Request) {
	var df *models.DataFrameCandle
	// limitでdfの抽出
	// startかendが指定されていない場合のみ、実行する
	strLimit := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(strLimit)
	if err != nil || limit < 0 || strLimit == "" {
		limit = 30
	}
	if err != nil {
		log.Println("error occured while making dataframe", err)
	}

	// start-endでdfの抽出
	// limitよりも優先される
	strStart := r.URL.Query().Get("start")
	startDate, err := time.Parse("2006-01-02", strStart)
	if err != nil {
		log.Println(err)
	}

	strEnd := r.URL.Query().Get("end")
	endDate, err := time.Parse("2006-01-02", strEnd)
	if err != nil {
		log.Println(err)
	}

	if strStart != "" && strEnd == "" {
		startDate, err := time.Parse("2006-01-02", strStart)
		if err != nil {
			log.Println(err)
		}
		df, err = models.GetCandlesAfterTime(startDate)
		if err != nil {
			log.Println(err)
		}
	} else if strStart == "" && strEnd == "" {
		df, err = models.GetCandlesByLimit(limit)
		if err != nil {
			log.Println(err)
		}
	} else {
		df, err = models.GetCandlesByBetween(startDate, endDate)
		if err != nil {
			log.Println(err)
		}
	}

	// SMA
	periodStr1 := r.URL.Query().Get("period1")
	periodStr2 := r.URL.Query().Get("period2")
	periodStr3 := r.URL.Query().Get("period3")

	period1, err := strconv.Atoi(periodStr1)
	if err != nil || period1 < 0 || periodStr1 == "" {
		period1 = 7
	}

	period2, err := strconv.Atoi(periodStr2)
	if err != nil || period2 < 0 || periodStr2 == "" {
		period2 = 14
	}

	period3, err := strconv.Atoi(periodStr3)
	if err != nil || period3 < 0 || periodStr3 == "" {
		period3 = 25
	}

	df.AddSma(period1)
	df.AddSma(period2)
	df.AddSma(period3)

	// EMA
	df.AddEma(period1)
	df.AddEma(period2)
	df.AddEma(period3)

	// BBands
	bbnStr := r.URL.Query().Get("bbn")
	bbkStr := r.URL.Query().Get("bbk")

	bbn, err := strconv.Atoi(bbnStr)
	if err != nil || bbn < 0 || bbnStr == "" {
		bbn = 20
	}

	bbk, err := strconv.Atoi(bbkStr)
	if err != nil || bbk < 0 || bbkStr == "" {
		bbk = 2
	}

	df.AddBBands(bbn, float64(bbk))

	// Signals関連
	// Todo: APIとExeSimWithStartDate()の連携

	models.DeleteSignals()
	df.AddSignals()
	df.ExeSimWithStartDate()
	fmt.Println(df.Results)

	// CORSの設定
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5500")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,UPDATE,OPTIONS")

	w.Header().Set("Content-Type", "application/json")

	js, err := json.Marshal(df)
	if err != nil {
		log.Println("error occured while marhaling:", err)
	}
	w.Write(js)
}
