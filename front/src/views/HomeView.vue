<template>
  <div class="container">
    <candle-chart v-bind:chartType="chartType" v-bind:chartData="chartData"></candle-chart>
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
      ['day', 'low', 'open', 'close', 'high'],
      ['Mon', 20, 28, 38, 45],
      ['Tue', 31, 38, 55, 66],
      ['Wed', 50, 55, 77, 80],
      ['Thu', 77, 77, 66, 50],
      ['Fri', 68, 66, 22, 15]
      ],
    }
  }
}
</script>
