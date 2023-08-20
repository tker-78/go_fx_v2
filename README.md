# Go FX V2


## front側の起動

```bash
yarn serve --port 5500
```

## app側の起動

```bash
go run main.go
```


## これはなに

go_fx_v1を再構築する. 
変化点は下記の通り. 

- 日足のトレード専用にする
- 扱う通貨はUSD/JPYのみ
- talibのテクニカルを使って、売買ルールを指定する
- 売買ルールのベースはRubyで構築したトレードシミュレーションと基本は同じにする
- Frontはおまけ
  - シミュレーションはバックエンドとデータベースだけで完結させる
  - Frontに表示したくなったらAPIHandlerを定義して、Vue.jsからAPI呼び出しをするようにする



その他

- databaseはpostgresql
- 日足データはクリック証券かな(??)から取ってきたものを使う


### config.iniの設定

```ini
[fxtrading]
currency_code = USD_JPY


[db]
SQLDriver = postgres
DbName = exchange_v2

```

config.go
```go
package config

import (
	"log"

	"gopkg.in/ini.v1"
)

type ConfigList struct {
	SQLDriver    string
	DbName       string
	CurrencyCode string
}

var Config ConfigList

func init() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalln(err)
	}

	Config = ConfigList{
		SQLDriver:    cfg.Section("db").Key("SQLDriver").String(),
		DbName:       cfg.Section("db").Key("DbName").String(),
		CurrencyCode: cfg.Section("fxtrading").Key("currency_code").String(),
		Port : cfg.Section("web").Key("port").MustInt(),
	}

}

```



### データベースの設定
postgresqlを使用する(sqlite3よりは使い慣れているので)

```go
import (
  "github.com/lib/pq"
)


func init() {
  conStr := "user=takuyakinoshita dbname=exchangerates sslmode=disable"
  DbConnection, err = sql.Open(config.Config.SQLDriver, conStr)
  // ここで := を使うとnil pointer dereference errorが出るので注意
}
```


### Databaseへのデータ取り込み
1dのローソクからデータベースを作成する

SQL.open時には、connectionとerrorはグローバル登録して、
型推論はしないようにしないとエラーが発生するので注意する. 

```go
// 下記はぬるぽ
DbConnection, err := sql.Open(config.Config.SQLDriver, connectionStr)
```


### DataFrameの作成
データベースから情報を取ってきてメモリに格納するDataFrameを作成する.

`candle.go`
データベースから値を読み出して、candle構造体に1本のcandle情報を格納する
- time
- open
- high
- low
- close
- swap

データベースから値を読み出して、DataFrameCandle構造体を返すグローバル関数もここに定義する.  



`dfcandle.go`
candleの配列の形をとる.  
- duration
- events(購買情報,後で追加)
- SMA(technical, 後で追加)
- EMA(technical, 後で追加)
- ...
- candles []candle



[models/dfcandle.go]
```go


```

[models/candle.go]


### SignalEventsの実装

抜き出したデータフレームにSignalEventsを渡す.  

個々のSignalEventは、データフレーム内で条件に合致した場合に
データベースに売買情報が保存される.  
なので、SignalEventsテーブルを参照することで、
売買ルールの有効性が確認できる.  



### シミュレーションの実行
- シミュレーションの開始時点を指定: startDate
- 保有ロット、建玉の価格、現在の価格profitを計算して、Accountの値をウォッチする
	(単純化のため、ひとまずswapは計算対象にしない)

startDataだけを指定して、startDate以降のデータをすべて抽出。  
startDateを最初の建玉を建てる日付(midPriceで購入)して、
一括決済に至るまでを一つのスパンとする.  

日付を進めながらシミュレーションを繰り返すことで、
包括的なシミュレーションを実行できる.  


### シミュレーションルール
- 指定した期間startで1lot購入
- 1円下がるごとに追加で1lot購入
- 最大10lotまで購入
- 利益が30%を超えたら売却
- 利益が-50%を超えたら売却
- 