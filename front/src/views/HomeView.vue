<template>
  <div class="container">
    <v-form>
      limit: <input type="text" v-model="limit" >
    </v-form>
    <v-btn @click="onclick">button</v-btn>

    <candle-chart v-bind:chartType="chartType" v-bind:chartData="candles" v-bind:chartOptions="chartOptions"></candle-chart>
  </div>
</template>

<script>
// @ is an alias to /src
import CandleChart from '@/components/CandleChart.vue'

export default {
  name: 'HomeView',
  components: {
    CandleChart,
  },
  data() {
    return {
      chartType: 'CandlestickChart',
      chartData: [
      ['time', 'low', 'open', 'close', 'high'],
      ['Mon', 20, 28, 38, 45],
      ['Tue', 31, 38, 55, 66],
      ['Wed', 50, 55, 77, 80],
      ['Thu', 77, 77, 66, 50],
      ['Fri', 68, 66, 22, 15]
      ],
      chartOptions: {
        width: '100%',
        height: 800,
      },
      limit: 10,
      candles: [],
    }
  },
  methods: {
    onclick() {
      fetch(`http://localhost:8080/api/candle/?limit=${this.limit}`)
      .then((response) => {
        console.log(response)
        return response.json()
      })
      .then((data) => {
        this.candles = []
        this.candles.push( ['time', 'low', 'open', 'close', 'high'])
        for (let candle of data.candles) {
          this.candles.push(
            [candle.time, parseFloat(candle.low), parseFloat(candle.open), parseFloat(candle.close), parseFloat(candle.high)]
          )
        }
        console.log(this.candles)
      })
    }
  }
}
</script>

<style scoped>
  .container {
    height: 800px;
  }
</style>
