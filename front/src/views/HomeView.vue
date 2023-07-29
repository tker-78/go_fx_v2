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
      chartType: 'PieChart',
      chartData: [
      ['年', '売上', '費用', '収益'],
      ['2014', 1000, 400, 200],
      ['2015', 1170, 460, 250],
      ['2016', 660, 1120, 300],
      ['2017', 1030, 540, 350]
      ],
    }
  }
}
</script>
