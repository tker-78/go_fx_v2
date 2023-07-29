<template>
  <p>This is candle chart</p>
  <GChart
    :type="chartType"
    :data="chartData"
    :options="chartOptions"
  ></GChart>

</template>

<script>
  import { GChart } from 'vue-google-charts'

  export default {
    components: {
      GChart
    },
    props: {
      chartType: {
        type: String,
        default: '',
      },
      chartData: {
        type: Array,
        default: () => {
          return []
        }
      },
      chartOptions: {
        type: Object, 
        default: () => {
          return {}
        }
      }
    }
  }
</script>