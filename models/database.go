package models

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"example.com/tker-78/fx2/config"
	_ "github.com/lib/pq"
)

var DbConnection *sql.DB

const tableName = "USD_JPY_1d"

func init() {
	var err error
	connectionStr := "user=takuyakinoshita dbname=exchange_2 sslmode=disable"
	DbConnection, err := sql.Open(config.Config.SQLDriver, connectionStr)
	if err != nil {
		log.Fatalln("Error occured while opening database file: ", err)
	}

	cmd := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			time TIMESTAMP PRIMARY KEY NOT NULL,
			open FLOAT,
			high FLOAT,
			low FLOAT,
			close FLOAT,
			swap FLOAT)
	`, tableName)

	_, err = DbConnection.Exec(cmd)
	if err != nil {
		log.Fatalln("Error occuredd while creating database table: ", err)
	}
}

// CSVファイルから為替データの流し込み
func LoadCSV() {
	// cmd := fmt.Sprintf(`
	// 	INSERT INTO %s (time, open, high, low, close, swap) VALUES ($1, $2, $3, $4, $5, $6)
	// `, tableName)

	// dataフォルダに格納しているファイルすべてのファイルパスをpathsに格納する

	files, err := ioutil.ReadDir("./data")
	if err != nil {
		log.Fatalln("Error occured while reading 'data' directory", err)
	}
	var paths []string

	for _, file := range files {
		paths = append(paths, filepath.Join("./data/", file.Name()))
	}
	fmt.Println(paths)

	for _, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			log.Println(err)
		}
		defer file.Close()
		r := csv.NewReader(file)
		rows, _ := r.ReadAll()

		for _, row := range rows {
			timeTime, err := time.Parse("2006-01-02", row[0])
			if err != nil {
				log.Fatalln("Error occured while parsing time: ", err)
			}

			timeString := timeTime.Format(time.RFC3339)

			open, err := strconv.ParseFloat(row[1], 64)
			if err != nil {
				log.Fatalln("Error occured while parsing open: ", err)
			}

			high, err := strconv.ParseFloat(row[2], 64)
			if err != nil {
				log.Fatalln("Error occured while parsing high: ", err)
			}

			low, err := strconv.ParseFloat(row[3], 64)
			if err != nil {
				log.Fatalln("Error occured while parsiing low: ", err)
			}

			close, err := strconv.ParseFloat(row[4], 64)
			if err != nil {
				log.Fatalln("Error occured while parsing close: ", err)
			}

			swap, err := strconv.ParseInt(row[5], 64, 64)
			if err != nil {
				log.Fatalln("Error occured while parsing swap: ", err)
			}
		}
	}

}
