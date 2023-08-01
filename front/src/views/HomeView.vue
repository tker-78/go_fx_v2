<template>
  <div class="container">
    <div>
      <label>Duration</label>
      <form>
        limit: <input type="text" v-model="limit" label="" style="width: 50px;">
        start: <input type="date" v-model="startDate">
        end: <input type="date" v-model="endDate">
      </form>
    </div>

    <div class="inputarea">
      <label>SMA</label>
      <form>
        <input type="checkbox" id="sma" v-model="sma.enabled">
        period1: <input type="text" v-model="sma.period1" style="width: 50px;"/>
        period2: <input type="text" v-model="sma.period2" style="width: 50px;"/>
        period3: <input type="text" v-model="sma.period3" style="width: 50px;"/> 
      </form>
    </div>

    <div class="inputarea">
      <label>EMA</label>
      <form>
        <input type="checkbox" id="ema" v-model="ema.enabled">
        period1: <input type="text" v-model="ema.period1" style="width: 50px;"/>
        period2: <input type="text" v-model="ema.period2" style="width: 50px;"/>
        period3: <input type="text" v-model="ema.period3" style="width: 50px;"/> 
      </form>
    </div>

    <div class="inputarea">
      <label>BBands</label>
      <form>
        <input type="checkbox" id="bbands" v-model="bband.enabled">
        N: <input type="text" v-model="bband.n" style="width: 50px;"/>
        K: <input type="text" v-model="bband.k" style="width: 50px;"/>
      </form>
    </div>

    <v-btn @click="onclick" variant="outlined">Draw candle chart</v-btn>

    <combo-chart v-bind:chartType="chartType" v-bind:chartData="candles" v-bind:chartOptions="chartOptions"></combo-chart>

    <signal-events :signals="signals[0]"></signal-events>
  </div>
</template>

<script>
// @ is an alias to /src
// import CandleChart from '@/components/CandleChart.vue'
import ComboChart from '@/components/ComboChart.vue'
import SignalEvents from '@/components/SignalEvents.vue'

