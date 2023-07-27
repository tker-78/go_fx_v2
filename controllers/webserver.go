package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"example.com/tker-78/fx2/config"
	"example.com/tker-78/fx2/models"
)

func StartServer() error {
	http.HandleFunc("/api/candle/", apiCandleHandler)
	return http.ListenAndServe(fmt.Sprintf(":%d", config.Config.Port), nil)
}

func apiCandleHandler(w http.ResponseWriter, r *http.Request) {

	strLimit := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(strLimit)
	if err != nil || limit < 0 || strLimit == "" {
		limit = 30
	}

	df, err := models.GetCandlesByLimit(limit)

	if err != nil {
		log.Println("error occured while making dataframe", err)
	}

	w.Header().Set("Content-Type", "application/json")

	js, err := json.Marshal(df)
	if err != nil {
		log.Println("error occured while marhaling:", err)
	}
	w.Write(js)
}