# Go FX V2


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

data配下のcsvファイルは、時刻のデータが
`00:01`となっており、time.Parse時に`00:00:01`として解析されるため、
置換が必要.  
Vimを使用すると一瞬でできる.  
```vim
:%s/([0-9]{2}):([0-9]{2})/\1:\2:00
```
()にマッチしたものを、置換後の文字列で`\1`, `\2`などで引き出せる.  

[Vimでregexpでマッチした文字列を使用して置換](https://penguing27.hatenablog.jp/entry/2023/01/12/232452)



parseする際は下記のようにする.  
```go
// 24H単位の時間の場合
timeTime, err := time.Parse("2006-01-02 15:04:05", "2023-05-01 20:01:12")
```
`03:04:05`にした場合は、time.Hourのout of rangeのエラーが発生するので注意する.  


### SignalEventsの実装

抜き出したデータフレームにSignalEventsを渡す.  

個々のSignalEventは、データフレーム内で条件に合致した場合に
データベースに売買情報が保存される.  
なので、SignalEventsテーブルを参照することで、
売買ルールの有効性が確認できる.  



### シミュレーションの実行
- シミュレーションの開始時点を指定: start_date
- シミュレーションの条件を指定: rule.go
- 

### シミュレーションルール
- 指定した期間startで1lot購入
- 1円下がるごとに追加で1lot購入
- 最大10lotまで購入
- 利益が30%を超えたら売却
- 利益が-50%を超えたら売却
- 