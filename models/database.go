package models

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"example.com/tker-78/fx2/config"
	_ "github.com/lib/pq"
)

var DbConnection *sql.DB

const (
	tableName                  = "USD_JPY_1D"
	signalEventsTableName      = "signal_events"
	simulationResultsTableName = "sim_results"
)

// "1m"を指定したら、"USD_JPY_1m0s"を出力
func GetTableName(duration string) string {
	durationName := config.Config.Durations[duration].String()

	return fmt.Sprintf(`USD_JPY_%s`, durationName)

}

func init() {
	var err error
	connectionStr := "user=takuyakinoshita dbname=exchange_2 sslmode=disable"
	DbConnection, err = sql.Open(config.Config.SQLDriver, connectionStr)
	if err != nil {
		log.Fatalln("Error occured while opening database file: ", err)
	}

	tableName := GetTableName("1m")

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

	c := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		time TIMESTAMP PRIMARY KEY NOT NULL,
		currency_code VARCHAR,
		side VARCHAR,
		price FLOAT,
		size FLOAT)
	`, signalEventsTableName)

	_, err = DbConnection.Exec(c)

	if err != nil {
		log.Fatalln("error occured while creating table:", err)
	}

	c2 := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		entry TIMESTAMP PRIMARY KEY NOT NULL, 
		exit TIMESTAMP,
		capital_profit FLOAT,
		swap_profit FLOAT,
		duration VARCHAR
	)
	`, simulationResultsTableName)

	_, err = DbConnection.Exec(c2)
	if err != nil {
		log.Fatalln("error occured while creating table:", err)
	}

}

// CSVファイルから為替データの流し込み
func LoadCSV() {
	c := fmt.Sprintf(`
		INSERT INTO %s (time, open, high, low, close, swap) VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT DO NOTHING
	`, tableName)

	// dataフォルダに格納しているファイルすべてのファイルパスをpathsに格納する

	files, err := os.ReadDir("./data")
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

		for _, row := range rows[1:] {
			timeTime, err := time.Parse("2006/01/02", row[0])
			if err != nil {
				log.Println("Error occured while parsing time: ", err)
			}

			timeString := timeTime.Format(time.RFC3339)

			open, err := strconv.ParseFloat(row[1], 64)
			if err != nil {
				log.Println("Error occured while parsing open: ", err)
			}

			high, err := strconv.ParseFloat(row[2], 64)
			if err != nil {
				log.Println("Error occured while parsing high: ", err)
			}

			low, err := strconv.ParseFloat(row[3], 64)
			if err != nil {
				log.Println("Error occured while parsiing low: ", err)
			}

			close, err := strconv.ParseFloat(row[4], 64)
			if err != nil {
				log.Println("Error occured while parsing close: ", err)
			}

			swap, err := strconv.Atoi(row[5])
			if err != nil {
				log.Println("Error occured while parsing swap: ", err)
			}
			if err != nil || row[5] == "" {
				swap = 0
			}

			_, err = DbConnection.Exec(c, timeString, open, high, low, close, swap)
			if err != nil {
				log.Println(err)
			}

		}
	}

}

// CSVファイルから為替データの流し込み
func LoadM1CSV() {
	table := GetTableName("1m")
	c := fmt.Sprintf(`
		INSERT INTO %s (time, open, high, low, close, swap) VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT DO NOTHING
	`, table)

	// dataフォルダに格納しているファイルすべてのファイルパスをpathsに格納する

	files, err := os.ReadDir("./data/m1")
	if err != nil {
		log.Fatalln("Error occured while reading 'data' directory", err)
	}
	var paths []string

	for _, file := range files {
		paths = append(paths, filepath.Join("./data/m1/", file.Name()))
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

		for _, row := range rows[1:] {
			timeTime, err := time.Parse("2006-01-02 15:04", row[0])
			if err != nil {
				log.Println("Error occured while parsing time: ", err)
			}

			timeString := timeTime.Format(time.RFC3339)

			open, err := strconv.ParseFloat(row[1], 64)
			if err != nil {
				log.Println("Error occured while parsing open: ", err)
			}

			high, err := strconv.ParseFloat(row[2], 64)
			if err != nil {
				log.Println("Error occured while parsing high: ", err)
			}

			low, err := strconv.ParseFloat(row[3], 64)
			if err != nil {
				log.Println("Error occured while parsiing low: ", err)
			}

			close, err := strconv.ParseFloat(row[4], 64)
			if err != nil {
				log.Println("Error occured while parsing close: ", err)
			}

			swap, err := strconv.Atoi(row[5])
			if err != nil {
				log.Println("Error occured while parsing swap: ", err)
			}
			if err != nil || row[5] == "" {
				swap = 0
			}

			_, err = DbConnection.Exec(c, timeString, open, high, low, close, swap)
			if err != nil {
				log.Println(err)
			}

		}
	}
	fmt.Println("load finished")

}