export default {
  name: 'HomeView',
  components: {
    ComboChart,
    SignalEvents
  },
  data() {
    return {
      chartType: 'ComboChart',
      chartData: [],
      chartOptions: {
        title: "combo chart",
        seriesType: 'candlesticks',
        series: {
          1: { type: 'scatter', targetAxisIndex: 1 },
          2: { type: 'line' }, // sma1 or ema1
          3: { type: 'line' }, // sma2 or ema2
          4: { type: 'line' }, // sma3 or ema3
        },
        width: '100%',
        height: 800,
      },
      limit: 40,
      startDate: "",
      endDate: "",
      candles: [],
      sma: {
        period1: 7,
        period2: 14,
        period3: 25, 
        enabled: false,
        sma1: [],
        sma2: [],
        sma3: [],
      },
      ema: {
        period1: 7,
        period2: 14,
        period3: 25,
        enabled: false,
        ema1: [],
        ema2: [],
        ema3: [],
      }, 
      bband: {
        n: 20, 
        k: 2,
        enabled: false,
        bbup: [],
        bbmid: [],
        bbdown: [],
      },
      signals: [],
      signal: {
        time: "",
        currency_code: "",
        side: "",
        price: 0,
        size: 0,
      },
    }
  },
  methods: {
    onclick() {
      fetch(`http://localhost:8080/api/candle/?limit=${this.limit}&start=${this.startDate}&end=${this.endDate}&period1=${this.sma.period1}&period2=${this.sma.period2}&period3=${this.sma.period3}&bbn=${this.bband.n}&bbk=${this.bband.k}`)
      .then((response) => {
        return response.json()
      })
      .then((data) => {
        this.candles = []
        var header = ['time', 'low', 'open', 'close', 'high', 'swap']
        if (this.sma.enabled === true) {
          // SMAを有効にする場合
          header.push('sma1', 'sma2', 'sma3')
          this.candles.push(header)
          this.sma.sma1 = []
          this.sma.sma2 = []
          this.sma.sma3 = []

          if ( data.smas[0] != "undefined" ) {
            this.sma.sma1.push(data.smas[0])
          }

          if (data.smas[1] != "undefined") {
            this.sma.sma2.push(data.smas[1])
          }

          if (data.smas[2] != "undefined") {
            this.sma.sma3.push(data.smas[2])
          }

          for (let candle of data.candles) {
            this.candles.push(
              [
                candle.time, 
                parseFloat(candle.low), 
                parseFloat(candle.open), 
                parseFloat(candle.close), 
                parseFloat(candle.high), 
                parseFloat(candle.swap),
                0,
                0,
                0
              ]
            )
          }

          for (let i = 1; i < this.candles.length; i++) {
            this.candles[i][6] = parseFloat(this.sma.sma1[0].value[i-1]);
            this.candles[i][7] = parseFloat(this.sma.sma2[0].value[i-1]);
            this.candles[i][8] = parseFloat(this.sma.sma3[0].value[i-1]);
            if (this.candles[i][6] == 0) {
              this.candles[i][6] = null
            }

            if (this.candles[i][7] == 0) {
              this.candles[i][7] = null
            }

            if (this.candles[i][8] == 0) {
              this.candles[i][8] = null
            }
          }
          console.log(this.candles)
          
        } else if (this.ema.enabled === true) {
          // EMAを有効にする場合、
            header.push('ema1', 'ema2', 'ema3')
            this.candles.push(header)
            this.ema.ema1 = []
            this.ema.ema2 = []
            this.ema.ema3 = []

            if ( data.emas[0] != "undefined" ) {
              this.ema.ema1.push(data.emas[0])
            }

            if (data.emas[1] != "undefined") {
              this.ema.ema2.push(data.emas[1])
            }

            if (data.emas[2] != "undefined") {
              this.ema.ema3.push(data.emas[2])
            }

            for (let candle of data.candles) {
              this.candles.push(
                [
                  candle.time, 
                  parseFloat(candle.low), 
                  parseFloat(candle.open), 
                  parseFloat(candle.close), 
                  parseFloat(candle.high), 
                  parseFloat(candle.swap),
                  0,
                  0,
                  0
                ]
              )
            }

            for (let i = 1; i < this.candles.length; i++) {
              this.candles[i][6] = parseFloat(this.ema.ema1[0].value[i-1]);
              this.candles[i][7] = parseFloat(this.ema.ema2[0].value[i-1]);
              this.candles[i][8] = parseFloat(this.ema.ema3[0].value[i-1]);
              if (this.candles[i][6] == 0) {
                this.candles[i][6] = null
              }

              if (this.candles[i][7] == 0) {
                this.candles[i][7] = null
              }

              if (this.candles[i][8] == 0) {
                this.candles[i][8] = null
              }
            }
            console.log(this.candles)
            



        } else if ( this.bband.enabled === true ) {
          // BBandsを有効にする場合
            header.push('bbup', 'bbmid', 'bbdown')
            this.candles.push(header)
            this.bband.bbup =  []
            this.bband.bbmid = []
            this.bband.bbdown = []

            if ( data.bbands[0] != "undefined" ) {
              this.bband.bbup = data.bbands.up
            }

            if (data.bbands[1] != "undefined") {
              this.bband.bbmid = data.bbands.mid
            }

            if (data.bbands[2] != "undefined") {
              this.bband.bbdown = data.bbands.down
            }

            for (let candle of data.candles) {
              this.candles.push(
                [
                  candle.time, 
                  parseFloat(candle.low), 
                  parseFloat(candle.open), 
                  parseFloat(candle.close), 
                  parseFloat(candle.high), 
                  parseFloat(candle.swap),
                  0,
                  0,
                  0
                ]
              )
            }

            for (let i = 1; i < this.candles.length; i++) {
              this.candles[i][6] = parseFloat(this.bband.bbup[i-1]);
              this.candles[i][7] = parseFloat(this.bband.bbmid[i-1]);
              this.candles[i][8] = parseFloat(this.bband.bbdown[i-1]);
              if (this.candles[i][6] == 0) {
                this.candles[i][6] = null
              }

              if (this.candles[i][7] == 0) {
                this.candles[i][7] = null
              }

              if (this.candles[i][8] == 0) {
                this.candles[i][8] = null
              }
            }

            console.log(this.bband.bbup)



        } else {
          this.candles.push(header)
          for (let candle of data.candles) {
            this.candles.push(
              [candle.time, parseFloat(candle.low), parseFloat(candle.open), parseFloat(candle.close), parseFloat(candle.high), parseFloat(candle.swap)]
            )
          }}

        this.signals = []
        this.signals.push(data.signals.signals)
        console.log(this.signals[0])

      })
    }, 
  },
  deleteSignals() {
    // todo
  }
}
</script>

<style scoped>
  .container {
    margin-top: 20px;
    height: 800px;
    margin-bottom: 50px;
  }

  .inputarea {
    margin-top: 20px;
    margin-bottom: 20px;
    text-align: center;
  }
  input {
    border: 1px solid gray;
    border-radius: 4px;
  }


  
</style>
