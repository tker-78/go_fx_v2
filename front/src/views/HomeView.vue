<template>
  <div class="container">
    <div>
      <label>Duration</label>
      <form>
        limit: <input type="text" v-model="limit" label="" style="width: 50px;">
        start: <input type="date" v-model="startDate" >
        end: <input type="date" v-model="endDate">
      </form>
    </div>

    <div class="inputarea">
      <label>SMA</label>
      <form>
        <input type="checkbox" id="sma" v-model="sma.enabled" @change="checkedSma">
        period1: <input type="text" v-model="sma.period1" style="width: 50px;"/>
        period2: <input type="text" v-model="sma.period2" style="width: 50px;"/>
        period3: <input type="text" v-model="sma.period3" style="width: 50px;"/> 
      </form>
    </div>

    <v-btn @click="onclick" variant="outlined">Draw candle chart</v-btn>

    <combo-chart v-bind:chartType="chartType" v-bind:chartData="candles" v-bind:chartOptions="chartOptions"></combo-chart>
  </div>
</template>

<script>
// @ is an alias to /src
// import CandleChart from '@/components/CandleChart.vue'
import ComboChart from '@/components/ComboChart.vue'

export default {
  name: 'HomeView',
  components: {
    ComboChart
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
          2: { type: 'line' }, // sma1
          3: { type: 'line' }, // sma2
          4: { type: 'line' }, // sma3
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
    }
  },
  methods: {
    onclick() {
      fetch(`http://localhost:8080/api/candle/?limit=${this.limit}&start=${this.startDate}&end=${this.endDate}&period1=${this.sma.period1}&period2=${this.sma.period2}&period3=${this.sma.period3}`)
      .then((response) => {
        return response.json()
      })
      .then((data) => {
        this.candles = []
        var header = ['time', 'low', 'open', 'close', 'high', 'swap']
        if (this.sma.enabled === true) {
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
          
        }else {
          this.candles.push(header)
          for (let candle of data.candles) {
            this.candles.push(
              [candle.time, parseFloat(candle.low), parseFloat(candle.open), parseFloat(candle.close), parseFloat(candle.high), parseFloat(candle.swap)]
            )
          }}
      })
    }, 
    checkedSma() {
      if ( this.sma.enabled === true ) {
        
        console.log("hello")
      }
    }
  }
}
</script>

<style scoped>
  .container {
    margin-top: 20px;
    height: 800px;
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
